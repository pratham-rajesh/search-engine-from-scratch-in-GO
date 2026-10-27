# search-engine-from-scratch-in-GO

A full-text web search engine written in Go from scratch: a polite crawler, a Wikipedia importer,
a positional inverted index with a custom compressed on-disk format, BM25/TF-IDF ranking, and an
HTTP serving layer with caching, rate limiting, autocomplete and analytics.

> **Status:** work in progress. This README will grow into the
> architecture overview, benchmarks and design-decision summary.

## Quickstart (development)

Requires Go 1.27+ and Docker.

```powershell
# Start Postgres
docker compose -f deploy/docker-compose.yml up -d

# Run the server
$env:SEARCH_DATABASE_URL = "postgres://search:search@localhost:5432/search"
go run ./cmd/searchd

# In another terminal
curl.exe localhost:8080/healthz

# Tests
go test ./...
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `SEARCH_PORT` | `8080` | HTTP listen port (1-65535) |
| `SEARCH_DATABASE_URL` | _required_ | Postgres connection string |
| `SEARCH_LOG_LEVEL` | `info` | `debug`, `info`, `warn` or `error` |
