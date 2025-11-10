package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"stocky/internal/config"
	"stocky/internal/database"
	"stocky/internal/handlers"
	"stocky/internal/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
	}

	// Initialize services
	rewardService := services.NewRewardService(db, logger)
	stockPriceService := services.NewStockPriceService(logger)
	portfolioService := services.NewPortfolioService(db, stockPriceService, logger)

	// Start hourly price update job
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go services.StartPriceUpdateJob(ctx, db, stockPriceService, logger)

	// Initialize handlers
	rewardHandler := handlers.NewRewardHandler(rewardService, logger)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService, logger)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/reward", rewardHandler.CreateReward)
		api.GET("/today-stocks/:userId", rewardHandler.GetTodayStocks)
		api.GET("/historical-inr/:userId", portfolioHandler.GetHistoricalINR)
		api.GET("/stats/:userId", portfolioHandler.GetStats)
		api.GET("/portfolio/:userId", portfolioHandler.GetPortfolio) // Bonus endpoint
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	logger.WithField("port", cfg.Port).Info("Server started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	cancel() // Stop price update job

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

