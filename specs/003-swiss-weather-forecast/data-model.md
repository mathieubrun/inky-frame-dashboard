# Data Model: Swiss Weather Forecast Integration

This document defines the entities and data structures for the Swiss weather feature.

## 1. Weather Data Structs

### `WeatherRecord` (internal/core/weather/types.go)
Represents the weather state at a specific point in time.

| Field | Type | Description |
|---|---|---|
| `Timestamp` | `time.Time` | Time of the record (UTC) |
| `Temperature` | `float64` | Air temperature in Celsius |
| `WindSpeed` | `float64` | Wind speed in km/h |
| `WindDirection` | `float64` | Wind direction in degrees (0-360) |
| `Precipitation` | `float64` | Total rain/snow amount in mm |
| `PrecipitationProb` | `float64` | Probability of precipitation (0-100%) |

### `WeatherForecast` (internal/core/weather/types.go)
Represents the current weather and the 24-hour hourly forecast.

| Field | Type | Description |
|---|---|---|
| `Location` | `Location` | Geographic information for the data |
| `Current` | `WeatherRecord` | Most recent weather measurement |
| `Hourly` | `[]WeatherRecord` | 24-hour sequence of forecast records |
| `FetchedAt` | `time.Time` | When the data was retrieved (for cache logic) |

### `Location` (internal/core/weather/types.go)
Represents a geographic point.

| Field | Type | Description |
|---|---|---|
| `City` | `string` | Name of the city (e.g., "Zurich") |
| `Latitude` | `float64` | Geographic latitude |
| `Longitude` | `float64` | Geographic longitude |
| `Country` | `string` | Country code (always "CH" for this feature) |

## 2. Persistence Model (Caching)

The cache will be stored in `~/.inky/cache/` (or a configured directory).

### `CacheFile` (JSON structure)
Each file corresponds to a single city (`weather_<city_normalized>.json`).

```json
{
  "city": "zurich",
  "lat": 47.3769,
  "lon": 8.5417,
  "fetched_at": "2026-03-05T14:30:00Z",
  "forecast": {
    "current": { ... },
    "hourly": [ ... ]
  }
}
```

## 3. Provider Interface

```go
package weather

type Provider interface {
    GetForecast(city string) (*WeatherForecast, error)
}
```

Implementations:
- `MeteoSwissProvider` (Real API)
- `MockProvider` (Hardcoded data)
- `CachedProvider` (Wrapper for file-based caching)
