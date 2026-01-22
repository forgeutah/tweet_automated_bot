# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an automated Twitter bot that tweets YouTube videos from the GoWest Conference on a scheduled basis. It manages multiple Twitter accounts (gowestconf and forgeutahbot), selecting random videos from a database and ensuring videos aren't repeated within a 3-month window.

## Development Commands

### Building and Testing
```bash
# Build the project
go build -v ./...

# Run tests (excluding examples)
go test `go list ./... | grep -v examples`

# Run tests with coverage
go test `go list ./... | grep -v examples` -coverprofile=coverage.txt -covermode=atomic
```

### Linting
```bash
# Run golangci-lint (uses config from .golangci.yml)
golangci-lint run

# The project uses these linters: bodyclose, deadcode, dogsled, errcheck, goconst,
# gocyclo, gofmt, gosimple, govet, importas, ineffassign, misspell, revive,
# rowserrcheck, sqlclosecheck, staticcheck, structcheck, stylecheck, typecheck, unused, varcheck
```

### Running the Application
```bash
# Set up authentication credentials first
source .secret  # or set CREDENTIAL_FILE env var

# Required environment variables:
# - CREDENTIAL_FILE: path to JSON file with Twitter OAuth credentials
# - DB_CLUSTER_ID: CockroachDB cluster ID
# - DB_USERNAME: database username
# - DB_PASSWORD: database password
# - DB_HOST: database host
# - DB_NAME: database name
# - PORT: (optional) HTTP server port, defaults to 8080

# Run the bot
go run main.go
```

## Architecture

### Core Components

**Main Application (main.go)**: Orchestrates the bot lifecycle. Creates two AutoBot instances (gowestconf and forgeutahbot) that run concurrently via goroutines. Starts an HTTP health check server on port 8080 (or PORT env var). Handles graceful shutdown via signal handling.

**AutoBot (internal/botguts/youtubevids.go)**: The scheduler responsible for automated tweeting. Each bot instance:
- Operates on a 7-day timer (time.Hour * 24 * 7)
- Selects random YouTube videos from the database
- Formats tweets with video title, presenter Twitter handle, conference year hashtag, and URL
- Updates the database to track when videos were last sent

**Client Package (client/)**: Wraps Twitter API interactions using dghubble/go-twitter. Supports multiple Twitter accounts via a map of clients keyed by twitter handle. Uses OAuth1 authentication. Credentials are loaded from a JSON file with structure:
```json
{
  "twitter_clients": [
    {
      "twitter_handle": "handle",
      "consumer_key": "key",
      "consumer_secret": "secret",
      "access_token": "token",
      "access_secret": "secret"
    }
  ]
}
```

**Database Package (db/)**: Manages PostgreSQL/CockroachDB connections and operations. Key tables:
- `yt_videos`: Stores YouTube video metadata (title, URL, playlist, conference year, presenter Twitter, last_sent_at)
- `tweets`: Tracks tweet status (queued, sent, failed)

The database connection downloads a root certificate from 1Password at connection time and removes it on close. The Migrate() method runs on startup to create tables and seed initial video data if needed.

### Video Selection Logic

Videos are selected randomly but filtered to ensure they haven't been sent in the last 3 months (db/video.go:79-80). The SelectOneRandomVideo query orders by random() and iterates through results until finding a valid video that passes the 3-month check.

### HTTP Server

Exposes a `/health` endpoint that returns "we are live" for uptime monitoring and container health checks.

## Database Schema

**yt_videos table**:
- twitter_username (VARCHAR 15)
- video_title (VARCHAR 255)
- video_playlist (VARCHAR 255)
- video_url (text)
- conference_year (VARCHAR 4)
- presenter_twitter_username (VARCHAR 64)
- created_at (timestamp)
- last_sent_at (timestamp)

**tweets table**:
- twitter_username (VARCHAR 15)
- tweet_text (text)
- links (text)
- send_time (timestamp)
- status (enum: queued, sent, failed)
- created_at (timestamp)

## Dependencies

Core dependencies managed via go.mod:
- github.com/dghubble/go-twitter: Twitter API client
- github.com/dghubble/oauth1: OAuth1 authentication
- github.com/jmoiron/sqlx: SQL extensions for Go
- github.com/lib/pq: PostgreSQL driver

## Testing

Test files use JSON fixtures for Twitter client setup (test_twitter_single_valid.json, etc.). Tests require database environment variables to be set even if not actively connecting.

## Notes

- The bot uses Go 1.18
- CockroachDB is used as the database backend with SSL certificate verification
- The scheduler runs indefinitely with weekly intervals
- Each bot instance operates independently on its assigned Twitter account
- Videos are never deleted, only their last_sent_at timestamp is updated
