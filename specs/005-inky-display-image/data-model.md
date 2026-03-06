# Data Model: Inky Frame Dashboard Display

This document defines the configuration variables and state for the MicroPython client.

## Entities

### Device Configuration (`firmware/.env.py`)
Variables stored on the device to control connectivity and behavior.

| Variable | Type | Description |
| :--- | :--- | :--- |
| `WIFI_SSID` | `string` | The SSID of the Wi-Fi network. |
| `WIFI_PASSWORD` | `string` | The password for the Wi-Fi network. |
| `DASHBOARD_URL` | `string` | The full URL to the Go API image endpoint. |
| `SLEEP_MINUTES` | `int` | Duration of deep sleep between updates (Default: 30). |
| `BATTERY_THRESHOLD` | `float` | Voltage threshold for low battery indicator (Default: 3.4). |

### Update State
Transient state during a wake cycle.

| Field | Type | Description |
| :--- | :--- | :--- |
| `Voltage` | `float` | Current VSYS voltage read from ADC. |
| `IsLowBattery` | `boolean` | True if `Voltage < BATTERY_THRESHOLD`. |
| `UpdateSuccess` | `boolean` | Outcome of image fetch and render. |

## Storage

- **Permanent**: `firmware/.env.py` (Flash memory on device).
- **Transient**: RAM buffer for PNG decoding.
