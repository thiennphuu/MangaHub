# Database Check Command Update

## Summary
The `db check` CLI command has been updated to check the remote database via HTTP API instead of checking a local database file.

## Changes Made

### 1. API Handler Updated ✓
- **File:** `internal/api/handler.go`
- **Route Added:** `GET /server/database/check` (protected, requires authentication)
- **Handler:** `GetDatabaseCheck()`
- **Helper Functions:**
  - `checkDatabaseIntegrity()` - Runs SQLite PRAGMA integrity_check
  - `verifyDatabaseTables()` - Verifies core tables exist

### 2. HTTP Client Method Added ✓
- **File:** `pkg/client/http_client.go`
- **Method:** `GetDatabaseCheck() (*DatabaseCheckResponse, error)`
- **Response Structure:**
```go
type DatabaseCheckResponse struct {
    Status    string // "healthy" or "unhealthy"
    Integrity struct {
        OK     bool     // true if integrity check passed
        Issues []string // list of integrity issues (if any)
    }
    Tables struct {
        Verified []string // list of verified tables
        Missing  []string // list of missing tables
    }
}
```

### 3. CLI Command Updated ✓
- **File:** `internal/cli/db/check.go`
- **Changes:**
  - Removed local SQLite database access
  - Now calls HTTP API endpoint
  - Requires authentication (uses session token)
  - Displays formatted results from server

### 4. API Server Main Updated ✓
- **File:** `cmd/api-server/main.go`
- **Health Endpoint Enhanced:** Now returns full server configuration

## API Endpoint Details

### Request
```http
GET /server/database/check
Authorization: Bearer {token}
```

### Response (Success)
```json
{
  "status": "healthy",
  "integrity": {
    "ok": true,
    "issues": []
  },
  "tables": {
    "verified": ["users", "manga", "user_progress"],
    "missing": []
  }
}
```

### Response (Issues Found)
```json
{
  "status": "unhealthy",
  "integrity": {
    "ok": false,
    "issues": ["corruption detected in table X"]
  },
  "tables": {
    "verified": ["users", "manga"],
    "missing": ["user_progress"]
  }
}
```

## CLI Command Usage

```bash
# Check database (requires login)
go run ./cmd/cli db check

# Check with specific profile
go run ./cmd/cli --profile user2 db check
```

## Expected Output

### Healthy Database
```
Checking remote database via HTTP API...

✓ Database Status: healthy

Integrity Check:
✓ PRAGMA integrity_check: ok

Table Verification:
✓ Table users exists
✓ Table manga exists
✓ Table user_progress exists

✓ Database check completed successfully
```

### Unhealthy Database
```
Checking remote database via HTTP API...

✓ Database Status: unhealthy

Integrity Check:
✗ Integrity issues found:
  - corruption detected

Table Verification:
✓ Table users exists
✗ Missing table: user_progress

Error: database check found issues
```

## Protocol Used
**HTTP REST API** ✓

## Current Status
- ✓ Client code updated and built successfully
- ✓ Server code updated and ready for deployment
- ⏳ Waiting for server deployment
- Current error: `404 page not found` (expected - server needs the new endpoint)

## Deployment Steps

### For Your Friend's Server (10.238.53.72)
1. Pull latest code from repository
2. Rebuild the API server:
   ```bash
   go build -o api-server ./cmd/api-server
   ```
3. Restart the API server
4. Test the endpoint:
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
        http://10.238.53.72:8080/server/database/check
   ```

## Files Modified
- ✓ `internal/api/handler.go` - Added endpoint and handler
- ✓ `pkg/client/http_client.go` - Added client method
- ✓ `internal/cli/db/check.go` - Updated to use HTTP API
- ✓ `cmd/api-server/main.go` - Enhanced health endpoint

## Notes
- The endpoint is **protected** and requires authentication
- Checks are performed on the server's database
- Returns both integrity check and table verification results
- Non-zero exit code if database has issues

