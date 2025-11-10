package models

import (
	"database/sql"
	"time"
)

// RewardEvent represents a stock reward event
type RewardEvent struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	StockSymbol    string    `json:"stock_symbol"`
	Quantity       string    `json:"quantity"` // NUMERIC as string for precision
	RewardTimestamp time.Time `json:"reward_timestamp"`
	EventID        string    `json:"event_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// LedgerEntry represents a double-entry ledger entry
type LedgerEntry struct {
	ID            string         `json:"id"`
	RewardEventID sql.NullString `json:"reward_event_id"`
	EntryType     string         `json:"entry_type"` // STOCK_CREDIT, CASH_DEBIT, etc.
	AccountType   string         `json:"account_type"` // STOCK, CASH, FEES
	StockSymbol   sql.NullString `json:"stock_symbol"`
	Quantity      sql.NullString `json:"quantity"`
	Amount        string         `json:"amount"` // NUMERIC as string
	Description   sql.NullString `json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
}

// StockPrice represents a stock price at a point in time
type StockPrice struct {
	ID            string    `json:"id"`
	StockSymbol   string    `json:"stock_symbol"`
	Price         string    `json:"price"`
	PriceTimestamp time.Time `json:"price_timestamp"`
	CreatedAt     time.Time `json:"created_at"`
}

// PortfolioSnapshot represents a daily snapshot of user holdings
type PortfolioSnapshot struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	SnapshotDate  time.Time `json:"snapshot_date"`
	StockSymbol   string    `json:"stock_symbol"`
	TotalQuantity string    `json:"total_quantity"`
	PricePerUnit  string    `json:"price_per_unit"`
	TotalINRValue string    `json:"total_inr_value"`
	CreatedAt     time.Time `json:"created_at"`
}

// TodayStock represents today's stock reward
type TodayStock struct {
	StockSymbol string `json:"stock_symbol"`
	Quantity    string `json:"quantity"`
	RewardedAt  string `json:"rewarded_at"`
}

// HistoricalINR represents historical INR value
type HistoricalINR struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// Stats represents user statistics
type Stats struct {
	TotalSharesToday map[string]string `json:"total_shares_today"` // stock_symbol -> quantity
	CurrentPortfolioValue string        `json:"current_portfolio_value"`
}

// Portfolio represents user portfolio
type Portfolio struct {
	Holdings []Holding `json:"holdings"`
	TotalValue string  `json:"total_value"`
}

// Holding represents a single stock holding
type Holding struct {
	StockSymbol string `json:"stock_symbol"`
	Quantity    string `json:"quantity"`
	CurrentPrice string `json:"current_price"`
	CurrentValue string `json:"current_value"`
}

