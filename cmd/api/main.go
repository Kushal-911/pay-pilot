package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/payments"
	"paypilot/internal/paypal"
)

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
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
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
	api.POST("/orders", payments.CreateOrder(db))
	api.GET("/orders", payments.ListOrders(db))
	api.GET("/orders/:id", payments.GetOrder(db))
	api.POST("/orders/:id/paypal-order", payments.CreatePayPalOrder(db, paypalClient))
	api.POST("/payments/:id/capture", payments.CapturePayPalOrder(db, paypalClient))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
