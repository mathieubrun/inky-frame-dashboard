# Research: Weather Data Image Generation

This document outlines the technical decisions and research findings for implementing the weather image generation feature.

## Image Generation Library

**Decision**: Use `github.com/fogleman/gg` (Go Graphics).

**Rationale**:
- Simple, high-level API for drawing onto an image canvas.
- Built on top of the Go standard `image` package.
- Supports TrueType fonts and basic shapes, making it ideal for layout generation.
- Lightweight compared to alternatives like `Cairo` (which requires CGO).

**Alternatives Considered**:
- `golang.org/x/image/font`: Too low-level for rapid layout development.
- `disintegration/imaging`: Great for resizing/filtering, but lacks drawing capabilities.

## Color Palette & Display Optimization

**Decision**: Target the 6-color Spectra 6 palette and apply dithering if necessary.

**Palette**:
| Color  | RGB Value |
| :----- | :-------- |
| Black  | (0, 0, 0) |
| White  | (255, 255, 255) |
| Red    | (255, 0, 0) |
| Green  | (0, 255, 0) |
| Blue   | (0, 0, 255) |
| Yellow | (255, 255, 0) |

**Optimization Strategy**:
- Use high-contrast layouts.
- Boost saturation and contrast of weather icons before rendering.
- For 6-color e-ink, specific pigments are used; the server will output a standard PNG, but colors will be constrained to this palette to ensure fidelity on the physical device.

## Geocoding (City/Postcode Lookup)

**Decision**: Utilize the existing `OpenMeteoProvider.geocode` logic, extending it to support postcodes.

**Rationale**:
- Already implemented and tested in `internal/core/weather/openmeteo.go`.
- Open-Meteo Geocoding API supports postal codes in its search parameter.
- Consolidates providers to a single external service.

## Weather Icons

**Decision**: Use a set of embedded PNG icons optimized for the Inky palette.

**Rationale**:
- Simple to load and render with `gg`.
- Can be pre-processed to match the target color depth.
- Avoids runtime SVG rendering complexity.

## Caching Strategy

**Decision**: Implement a file-based cache in `internal/core/weather/cache.go`.

**Rationale**:
- Meets the 2.5s response goal by serving pre-rendered images.
- Persistent across application restarts.
- Simple to implement using the existing cache patterns in the codebase.
- Cache key: `city_name_resolution_palette.png`.
- TTL: 15 minutes (aligned with weather data freshness).
