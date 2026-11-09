package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/paypal"
)

func TestLivenessAndReadinessAreDistinct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := pgxpool.New(t.Context(), "postgres://missing:missing@127.0.0.1:1/missing?connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	router := newRouter(db, paypal.New("", "", "", nil), nil)

	live := httptest.NewRecorder()
	router.ServeHTTP(live, httptest.NewRequest(http.MethodGet, "/livez", nil))
	if live.Code != http.StatusOK || live.Body.String() != `{"status":"alive"}` {
		t.Fatalf("liveness = %d %s", live.Code, live.Body)
	}
	ready := httptest.NewRecorder()
	router.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if ready.Code != http.StatusServiceUnavailable || !strings.Contains(ready.Body.String(), `"dependency":"database"`) {
		t.Fatalf("readiness = %d %s", ready.Code, ready.Body)
	}
	legacyHealth := httptest.NewRecorder()
	router.ServeHTTP(legacyHealth, httptest.NewRequest(http.MethodGet, "/health", nil))
	if legacyHealth.Code != http.StatusServiceUnavailable {
		t.Fatalf("legacy health alias status = %d, want database readiness status", legacyHealth.Code)
	}
}

func TestPublicDemoIsSyntheticAndDoesNotExposeLiveAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newPublicDemoRouter()

	ready := httptest.NewRecorder()
	router.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if ready.Code != http.StatusOK || !strings.Contains(ready.Body.String(), `"mode":"synthetic-demo"`) {
		t.Fatalf("public demo readiness = %d %s", ready.Code, ready.Body)
	}

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/", nil))
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), "PAYPILOT_PUBLIC_DEMO=true") {
		t.Fatalf("public demo page = %d; demo mode flag missing", page.Code)
	}

	for _, request := range []struct{ method, path string }{
		{http.MethodGet, "/api/v1/orders"},
		{http.MethodPost, "/api/v1/orders"},
		{http.MethodPost, "/api/v1/orders/order-id/paypal-order"},
		{http.MethodPost, "/api/v1/payments/payment-id/capture"},
		{http.MethodPost, "/api/v1/webhooks/paypal"},
		{http.MethodPost, "/api/v1/investigations"},
	} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(request.method, request.path, nil))
		if response.Code != http.StatusNotFound {
			t.Errorf("public demo route %s %s status = %d, want 404", request.method, request.path, response.Code)
		}
	}
}
