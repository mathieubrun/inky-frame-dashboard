# Quickstart: Inky Frame Dashboard Display

This guide provides instructions for deploying the MicroPython client from the `firmware/` folder to your Inky Frame 7.3".

## Prerequisites

- Raspberry Pi Pico W (integrated into Inky Frame 7.3").
- [Pimoroni MicroPython Firmware](https://github.com/pimoroni/pimoroni-pico/releases) installed.
- [Thonny IDE](https://thonny.org/) for file transfer.

## Installation Steps

1. **Configure credentials**:
   In your local repository, navigate to `firmware/` and create/edit `.env.py`:
   ```python
   WIFI_SSID = "Your_SSID"
   WIFI_PASSWORD = "Your_Password"
   DASHBOARD_URL = "http://your-server-ip:8080/api/v1/weather/image?location=Zurich"
   SLEEP_MINUTES = 30
   BATTERY_THRESHOLD = 3.4
   ```

2. **Upload files**:
   Using Thonny, upload the following files from the `firmware/` folder to the root of the Pico W:
   - `main.py`
   - `.env.py`

3. **Verify first run**:
   Press the **RESET** button on the Inky Frame.
   - The device should wake up and connect to Wi-Fi.
   - The e-ink display should refresh with the dashboard image.
   - If the battery is low, a red indicator will appear in the top right.
   - The device will enter deep sleep for the configured interval.

## Troubleshooting

- **Connection Refused**: Ensure the dashboard API is running and reachable from the device's network.
- **Low Battery**: If the indicator persists despite fresh batteries, check the `BATTERY_THRESHOLD` in `.env.py`.
