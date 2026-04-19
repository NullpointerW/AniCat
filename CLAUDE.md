# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
make windows           # Build Windows amd64 binary
make linux             # Build Linux amd64 binary
make windows-service   # Build Windows service binary
make all               # Build all targets + Docker image

# Run
go run main.go -e env.yaml          # Run with config file
go run main.go -e env.yaml -d       # Run with debug logging

# Test
go test ./...                        # Run all tests
go test ./subject/...                # Run a specific package's tests
go test ./test/...                   # Run integration tests

# Docker
make docker-build      # Build Docker image (wmooon/anicat)
make docker-push       # Push to registry
```

## Architecture

AniCat is an anime subscription/download automation tool. Users subscribe to anime titles; the app monitors RSS feeds, triggers downloads via qBittorrent (or a built-in torrent client), renames completed files for media library scrapers (Plex/Jellyfin), and sends email notifications.

**Entry point:** `main.go` starts four goroutines: scan loop, management server, download detector, and two network servers (TCP CLI + Web UI).

### Core Layers

**`subject/`** — The heart of the application.
- `subject.go`: The `Subject` struct (subscription model) and its operations.
- `manager.go`: Lifecycle management — creating and deleting subscriptions; persists state to JSON files.
- `runner.go`: The main polling loop — checks RSS, matches new episodes, dispatches downloads.
- `filter.go`: Episode filtering by regex, subtitle group, and keyword rules from `env.yaml`.
- `rename.go`: Post-download file renaming using TMDB metadata to standard `S01E01` format.
- `builtin_download.go`: Native torrent downloading via `anacrolix/torrent`.

**`crawl/`** — External data sources.
- `resource/mikan.go`: Scrapes Mikan RSS feeds for torrent links.
- `information/tmdb.go`, `information/bgmi.go`: Fetch anime metadata (title, season, episode count).
- `cover/`: Downloads poster images from TMDB and Douban.

**`downloader/`** — Download orchestration.
- `qbtcli.go`: qBittorrent Web API client.
- `rss/`: RSS feed reader that passes torrents directly to qBittorrent.
- `builtin/`: Wrapper around the `anacrolix/torrent` native client.
- `detector/`: Polls qBittorrent/built-in client to detect when downloads complete, then triggers renaming.

**`net/`** — TCP-based CLI server. Clients connect over TCP; `cmd/` handles commands like Add, Remove, List. Protocol defined in `transport.go`.

**`web/`** — HTTP REST API + static frontend served from `web/static/`.

**`conf/`** — `env.go` parses `env.yaml` into the `Environment` struct; `flag.go` handles `-e` (config path) and `-d` (debug) flags.

### Data Flow

```
User (CLI/Web) → Subscribe → Subject created & persisted (JSON)
                                    ↓
                          runner.go polling loop
                                    ↓
                    Mikan RSS scraped → episodes filtered
                                    ↓
                  qBittorrent API  OR  built-in torrent client
                                    ↓
                  detector/ polls completion
                                    ↓
                  rename.go renames files (TMDB S01E01 format)
                                    ↓
                       Email notification sent
```

### Configuration

`env.yaml` (path passed via `-e` flag) controls: TCP/Web ports, qBittorrent URL and credentials, download directory, global RSS filters, SMTP settings, HTTP proxies, and log path. The `conf.Environment` struct is the single source of truth passed through the app.

## Key Patterns

- Subscriptions are persisted as JSON files on disk (not a database); `manager.go` owns read/write.
- Platform-specific code (Windows service, path separators) is isolated in `service-exec/` and `utils/`.
- `errs/` wraps errors with context; use its helpers rather than `fmt.Errorf` for consistency.
- Tests in `test/` are integration-style; tests in `*_test.go` alongside source are unit tests.
