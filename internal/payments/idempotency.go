package payments

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type capturedWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *capturedWriter) WriteHeader(code int) { w.ResponseWriter.WriteHeader(code) }
func (w *capturedWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
func (w *capturedWriter) WriteString(data string) (int, error) {
	w.body.WriteString(data)
	return w.ResponseWriter.WriteString(data)
}

// Idempotency protects mutation routes from duplicate submissions. The same key
// and request replays the original response; a key reused for another request conflicts.
func Idempotency(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Idempotency-Key")
		if len(key) < 8 || len(key) > 128 {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorBody("invalid_idempotency_key", "Idempotency-Key must contain 8 to 128 characters"))
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorBody("invalid_request", "Could not read request body"))
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		hash := sha256.Sum256(append([]byte(c.Request.Method+"\n"+c.Request.URL.Path+"\n"), body...))
		hashString := hex.EncodeToString(hash[:])
		ctx := c.Request.Context()
		var inserted string
		err = db.QueryRow(ctx, `INSERT INTO idempotency_keys(key,request_hash,state) VALUES($1,$2,'processing') ON CONFLICT DO NOTHING RETURNING key`, key, hashString).Scan(&inserted)
		if err != nil && err != pgx.ErrNoRows {
			databaseError(c)
			return
		}
		if inserted == "" {
			var priorHash, state string
			var responseStatus *int
			var response []byte
			err = db.QueryRow(ctx, `SELECT request_hash,state,response_status,stored_response FROM idempotency_keys WHERE key=$1`, key).Scan(&priorHash, &state, &responseStatus, &response)
			if err != nil {
				databaseError(c)
				return
			}
			if priorHash != hashString {
				c.AbortWithStatusJSON(http.StatusConflict, errorBody("idempotency_conflict", "Idempotency-Key was already used for a different request"))
				return
			}
			if state == "completed" && responseStatus != nil {
				c.Data(*responseStatus, "application/json; charset=utf-8", response)
				return
			}
			// A process may have exited while handling the request. Once its bounded
			// provider timeout has elapsed, the operation can be retried safely because
			// the PayPal request ID is deterministic for the order/payment.
			var reclaimed string
			err = db.QueryRow(ctx, `UPDATE idempotency_keys SET updated_at=now() WHERE key=$1 AND state='processing' AND updated_at < now() - interval '2 minutes' RETURNING key`, key).Scan(&reclaimed)
			if err == nil {
				inserted = reclaimed
			} else {
				c.AbortWithStatusJSON(http.StatusConflict, errorBody("idempotency_in_progress", "A request with this key is already being processed"))
				return
			}
		}

		writer := &capturedWriter{ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()
		status := writer.Status()
		if status == 0 {
			status = http.StatusOK
		}
		responseJSON := writer.body.Bytes()
		if len(responseJSON) == 0 {
			responseJSON = []byte(`{}`)
		}
		if !json.Valid(responseJSON) {
			responseJSON, _ = json.Marshal(map[string]string{"message": "response unavailable"})
		}
		_, err = db.Exec(ctx, `UPDATE idempotency_keys SET state='completed',response_status=$2,stored_response=$3::jsonb,updated_at=now() WHERE key=$1`, key, status, responseJSON)
		if err != nil {
			// The mutation result was already sent. Logs make the cache failure visible;
			// provider-side operations also carry a deterministic PayPal request ID.
			c.Error(err)
		}
	}
}

func errorBody(code, message string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": message}}
}

func providerRequestID(prefix, internalID string) string {
	hash := sha256.Sum256([]byte(prefix + ":" + internalID))
	return hex.EncodeToString(hash[:16])
}
