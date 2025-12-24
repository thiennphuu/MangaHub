# MangaHub UDP Service - System Architecture Diagram

## Overview

UDP Service provides lightweight chapter release notifications using UDP protocol. It broadcasts notifications to all registered clients in real-time without requiring persistent connections.

## UDP Service Architecture (Mermaid Diagram)

```mermaid
graph LR
    subgraph Clients["CLI Clients"]
        CLI1[CLI Client 1<br/>notify subscribe]
        CLI2[CLI Client 2<br/>notify subscribe]
        CLI3[CLI Client N<br/>notify unsubscribe]
    end

    subgraph UDPClient["UDP Client Layer"]
        Connect[Connect<br/>net.DialUDP]
        Register[Register<br/>send 'register']
        Unregister[Unregister<br/>send 'unregister']
        SendNotification[SendNotification<br/>JSON marshal]
        Listen[Listen<br/>ReadFromUDP<br/>1s timeout]
    end

    subgraph UDPServer["UDP Server<br/>Port 9091"]
        UDPListener[UDP Listener<br/>net.ListenUDP<br/>ReadFromUDP loop]
        ClientsMap[Clients Map<br/>map[string]*net.UDPAddr<br/>mutex sync.RWMutex]
        NotificationQueue[Notification Queue<br/>chan NotificationPayload<br/>buffer: 100]
        MessageHandler[Message Handler<br/>ReadFromUDP loop<br/>switch message type]
        BroadcastHandler[Broadcast Handler<br/>single goroutine<br/>WriteToUDP to all]
    end

    subgraph Data["Data Model"]
        NotificationPayload[NotificationPayload<br/>Type, MangaID<br/>Message, Timestamp<br/>JSON]
    end

    CLI1 -->|uses| Register
    CLI2 -->|uses| Listen
    CLI3 -->|uses| Unregister

    Connect -->|UDP| UDPListener
    Register -->|UDP 'register'| UDPListener
    Unregister -->|UDP 'unregister'| UDPListener
    SendNotification -->|UDP + JSON| UDPListener

    UDPListener -->|read| MessageHandler
    MessageHandler -->|'register'| ClientsMap
    MessageHandler -->|'unregister'| ClientsMap
    MessageHandler -->|JSON payload| NotificationQueue
    MessageHandler -->|auto-register| ClientsMap

    NotificationQueue -->|receive| BroadcastHandler
    BroadcastHandler -->|marshal JSON| NotificationPayload
    BroadcastHandler -->|WriteToUDP| ClientsMap
    ClientsMap -->|send to all| Listen

    style Clients fill:#e1f5ff
    style UDPClient fill:#fff4e1
    style UDPServer fill:#e8f5e9
    style Data fill:#f3e5f5
```

## UDP Service Architecture (Text Diagram)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                         │
│  │ CLI Client 1 │  │ CLI Client 2 │  │ CLI Client N │                         │
│  │ (Device A)   │  │ (Device B)   │  │ (Device C)   │                         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                         │
└─────────┼──────────────────┼──────────────────┼───────────────────────────────┘
          │                  │                  │
          │              UDP Connection (Port 9091)                                │
          │                  │                  │
┌─────────▼──────────────────▼──────────────────▼─────────────────────────────────┐
│                         UDP SERVER LAYER                                        │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  UDP Listener (net.ListenUDP)                                            │  │
│  │  ├─ Resolves UDP address (0.0.0.0:9091)                                  │  │
│  │  ├─ Listens for incoming UDP packets                                     │  │
│  │  └─ ReadFromUDP loop (1024 byte buffer)                                  │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Client Manager                                                           │  │
│  │  ├─ Clients: map[string]*net.UDPAddr                                     │  │
│  │  ├─ Thread-safe: sync.RWMutex                                            │  │
│  │  ├─ Register: Add client address (clientID = remoteAddr.String())        │  │
│  │  └─ Unregister: Remove client from map                                   │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Notification Queue                                                      │  │
│  │  ├─ Type: chan models.NotificationPayload                                │  │
│  │  ├─ Buffer: 100                                                          │  │
│  │  └─ Receives notifications for broadcasting                              │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
          │                  │                  │
          │                  │                  │
┌─────────▼──────────────────▼──────────────────▼─────────────────────────────────┐
│                         HANDLER LAYER                                          │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Message Handler (ReadFromUDP loop)                                       │  │
│  │  ├─ Reads: ReadFromUDP(buffer) - 1024 bytes                              │  │
│  │  ├─ Message Types:                                                        │  │
│  │  │  ├─ "register" → Add to Clients map                                   │  │
│  │  │  ├─ "unregister" → Remove from Clients map                            │  │
│  │  │  └─ JSON payload → Parse NotificationPayload, auto-register, broadcast│  │
│  │  └─ Auto-registration: Register client if sending JSON payload            │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  Broadcast Handler (handleBroadcast)                                      │  │
│  │  ├─ Single goroutine (runs continuously)                                 │  │
│  │  ├─ Receives: NotificationPayload from Queue channel                     │  │
│  │  ├─ Marshals: json.Marshal(notification)                                 │  │
│  │  ├─ Iterates: All clients in Clients map                                 │  │
│  │  └─ Sends: WriteToUDP(data, addr) to each registered client              │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
          │
          │ NotificationPayload Data
          │
┌─────────▼─────────────────────────────────────────────────────────────────────┐
│                         DATA MODEL                                             │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  NotificationPayload (models.NotificationPayload)                       │  │
│  │  ├─ Type: string (e.g., "chapter_release")                               │  │
│  │  ├─ MangaID: string                                                       │  │
│  │  ├─ Message: string                                                       │  │
│  │  └─ Timestamp: int64 (Unix timestamp)                                     │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. UDP Server (`internal/udp/server.go`)

**Server Struct:**

- `Port`: Server port (default: 9091)
- `Clients`: Map of client ID (string) to `*net.UDPAddr` (in-memory, ephemeral)
- `Queue`: Channel for notification payloads (buffer: 100)
- `mutex`: `sync.RWMutex` for thread-safe access to Clients map
- `done`: Channel for shutdown signal
- `logger`: Logger instance

**Key Methods:**

- `NewServer(port, logger)`: Creates new UDP server
- `Start()`: Starts listening and handling UDP packets
- `handleBroadcast(conn)`: Broadcasts notifications to all registered clients
- `SendNotification(notification)`: Adds notification to broadcast queue
- `RegisterClient(clientID, addr)`: Manually register a client
- `UnregisterClient(clientID)`: Manually unregister a client
- `GetClientCount()`: Returns number of registered clients
- `Stop()`: Gracefully shuts down server

### 2. Message Handling

**ReadFromUDP Loop:**

- Continuously reads UDP packets (1024 byte buffer)
- Parses message type:
  - `"register"`: Add client to Clients map, send confirmation
  - `"unregister"`: Remove client from Clients map, send confirmation
  - JSON payload: Parse as NotificationPayload, auto-register if needed, broadcast

**Auto-Registration:**

- If client sends JSON payload without prior registration, automatically registers
- Client ID is derived from `remoteAddr.String()`

**Broadcast Handler:**

- Single goroutine running continuously
- Receives notifications from Queue channel
- Marshals notification to JSON
- Iterates all registered clients and sends via `WriteToUDP`
- Thread-safe read using `mutex.RLock()`

### 3. UDP Client (`pkg/client/udp_client.go`)

**UDPClient Struct:**

- `ServerAddr`: Server address (host:port)
- `conn`: UDP connection (`*net.UDPConn`)
- `localAddr`: Local UDP address
- `Done`: Channel for stopping listener

**Key Methods:**

- `NewUDPClient(serverAddr)`: Creates new UDP client
- `Connect()`: Establishes UDP connection via `net.DialUDP`
- `Register()`: Sends "register" message to server
- `Unregister()`: Sends "unregister" message to server
- `SendNotification(payload)`: Sends JSON notification payload to server
- `Listen(callback)`: Continuously listens for notifications with timeout
- `Close()`: Closes UDP connection
- `GetLocalAddr()`: Returns local address string

### 4. CLI Commands (`internal/cli/notify/`)

**Commands:**

- `mangahub notify subscribe`: Subscribe to notifications (register with UDP server)
- `mangahub notify unsubscribe`: Unsubscribe from notifications (unregister from server)
- `mangahub notify preferences`: View and update notification preferences
- `mangahub notify test`: Test notification system by sending test notification

## Communication Flow

### Client Registration Flow:

```
CLI Client → UDP Connect → Send "register" → Server ReadFromUDP
  → Add to Clients map → Send confirmation JSON → Client receives confirmation
```

### Notification Broadcast Flow:

```
Client/Server → Send JSON NotificationPayload → Server ReadFromUDP
  → Parse JSON → Auto-register if needed → Send to Queue channel
  → Broadcast Handler → Marshal JSON → WriteToUDP to all registered clients
  → Clients receive notification via Listen()
```

### Client Unregistration Flow:

```
CLI Client → Send "unregister" → Server ReadFromUDP
  → Remove from Clients map → Send confirmation JSON → Client receives confirmation
```

## Protocol Details

### Message Format

- **Protocol**: JSON over UDP (stateless)
- **Encoding**: UTF-8 JSON
- **Structure**: `models.NotificationPayload` or plain text commands
- **Buffer Size**: 1024 bytes (server), 4096 bytes (client)

### Message Types

**Text Commands:**

- `"register"`: Register client for notifications
- `"unregister"`: Unregister client from notifications

**JSON Payload:**

```json
{
  "type": "chapter_release",
  "manga_id": "one-piece",
  "message": "Chapter 1101 released!",
  "timestamp": 1705312200
}
```

### Example Notification:

```json
{
  "type": "chapter_release",
  "manga_id": "one-piece",
  "message": "Chapter 1101 released!",
  "timestamp": 1705312200
}
```

### Connection Lifecycle

1. **Connect**: Client establishes UDP connection via `net.DialUDP`
2. **Register**: Client sends "register" message, server adds to Clients map
3. **Listen**: Client calls `Listen()` with callback to receive notifications
4. **Receive Notifications**: Server broadcasts notifications via `WriteToUDP`
5. **Send Notification**: Client or server sends JSON payload, server broadcasts to all
6. **Unregister**: Client sends "unregister" message, server removes from map
7. **Disconnect**: Client closes connection via `Close()`

## Concurrency Model

- **Main Loop**: Single goroutine reading from UDP socket
- **Broadcast Handler**: Single goroutine for broadcasting notifications
- **Thread Safety**: `sync.RWMutex` protects `Clients` map
- **Channel Buffer**: Notification queue buffer (100) prevents blocking
- **Stateless**: No persistent connections, clients identified by UDP address

## Error Handling

- **UDP Read Errors**: Logged, continue reading
- **JSON Parse Errors**: Logged, message ignored
- **Broadcast Errors**: Logged per client, continue to other clients
- **Connection Timeouts**: Client `Listen()` uses 1-second timeout for non-blocking reads
- **Graceful Shutdown**: Close `done` channel on `Stop()`

## Configuration

**From `config.yaml`:**

```yaml
udp:
  host: 0.0.0.0
  port: 9091
```

## Use Cases

1. **Chapter Release Notifications**: Broadcast new chapter releases to subscribed users
2. **Real-time Alerts**: Lightweight notifications without persistent connections
3. **Multi-Client Broadcasting**: Send notifications to multiple clients simultaneously
4. **Stateless Communication**: No connection overhead, suitable for ephemeral notifications

## Performance Considerations

- **Stateless Protocol**: No connection overhead, efficient for one-way notifications
- **In-Memory Client List**: Fast lookups, but lost on server restart
- **Channel Buffering**: Prevents blocking on notification queue
- **UDP Characteristics**: Fast but unreliable (no delivery guarantee)
- **No Database Integration**: Clients managed in memory only (as per specification)

## Differences from TCP Service

- **Stateless**: No persistent connections, clients identified by address
- **No Database**: Client list and notifications not persisted
- **Fire-and-Forget**: No delivery confirmation
- **Lighter Weight**: Lower overhead, suitable for notifications
- **Auto-Registration**: Clients can auto-register by sending JSON payload
