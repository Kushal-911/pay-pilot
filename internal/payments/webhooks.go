package payments

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/paypal"
)

type webhookEvent struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	Resource  struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Amount struct {
			Value        string `json:"value"`
			CurrencyCode string `json:"currency_code"`
		} `json:"amount"`
		SupplementaryData struct {
			RelatedIDs struct {
				OrderID string `json:"order_id"`
			} `json:"related_ids"`
		} `json:"supplementary_data"`
	} `json:"resource"`
}

func PayPalWebhook(db *pgxpool.Pool, client *paypal.Client, webhookID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if webhookID == "" {
			c.JSON(http.StatusServiceUnavailable, errorBody("webhook_unconfigured", "PayPal webhook ID is not configured"))
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, errorBody("payload_too_large", "Webhook payload exceeds 1 MiB"))
			return
		}
		var event webhookEvent
		if err := json.Unmarshal(raw, &event); err != nil || event.ID == "" || len(event.ID) > 200 || event.EventType == "" || len(event.EventType) > 200 {
			c.JSON(http.StatusBadRequest, errorBody("invalid_webhook", "Webhook body is invalid"))
			return
		}
		headers := paypal.WebhookHeaders{AuthAlgo: c.GetHeader("paypal-auth-algo"), CertURL: c.GetHeader("paypal-cert-url"), TransmissionID: c.GetHeader("paypal-transmission-id"), TransmissionSig: c.GetHeader("paypal-transmission-sig"), TransmissionTime: c.GetHeader("paypal-transmission-time")}
		if headers.AuthAlgo == "" || headers.CertURL == "" || headers.TransmissionID == "" || headers.TransmissionSig == "" || headers.TransmissionTime == "" {
			c.JSON(http.StatusUnauthorized, errorBody("invalid_signature", "Required PayPal transmission headers are missing"))
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
		defer cancel()
		valid, err := client.VerifyWebhook(ctx, webhookID, headers, json.RawMessage(raw))
		if err != nil {
			c.JSON(http.StatusBadGateway, errorBody("verification_unavailable", "PayPal signature verification could not be completed"))
			return
		}
		if !valid {
			c.JSON(http.StatusUnauthorized, errorBody("invalid_signature", "PayPal webhook signature is invalid"))
			return
		}
		if len(event.Resource.Amount.Value) > 0 {
			if _, err := parseMinorAmount(event.Resource.Amount.Value); err != nil {
				c.JSON(http.StatusBadRequest, errorBody("invalid_webhook", "Webhook amount format is invalid"))
				return
			}
		}
		if len(event.Resource.ID) > 200 || len(event.Resource.SupplementaryData.RelatedIDs.OrderID) > 200 {
			c.JSON(http.StatusBadRequest, errorBody("invalid_webhook", "Webhook provider identifier is invalid"))
			return
		}
		if err := persistAndReconcileWebhook(ctx, db, event); err != nil {
			c.JSON(http.StatusInternalServerError, errorBody("webhook_processing_error", "Verified webhook could not be persisted"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": true, "event_id": event.ID})
	}
}

func persistAndReconcileWebhook(ctx context.Context, db *pgxpool.Pool, event webhookEvent) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var orderID *string
	var attemptID, attemptOrderID, attemptStatus string
	var attemptAmount int64
	var attemptCurrency string
	var captureID *string
	if event.Resource.SupplementaryData.RelatedIDs.OrderID != "" {
		err = tx.QueryRow(ctx, `SELECT id,order_id,status,amount_minor,currency,paypal_capture_id FROM payment_attempts WHERE paypal_order_id=$1 FOR UPDATE`, event.Resource.SupplementaryData.RelatedIDs.OrderID).Scan(&attemptID, &attemptOrderID, &attemptStatus, &attemptAmount, &attemptCurrency, &captureID)
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		} else if err != nil {
			return err
		} else {
			orderID = &attemptOrderID
		}
	}
	if attemptID == "" && event.Resource.ID != "" && event.Resource.SupplementaryData.RelatedIDs.OrderID == "" {
		err = tx.QueryRow(ctx, `SELECT id,order_id,status,amount_minor,currency,paypal_capture_id FROM payment_attempts WHERE paypal_capture_id=$1 FOR UPDATE`, event.Resource.ID).Scan(&attemptID, &attemptOrderID, &attemptStatus, &attemptAmount, &attemptCurrency, &captureID)
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		} else if err != nil {
			return err
		} else {
			orderID = &attemptOrderID
		}
	}
	minorAmount := int64(0)
	if event.Resource.Amount.Value != "" {
		minorAmount, err = parseMinorAmount(event.Resource.Amount.Value)
		if err != nil {
			return err
		}
	}
	payload, _ := json.Marshal(map[string]any{"resource_status": event.Resource.Status, "amount_minor": minorAmount, "currency": strings.ToUpper(event.Resource.Amount.CurrencyCode)})
	var eventRecordID string
	err = tx.QueryRow(ctx, `INSERT INTO provider_events(provider_event_id,order_id,event_type,paypal_order_id,paypal_capture_id,payload,processing_status) VALUES($1,$2,$3,NULLIF($4,''),NULLIF($5,''),$6::jsonb,'received') ON CONFLICT(provider_event_id) DO NOTHING RETURNING id`, event.ID, orderID, event.EventType, event.Resource.SupplementaryData.RelatedIDs.OrderID, event.Resource.ID, payload).Scan(&eventRecordID)
	if errors.Is(err, pgx.ErrNoRows) {
		return tx.Commit(ctx)
	}
	if err != nil {
		return err
	}
	if attemptID == "" {
		return tx.Commit(ctx)
	}

	newAttemptStatus, processingStatus, apply := classifyCaptureWebhook(event, attemptStatus, captureID, attemptAmount, attemptCurrency, minorAmount)
	if apply {
		switch newAttemptStatus {
		case "completed":
			_, err = tx.Exec(ctx, `UPDATE payment_attempts SET paypal_capture_id=$1,status='completed',updated_at=now() WHERE id=$2`, event.Resource.ID, attemptID)
			if err == nil {
				_, err = tx.Exec(ctx, `UPDATE orders SET status='paid',updated_at=now() WHERE id=$1`, attemptOrderID)
			}
		case "unknown":
			_, err = tx.Exec(ctx, `UPDATE payment_attempts SET status='unknown',updated_at=now() WHERE id=$1 AND status NOT IN ('completed','failed')`, attemptID)
		case "failed":
			_, err = tx.Exec(ctx, `UPDATE payment_attempts SET status='failed',updated_at=now() WHERE id=$1 AND status<>'completed'`, attemptID)
			if err == nil {
				_, err = tx.Exec(ctx, `UPDATE orders SET status='failed',updated_at=now() WHERE id=$1 AND status='pending' AND EXISTS(SELECT 1 FROM payment_attempts WHERE order_id=$1) AND NOT EXISTS(SELECT 1 FROM payment_attempts WHERE order_id=$1 AND status<>'failed')`, attemptOrderID)
			}
		}
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE provider_events SET processing_status=$1,processed_at=now() WHERE id=$2`, processingStatus, eventRecordID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func classifyCaptureWebhook(event webhookEvent, currentStatus string, currentCaptureID *string, expectedAmount int64, expectedCurrency string, eventAmount int64) (string, string, bool) {
	amountMatches := eventAmount > 0 && eventAmount == expectedAmount && strings.EqualFold(event.Resource.Amount.CurrencyCode, expectedCurrency)
	switch event.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		if event.Resource.Status != "COMPLETED" || event.Resource.ID == "" || !amountMatches {
			return currentStatus, "failed", false
		}
		if currentStatus == "completed" && currentCaptureID != nil && *currentCaptureID != event.Resource.ID {
			return currentStatus, "failed", false
		}
		return "completed", "processed", true
	case "PAYMENT.CAPTURE.PENDING":
		if event.Resource.Status != "PENDING" || !amountMatches {
			return currentStatus, "failed", false
		}
		if currentStatus == "completed" || currentStatus == "failed" {
			return currentStatus, "processed", false
		}
		return "unknown", "processed", true
	case "PAYMENT.CAPTURE.DENIED":
		if event.Resource.Status != "DENIED" || !amountMatches {
			return currentStatus, "failed", false
		}
		if currentStatus == "completed" {
			return currentStatus, "processed", false
		}
		return "failed", "processed", true
	default:
		return currentStatus, "processed", false
	}
}

var decimalAmountPattern = regexp.MustCompile(`^[0-9]+(?:\.[0-9]{1,2})?$`)

func parseMinorAmount(value string) (int64, error) {
	if !decimalAmountPattern.MatchString(value) {
		return 0, errors.New("invalid amount format")
	}
	wholePart, fractionPart, hasFraction := strings.Cut(value, ".")
	whole, err := strconv.ParseInt(wholePart, 10, 64)
	if err != nil || whole > (math.MaxInt64/100) {
		return 0, errors.New("amount out of range")
	}
	minor := whole * 100
	if hasFraction {
		if len(fractionPart) == 1 {
			fractionPart += "0"
		}
		fraction, err := strconv.ParseInt(fractionPart, 10, 64)
		if err != nil {
			return 0, errors.New("invalid amount format")
		}
		if minor > math.MaxInt64-fraction {
			return 0, errors.New("amount out of range")
		}
		minor += fraction
	}
	return minor, nil
}
