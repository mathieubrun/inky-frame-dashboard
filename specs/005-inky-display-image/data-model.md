# Data Model: Inky Frame Dashboard Display

This document defines the configuration and state model for the MicroPython client.

## Entities

### Device Configuration (`.env.py`)
Variables stored on the device to control connectivity and behavior.

| Variable | Type | Description |
| :--- | :--- | :--- |
| `WIFI_SSID` | `string` | The SSID of the Wi-Fi network. |
| `WIFI_PASSWORD` | `string` | The password for the Wi-Fi network. |
| `DASHBOARD_URL` | `string` | The full URL to the Go API image endpoint. |
| `SLEEP_MINUTES` | `int` | Duration of deep sleep between updates. |
| `DISPLAY_TYPE` | `string` | Physical hardware variant (e.g., `spectra`). |
| `BATTERY_THRESHOLD` | `float` | Voltage threshold for low battery indicator (e.g. 3.4). |

### Update State
Transient state during an update cycle.

| Field | Type | Description |
| :--- | :--- | :--- |
| `WokenBy` | `string` | Source of wake (RTC or Button). |
| `LastResult` | `string` | Outcome of last update (Success/Error). |
| `BatteryLevel` | `float` | Estimated battery voltage. |

## Storage

- **Permanent**: `.env.py` (Flash memory).
- **Temporary**: `dashboard.png` (Buffer or temporary file in Flash).
