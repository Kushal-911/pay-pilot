package payments

import (
	"testing"
)

func TestParseMinorAmount(t *testing.T) {
	cases := []struct {
		value string
		want  int64
		valid bool
	}{{"25.00", 2500, true}, {"25.5", 2550, true}, {"0.01", 1, true}, {"25", 2500, true}, {"25.000", 0, false}, {"-1.00", 0, false}, {"92233720368547758.08", 0, false}}
	for _, tc := range cases {
		got, err := parseMinorAmount(tc.value)
		if tc.valid && (err != nil || got != tc.want) {
			t.Errorf("parseMinorAmount(%q) = %d,%v; want %d", tc.value, got, err, tc.want)
		}
		if !tc.valid && err == nil {
			t.Errorf("parseMinorAmount(%q) unexpectedly succeeded: %d", tc.value, got)
		}
	}
}

func TestCaptureWebhookTransitionsAreMonotonic(t *testing.T) {
	captureID := "CAPTURE-1"
	base := webhookEvent{ID: "EVENT-1", EventType: "PAYMENT.CAPTURE.COMPLETED"}
	base.Resource.ID = "CAPTURE-1"
	base.Resource.Status = "COMPLETED"
	base.Resource.Amount.CurrencyCode = "USD"
	current, processing, apply := classifyCaptureWebhook(base, "failed", nil, 2500, "USD", 2500)
	if current != "completed" || processing != "processed" || !apply {
		t.Fatalf("verified completion did not recover failed local state: %s %s %t", current, processing, apply)
	}

	pending := base
	pending.EventType = "PAYMENT.CAPTURE.PENDING"
	pending.Resource.Status = "PENDING"
	current, processing, apply = classifyCaptureWebhook(pending, "completed", &captureID, 2500, "USD", 2500)
	if current != "completed" || processing != "processed" || apply {
		t.Fatalf("out-of-order pending event downgraded completed attempt: %s %s %t", current, processing, apply)
	}

	denied := base
	denied.EventType = "PAYMENT.CAPTURE.DENIED"
	denied.Resource.Status = "DENIED"
	current, processing, apply = classifyCaptureWebhook(denied, "completed", &captureID, 2500, "USD", 2500)
	if current != "completed" || processing != "processed" || apply {
		t.Fatalf("out-of-order denied event downgraded completed attempt: %s %s %t", current, processing, apply)
	}

	current, processing, apply = classifyCaptureWebhook(base, "unknown", nil, 2600, "USD", 2500)
	if current != "unknown" || processing != "failed" || apply {
		t.Fatalf("amount mismatch was applied: %s %s %t", current, processing, apply)
	}
}
