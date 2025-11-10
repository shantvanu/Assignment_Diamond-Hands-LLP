package services

import (
	"database/sql"
	"errors"
	"fmt"
	"stocky/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrDuplicateEvent = errors.New("duplicate reward event")
)

type RewardService struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewRewardService(db *sql.DB, logger *logrus.Logger) *RewardService {
	return &RewardService{
		db:     db,
		logger: logger,
	}
}

// CreateReward creates a reward event and corresponding ledger entries
func (s *RewardService) CreateReward(userID, stockSymbol, quantity, eventID string, rewardTimestamp time.Time) (*models.RewardEvent, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check for duplicate event_id
	var existingID string
	err = tx.QueryRow("SELECT id FROM reward_events WHERE event_id = $1", eventID).Scan(&existingID)
	if err == nil {
		// Duplicate found
		s.logger.WithFields(logrus.Fields{
			"event_id": eventID,
			"existing_id": existingID,
		}).Warn("Duplicate reward event detected")
		return nil, ErrDuplicateEvent
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Insert reward event
	rewardID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO reward_events (id, user_id, stock_symbol, quantity, reward_timestamp, event_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rewardID, userID, stockSymbol, quantity, rewardTimestamp, eventID)
	if err != nil {
		return nil, err
	}

	// Get current stock price for calculations
	var currentPrice string
	err = tx.QueryRow(`
		SELECT price FROM stock_prices 
		WHERE stock_symbol = $1 
		ORDER BY price_timestamp DESC 
		LIMIT 1
	`, stockSymbol).Scan(&currentPrice)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// If no price found, use a default price (in production, this would fetch from API)
	if err == sql.ErrNoRows {
		currentPrice = "100.0000" // Default price
		s.logger.WithField("stock_symbol", stockSymbol).Warn("No price found, using default")
	}

	// Calculate fees (hypothetical values)
	// In production, these would be calculated based on actual NSE/BSE rates
	brokerage := calculateBrokerage(quantity, currentPrice)
	stt := calculateSTT(quantity, currentPrice)
	gst := calculateGST(brokerage)
	otherFees := calculateOtherFees(quantity, currentPrice)
	totalFees := addAmounts(brokerage, stt, gst, otherFees)

	// Calculate total cash outflow
	totalCashOutflow := addAmounts(multiplyAmounts(quantity, currentPrice), totalFees)

	// Create ledger entries
	ledgerEntries := []struct {
		entryType   string
		accountType string
		stockSymbol sql.NullString
		quantity    sql.NullString
		amount      string
		description string
	}{
		{"STOCK_CREDIT", "STOCK", sql.NullString{String: stockSymbol, Valid: true}, sql.NullString{String: quantity, Valid: true}, "0.0000", "Stock reward credit"},
		{"CASH_DEBIT", "CASH", sql.NullString{Valid: false}, sql.NullString{Valid: false}, totalCashOutflow, "Cash outflow for stock purchase"},
		{"BROKERAGE_FEE", "FEES", sql.NullString{Valid: false}, sql.NullString{Valid: false}, brokerage, "Brokerage fee"},
		{"STT_FEE", "FEES", sql.NullString{Valid: false}, sql.NullString{Valid: false}, stt, "Securities Transaction Tax"},
		{"GST_FEE", "FEES", sql.NullString{Valid: false}, sql.NullString{Valid: false}, gst, "GST on brokerage"},
		{"OTHER_FEE", "FEES", sql.NullString{Valid: false}, sql.NullString{Valid: false}, otherFees, "Other regulatory fees"},
	}

	for _, entry := range ledgerEntries {
		_, err = tx.Exec(`
			INSERT INTO ledger_entries (id, reward_event_id, entry_type, account_type, stock_symbol, quantity, amount, description)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
		`, rewardID, entry.entryType, entry.accountType, entry.stockSymbol, entry.quantity, entry.amount, entry.description)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"reward_id": rewardID,
		"user_id": userID,
		"stock_symbol": stockSymbol,
		"quantity": quantity,
	}).Info("Reward created successfully")

	return &models.RewardEvent{
		ID:             rewardID,
		UserID:         userID,
		StockSymbol:    stockSymbol,
		Quantity:       quantity,
		RewardTimestamp: rewardTimestamp,
		EventID:        eventID,
	}, nil
}

// GetTodayStocks returns all stock rewards for a user today
func (s *RewardService) GetTodayStocks(userID string) ([]models.TodayStock, error) {
	rows, err := s.db.Query(`
		SELECT stock_symbol, quantity, reward_timestamp
		FROM reward_events
		WHERE user_id = $1 
		AND DATE(reward_timestamp) = CURRENT_DATE
		ORDER BY reward_timestamp DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.TodayStock
	for rows.Next() {
		var stock models.TodayStock
		var timestamp time.Time
		if err := rows.Scan(&stock.StockSymbol, &stock.Quantity, &timestamp); err != nil {
			return nil, err
		}
		stock.RewardedAt = timestamp.Format(time.RFC3339)
		stocks = append(stocks, stock)
	}

	return stocks, rows.Err()
}

// Helper functions for fee calculations
func calculateBrokerage(quantity, price string) string {
	// Hypothetical: 0.05% of transaction value
	value := multiplyAmounts(quantity, price)
	return divideAmount(value, "10000") // 0.05% = divide by 10000
}

func calculateSTT(quantity, price string) string {
	// Hypothetical: 0.1% of transaction value
	value := multiplyAmounts(quantity, price)
	return divideAmount(value, "1000") // 0.1% = divide by 1000
}

func calculateGST(brokerage string) string {
	// Hypothetical: 18% GST on brokerage
	// GST = brokerage * 0.18
	return multiplyAmounts(brokerage, "0.18")
}

func calculateOtherFees(quantity, price string) string {
	// Hypothetical: Fixed 10 INR per transaction
	return "10.0000"
}

// Amount arithmetic helpers (using string to maintain precision)
func addAmounts(amounts ...string) string {
	// Simple addition - in production, use proper decimal library
	// For now, we'll use a simple approach
	var total float64
	for _, amt := range amounts {
		var val float64
		_, _ = fmt.Sscanf(amt, "%f", &val)
		total += val
	}
	return fmt.Sprintf("%.4f", total)
}

func multiplyAmounts(a, b string) string {
	var valA, valB float64
	_, _ = fmt.Sscanf(a, "%f", &valA)
	_, _ = fmt.Sscanf(b, "%f", &valB)
	return fmt.Sprintf("%.4f", valA*valB)
}

func divideAmount(a, b string) string {
	var valA, valB float64
	_, _ = fmt.Sscanf(a, "%f", &valA)
	_, _ = fmt.Sscanf(b, "%f", &valB)
	if valB == 0 {
		return "0.0000"
	}
	return fmt.Sprintf("%.4f", valA/valB)
}

