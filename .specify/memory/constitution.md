<!--
  Sync Impact Report:
  - Version change: N/A → 1.0.0
  - List of modified principles (old title → new title if renamed):
    - [PRINCIPLE_1_NAME] → Logic Offloading (Server-side rendering)
    - [PRINCIPLE_2_NAME] → Energy-First Lifecycle
    - [PRINCIPLE_3_NAME] → Data Integrity & Freshness
    - [PRINCIPLE_4_NAME] → Resource-Conscious Image Delivery
    - [PRINCIPLE_5_NAME] → API-First Development
  - Added sections:
    - Technical Constraints
    - Development Workflow
  - Removed sections: N/A
  - Templates requiring updates (✅ updated / ⚠ pending) with file paths:
    - ✅ .specify/memory/constitution.md (initial adoption)
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

## Technical Constraints

- **Client**: MicroPython on Raspberry Pi Pico W (Inky Frame).
- **Server**: Python 3.10+ (FastAPI recommended).
- **Display**: 7-color E-Ink (Pimoroni Inky Frame).
- **Protocol**: HTTP/HTTPS returning BMP or PNG optimized for the Inky palette.

## Development Workflow

- **Image validation**: Layout changes SHOULD be validated using local Python scripts and previewed as standard images before being integrated into the API.
- **Testing**: The Python API MUST include unit tests for data parsing and layout generation.
- **Contract verification**: Every change to the API that affects the image output MUST be manually verified with a sample image simulating the Inky display constraints.

## Governance

- This constitution supersedes all other development practices in this project.
- Amendments require a version bump following semantic versioning (MAJOR for breaking changes, MINOR for additions, PATCH for clarifications).
- All implementation plans must include a "Constitution Check" to verify alignment with these principles.

**Version**: 1.0.0 | **Ratified**: 2026-03-05 | **Last Amended**: 2026-03-05
