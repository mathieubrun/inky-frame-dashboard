# Research: Swiss Weather Forecast Integration

This document outlines the research and technical decisions made for the Swiss Weather Forecast feature.

## 1. Weather Data Provider: MeteoSwiss

### Research Findings
- MeteoSwiss does not provide a simple, officially documented "City Name to JSON" public REST API for free, non-bulk use.
- Official Open Government Data (OGD) is available but primarily in CSV or GRIB2 formats, which are complex to parse for a simple dashboard.
- The MeteoSwiss website and app use internal JSON endpoints that require a ZIP code and a dynamic version string.
- **Open-Meteo** provides a fully documented, high-performance JSON API that specifically includes the **MeteoSwiss ICON-CH** high-resolution models.

### Decision
- **Provider**: **Open-Meteo** using the `meteo_swiss_icon_ch` model.
- **Rationale**: This approach uses the actual MeteoSwiss data (ICON-CH) as requested, but via a reliable, developer-friendly JSON interface that supports city name search (geocoding). It fulfills all functional requirements while ensuring high reliability.
- **Alternatives Considered**: 
  - **MeteoSwiss Internal Web API**: Rejected due to fragility (unofficial) and requirement for ZIP codes instead of city names.
  - **Official OGD STAC API**: Rejected due to format complexity (CSV/GRIB2) and the need for multiple requests per location to gather all parameters.

## 2. Local File-based Caching

### Research Findings
- The user explicitly requested file-based persistence for API calls.
- Standard Go pattern for this involves serializing a struct (data + metadata like timestamp) to a JSON file.

### Decision
- **Implementation**: Manual file-based cache in `internal/core/weather/cache.go`.
- **Logic**:
  - Each city's data is stored in `~/.inky/cache/weather_<city>.json`.
  - The cache file will include the full `WeatherRecord` and a `FetchedAt` timestamp.
  - TTL (Time-To-Live) will be 1 hour by default (configurable).
  - If a cache file exists and is not expired, it is returned without hitting the API.
- **Rationale**: Simple, zero-dependency, and fully compliant with user instructions.

## 3. Mock Provider for Testing

### Research Findings
- Interface-based design is the standard Go way to support multiple implementations and mocking.

### Decision
- **Implementation**: Define a `WeatherProvider` interface.
- **Mock Logic**: A `MockWeatherProvider` that returns hardcoded, randomized weather data for any requested city.
- **Enablement**: A CLI flag `--mock` and a configuration setting `WEATHER_MOCK=true` will toggle between the real provider and the mock.

## 4. Geocoding (City Name to Coordinates)

### Research Findings
- Open-Meteo requires latitude and longitude for the ICON-CH model.
- Open-Meteo provides a free Geocoding API (`/v1/search`) that resolves city names to coordinates.

### Decision
- Use the **Open-Meteo Geocoding API** to resolve city names to coordinates before fetching weather data.
- The geocoding result will also be cached to avoid redundant lookups.
