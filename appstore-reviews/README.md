# appstore-reviews

Polls the App Store RSS feed for reviews, stores them to disk, and serves them over HTTP. A small React app displays them.

## How it works

**Backend** — a poller fetches each app's RSS feed on an interval and saves new reviews to disk; an HTTP API reads them back.

```
poller ──fetch──> appstore ──parse──> store (data/*.json)
                                         │
                                api ─────┘──> GET /reviews/{appID}
```

- `cmd/server` — wires it together, starts the poller + HTTP server
- `internal/appstore` — fetches & parses the RSS feed
- `internal/poller` — background loop, polls on an interval
- `internal/store` — JSON file storage, dedup by review ID
- `internal/api` — HTTP handler, filters to the review window
- `internal/review` / `internal/config` — domain type & config loading

**Frontend** — one React component fetches `/reviews/{appID}` and renders the cards.

```
App.tsx ──fetch──> GET /reviews/{appID} ──> <ReviewCard /> list
```

## Requirements

- Go 1.22+
- Node + pnpm (for the frontend)

## Run

**Backend** — from `appstore-reviews/`:

```bash
go run ./cmd/server
```

Serves on `http://localhost:8080`.

**Frontend** — in another terminal, from `frontend/`:

```bash
pnpm install && pnpm dev
```

Open the URL Vite prints

## Configure

Edit `config.json`:

```json
{
  "addr": ":8080",
  "pollIntervalSeconds": 15,
  "reviewWindowHours": 720,
  "appIDs": ["595068606", "284882215"]
}
```

- `addr` — HTTP listen address
- `pollIntervalSeconds` — how often to fetch the feed
- `reviewWindowHours` — how far back the API returns reviews
- `appIDs` — App Store app IDs to poll

Reviews are persisted as JSON under `data/` and reloaded on restart.

## API

```
GET /reviews/{appID}
```

Returns reviews for that app within the configured window which is 720 because 48 hours didn't have enough data, sorted by newest first.

## Test

```bash
go test ./...
```

## What I left out for simplicity

A few things I deliberately chose not to do here to keep it minimal. In production code I would add them:

- RSS edge cases & pagination
- Memory/disk consistency on write failure
- Corrupted/empty data files
- Proper CORS (restricted origins, methods, headers)
- Graceful HTTP shutdown
- Config validation
- HTTP client timeout
- An interface for the store

## How I used AI

I used AI on this README, I wrote it loosely with my own words and asked AI to organize it.

I started by hand-creating the minimal folder/file structure I thought this needed. Then I tested fetching the Apple RSS endpoint directly, pasted the raw response into Claude, and asked it to generate the Go struct for it.

From there I built the backend in small, testable increments — one per commit. For each step I described how I wanted that part to work, had Claude generate it, then tested and trimmed: AI tends to go overboard, so I kept removing anything unnecessary to stay scoped to the minimum. Early steps were intentionally not dynamic or scalable — some hardcoded values — just enough to confirm each piece worked. At the end I made it dynamic and wired the separate parts together. Each commit is a small, working increment and the design only became dynamic and fully connected at the last few steps.

### Frontend

Since I only needed a small frontend just to display the reviews, I generated it with `npm create vite@latest frontend -- --template react`, which gave me a minimal structure with everything I needed already wired up. Then I used AI to generate a display with a minimal but decent UI, and manually added the code to fetch the backend with plain `fetch` instead of the what I usually use which is Axios, since the goal was to keep libraries minimal. After that I deleted the extra unnecessary files and code that came with the CLI boilerplate.

### Testing

Because I didn't want to go too far over the 3h limit on this, I listed the places I thought would be the minimal that would benefit from tests and asked AI to generate minimal tests for just those:

- `Save` dedup
- `Get` sort order
- load/reload round-trips
- parsing with bad rating or bad timestamp inputs in `client.go`
- the time-window filtering logic
- the app ID validation
- CORS headers

For production I would definitely add more, for example:

- the poller's loop, `FetchReviews` against a mock HTTP server, reviews mapping end to end, config loading, a missing or malformed `config.json`, and that the seconds/hours values convert to the right `time.Duration`


