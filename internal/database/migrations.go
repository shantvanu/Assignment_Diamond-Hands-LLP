package database

import (
	"database/sql"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createRewardEventsTable,
		createLedgerEntriesTable,
		createStockPricesTable,
		createPortfolioSnapshotsTable,
		createIndexes,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

const createRewardEventsTable = `
CREATE TABLE IF NOT EXISTS reward_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    stock_symbol VARCHAR(50) NOT NULL,
    quantity NUMERIC(18, 6) NOT NULL,
    reward_timestamp TIMESTAMP NOT NULL,
    event_id VARCHAR(255) UNIQUE NOT NULL, -- For duplicate detection
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createLedgerEntriesTable = `
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reward_event_id UUID REFERENCES reward_events(id) ON DELETE CASCADE,
    entry_type VARCHAR(50) NOT NULL, -- 'STOCK_CREDIT', 'CASH_DEBIT', 'BROKERAGE_FEE', 'STT_FEE', 'GST_FEE', 'OTHER_FEE'
    account_type VARCHAR(50) NOT NULL, -- 'STOCK', 'CASH', 'FEES'
    stock_symbol VARCHAR(50), -- NULL for cash/fee entries
    quantity NUMERIC(18, 6), -- NULL for cash/fee entries
    amount NUMERIC(18, 4) NOT NULL, -- INR amount
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createStockPricesTable = `
CREATE TABLE IF NOT EXISTS stock_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_symbol VARCHAR(50) NOT NULL,
    price NUMERIC(18, 4) NOT NULL,
    price_timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(stock_symbol, price_timestamp)
);
`

const createPortfolioSnapshotsTable = `
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    snapshot_date DATE NOT NULL,
    stock_symbol VARCHAR(50) NOT NULL,
    total_quantity NUMERIC(18, 6) NOT NULL,
    price_per_unit NUMERIC(18, 4) NOT NULL,
    total_inr_value NUMERIC(18, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, snapshot_date, stock_symbol)
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_reward_events_user_id ON reward_events(user_id);
CREATE INDEX IF NOT EXISTS idx_reward_events_timestamp ON reward_events(reward_timestamp);
CREATE INDEX IF NOT EXISTS idx_reward_events_event_id ON reward_events(event_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_reward_event_id ON ledger_entries(reward_event_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_type ON ledger_entries(account_type);
CREATE INDEX IF NOT EXISTS idx_stock_prices_symbol_timestamp ON stock_prices(stock_symbol, price_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_user_date ON portfolio_snapshots(user_id, snapshot_date DESC);
`

