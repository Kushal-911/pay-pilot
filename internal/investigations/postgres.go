package investigations

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct{ DB *pgxpool.Pool }

func (s PostgresStore) Load(ctx context.Context, orderID string) (Snapshot, error) {
	var out Snapshot
	var orderCreated, orderUpdated time.Time
	err := s.DB.QueryRow(ctx, `SELECT id,merchant_reference,amount_minor,currency,description,status,created_at,updated_at FROM orders WHERE id=$1`, orderID).Scan(
		&out.Order.ID, &out.Order.MerchantReference, &out.Order.AmountMinor, &out.Order.Currency, &out.Order.Description, &out.Order.Status, &orderCreated, &orderUpdated)
	if err != nil {
		return out, err
	}
	out.Order.CreatedAt, out.Order.UpdatedAt = orderCreated.UTC().Format(time.RFC3339Nano), orderUpdated.UTC().Format(time.RFC3339Nano)
	rows, err := s.DB.Query(ctx, `SELECT id,paypal_order_id,paypal_capture_id,amount_minor,currency,status,created_at,updated_at FROM payment_attempts WHERE order_id=$1 ORDER BY created_at ASC,id ASC`, orderID)
	if err != nil {
		return out, err
	}
	for rows.Next() {
		var attempt AttemptRecord
		var created, updated time.Time
		if err := rows.Scan(&attempt.ID, &attempt.PayPalOrderID, &attempt.PayPalCaptureID, &attempt.AmountMinor, &attempt.Currency, &attempt.Status, &created, &updated); err != nil {
			rows.Close()
			return out, err
		}
		attempt.CreatedAt, attempt.UpdatedAt = created.UTC().Format(time.RFC3339Nano), updated.UTC().Format(time.RFC3339Nano)
		out.Attempts = append(out.Attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return out, err
	}
	rows.Close()
	orderIDs, captureIDs := make([]string, 0, len(out.Attempts)), make([]string, 0, len(out.Attempts))
	for _, attempt := range out.Attempts {
		if attempt.PayPalOrderID != nil {
			orderIDs = append(orderIDs, *attempt.PayPalOrderID)
		}
		if attempt.PayPalCaptureID != nil {
			captureIDs = append(captureIDs, *attempt.PayPalCaptureID)
		}
	}
	eventRows, err := s.DB.Query(ctx, `SELECT id,provider_event_id,event_type,paypal_order_id,paypal_capture_id,processing_status,received_at FROM provider_events WHERE order_id=$1 OR paypal_order_id=ANY($2::text[]) OR paypal_capture_id=ANY($3::text[]) ORDER BY received_at ASC,id ASC`, orderID, orderIDs, captureIDs)
	if err != nil {
		return out, err
	}
	defer eventRows.Close()
	for eventRows.Next() {
		var event EventRecord
		var received time.Time
		if err := eventRows.Scan(&event.ID, &event.ProviderEventID, &event.EventType, &event.PayPalOrderID, &event.PayPalCaptureID, &event.ProcessingStatus, &received); err != nil {
			return out, err
		}
		event.ReceivedAt = received.UTC().Format(time.RFC3339Nano)
		out.Events = append(out.Events, event)
	}
	if eventRows.Err() != nil {
		return out, eventRows.Err()
	}
	return out, nil
}

func BuildEvidence(snapshot Snapshot) []Evidence {
	evidence := []Evidence{{ID: "order:" + snapshot.Order.ID, RecordType: "order", RecordID: snapshot.Order.ID, Summary: "Merchant order record", Data: map[string]any{
		"merchant_reference": snapshot.Order.MerchantReference, "amount_minor": snapshot.Order.AmountMinor, "currency": snapshot.Order.Currency, "description": snapshot.Order.Description, "status": snapshot.Order.Status, "created_at": snapshot.Order.CreatedAt, "updated_at": snapshot.Order.UpdatedAt,
	}}}
	for _, attempt := range snapshot.Attempts {
		evidence = append(evidence, Evidence{ID: "attempt:" + attempt.ID, RecordType: "payment_attempt", RecordID: attempt.ID, Summary: "PayPal payment attempt", Data: map[string]any{
			"paypal_order_id": attempt.PayPalOrderID, "paypal_capture_id": attempt.PayPalCaptureID, "amount_minor": attempt.AmountMinor, "currency": attempt.Currency, "status": attempt.Status, "created_at": attempt.CreatedAt, "updated_at": attempt.UpdatedAt,
		}})
	}
	for _, event := range snapshot.Events {
		evidence = append(evidence, Evidence{ID: "event:" + event.ID, RecordType: "provider_event", RecordID: event.ID, Summary: "Persisted provider event", Data: map[string]any{
			"provider_event_id": event.ProviderEventID, "event_type": event.EventType, "paypal_order_id": event.PayPalOrderID, "paypal_capture_id": event.PayPalCaptureID, "processing_status": event.ProcessingStatus, "received_at": event.ReceivedAt,
		}})
	}
	return evidence
}
