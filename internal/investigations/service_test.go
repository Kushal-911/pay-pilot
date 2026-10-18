package investigations

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

const testOrderID = "123e4567-e89b-42d3-a456-426614174000"

type testStore struct {
	snapshot Snapshot
	err      error
}

func (s testStore) Load(context.Context, string) (Snapshot, error) { return s.snapshot, s.err }

type testModel struct {
	report Report
	err    error
}

func (m testModel) Analyze(context.Context, string, Snapshot, []Evidence) (Report, error) {
	return m.report, m.err
}

func sampleSnapshot() Snapshot {
	return Snapshot{Order: OrderRecord{ID: testOrderID, MerchantReference: "order-7", AmountMinor: 2500, Currency: "USD", Status: "pending"}, Attempts: []AttemptRecord{{ID: "223e4567-e89b-42d3-a456-426614174000", PayPalOrderID: stringPtr("PAYPAL-1"), AmountMinor: 2500, Currency: "USD", Status: "completed"}}}
}
func stringPtr(s string) *string { return &s }

func TestDeterministicReportCitesRealMismatchedRecords(t *testing.T) {
	snapshot := sampleSnapshot()
	evidence := BuildEvidence(snapshot)
	report := deterministicReport(snapshot, evidence)
	if report.Finding != "order_attempt_status_mismatch" {
		t.Fatalf("finding = %q", report.Finding)
	}
	if report.AnalysisSource != "deterministic" || report.Recommendation != "review_records" {
		t.Fatalf("unexpected report policy: %+v", report)
	}
	if err := validateModelReport(report, evidence); err != nil {
		t.Fatalf("deterministic report should have valid references: %v", err)
	}
}

func TestModelFailureReturnsDeterministicReport(t *testing.T) {
	snapshot := sampleSnapshot()
	service := Service{Store: testStore{snapshot: snapshot}, Model: testModel{err: errors.New("provider secret details should not reach response")}}
	report, err := service.Investigate(context.Background(), testOrderID, "Why pending?")
	if err != nil {
		t.Fatal(err)
	}
	if report.ModelStatus != "unavailable" || !strings.Contains(report.ModelMessage, "deterministic") {
		t.Fatalf("missing understandable fallback status: %+v", report)
	}
	if report.Finding != "order_attempt_status_mismatch" {
		t.Fatalf("fallback lost finding: %q", report.Finding)
	}
}

func TestInvalidModelEvidenceFallsBackAndDoesNotDowngradeDeterministicReview(t *testing.T) {
	snapshot := sampleSnapshot()
	modelReport := Report{Summary: "All records look good.", Explanation: Claim{Text: "No issue exists.", EvidenceIDs: []string{"order:" + testOrderID}}, ConfirmedFacts: []Claim{{Text: "The payment succeeded.", EvidenceIDs: []string{"attempt:invented"}}}, Confidence: .99, ConfidenceExplanation: "The model is certain.", Recommendation: "no_action"}
	service := Service{Store: testStore{snapshot: snapshot}, Model: testModel{report: modelReport}}
	report, err := service.Investigate(context.Background(), testOrderID, "Why pending?")
	if err != nil {
		t.Fatal(err)
	}
	if report.AnalysisSource != "deterministic" || report.ModelStatus != "unavailable" {
		t.Fatalf("invalid model output was accepted: %+v", report)
	}
	if report.Recommendation != "review_records" || report.Finding != "order_attempt_status_mismatch" {
		t.Fatalf("model downgraded deterministic finding: %+v", report)
	}
}

func TestEvidenceValidatorRejectsFabricatedReferences(t *testing.T) {
	evidence := BuildEvidence(sampleSnapshot())
	report := Report{Summary: "The payment records have a discrepancy.", Confidence: .8, Recommendation: "review_records", ConfirmedFacts: []Claim{{Text: "Unsupported fact", EvidenceIDs: []string{"event:invented"}}}}
	if err := validateModelReport(report, evidence); err == nil {
		t.Fatal("expected fabricated evidence to be rejected")
	}
}

func TestReadOnlyToolValidatesOrderScopeAndArguments(t *testing.T) {
	snapshot := sampleSnapshot()
	valid := json.RawMessage(`{"order_id":"` + testOrderID + `"}`)
	if _, err := executeTool("get_payment", valid, testOrderID, snapshot); err != nil {
		t.Fatalf("valid tool call rejected: %v", err)
	}
	wrongOrder := json.RawMessage(`{"order_id":"00000000-0000-0000-0000-000000000000"}`)
	if _, err := executeTool("get_payment", wrongOrder, testOrderID, snapshot); err == nil {
		t.Fatal("wrong order scope was accepted")
	}
	extra := json.RawMessage(`{"order_id":"` + testOrderID + `","sql":"select *"}`)
	if _, err := executeTool("get_payment", extra, testOrderID, snapshot); err == nil {
		t.Fatal("unexpected arguments were accepted")
	}
}

func TestDeterministicUnknownCaptureRequiresRefresh(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Order.Status = "pending"
	snapshot.Attempts[0].Status = "unknown"
	report := deterministicReport(snapshot, BuildEvidence(snapshot))
	if report.Recommendation != "refresh_provider_status" || report.RequiresHumanApproval {
		t.Fatalf("unsafe recommendation: %+v", report)
	}
}

func TestOpenAIModelRunsReadOnlyToolLoop(t *testing.T) {
	snapshot := sampleSnapshot()
	evidence := BuildEvidence(snapshot)
	toolOutput := make([]map[string]any, 0, 4)
	for _, name := range readOnlyToolNames {
		toolOutput = append(toolOutput, map[string]any{"type": "function_call", "call_id": "call-" + name, "name": name, "arguments": "{\"order_id\":\"" + testOrderID + "\"}", "id": "fc-" + name, "status": "completed"})
	}
	reportJSON, _ := json.Marshal(Report{Summary: "The order and payment attempt disagree.", Explanation: Claim{Text: "Order record is pending while its attempt is completed.", EvidenceIDs: []string{"order:" + testOrderID, "attempt:" + snapshot.Attempts[0].ID}}, ConfirmedFacts: []Claim{{Text: "Order is pending while its attempt is completed.", EvidenceIDs: []string{"order:" + testOrderID, "attempt:" + snapshot.Attempts[0].ID}}}, Hypotheses: []Claim{{Text: "A state update may have been delayed.", EvidenceIDs: []string{"order:" + testOrderID, "attempt:" + snapshot.Attempts[0].ID}}}, Confidence: .8, ConfidenceExplanation: "Both persisted records are available.", Recommendation: "review_records", UnresolvedQuestions: []string{}})
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/responses" || req.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("unexpected request %s %s", req.Method, req.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["store"] != false {
			t.Errorf("response storage must be disabled")
		}
		if _, ok := payload["tools"]; ok {
			return fakeResponse(http.StatusOK, map[string]any{"output": toolOutput}), nil
		}
		input, _ := json.Marshal(payload["input"])
		if !strings.Contains(string(input), "function_call_output") {
			t.Errorf("second request did not return tool outputs")
		}
		return fakeResponse(http.StatusOK, map[string]any{"output_text": string(reportJSON)}), nil
	})
	model := OpenAIModel{APIKey: "test-key", Model: "test-model", BaseURL: "https://api.openai.test/v1", HTTPClient: &http.Client{Transport: transport}}
	got, err := model.Analyze(context.Background(), "Why pending?", snapshot, evidence)
	if err != nil {
		t.Fatal(err)
	}
	if got.Summary != "The order and payment attempt disagree." {
		t.Fatalf("unexpected report: %+v", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func fakeResponse(status int, value any) *http.Response {
	body, _ := json.Marshal(value)
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(body))}
}
