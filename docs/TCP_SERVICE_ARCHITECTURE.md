# MangaHub TCP Service - System Architecture Diagram

## Overview

TCP Service provides real-time progress synchronization across multiple devices using TCP protocol. It broadcasts reading progress updates to all connected clients.

## TCP Service Architecture (Mermaid Diagram)

```mermaid
graph LR
    subgraph Clients["CLI Clients"]
        CLI1[CLI Client 1<br/>sync connect]
        CLI2[CLI Client 2<br/>sync monitor]
        CLI3[CLI Client N<br/>sync status]
    end

    subgraph TCPClient["TCP Client Layer"]
        Connect[Connect<br/>net.Dial<br/>5s timeout]
        SendUpdate[SendUpdate<br/>JSON marshal<br/>fmt.Fprintf]
        MonitorUpdates[MonitorUpdates<br/>bufio.Scanner<br/>long-lived]
        CheckStatus[CheckStatus<br/>test connection]
    end

    subgraph TCPServer["TCP Server<br/>Port 9090"]
        Listener[TCP Listener<br/>net.Listen<br/>Accept loop]
        ConnMap[Connections Map<br/>map[string]net.Conn<br/>mutex sync.RWMutex]
        BroadcastChan[Broadcast Channel<br/>chan ProgressUpdate<br/>buffer: 100]
        ConnHandler[handleConnection<br/>goroutine per client<br/>bufio.Reader<br/>ReadString '\n']
        BroadcastHandler[handleBroadcast<br/>single goroutine<br/>range channel]
    end

    subgraph Data["Data & Persistence"]
        ProgressUpdate[ProgressUpdate<br/>UserID, MangaID<br/>Chapter, Timestamp<br/>DeviceID]
        Database[(SQLite Database<br/>user_progress table)]
    end

    CLI1 -->|uses| Connect
    CLI2 -->|uses| MonitorUpdates
    CLI3 -->|uses| CheckStatus

    Connect -->|TCP| Listener
    SendUpdate -->|TCP + JSON| Listener
    MonitorUpdates -->|TCP| Listener

    Listener -->|accept| ConnMap
    Listener -->|spawn| ConnHandler

    ConnHandler -->|read JSON| ProgressUpdate
    ProgressUpdate -->|send| BroadcastChan
    BroadcastChan -->|receive| BroadcastHandler
    BroadcastHandler -->|marshal JSON| ProgressUpdate
    BroadcastHandler -->|fmt.Fprintf| ConnMap
    ConnMap -->|send to all| MonitorUpdates

    ConnHandler -->|save| Database

    style Clients fill:#e1f5ff
    style TCPClient fill:#fff4e1
    style TCPServer fill:#e8f5e9
    style Data fill:#f3e5f5
```

## TCP Service Architecture (Text Diagram)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                         │
│  │ CLI Client 1 │  │ CLI Client 2 │  │ CLI Client N │                         │
│  │ (Device A)   │  │ (Device B)   │  │ (Device C)   │                         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                         │
└─────────┼──────────────────┼──────────────────┼───────────────────────────────┘
          │                  │                  │
          │              TCP Connection (Port 9090)                                │
          │                  │                  │
┌─────────▼──────────────────▼──────────────────▼─────────────────────────────────┐
│                         TCP SERVER LAYER                                        │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  TCP Listener (net.Listen)                                                │  │
│  │  ├─ Accepts incoming connections                                          │  │
│  │  ├─ Creates connection ID (conn_0, conn_1, ...)                          │  │
│  │  └─ Spawns goroutine per connection                                       │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Connection Manager                                                       │  │
│  │  ├─ Connections: map[string]net.Conn                                     │  │
│  │  ├─ Thread-safe: sync.RWMutex                                            │  │
│  │  ├─ Add: Directly in Accept loop (conn_0, conn_1, ...)                  │  │
│  │  └─ Remove: In handleConnection defer (on disconnect)                   │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Broadcast Channel                                                       │  │
│  │  ├─ Type: chan models.ProgressUpdate                                     │  │
│  │  ├─ Buffer: 100                                                          │  │
│  │  └─ Receives updates from all connections                                │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
          │                  │                  │
          │                  │                  │
┌─────────▼──────────────────▼──────────────────▼─────────────────────────────────┐
│                         HANDLER LAYER                                          │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Connection Handler (handleConnection)                                    │  │
│  │  ├─ Per-connection goroutine                                             │  │
│  │  ├─ Reads: bufio.NewReader(conn)                                          │  │
│  │  ├─ Protocol: JSON over TCP (newline-delimited)                         │  │
│  │  ├─ Parses: json.Unmarshal → ProgressUpdate                              │  │
│  │  ├─ Saves: saveProgressUpdate() → Database (user_progress table)          │  │
│  │  └─ Sends: update → Broadcast channel                                     │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Broadcast Handler (handleBroadcast)                                      │  │
│  │  ├─ Single goroutine (runs continuously)                                 │  │
│  │  ├─ Receives: ProgressUpdate from Broadcast channel                      │  │
│  │  ├─ Marshals: json.Marshal(update)                                       │  │
│  │  ├─ Iterates: All connections in map                                     │  │
│  │  └─ Sends: fmt.Fprintf(conn, "%s\n", data) to each client                 │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
          │
          │ ProgressUpdate Data
          │
┌─────────▼─────────────────────────────────────────────────────────────────────┐
│                         DATA MODEL                                             │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  ProgressUpdate (models.ProgressUpdate)                                   │  │
│  │  ├─ UserID: string                                                        │  │
│  │  ├─ MangaID: string                                                       │  │
│  │  ├─ Chapter: int                                                          │  │
│  │  ├─ Timestamp: int64 (Unix timestamp)                                     │  │
│  │  └─ DeviceID: string                                                      │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. TCP Server (`internal/tcp/server.go`)

**Server Struct:**

- `Port`: Server port (default: 9090)
- `Connections`: Map of connection ID to `net.Conn` (conn_0, conn_1, ...)
- `Broadcast`: Channel for progress updates (buffer: 100)
- `Register`: Channel for new connections (defined but not used in current implementation)
- `Unregister`: Channel for disconnected connections (defined but not used in current implementation)
- `mutex`: `sync.RWMutex` for thread-safe access to Connections map
- `done`: Channel for shutdown signal
- `logger`: Logger instance
- `db`: Database connection for persisting progress updates

**Key Methods:**

- `NewServer(port, logger, db)`: Creates new TCP server with database connection
- `Start()`: Starts listening and accepting connections
- `handleConnection(connID, conn)`: Handles individual client connection
- `handleBroadcast()`: Broadcasts updates to all clients
- `Stop()`: Gracefully shuts down server
- `GetConnectionCount()`: Returns active connection count
- `saveProgressUpdate(update)`: Saves progress update to database (INSERT or UPDATE)

### 2. Connection Handling

**Accept Loop:**

```
Listener.Accept() → Create connID → Add to Connections map → Spawn handleConnection goroutine
```

**Connection Handler:**

- Reads JSON messages line-by-line using `bufio.NewReader`
- Parses `ProgressUpdate` from JSON
- Saves progress update to database (via `saveProgressUpdate()`)
- Sends update to `Broadcast` channel
- Handles connection cleanup on disconnect

**Broadcast Handler:**

- Single goroutine running continuously
- Receives updates from `Broadcast` channel
- Marshals update to JSON
- Iterates all connections and sends update
- Thread-safe read using `mutex.RLock()`

### 3. TCP Client (`pkg/client/tcp_client.go`)

**TCPClient Struct:**

- `Addr`: Server address (host:port)

**Key Methods:**

- `NewTCPClient(host, port)`: Creates new client
- `Connect()`: Establishes TCP connection
- `CheckStatus()`: Verifies server reachability
- `SendUpdate(update)`: Sends progress update to server
- `MonitorUpdates(stop, handler)`: Continuously receives broadcast updates

### 4. CLI Commands (`internal/cli/sync/`)

**Commands:**

- `mangahub sync connect`: Connect to TCP server
- `mangahub sync disconnect`: Disconnect from server
- `mangahub sync status`: Check server status
- `mangahub sync monitor`: Monitor real-time updates

## Communication Flow

### Client Sends Update Flow:

```
CLI Client → TCP Connect → Server Accept → Connection Handler
  → Read JSON → Parse ProgressUpdate → Save to Database
  → Send to Broadcast Channel → Broadcast Handler → Marshal JSON → Send to All Clients
```

### Client Receives Update Flow:

```
Server Broadcast Handler → Marshal Update → Send to All Connections
  → CLI Client → Read JSON → Parse ProgressUpdate → Display/Process
```

## Protocol Details

### Message Format

- **Protocol**: JSON over TCP
- **Delimiter**: Newline (`\n`)
- **Encoding**: UTF-8 JSON
- **Structure**: `models.ProgressUpdate`

### Example Message:

```json
{
  "user_id": "user123",
  "manga_id": "manga456",
  "chapter": 42,
  "timestamp": 1705312200,
  "device_id": "device-abc-123"
}
```

### Connection Lifecycle

1. **Connect**: Client establishes TCP connection via `net.Dial`
2. **Accept**: Server accepts connection, creates `connID` (conn_0, conn_1, ...), adds to Connections map
3. **Handle**: Server spawns `handleConnection` goroutine for each client
4. **Send Update**: Client sends JSON ProgressUpdate (newline-delimited) via `fmt.Fprintf`
5. **Receive & Save**: Server reads JSON, unmarshals, saves to database (INSERT or UPDATE `user_progress` table)
6. **Broadcast**: Server sends update to Broadcast channel, then broadcasts to all connections
7. **Monitor**: Clients using `MonitorUpdates` receive broadcast updates via `bufio.Scanner`
8. **Disconnect**: Connection closed (client or server), removed from map in `handleConnection` defer

## Concurrency Model

- **Main Loop**: Single goroutine accepting connections
- **Connection Handlers**: One goroutine per client connection
- **Broadcast Handler**: Single goroutine for broadcasting
- **Thread Safety**: `sync.RWMutex` protects `Connections` map
- **Channel Buffer**: Broadcast channel buffer (100) prevents blocking

## Error Handling

- **Connection Errors**: Logged, connection removed from map
- **JSON Parse Errors**: Logged, message skipped
- **Database Errors**: Logged, update still broadcasted to clients (non-blocking)
- **Broadcast Errors**: Logged, update skipped for that client
- **Graceful Shutdown**: All connections closed on `Stop()`

## Configuration

**From `config.yaml`:**

```yaml
tcp:
  host: 0.0.0.0
  port: 9090
  max_connections: 50
  read_buffer_size: 1024
  write_buffer_size: 1024
  keep_alive: true
  keep_alive_interval: 30
```

## Use Cases

1. **Multi-Device Sync**: User updates progress on Device A, all other devices receive update
2. **Real-time Monitoring**: Monitor all progress updates across all users
3. **Progress Broadcasting**: Share reading progress with connected clients
4. **Status Checking**: Verify TCP server availability

## Performance Considerations

- **Concurrent Connections**: Handles multiple clients simultaneously
- **Channel Buffering**: Prevents blocking on broadcast channel
- **Goroutine Per Connection**: Efficient I/O handling
- **Thread-Safe Operations**: Mutex protects shared state
- **Connection Limits**: Configurable max connections (default: 50)
