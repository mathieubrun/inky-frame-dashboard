# dogmas

- gemini is run in yolo mode
- of course I commit on master until asked by the model !
- 1 prompt = 1 commit with prompt as a message

constitution, specify, plan, implement : no manual modifications, only through spec kit

# Inky Frame Dashboard

A centralized dashboard for managing the Inky Frame (Raspberry Pi Pico W) and other devices.

## Features
- Retrieve application version via CLI and API.
- Battery level monitoring for Inky Frame devices.
- Persistent battery history in CSV format.

## CLI Usage
```bash
./inky version
./inky battery history # View full history
./inky battery clear   # Reset history
```

## API Usage
```bash
./inky serve --port 8080
curl http://localhost:8080/version
curl http://localhost:8080/battery/status  # Get latest status (JSON)
curl http://localhost:8080/battery/history # Get full history (CSV)
```
