# Data Model: Google Agenda Display

This document defines the data structures used for the Google Agenda integration.

## Entities

### Agenda Event
Represents a single event from a Google Calendar.

| Field | Type | Description |
| :--- | :--- | :--- |
| `summary` | `string` | The title of the event. |
| `start_time` | `time.Time` | When the event starts. |
| `end_time` | `time.Time` | When the event ends. |
| `location` | `string` | Where the event takes place (optional). |

### Agenda Forecast
Represents a collection of upcoming events for a specific period.

| Field | Type | Description |
| :--- | :--- | :--- |
| `events` | `[]AgendaEvent` | List of upcoming events. |
| `fetched_at` | `time.Time` | When the data was last updated from Google. |

### Dashboard Parameters
Represents the combined inputs for generating the full dashboard image.

| Field | Type | Description |
| :--- | :--- | :--- |
| `weather_req` | `ImageRequest` | Weather parameters (location, etc.). |
| `agenda_req` | `AgendaRequest` | Agenda parameters (calendar_id, count). |

## Relationships

- An **Agenda Forecast** contains multiple **Agenda Events**.
- A **Dashboard Image** is composed of one **Weather Forecast** and one **Agenda Forecast**.

## Validation Rules

- **Event Summary**: MUST NOT be empty.
- **Event Time**: Start time MUST be before end time.
- **Event Limit**: Maximum of 8 events for the dashboard display.
