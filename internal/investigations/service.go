package investigations

import (
	"context"
	"fmt"
	"math"
)

type Service struct {
	Store Store
	Model Model
}

func (s Service) Investigate(ctx context.Context, orderID, question string) (Report, error) {
	snapshot, err := s.Store.Load(ctx, orderID)
	if err != nil {
		return Report{}, err
	}
	evidence := BuildEvidence(snapshot)
	fallback := deterministicReport(snapshot, evidence)
	if s.Model == nil {
		fallback.ModelStatus = "disabled"
		return fallback, nil
	}
	modelReport, err := s.Model.Analyze(ctx, question, snapshot, evidence)
	if err != nil {
		fallback.ModelStatus = "unavailable"
		fallback.ModelMessage = "AI analysis was unavailable; this report uses deterministic record checks."
		return fallback, nil
	}
	if err := validateModelReport(modelReport, evidence); err != nil {
		fallback.ModelStatus = "unavailable"
		fallback.ModelMessage = "AI analysis did not pass evidence validation; this report uses deterministic record checks."
		return fallback, nil
	}
	modelReport.Finding = fallback.Finding
	modelReport.Summary = fallback.Summary
	modelReport.ConfirmedFacts = fallback.ConfirmedFacts
	modelReport.Evidence = evidence
	modelReport.AnalysisSource = "openai"
	modelReport.ModelStatus = "used"
	// Deterministic payment-state checks take precedence over model suggestions.
	if fallback.Recommendation != "no_action" {
		modelReport.Recommendation = fallback.Recommendation
	} else if !allowedRecommendation(modelReport.Recommendation) {
		modelReport.Recommendation = fallback.Recommendation
	}
	if modelReport.Recommendation == "manual_review" {
		modelReport.RequiresHumanApproval = true
	} else {
		modelReport.RequiresHumanApproval = false
	}
	modelReport.Evidence = evidence
	return modelReport, nil
}

func deterministicReport(snapshot Snapshot, evidence []Evidence) Report {
	findings := compare(snapshot)
	report := Report{Finding: "records_consistent", Summary: "No discrepancy was detected across the available order, attempt, and provider-event records.", Confidence: 0.88, ConfidenceExplanation: "The merchant order and available payment attempts agree under deterministic checks. Provider-event coverage may be incomplete.", Recommendation: "no_action", UnresolvedQuestions: make([]string, 0), Evidence: evidence, AnalysisSource: "deterministic"}
	orderEvidence := "order:" + snapshot.Order.ID
	report.ConfirmedFacts = append(report.ConfirmedFacts, Claim{Text: fmt.Sprintf("Merchant order %s is %s for %d %s minor units.", snapshot.Order.MerchantReference, snapshot.Order.Status, snapshot.Order.AmountMinor, snapshot.Order.Currency), EvidenceIDs: []string{orderEvidence}})
	for _, attempt := range snapshot.Attempts {
		report.ConfirmedFacts = append(report.ConfirmedFacts, Claim{Text: fmt.Sprintf("Payment attempt %s is %s for %d %s minor units.", attempt.ID, attempt.Status, attempt.AmountMinor, attempt.Currency), EvidenceIDs: []string{"attempt:" + attempt.ID}})
	}
	for _, event := range snapshot.Events {
		report.ConfirmedFacts = append(report.ConfirmedFacts, Claim{Text: fmt.Sprintf("Provider event %s (%s) was received with processing status %s.", event.ProviderEventID, event.EventType, event.ProcessingStatus), EvidenceIDs: []string{"event:" + event.ID}})
	}
	if len(snapshot.Events) == 0 {
		report.UnresolvedQuestions = append(report.UnresolvedQuestions, "No matching provider events are stored, so webhook delivery and event ordering cannot be assessed.")
	}
	if len(findings) == 0 {
		report.Explanation = Claim{Text: report.Summary, EvidenceIDs: allAttemptIDs(snapshot)}
		return report
	}
	if len(findings) == 1 {
		report.Finding = findings[0].Code
	} else {
		report.Finding = "multiple_discrepancies"
	}
	report.Summary = findings[0].Message
	if len(findings) > 1 {
		report.Summary = fmt.Sprintf("%s Additional discrepancy checks also found %d issue(s).", findings[0].Message, len(findings)-1)
	}
	report.Explanation = Claim{Text: findings[0].Message, EvidenceIDs: findings[0].EvidenceIDs}
	report.Confidence = 0.96
	report.ConfidenceExplanation = "The discrepancy is derived from directly compared persisted records; the cause of the mismatch may still be uncertain."
	report.Recommendation = "review_records"
	if len(report.Hypotheses) == 0 {
		report.Hypotheses = make([]Claim, 0)
	}
	hasAmountMismatch, hasUnknown := false, false
	for _, finding := range findings {
		report.ConfirmedFacts = append(report.ConfirmedFacts, Claim{Text: finding.Message, EvidenceIDs: finding.EvidenceIDs})
		if finding.Code == "capture_outcome_unknown" {
			hasUnknown = true
			report.Hypotheses = append(report.Hypotheses, Claim{Text: "The capture request may have succeeded remotely without the application recording a confirmed response.", EvidenceIDs: finding.EvidenceIDs})
			report.UnresolvedQuestions = append(report.UnresolvedQuestions, "What status does PayPal currently report for this order and its captures?")
		}
		if finding.Code == "order_attempt_status_mismatch" {
			report.Hypotheses = append(report.Hypotheses, Claim{Text: "A state update may have been delayed or lost between PayPal and the merchant database.", EvidenceIDs: finding.EvidenceIDs})
		}
		if finding.Code == "amount_currency_mismatch" {
			hasAmountMismatch = true
		}
	}
	if hasAmountMismatch {
		report.Recommendation = "manual_review"
		report.RequiresHumanApproval = true
	} else if hasUnknown {
		report.Recommendation = "refresh_provider_status"
		report.RequiresHumanApproval = false
	}
	return report
}

func validateModelReport(r Report, evidence []Evidence) error {
	if len(r.Summary) < 8 || len(r.Summary) > 2000 || math.IsNaN(r.Confidence) || math.IsInf(r.Confidence, 0) || r.Confidence < 0 || r.Confidence > 1 {
		return fmt.Errorf("invalid report fields")
	}
	if !allowedRecommendation(r.Recommendation) || len(r.ConfirmedFacts) == 0 || len(r.ConfirmedFacts) > 20 || len(r.Hypotheses) > 20 || len(r.UnresolvedQuestions) > 20 {
		return fmt.Errorf("invalid report shape")
	}
	valid := make(map[string]bool, len(evidence))
	for _, e := range evidence {
		valid[e.ID] = true
	}
	for _, claim := range append(append([]Claim{}, r.ConfirmedFacts...), r.Hypotheses...) {
		if len(claim.Text) == 0 || len(claim.Text) > 1000 || len(claim.EvidenceIDs) == 0 {
			return fmt.Errorf("claim missing evidence")
		}
		for _, id := range claim.EvidenceIDs {
			if !valid[id] {
				return fmt.Errorf("claim cites unknown evidence")
			}
		}
	}
	if len(r.Explanation.Text) < 8 || len(r.Explanation.Text) > 2000 || len(r.Explanation.EvidenceIDs) == 0 {
		return fmt.Errorf("explanation missing evidence")
	}
	for _, id := range r.Explanation.EvidenceIDs {
		if !valid[id] {
			return fmt.Errorf("explanation cites unknown evidence")
		}
	}
	for _, q := range r.UnresolvedQuestions {
		if len(q) == 0 || len(q) > 500 {
			return fmt.Errorf("invalid unresolved question")
		}
	}
	return nil
}

func allowedRecommendation(value string) bool {
	switch value {
	case "no_action", "refresh_provider_status", "review_records", "manual_review":
		return true
	}
	return false
}
