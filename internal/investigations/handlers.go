package investigations

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type createInput struct {
	OrderID  string `json:"order_id" binding:"required"`
	Question string `json:"question" binding:"required,min=3,max=2000"`
}

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func CreateHandler(db *pgxpool.Pool, model Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse("invalid_request", "Provide an order_id UUID and a question up to 2000 characters."))
			return
		}
		if !uuidPattern.MatchString(input.OrderID) {
			c.JSON(http.StatusBadRequest, errorResponse("invalid_order_id", "order_id must be a UUID."))
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
		defer cancel()
		service := Service{Store: PostgresStore{DB: db}, Model: model}
		report, err := service.Investigate(ctx, input.OrderID, input.Question)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse("order_not_found", "Order not found."))
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("investigation_error", "Could not load payment evidence."))
			return
		}
		reportJSON, err := json.Marshal(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("investigation_error", "Could not encode investigation report."))
			return
		}
		evidenceJSON, err := json.Marshal(report.Evidence)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("investigation_error", "Could not encode investigation evidence."))
			return
		}
		var investigationID string
		err = db.QueryRow(ctx, `INSERT INTO investigations(order_id,question,finding,confidence,recommendation,evidence,report) VALUES($1,$2,$3,$4,$5,$6::jsonb,$7::jsonb) RETURNING id`, input.OrderID, input.Question, report.Finding, report.Confidence, report.Recommendation, evidenceJSON, reportJSON).Scan(&investigationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("persistence_error", "Could not save the investigation."))
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": investigationID, "order_id": input.OrderID, "question": input.Question, "status": "open", "report": report})
	}
}

func GetHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !uuidPattern.MatchString(c.Param("id")) {
			c.JSON(http.StatusBadRequest, errorResponse("invalid_investigation_id", "id must be a UUID."))
			return
		}
		var id, orderID, question, finding, recommendation, status string
		var confidence float64
		var evidence, report []byte
		var created time.Time
		err := db.QueryRow(c.Request.Context(), `SELECT id,order_id,question,finding,confidence,recommendation,status,evidence,report,created_at FROM investigations WHERE id=$1`, c.Param("id")).Scan(&id, &orderID, &question, &finding, &confidence, &recommendation, &status, &evidence, &report, &created)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse("not_found", "Investigation not found."))
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("database_error", "Could not load investigation."))
			return
		}
		var decodedReport any
		if err := json.Unmarshal(report, &decodedReport); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("database_error", "Stored report could not be decoded."))
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "order_id": orderID, "question": question, "finding": finding, "confidence": confidence, "recommendation": recommendation, "status": status, "evidence": json.RawMessage(evidence), "report": decodedReport, "created_at": created})
	}
}

func ListForOrderHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !uuidPattern.MatchString(c.Param("id")) {
			c.JSON(http.StatusBadRequest, errorResponse("invalid_order_id", "id must be a UUID."))
			return
		}
		rows, err := db.Query(c.Request.Context(), `SELECT id,finding,confidence,recommendation,status,created_at FROM investigations WHERE order_id=$1 ORDER BY created_at DESC`, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("database_error", "Could not list investigations."))
			return
		}
		defer rows.Close()
		items := make([]gin.H, 0)
		for rows.Next() {
			var id, finding, recommendation, status string
			var confidence float64
			var created time.Time
			if err := rows.Scan(&id, &finding, &confidence, &recommendation, &status, &created); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse("database_error", "Could not list investigations."))
				return
			}
			items = append(items, gin.H{"id": id, "finding": finding, "confidence": confidence, "recommendation": recommendation, "status": status, "created_at": created})
		}
		if rows.Err() != nil {
			c.JSON(http.StatusInternalServerError, errorResponse("database_error", "Could not list investigations."))
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	}
}

func errorResponse(code, message string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": message}}
}
