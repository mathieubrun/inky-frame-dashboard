# Quickstart: Battery Refresh Optimization

## Server-Side (Go)

### Test ETag support manually
Run the API server:
```bash
inky serve --port 8080
```

Request the image for the first time:
```bash
curl -v http://localhost:8080/dashboard/image?location=Zurich
```
Copy the `ETag` value from the response headers (e.g., `"abcdef123"`).

Request again with the ETag:
```bash
curl -v http://localhost:8080/dashboard/image?location=Zurich \
  -H 'If-None-Match: "abcdef123"'
```
Verify the server returns `HTTP/1.1 304 Not Modified` with an empty body.

## Client-Side (MicroPython)

### Integration
1. The new `main.py` will automatically manage `/etag.txt`.
2. To force a full refresh, manually delete `/etag.txt` from the device:
   ```python
   import os
   os.remove("etag.txt")
   ```

### Debugging
- Observe the device logs (via Thonny or serial console).
- Look for messages like "Status: 304, skipping screen refresh" or "Status: 200, updating display".
