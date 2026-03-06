# Implementation Plan: Inky Frame Dashboard Display

**Branch**: `005-inky-display-image` | **Date**: 2026-03-05 | **Spec**: [specs/005-inky-display-image/spec.md](spec.md)
**Input**: Feature specification from `/specs/005-inky-display-image/spec.md` + Low Battery Indicator requirement.

## Summary

This feature implements the MicroPython client for the Inky Frame 7.3". The device will wake up at a configurable interval (default 30 minutes), connect to Wi-Fi using credentials from `.env.py`, fetch a pre-rendered 800x480 PNG from the dashboard API, and update its e-ink display. The implementation includes battery voltage monitoring via ADC; if the voltage falls below a configurable threshold (e.g., 3.4V), a low battery indicator will be overlaid on the top right corner of the screen.

## Technical Context

**Language/Version**: MicroPython (Pimoroni Firmware for RP2040)
**Primary Dependencies**: inky_frame, picographics, pngdec, network, urequests, machine (ADC)
**Storage**: Internal Flash (for script and config)
**Testing**: Hardware verification (manual)
**Target Platform**: Raspberry Pi Pico W (Inky Frame 7.3")
**Project Type**: Embedded Client
**Performance Goals**: < 60s total active time per update cycle
**Constraints**: 264KB RAM, limited battery capacity (AA batteries)
**Scale/Scope**: Single-device client

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Validation |
| :--- | :--- | :--- |
| I. Logic Offloading | ✅ PASS | Device only handles fetching and rendering; Go API does all complex logic. |
| II. Energy-First | ✅ PASS | Uses `inky_frame.sleep_for()` to truly power off between updates. |
| III. Data Freshness | ✅ PASS | Image fetched on every wake cycle (30 min interval). |
| IV. Resource-Conscious | ✅ PASS | Uses hardware-accelerated `pngdec` and native `picographics`. |
| VII. Modular Architecture | ✅ PASS | Self-contained `main.py` client. |
| VIII. Flexible Configuration | ✅ PASS | Configured via `.env.py` (analogous to env vars for MicroPython). |

## Project Structure

### Documentation (this feature)

```text
specs/005-inky-display-image/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── quickstart.md        # Phase 1 output
```

### Source Code (Inky Frame Root)

```text
main.py      # Entry point and update logic
.env.py      # Local configuration (WIFI, API URL, Refresh Interval, Battery Threshold)
```

**Structure Decision**: The implementation will be a single MicroPython script `main.py` deployed to the root of the Inky Frame's flash memory, along with a `.env.py` for configuration. The script will perform a linear sequence: Wake -> Check Battery -> Connect WiFi -> Fetch Image -> Overlay Indicator (if needed) -> Render -> Sleep.
