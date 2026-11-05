package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func TestCreateOrderUsesOAuthAndMinorUnits(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/v1/oauth2/token":
			if r.Method != http.MethodPost {
				t.Errorf("token method = %s", r.Method)
			}
			return response(http.StatusOK, `{"access_token":"test-token"}`), nil
		case "/v2/checkout/orders":
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("missing bearer token")
			}
			if got := r.Header.Get("PayPal-Request-Id"); got != "stable-order-request" {
				t.Errorf("PayPal-Request-Id = %q", got)
			}
			var body struct {
				PurchaseUnits []struct {
					Amount struct {
						Value string `json:"value"`
					} `json:"amount"`
				} `json:"purchase_units"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if got := body.PurchaseUnits[0].Amount.Value; got != "25.00" {
				t.Errorf("amount = %q, want 25.00", got)
			}
			return response(http.StatusCreated, `{"id":"ORDER-1","status":"CREATED","links":[{"rel":"approve","href":"https://paypal.test/approve"}]}`), nil
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			return response(http.StatusNotFound, ""), nil
		}
	})}
	api := New("https://paypal.test", "client", "secret", client)
	result, err := api.CreateOrder(context.Background(), "ref-1", "Test", "USD", 2500, "https://merchant.test/return", "https://merchant.test/cancel", "stable-order-request")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "ORDER-1" || result.Status != "CREATED" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestGetOrderRefreshesProviderState(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/v1/oauth2/token" {
			return response(http.StatusOK, `{"access_token":"test-token"}`), nil
		}
		if r.Method != http.MethodGet || r.URL.Path != "/v2/checkout/orders/ORDER-2" {
			t.Fatalf("unexpected refresh request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("PayPal-Request-Id"); got != "" {
			t.Errorf("GET refresh unexpectedly has request ID %q", got)
		}
		return response(http.StatusOK, `{"id":"ORDER-2","status":"COMPLETED","purchase_units":[{"amount":{"currency_code":"USD","value":"25.00"},"payments":{"captures":[{"id":"CAPTURE-2","status":"COMPLETED","amount":{"currency_code":"USD","value":"25.00"}}]}}]}`), nil
	})}
	api := New("https://paypal.test", "client", "secret", client)
	result, err := api.GetOrder(context.Background(), "ORDER-2")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "COMPLETED" || result.PurchaseUnits[0].Payments.Captures[0].ID != "CAPTURE-2" {
		t.Fatalf("unexpected refreshed state: %+v", result)
	}
}

func TestCaptureUsesSuppliedIdempotencyRequestID(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/v1/oauth2/token" {
			return response(http.StatusOK, `{"access_token":"test-token"}`), nil
		}
		if r.Method != http.MethodPost || r.URL.Path != "/v2/checkout/orders/ORDER-3/capture" {
			t.Fatalf("unexpected capture request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("PayPal-Request-Id"); got != "stable-capture-request" {
			t.Errorf("PayPal-Request-Id = %q", got)
		}
		return response(http.StatusCreated, `{"id":"ORDER-3","status":"COMPLETED","purchase_units":[{"payments":{"captures":[{"id":"CAPTURE-3","status":"COMPLETED","amount":{"currency_code":"USD","value":"25.00"}}]}}]}`), nil
	})}
	api := New("https://paypal.test", "client", "secret", client)
	result, err := api.Capture(context.Background(), "ORDER-3", "stable-capture-request")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "COMPLETED" || result.PurchaseUnits[0].Payments.Captures[0].Amount.Value != "25.00" {
		t.Fatalf("capture status = %q", result.Status)
	}
}

func TestDoRejectsOversizedAndMalformedProviderResponses(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
	}{
		{name: "oversized", body: strings.Repeat("x", maxProviderResponseBytes+1), want: "exceeds"},
		{name: "malformed", body: "not-json", want: "decode PayPal response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/v1/oauth2/token" {
					return response(http.StatusOK, `{"access_token":"test-token"}`), nil
				}
				return response(http.StatusOK, tc.body), nil
			})}
			api := New("https://paypal.test", "client", "secret", client)
			_, err := api.GetOrder(context.Background(), "ORDER-4")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func TestDoHonorsRequestContextDeadline(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done()
		return nil, r.Context().Err()
	})}
	api := New("https://paypal.test", "client", "secret", client)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	if _, err := api.GetOrder(ctx, "ORDER-5"); err == nil {
		t.Fatal("expected request deadline error")
	}
}

func TestVerifyWebhookPostsRawEventAndHeaderMetadata(t *testing.T) {
	rawEvent := json.RawMessage(`{ "id":"WH-EVENT-1", "event_type":"PAYMENT.CAPTURE.COMPLETED" }`)
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/v1/oauth2/token" {
			return response(http.StatusOK, `{"access_token":"test-token"}`), nil
		}
		if r.URL.Path != "/v1/notifications/verify-webhook-signature" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		var body struct {
			AuthAlgo  string          `json:"auth_algo"`
			WebhookID string          `json:"webhook_id"`
			Event     json.RawMessage `json:"webhook_event"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.AuthAlgo != "SHA256withRSA" || body.WebhookID != "WEBHOOK-123" {
			t.Errorf("missing verification values: %+v", body)
		}
		if !bytes.Equal(body.Event, rawEvent) {
			t.Errorf("raw webhook event changed: %s", body.Event)
		}
		return response(http.StatusOK, `{"verification_status":"SUCCESS"}`), nil
	})}
	api := New("https://paypal.test", "client", "secret", client)
	valid, err := api.VerifyWebhook(context.Background(), "WEBHOOK-123", WebhookHeaders{AuthAlgo: "SHA256withRSA", CertURL: "https://api.paypal.test/cert", TransmissionID: "t-1", TransmissionSig: "sig", TransmissionTime: "2026-10-18T00:00:00Z"}, rawEvent)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("expected verified webhook")
	}
}

func TestProviderFailureDoesNotLeakBody(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return response(http.StatusUnauthorized, "sensitive provider detail"), nil
	})
	client := New("https://paypal.test", "client", "secret", &http.Client{Transport: transport})
	_, err := client.CreateOrder(context.Background(), "r", "d", "USD", 100, "r", "c", "stable-order-request")
	if err == nil {
		t.Fatal("expected OAuth error")
	}
	if got := err.Error(); got == "" || strings.Contains(got, "sensitive provider detail") {
		t.Fatalf("unsafe or empty error: %q", got)
	}
}
