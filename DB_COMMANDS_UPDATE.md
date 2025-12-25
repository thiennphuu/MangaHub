# Database Commands Update

## Summary
All database management CLI commands (`check`, `optimize`, `stats`, `repair`) have been updated to use HTTP REST API instead of accessing local database files.

## Commands Updated

### 1. `db check` ✓
- **Endpoint:** `GET /server/database/check`
- **Function:** Checks database integrity and verifies core tables
- **Response:** Integrity status, issues, table verification results

### 2. `db optimize` ✓
- **Endpoint:** `POST /server/database/optimize`
- **Function:** Runs ANALYZE, REINDEX, and VACUUM commands
- **Response:** List of completed steps and any errors

### 3. `db stats` ✓
- **Endpoint:** `GET /server/database/stats`
- **Function:** Returns database file size and row counts
- **Response:** File size in MB and table row counts

### 4. `db repair` ✓
- **Endpoint:** `POST /server/database/repair`
- **Function:** Runs quick_check, VACUUM, and schema re-initialization
- **Response:** List of repair steps and any issues

## API Endpoints Added

### GET /server/database/check
```json
Response:
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

### POST /server/database/optimize
```json
Response:
{
  "status": "success",
  "steps": [
    "ANALYZE completed",
    "REINDEX completed",
    "VACUUM completed"
  ],
  "errors": []
}
```

### GET /server/database/stats
```json
Response:
{
  "file_size_bytes": 102400,
  "file_size_mb": 0.1,
  "tables": {
    "users": 5,
    "manga": 100,
    "user_progress": 50,
    "chat_messages": 200,
    "notifications": 10
  }
}
```

### POST /server/database/repair
```json
Response:
{
  "status": "success",
  "steps": [
    "PRAGMA quick_check: ok",
    "VACUUM completed",
    "Schema re-initialized"
  ],
  "errors": []
}
```

## CLI Command Usage

All commands now require authentication:

```bash
# Check database
go run ./cmd/cli --profile user2 db check

# Optimize database
go run ./cmd/cli --profile user2 db optimize

# Show database statistics
go run ./cmd/cli --profile user2 db stats

# Repair database
go run ./cmd/cli --profile user2 db repair
```

## Expected Output Examples

### db check (healthy)
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

### db optimize
```
Optimizing remote database via HTTP API...

✓ ANALYZE completed
✓ REINDEX completed
✓ VACUUM completed

✓ Database optimization completed successfully
```

### db stats
```
Fetching database statistics via HTTP API...

File size: 0.10 MB

Row counts:
  users: 5 rows
  manga: 100 rows
  user_progress: 50 rows
  chat_messages: 200 rows
  notifications: 10 rows
```

### db repair
```
Repairing remote database via HTTP API...

✓ PRAGMA quick_check: ok
✓ VACUUM completed
✓ Schema re-initialized

✓ Database repair/maintenance completed
```

## Protocol Used
**HTTP REST API** ✓

## Files Modified

### API Handler
- **File:** `internal/api/handler.go`
- Added 4 new endpoints and handler functions
- Added helper methods for database operations

### HTTP Client
- **File:** `pkg/client/http_client.go`
- Added 4 new methods and response types:
  - `GetDatabaseCheck()`
  - `OptimizeDatabase()`
  - `GetDatabaseStats()`
  - `RepairDatabase()`

### CLI Commands
- **File:** `internal/cli/db/check.go` - Updated to use HTTP API
- **File:** `internal/cli/db/optimize.go` - Updated to use HTTP API
- **File:** `internal/cli/db/stats.go` - Updated to use HTTP API
- **File:** `internal/cli/db/repair.go` - Updated to use HTTP API

## Current Status
- ✅ Client code updated and built successfully
- ✅ Server code updated and ready for deployment
- ⏳ Waiting for server deployment on `10.238.53.72`
- Current error: `404 page not found` (expected - endpoints need deployment)

## Deployment Steps for Server

1. **Pull latest code** from repository
2. **Rebuild the API server:**
   ```bash
   cd /path/to/MangaHub
   go build -o api-server ./cmd/api-server
   ```
3. **Restart the server** (or restart Docker container)
4. **Test endpoints:**
   ```bash
   # Check
   curl -H "Authorization: Bearer TOKEN" http://10.238.53.72:8080/server/database/check
   
   # Stats
   curl -H "Authorization: Bearer TOKEN" http://10.238.53.72:8080/server/database/stats
   
   # Optimize
   curl -X POST -H "Authorization: Bearer TOKEN" http://10.238.53.72:8080/server/database/optimize
   
   # Repair
   curl -X POST -H "Authorization: Bearer TOKEN" http://10.238.53.72:8080/server/database/repair
   ```

## Security Notes
- All endpoints are **protected** and require authentication
- Operations are performed on the server's database
- Users can only manage the remote database (no local database access)
- Commands require valid JWT token from login

## Benefits
1. ✅ Centralized database management
2. ✅ No need for local database files on client
3. ✅ Consistent database operations across all clients
4. ✅ Proper authentication and authorization
5. ✅ Server-side validation and error handling

