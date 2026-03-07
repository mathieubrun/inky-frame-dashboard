# Data Model: Battery Refresh Optimization

## Entities

### DashboardState
Represents the current state of the dashboard as delivered to a device.

| Field | Type | Description |
|-------|------|-------------|
| ETag | string | MD5 hash of the rendered PNG image. |
| LastRefresh | Time | When the last image was successfully generated. |
| WeatherFreshUntil | Time | The calculated next 04:00 AM timestamp. |

## Storage Strategy

### Server-Side (Go)
- **ETag**: Calculated on-the-fly or cached in memory for the last request.
- **Weather Cache**: Persisted in `.inky/cache/weather_*.json`. The `FetchedAt` timestamp is used to check against the "Last 04:00 AM" logic.

### Client-Side (MicroPython)
- **ETag File**: `/etag.txt` on internal flash. Contains the last received ETag string.

## State Transitions
1. **Request Received**: API checks `If-None-Match` header.
2. **Freshness Check**: Server checks if weather data is stale (passed 04:00 AM) or calendar data is stale (passed TTL).
3. **Image Generation**: If stale, generate new image and new ETag.
4. **ETag Comparison**:
   - If new ETag == `If-None-Match`: Return 304.
   - If new ETag != `If-None-Match`: Return 200 + image + new ETag header.
5. **Device Update**: If 200, device saves new ETag and refreshes screen.
