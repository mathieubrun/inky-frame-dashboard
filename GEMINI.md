# inky-frame-dashboard Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-05

## Active Technologies
- Go (Golang) 1.22+ (Managed by `go mod`) + standard library `net/http`, `spf13/cobra`, `spf13/viper` (003-get-app-version)
- Local JSON files for weather data caching. (003-swiss-weather-forecast)
- MicroPython (Pimoroni Firmware for RP2040) + inky_frame, picographics, pngdec, network, urequests, machine (ADC) (005-inky-display-image)
- Internal Flash (for script and config) (005-inky-display-image)
- File-based JSON cache for calendar events. (006-google-agenda-display)
- Local CSV file (e.g., `data/battery.csv`) (007-battery-level-monitoring)

- Go (Golang) 1.22+ + net/http, spf13/cobra, viper (003-get-app-version)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go (Golang) 1.22+

## Code Style

Go (Golang) 1.22+: Follow standard conventions

## Recent Changes
- 008-battery-refresh-optimization: Added Go (Golang) 1.22+ + net/http, spf13/cobra, viper
- 007-battery-level-monitoring: Added Go (Golang) 1.22+ + net/http, spf13/cobra, viper
- 006-google-agenda-display: Added Go (Golang) 1.22+


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
