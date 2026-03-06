# Quickstart: Inky Frame Dashboard Display

This guide provides instructions for deploying the MicroPython client to your Inky Frame 7.3".

## Prerequisites

- Raspberry Pi Pico W (integrated into Inky Frame 7.3").
- [Pimoroni MicroPython Firmware](https://github.com/pimoroni/pimoroni-pico/releases) installed.
- [Thonny IDE](https://thonny.org/) or similar tool for file transfer.

## Installation Steps

1. **Configure credentials**:
   Create a file named `.env.py` on your computer with the following content:
   ```python
   WIFI_SSID = "Your_SSID"
   WIFI_PASSWORD = "Your_Password"
   DASHBOARD_URL = "http://your-server-ip:8080/api/v1/weather/image?location=Zurich"
   SLEEP_MINUTES = 30
   DISPLAY_TYPE = "spectra" # or "standard"
   ```

2. **Upload files**:
   Transfer the following files to the root directory of the Pico W:
   - `main.py` (The client script)
   - `.env.py` (Your configuration)

3. **Verify first run**:
   Press the **RESET** button on the back of the Inky Frame or power it via the JST connector.
   - The onboard LED should blink during the update.
   - The e-ink display should refresh with the weather dashboard.
   - The device will enter deep sleep after completion.

## Troubleshooting

- **No image?**: Verify the `DASHBOARD_URL` is accessible from your Wi-Fi network.
- **Immediate sleep?**: If powered via USB, the device might not truly power off, but it will wait for the interval. Test on battery for true deep sleep.
- **Memory Error**: Ensure the Go API is returning a 800x480 PNG and not a larger file.
