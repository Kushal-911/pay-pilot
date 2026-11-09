package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Public hosting is restricted to the synthetic, browser-only workbench.
	// The live API remains available only in the normal database-backed mode.
	if os.Getenv("PUBLIC_DEMO_MODE") == "true" {
		if err := newPublicDemoRouter().Run(":" + port); err != nil {
			log.Fatal(err)
		}
		return
	}

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
	paypalClient := paypal.New(os.Getenv("PAYPAL_BASE_URL"), os.Getenv("PAYPAL_CLIENT_ID"), os.Getenv("PAYPAL_CLIENT_SECRET"), &http.Client{Timeout: 25 * time.Second})
	var investigationModel investigations.Model
	if os.Getenv("OPENAI_API_KEY") != "" || os.Getenv("OPENAI_MODEL") != "" {
		investigationModel = investigations.OpenAIModel{APIKey: os.Getenv("OPENAI_API_KEY"), Model: os.Getenv("OPENAI_MODEL"), BaseURL: os.Getenv("OPENAI_BASE_URL"), HTTPClient: &http.Client{Timeout: 20 * time.Second}}
	}
	r := newRouter(db, paypalClient, investigationModel)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func newRouter(db *pgxpool.Pool, paypalClient *paypal.Client, investigationModel investigations.Model) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	registerWebRoutes(r, false)
	r.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	ready := func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "dependency": "database"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
	r.GET("/readyz", ready)
	r.GET("/health", ready)
	api := r.Group("/api/v1")
	api.POST("/orders", payments.Idempotency(db), payments.CreateOrder(db))
	api.GET("/orders", payments.ListOrders(db))
	api.GET("/orders/:id", payments.GetOrder(db))
	api.GET("/orders/:id/payments", payments.ListOrderPayments(db))
	api.GET("/orders/:id/events", payments.ListOrderEvents(db))
	api.POST("/orders/:id/paypal-order", payments.Idempotency(db), payments.CreatePayPalOrder(db, paypalClient))
	api.POST("/payments/:id/capture", payments.Idempotency(db), payments.CapturePayPalOrder(db, paypalClient))
	api.GET("/payments/:id/refresh", payments.RefreshPayPalOrder(db, paypalClient))
	api.POST("/webhooks/paypal", payments.PayPalWebhook(db, paypalClient, os.Getenv("PAYPAL_WEBHOOK_ID")))
	api.POST("/investigations", payments.Idempotency(db), investigations.CreateHandler(db, investigationModel))
	api.GET("/investigations/:id", investigations.GetHandler(db))
	api.GET("/orders/:id/investigations", investigations.ListForOrderHandler(db))
	return r
}

func newPublicDemoRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}
	registerWebRoutes(r, true)
	r.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready", "mode": "synthetic-demo"})
	})
	return r
}

func registerWebRoutes(r *gin.Engine, publicDemo bool) {
	r.GET("/", func(c *gin.Context) {
		index := mustReadWebFile("web/index.html")
		if publicDemo {
			index = []byte(strings.Replace(string(index), "</head>", `<script>window.PAYPILOT_PUBLIC_DEMO=true;</script></head>`, 1))
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	r.GET("/app.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css; charset=utf-8", mustReadWebFile("web/app.css"))
	})
	r.GET("/app.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", mustReadWebFile("web/app.js"))
	})
}

func mustReadWebFile(name string) []byte {
	data, err := webFiles.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}
