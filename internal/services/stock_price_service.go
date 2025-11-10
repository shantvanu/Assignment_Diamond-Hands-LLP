package services

import (
	"database/sql"
	"fmt"
	"math/rand"
	"stocky/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

type StockPriceService struct {
	logger *logrus.Logger
	// Cache for prices to avoid frequent DB queries
	priceCache map[string]string
}

func NewStockPriceService(logger *logrus.Logger) *StockPriceService {
	return &StockPriceService{
		logger:     logger,
		priceCache: make(map[string]string),
	}
}

// GetCurrentPrice returns the current price for a stock symbol
// This is a hypothetical service that returns random prices
func (s *StockPriceService) GetCurrentPrice(stockSymbol string) (string, error) {
	// Check cache first
	if price, ok := s.priceCache[stockSymbol]; ok {
		return price, nil
	}

	// Generate a random price between 100 and 5000 INR
	// In production, this would fetch from NSE/BSE API
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	basePrice := 100.0 + r.Float64()*4900.0
	price := fmt.Sprintf("%.4f", basePrice)

	// Cache the price
	s.priceCache[stockSymbol] = price

	s.logger.WithFields(logrus.Fields{
		"stock_symbol": stockSymbol,
		"price": price,
	}).Debug("Generated stock price")

	return price, nil
}

// UpdatePrices fetches and stores latest prices for all stocks
func (s *StockPriceService) UpdatePrices(db *sql.DB) error {
	// Get all unique stock symbols from reward events
	rows, err := db.Query(`
		SELECT DISTINCT stock_symbol 
		FROM reward_events
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return err
		}
		symbols = append(symbols, symbol)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// Update prices for each symbol
	now := time.Now()
	for _, symbol := range symbols {
		price, err := s.GetCurrentPrice(symbol)
		if err != nil {
			s.logger.WithError(err).WithField("stock_symbol", symbol).Error("Failed to get price")
			continue
		}

		// Store price in database
		_, err = db.Exec(`
			INSERT INTO stock_prices (stock_symbol, price, price_timestamp)
			VALUES ($1, $2, $3)
			ON CONFLICT (stock_symbol, price_timestamp) DO UPDATE
			SET price = EXCLUDED.price
		`, symbol, price, now)
		if err != nil {
			s.logger.WithError(err).WithField("stock_symbol", symbol).Error("Failed to store price")
			continue
		}

		// Update cache
		s.priceCache[symbol] = price
	}

	s.logger.WithField("count", len(symbols)).Info("Updated stock prices")
	return nil
}

// GetLatestPrice gets the latest price from database
func (s *StockPriceService) GetLatestPrice(db *sql.DB, stockSymbol string) (string, error) {
	var price string
	err := db.QueryRow(`
		SELECT price FROM stock_prices 
		WHERE stock_symbol = $1 
		ORDER BY price_timestamp DESC 
		LIMIT 1
	`, stockSymbol).Scan(&price)

	if err == sql.ErrNoRows {
		// No price in DB, generate one
		return s.GetCurrentPrice(stockSymbol)
	}

	return price, err
}

