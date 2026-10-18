package investigations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

var readOnlyToolNames = []string{"get_payment", "get_attempt_history", "list_provider_events", "compare_order_records"}

type toolArguments struct {
	OrderID string `json:"order_id"`
}

type toolOutput struct {
	Tool        string   `json:"tool"`
	OrderID     string   `json:"order_id"`
	Data        any      `json:"data"`
	EvidenceIDs []string `json:"evidence_ids"`
}

func executeTool(name string, raw json.RawMessage, expectedOrderID string, snapshot Snapshot) (toolOutput, error) {
	if !contains(readOnlyToolNames, name) {
		return toolOutput{}, fmt.Errorf("unsupported investigation tool")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var args toolArguments
	if err := decoder.Decode(&args); err != nil {
		return toolOutput{}, fmt.Errorf("invalid tool arguments")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return toolOutput{}, fmt.Errorf("invalid tool arguments")
	}
	if args.OrderID == "" || args.OrderID != expectedOrderID {
		return toolOutput{}, fmt.Errorf("tool order_id is outside this investigation")
	}
	out := toolOutput{Tool: name, OrderID: expectedOrderID, EvidenceIDs: make([]string, 0)}
	switch name {
	case "get_payment":
		out.Data = snapshot.Order
		out.EvidenceIDs = append(out.EvidenceIDs, "order:"+snapshot.Order.ID)
		if len(snapshot.Attempts) > 0 {
			latest := snapshot.Attempts[len(snapshot.Attempts)-1]
			out.Data = map[string]any{"order": snapshot.Order, "latest_attempt": latest}
			out.EvidenceIDs = append(out.EvidenceIDs, "attempt:"+latest.ID)
		}
	case "get_attempt_history":
		out.Data = snapshot.Attempts
		for _, attempt := range snapshot.Attempts {
			out.EvidenceIDs = append(out.EvidenceIDs, "attempt:"+attempt.ID)
		}
	case "list_provider_events":
		out.Data = snapshot.Events
		for _, event := range snapshot.Events {
			out.EvidenceIDs = append(out.EvidenceIDs, "event:"+event.ID)
		}
	case "compare_order_records":
		out.Data = compare(snapshot)
		out.EvidenceIDs = append(out.EvidenceIDs, "order:"+snapshot.Order.ID)
		for _, attempt := range snapshot.Attempts {
			out.EvidenceIDs = append(out.EvidenceIDs, "attempt:"+attempt.ID)
		}
		for _, event := range snapshot.Events {
			out.EvidenceIDs = append(out.EvidenceIDs, "event:"+event.ID)
		}
	}
	return out, nil
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

type Comparison struct {
	Code        string   `json:"code"`
	Message     string   `json:"message"`
	EvidenceIDs []string `json:"evidence_ids"`
}

func compare(snapshot Snapshot) []Comparison {
	findings := make([]Comparison, 0)
	orderID := "order:" + snapshot.Order.ID
	hasCompleted, hasUnknown := false, false
	for _, attempt := range snapshot.Attempts {
		attemptID := "attempt:" + attempt.ID
		if attempt.AmountMinor != snapshot.Order.AmountMinor || attempt.Currency != snapshot.Order.Currency {
			findings = append(findings, Comparison{Code: "amount_currency_mismatch", Message: "Payment attempt amount or currency differs from its merchant order", EvidenceIDs: []string{orderID, attemptID}})
		}
		if attempt.Status == "completed" {
			hasCompleted = true
		}
		if attempt.Status == "unknown" || attempt.Status == "capturing" {
			hasUnknown = true
		}
	}
	if snapshot.Order.Status == "pending" && hasCompleted {
		findings = append(findings, Comparison{Code: "order_attempt_status_mismatch", Message: "The order is pending while at least one payment attempt is completed", EvidenceIDs: allAttemptIDs(snapshot)})
	}
	if snapshot.Order.Status == "paid" && !hasCompleted {
		findings = append(findings, Comparison{Code: "order_attempt_status_mismatch", Message: "The order is paid but no completed payment attempt is recorded", EvidenceIDs: allAttemptIDs(snapshot)})
	}
	if hasUnknown {
		findings = append(findings, Comparison{Code: "capture_outcome_unknown", Message: "At least one payment attempt has an unknown or in-progress capture outcome", EvidenceIDs: allAttemptIDs(snapshot)})
	}
	for _, event := range snapshot.Events {
		if event.EventType == "PAYMENT.CAPTURE.COMPLETED" {
			matchingCompletedAttempt := false
			for _, attempt := range snapshot.Attempts {
				matchesOrder := event.PayPalOrderID != nil && attempt.PayPalOrderID != nil && *event.PayPalOrderID == *attempt.PayPalOrderID
				matchesCapture := event.PayPalCaptureID != nil && attempt.PayPalCaptureID != nil && *event.PayPalCaptureID == *attempt.PayPalCaptureID
				if (matchesOrder || matchesCapture) && attempt.Status == "completed" && (event.PayPalCaptureID == nil || matchesCapture) {
					matchingCompletedAttempt = true
				}
			}
			if !matchingCompletedAttempt {
				findings = append(findings, Comparison{Code: "event_attempt_mismatch", Message: "A persisted capture-completed event has no completed payment attempt", EvidenceIDs: append(allAttemptIDs(snapshot), "event:"+event.ID)})
			}
		}
	}
	return findings
}

func allAttemptIDs(snapshot Snapshot) []string {
	ids := []string{"order:" + snapshot.Order.ID}
	for _, attempt := range snapshot.Attempts {
		ids = append(ids, "attempt:"+attempt.ID)
	}
	return ids
}

func toolDefinitions() []map[string]any {
	defs := make([]map[string]any, 0, len(readOnlyToolNames))
	descriptions := map[string]string{
		"get_payment":           "Retrieve the merchant order and latest payment attempt for the order under investigation.",
		"get_attempt_history":   "Retrieve all payment attempts for the order under investigation.",
		"list_provider_events":  "Retrieve stored provider event metadata associated with the order's PayPal identifiers.",
		"compare_order_records": "Compare merchant order fields, payment attempts, and provider event evidence using deterministic rules.",
	}
	for _, name := range readOnlyToolNames {
		defs = append(defs, map[string]any{"type": "function", "name": name, "description": descriptions[name], "parameters": map[string]any{"type": "object", "properties": map[string]any{"order_id": map[string]any{"type": "string", "description": "The UUID of the order currently being investigated."}}, "required": []string{"order_id"}, "additionalProperties": false}, "strict": true})
	}
	return defs
}
