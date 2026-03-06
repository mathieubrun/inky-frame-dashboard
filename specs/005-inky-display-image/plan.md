# Implementation Plan: Inky Frame Dashboard Display

**Branch**: `005-inky-display-image` | **Date**: 2026-03-05 | **Spec**: [specs/005-inky-display-image/spec.md](spec.md)
**Input**: Feature specification from `/specs/005-inky-display-image/spec.md` + "The inky micropython code resides in a firmware folder."

## Summary

This feature implements the MicroPython firmware for the Inky Frame 7.3". The client code, residing in the `firmware/` folder, will handle waking the device every 30 minutes, connecting to Wi-Fi, fetching a pre-rendered weather image from the Golang API, and rendering it to the e-ink display using hardware acceleration. It also includes battery monitoring to overlay a low-battery warning when the voltage drops below a specified threshold.

## Technical Context

**Language/Version**: MicroPython (Pimoroni Firmware for RP2040)
**Primary Dependencies**: inky_frame, picographics, pngdec, network, urequests, machine (ADC)
**Storage**: Internal Flash (for script and config)
**Testing**: Manual hardware verification
**Target Platform**: Raspberry Pi Pico W (Inky Frame 7.3")
**Project Type**: Embedded Client
**Performance Goals**: < 60s active time per cycle
**Constraints**: 264KB RAM, limited battery power
**Scale/Scope**: Single device firmware

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Validation |
| :--- | :--- | :--- |
| I. Logic Offloading | ✅ PASS | All rendering layout logic is handled by the server; device only displays the final PNG. |
| II. Energy-First | ✅ PASS | Uses `inky_frame.sleep_for()` to power off between cycles. |
| III. Data Freshness | ✅ PASS | Fetches fresh data on every wake cycle. |
| IV. Resource-Conscious | ✅ PASS | Uses `pngdec` for efficient image handling on limited RAM. |
| VII. Modular & Unified | ✅ PASS | Embedded code isolated in `firmware/` directory. |
| VIII. Flexible Configuration | ✅ PASS | Configured via `.env.py` template. |

## Project Structure

### Documentation (this feature)

```text
specs/005-inky-display-image/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── quickstart.md        # Phase 1 output
```

### Source Code (repository root)

```text
cmd/
└── inky/     # Golang application

firmware/     # NEW: MicroPython firmware for Inky Frame
├── main.py   # Entry point
└── .env.py   # Local configuration (WIFI, URL, etc.)

internal/     # Golang core logic
```

**Structure Decision**: MicroPython code is strictly separated into the `firmware/` folder to maintain clear separation of concerns between the server and client components.
