# MangaHub CLI - Testing Guide

## Quick Start

### Build the CLI (Optional)
```powershell
go build -o mangahub.exe ./cmd/cli
```

### Run Commands Directly
```powershell
go run ./cmd/cli --help
go run ./cmd/cli --version
```

## Command Testing

### Authentication
```powershell
# Register new user
go run ./cmd/cli auth register --username john --email john@example.com

# Login (prompts for password)
go run ./cmd/cli auth login --username john

# Check status
go run ./cmd/cli auth status

# Logout
go run ./cmd/cli auth logout
```

### Manga Search
```powershell
# Basic search
go run ./cmd/cli manga search "attack on titan"

# Search with filters
go run ./cmd/cli manga search "action" --genre action --status ongoing --limit 5

# View details
go run ./cmd/cli manga info attack-on-titan

# List manga
go run ./cmd/cli manga list --page 1 --limit 20

# List with filters
go run ./cmd/cli manga list --status ongoing --limit 10
go run ./cmd/cli manga list --genre shounen --limit 15

# Advanced search
go run ./cmd/cli manga advanced-search --genre shounen --status completed --limit 10
go run ./cmd/cli manga advanced-search "romance" --min-chapters 50 --sort-by title --order asc
```

### Library Management
```powershell
# Add to library
go run ./cmd/cli library add --manga-id one-piece --status reading --rating 9

# View library
go run ./cmd/cli library list

# Filter library
go run ./cmd/cli library list --status reading
go run ./cmd/cli library list --status completed --sort-by title --order asc

# Update entry
go run ./cmd/cli library update --manga-id one-piece --status completed --rating 10

# Remove from library
go run ./cmd/cli library remove --manga-id one-piece
```

### Progress Tracking
```powershell
# Update progress
go run ./cmd/cli progress update --manga-id naruto --chapter 700

# With additional info
go run ./cmd/cli progress update --manga-id one-piece --chapter 1095 --notes "Epic chapter!"

# View progress history
go run ./cmd/cli progress history
```

### Notifications
```powershell
# Subscribe
go run ./cmd/cli notify subscribe

# View preferences
go run ./cmd/cli notify preferences

# Test notification
go run ./cmd/cli notify test

# Unsubscribe
go run ./cmd/cli notify unsubscribe
```

### Chat
```powershell
# Join general chat
go run ./cmd/cli chat join

# Join manga-specific chat
go run ./cmd/cli chat join --room one-piece

# Send message
go run ./cmd/cli chat send "Hello everyone!"

# View history
go run ./cmd/cli chat history
go run ./cmd/cli chat history --room one-piece --limit 50
```

### gRPC Operations
```powershell
# Query manga
go run ./cmd/cli grpc manga get --id one-piece

# Search
go run ./cmd/cli grpc manga search --query "naruto"

# Update progress
go run ./cmd/cli grpc progress update --manga-id one-piece --chapter 1095
```

### Statistics
```powershell
# Overview
go run ./cmd/cli stats overview

# Detailed
go run ./cmd/cli stats detailed
```

### Data Export
```powershell
# Export library
go run ./cmd/cli export library --format json --output my-library.json

# Export progress
go run ./cmd/cli export progress --format csv --output progress.csv

# Full backup
go run ./cmd/cli export all --output mangahub-backup.tar.gz
```

### Configuration
```powershell
# View config
go run ./cmd/cli config show

# Set value
go run ./cmd/cli config set api-url http://localhost:8080

# Reset
go run ./cmd/cli config reset
```

### Server Management
```powershell
# Check server health
go run ./cmd/cli server health

# View server logs
go run ./cmd/cli server logs
```

### Database Management
```powershell
# Check database
go run ./cmd/cli db check

# Database stats
go run ./cmd/cli db stats

# Repair database
go run ./cmd/cli db repair
```

## Testing Output Examples

### Successful Search
```
Searching for "attack on titan"...

Found 1 results:

┌──────────────────────────────────────────────────────────────────────────────────────────┐
│ ID   │ TITLE                          │ AUTHOR               │ STATUS     │ CHAPTERS │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│ a... │ Attack on Titan                │ Hajime Isayama       │ completed  │      141 │
└──────────────────────────────────────────────────────────────────────────────────────────┘

Use 'mangahub manga info <id>' to view details
Use 'mangahub library add --manga-id <id>' to add to your library
```

### Manga List
```
Listing manga (page 1, limit 15)

┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ ID     │ TITLE                               │ AUTHOR                    │ STATUS     │ CHAPTERS │
├────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ 20t... │ 20th Century Boys                   │ Naoki Urasawa             │ completed  │      249 │
│ one... │ One Piece                           │ Eiichiro Oda              │ ongoing    │     1100 │
│ nar... │ Naruto                              │ Masashi Kishimoto         │ completed  │      700 │
└────────────────────────────────────────────────────────────────────────────────────────────────────────┘

Showing 15 manga (page 1)
Use --page <n> to see more results
```

### Manga Info
```
┌──────────────────────────────────────────────────────────────────────┐
│ ONE PIECE                                                            │
└──────────────────────────────────────────────────────────────────────┘

Basic Information:
  ID:      one-piece
  Title:   One Piece
  Author:  Eiichiro Oda
  Genres:  Action, Adventure, Comedy, Shounen
  Status:  ongoing

Progress:
  Total Chapters: 1100
  Created:        2025-12-17
  Updated:        2025-12-17

Description:
  Monkey D. Luffy sets off on an adventure to find the legendary One Piece
  treasure and become the Pirate King.

Actions:
  Update Progress: mangahub progress update --manga-id one-piece --chapter <n>
  Add to Library:  mangahub library add --manga-id one-piece
  Rate/Review:     mangahub library update --manga-id one-piece --rating <1-10>
```

### Auth Status
```
Authentication Status: ✓ Logged in

User ID:  user_abc123
Username: johndoe
Email:    john@example.com
Expires:  2025-12-19 09:30:00 UTC
```

## Global Flags Testing

### All commands support:
```powershell
# With auth token
go run ./cmd/cli manga search "naruto" --token your-token-here

# With custom API URL
go run ./cmd/cli manga search "naruto" --api http://localhost:8080

# With verbose output
go run ./cmd/cli manga search "naruto" --verbose

# Show help
go run ./cmd/cli manga search --help
```

## Error Testing

### Invalid flags
```powershell
go run ./cmd/cli library add                    # Error: --manga-id required
go run ./cmd/cli progress update                # Error: --manga-id and --chapter required
go run ./cmd/cli auth login                     # Error: --username or --email required
```

### Command help
```powershell
go run ./cmd/cli --help                         # Show all commands
go run ./cmd/cli manga --help                   # Show manga subcommands
go run ./cmd/cli manga search --help            # Show search flags
go run ./cmd/cli auth --help                    # Show auth subcommands
```

## Integration Testing Checklist

- [x] All commands execute without syntax errors
- [x] Help text displays correctly
- [x] Flags are parsed correctly
- [x] Global flags work on all commands
- [x] Error messages are user-friendly
- [x] Output formatting is readable
- [x] Manga search queries database
- [x] Manga list with pagination works
- [x] Manga info displays details
- [x] Auth register creates users
- [x] Auth login validates credentials
- [x] Auth status shows session info
- [ ] Library management (needs user login)
- [ ] Progress tracking (needs user login)
- [ ] Chat functionality (needs WebSocket server)
- [ ] gRPC operations (needs gRPC server)

## Running Servers

```powershell
# API Server (HTTP REST) - Port 8080
go run ./cmd/api-server

# gRPC Server - Port 9092
go run ./cmd/grpc-server

# TCP Server - Port 9090
go run ./cmd/tcp-server

# UDP Server - Port 9091
go run ./cmd/udp-server

# WebSocket Server - Port 9093
go run ./cmd/websocket-server
```

## Quick Test Commands

```powershell
# Test manga functionality
go run ./cmd/cli manga list --limit 10
go run ./cmd/cli manga search "one piece"
go run ./cmd/cli manga info one-piece

# Test auth functionality
go run ./cmd/cli auth status
go run ./cmd/cli auth register --username testuser --email test@example.com
go run ./cmd/cli auth login --username testuser

# Test advanced search
go run ./cmd/cli manga advanced-search --genre shounen --status completed --limit 5
```

---

**Database**: 200 manga entries (100 manual + 100 API)
**Status**: CLI connected to SQLite database
**Config**: config.yaml for server settings
