# MangaHub gRPC Service - System Architecture Diagram

## Overview

gRPC Service provides high-performance, type-safe service-to-service communication for manga operations. It uses Protocol Buffers for efficient serialization and supports JSON encoding for easier debugging and interoperability.

## gRPC Service Architecture (Mermaid Diagram)

```mermaid
graph LR
    subgraph Clients["CLI Clients"]
        CLI[grpc commands<br/>manga get/search/progress]
    end

    subgraph GRPCClient["gRPC Client<br/>pkg/client/grpc_client.go"]
        GRPCClientNode[GRPCClient<br/>Connect + RPC calls]
    end

    subgraph GRPCServer["gRPC Server<br/>cmd/grpc-server/main.go"]
        GRPCServerNode[gRPC Server<br/>tcp :9092]
    end

    subgraph GRPCService["gRPC Service<br/>internal/grpc/service/manga.go"]
        GRPCMangaService[MangaService<br/>GetManga / Search / UpdateProgress / GetTop10]
    end

    subgraph MangaServiceLayer["Manga Service<br/>internal/manga/service.go"]
        MangaSvcNode[Manga Service<br/>DB queries]
    end

    subgraph Database["SQLite Database"]
        MangaTable[manga table]
    end

    subgraph Proto["Proto Definitions<br/>proto/manga.proto"]
        ProtoNode[MangaService + messages]
    end

    CLI -->|commands| GRPCClientNode
    GRPCClientNode -->|gRPC (JSON)| GRPCServerNode
    GRPCServerNode -->|RPC dispatch| GRPCMangaService
    GRPCMangaService -->|business calls| MangaSvcNode
    MangaSvcNode -->|SQL| MangaTable

    ProtoNode -->|proto → Go code| GRPCServerNode
    ProtoNode -->|proto → Go code| GRPCClientNode

    %% Styling
    classDef client fill:#E3F2FD,stroke:#1E88E5,stroke-width:1px,color:#0D47A1;
    classDef grpcClient fill:#E8F5E9,stroke:#43A047,stroke-width:1px,color:#1B5E20;
    classDef grpcServer fill:#FFF3E0,stroke:#FB8C00,stroke-width:1px,color:#E65100;
    classDef service fill:#F3E5F5,stroke:#8E24AA,stroke-width:1px,color:#4A148C;
    classDef db fill:#FBE9E7,stroke:#D84315,stroke-width:1px,color:#BF360C;
    classDef proto fill:#ECEFF1,stroke:#546E7A,stroke-width:1px,color:#263238;

    class CLI client;
    class GRPCClientNode grpcClient;
    class GRPCServerNode grpcServer;
    class GRPCMangaService service;
    class MangaSvcNode service;
    class MangaTable db;
    class ProtoNode proto;
```

## Text Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI Clients                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │ grpc manga   │  │ grpc manga   │  │ grpc        │           │
│  │ get          │  │ search       │  │ progress    │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ uses
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    gRPC Client Layer                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ GRPCClient (pkg/client/grpc_client.go)                   │  │
│  │  - Connect() - grpc.DialContext with JSON codec          │  │
│  │  - GetManga(mangaID) - MangaRequest                      │  │
│  │  - SearchManga(query, limit) - SearchRequest             │  │
│  │  - UpdateProgress(userID, mangaID, chapter)             │  │
│  │  - GetTop10Manga() - Empty                               │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ gRPC + JSON over TCP
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      gRPC Server                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ TCP Listener (net.Listen "tcp :9092")                   │  │
│  │                                                           │  │
│  │ gRPC Server (grpc.NewServer)                            │  │
│  │  - JSONCodec (custom JSON encoding)                     │  │
│  │  - Service Registry (RegisterMangaServiceServer)         │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ route to service
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    gRPC Service Layer                            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ MangaService (internal/grpc/service/manga.go)           │  │
│  │  - mangaService *manga.Service                           │  │
│  │  - logger *utils.Logger                                  │  │
│  │                                                           │  │
│  │ Methods:                                                 │  │
│  │  1. GetManga(ctx, MangaRequest) → MangaResponse         │  │
│  │     → mangaService.GetByID()                             │  │
│  │                                                           │  │
│  │  2. SearchManga(ctx, SearchRequest) → SearchResponse    │  │
│  │     → mangaService.Search(filter)                       │  │
│  │                                                           │  │
│  │  3. UpdateProgress(ctx, UpdateProgressRequest)          │  │
│  │     → TODO: implement database save                      │  │
│  │                                                           │  │
│  │  4. GetTop10Manga(ctx, Empty) → Top10Response          │  │
│  │     → mangaService.List(10, 0)                          │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ delegate to business service
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Manga Service                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ Service (internal/manga/service.go)                      │  │
│  │  - GetByID(id) → SELECT manga WHERE id = ?               │  │
│  │  - Search(filter) → SELECT with WHERE filters           │  │
│  │  - List(limit, offset) → SELECT LIMIT/OFFSET            │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ SQL queries
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    SQLite Database                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ manga table                                               │  │
│  │  - id (TEXT PRIMARY KEY)                                 │  │
│  │  - title, author, genres, status                         │  │
│  │  - chapters, description, cover_url                      │  │
│  │  - created_at, updated_at                                │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│              Protocol Buffers Definition                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ proto/manga.proto                                        │  │
│  │  - Messages: MangaRequest, MangaResponse,                │  │
│  │             SearchRequest, SearchResponse,                │  │
│  │             UpdateProgressRequest, UpdateProgressResponse,│  │
│  │             Top10Response, Empty                         │  │
│  │  - Service: MangaService (4 RPC methods)                 │  │
│  │  - Code generation: Go code for server & client         │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. CLI Clients (`internal/cli/grpc/`)

**Commands:**

- `grpc manga get --id <manga-id>` - Get manga details by ID
- `grpc manga search --query <query> --limit <n>` - Search manga
- `grpc progress update --manga-id <id> --chapter <n> --user-id <id>` - Update progress
- `grpc top10` - Get top 10 manga (if implemented)

**Features:**

- Connects to gRPC server using `pkg/client.GRPCClient`
- Displays formatted results
- Error handling and connection management

### 2. gRPC Client (`pkg/client/grpc_client.go`)

**GRPCClient Struct:**

- `ServerAddr string` - Server address (default: localhost:9092)
- `conn *grpc.ClientConn` - gRPC connection
- `client pb.MangaServiceClient` - Generated client stub

**Key Methods:**

- `Connect()` - Establishes connection with JSON codec
  - Uses `grpc.DialContext` with `insecure.NewCredentials()`
  - Sets `CallContentSubtype("json")` for JSON encoding
  - 5-second connection timeout
- `GetManga(mangaID)` - Calls `GetManga` RPC
  - Creates `MangaRequest` with ID
  - Returns `GRPCMangaResponse` wrapper
  - 10-second timeout
- `SearchManga(query, limit)` - Calls `SearchManga` RPC
  - Creates `SearchRequest` with title and limit
  - Returns `SearchResponse` with results array
- `UpdateProgress(userID, mangaID, chapter)` - Calls `UpdateProgress` RPC
  - Creates `UpdateProgressRequest`
  - Returns `UpdateProgressResponse`
- `GetTop10Manga()` - Calls `GetTop10Manga` RPC
  - Uses `Empty` message
  - Returns `Top10Response` with rankings

**JSON Codec:**

- Custom `JSONCodec` implements `grpc.encoding.Codec`
- Uses standard `json.Marshal`/`Unmarshal`
- Registered in `init()` function

### 3. gRPC Server (`cmd/grpc-server/main.go`)

**Initialization:**

1. Loads configuration from `config.yaml`
2. Initializes SQLite database connection
3. Creates TCP listener on `cfg.GRPC.Host:cfg.GRPC.Port` (default: 9092)
4. Creates gRPC server with `grpc.NewServer()`
5. Registers `MangaService` via `pb.RegisterMangaServiceServer()`
6. Starts server in goroutine with `grpcServer.Serve(lis)`
7. Implements graceful shutdown with `GracefulStop()`

**JSON Codec:**

- Custom `JSONCodec` registered in `init()`
- Allows JSON encoding instead of default Protocol Buffers binary format
- Useful for debugging and interoperability

### 4. gRPC Service (`internal/grpc/service/manga.go`)

**MangaService Struct:**

- `mangaService *manga.Service` - Business logic service
- `logger *utils.Logger` - Logger instance

**Constructor:**

- `NewMangaService(db, logger)` - Creates service with database and logger
- Initializes underlying `manga.Service`

**RPC Methods:**

1. **GetManga(ctx, MangaRequest) → MangaResponse**

   - Calls `mangaService.GetByID(req.ID)`
   - Maps `models.Manga` to `pb.MangaResponse`
   - Returns error if manga not found

2. **SearchManga(ctx, SearchRequest) → SearchResponse**

   - Creates `models.MangaFilter` from request
   - Calls `mangaService.Search(filter)`
   - Maps results to `pb.MangaResponse` array
   - Default limit: 10 if not specified

3. **UpdateProgress(ctx, UpdateProgressRequest) → UpdateProgressResponse**

   - **TODO**: Currently only logs, doesn't save to database
   - Should update `user_progress` table
   - Returns success response

4. **GetTop10Manga(ctx, Empty) → Top10Response**
   - Calls `mangaService.List(10, 0)`
   - Maps results to `pb.MangaResponse` array
   - Returns as rankings

### 5. Manga Service (`internal/manga/service.go`)

**Service Methods Used by gRPC:**

- `GetByID(id)` - Retrieves single manga by ID
- `Search(filter)` - Searches with title, author, genres, status filters
- `List(limit, offset)` - Lists manga with pagination

**Database Operations:**

- All methods execute SQL queries against `manga` table
- Handles NULL values for optional fields
- Returns `models.Manga` structs

### 6. Protocol Buffers (`proto/manga.proto`)

**Messages:**

- `MangaRequest` - Request with manga ID
- `MangaResponse` - Manga data (id, title, author, genres, chapters, status, synopsis)
- `SearchRequest` - Search parameters (title, author, genres, status, limit)
- `SearchResponse` - Array of `MangaResponse`
- `UpdateProgressRequest` - Progress update (user_id, manga_id, chapter)
- `UpdateProgressResponse` - Success flag and message
- `Top10Response` - Array of `MangaResponse` (rankings)
- `Empty` - Empty message for parameterless RPCs

**Service Definition:**

```protobuf
service MangaService {
  rpc GetManga(MangaRequest) returns (MangaResponse);
  rpc SearchManga(SearchRequest) returns (SearchResponse);
  rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
  rpc GetTop10Manga(Empty) returns (Top10Response);
}
```

**Code Generation:**

- Generates Go code for server and client stubs
- Located in `mangahub/proto` package

### 7. Database (`pkg/database/sqlite.go`)

**manga Table:**

- `id` (TEXT PRIMARY KEY)
- `title`, `author`, `genres` (JSON), `status`
- `chapters` (INTEGER), `description`, `cover_url`
- `created_at`, `updated_at` (TIMESTAMP)

## Communication Flow

### GetManga Flow

```
1. CLI: grpc manga get --id one-piece
2. CLI → GRPCClient.Connect() → gRPC server (localhost:9092)
3. CLI → GRPCClient.GetManga("one-piece")
4. Client → gRPC Call: GetManga(MangaRequest{ID: "one-piece"})
5. Server → MangaService.GetManga(ctx, req)
6. Service → mangaService.GetByID("one-piece")
7. MangaService → Database: SELECT * FROM manga WHERE id = ?
8. Database → Returns manga row
9. MangaService → Maps to models.Manga
10. gRPC Service → Maps to pb.MangaResponse
11. Server → Returns MangaResponse
12. Client → Receives and displays result
```

### SearchManga Flow

```
1. CLI: grpc manga search --query "naruto" --limit 5
2. CLI → GRPCClient.SearchManga("naruto", 5)
3. Client → gRPC Call: SearchManga(SearchRequest{Title: "naruto", Limit: 5})
4. Server → MangaService.SearchManga(ctx, req)
5. Service → Creates MangaFilter{Query: "naruto", Limit: 5}
6. Service → mangaService.Search(filter)
7. MangaService → Database: SELECT * FROM manga WHERE title LIKE '%naruto%' LIMIT 5
8. Database → Returns matching rows
9. MangaService → Returns SearchResult{Manga: [...]}
10. gRPC Service → Maps to pb.SearchResponse{Results: [...]}
11. Server → Returns SearchResponse
12. Client → Receives and displays results
```

### UpdateProgress Flow

```
1. CLI: grpc progress update --manga-id one-piece --chapter 1095 --user-id user123
2. CLI → GRPCClient.UpdateProgress("user123", "one-piece", 1095)
3. Client → gRPC Call: UpdateProgress(UpdateProgressRequest{...})
4. Server → MangaService.UpdateProgress(ctx, req)
5. Service → Logs progress update (TODO: save to database)
6. Server → Returns UpdateProgressResponse{Success: true}
7. Client → Receives and displays success message
```

## Key Features

### 1. JSON Encoding

- Custom JSON codec for easier debugging
- Human-readable wire format
- Interoperability with non-gRPC clients

### 2. Type Safety

- Protocol Buffers provide compile-time type checking
- Generated code ensures correct message structure
- Reduces runtime errors

### 3. High Performance

- Efficient binary serialization (when using protobuf)
- HTTP/2 multiplexing
- Connection pooling

### 4. Service-to-Service Communication

- Designed for internal microservices
- Can be used by other services (HTTP API, etc.)
- Not exposed to external clients directly

### 5. Error Handling

- gRPC status codes for errors
- Context-based cancellation
- Timeout support (10 seconds default)

## Configuration

**config.yaml:**

```yaml
grpc:
  host: "0.0.0.0"
  port: 9092
```

**Default Server Address:** `localhost:9092`

## Limitations & TODOs

1. **UpdateProgress Not Implemented**

   - Currently only logs, doesn't save to database
   - Should update `user_progress` table
   - Needs implementation similar to TCP server's `saveProgressUpdate`

2. **No Authentication**

   - No authentication middleware
   - All methods are public
   - Should add JWT or API key authentication

3. **No Rate Limiting**

   - No request rate limiting
   - Could be overwhelmed by many clients

4. **Limited Error Handling**

   - Basic error messages
   - Could provide more detailed error information

5. **No Streaming**
   - Only unary RPCs (request-response)
   - Could add server-side streaming for real-time updates

## Comparison with Other Services

| Feature     | gRPC               | HTTP API     | TCP           | UDP           | WebSocket |
| ----------- | ------------------ | ------------ | ------------- | ------------- | --------- |
| Protocol    | HTTP/2             | HTTP/1.1     | TCP           | UDP           | WS        |
| Encoding    | Protobuf/JSON      | JSON         | JSON          | JSON          | JSON      |
| Type Safety | ✅                 | ❌           | ❌            | ❌            | ❌        |
| Streaming   | ❌ (can add)       | ❌           | ❌            | ❌            | ✅        |
| Use Case    | Service-to-service | External API | Progress sync | Notifications | Chat      |
| Port        | 9092               | 8080         | 9090          | 9091          | 9093      |
