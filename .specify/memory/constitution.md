<!--
  Sync Impact Report:
  - Version change: 1.0.0 → 1.1.0
  - List of modified principles (old title → new title if renamed):
    - Added VI. Tooling Consistency
  - Added sections:
    - Technical Stack (merged/expanded from Constraints)
  - Removed sections: N/A
  - Templates requiring updates (✅ updated / ⚠ pending) with file paths:
    - ✅ .specify/memory/constitution.md
  - Follow-up TODOs if any placeholders intentionally deferred: N/A
-->

# Inky Frame Dashboard Constitution

## Core Principles

### I. Logic Offloading (Server-side rendering)
The Inky Frame MUST be treated as a "dumb" display. All complex logic, including data fetching from weather/calendar APIs and image layout generation, MUST be performed by the Python API. The Inky Frame SHOULD only make a simple HTTP request to receive a ready-to-display image.
*Rationale*: This minimizes the power-hungry processing and Wi-Fi on-time for the battery-powered Inky Frame, while simplifying the MicroPython code on the device.

### II. Energy-First Lifecycle
Development MUST prioritize battery longevity. The Inky Frame MUST enter deep sleep between scheduled updates. Network requests SHOULD be consolidated and timeouts MUST be strictly enforced to prevent excessive battery drain during connectivity issues.
*Rationale*: Inky Frames are typically battery-operated; inefficient code leads to frequent charging and a poor user experience.

### III. Data Integrity & Freshness
The Python API MUST ensure that the returned image contains accurate and up-to-date information. If an upstream data source (e.g., weather API) is unavailable, the image SHOULD clearly indicate the stale state or the time of last successful update to avoid misleading the user.
*Rationale*: A dashboard is only useful if its information is trustworthy and its current state is transparent.

### IV. Resource-Conscious Image Delivery
Images delivered to the Inky Frame MUST be optimized for its specific display capabilities (e.g., 7-color palette, fixed dimensions). The Python API SHOULD handle all dithering and color mapping to ensure the best possible visual quality with minimal client-side decoding.
*Rationale*: Offloading image processing ensures faster refresh times and higher quality visuals on the E-Ink display without overtaxing the MicroPython environment.

### V. API-First Development
All new dashboard features MUST start with an update to the Python API and its image generation logic. The communication contract between the Inky Frame and the API MUST be stable and ideally versioned to prevent breaking the client during server-side updates.
*Rationale*: Decoupling the data presentation from the display hardware allows for rapid iteration and testing without requiring frequent firmware updates to the Inky Frame.

### VI. Tooling Consistency
Development MUST adhere to the project's selected modern Python toolchain. All dependencies MUST be managed via `uv`, and all code MUST be linted and formatted using `ruff`.
*Rationale*: Using a consistent, modern, and fast toolchain ensures developer productivity and maintains high code quality across the project.

## Technical Stack

- **Python Version**: 3.13 (Mandatory for server-side logic).
- **Package Management**: `uv` (Required for performance and reproducibility).
- **Linting & Formatting**: `ruff` (Strict compliance required).
- **CLI Framework**: `typer`.
- **API Framework**: `fastapi`.
- **Testing Framework**: `pytest` with `pytest-cov` (Targeting >80% coverage).
- **Client**: MicroPython on Raspberry Pi Pico W (Inky Frame).
- **Display**: 7-color E-Ink (Pimoroni Inky Frame).

## Development Workflow

- **Dependency Management**: Use `uv add` for new dependencies and `uv lock` to maintain deterministic environments.
- **Linting**: Run `ruff check` and `ruff format` before every commit.
- **Testing**: The Python API MUST include unit tests for data parsing and layout generation. Run `pytest --cov` to verify coverage.
- **Image validation**: Layout changes SHOULD be validated using local Python scripts and previewed as standard images before being integrated into the API.
- **Contract verification**: Every change to the API that affects the image output MUST be manually verified with a sample image simulating the Inky display constraints.

## Governance

- This constitution supersedes all other development practices in this project.
- Amendments require a version bump following semantic versioning (MAJOR for breaking changes, MINOR for additions, PATCH for clarifications).
- All implementation plans must include a "Constitution Check" to verify alignment with these principles.

**Version**: 1.1.0 | **Ratified**: 2026-03-05 | **Last Amended**: 2026-03-05
