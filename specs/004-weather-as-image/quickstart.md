# Quickstart: Weather Data Image Generation

This guide provides instructions for developers to test and verify the weather image generation feature.

## Prerequisites

- Go (Golang) 1.22+ installed.
- Internet access for fetching weather data from Open-Meteo.
- An image viewer that supports PNG files.

## Running the API Server

1. **Start the server**:
```bash
go run cmd/inky/main.go serve
```

2. **Fetch a weather image using `curl`**:
```bash
curl -o weather_zurich.png "http://localhost:8080/api/v1/weather/image?location=Zurich"
```

3. **Verify the image**:
Open `weather_zurich.png` to confirm it contains weather information and matches the 800x480 resolution.

## Using the CLI

1. **Generate and save an image locally**:
```bash
go run cmd/inky/main.go weather image "Zurich" --output weather_test.png
```

2. **Generate an image with custom resolution**:
```bash
go run cmd/inky/main.go weather image "8001" --width 400 --height 240 --output weather_small.png
```

## Testing with Bruno

1. Open the [Bruno app](https://usebruno.com/).
2. Load the `bruno/` collection from this repository.
3. Select the `Get Weather Image` request.
4. Update the `location` query parameter if needed and click **Send**.
5. Bruno will display the binary PNG response in the Preview tab.

## Caching Behavior

- Images are cached on disk for 15 minutes.
- To force a refresh, delete the corresponding file in the `cache/` directory (or wherever the cache is configured) or wait for expiration.
