# Server Logs API Update

## Summary
The `server logs` CLI command has been updated to fetch logs from the server via HTTP API instead of reading from a local file.

## Client-Side Changes (COMPLETED ✓)

### 1. CLI Command Updated
- **File:** `internal/cli/server/logs.go`
- **Changes:**
  - Removed local file reading logic
  - Now uses HTTP API to fetch logs from remote server
  - Requires authentication (uses session token)
  - Supports filtering by log level and max lines

### 2. HTTP Client Method Added
- **File:** `pkg/client/http_client.go`
- **Method:** `GetServerLogs(maxLines int, level string) (*ServerLogsResponse, error)`
- **Endpoint:** `GET /server/logs?max_lines={n}&level={level}`

### 3. Utility Functions Added
- **File:** `pkg/utils/logger.go`
- **Functions:**
  - `GetLogFilePath()` - Returns standard log file path
  - `OpenLogFile(path)` - Opens log file for reading
  - `ReadLogLines(file, level, maxLines)` - Reads and filters log lines

## Server-Side Changes (REQUIRED - Your friend needs to deploy these)

### 1. API Handler Updated
- **File:** `internal/api/handler.go`
- **Changes:**
  - Added route: `GET /server/logs` (protected, requires authentication)
  - Added handler: `GetServerLogs(c *gin.Context)`
  - Reads server log file and returns filtered logs as JSON

### 2. Route Registration
```go
// Server management routes
server := protected.Group("/server")
{
    server.GET("/logs", h.GetServerLogs)
}
```

### 3. API Response Format
```json
{
  "logs": ["[INFO] 2024-01-20 ...", "[ERROR] 2024-01-20 ..."],
  "count": 42,
  "max_lines": 100,
  "level": "error"
}
```

### 4. Query Parameters
- `max_lines` (int, default: 100) - Maximum number of log lines to return
- `level` (string, optional) - Filter by log level (debug, info, warn, error)

## Testing

### CLI Command Usage
```bash
# View last 100 logs (default)
go run ./cmd/cli server logs

# View last 50 logs
go run ./cmd/cli server logs --max-lines 50

# Filter by error level
go run ./cmd/cli server logs --level error

# Combine filters
go run ./cmd/cli server logs --max-lines 200 --level info
```

### Expected Output
```
Fetching server logs via HTTP API...
✓ Retrieved 42 log entries from server

Recent server logs (max: 100, level: all):

[INFO] 2024-01-20 10:15:30 Server started on :8080
[INFO] 2024-01-20 10:15:31 TCP server listening on :9090
[INFO] 2024-01-20 10:15:32 UDP server listening on :9091
...
```

## Deployment Steps for Server

1. **Pull latest code** from your repository
2. **Rebuild the server:**
   ```bash
   go build -o server ./cmd/server
   ```
3. **Restart the Docker container** (if using Docker) or restart the server process
4. **Verify the endpoint:**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
        http://10.238.53.72:8080/server/logs?max_lines=10
   ```

## Files Modified

### Client Side (Already Updated)
- `internal/cli/server/logs.go` - CLI command
- `pkg/client/http_client.go` - HTTP client method
- `pkg/utils/logger.go` - Utility functions

### Server Side (Need Deployment)
- `internal/api/handler.go` - API endpoint and handler

## Notes
- The endpoint is **protected** and requires authentication
- Logs are read from the server's log file at `~/.mangahub/logs/server.log`
- If the log file doesn't exist, a proper error is returned
- The API returns the most recent N lines (based on `max_lines` parameter)
- Log level filtering is optional (if omitted, all logs are returned)

## Current Status
- ✓ Client code updated and tested
- ✓ Build successful
- ⏳ Waiting for server deployment
- Current error: `404 page not found` (expected - server doesn't have the endpoint yet)

## Next Steps
1. Your friend needs to pull the latest code
2. They need to rebuild and restart their server
3. Once deployed, test with: `.\mangahub.exe server logs`

