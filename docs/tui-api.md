# AniCat TUI API Contract

> This document describes the JSON-mode TCP protocol for building a TUI client against the AniCat server.

---

## Transport

- **Protocol**: TCP
- **Default port**: `12314` (configurable via `env.yaml` → `port`)
- **Framing**: each message is a single line terminated by `\r\n`
- **Encoding**: JSON

---

## Request Format

Every request is a JSON object sent as one line:

```json
{"cmd": <int>, "arg": "<string>", "raw": <object|null>}
```

| Field | Type           | Description                                      |
|-------|----------------|--------------------------------------------------|
| `cmd` | int            | Command type (see table below)                   |
| `arg` | string         | Primary argument (name, sid, etc.)               |
| `raw` | object or null | Command-specific flags (see per-command section) |

### Command Types

| Value | Name     | Purpose                        |
|-------|----------|--------------------------------|
| 0     | Add      | Subscribe by anime name        |
| 1     | AddFeed  | Subscribe via RSS feed URL     |
| 2     | Remove   | Delete subscription            |
| 3     | Ls       | List all subscriptions         |
| 4     | LsItems  | Browse torrent resource list   |
| 5     | Status   | Get download progress          |
| 6     | Stop     | Shutdown server                |
| 7     | Rename   | Rename subject files           |

---

## Commands

### Ls — List Subscriptions

**Request**
```json
{"cmd": 3, "arg": "", "raw": null}
```

**Response** — JSON array of subscription objects:
```json
[
  {
    "sid": 3,
    "type": "feed",
    "name": "葬送的芙莉莲",
    "epi": "28",
    "status": "updating",
    "compl": "N"
  },
  {
    "sid": 2,
    "type": "single",
    "name": "孤独摇滚！",
    "epi": "12",
    "status": "fin",
    "compl": "Y"
  }
]
```

| Field    | Values                  | Description                          |
|----------|-------------------------|--------------------------------------|
| `sid`    | int                     | Subscription ID                      |
| `type`   | `"feed"` / `"single"`   | Subscription type                    |
| `name`   | string                  | Anime title                          |
| `epi`    | string (number or `"0"`)| Total episode count (0 = unknown)    |
| `status` | `"updating"` / `"fin"`  | Whether still airing                 |
| `compl`  | `"Y"` / `"N"`           | Whether fully downloaded and done    |

---

### Status — Download Progress

**Request**
```json
{"cmd": 5, "arg": "3", "raw": null}
```

`arg` is the subscription `sid` as a string.

**Response (qBittorrent mode)** — JSON object:
```json
{
  "path": "/bangumi/葬送的芙莉莲/",
  "total_size": "12GB",
  "torrents": [
    {"file": "葬送的芙莉莲 S01E25.mkv", "size": "350MB", "progress": "100%"},
    {"file": "葬送的芙莉莲 S01E26.mkv", "size": "348MB", "progress": "52%"},
    {"file": "葬送的芙莉莲 S01E27.mkv", "size": "349MB", "progress": "0%"}
  ]
}
```

| Field        | Type            | Description                     |
|--------------|-----------------|---------------------------------|
| `path`       | string          | Absolute download directory     |
| `total_size` | string          | Sum of all torrent sizes        |
| `torrents`   | array           | Per-file progress entries       |
| `torrents[].file`     | string | Filename                   |
| `torrents[].size`     | string | File size (e.g. `"350MB"`) |
| `torrents[].progress` | string | Completion (e.g. `"52%"`)  |

**Response (built-in downloader mode)**

When the subscription uses the built-in torrent client (`builtin-downloader: on`), the connection is **kept alive** and the server streams updates:

1. Server first sends: `keep-alive\r\n`
2. Then sends progress JSON every 15 seconds until complete:

```json
{"list": [{"name": "葬送的芙莉莲 S01E25.mkv", "progress": 100}, ...], "fin": false}
```

3. Final message has `"fin": true`, then server closes the connection.

> The TUI client should detect `keep-alive` and switch to streaming read mode.

---

### LsItems — Browse Torrent Resource List

**Request**
```json
{"cmd": 4, "arg": "葬送的芙莉莲", "raw": {"searchList": false}}
```

Set `"searchList": true` to return a flat item list instead of grouped by subtitle group.

**Response (searchList: false)** — grouped by subtitle group:
```json
{
  "SubsPlease": [
    {"name": "[SubsPlease] Frieren - 28 (1080p)", "size": "350MB", "uptime": "04-17 12:00"},
    {"name": "[SubsPlease] Frieren - 27 (1080p)", "size": "348MB", "uptime": "04-10 12:00"}
  ],
  "Erai-raws": [
    {"name": "[Erai-raws] 葬送的芙莉莲 - 28 [1080p]", "size": "345MB", "uptime": "04-17 13:00"}
  ]
}
```

**Response (searchList: true)** — flat array:
```json
[
  {"name": "[SubsPlease] Frieren - 28 (1080p)", "size": "350MB", "uptime": "04-17 12:00"},
  {"name": "[Erai-raws] 葬送的芙莉莲 - 28 [1080p]", "size": "345MB", "uptime": "04-17 13:00"}
]
```

---

### Add — Subscribe by Name

**Request**
```json
{"cmd": 0, "arg": "葬送的芙莉莲", "raw": {"mustContain": "1080p,简中", "mustNotContain": "外挂", "useRegexp": false, "group": "", "index": 0, "feedInfoName": ""}}
```

`raw` may be `null` if no filters needed.

**Response** — assigned `sid` as a string on success:
```json
"3"
```

On error: plain error message string.

---

### AddFeed — Subscribe via RSS Feed URL

**Request**
```json
{"cmd": 1, "arg": "https://mikanani.me/RSS/Bangumi?...", "raw": {"mustContain": "1080p", "mustNotContain": "", "useRegexp": false, "group": "", "index": 0, "feedInfoName": "葬送的芙莉莲"}}
```

**Response** — same as Add: assigned `sid` string or error.

---

### Remove — Delete Subscription

**Request** (single):
```json
{"cmd": 2, "arg": "3", "raw": null}
```

**Request** (all):
```json
{"cmd": 2, "arg": "*", "raw": null}
```

**Response**: `"ok"` or error string.

---

### Rename — Rename Subject Files

**Request**
```json
{"cmd": 7, "arg": "3", "raw": "葬送的芙莉莲"}
```

`arg` is the `sid`; `raw` is the new name string (JSON-encoded).

**Response**: `"ok"` or error string.

---

### Stop — Shutdown Server

**Request**
```json
{"cmd": 6, "arg": "", "raw": null}
```

**Response**: `"exited."`

---

## Error Handling

All errors are returned as plain text strings (not JSON objects), regardless of `format`:

```
subject not found
```

The client should treat any response that fails JSON parsing as an error message.

---

## Connection Lifecycle

```
client                      server
  |                           |
  |── connect TCP :12314 ────>|
  |── {"cmd":3,...}\r\n ─────>|
  |<── [...]\r\n ─────────────|
  |── close ─────────────────>|
```

Each command uses a **new TCP connection** (one request per connection), except for `Status` in built-in mode which stays open for streaming.
