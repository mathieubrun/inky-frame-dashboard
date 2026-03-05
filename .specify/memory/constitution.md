<!--
  Sync Impact Report:
  - Version change: 2.0.0 → 2.1.0
  - List of modified principles (old title → new title if renamed):
    - VII. Modular & Unified Architecture: Added single-binary requirement.
    - Added VIII. Flexible Configuration: Flags and Environment variables.
  - Added sections: N/A
  - Removed sections: N/A
  - Templates requiring updates (✅ updated / ⚠ pending) with file paths:
    - ✅ .specify/memory/constitution.md
    - ✅ .specify/templates/plan-template.md
    - ✅ .specify/templates/tasks-template.md
  - Follow-up TODOs if any placeholders intentionally deferred: N/A
-->

# Inky Frame Dashboard Constitution

## Core Principles

### I. Logic Offloading (Server-side rendering)
The Inky Frame MUST be treated as a "dumb" display. All complex logic, including data fetching from weather/calendar APIs and image layout generation, MUST be performed by the Golang API. The Inky Frame SHOULD only make a simple HTTP request to receive a ready-to-display image.
*Rationale*: This minimizes the power-hungry processing and Wi-Fi on-time for the battery-powered Inky Frame, while simplifying the MicroPython code on the device.

### II. Energy-First Lifecycle
Development MUST prioritize battery longevity. The Inky Frame MUST enter deep sleep between scheduled updates. Network requests SHOULD be consolidated and timeouts MUST be strictly enforced to prevent excessive battery drain during connectivity issues.
*Rationale*: Inky Frames are typically battery-operated; inefficient code leads to frequent charging and a poor user experience.

### III. Data Integrity & Freshness
The Golang API MUST ensure that the returned image contains accurate and up-to-date information. If an upstream data source (e.g., weather API) is unavailable, the image SHOULD clearly indicate the stale state or the time of last successful update to avoid misleading the user.
*Rationale*: A dashboard is only useful if its information is trustworthy and its current state is transparent.

### IV. Resource-Conscious Image Delivery
Images delivered to the Inky Frame MUST be optimized for its specific display capabilities (e.g., 7-color palette, fixed dimensions). The Golang server SHOULD handle all dithering and color mapping to ensure the best possible visual quality with minimal client-side decoding.
*Rationale*: Offloading image processing ensures faster refresh times and higher quality visuals on the E-Ink display without overtaxing the MicroPython environment.

### V. API-First Development
All new dashboard features MUST start with an update to the Golang API and its image generation logic. The communication contract between the Inky Frame and the API MUST be stable and ideally versioned to prevent breaking the client during server-side updates.
*Rationale*: Decoupling the data presentation from the display hardware allows for rapid iteration and testing without requiring frequent firmware updates to the Inky Frame.

### VI. Tooling Consistency
Development MUST adhere to standard Go best practices. All dependencies MUST be managed via `go modules`, and all code MUST be linted and formatted using `gofmt` and `golangci-lint`.
*Rationale*: Using a consistent, standard toolchain ensures developer productivity and maintains high code quality across the project.

### VII. Modular & Unified Architecture
The source code MUST be organized into a clear hierarchy that separates presentation (API/CLI) from core logic. Every piece of functionality exposed via the API endpoint MUST also be accessible through the CLI. The release process MUST produce a single binary that contains both the HTTP server (e.g., via a `serve` subcommand) and all CLI functionality.
*Rationale*: A single binary simplifies deployment and ensures that the system is easily testable, maintainable, and verifiable in any environment.

### VIII. Flexible Configuration
The application MUST be configurable via both command-line flags and environment variables. Flags MUST take precedence over environment variables, which MUST take precedence over default values.
*Rationale*: This ensures the application can be easily configured across different deployment environments (e.g., local development, Docker, bare metal) without requiring code changes.

## Technical Stack

- **Language**: Go (Golang) 1.22+
- **Package Management**: `go mod`
- **Linting & Formatting**: `gofmt` and `golangci-lint` (Strict compliance required)
- **CLI Framework**: `cobra` (spf13/cobra)
- **Configuration**: `viper` (or equivalent supporting flags and env vars)
- **API Framework**: standard library `net/http`
- **Testing Framework**: standard library `testing` (Targeting >80% coverage)
- **Client**: MicroPython on Raspberry Pi Pico W (Inky Frame)
- **Display**: 7-color E-Ink (Pimoroni Inky Frame)

## Project Layout

Following standard Go idioms, the codebase MUST be structured as follows:
- `cmd/inky/`: The single application entry point.
- `internal/api/`: HTTP handlers, routes, and API-specific logic.
- `internal/cli/`: Cobra command definitions (including the root and subcommands).
- `internal/core/`: Common business logic, data models, and image processing shared by both API and CLI.
- `internal/config/`: Configuration loading logic (flags/env).
- **Tests**: `*_test.go` files MUST be placed alongside the code they test in the same directory.

## Development Workflow

- **Dependency Management**: Use `go get` for new dependencies and `go mod tidy` to clean up the module file.
- **Linting**: Run `golangci-lint run` and ensure code is formatted with `gofmt` before every commit.
- **Testing**: The Golang API MUST include unit tests for data parsing and layout generation. Run `go test ./... -cover` to verify coverage.
- **Image validation**: Layout changes SHOULD be validated using local tests and previewed as standard images before being integrated into the API.
- **Contract verification**: Every change to the API that affects the image output MUST be manually verified with a sample image simulating the Inky display constraints.

## Governance

- This constitution supersedes all other development practices in this project.
- Amendments require a version bump following semantic versioning (MAJOR for breaking changes, MINOR for additions, PATCH for clarifications).
- All implementation plans must include a "Constitution Check" to verify alignment with these principles.

**Version**: 2.1.0 | **Ratified**: 2026-03-05 | **Last Amended**: 2026-03-05
