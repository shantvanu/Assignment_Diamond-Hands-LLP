# Postman Testing Guide - Stocky API

## Base URL
```
http://localhost:8080
```

## Environment Variables
Set these in Postman:
- `base_url`: `http://localhost:8080`

---

## 1. Health Check

**Method:** `GET`  
**URL:** `{{base_url}}/health`  
**Headers:** None  
**Body:** None

**Expected Response:**
```json
{
  "status": "ok"
}
```

**Status Code:** `200 OK`

---

## 2. Create Reward

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/reward`  
**Headers:**
```
Content-Type: application/json
```
 
 
**Body (JSON):**
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-001"
}
```

**Required Fields:**
- `user_id` (string, required)
- `stock_symbol` (string, required)
- `quantity` (string, required, numeric)

**Optional Fields:**
- `reward_timestamp` (string, RFC3339 format, defaults to current time)
- `event_id` (string, auto-generated if not provided)

**Example Request (Minimal):**
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5"
}
```

**Expected Response (Success):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-001"
}
```

**Status Code:** `201 Created`

**Error Responses:**

**400 Bad Request** (Missing required field):
```json
{
  "error": "Key: 'CreateRewardRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag"
}
```

**409 Conflict** (Duplicate event_id):
```json
{
  "error": "duplicate reward event"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to create reward"
}
```

---

## 3. Get Today's Stocks

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/today-stocks/:userId`  
**Headers:** None  
**Body:** None

**Path Parameters:**
- `userId` (string, required) - Replace `:userId` in URL

**Example URL:**
```
{{base_url}}/api/v1/today-stocks/user123
```

**Expected Response (Success):**
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

**Empty Response (No rewards today):**
```json
[]
```

**Status Code:** `200 OK`

**Error Responses:**

**400 Bad Request:**
```json
{
  "error": "user_id is required"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to get today's stocks"
}
```

---

## 4. Get Historical INR

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/historical-inr/:userId`  
**Headers:** None  
**Body:** None

**Path Parameters:**
- `userId` (string, required) - Replace `:userId` in URL

**Example URL:**
```
{{base_url}}/api/v1/historical-inr/user123
```

**Expected Response (Success):**
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

**Empty Response (No historical data):**
```json
[]
```

**Status Code:** `200 OK`

**Error Responses:**

**400 Bad Request:**
```json
{
  "error": "user_id is required"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to get historical INR"
}
```

---

## 5. Get Stats

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/stats/:userId`  
**Headers:** None  
**Body:** None

**Path Parameters:**
- `userId` (string, required) - Replace `:userId` in URL

**Example URL:**
```
{{base_url}}/api/v1/stats/user123
```

**Expected Response (Success):**
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

**Status Code:** `200 OK`

**Error Responses:**

**400 Bad Request:**
```json
{
  "error": "user_id is required"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to get stats"
}
```

---

## 6. Get Portfolio (Bonus)

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/portfolio/:userId`  
**Headers:** None  
**Body:** None

**Path Parameters:**
- `userId` (string, required) - Replace `:userId` in URL

**Example URL:**
```
{{base_url}}/api/v1/portfolio/user123
```

**Expected Response (Success):**
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

**Empty Response (No holdings):**
```json
{
  "holdings": [],
  "total_value": "0.0000"
}
```

**Status Code:** `200 OK`

**Error Responses:**

**400 Bad Request:**
```json
{
  "error": "user_id is required"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to get portfolio"
}
```

---

## Testing Workflow

### Step 1: Health Check
1. Test the health endpoint to verify the server is running
2. Should return `{"status": "ok"}`

### Step 2: Create Rewards
1. Create multiple rewards for the same user with different stocks
2. Create rewards for different users
3. Test with and without optional fields
4. Test duplicate event_id (should return 409)

**Example Test Cases:**

**Test Case 1: Create Reward with All Fields**
```json
POST {{base_url}}/api/v1/reward
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "reward_timestamp": "2024-01-15T10:30:00Z",
  "event_id": "event-001"
}
```

**Test Case 2: Create Reward with Minimal Fields**
```json
POST {{base_url}}/api/v1/reward
{
  "user_id": "user123",
  "stock_symbol": "TCS",
  "quantity": "5.25"
}
```

**Test Case 3: Create Reward for Different User**
```json
POST {{base_url}}/api/v1/reward
{
  "user_id": "user456",
  "stock_symbol": "INFOSYS",
  "quantity": "20.0",
  "event_id": "event-002"
}
```

**Test Case 4: Test Duplicate Event ID**
```json
POST {{base_url}}/api/v1/reward
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": "10.5",
  "event_id": "event-001"
}
```
Should return `409 Conflict` with error message.

### Step 3: Get Today's Stocks
1. Test with user who has rewards today
2. Test with user who has no rewards today (should return empty array)
3. Test with invalid user_id format

**Example:**
```
GET {{base_url}}/api/v1/today-stocks/user123
```

### Step 4: Get Stats
1. Test with user who has rewards
2. Test with user who has no rewards
3. Verify current_portfolio_value is calculated correctly

**Example:**
```
GET {{base_url}}/api/v1/stats/user123
```

### Step 5: Get Portfolio
1. Test with user who has holdings
2. Test with user who has no holdings
3. Verify current prices and values are calculated correctly

**Example:**
```
GET {{base_url}}/api/v1/portfolio/user123
```

### Step 6: Get Historical INR
1. Test with user who has historical data
2. Test with user who has no historical data
3. Note: Historical data is created by the hourly job for past dates

**Example:**
```
GET {{base_url}}/api/v1/historical-inr/user123
```

---

## Common Test Scenarios

### Scenario 1: Complete User Journey
1. Create reward for user123 with RELIANCE stock
2. Create another reward for user123 with TCS stock
3. Get today's stocks for user123 (should show both)
4. Get stats for user123 (should show both stocks)
5. Get portfolio for user123 (should show both holdings)

### Scenario 2: Multiple Users
1. Create rewards for user123
2. Create rewards for user456
3. Verify each user's data is isolated
4. Test stats and portfolio for each user separately

### Scenario 3: Edge Cases
1. Test duplicate event_id (should fail)
2. Test invalid stock symbol format
3. Test negative quantity (if allowed)
4. Test very large quantities
5. Test fractional shares (e.g., 0.000001)

### Scenario 4: Error Handling
1. Test missing required fields
2. Test invalid data types
3. Test invalid date format
4. Test non-existent user_id

---

## Postman Collection Setup

### Import Collection
1. Import `Stocky.postman_collection.json` into Postman
2. Set environment variable `base_url` to `http://localhost:8080`

### Create Environment
1. Create a new environment in Postman
2. Add variable:
   - Variable: `base_url`
   - Initial Value: `http://localhost:8080`
   - Current Value: `http://localhost:8080`

### Pre-request Scripts (Optional)
You can add pre-request scripts to generate dynamic data:

```javascript
// Generate random user_id
pm.environment.set("user_id", "user" + Math.floor(Math.random() * 1000));

// Generate random event_id
pm.environment.set("event_id", "event-" + Date.now());
```

### Tests (Optional)
Add tests to verify responses:

```javascript
// Test for successful response
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

// Test for response structure
pm.test("Response has correct structure", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('user_id');
    pm.expect(jsonData).to.have.property('stock_symbol');
});
```

---

## Quick Reference

| Endpoint | Method | Headers | Body Required |
|----------|--------|---------|---------------|
| `/health` | GET | None | No |
| `/api/v1/reward` | POST | `Content-Type: application/json` | Yes |
| `/api/v1/today-stocks/:userId` | GET | None | No |
| `/api/v1/historical-inr/:userId` | GET | None | No |
| `/api/v1/stats/:userId` | GET | None | No |
| `/api/v1/portfolio/:userId` | GET | None | No |

---

## Notes

1. **Server must be running** before testing
2. **Database must be set up** and migrations run
3. **Port 8080** is default (change if needed)
4. **All timestamps** should be in RFC3339 format (e.g., `2024-01-15T10:30:00Z`)
5. **Quantities** are strings but must be numeric values
6. **Stock symbols** are typically uppercase (e.g., `RELIANCE`, `TCS`, `INFOSYS`)

---

## Troubleshooting

### Connection Refused
- Check if server is running
- Verify port number (default: 8080)
- Check firewall settings

### 500 Internal Server Error
- Check database connection
- Verify database migrations ran successfully
- Check server logs for details

### 400 Bad Request
- Verify request body format (must be valid JSON)
- Check required fields are present
- Verify data types match expected format

### Empty Responses
- Verify data exists in database
- Check date filters (today's stocks only shows today's data)
- Historical data requires hourly job to run

---

## Sample Test Data

### User IDs
- `user123`
- `user456`
- `user789`

### Stock Symbols
- `RELIANCE`
- `TCS`
- `INFOSYS`
- `HDFCBANK`
- `ICICIBANK`

### Quantities
- `10.5` (fractional share)
- `5.25`
- `100.0`
- `0.5`

### Event IDs
- `event-001`
- `event-002`
- `event-003`
- Or use UUIDs: `550e8400-e29b-41d4-a716-446655440000`

