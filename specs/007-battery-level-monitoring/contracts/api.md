# API Contract: Battery Monitoring

### POST `/battery`
Report current battery voltage.

- **Request Body**: JSON
  ```json
  {
    "voltage": 3.75
  }
  ```
- **Responses**:
  - `201 Created`: Successfully recorded.
  - `400 Bad Request`: Invalid payload or negative voltage.
  - `500 Internal Server Error`: Disk write failure.

### GET `/battery/history`
Retrieve full battery report history as raw CSV text.

- **Responses**:
  - `200 OK`: Success.
  - **Content-Type**: `text/plain; charset=utf-8`
  - **Body**:
    ```csv
    Timestamp,Voltage
    2026-03-07T12:00:00Z,3.75
    2026-03-07T13:00:00Z,3.72
    ```
