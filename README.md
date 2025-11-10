# Stocky - Stock Rewards System

Stocky is a system where users can earn shares of Indian stocks (e.g., Reliance, TCS, Infosys) as incentives for actions like onboarding, referrals, or trading milestones.

## Features

- **Reward Management**: Record stock rewards for users with full audit trail
- **Double-Entry Ledger**: Track stock units, INR cash outflow, and company-incurred fees (brokerage, STT, GST, etc.)
- **Real-time Portfolio Valuation**: Hourly price updates to calculate current INR value of holdings
- **Historical Tracking**: Track historical INR values for all past days
- **Statistics**: Get total shares rewarded today and current portfolio value

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Logging**: Logrus
- **Database**: PostgreSQL
- **ORM**: Native database/sql with manual migrations

## Project Structure

```
.
├── main.go                          # Application entry point
├── go.mod                           # Go module dependencies
├── internal/
│   ├── config/                      # Configuration management
│   ├── database/                    # Database connection and migrations
│   ├── handlers/                    # HTTP request handlers
│   ├── models/                      # Data models
│   └── services/                    # Business logic services
├── README.md                        # This file
├── .env.example                     # Environment variables template
└── Stocky.postman_collection.json   # Postman collection for API testing
```

## Database Schema

### Tables

1. **reward_events**: Stores all reward events
   - `id` (UUID): Primary key
   - `user_id` (VARCHAR): User identifier
   - `stock_symbol` (VARCHAR): Stock symbol (e.g., RELIANCE, TCS)
   - `quantity` (NUMERIC(18,6)): Number of shares (supports fractional)
   - `reward_timestamp` (TIMESTAMP): When the reward was given
   - `event_id` (VARCHAR): Unique event identifier for duplicate detection
   - `created_at`, `updated_at` (TIMESTAMP): Audit fields

2. **ledger_entries**: Double-entry ledger system
   - `id` (UUID): Primary key
   - `reward_event_id` (UUID): Foreign key to reward_events
   - `entry_type` (VARCHAR): STOCK_CREDIT, CASH_DEBIT, BROKERAGE_FEE, STT_FEE, GST_FEE, OTHER_FEE
   - `account_type` (VARCHAR): STOCK, CASH, FEES
   - `stock_symbol` (VARCHAR): NULL for cash/fee entries
   - `quantity` (NUMERIC(18,6)): NULL for cash/fee entries
   - `amount` (NUMERIC(18,4)): INR amount
   - `description` (TEXT): Entry description

3. **stock_prices**: Historical stock prices
   - `id` (UUID): Primary key
   - `stock_symbol` (VARCHAR): Stock symbol
   - `price` (NUMERIC(18,4)): Price in INR
   - `price_timestamp` (TIMESTAMP): When the price was recorded
   - Unique constraint on (stock_symbol, price_timestamp)

4. **portfolio_snapshots**: Daily snapshots of user holdings
   - `id` (UUID): Primary key
   - `user_id` (VARCHAR): User identifier
   - `snapshot_date` (DATE): Date of snapshot
   - `stock_symbol` (VARCHAR): Stock symbol
   - `total_quantity` (NUMERIC(18,6)): Total shares held
   - `price_per_unit` (NUMERIC(18,4)): Price at snapshot time
   - `total_inr_value` (NUMERIC(18,4)): Total value in INR
   - Unique constraint on (user_id, snapshot_date, stock_symbol)

## API Endpoints

### 1. POST /api/v1/reward
Record that a user has been rewarded X shares of a stock.

**Request Body:**
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",  // Optional, defaults to now
  "event_id": "event-123"  // Optional, auto-generated if not provided
}
```

**Response:** 201 Created
```json
{
  "id": "uuid",
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-123"
}
```

### 2. GET /api/v1/today-stocks/:userId
Return all stock rewards for the user for today.

**Response:** 200 OK
```json
[
  {
    "stock_symbol": "RELIANCE",
    "quantity": "10.5",
    "rewarded_at": "2024-01-15T10:30:00Z"
  },
  {
    "stock_symbol": "TCS",
    "quantity": "5.25",
    "rewarded_at": "2024-01-15T11:00:00Z"
  }
]
```

### 3. GET /api/v1/historical-inr/:userId
Return the INR value of the user's stock rewards for all past days (up to yesterday).

**Response:** 200 OK
```json
[
  {
    "date": "2024-01-14",
    "value": "125000.5000"
  },
  {
    "date": "2024-01-13",
    "value": "120000.2500"
  }
]
```

### 4. GET /api/v1/stats/:userId
Return statistics for the user.

**Response:** 200 OK
```json
{
  "total_shares_today": {
    "RELIANCE": "10.5",
    "TCS": "5.25"
  },
  "current_portfolio_value": "150000.7500"
}
```

### 5. GET /api/v1/portfolio/:userId (Bonus)
Return holdings per stock symbol with current INR value.

**Response:** 200 OK
```json
{
  "holdings": [
    {
      "stock_symbol": "RELIANCE",
      "quantity": "50.5",
      "current_price": "2450.5000",
      "current_value": "123750.2500"
    },
    {
      "stock_symbol": "TCS",
      "quantity": "25.25",
      "current_price": "3500.0000",
      "current_value": "88375.0000"
    }
  ],
  "total_value": "212125.2500"
}
```

### 6. GET /health
Health check endpoint.

**Response:** 200 OK
```json
{
  "status": "ok"
}
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make sure PostgreSQL is running and accessible

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Assignment
   ```

2. **Create PostgreSQL database**
   ```sql
   CREATE DATABASE assignment;
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on port 8080 (or the port specified in your `.env` file).

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/assignment?sslmode=disable
```

## Edge Cases Handled

### 1. Duplicate Reward Events / Replay Attacks
- Each reward event has a unique `event_id` field
- The system checks for duplicate `event_id` before creating a reward
- Returns HTTP 409 Conflict if duplicate is detected

### 2. Stock Splits, Mergers, or Delisting
- The system tracks stock symbols and quantities separately
- Historical data is preserved in `portfolio_snapshots`
- For splits/mergers, manual adjustments can be made via ledger entries
- Delisted stocks can be handled by setting price to 0 or marking as inactive

### 3. Rounding Errors in INR Valuation
- Uses `NUMERIC(18,4)` for INR amounts to maintain precision
- Uses `NUMERIC(18,6)` for stock quantities to support fractional shares
- All calculations are done with proper decimal precision

### 4. Price API Downtime or Stale Data
- The system caches the latest price for each stock
- If price API is down, uses the last known price from database
- If no price exists, uses a default price (logged as warning)
- Hourly job retries price updates automatically

### 5. Adjustments/Refunds of Previously Given Rewards
- Can be handled by creating negative reward events
- Or by creating adjustment ledger entries
- Historical snapshots preserve the state at each point in time

## Background Jobs

### Hourly Price Update Job
- Runs every hour to fetch and update stock prices
- Updates `stock_prices` table with latest prices
- Creates daily `portfolio_snapshots` for yesterday's data
- Handles errors gracefully and logs them

## Scaling Considerations

1. **Database Indexing**: All frequently queried columns are indexed
2. **Connection Pooling**: Database connection pool is configured
3. **Caching**: Stock prices are cached in memory to reduce DB queries
4. **Background Jobs**: Price updates run asynchronously
5. **Transaction Safety**: All reward creation uses database transactions
6. **Error Handling**: Comprehensive error handling and logging

## Testing

Use the provided Postman collection (`Stocky.postman_collection.json`) to test all endpoints.

1. Import the collection into Postman
2. Update the base URL if needed
3. Test each endpoint with sample data

## License

This project is part of an assignment.

