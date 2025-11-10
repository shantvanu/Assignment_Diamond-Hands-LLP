package services

import (
	"database/sql"
	"fmt"
	"stocky/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

type PortfolioService struct {
	db                *sql.DB
	stockPriceService *StockPriceService
	logger            *logrus.Logger
}

func NewPortfolioService(db *sql.DB, stockPriceService *StockPriceService, logger *logrus.Logger) *PortfolioService {
	return &PortfolioService{
		db:                db,
		stockPriceService: stockPriceService,
		logger:            logger,
	}
}

// GetHistoricalINR returns the INR value of user's stock rewards for all past days
func (s *PortfolioService) GetHistoricalINR(userID string) ([]models.HistoricalINR, error) {
	rows, err := s.db.Query(`
		SELECT snapshot_date, SUM(total_inr_value) as daily_value
		FROM portfolio_snapshots
		WHERE user_id = $1 
		AND snapshot_date < CURRENT_DATE
		GROUP BY snapshot_date
		ORDER BY snapshot_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historical []models.HistoricalINR
	for rows.Next() {
		var h models.HistoricalINR
		var date time.Time
		var value string
		if err := rows.Scan(&date, &value); err != nil {
			return nil, err
		}
		h.Date = date.Format("2006-01-02")
		h.Value = value
		historical = append(historical, h)
	}

	return historical, rows.Err()
}

// GetStats returns user statistics
func (s *PortfolioService) GetStats(userID string) (*models.Stats, error) {
	stats := &models.Stats{
		TotalSharesToday: make(map[string]string),
	}

	// Get total shares rewarded today (grouped by stock symbol)
	rows, err := s.db.Query(`
		SELECT stock_symbol, SUM(quantity) as total_quantity
		FROM reward_events
		WHERE user_id = $1 
		AND DATE(reward_timestamp) = CURRENT_DATE
		GROUP BY stock_symbol
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var symbol, quantity string
		if err := rows.Scan(&symbol, &quantity); err != nil {
			return nil, err
		}
		stats.TotalSharesToday[symbol] = quantity
	}
	rows.Close()

	// Get current portfolio value
	currentValue, err := s.calculateCurrentPortfolioValue(userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentPortfolioValue = currentValue

	return stats, nil
}

// GetPortfolio returns user portfolio with holdings per stock
func (s *PortfolioService) GetPortfolio(userID string) (*models.Portfolio, error) {
	// Get all holdings for user
	rows, err := s.db.Query(`
		SELECT stock_symbol, SUM(quantity) as total_quantity
		FROM reward_events
		WHERE user_id = $1
		GROUP BY stock_symbol
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	portfolio := &models.Portfolio{
		Holdings: []models.Holding{},
	}

	var totalValue float64
	for rows.Next() {
		var symbol, quantity string
		if err := rows.Scan(&symbol, &quantity); err != nil {
			return nil, err
		}

		// Get current price
		price, err := s.stockPriceService.GetLatestPrice(s.db, symbol)
		if err != nil {
			s.logger.WithError(err).WithField("stock_symbol", symbol).Warn("Failed to get price")
			price = "0.0000"
		}

		// Calculate current value
		value := multiplyAmounts(quantity, price)

		holding := models.Holding{
			StockSymbol:  symbol,
			Quantity:     quantity,
			CurrentPrice: price,
			CurrentValue: value,
		}
		portfolio.Holdings = append(portfolio.Holdings, holding)

		var val float64
		_, _ = fmt.Sscanf(value, "%f", &val)
		totalValue += val
	}

	portfolio.TotalValue = fmt.Sprintf("%.4f", totalValue)
	return portfolio, rows.Err()
}

// calculateCurrentPortfolioValue calculates the current INR value of user's portfolio
func (s *PortfolioService) calculateCurrentPortfolioValue(userID string) (string, error) {
	// Get all holdings
	rows, err := s.db.Query(`
		SELECT stock_symbol, SUM(quantity) as total_quantity
		FROM reward_events
		WHERE user_id = $1
		GROUP BY stock_symbol
	`, userID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var totalValue float64
	for rows.Next() {
		var symbol, quantity string
		if err := rows.Scan(&symbol, &quantity); err != nil {
			return "", err
		}

		// Get current price
		price, err := s.stockPriceService.GetLatestPrice(s.db, symbol)
		if err != nil {
			s.logger.WithError(err).WithField("stock_symbol", symbol).Warn("Failed to get price, using 0")
			price = "0.0000"
		}

		// Calculate value
		value := multiplyAmounts(quantity, price)
		var val float64
		_, _ = fmt.Sscanf(value, "%f", &val)
		totalValue += val
	}

	return fmt.Sprintf("%.4f", totalValue), rows.Err()
}

// UpdatePortfolioSnapshots creates daily snapshots for all users
func (s *PortfolioService) UpdatePortfolioSnapshots() error {
	// Get all unique user IDs
	rows, err := s.db.Query(`SELECT DISTINCT user_id FROM reward_events`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// Get yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Create snapshots for each user
	for _, userID := range userIDs {
		// Get holdings per stock
		holdingsRows, err := s.db.Query(`
			SELECT stock_symbol, SUM(quantity) as total_quantity
			FROM reward_events
			WHERE user_id = $1
			AND DATE(reward_timestamp) <= $2
			GROUP BY stock_symbol
		`, userID, yesterday)
		if err != nil {
			s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get holdings")
			continue
		}

		for holdingsRows.Next() {
			var symbol, quantity string
			if err := holdingsRows.Scan(&symbol, &quantity); err != nil {
				holdingsRows.Close()
				return err
			}

			// Get price for yesterday (or latest available)
			var price string
			err = s.db.QueryRow(`
				SELECT price FROM stock_prices 
				WHERE stock_symbol = $1 
				AND DATE(price_timestamp) <= $2
				ORDER BY price_timestamp DESC 
				LIMIT 1
			`, symbol, yesterday).Scan(&price)

			if err == sql.ErrNoRows {
				// No price available, use 0
				price = "0.0000"
				s.logger.WithFields(logrus.Fields{
					"user_id": userID,
					"stock_symbol": symbol,
					"date": yesterday,
				}).Warn("No price found for snapshot")
			} else if err != nil {
				holdingsRows.Close()
				return err
			}

			// Calculate total value
			totalValue := multiplyAmounts(quantity, price)

			// Insert or update snapshot
			_, err = s.db.Exec(`
				INSERT INTO portfolio_snapshots (user_id, snapshot_date, stock_symbol, total_quantity, price_per_unit, total_inr_value)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (user_id, snapshot_date, stock_symbol) DO UPDATE
				SET total_quantity = EXCLUDED.total_quantity,
				    price_per_unit = EXCLUDED.price_per_unit,
				    total_inr_value = EXCLUDED.total_inr_value
			`, userID, yesterday, symbol, quantity, price, totalValue)
			if err != nil {
				holdingsRows.Close()
				return err
			}
		}
		holdingsRows.Close()
	}

	s.logger.Info("Portfolio snapshots updated")
	return nil
}

// Helper function (shared with reward_service)
func multiplyAmounts(a, b string) string {
	var valA, valB float64
	_, _ = fmt.Sscanf(a, "%f", &valA)
	_, _ = fmt.Sscanf(b, "%f", &valB)
	return fmt.Sprintf("%.4f", valA*valB)
}

