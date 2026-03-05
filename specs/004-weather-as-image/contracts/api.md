# API Contract: Weather Image Endpoint

This document defines the interface for fetching pre-rendered weather images.

## Endpoint: GET `/api/v1/weather/image`

Retrieves a weather forecast image for a specific location, optimized for e-ink displays.

### Request Parameters (Query)

| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `location` | `string` | **Yes** | N/A | City name or postcode (e.g., `Zurich`, `8001`). |
| `width` | `int` | No | `800` | Target image width. |
| `height` | `int` | No | `480` | Target image height. |
| `palette` | `string` | No | `spectra6` | The target display palette (Options: `spectra6`, `grayscale`). |

### Success Response

- **Status Code**: `200 OK`
- **Content-Type**: `image/png`
- **Body**: Binary PNG image data.

### Error Responses

#### Missing Location
- **Status Code**: `400 Bad Request`
- **Content-Type**: `application/json`
- **Body**:
```json
{
  "error": "missing location parameter"
}
```

#### Location Not Found
- **Status Code**: `404 Not Found`
- **Content-Type**: `application/json`
- **Body**:
```json
{
  "error": "location not found: Zurich"
}
```

#### Upstream Service Error
- **Status Code**: `502 Bad Gateway`
- **Content-Type**: `application/json`
- **Body**:
```json
{
  "error": "failed to fetch weather data from source"
}
```

## CLI Interface

The CLI will also support generating and saving weather images locally.

### Command: `inky weather image`

```bash
inky weather image [location] --width 800 --height 480 --output weather.png
```

| Flag | Shorthand | Description |
| :--- | :--- | :--- |
| `--width` | `-w` | Target width (Default: 800). |
| `--height` | `-h` | Target height (Default: 480). |
| `--output` | `-o` | Output file path (Default: `weather.png`). |
| `--palette` | `-p` | Target palette (Default: `spectra6`). |
