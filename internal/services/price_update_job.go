package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

// StartPriceUpdateJob starts a background job that updates stock prices every hour
func StartPriceUpdateJob(ctx context.Context, db *sql.DB, stockPriceService *StockPriceService, logger *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start
	updatePrices(db, stockPriceService, logger)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Price update job stopped")
			return
		case <-ticker.C:
			updatePrices(db, stockPriceService, logger)
		}
	}
}

func updatePrices(db *sql.DB, stockPriceService *StockPriceService, logger *logrus.Logger) {
	logger.Info("Starting hourly price update")

	// Update stock prices
	if err := stockPriceService.UpdatePrices(db); err != nil {
		logger.WithError(err).Error("Failed to update stock prices")
		return
	}

	// Update portfolio snapshots (for yesterday's data)
	portfolioService := NewPortfolioService(db, stockPriceService, logger)
	if err := portfolioService.UpdatePortfolioSnapshots(); err != nil {
		logger.WithError(err).Error("Failed to update portfolio snapshots")
		return
	}

	logger.Info("Hourly price update completed")
}

