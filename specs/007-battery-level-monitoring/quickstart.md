# Quickstart: Battery Level Monitoring

## Server-Side (API)

### Start the Server
Run the API server with default configuration:
```bash
inky serve --port 8080
```

### Report Battery (Manual Test)
Simulate a device reporting battery level using `curl`:
```bash
curl -X POST http://localhost:8080/api/v1/battery \
  -H "Content-Type: application/json" \
  -d '{"voltage": 3.75}'
```

### View History (Manual Test)
Retrieve the recorded history:
```bash
curl http://localhost:8080/api/v1/battery/history
```

## Inky Frame (MicroPython)

### Integration Example
Add this function to your `main.py`:
```python
import urequests
import machine

def report_battery(server_url):
    adc = machine.ADC(29) # Vsys monitoring pin on Pico W
    conversion_factor = 3 * 3.3 / 65535
    voltage = adc.read_u16() * conversion_factor
    try:
        r = urequests.post(f"{server_url}/api/v1/battery", json={"voltage": voltage})
        r.close()
    except Exception as e:
        print(f"Failed to report battery: {e}")

# In your main execution loop:
# report_battery("http://your-server-ip:8080")
# get_dashboard_image()
```

## CLI Interface
View the battery history directly from the server CLI:
```bash
inky battery history
```
Clear all history (be careful!):
```bash
inky battery clear
```
