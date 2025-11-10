# API Specification

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Health Check

**GET** `/health`

Check if the API is running.

**Response:**
```json
{
  "status": "ok"
}
```

**Status Code:** 200 OK

---

### 2. Create Reward

**POST** `/api/v1/reward`

Record that a user has been rewarded X shares of a stock.

**Request Body:**
```json
{
  "user_id": "string (required)",
  "stock_symbol": "string (required)",
  "quantity": "string (required, numeric)",
  "reward_timestamp": "string (optional, RFC3339 format)",
  "event_id": "string (optional, auto-generated if not provided)"
}
```

**Example Request:**
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-123"
}
```

**Response:** 201 Created
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-123"
}
```

**Error Responses:**

- **400 Bad Request:** Invalid request body
  ```json
  {
    "error": "Key: 'CreateRewardRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag"
  }
  ```

- **409 Conflict:** Duplicate event_id
  ```json
  {
    "error": "duplicate reward event"
  }
  ```

- **500 Internal Server Error:** Server error
  ```json
  {
    "error": "failed to create reward"
  }
  ```

---

### 3. Get Today's Stocks

**GET** `/api/v1/today-stocks/:userId`

Return all stock rewards for the user for today.

**Path Parameters:**
- `userId` (string, required): User identifier

**Example Request:**
```
GET /api/v1/today-stocks/user123
```

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

**Empty Response:** 200 OK
```json
[]
```

**Error Responses:**

- **400 Bad Request:** Missing user_id
  ```json
  {
    "error": "user_id is required"
  }
  ```

- **500 Internal Server Error:** Server error
  ```json
  {
    "error": "failed to get today's stocks"
  }
  ```

---

### 4. Get Historical INR

**GET** `/api/v1/historical-inr/:userId`

Return the INR value of the user's stock rewards for all past days (up to yesterday).

**Path Parameters:**
- `userId` (string, required): User identifier

**Example Request:**
```
GET /api/v1/historical-inr/user123
```

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
  },
  {
    "date": "2024-01-12",
    "value": "115000.7500"
  }
]
```

**Empty Response:** 200 OK
```json
[]
```

**Error Responses:**

- **400 Bad Request:** Missing user_id
  ```json
  {
    "error": "user_id is required"
  }
  ```

- **500 Internal Server Error:** Server error
  ```json
  {
    "error": "failed to get historical INR"
  }
  ```

---

### 5. Get Stats

**GET** `/api/v1/stats/:userId`

Return statistics for the user:
- Total shares rewarded today (grouped by stock symbol)
- Current INR value of the user's portfolio

**Path Parameters:**
- `userId` (string, required): User identifier

**Example Request:**
```
GET /api/v1/stats/user123
```

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

**Example Response (No rewards today):**
```json
{
  "total_shares_today": {},
  "current_portfolio_value": "150000.7500"
}
```

**Error Responses:**

- **400 Bad Request:** Missing user_id
  ```json
  {
    "error": "user_id is required"
  }
  ```

- **500 Internal Server Error:** Server error
  ```json
  {
    "error": "failed to get stats"
  }
  ```

---

### 6. Get Portfolio (Bonus)

**GET** `/api/v1/portfolio/:userId`

Return holdings per stock symbol with current INR value.

**Path Parameters:**
- `userId` (string, required): User identifier

**Example Request:**
```
GET /api/v1/portfolio/user123
```

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

**Empty Response:** 200 OK
```json
{
  "holdings": [],
  "total_value": "0.0000"
}
```

**Error Responses:**

- **400 Bad Request:** Missing user_id
  ```json
  {
    "error": "user_id is required"
  }
  ```

- **500 Internal Server Error:** Server error
  ```json
  {
    "error": "failed to get portfolio"
  }
  ```

---

## Data Types

### Stock Symbol
- Type: String
- Format: Uppercase stock symbols (e.g., "RELIANCE", "TCS", "INFOSYS")
- Length: Up to 50 characters

### Quantity
- Type: String (numeric)
- Format: Decimal number with up to 6 decimal places
- Example: "10.5", "0.123456", "1000.0"
- Supports fractional shares

### INR Amount
- Type: String (numeric)
- Format: Decimal number with up to 4 decimal places
- Example: "1000.50", "123456.7890"
- Precision: NUMERIC(18,4)

### Timestamp
- Type: String
- Format: RFC3339 (ISO 8601)
- Example: "2024-01-15T10:30:00Z"
- Timezone: UTC

### Date
- Type: String
- Format: YYYY-MM-DD
- Example: "2024-01-15"

---

## Error Handling

All endpoints follow consistent error handling:

1. **400 Bad Request:** Invalid input or missing required fields
2. **409 Conflict:** Duplicate resource (e.g., duplicate event_id)
3. **500 Internal Server Error:** Server-side error

Error responses include an `error` field with a descriptive message.

---

## Rate Limiting

Currently, there is no rate limiting implemented. For production, consider:
- Per-user rate limiting
- Per-IP rate limiting
- API key-based rate limiting

---

## Authentication

Currently, there is no authentication implemented. For production, consider:
- API key authentication
- JWT tokens
- OAuth 2.0

---

## Testing

Use the provided Postman collection (`Stocky.postman_collection.json`) to test all endpoints.

1. Import the collection into Postman
2. Set the `base_url` variable (default: `http://localhost:8080`)
3. Test each endpoint with sample data

---

## Example Workflow

1. **Create a reward:**
   ```bash
   POST /api/v1/reward
   {
     "user_id": "user123",
     "stock_symbol": "RELIANCE",
     "quantity": "10.5",
     "event_id": "event-001"
   }
   ```

2. **Check today's stocks:**
   ```bash
   GET /api/v1/today-stocks/user123
   ```

3. **Get portfolio value:**
   ```bash
   GET /api/v1/stats/user123
   ```

4. **View full portfolio:**
   ```bash
   GET /api/v1/portfolio/user123
   ```

5. **Check historical values:**
   ```bash
   GET /api/v1/historical-inr/user123
   ```

