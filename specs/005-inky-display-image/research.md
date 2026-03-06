# Research: Inky Frame Dashboard Display

This document outlines the technical research findings for implementing the MicroPython client on the Inky Frame 7.3".

## MicroPython Environment

**Decision**: Use the standard Pimoroni MicroPython firmware for Raspberry Pi Pico W.

**Rationale**:
- Provides optimized libraries for Inky Frame hardware (`inky_frame`, `picographics`).
- Includes `pngdec` for hardware-accelerated PNG decoding.
- Reliable Wi-Fi support via the `network` module.

## Image Display

**Decision**: Use `picographics` and `pngdec` to render the 800x480 dashboard image.

**Implementation Details**:
- Target display: `DISPLAY_INKY_FRAME_7` (or `DISPLAY_INKY_FRAME_7_SPECTRA` for Spectra 6 models).
- Process: Download PNG -> Save to temporary file -> `pngdec.decode(0, 0)` -> `display.update()`.
- Resolution: Fixed at 800x480.

## Power Management (Deep Sleep)

**Decision**: Use `inky_frame.sleep_for(minutes)` for scheduled updates.

**Implementation Details**:
- This function sets the onboard PCF85063A RTC alarm and cuts power to the Pico W.
- The device performs a hardware reset upon waking, meaning the entry point must be `main.py`.
- Interval: Default 30 minutes (configurable via `.env.py`).

## Configuration via `.env.py`

**Decision**: Store credentials and URL in a `.env.py` file.

**Rationale**:
- Native MicroPython support (it's just a Python module).
- Easy to import: `import env`.
- Structure:
  ```python
  WIFI_SSID = "Your SSID"
  WIFI_PASSWORD = "Your Password"
  DASHBOARD_URL = "http://your-server:8080/api/v1/weather/image?location=Zurich"
  UPDATE_INTERVAL_MINS = 30
  ```

## Network Requests

**Decision**: Use `urequests` for HTTP GET calls.

**Optimization**:
- Set a strict timeout (e.g., 15 seconds) to prevent battery drain on poor connections.
- Perform `gc.collect()` before and after large downloads to manage limited RAM.

## Visual Feedback

**Decision**: Use the Inky Frame's onboard LEDs or a small text overlay during the "Active" phase.

**Rationale**:
- E-ink updates are slow (~40s). Users need to know the device is working.
- Onboard LEDs provide immediate feedback without a full e-ink refresh.

## Battery Monitoring

**Decision**: Read VSYS voltage via RP2040 ADC.

**Implementation Details**:
- **ADC Pin**: 29 (VSYS).
- **Control Pins**: Pin 2 (VSYS Hold) and Pin 25 (VSYS Read Enable) must be HIGH.
- **Conversion**: VSYS voltage is divided by 3. Formula: `Reading * 3 * 3.3 / 65535`.
- **Threshold**: Default low battery threshold is 3.4V (for 3xAA alkaline).
- **Overlay**: Use `picographics` to draw a small battery icon or "LOW BATT" text in the top right corner if below threshold.
