# Setup Instructions

## Prerequisites

1. **Go 1.21 or higher**
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **PostgreSQL 12 or higher**
   - Download from: https://www.postgresql.org/download/
   - Verify installation: `psql --version`

3. **Make sure PostgreSQL is running and accessible**

## Installation Steps

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Assignment
```

### 2. Create PostgreSQL Database

Connect to PostgreSQL and create the database:

```bash
psql -U postgres
```

Then run:

```sql
CREATE DATABASE assignment;
\q
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/assignment?sslmode=disable
```

**Note:** Update the `DATABASE_URL` with your PostgreSQL credentials:
- Replace `postgres` (username) with your PostgreSQL username
- Replace `postgres` (password) with your PostgreSQL password
- Replace `localhost:5432` with your PostgreSQL host and port if different

### 4. Install Dependencies

```bash
go mod download
```

This will download all required dependencies:
- `github.com/gin-gonic/gin` - Web framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/sirupsen/logrus` - Logging
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading

### 5. Run the Application

```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your `.env` file).

You should see output like:
```
{"level":"info","msg":"Server started","port":"8080"}
```

### 6. Verify Installation

Open a new terminal and test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## Database Migrations

Database migrations run automatically on application startup. The following tables will be created:

1. `reward_events` - Stores reward events
2. `ledger_entries` - Double-entry ledger
3. `stock_prices` - Historical stock prices
4. `portfolio_snapshots` - Daily portfolio snapshots

## Testing the API

### Using Postman

1. Import the `Stocky.postman_collection.json` file into Postman
2. Set the `base_url` variable to `http://localhost:8080`
3. Test each endpoint

### Using cURL

**1. Create a reward:**
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": "10.5",
    "event_id": "event-001"
  }'
```

**2. Get today's stocks:**
```bash
curl http://localhost:8080/api/v1/today-stocks/user123
```

**3. Get stats:**
```bash
curl http://localhost:8080/api/v1/stats/user123
```

**4. Get portfolio:**
```bash
curl http://localhost:8080/api/v1/portfolio/user123
```

**5. Get historical INR:**
```bash
curl http://localhost:8080/api/v1/historical-inr/user123
```

## Troubleshooting

### Database Connection Error

**Error:** `Failed to connect to database`

**Solution:**
1. Verify PostgreSQL is running: `pg_isready`
2. Check database credentials in `.env` file
3. Verify database exists: `psql -U postgres -l`
4. Check PostgreSQL is listening on the correct port

### Port Already in Use

**Error:** `bind: address already in use`

**Solution:**
1. Change the `PORT` in `.env` file
2. Or stop the process using port 8080

### Migration Errors

**Error:** `Failed to run migrations`

**Solution:**
1. Check database connection
2. Verify user has CREATE TABLE permissions
3. Check for existing tables that might conflict

### Module Not Found

**Error:** `cannot find module`

**Solution:**
```bash
go mod tidy
go mod download
```

## Development

### Project Structure

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
├── README.md                        # Main documentation
├── API_SPECIFICATION.md             # API documentation
├── DATABASE_SCHEMA.md               # Database schema documentation
├── EDGE_CASES.md                    # Edge cases and scaling
├── SETUP.md                         # This file
├── Stocky.postman_collection.json   # Postman collection
└── .env                             # Environment variables (create this)
```

### Running in Development Mode

The application runs in development mode by default. For production:

1. Set `GIN_MODE=release` environment variable
2. Use proper logging configuration
3. Set up proper error handling
4. Configure database connection pooling

### Building for Production

```bash
go build -o stocky main.go
./stocky
```

## Background Jobs

The application runs a background job that:
- Updates stock prices every hour
- Creates daily portfolio snapshots

This job runs automatically when the application starts.

## Next Steps

1. Test all API endpoints using Postman
2. Create sample rewards for different users
3. Verify portfolio calculations
4. Check historical data after running for a few days

## Support

For issues or questions, refer to:
- `README.md` - Main documentation
- `API_SPECIFICATION.md` - API details
- `DATABASE_SCHEMA.md` - Database structure
- `EDGE_CASES.md` - Edge cases and scaling

