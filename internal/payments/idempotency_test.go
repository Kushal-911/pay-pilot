package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestIdempotencyReplayAndPayloadConflict(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL idempotency integration test")
	}
	ctx := t.Context()
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("paypilot_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`) }()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(ctx, `CREATE TABLE idempotency_keys(key TEXT PRIMARY KEY,request_hash TEXT NOT NULL,state TEXT NOT NULL,response_status INTEGER,stored_response JSONB,created_at TIMESTAMPTZ DEFAULT now(),updated_at TIMESTAMPTZ DEFAULT now())`)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	var calls atomic.Int32
	router.POST("/mutate", Idempotency(db), func(c *gin.Context) {
		calls.Add(1)
		var input map[string]string
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, errorBody("bad_request", "invalid"))
			return
		}
		c.JSON(http.StatusCreated, gin.H{"value": input["value"]})
	})

	first := sendIdempotentRequest(router, "idem-test-001", `{"value":"one"}`)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d, body = %s", first.Code, first.Body)
	}
	replay := sendIdempotentRequest(router, "idem-test-001", `{"value":"one"}`)
	if replay.Code != first.Code {
		t.Fatalf("replay status = %d, want %d", replay.Code, first.Code)
	}
	var firstBody, replayBody map[string]any
	if err := json.Unmarshal(first.Body.Bytes(), &firstBody); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(replay.Body.Bytes(), &replayBody); err != nil {
		t.Fatal(err)
	}
	if firstBody["value"] != replayBody["value"] || calls.Load() != 1 {
		t.Fatalf("replay mismatch or handler reran: first=%v replay=%v calls=%d", firstBody, replayBody, calls.Load())
	}
	conflict := sendIdempotentRequest(router, "idem-test-001", `{"value":"two"}`)
	if conflict.Code != http.StatusConflict {
		t.Fatalf("payload conflict status = %d, want 409", conflict.Code)
	}
	if calls.Load() != 1 {
		t.Fatalf("conflicting payload reached handler; calls=%d", calls.Load())
	}
}

func sendIdempotentRequest(router http.Handler, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mutate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}
