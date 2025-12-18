# MangaHub CLI Application - User Manual

## Table of Contents

1. [Introduction](#introduction)
2. [Installation and Setup](#installation-and-setup)
3. [Getting Started](#getting-started)
4. [Authentication Commands](#authentication-commands)
5. [Manga Management](#manga-management)
6. [Library Operations](#library-operations)
7. [Network Protocol Features](#network-protocol-features)
8. [Chat System](#chat-system)
9. [Configuration](#configuration)
10. [Troubleshooting](#troubleshooting)
11. [Advanced Features](#advanced-features)

## Introduction

MangaHub CLI is a command-line interface for the MangaHub manga tracking system. It provides access to all core features including manga discovery, reading progress tracking, real-time synchronization, and community chat functionality.

### System Requirements

- Go 1.19 or later
- SQLite 3.x
- Network connectivity for synchronization features
- Terminal with UTF-8 support

### Supported Platforms

- Linux (x64, ARM)
- macOS (Intel, Apple Silicon)
- Windows (x64)

## Installation and Setup

### Download and Install

```bash
# Download the latest release
wget https://github.com/yourorg/mangahub/releases/latest/mangahub-cli

# Make executable (Linux/macOS)
chmod +x mangahub-cli

# Move to system path
sudo mv mangahub-cli /usr/local/bin/mangahub

# Verify installation
mangahub version
```

### First-Time Setup

```bash
# Initialize configuration
mangahub init

# This creates:
# ~/.mangahub/config.yaml
# ~/.mangahub/data.db
# ~/.mangahub/logs/
```

### Configuration File

The default configuration is created at `~/.mangahub/config.yaml`:

```yaml
server:
  host: "localhost"
  http_port: 8080
  tcp_port: 9090
  udp_port: 9091
  grpc_port: 9092
  websocket_port: 9093

database:
  path: "~/.mangahub/data.db"

user:
  username: ""
  token: ""

sync:
  auto_sync: true
  conflict_resolution: "last_write_wins"

notifications:
  enabled: true
  sound: false

logging:
  level: "info"
  path: "~/.mangahub/logs/"
```

## Getting Started

### Quick Start Guide

```bash
# 1. Start the MangaHub server
mangahub server start

# 2. In another terminal, register a new account
mangahub auth register --username myuser --email myuser@example.com

# 3. Login to get authentication token
mangahub auth login --username myuser

# 4. Search for manga
mangahub manga search "one piece"

# 5. Add manga to library
mangahub library add --manga-id one-piece --status reading

# 6. Update reading progress
mangahub progress update --manga-id one-piece --chapter 1095
```

### Command Structure

All commands follow the pattern:

```
mangahub <command> <subcommand> [flags] [arguments]
```

### Global Flags

- `--config`: Specify config file path
- `--verbose`: Enable verbose output
- `--quiet`: Suppress non-error output
- `--help`: Show help information

## Authentication Commands

### Register New Account

```bash
mangahub auth register --username <username> --email <email>
# Prompts for password securely
```

**Example:**
```bash
mangahub auth register --username johndoe --email john@example.com
```

**Expected Output:**
```
Password: [hidden input]
Confirm password: [hidden input]
✓ Account created successfully!
User ID: usr_1a2b3c4d5e
Username: johndoe
Email: john@example.com
Created: 2024-01-20 10:30:00 UTC

Please login to start using MangaHub:
 mangahub auth login --username johndoe
```

### Login

```bash
mangahub auth login --username <username>
# Prompts for password

# Alternative: login with email
mangahub auth login --email <email>
```

**Expected Output:**
```
Password: [hidden input]
✓ Login successful!
Welcome back, johndoe!

Session Details:
 Token expires: 2024-01-21 10:30:00 UTC (24 hours)
 Permissions: read, write, sync

Auto-sync: enabled
Notifications: enabled

Ready to use MangaHub! Try:
 mangahub manga search "your favorite manga"
```

### Logout

```bash
mangahub auth logout
# Removes stored authentication token
```

### Check Authentication Status

```bash
mangahub auth status
# Shows current login status and user information
```

### Change Password

```bash
mangahub auth change-password
# Prompts for current password and new password
```

## Manga Management

### Search Manga

```bash
# Basic search
mangahub manga search <query>

# Search with filters
mangahub manga search <query> --genre <genre> --status <status>

# Examples
mangahub manga search "attack on titan"
mangahub manga search "romance" --genre romance --status completed
mangahub manga search "naruto" --limit 5
```

### View Manga Details

```bash
mangahub manga info <manga-id>

# Example
mangahub manga info one-piece
```

### List All Manga

```bash
# List all manga in database
mangahub manga list

# List with pagination
mangahub manga list --page 2 --limit 20

# Filter by genre
mangahub manga list --genre shounen
```

### Advanced Search Options

```bash
mangahub manga search "keyword" \
  --genre "action,adventure" \
  --status "ongoing" \
  --author "author name" \
  --year-from 2020 \
  --year-to 2024 \
  --min-chapters 50 \
  --sort-by "popularity" \
  --order "desc"
```

## Library Operations

### Add Manga to Library

```bash
mangahub library add --manga-id <id> --status <status>
# Status options: reading, completed, plan-to-read, on-hold, dropped

# Examples
mangahub library add --manga-id one-piece --status reading
mangahub library add --manga-id death-note --status completed --rating 9
```

### View Library

```bash
# View entire library
mangahub library list

# Filter by status
mangahub library list --status reading
mangahub library list --status completed

# Sort options
mangahub library list --sort-by title
mangahub library list --sort-by last-updated --order desc
```

### Remove from Library

```bash
mangahub library remove --manga-id <id>

# Example
mangahub library remove --manga-id completed-series
```

### Update Library Entry

```bash
mangahub library update --manga-id <id> --status <new-status>

# Example
mangahub library update --manga-id one-piece --status completed --rating 10
```

### Batch Operations

```bash
# Batch add manga to library
mangahub library batch-add --file manga-list.txt --status plan-to-read

# Batch update progress
mangahub progress batch-update --file progress-updates.csv
```

## Progress Tracking

### Update Reading Progress

```bash
mangahub progress update --manga-id <id> --chapter <number>

# With additional info
mangahub progress update --manga-id <id> --chapter <number> --volume <number>

# Examples
mangahub progress update --manga-id one-piece --chapter 1095
mangahub progress update --manga-id naruto --chapter 700 --volume 72 --notes "Great ending!"
```

### View Progress History

```bash
mangahub progress history --manga-id <id>

# View all progress updates
mangahub progress history
```

### Sync Progress

```bash
# Manual sync with server
mangahub progress sync

# Check sync status
mangahub progress sync-status
```

## Network Protocol Features

### TCP Progress Synchronization

```bash
# Connect to TCP sync server
mangahub sync connect

# Disconnect from sync server
mangahub sync disconnect

# Check connection status
mangahub sync status

# View real-time progress updates
mangahub sync monitor
```

### UDP Notifications

```bash
# Subscribe to chapter release notifications
mangahub notify subscribe

# Unsubscribe from notifications
mangahub notify unsubscribe

# View notification preferences
mangahub notify preferences

# Test notification system
mangahub notify test
```

### gRPC Service Operations

```bash
# Query manga via gRPC
mangahub grpc manga get --id <manga-id>

# Search via gRPC
mangahub grpc manga search --query <search-term>

# Update progress via gRPC
mangahub grpc progress update --manga-id <id> --chapter <number>
```

## Chat System

### Connect to Chat

```bash
# Join general chat
mangahub chat join

# Join specific manga discussion
mangahub chat join --manga-id <id>

# Example
mangahub chat join --manga-id one-piece
```

### Send Messages

```bash
# Send message to current chat
mangahub chat send "Hello everyone!"

# Send message to specific manga chat
mangahub chat send "Great chapter!" --manga-id one-piece
```

### Chat Commands (Interactive Mode)

When in chat mode, use these commands:

- `/help` - Show chat commands
- `/users` - List online users
- `/quit` - Leave chat
- `/pm <username> <message>` - Private message
- `/manga <id>` - Switch to manga-specific chat
- `/history` - Show recent history
- `/status` - Connection status

### View Chat History

```bash
# View recent messages
mangahub chat history

# View messages for specific manga
mangahub chat history --manga-id one-piece --limit 50
```

## Statistics and Analytics

### Reading Statistics

```bash
# View personal reading statistics
mangahub stats overview

# Detailed breakdown
mangahub stats detailed

# Stats for specific time period
mangahub stats --from 2024-01-01 --to 2024-12-31
```

## Export Data

```bash
# Export library to JSON
mangahub export library --format json --output library.json

# Export reading progress
mangahub export progress --format csv --output progress.csv

# Full data export
mangahub export all --output mangahub-backup.tar.gz
```

## Server Management

### Start Server Components

```bash
# Start all servers
mangahub server start

# Start specific servers
mangahub server start --http-only
mangahub server start --tcp-only
mangahub server start --udp-only
```

### Check Server Status

```bash
# Check server status
mangahub server status

# Detailed health check
mangahub server health
```

### Stop Servers

```bash
# Stop all servers
mangahub server stop

# Stop specific server
mangahub server stop --component http
```

### View Server Logs

```bash
# View server logs
mangahub server logs

# Follow logs in real-time
mangahub server logs --follow

# Filter logs by level
mangahub server logs --level error
```

## Configuration

### View Configuration

```bash
# Show current configuration
mangahub config show

# Show specific section
mangahub config show server
```

### Update Configuration

```bash
# Set configuration value
mangahub config set server.host "192.168.1.100"
mangahub config set notifications.enabled false

# Reset to defaults
mangahub config reset
```

## Profile Management

```bash
# Create new profile
mangahub profile create --name work

# Switch profiles
mangahub profile switch --name work

# List profiles
mangahub profile list
```

## Backup and Restore

```bash
# Create backup
mangahub backup create --output backup-2024.tar.gz

# Restore from backup
mangahub backup restore --input backup-2024.tar.gz
```

## Database Operations

```bash
# Database integrity check
mangahub db check

# Optimize database
mangahub db optimize

# Database statistics
mangahub db stats

# Repair database
mangahub db repair
```

## Troubleshooting

### Authentication Problems

```bash
# Clear authentication data
mangahub auth clear

# Re-register if needed
mangahub auth register --username <username> --email <email>
```

### Connection Issues

```bash
# Test server connectivity
mangahub server ping

# Reset network connections
mangahub sync reconnect
```

### Database Issues

```bash
# Repair database
mangahub db repair

# Reinitialize if needed
mangahub init --force
```

### Debug Mode

```bash
# Run with debug logging
mangahub --verbose <command>

# Enable trace logging
mangahub config set logging.level trace
```

### Log Analysis

```bash
# View error logs
mangahub logs errors

# Search logs
mangahub logs search "connection failed"

# Clear old logs
mangahub logs clean --older-than 30d
```

## Advanced Features

### Getting Help

```bash
# General help
mangahub help

# Command-specific help
mangahub manga help
mangahub library help
```

### Version Information

```bash
# Check version
mangahub version

# Check for updates
mangahub update check

# Update to latest version
mangahub update install
```

### System Information

```bash
# Display system information
mangahub system info

# Test system connectivity
mangahub system ping
```

---

**MangaHub CLI** - A comprehensive manga tracking system with multi-protocol support and real-time synchronization.

For additional features or specific use cases, refer to the built-in help system or consult the online documentation.
