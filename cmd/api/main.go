package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/investigations"
	"paypilot/internal/payments"
	"paypilot/internal/paypal"
)

//go:embed web/index.html web/app.css web/app.js
var webFiles embed.FS

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	paypalClient := paypal.New(os.Getenv("PAYPAL_BASE_URL"), os.Getenv("PAYPAL_CLIENT_ID"), os.Getenv("PAYPAL_CLIENT_SECRET"), http.DefaultClient)
	var investigationModel investigations.Model
	if os.Getenv("OPENAI_API_KEY") != "" || os.Getenv("OPENAI_MODEL") != "" {
		investigationModel = investigations.OpenAIModel{APIKey: os.Getenv("OPENAI_API_KEY"), Model: os.Getenv("OPENAI_MODEL"), BaseURL: os.Getenv("OPENAI_BASE_URL"), HTTPClient: &http.Client{Timeout: 20 * time.Second}}
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", mustReadWebFile("web/index.html"))
	})
	r.GET("/app.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css; charset=utf-8", mustReadWebFile("web/app.css"))
	})
	r.GET("/app.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", mustReadWebFile("web/app.js"))
	})
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api := r.Group("/api/v1")
	api.POST("/orders", payments.Idempotency(db), payments.CreateOrder(db))
	api.GET("/orders", payments.ListOrders(db))
	api.GET("/orders/:id", payments.GetOrder(db))
	api.GET("/orders/:id/payments", payments.ListOrderPayments(db))
	api.POST("/orders/:id/paypal-order", payments.Idempotency(db), payments.CreatePayPalOrder(db, paypalClient))
	api.POST("/payments/:id/capture", payments.Idempotency(db), payments.CapturePayPalOrder(db, paypalClient))
	api.GET("/payments/:id/refresh", payments.RefreshPayPalOrder(db, paypalClient))
	api.POST("/webhooks/paypal", payments.PayPalWebhook(db, paypalClient, os.Getenv("PAYPAL_WEBHOOK_ID")))
	api.POST("/investigations", payments.Idempotency(db), investigations.CreateHandler(db, investigationModel))
	api.GET("/investigations/:id", investigations.GetHandler(db))
	api.GET("/orders/:id/investigations", investigations.ListForOrderHandler(db))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func mustReadWebFile(name string) []byte {
	data, err := webFiles.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}
