package payments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListOrderEventsIsScopedSanitizedAndOrdered(t *testing.T) {
	db := newPostgresTestSchema(t)
	ctx := t.Context()
	orderID := "223e4567-e89b-42d3-a456-426614174000"
	otherOrderID := "323e4567-e89b-42d3-a456-426614174000"
	emptyOrderID := "423e4567-e89b-42d3-a456-426614174000"
	for i, id := range []string{orderID, otherOrderID, emptyOrderID} {
		if _, err := db.Exec(ctx, `INSERT INTO orders(id,merchant_reference,amount_minor,currency,description,status) VALUES($1,$2,2500,'USD','event integration','pending')`, id, fmt.Sprintf("event-order-%d", i)); err != nil {
			t.Fatal(err)
		}
	}
	for _, event := range []struct {
		id, eventID, typ, status, received string
	}{
		{"523e4567-e89b-42d3-a456-426614174000", "EVT-2", "PAYMENT.CAPTURE.PENDING", "processed", "2026-11-05T10:00:02Z"},
		{"623e4567-e89b-42d3-a456-426614174000", "EVT-1", "PAYMENT.CAPTURE.COMPLETED", "processed", "2026-11-05T10:00:01Z"},
	} {
		_, err := db.Exec(ctx, `INSERT INTO provider_events(id,provider_event_id,order_id,event_type,paypal_order_id,paypal_capture_id,payload,processing_status,received_at) VALUES($1,$2,$3,$4,'PAYPAL-ORDER','CAPTURE-1',$5::jsonb,$6,$7)`, event.id, event.eventID, orderID, event.typ, `{"resource_status":"COMPLETED","amount_minor":2500,"currency":"USD","raw_secret":"must-not-leak"}`, event.status, event.received)
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err := db.Exec(ctx, `INSERT INTO provider_events(id,provider_event_id,order_id,event_type,payload,processing_status,received_at) VALUES($1,'EVT-OTHER',$2,'PAYMENT.CAPTURE.COMPLETED','{}','processed',now())`, "723e4567-e89b-42d3-a456-426614174000", otherOrderID)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/orders/:id/events", ListOrderEvents(db))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID+"/events", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("event list status = %d, body = %s", response.Code, response.Body)
	}
	var result struct {
		Data []providerEventView `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 2 || result.Data[0].ProviderEventID != "EVT-1" || result.Data[1].ProviderEventID != "EVT-2" {
		t.Fatalf("events are not scoped and ordered by received_at: %+v", result.Data)
	}
	if strings.Contains(response.Body.String(), "raw_secret") || strings.Contains(response.Body.String(), "must-not-leak") {
		t.Fatalf("raw webhook data leaked in response: %s", response.Body)
	}

	empty := httptest.NewRecorder()
	router.ServeHTTP(empty, httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+emptyOrderID+"/events", nil))
	if empty.Code != http.StatusOK || !strings.Contains(empty.Body.String(), `"data":[]`) {
		t.Fatalf("empty event response = %d %s", empty.Code, empty.Body)
	}
	missing := httptest.NewRecorder()
	router.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/v1/orders/00000000-0000-4000-8000-000000000000/events", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing order status = %d, want 404", missing.Code)
	}
}
