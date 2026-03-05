# API Contract: Swiss Weather Forecast

This document defines the interface for the `/weather/swiss` endpoint.

## GET `/weather/swiss`

Retrieves current weather and a 24-hour hourly forecast for a Swiss city.

### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `city` | `string` | Yes | Name of the Swiss city (e.g., "Zurich") |
| `mock` | `boolean` | No | If "true", return mock data instead of live data |

### Success Response (200 OK)

```json
{
  "location": {
    "city": "Zurich",
    "latitude": 47.3769,
    "longitude": 8.5417,
    "country": "CH"
  },
  "current": {
    "timestamp": "2026-03-05T14:30:00Z",
    "temperature": 12.5,
    "wind_speed": 15.2,
    "wind_direction": 220.0,
    "precipitation": 0.0,
    "precipitation_prob": 10.0
  },
  "hourly": [
    {
      "timestamp": "2026-03-05T15:00:00Z",
      "temperature": 11.8,
      "wind_speed": 14.5,
      "wind_direction": 215.0,
      "precipitation": 0.2,
      "precipitation_prob": 25.0
    },
    ...
  ],
  "fetched_at": "2026-03-05T14:35:00Z",
  "source": "MeteoSwiss (ICON-CH via Open-Meteo)"
}
```

### Error Responses

#### 400 Bad Request
If `city` parameter is missing.
```json
{ "error": "missing city parameter" }
```

#### 404 Not Found
If the city cannot be resolved by the geocoder.
```json
{ "error": "city not found: UnknownCity" }
```

#### 503 Service Unavailable
If the external weather provider is down.
```json
{ "error": "weather service unavailable" }
```
