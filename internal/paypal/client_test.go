package paypal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
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
	result, err := api.CreateOrder(context.Background(), "ref-1", "Test", "USD", 2500, "https://merchant.test/return", "https://merchant.test/cancel")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "ORDER-1" || result.Status != "CREATED" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestProviderFailureDoesNotLeakBody(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return response(http.StatusUnauthorized, "sensitive provider detail"), nil
	})
	client := New("https://paypal.test", "client", "secret", &http.Client{Transport: transport})
	_, err := client.CreateOrder(context.Background(), "r", "d", "USD", 100, "r", "c")
	if err == nil {
		t.Fatal("expected OAuth error")
	}
	if got := err.Error(); got == "" || strings.Contains(got, "sensitive provider detail") {
		t.Fatalf("unsafe or empty error: %q", got)
	}
}
