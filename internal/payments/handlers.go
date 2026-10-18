package payments

import (
	"context"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/paypal"
)

type orderInput struct {
	MerchantReference string `json:"merchant_reference" binding:"required,max=120"`
	AmountMinor       int64  `json:"amount_minor" binding:"required,gt=0"`
	Currency          string `json:"currency" binding:"required,len=3"`
	Description       string `json:"description" binding:"required,max=255"`
}

var referencePattern = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

func CreateOrder(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in orderInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": "Invalid order fields"}})
			return
		}
		in.Currency = strings.ToUpper(in.Currency)
		if !referencePattern.MatchString(in.MerchantReference) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(in.Currency) {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": "Invalid merchant reference or currency"}})
			return
		}
		var id string
		err := db.QueryRow(c.Request.Context(), `INSERT INTO orders(merchant_reference,amount_minor,currency,description,status) VALUES($1,$2,$3,$4,'pending') ON CONFLICT(merchant_reference) DO UPDATE SET merchant_reference=EXCLUDED.merchant_reference WHERE orders.amount_minor=EXCLUDED.amount_minor AND orders.currency=EXCLUDED.currency AND orders.description=EXCLUDED.description RETURNING id`, in.MerchantReference, in.AmountMinor, in.Currency, in.Description).Scan(&id)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "order_conflict", "message": "Could not create order; merchant reference may already exist"}})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id, "merchant_reference": in.MerchantReference, "amount_minor": in.AmountMinor, "currency": in.Currency, "description": in.Description, "status": "pending"})
	}
}

func ListOrders(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(c.Request.Context(), `SELECT id,merchant_reference,amount_minor,currency,description,status,created_at,updated_at FROM orders ORDER BY created_at DESC LIMIT 100`)
		if err != nil {
			databaseError(c)
			return
		}
		defer rows.Close()
		orders := make([]map[string]any, 0)
		for rows.Next() {
			var id, ref, currency, description, status string
			var amount int64
			var created, updated time.Time
			if err := rows.Scan(&id, &ref, &amount, &currency, &description, &status, &created, &updated); err != nil {
				databaseError(c)
				return
			}
			orders = append(orders, map[string]any{"id": id, "merchant_reference": ref, "amount_minor": amount, "currency": currency, "description": description, "status": status, "created_at": created, "updated_at": updated})
		}
		if rows.Err() != nil {
			databaseError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders})
	}
}

func GetOrder(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id, ref, currency, description, status string
		var amount int64
		var created, updated time.Time
		err := db.QueryRow(c.Request.Context(), `SELECT id,merchant_reference,amount_minor,currency,description,status,created_at,updated_at FROM orders WHERE id=$1`, c.Param("id")).Scan(&id, &ref, &amount, &currency, &description, &status, &created, &updated)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Order not found"}})
			return
		}
		if err != nil {
			databaseError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "merchant_reference": ref, "amount_minor": amount, "currency": currency, "description": description, "status": status, "created_at": created, "updated_at": updated})
	}
}

func ListOrderPayments(db *pgxpool.Pool) gin.HandlerFunc {
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
		rows, err := db.Query(c.Request.Context(), `SELECT id,paypal_order_id,paypal_capture_id,amount_minor,currency,status,created_at,updated_at FROM payment_attempts WHERE order_id=$1 ORDER BY created_at DESC`, c.Param("id"))
		if err != nil {
			databaseError(c)
			return
		}
		defer rows.Close()
		attempts := make([]map[string]any, 0)
		for rows.Next() {
			var id, status, currency string
			var paypalOrderID, captureID *string
			var amount int64
			var created, updated time.Time
			if err := rows.Scan(&id, &paypalOrderID, &captureID, &amount, &currency, &status, &created, &updated); err != nil {
				databaseError(c)
				return
			}
			attempts = append(attempts, map[string]any{"id": id, "paypal_order_id": paypalOrderID, "paypal_capture_id": captureID, "amount_minor": amount, "currency": currency, "status": status, "created_at": created, "updated_at": updated})
		}
		if rows.Err() != nil {
			databaseError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": attempts})
	}
}

func CreatePayPalOrder(db *pgxpool.Pool, client *paypal.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
		defer cancel()
		var id, ref, currency, description, status string
		var amount int64
		err := db.QueryRow(ctx, `SELECT id,merchant_reference,currency,description,status,amount_minor FROM orders WHERE id=$1`, c.Param("id")).Scan(&id, &ref, &currency, &description, &status, &amount)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Order not found"}})
			return
		}
		if err != nil {
			databaseError(c)
			return
		}
		if status != "pending" {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "invalid_state", "message": "Order is not eligible for checkout"}})
			return
		}
		returnURL, cancelURL := os.Getenv("PAYPAL_RETURN_URL"), os.Getenv("PAYPAL_CANCEL_URL")
		if returnURL == "" || cancelURL == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{"code": "provider_unconfigured", "message": "PayPal return and cancel URLs are required"}})
			return
		}
		result, err := client.CreateOrder(ctx, ref, description, currency, amount, returnURL, cancelURL, providerRequestID("order", id))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "provider_error", "message": "PayPal order creation failed"}})
			return
		}
		var attemptID string
		err = db.QueryRow(ctx, `INSERT INTO payment_attempts(order_id,paypal_order_id,amount_minor,currency,status) VALUES($1,$2,$3,$4,'created') RETURNING id`, id, result.ID, amount, currency).Scan(&attemptID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "persistence_error", "message": "PayPal order created but could not be recorded; contact support before retrying"}})
			return
		}
		approvalURL := ""
		for _, link := range result.Links {
			if link.Rel == "approve" {
				approvalURL = link.Href
			}
		}
		c.JSON(http.StatusCreated, gin.H{"payment_id": attemptID, "paypal_order_id": result.ID, "status": result.Status, "approval_url": approvalURL})
	}
}

func CapturePayPalOrder(db *pgxpool.Pool, client *paypal.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()
		var id, orderID, status string
		err := db.QueryRow(ctx, `SELECT id,paypal_order_id,status FROM payment_attempts WHERE id=$1`, c.Param("id")).Scan(&id, &orderID, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Payment attempt not found"}})
			return
		}
		if err != nil {
			databaseError(c)
			return
		}
		if status == "completed" {
			var captureID string
			if err := db.QueryRow(ctx, `SELECT paypal_capture_id FROM payment_attempts WHERE id=$1`, id).Scan(&captureID); err != nil {
				databaseError(c)
				return
			}
			c.JSON(http.StatusOK, gin.H{"payment_id": id, "paypal_order_id": orderID, "paypal_capture_id": captureID, "status": "completed"})
			return
		}
		if status != "created" {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "invalid_state", "message": "Payment attempt is not eligible for capture"}})
			return
		}
		var locked string
		err = db.QueryRow(ctx, `UPDATE payment_attempts SET status='capturing',updated_at=now() WHERE id=$1 AND status='created' RETURNING id`, id).Scan(&locked)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "invalid_state", "message": "Payment attempt is already being processed or changed state"}})
			return
		}
		result, err := client.Capture(ctx, orderID, providerRequestID("capture", id))
		if err != nil {
			_, _ = db.Exec(ctx, `UPDATE payment_attempts SET status='unknown',updated_at=now() WHERE id=$1 AND status='capturing'`, id)
			c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "provider_error", "message": "PayPal capture outcome could not be confirmed; refresh provider state before retrying"}})
			return
		}
		captureID, captureStatus := "", ""
		if len(result.PurchaseUnits) > 0 && len(result.PurchaseUnits[0].Payments.Captures) > 0 {
			captureID = result.PurchaseUnits[0].Payments.Captures[0].ID
			captureStatus = result.PurchaseUnits[0].Payments.Captures[0].Status
		}
		if result.Status != "COMPLETED" || captureStatus != "COMPLETED" || captureID == "" {
			_, _ = db.Exec(ctx, `UPDATE payment_attempts SET status='unknown',updated_at=now() WHERE id=$1 AND status='capturing'`, id)
			c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "capture_incomplete", "message": "PayPal did not confirm a completed capture; refresh provider state before retrying"}, "provider_status": result.Status})
			return
		}
		tx, err := db.Begin(ctx)
		if err != nil {
			databaseError(c)
			return
		}
		defer tx.Rollback(ctx)
		_, err = tx.Exec(ctx, `UPDATE payment_attempts SET paypal_capture_id=$1,status='completed',updated_at=now() WHERE id=$2 AND status='capturing'`, captureID, id)
		if err == nil {
			_, err = tx.Exec(ctx, `UPDATE orders SET status='paid',updated_at=now() WHERE id=(SELECT order_id FROM payment_attempts WHERE id=$1)`, id)
		}
		if err == nil {
			err = tx.Commit(ctx)
		}
		if err != nil {
			databaseError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"payment_id": id, "paypal_order_id": orderID, "paypal_capture_id": captureID, "status": "completed"})
	}
}

func RefreshPayPalOrder(db *pgxpool.Pool, client *paypal.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
		defer cancel()
		var id, orderID, status string
		err := db.QueryRow(ctx, `SELECT id,paypal_order_id,status FROM payment_attempts WHERE id=$1`, c.Param("id")).Scan(&id, &orderID, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "Payment attempt not found"}})
			return
		}
		if err != nil {
			databaseError(c)
			return
		}
		result, err := client.GetOrder(ctx, orderID)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": "provider_error", "message": "Could not refresh PayPal order state"}})
			return
		}
		captureID, captureStatus := "", ""
		if len(result.PurchaseUnits) > 0 && len(result.PurchaseUnits[0].Payments.Captures) > 0 {
			captureID = result.PurchaseUnits[0].Payments.Captures[0].ID
			captureStatus = result.PurchaseUnits[0].Payments.Captures[0].Status
		}
		if result.Status == "COMPLETED" && captureID != "" && captureStatus == "COMPLETED" {
			tx, err := db.Begin(ctx)
			if err != nil {
				databaseError(c)
				return
			}
			defer tx.Rollback(ctx)
			var updatedID string
			err = tx.QueryRow(ctx, `UPDATE payment_attempts SET paypal_capture_id=$1,status='completed',updated_at=now() WHERE id=$2 AND status IN ('unknown','capturing','created','completed') RETURNING id`, captureID, id).Scan(&updatedID)
			if err == nil {
				_, err = tx.Exec(ctx, `UPDATE orders SET status='paid',updated_at=now() WHERE id=(SELECT order_id FROM payment_attempts WHERE id=$1)`, id)
			}
			if err == nil {
				err = tx.Commit(ctx)
			}
			if err != nil {
				databaseError(c)
				return
			}
			status = "completed"
		} else if result.Status == "APPROVED" && (status == "unknown" || status == "capturing") {
			var updatedStatus string
			err = db.QueryRow(ctx, `UPDATE payment_attempts SET status='created',updated_at=now() WHERE id=$1 AND status IN ('unknown','capturing') RETURNING status`, id).Scan(&updatedStatus)
			if err == nil {
				status = updatedStatus
			} else if errors.Is(err, pgx.ErrNoRows) {
				err = db.QueryRow(ctx, `SELECT status FROM payment_attempts WHERE id=$1`, id).Scan(&status)
				if err != nil {
					databaseError(c)
					return
				}
			} else {
				databaseError(c)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"payment_id": id, "paypal_order_id": orderID, "provider_status": result.Status, "status": status, "paypal_capture_id": captureID})
	}
}

func databaseError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "database_error", "message": "Database operation failed"}})
}
