# Research: Inky Frame Dashboard Display

This document outlines the technical research findings for implementing the MicroPython client on the Inky Frame 7.3".

## Project Organization

**Decision**: All MicroPython source code will reside in a `firmware/` folder in the repository.

**Rationale**:
- Separates embedded client code from the Golang API/CLI source.
- Simplifies deployment scripts and documentation.
- Follows the user constraint: "The inky micropython code resides in a firmware folder."

## MicroPython Environment

**Decision**: Use the standard Pimoroni MicroPython firmware for Raspberry Pi Pico W.

**Rationale**:
- Provides optimized libraries for Inky Frame hardware (`inky_frame`, `picographics`).
- Includes `pngdec` for hardware-accelerated PNG decoding.
- Reliable Wi-Fi support via the `network` module.

## Image Display & Rendering

**Decision**: Use `picographics` and `pngdec` to render the 800x480 dashboard image.

**Implementation Details**:
- Target display: `DISPLAY_INKY_FRAME_7`.
- Process: Download PNG -> `pngdec.decode(0, 0)` -> `display.update()`.
- The battery overlay will be drawn using `picographics` primitives (rectangles/text) on top of the decoded PNG buffer before calling `update()`.

## Power Management (Deep Sleep)

**Decision**: Use `inky_frame.sleep_for(minutes)` for scheduled updates.

**Implementation Details**:
- This function sets the onboard PCF85063A RTC alarm and cuts power to the Pico W.
- The device performs a hardware reset upon waking, meaning the entry point must be `main.py` (deployed from `firmware/main.py`).
- Interval: Default 30 minutes (configurable via `.env.py`).

## Configuration via `.env.py`

**Decision**: Store credentials and URL in a `.env.py` file (deployed from `firmware/.env.py`).

**Structure**:
```python
WIFI_SSID = "Your SSID"
WIFI_PASSWORD = "Your Password"
DASHBOARD_URL = "http://your-server:8080/api/v1/weather/image?location=Zurich"
UPDATE_INTERVAL_MINS = 30
BATTERY_THRESHOLD = 3.4
```

## Battery Monitoring

**Decision**: Read VSYS voltage via RP2040 ADC (Pin 29).

**Implementation Details**:
- **Control Pins**: Pin 2 (VSYS Hold) and Pin 25 (VSYS Read Enable) must be HIGH.
- **Conversion**: Formula `Reading * 3 * 3.3 / 65535`.
- **Threshold**: 3.4V for alkaline batteries.
