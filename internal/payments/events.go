package payments

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type providerEventView struct {
	ID               string     `json:"id"`
	ProviderEventID  string     `json:"provider_event_id"`
	EventType        string     `json:"event_type"`
	PayPalOrderID    *string    `json:"paypal_order_id"`
	PayPalCaptureID  *string    `json:"paypal_capture_id"`
	ResourceStatus   string     `json:"resource_status,omitempty"`
	AmountMinor      int64      `json:"amount_minor,omitempty"`
	Currency         string     `json:"currency,omitempty"`
	ProcessingStatus string     `json:"processing_status"`
	ReceivedAt       time.Time  `json:"received_at"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
}

// ListOrderEvents returns sanitized, order-scoped provider event evidence.
func ListOrderEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var exists bool
		if err := db.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)`, c.Param("id")).Scan(&exists); err != nil {
			databaseError(c)
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, errorBody("not_found", "Order not found"))
			return
		}

		rows, err := db.Query(c.Request.Context(), `
			SELECT id,provider_event_id,event_type,paypal_order_id,paypal_capture_id,
			       COALESCE(payload->>'resource_status',''),COALESCE((payload->>'amount_minor')::bigint,0),
			       COALESCE(payload->>'currency',''),processing_status,received_at,processed_at
			FROM provider_events
			WHERE order_id=$1
			ORDER BY received_at ASC,id ASC`, c.Param("id"))
		if err != nil {
			databaseError(c)
			return
		}
		defer rows.Close()

		events := make([]providerEventView, 0)
		for rows.Next() {
			var event providerEventView
			if err := rows.Scan(&event.ID, &event.ProviderEventID, &event.EventType, &event.PayPalOrderID, &event.PayPalCaptureID, &event.ResourceStatus, &event.AmountMinor, &event.Currency, &event.ProcessingStatus, &event.ReceivedAt, &event.ProcessedAt); err != nil {
				databaseError(c)
				return
			}
			events = append(events, event)
		}
		if err := rows.Err(); err != nil {
			databaseError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": events})
	}
}
