# Implementation Plan: Weather Data Image Generation

**Branch**: `004-weather-as-image` | **Date**: 2026-03-05 | **Spec**: [specs/004-weather-as-image/spec.md](spec.md)
**Input**: Feature specification from `/specs/004-weather-as-image/spec.md`

## Summary

The Inky Frame Dashboard requires a way to serve pre-rendered weather images to low-power e-ink displays. This feature implements an image generation engine using the `fogleman/gg` library, a new HTTP endpoint for image retrieval by city name or postcode, and a file-based caching mechanism to ensure fast response times and minimal server-side load. The images will be optimized for the Inky Frame's 6-color Spectra 6 palette at 800x480 resolution.

## Technical Context

**Language/Version**: Go (Golang) 1.22+
**Primary Dependencies**: net/http, spf13/cobra, viper, github.com/fogleman/gg
**Storage**: File-based PNG cache
**Testing**: standard library testing (Required, >80% coverage)
**Target Platform**: Raspberry Pi Pico W (Inky Frame) + Linux Server (API)
**Project Type**: web-service/cli
**Performance Goals**: < 2.5s response time for image generation/delivery
**Constraints**: < 200KB image size, 800x480 resolution, 6-color palette
**Scale/Scope**: 10 concurrent requests handled without failure

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Validation |
| :--- | :--- | :--- |
| I. Logic Offloading | ✅ PASS | All image rendering is performed on the Golang server. |
| II. Energy-First | ✅ PASS | Pre-rendered images minimize Inky Frame processing and Wi-Fi time. |
| III. Data Freshness | ✅ PASS | Images include update timestamps and indicate stale state. |
| IV. Resource-Conscious | ✅ PASS | Images optimized for 6-color palette and correct resolution. |
| V. API-First | ✅ PASS | Defined a stable `/api/v1/weather/image` contract. |
| VII. Modular & Unified | ✅ PASS | Functionality available via both `serve` and CLI commands. |
| IX. API Testing | ✅ PASS | Bruno collection included in `bruno/`. |
| X. Mandatory Testing | ✅ PASS | Targeting >80% coverage for rendering and caching logic. |

## Project Structure

### Documentation (this feature)

```text
specs/004-weather-as-image/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── api.md
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
cmd/
└── inky/     # Single application entry point (CLI & HTTP Server)

internal/
├── api/      # HTTP handlers, routes, and API-specific logic
├── cli/      # Cobra command definitions (including subcommands)
├── config/   # Configuration loading (flags, environment variables)
├── core/     # Shared business logic and image processing
│   └── weather/
│       ├── image.go     # NEW: Image rendering logic
│       ├── icons/       # NEW: Embedded weather icons
│       └── cache.go     # UPDATED: Image caching logic
```

**Structure Decision**: The implementation will reside primarily in `internal/core/weather` to allow both the API and CLI to share the same rendering and caching logic, adhering to Principle VII.
