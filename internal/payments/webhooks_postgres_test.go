package payments

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"paypilot/internal/paypal"
)

type webhookRoundTripper func(*http.Request) (*http.Response, error)

func (f webhookRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPayPalWebhookPersistsDeduplicatesAndKeepsCompletedState(t *testing.T) {
	db := newPostgresTestSchema(t)
	ctx := context.Background()
	var orderID string
	err := db.QueryRow(ctx, `INSERT INTO orders(merchant_reference,amount_minor,currency,description,status) VALUES('webhook-integration',2500,'USD','webhook integration','pending') RETURNING id`).Scan(&orderID)
	if err != nil {
		t.Fatal(err)
	}
	const paypalOrderID = "PAYPAL-ORDER-INTEGRATION"
	if _, err := db.Exec(ctx, `INSERT INTO payment_attempts(order_id,paypal_order_id,amount_minor,currency,status) VALUES($1,$2,2500,'USD','unknown')`, orderID, paypalOrderID); err != nil {
		t.Fatal(err)
	}
	verification := "SUCCESS"
	client := paypal.New("https://paypal.test", "client", "secret", &http.Client{Transport: webhookRoundTripper(func(r *http.Request) (*http.Response, error) {
		body := `{"access_token":"token"}`
		if strings.HasSuffix(r.URL.Path, "/verify-webhook-signature") {
			body = fmt.Sprintf(`{"verification_status":%q}`, verification)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})})
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/webhook", PayPalWebhook(db, client, "WEBHOOK-ID"))
	post := func(eventID, eventType, status, captureID, amount string) *httptest.ResponseRecorder {
		body := fmt.Sprintf(`{"id":%q,"event_type":%q,"resource":{"id":%q,"status":%q,"amount":{"value":%q,"currency_code":"USD"},"supplementary_data":{"related_ids":{"order_id":%q}}}}`, eventID, eventType, captureID, status, amount, paypalOrderID)
		request := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
		request.Header.Set("paypal-auth-algo", "SHA256withRSA")
		request.Header.Set("paypal-cert-url", "https://api-m.paypal.com/cert")
		request.Header.Set("paypal-transmission-id", "transmission")
		request.Header.Set("paypal-transmission-sig", "signature")
		request.Header.Set("paypal-transmission-time", "2026-11-05T10:00:00Z")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		return response
	}
	assertAccepted := func(response *httptest.ResponseRecorder, eventID string) {
		t.Helper()
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), eventID) {
			t.Fatalf("webhook response = %d %s", response.Code, response.Body)
		}
	}

	completed := post("EVENT-COMPLETED", "PAYMENT.CAPTURE.COMPLETED", "COMPLETED", "CAPTURE-1", "25.00")
	assertAccepted(completed, "EVENT-COMPLETED")
	var attemptStatus, orderStatus, captureID string
	if err := db.QueryRow(ctx, `SELECT p.status,o.status,p.paypal_capture_id FROM payment_attempts p JOIN orders o ON o.id=p.order_id WHERE p.order_id=$1`, orderID).Scan(&attemptStatus, &orderStatus, &captureID); err != nil {
		t.Fatal(err)
	}
	if attemptStatus != "completed" || orderStatus != "paid" || captureID != "CAPTURE-1" {
		t.Fatalf("completion state = attempt %q, order %q, capture %q", attemptStatus, orderStatus, captureID)
	}
	assertAccepted(post("EVENT-COMPLETED", "PAYMENT.CAPTURE.COMPLETED", "COMPLETED", "CAPTURE-1", "25.00"), "EVENT-COMPLETED")
	var duplicateCount int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM provider_events WHERE provider_event_id='EVENT-COMPLETED'`).Scan(&duplicateCount); err != nil || duplicateCount != 1 {
		t.Fatalf("duplicate event rows = %d, err = %v", duplicateCount, err)
	}
	assertAccepted(post("EVENT-LATE-PENDING", "PAYMENT.CAPTURE.PENDING", "PENDING", "CAPTURE-1", "25.00"), "EVENT-LATE-PENDING")
	if err := db.QueryRow(ctx, `SELECT p.status,o.status FROM payment_attempts p JOIN orders o ON o.id=p.order_id WHERE p.order_id=$1`, orderID).Scan(&attemptStatus, &orderStatus); err != nil {
		t.Fatal(err)
	}
	if attemptStatus != "completed" || orderStatus != "paid" {
		t.Fatalf("late pending event downgraded state: attempt %q, order %q", attemptStatus, orderStatus)
	}
	assertAccepted(post("EVENT-WRONG-AMOUNT", "PAYMENT.CAPTURE.COMPLETED", "COMPLETED", "CAPTURE-2", "24.99"), "EVENT-WRONG-AMOUNT")
	var processingStatus string
	if err := db.QueryRow(ctx, `SELECT processing_status FROM provider_events WHERE provider_event_id='EVENT-WRONG-AMOUNT'`).Scan(&processingStatus); err != nil || processingStatus != "failed" {
		t.Fatalf("mismatched event processing status = %q, err = %v", processingStatus, err)
	}
	if err := db.QueryRow(ctx, `SELECT p.status,o.status FROM payment_attempts p JOIN orders o ON o.id=p.order_id WHERE p.order_id=$1`, orderID).Scan(&attemptStatus, &orderStatus); err != nil {
		t.Fatal(err)
	}
	if attemptStatus != "completed" || orderStatus != "paid" {
		t.Fatalf("mismatched event changed paid state: attempt %q, order %q", attemptStatus, orderStatus)
	}

	verification = "FAILURE"
	invalid := post("EVENT-INVALID-SIG", "PAYMENT.CAPTURE.COMPLETED", "COMPLETED", "CAPTURE-3", "25.00")
	if invalid.Code != http.StatusUnauthorized {
		t.Fatalf("invalid signature status = %d, body = %s", invalid.Code, invalid.Body)
	}
	if err := db.QueryRow(ctx, `SELECT count(*) FROM provider_events WHERE provider_event_id='EVENT-INVALID-SIG'`).Scan(&duplicateCount); err != nil || duplicateCount != 0 {
		t.Fatalf("invalid-signature event persisted count = %d, err = %v", duplicateCount, err)
	}
}
