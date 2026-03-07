# Research: Battery Level Monitoring

## Findings

### CSV Storage and Concurrency
- **Location**: `.inky/battery.csv` (aligns with other cache files in the project).
- **Concurrency**: Go's `os.OpenFile` with `os.O_APPEND|os.O_WRONLY|os.O_CREATE` is generally atomic for small writes on Linux. However, to ensure total reliability within the application (especially if multiple devices are added), a `sync.Mutex` in the `core/battery` package will be used.
- **Format**: `Timestamp,Voltage` (e.g., `2026-03-07T12:00:00Z,3.75`). RFC3339 is standard and easy to parse.

### API Routes
- **Endpoint 1**: `POST /api/v1/battery`
  - Payload: `{"voltage": 3.75}`
  - Status: 201 Created on success, 400 Bad Request on invalid input.
- **Endpoint 2**: `GET /api/v1/battery/history`
  - Response: Raw CSV text (Content-Type: `text/plain`).

### MicroPython Implementation
- Use `urequests.post(url, json={"voltage": v})`.
- The voltage can be read using `machine.ADC` on Inky Frame.
- Integration: Add `report_battery()` call in `main.py` before `get_dashboard_image()`.

### CLI Command
- `inky battery history`: Prints the content of the CSV file.
- `inky battery clear`: (Optional) Deletes the CSV file.

## Decisions

### Decision: Mutex-guarded CSV appending
- **Rationale**: While a single device is unlikely to cause conflicts, using a Mutex in the core logic ensures the API is thread-safe as per Go best practices.
- **Alternatives considered**: SQLite (overkill for this), Single-threaded handler (not idiomatic for Go net/http).

### Decision: JSON for POST payload
- **Rationale**: Easiest to handle in both MicroPython and Go, and allows for future expansion (e.g., adding device IDs).
- **Alternatives considered**: Form-urlencoded, URL query parameters.

### Decision: Raw text for history
- **Rationale**: Direct export of the CSV content satisfies the requirement for "raw text" and is easy for users/admins to analyze.
- **Alternatives considered**: JSON list of objects (more complex for MVP).
