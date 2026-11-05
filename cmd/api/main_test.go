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
