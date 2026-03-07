# Research: Battery Refresh Optimization

## Findings

### ETag / 304 Not Modified
- **Server-side (Go)**:
  - Generate an ETag by hashing the rendered PNG data (MD5 is sufficient and fast).
  - Use `w.Header().Set("ETag", etag)` and check the `If-None-Match` header from the request.
  - If they match, return `http.StatusNotModified` (304) and no body.
- **Client-side (MicroPython)**:
  - `urequests` supports custom headers.
  - Store the last received ETag in a local file (e.g., `/etag.txt`).
  - Send `If-None-Match: <stored_etag>` in the request.
  - If response is 304, skip the screen refresh.
  - If response is 200, update the stored ETag from the `ETag` response header.

### Weather Refresh at 04:00
- **Logic**: Instead of a fixed TTL (e.g., 1 hour), the weather data should be considered "fresh" until the next 04:00 occurs.
- **Implementation**:
  - In `internal/core/weather/cache.go`, modify the freshness check.
  - A simple way: Calculate the "last 04:00" relative to now. If `FetchedAt` is before that "last 04:00", it is stale.
  - This ensures that after 04:00, the first request will trigger a refresh.

### Calendar Refresh
- The specification says "Calendar data is refreshed when calendar is updated."
- Currently, `agenda.CachedProvider` uses a fixed TTL (default 15 minutes).
- To support "when updated", we would ideally need webhooks, but since the Inky Frame polls every 60 minutes anyway, the current TTL approach is acceptable as long as it's shorter than the wake interval.
- **Optimization**: The server-side dashboard rendering will produce the same image (and thus same ETag) if the calendar data hasn't changed, fulfilling the requirement to limit refreshes.

## Decisions

### Decision: MD5 for ETag
- **Rationale**: Fast to calculate on the server and provides a unique signature for the image content.
- **Alternatives considered**: SHA1 (slower, not needed for this security level), Last-Modified header (less reliable for dynamic images).

### Decision: Local file `etag.txt` for MicroPython
- **Rationale**: Simple persistence across deep sleep cycles.
- **Alternatives considered**: Persistent variable in memory (lost during deep sleep).

### Decision: Daily 04:00 logic in Go
- **Rationale**: Centralizing the scheduling logic in the Go API avoids complex time calculations on the MicroPython device (which may have clock drift).
- **Alternatives considered**: Device-side scheduling.
