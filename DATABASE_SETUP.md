# Database Setup Instructions

## Error Message
```
dial tcp [::1]:5432: connectex: No connection could be made because the target machine actively refused it.
```

This error means PostgreSQL is not running or not accessible on port 5432.

## Solution Steps

### 1. Check if PostgreSQL is Installed

**Windows:**
```powershell
# Check if PostgreSQL service exists
Get-Service -Name "*postgres*"

# Or check if psql is available
psql --version
```

**If PostgreSQL is not installed:**
- Download from: https://www.postgresql.org/download/windows/
- Install PostgreSQL (default port is 5432)
- Remember the password you set for the `postgres` user

### 2. Start PostgreSQL Service

**Windows:**
```powershell
# Start PostgreSQL service
Start-Service -Name "postgresql-x64-*"  # Replace * with your version number

# Or use Services GUI:
# Win + R -> services.msc -> Find PostgreSQL -> Right-click -> Start
```

**Check if it's running:**
```powershell
Get-Service -Name "*postgres*" | Select-Object Name, Status
```

### 3. Create the Database

**Connect to PostgreSQL:**
```powershell
psql -U postgres
```

**Or use pgAdmin (GUI tool that comes with PostgreSQL)**

**Create the database:**
```sql
CREATE DATABASE assignment;
\q
```

### 4. Update .env File

The `.env` file has been created with default values. Update it with your PostgreSQL credentials:

```env
PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/assignment?sslmode=disable
```

**Replace:**
- `username` with your PostgreSQL username (default: `postgres`)
- `password` with your PostgreSQL password
- `localhost:5432` if your PostgreSQL is on a different host/port

### 5. Test Connection

**Test PostgreSQL connection:**
```powershell
psql -U postgres -d assignment -c "SELECT version();"
```

**If connection works, run the application:**
```powershell
go run main.go
```

## Alternative: Use Docker (If PostgreSQL is not installed)

If you don't want to install PostgreSQL, you can use Docker:

**1. Install Docker Desktop:**
- Download from: https://www.docker.com/products/docker-desktop

**2. Run PostgreSQL in Docker:**
```powershell
docker run --name postgres-stocky -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=assignment -p 5432:5432 -d postgres
```

**3. Update .env file:**
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/assignment?sslmode=disable
```

**4. Run the application:**
```powershell
go run main.go
```

## Troubleshooting

### Port 5432 Already in Use
If port 5432 is already in use:
1. Find what's using it:
   ```powershell
   netstat -ano | findstr :5432
   ```
2. Either stop that service or change PostgreSQL port
3. Update `.env` file with new port

### Connection Refused
- Check if PostgreSQL service is running
- Check firewall settings
- Verify port number in connection string

### Authentication Failed
- Verify username and password in `.env` file
- Check PostgreSQL authentication settings in `pg_hba.conf`

### Database Does Not Exist
- Create the database: `CREATE DATABASE assignment;`
- Verify database name in `.env` file

## Quick Setup Script (PowerShell)

```powershell
# Check if PostgreSQL is installed
if (Get-Command psql -ErrorAction SilentlyContinue) {
    Write-Host "PostgreSQL is installed"
    
    # Create database
    $env:PGPASSWORD = "postgres"  # Set your password here
    psql -U postgres -c "CREATE DATABASE assignment;" 2>$null
    
    Write-Host "Database 'assignment' created (or already exists)"
} else {
    Write-Host "PostgreSQL is not installed. Please install it first."
    Write-Host "Download from: https://www.postgresql.org/download/windows/"
}
```

## Verify Setup

After setting up, test the connection:

```powershell
# Test connection
psql -U postgres -d assignment -c "SELECT 1;"

# If successful, run the application
go run main.go
```

The application should start without database connection errors.

