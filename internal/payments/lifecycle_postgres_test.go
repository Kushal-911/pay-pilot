package payments

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"paypilot/internal/investigations"
	"paypilot/internal/paypal"
)

func TestOrderLifecyclePersistsAcrossPostgresHandlers(t *testing.T) {
	db := newPostgresTestSchema(t)
	t.Setenv("PAYPAL_RETURN_URL", "http://localhost/return")
	t.Setenv("PAYPAL_CANCEL_URL", "http://localhost/cancel")
	paypalClient := paypal.New("https://paypal.test", "client", "secret", &http.Client{Transport: webhookRoundTripper(func(r *http.Request) (*http.Response, error) {
		body := `{"access_token":"token"}`
		switch {
		case strings.HasSuffix(r.URL.Path, "/v1/oauth2/token"):
			// Return the OAuth token response above.
		case strings.HasSuffix(r.URL.Path, "/v2/checkout/orders") && r.Method == http.MethodPost:
			body = `{"id":"PAYPAL-LIFECYCLE-ORDER","status":"APPROVED","links":[{"rel":"approve","href":"https://paypal.test/approve"}]}`
		case strings.HasSuffix(r.URL.Path, "/capture"):
			body = `{"id":"PAYPAL-LIFECYCLE-ORDER","status":"COMPLETED","purchase_units":[{"amount":{"currency_code":"USD","value":"25.00"},"payments":{"captures":[{"id":"PAYPAL-LIFECYCLE-CAPTURE","status":"COMPLETED","amount":{"currency_code":"USD","value":"25.00"}}]}}]}`
		case strings.HasSuffix(r.URL.Path, "/v2/checkout/orders/PAYPAL-LIFECYCLE-ORDER"):
			body = `{"id":"PAYPAL-LIFECYCLE-ORDER","status":"COMPLETED","purchase_units":[{"amount":{"currency_code":"USD","value":"25.00"},"payments":{"captures":[{"id":"PAYPAL-LIFECYCLE-CAPTURE","status":"COMPLETED","amount":{"currency_code":"USD","value":"25.00"}}]}}]}`
		default:
			return nil, io.ErrUnexpectedEOF
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})})
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	api.POST("/orders", Idempotency(db), CreateOrder(db))
	api.GET("/orders/:id", GetOrder(db))
	api.POST("/orders/:id/paypal-order", Idempotency(db), CreatePayPalOrder(db, paypalClient))
	api.POST("/payments/:id/capture", Idempotency(db), CapturePayPalOrder(db, paypalClient))
	api.GET("/payments/:id/refresh", RefreshPayPalOrder(db, paypalClient))
	api.POST("/investigations", Idempotency(db), investigations.CreateHandler(db, nil))
	api.GET("/investigations/:id", investigations.GetHandler(db))
	api.GET("/orders/:id/investigations", investigations.ListForOrderHandler(db))
	request := func(method, path, body, key string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		return res
	}
	decode := func(response *httptest.ResponseRecorder, target any) {
		t.Helper()
		if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
			t.Fatalf("decode response %q: %v", response.Body.String(), err)
		}
	}
	createdOrder := request(http.MethodPost, "/api/v1/orders", `{"merchant_reference":"lifecycle-001","amount_minor":2500,"currency":"USD","description":"Lifecycle integration"}`, "lifecycle-create-order")
	if createdOrder.Code != http.StatusCreated {
		t.Fatalf("create order = %d %s", createdOrder.Code, createdOrder.Body)
	}
	var order struct {
		ID string `json:"id"`
	}
	decode(createdOrder, &order)
	persistedOrder := request(http.MethodGet, "/api/v1/orders/"+order.ID, "", "")
	if persistedOrder.Code != http.StatusOK || !strings.Contains(persistedOrder.Body.String(), `"status":"pending"`) {
		t.Fatalf("read persisted order = %d %s", persistedOrder.Code, persistedOrder.Body)
	}
	createdAttempt := request(http.MethodPost, "/api/v1/orders/"+order.ID+"/paypal-order", `{}`, "lifecycle-paypal-order")
	if createdAttempt.Code != http.StatusCreated {
		t.Fatalf("create provider order = %d %s", createdAttempt.Code, createdAttempt.Body)
	}
	var attempt struct {
		ID string `json:"payment_id"`
	}
	decode(createdAttempt, &attempt)
	createdCapture := request(http.MethodPost, "/api/v1/payments/"+attempt.ID+"/capture", `{}`, "lifecycle-capture")
	if createdCapture.Code != http.StatusOK || !strings.Contains(createdCapture.Body.String(), `"status":"completed"`) {
		t.Fatalf("capture = %d %s", createdCapture.Code, createdCapture.Body)
	}
	refreshed := request(http.MethodGet, "/api/v1/payments/"+attempt.ID+"/refresh", "", "")
	if refreshed.Code != http.StatusOK || !strings.Contains(refreshed.Body.String(), `"status":"completed"`) {
		t.Fatalf("refresh = %d %s", refreshed.Code, refreshed.Body)
	}
	paidOrder := request(http.MethodGet, "/api/v1/orders/"+order.ID, "", "")
	if paidOrder.Code != http.StatusOK || !strings.Contains(paidOrder.Body.String(), `"status":"paid"`) {
		t.Fatalf("read reconciled order = %d %s", paidOrder.Code, paidOrder.Body)
	}
	createdInvestigation := request(http.MethodPost, "/api/v1/investigations", `{"order_id":"`+order.ID+`","question":"Why is this payment complete?"}`, "lifecycle-investigation")
	if createdInvestigation.Code != http.StatusCreated {
		t.Fatalf("create investigation = %d %s", createdInvestigation.Code, createdInvestigation.Body)
	}
	var investigation struct {
		ID string `json:"id"`
	}
	decode(createdInvestigation, &investigation)
	readInvestigation := request(http.MethodGet, "/api/v1/investigations/"+investigation.ID, "", "")
	listInvestigations := request(http.MethodGet, "/api/v1/orders/"+order.ID+"/investigations", "", "")
	if readInvestigation.Code != http.StatusOK || listInvestigations.Code != http.StatusOK || !strings.Contains(listInvestigations.Body.String(), investigation.ID) {
		t.Fatalf("investigation persistence/read failed: get=%d %s list=%d %s", readInvestigation.Code, readInvestigation.Body, listInvestigations.Code, listInvestigations.Body)
	}
}
