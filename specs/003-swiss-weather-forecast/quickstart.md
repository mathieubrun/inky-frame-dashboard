# Quickstart: Swiss Weather Forecast Integration

This document provides immediate usage instructions for the Swiss weather feature.

## 1. CLI Usage

### Basic Command
```bash
inky weather --city "Zurich"
```

### Mocking for Development
To use the mock provider instead of the live API:
```bash
inky weather --city "Zurich" --mock
```

### JSON Output (for scripts)
```bash
inky weather --city "Zurich" --json | jq .current.temperature
```

## 2. API Usage

### Start the Server
```bash
inky serve --port 8080
```

### Fetch Weather via HTTP
```bash
curl "http://localhost:8080/weather/swiss?city=Geneva"
```

## 3. Configuration

The feature can be configured via environment variables or CLI flags:

- `WEATHER_CACHE_DIR`: Location for the local file-based cache (default: `~/.inky/cache`).
- `WEATHER_CACHE_TTL`: How long to keep cached data (default: `1h`).
- `WEATHER_MOCK`: Set to `true` to use mock provider by default.

## 4. Troubleshooting

- **City Not Found**: Ensure you are using a valid Swiss city name.
- **Cache Permission Error**: Ensure the directory specified in `WEATHER_CACHE_DIR` is writable.
- **Network Error**: Verify connectivity to Open-Meteo (`api.open-meteo.com`).
