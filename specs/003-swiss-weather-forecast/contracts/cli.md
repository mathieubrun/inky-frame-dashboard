# CLI Contract: Swiss Weather Forecast

This document defines the CLI output for the `inky weather` command.

## Command: `inky weather`

Retrieves and displays weather for a Swiss city.

### Arguments

| Argument | Description |
|---|---|
| `--city "Name"` | The Swiss city to fetch weather for. |
| `--mock` | Use mock data instead of live data. |
| `--json` | Output the raw JSON response instead of a formatted table. |

### Default Text Output

```text
Weather for Zurich, CH (MeteoSwiss)
------------------------------------
Current: 12.5°C | Wind: 15.2 km/h (SW) | Rain: 0.0 mm (10%)

Next 24h Forecast:
- 15:00: 11.8°C | Rain: 0.2 mm (25%)
- 16:00: 11.2°C | Rain: 0.5 mm (40%)
- 17:00: 10.5°C | Rain: 0.1 mm (30%)
...
(Fetched from cache: 2026-03-05 14:35:00Z)
```

### JSON Output (`--json`)

Should match the API response structure defined in `api.md`.
```json
{
  "location": { ... },
  "current": { ... },
  "hourly": [ ... ],
  "fetched_at": "...",
  "source": "..."
}
```
