# Database Schema Documentation

## Overview

The Stocky system uses PostgreSQL with a double-entry ledger system to track stock rewards, cash flows, and fees.

## Tables

### 1. reward_events

Stores all stock reward events given to users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier for the reward event |
| user_id | VARCHAR(255) | NOT NULL | User identifier |
| stock_symbol | VARCHAR(50) | NOT NULL | Stock symbol (e.g., RELIANCE, TCS, INFOSYS) |
| quantity | NUMERIC(18,6) | NOT NULL | Number of shares (supports fractional shares) |
| reward_timestamp | TIMESTAMP | NOT NULL | When the reward was given |
| event_id | VARCHAR(255) | UNIQUE, NOT NULL | Unique event identifier for duplicate detection |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record update timestamp |

**Indexes:**
- `idx_reward_events_user_id` on `user_id`
- `idx_reward_events_timestamp` on `reward_timestamp`
- `idx_reward_events_event_id` on `event_id` (unique)

**Relationships:**
- One-to-many with `ledger_entries` (via `reward_event_id`)

### 2. ledger_entries

Double-entry ledger system tracking all financial transactions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier for the ledger entry |
| reward_event_id | UUID | FOREIGN KEY | Reference to reward_events.id (nullable) |
| entry_type | VARCHAR(50) | NOT NULL | Type: STOCK_CREDIT, CASH_DEBIT, BROKERAGE_FEE, STT_FEE, GST_FEE, OTHER_FEE |
| account_type | VARCHAR(50) | NOT NULL | Account: STOCK, CASH, FEES |
| stock_symbol | VARCHAR(50) | NULL | Stock symbol (NULL for cash/fee entries) |
| quantity | NUMERIC(18,6) | NULL | Number of shares (NULL for cash/fee entries) |
| amount | NUMERIC(18,4) | NOT NULL | INR amount |
| description | TEXT | NULL | Entry description |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |

**Indexes:**
- `idx_ledger_entries_reward_event_id` on `reward_event_id`
- `idx_ledger_entries_account_type` on `account_type`

**Entry Types:**
- `STOCK_CREDIT`: Stock units credited to user
- `CASH_DEBIT`: Cash outflow for stock purchase
- `BROKERAGE_FEE`: Brokerage fee charged
- `STT_FEE`: Securities Transaction Tax
- `GST_FEE`: GST on brokerage
- `OTHER_FEE`: Other regulatory fees

**Account Types:**
- `STOCK`: Stock holdings account
- `CASH`: Cash account
- `FEES`: Fees account

### 3. stock_prices

Historical stock prices for valuation.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier |
| stock_symbol | VARCHAR(50) | NOT NULL | Stock symbol |
| price | NUMERIC(18,4) | NOT NULL | Price in INR |
| price_timestamp | TIMESTAMP | NOT NULL | When the price was recorded |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |

**Indexes:**
- `idx_stock_prices_symbol_timestamp` on `(stock_symbol, price_timestamp DESC)`

**Unique Constraint:**
- `(stock_symbol, price_timestamp)` - Prevents duplicate price entries for same timestamp

### 4. portfolio_snapshots

Daily snapshots of user holdings for historical tracking.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier |
| user_id | VARCHAR(255) | NOT NULL | User identifier |
| snapshot_date | DATE | NOT NULL | Date of snapshot |
| stock_symbol | VARCHAR(50) | NOT NULL | Stock symbol |
| total_quantity | NUMERIC(18,6) | NOT NULL | Total shares held |
| price_per_unit | NUMERIC(18,4) | NOT NULL | Price at snapshot time |
| total_inr_value | NUMERIC(18,4) | NOT NULL | Total value in INR |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |

**Indexes:**
- `idx_portfolio_snapshots_user_date` on `(user_id, snapshot_date DESC)`

**Unique Constraint:**
- `(user_id, snapshot_date, stock_symbol)` - One snapshot per user per stock per day

## Data Types

### NUMERIC Precision

- **Stock Quantities**: `NUMERIC(18,6)` - Supports up to 12 integer digits and 6 decimal places (fractional shares)
- **INR Amounts**: `NUMERIC(18,4)` - Supports up to 14 integer digits and 4 decimal places (precise currency)

### Example Values

- Stock quantity: `123456789012.123456` (max value)
- INR amount: `12345678901234.1234` (max value)

## Relationships

```
reward_events (1) ----< (many) ledger_entries
```

Each reward event creates multiple ledger entries:
1. STOCK_CREDIT entry
2. CASH_DEBIT entry
3. BROKERAGE_FEE entry
4. STT_FEE entry
5. GST_FEE entry
6. OTHER_FEE entry

## Query Patterns

### Get user's current holdings
```sql
SELECT stock_symbol, SUM(quantity) as total_quantity
FROM reward_events
WHERE user_id = $1
GROUP BY stock_symbol;
```

### Get latest price for a stock
```sql
SELECT price FROM stock_prices
WHERE stock_symbol = $1
ORDER BY price_timestamp DESC
LIMIT 1;
```

### Get historical portfolio value
```sql
SELECT snapshot_date, SUM(total_inr_value) as daily_value
FROM portfolio_snapshots
WHERE user_id = $1
AND snapshot_date < CURRENT_DATE
GROUP BY snapshot_date
ORDER BY snapshot_date DESC;
```

### Get all ledger entries for a reward
```sql
SELECT * FROM ledger_entries
WHERE reward_event_id = $1
ORDER BY created_at;
```

## Migration Notes

All migrations are run automatically on application startup via `database.RunMigrations()`.

The migrations use `CREATE TABLE IF NOT EXISTS` to be idempotent, allowing safe re-runs.

