# Inky Frame Client Firmware

The design, configuration, and execution lifecycle of the MicroPython firmware running on the physical Inky Frame device are described by this document.

---

## 1. Purpose

The client firmware is executed on a Raspberry Pi Pico W-powered Pimoroni Inky Frame (7" e-ink display). Waking up periodically, reporting battery metrics, downloading the latest dashboard image from the server, rendering it to the screen, and returning to an ultra-low-power deep sleep state are handled as its main job.

---

## 2. Configuration (`env.py`)

A local configuration file `firmware/env.py` (based on `firmware/env.template.py`) is relied upon by the firmware.

| Configuration Key | Description | Example Value |
| :--- | :--- | :--- |
| `WIFI_SSID` | The Wi-Fi network name. | `"HomeNetwork"` |
| `WIFI_PASSWORD` | The Wi-Fi network password. | `"secretpassword"` |
| `DASHBOARD_URL` | The HTTP endpoint by which a standard PNG image is rendered and returned. | `"http://192.168.1.100:8080/dashboard/image"` |
| `BATTERY_URL` | The HTTP endpoint to which battery voltage is reported. | `"http://192.168.1.100:8080/battery"` |
| `SLEEP_MINUTES` | The duration of deep sleep between updates. | `60` |
| `BATTERY_THRESHOLD` | The voltage threshold below which a "LOW BATT" warning is overlaid. | `3.4` |

---

## 3. Execution Lifecycle

A state machine, by which the following linear steps are executed, is implemented by the main loop (`firmware/main.py`):

```
    ┌─────────────────────────┐
    │     Device Wakeup       │
    └───────────┬─────────────┘
                │
                ▼
    ┌─────────────────────────┐
    │  Measure Battery (VSYS) │
    └───────────┬─────────────┘
                │
                ▼
    ┌─────────────────────────┐
    │    Connect to Wi-Fi     │
    └───────────┬─────────────┘
                │
                ▼
    ┌─────────────────────────┐
    │  POST Battery Voltage   │
    └───────────┬─────────────┘
                │
                ▼
    ┌─────────────────────────┐
    │   Fetch PNG via GET     ├──────────┐ [HTTP 304 / Not Modified]
    └───────────┬─────────────┘          │
                │ [HTTP 200]             ▼
                ▼                ┌─────────────────────────┐
    ┌─────────────────────────┐  │      Skip Rendering     │
    │  Decode PNG (pngdec)    │  │ (Save Battery / Flash)  │
    └───────────┬─────────────┘  └───────┬─────────────────┘
                │                        │
                ▼                        │
    ┌─────────────────────────┐          │
    │   Update E-Ink Screen   │◄─────────┘
    └───────────┬─────────────┘
                │
                ▼
    ┌─────────────────────────┐
    │    Go to Deep Sleep     │
    └─────────────────────────┘
```

### Step 1: Battery Monitoring & VSYS
* The raw ADC value of the VSYS battery line (ADC pin 29) is read by the device so that LiPo battery voltage is measured.
* The factor `3 * 3.3 / 65535` is used for voltage conversion.
* Readings are averaged across 10 measurements for electrical stability.

### Step 2: Network Connection
* The Pico W station interface (`network.WLAN(network.STA_IF)`) is activated.
* Power management (`pm=0xa11140`) is temporarily disabled so that connection dropouts with slow or legacy access points are prevented.
* Up to 3 connection attempts are performed before the error state is fallen back to.

### Step 3: Battery Reporting
* A non-blocking HTTP `POST` payload (`{"voltage": float}`) is submitted to the configured `BATTERY_URL`.
* If reporting fails (e.g. server temporary timeout), normal execution is proceeded with.

### Step 4: Caching & Efficient Download
* The dashboard image is fetched from the `DASHBOARD_URL` with caching optimizations:
  * ETag Support: The cached ETag is read from `etag.txt` in internal flash and is sent in the `If-None-Match` header.
  * HTTP 304 Handling: If `304 Not Modified` is returned by the API, downloading and rendering are skipped by the firmware so that battery and flash memory wear are saved.
  * HTTP 200 Handling: The raw PNG stream is downloaded, `etag.txt` is updated with the new header ETag, and the image is saved to `dashboard.png`.

### Step 5: Screen Rendering
* The hardware-accelerated **`PicoGraphics`** display driver and the **`pngdec`** decoder library are used.
* Low Battery Failsafe Overlay: If the measured voltage is below the configured `BATTERY_THRESHOLD`, a high-visibility, solid red rectangle with white text `"LOW BATT"` is rendered in the top-right corner of the display buffer.
* The buffer is flushed to the screen using `display.update()`. (The physical e-ink pigments are triggered by this operation, which takes approximately 40 seconds).

---

## 4. Error Handling and Failsafes

MicroPython is run in a highly resource-constrained environment (especially concerning heap memory and network timeouts).

### Failsafe Safeguards
* Memory Management: The garbage collector is invoked explicitly via `gc.collect()` before memory-heavy routines (fetching the image and decoding the PNG buffer) so that `MemoryError` allocation failures are prevented.
* Network Timeouts: All connection fetches are strictly bound by timeouts (10 seconds for battery reports, 30 seconds for image downloads) so that a permanent battery drain by a hung request is prevented.
* On-Screen Failure Diagnostics: If a critical error occurs (Wi-Fi fails, download fails, PNG is corrupted):
  1. The stored `etag.txt` is deleted so that a clean reload is forced during the next wake cycle.
  2. The display buffer is cleared to solid White.
  3. A clean error view is displayed: `"Update Failed"`, the specific error message, and the exact date/time from the internal Real-Time Clock (RTC) are shown.
  4. The screen is updated and the device is returned to deep sleep so that battery depletion is avoided.
