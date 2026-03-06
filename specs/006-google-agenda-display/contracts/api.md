# API Contract: Google Agenda & Dashboard

This document defines the new endpoints and CLI commands for the Google Agenda integration.

## API Endpoints

### GET `/api/v1/agenda`
Retrieves upcoming calendar events in JSON format.

**Query Parameters**:
- `calendar_id` (string, optional): Overrides the default configured calendar.
- `count` (int, default=10): Number of events to retrieve.

**Response (200 OK)**:
```json
{
  "events": [
    {
      "summary": "Team Sync",
      "start_time": "2026-03-06T10:00:00Z",
      "end_time": "2026-03-06T11:00:00Z",
      "location": "Meeting Room A"
    }
  ],
  "fetched_at": "2026-03-06T09:00:00Z"
}
```

### GET `/api/v1/dashboard/image`
Retrieves the combined weather and agenda image (800x480).

**Query Parameters**:
- `location` (string, required): Weather location.
- `calendar_id` (string, optional): Agenda calendar ID.
- `palette` (string, default="spectra6"): Target display palette.

**Response (200 OK)**:
- `Content-Type`: `image/png`
- `Body`: Binary PNG data.

## CLI Interface

### Command: `inky agenda list`
Displays upcoming events in a table format.

**Flags**:
- `--count`, `-n`: Number of events (Default: 10).
- `--mock`: Use mock agenda data.

### Command: `inky dashboard image`
Generates and saves the combined dashboard PNG.

**Flags**:
- `--location`, `-l`: Weather location (Required).
- `--output`, `-o`: File path (Default: `dashboard.png`).
- `--mock`: Use mock weather and agenda data.
