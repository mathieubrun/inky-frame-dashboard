# Battery API Test Specification

## 1. POST /api/battery
### Given
A mocked `BatteryService` is provided. A parameterized payload is prepared, including:
  * A valid `BatteryReportRequest` payload.
  * An invalid payload missing required fields (e.g., `{}`).
  * An invalid payload with wrong data types (e.g., `{"voltage": "high"}`).
### When
An HTTP POST request is executed to `/api/battery`.
### Then
* For valid payloads, a status code of `201 Created` and the `BatteryReportResponse` model are returned.
* For invalid payloads, a status code of `422 Unprocessable Entity` is returned.
* On internal execution failure, a status code of `400 Bad Request` is returned.

## 2. GET /api/battery/status
### Given
A mocked `BatteryService` is provided. Both a valid request and a malformed request (e.g., unsupported HTTP method) are prepared.
### When
An HTTP GET request is executed to `/api/battery/status`.
### Then
On success, `200 OK` and the `BatteryStatusResponse` model are returned. If no history is available, `404 Not Found` is returned. On a malformed request, `405 Method Not Allowed` is returned. On other failures, `500 Internal Server Error` is returned.

## 3. GET /api/battery/history (Parameterized)
### Given
A mocked `BatteryService` containing history data is provided, alongside parameterized query strings for pagination (e.g., `?limit=10&offset=0`, `?limit=0`, or `?limit=150`).
### When
An HTTP GET request is executed to `/api/battery/history` with the parameterized query strings.
### Then
* For valid pagination bounds, a status code of `200 OK` and a JSON array of `BatteryReportResponse` models are returned by the endpoint.
* For invalid pagination parameters (e.g., limit below 1 or exceeding 100), a status code of `422 Unprocessable Entity` is returned.
