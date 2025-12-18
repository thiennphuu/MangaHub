# MangaHub - Network Programming Project Specification

## Course Information

- **Course**: Net-centric Programming (IT096IU)
- **Programming Language**: Go
- **Team Size**: 2 students per group
- **Timeline**: 10-11 weeks
- **Instructor**: Lê Thanh Sơn - Nguyễn Trung Nghĩa

## Project Objectives

- Gain practical experience in network application development using Go
- Implement and understand all five required communication protocols (TCP, UDP, HTTP, gRPC, WebSocket)
- Strengthen understanding of networking concepts through hands-on implementation
- Develop foundational skills in concurrent programming and distributed systems
- Create a working system demonstrating network programming competency

## Project Overview: MangaHub

MangaHub is a manga tracking system that demonstrates network programming concepts through practical implementation. The system uses all five required protocols in a cohesive application.

### Core Requirements

- **Implementation Language**: Go
- **Required Protocols**: TCP, UDP, HTTP, gRPC, WebSocket
- **Database**: SQLite
- **Data Format**: JSON

## Deliverables

1. **Source Code**: All implementation files submitted on Blackboard
2. **Documentation**: Comprehensive technical documentation
3. **Demonstration**: Live demonstration showing all five protocols working together
4. **File Naming**: `GroupXX_MangaHub.zip` (e.g., Group01_MangaHub.zip)

### Submission Deadline

- **Final Submission**: 23:59 on demo day
- **Demo Session**: Scheduled separately
- **Note**: Failure to attend demonstration results in ZERO grading

## Manga Database Requirements

### Data Collection

- **Manual Entry**: 100 popular manga series with essential metadata
- **API Integration**: 100 additional series from MangaDx API or other legal sources
- **Educational Practice**: Limited web scraping from practice sites
- **Storage Format**: JSON

### Minimum Database Coverage

- At least 30-40 different manga series across major genres
- At least 15-20 series per major genre (shounen, shoujo, seinen, josei, etc.)
- Basic metadata per series:
  - Title
  - Author
  - Genres
  - Status
  - Chapter count
  - Description

### Data Structure Example

```json
{
  "id": "one-piece",
  "title": "One Piece",
  "author": "Oda Eiichiro",
  "genres": ["Action", "Adventure", "Shounen"],
  "status": "ongoing",
  "total_chapters": 1100,
  "description": "A young pirate's adventure...",
  "cover_url": "https://example.com/covers/one-piece.jpg"
}
```

## System Architecture

### 1. HTTP REST API Server (25 points)

Basic RESTful service with essential endpoints.

**Essential Endpoints**:
- `POST /auth/register` - User registration
- `POST /auth/login` - User authentication
- `GET /manga` - Search manga with basic filters
- `GET /manga/{id}` - Get manga details
- `POST /users/library` - Add manga to library
- `GET /users/library` - Get user's library
- `PUT /users/progress` - Update reading progress

**Requirements**:
- JWT-based authentication
- SQLite database integration
- JSON request/response handling
- Basic error handling and logging
- Input validation

### 2. TCP Progress Sync Server (20 points)

Simple TCP server for basic progress broadcasting.

**Features**:
- Accept multiple TCP connections
- Broadcast progress updates to connected clients
- Handle client connections and disconnections
- Basic JSON message protocol
- Concurrent connection handling with goroutines

**Message Structure**:
```go
type ProgressUpdate struct {
    UserID    string `json:"user_id"`
    MangaID   string `json:"manga_id"`
    Chapter   int    `json:"chapter"`
    Timestamp int64  `json:"timestamp"`
}
```

### 3. UDP Notification System (15 points)

Basic UDP broadcaster for chapter notifications.

**Features**:
- UDP server listening for client registrations
- Broadcast chapter release notifications
- Client list management
- Basic error logging

**Message Structure**:
```go
type Notification struct {
    Type      string `json:"type"`
    MangaID   string `json:"manga_id"`
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}
```

### 4. WebSocket Chat System (15 points)

Simple real-time chat for manga discussions.

**Features**:
- WebSocket connection handling
- Real-time message broadcasting
- User join/leave functionality
- Basic connection management

**Message Structure**:
```go
type ChatMessage struct {
    UserID    string `json:"user_id"`
    Username  string `json:"username"`
    Message   string `json:"message"`
    Timestamp int64  `json:"timestamp"`
}
```

### 5. gRPC Internal Service (10 points)

Simple gRPC service for internal communication.

**Services**:
- `GetManga(GetMangaRequest) returns (MangaResponse)`
- `SearchManga(SearchRequest) returns (SearchResponse)`
- `UpdateProgress(ProgressRequest) returns (ProgressResponse)`

**Requirements**:
- Protocol Buffer definitions for 2-3 services
- Basic gRPC server implementation
- Simple client integration
- Unary RPC calls

### 6. Database Layer (10 points)

SQLite database with proper schema and relationships.

**Tables**:
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE manga (
    id TEXT PRIMARY KEY,
    title TEXT,
    author TEXT,
    genres TEXT,
    status TEXT,
    total_chapters INTEGER,
    description TEXT
);

CREATE TABLE user_progress (
    user_id TEXT,
    manga_id TEXT,
    current_chapter INTEGER,
    status TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id)
);
```

## Network Communication Requirements

### Protocol Implementation Standards

#### HTTP Services
- RESTful API with proper HTTP methods and status codes
- Basic JWT authentication
- Error handling with appropriate HTTP responses
- Simple CORS support for web clients

#### TCP Socket Communication
- Basic server accepting multiple connections
- JSON-based message protocol
- Concurrent connection handling with goroutines
- Graceful connection termination

#### UDP Broadcasting
- Simple UDP server for notifications
- Client registration mechanism
- Basic broadcast functionality
- Error handling for network failures

#### gRPC Services
- Protocol Buffer message definitions
- Basic service implementation
- Client-server communication
- Simple error handling

#### WebSocket Connections
- WebSocket upgrade handling
- Real-time message broadcasting
- Connection lifecycle management
- Basic client management

## Performance Requirements

### Scalability Targets

- Support 50-100 concurrent users during testing
- Handle 30-40 manga series in database
- Process basic search queries within 500ms
- Support 20-30 concurrent TCP connections
- WebSocket chat with 10-20 simultaneous users

### Reliability Standards

- 80-90% uptime during demonstration period
- Basic error handling and recovery
- Simple logging for debugging
- Graceful degradation when services unavailable

## Development Timeline (10 Weeks)

### Phase 1: Foundation (Weeks 1-3)

**Week 1: Project Setup & HTTP Basics**
- Go project structure setup
- Basic HTTP server with Gin framework
- User registration and login endpoints
- SQLite database setup

**Week 2: Core HTTP API**
- Manga data model and CRUD endpoints
- User library management endpoints
- Basic JWT authentication middleware
- API testing and validation

**Week 3: Data Collection & Integration**
- Manual manga data entry (20-30 series)
- Simple MangaDx API integration
- Data validation and storage
- API endpoint completion

### Phase 2: Network Protocols (Weeks 3-7)

**Week 3-4: TCP Implementation**
- Basic TCP server setup
- Connection handling with goroutines
- Simple message protocol design
- Progress update broadcasting

**Week 5: UDP Notification System**
- UDP server implementation
- Client registration mechanism
- Basic notification broadcasting
- Integration testing

**Week 6: WebSocket Chat**
- WebSocket server setup
- Basic chat functionality
- Connection management
- Real-time message broadcasting

**Week 7: gRPC Service**
- Protocol Buffer definitions
- Basic gRPC service implementation
- Client integration
- Service testing

### Phase 3: Integration & Testing (Weeks 8-10)

**Week 8: System Integration**
- Connect all protocols together
- End-to-end testing
- Bug fixes and stability improvements

**Week 9: User Interface & Documentation**
- Simple web interface (optional)
- API documentation
- Code documentation and comments

**Week 10: Demo Preparation**
- Demo script preparation
- Live demonstration practice
- Final code review and cleanup

## Recommended Go Libraries

### Core Framework
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/golang-jwt/jwt/v4` - JWT authentication
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/mattn/go-sqlite3` - SQLite database driver

### gRPC
- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers

### Testing
- `github.com/stretchr/testify` - Testing utilities

## Grading Criteria (Total: 100 points)

### Core Protocol Implementation (40 points)

- **HTTP REST API** (15 pts): Complete endpoints with authentication and database integration
- **TCP Progress Sync** (13 pts): Working server with concurrent connections and broadcasting
- **UDP Notifications** (5 pts): Basic notification system with client management
- **WebSocket Chat** (10 pts): Real-time messaging with connection handling
- **gRPC Service** (7 pts): Basic service with 2-3 working methods

### System Integration & Architecture (20 points)

- **Database Integration** (8 pts): Working data persistence with proper schema
- **Service Communication** (7 pts): All protocols integrated and working seamlessly
- **Error Handling & Logging** (3 pts): Comprehensive error handling across components
- **Code Structure & Organization** (2 pts): Proper Go project organization and modularity

### Code Quality & Testing (10 points)

- **Go Code Quality** (5 pts): Proper Go idioms and error handling patterns
- **Testing Coverage** (3 pts): Unit tests for core functionality
- **Code Documentation** (2 pts): Clear comments and function documentation

### Documentation & Demo (10 points)

- **Technical Documentation** (5 pts): API docs, setup instructions, architecture
- **Live Demonstration** (5 pts): Successfully demonstrate all five protocols with Q&A

## Bonus Features (Extra Credit - Up to 20 points)

### Advanced Protocol Features (5-10 points)

- **Enhanced TCP Synchronization** (10 pts): Implement conflict resolution
- **WebSocket Room Management** (10 pts): Multiple chat rooms for different topics
- **UDP Delivery Confirmation** (5 pts): Implement acknowledgment system
- **gRPC Streaming** (10 pts): Add server-side streaming for real-time updates

### Enhanced Data Management (5-10 points)

- **Advanced Search & Filtering** (5 pts): Full-text search with multiple filters
- **Data Caching with Redis** (10 pts): Redis for frequently accessed data
- **Recommendation System** (10 pts): Collaborative filtering based on patterns

### Social & Community Features (5-10 points)

- **User Reviews & Ratings** (8 pts): User reviews and ratings system
- **Friend System** (5 pts): Add/remove friends and view activity
- **Reading Lists Sharing** (6 pts): Share reading lists with other users
- **Activity Feed** (7 pts): Recent activities from friends

### Performance & Scalability (5-10 points)

- **Connection Pooling** (6 pts): Proper connection pooling
- **Rate Limiting** (5 pts): API endpoint rate limiting
- **Horizontal Scaling** (8 pts): Support multiple server instances
- **Performance Monitoring** (7 pts): Metrics collection and monitoring
- **Load Balancing** (10 pts): Load balancing for service instances

### Advanced User Features (5-12 points)

- **Reading Statistics** (8 pts): Detailed reading analytics
- **Notification Preferences** (5 pts): Customizable settings per user
- **Reading Goals & Achievements** (10 pts): Goals and achievements system
- **Data Export/Import** (10 pts): Export to JSON/CSV, import from services
- **Multiple Reading Lists** (5 pts): Custom reading lists

### API & Integration Enhancements (5-10 points)

- **External API Integration** (10 pts): Integrate with additional manga APIs
- **Webhook System** (10 pts): Notifications via webhooks
- **API Versioning** (10 pts): Proper versioning with backward compatibility
- **OpenAPI Documentation** (5 pts): Interactive API documentation
- **Mobile-Optimized Endpoints** (10 pts): Specialized mobile endpoints

### Security & Reliability (5-10 points)

- **Advanced Authentication** (10 pts): Refresh tokens and session management
- **Input Sanitization** (5 pts): Comprehensive input validation
- **Automated Backups** (10 pts): Automated database backup system
- **Health Checks** (5 pts): Health check endpoints
- **Graceful Shutdown** (10 pts): Graceful shutdown for all servers

### Development & Deployment (5-10 points)

- **Docker Compose Setup** (10 pts): Complete containerization
- **CI/CD Pipeline** (10 pts): Automated testing and deployment
- **Environment Configuration** (5 pts): Environment-based configuration
- **Database Migrations** (7 pts): Automated schema migration system
- **Monitoring & Alerting** (8 pts): System monitoring with alerts

## Success Criteria - Minimum Requirements

### Must Have Features

1. All five network protocols implemented and functional
2. Basic user authentication and authorization
3. Manga data storage and retrieval
4. Progress tracking and synchronization
5. Real-time chat functionality
6. Successful live demonstration

## Expected Learning Outcomes

- Understanding of network programming concepts in Go
- Experience with concurrent programming using goroutines
- Knowledge of different communication protocols and use cases
- Basic distributed system integration skills
- Foundation for advanced network programming concepts

## Project Structure

```
mangahub/
├── cmd/
│   ├── api-server/main.go
│   ├── tcp-server/main.go
│   ├── udp-server/main.go
│   └── grpc-server/main.go
├── internal/
│   ├── auth/
│   ├── manga/
│   ├── user/
│   ├── tcp/
│   ├── udp/
│   ├── websocket/
│   └── grpc/
├── pkg/
│   ├── models/
│   ├── database/
│   └── utils/
├── proto/
├── data/
├── docs/
├── docker-compose.yml
└── README.md
```

## Regulations on AI Chatbot Usage

### Permitted Uses

- Brainstorming ideas and exploring programming approaches
- Language refinement, grammar checking, and summarization
- Clarifying complex concepts as study guides
- Requesting explanations of syntax errors

### Prohibited Uses

- Submitting AI-generated code without meaningful modification
- Using AI to solve entire project tasks
- Bypassing learning objectives with AI assistance

### Transparency Requirements

- Acknowledge if AI tools were used
- Describe briefly how AI assistance was applied
- Student teams bear full responsibility for accuracy and originality

---

**Total Maximum Points**: 120 points (100 core + 20 bonus)

**Final Grade Calculation**: `min(Total Points, 100)` for the 30% course component
