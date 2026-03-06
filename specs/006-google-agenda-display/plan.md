# Implementation Plan: Google Agenda Display

**Branch**: `006-google-agenda-display` | **Date**: 2026-03-06 | **Spec**: [specs/006-google-agenda-display/spec.md](spec.md)
**Input**: Feature specification from `/specs/006-google-agenda-display/spec.md` + "like weather data implementation, a mock api must be provided for testing purposes."

## Summary

This feature integrates Google Calendar events into the Inky Frame Dashboard. It provides a server-side synchronization engine using Google Service Accounts, exposing the data via a new `/api/v1/agenda` endpoint and a corresponding `inky agenda list` CLI command. Most importantly, it introduces a combined dashboard image generator that merges weather data (left panel) and upcoming calendar events (right panel) into a single 800x480 PNG optimized for e-ink displays.

## Technical Context

**Language/Version**: Go (Golang) 1.22+
**Primary Dependencies**: 
- `google.golang.org/api/calendar/v3`
- `google.golang.org/api/option`
- `github.com/fogleman/gg` (Existing)
- `spf13/cobra`, `spf13/viper` (Existing)
**Storage**: File-based JSON cache for calendar events.
**Testing**: standard library testing + `httptest` for mocking Google APIs.
**Target Platform**: Linux Server (API)
**Project Type**: web-service/cli
**Performance Goals**: Combined image generation in < 3s.
**Constraints**: 800x480 resolution, 6-color palette (Spectra 6).
**Scale/Scope**: Support for multiple shared calendars via service account.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Validation |
| :--- | :--- | :--- |
| I. Logic Offloading | ✅ PASS | All Google API communication and image composition is handled on the Go server. |
| II. Energy-First | ✅ PASS | The combined image reduces the number of HTTP requests and processing time for the Inky Frame. |
| III. Data Integrity | ✅ PASS | Last fetch time is tracked; "Calendar Unavailable" states are handled gracefully in the UI. |
| IV. Resource-Conscious | ✅ PASS | Single 800x480 PNG optimized for the target display. |
| V. API-First | ✅ PASS | Defined new endpoints for both raw data and processed images. |
| VII. Modular Architecture | ✅ PASS | Shared rendering logic between API and CLI. |
| IX. API Testing | ✅ PASS | New `.bru` files planned for agenda and dashboard endpoints. |
| X. Mandatory Testing | ✅ PASS | Mock provider included for local development and CI verification. |

## Project Structure

### Documentation (this feature)

```text
specs/006-google-agenda-display/
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
internal/
├── core/
│   ├── agenda/          # NEW: Google Calendar logic
│   │   ├── google.go    # API communication
│   │   ├── mock.go      # Mock provider for testing
│   │   ├── cache.go     # Event caching
│   │   └── types.go     # Data structures
│   └── dashboard/       # NEW: Combined image rendering
│       └── render.go    # Split-screen logic
├── api/
│   ├── agenda.go        # NEW: Agenda handlers
│   └── dashboard.go     # NEW: Combined image handler
└── cli/
    ├── agenda.go        # NEW: `inky agenda` commands
    └── dashboard.go     # NEW: `inky dashboard` commands
```

**Structure Decision**: I am creating a new `core/agenda` package to keep calendar logic separate from weather. A `core/dashboard` package will handle the high-level task of composing the final image from various data sources.
