# Dashboard API Specification

### 1. Purpose
The HTTP endpoints for retrieving dynamically generated graphical representations (weather, calendar, and combined views) are exposed by the dashboard API module.

### 2. Data Models

#### DashboardImageRequest
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| width | integer | No | 600 | The width of the requested image in pixels. |
| height | integer | No | 448 | The height of the requested image in pixels. |

### 3. Endpoints
| Endpoint Name | Verb | URL with query parameters | Request Model | Response Type | Response Model |
|---|---|---|---|---|---|
| Weather Image | GET | `/api/dashboard/weather?width={width}&height={height}` | `DashboardImageRequest` | image/png | None |
| Calendar Image | GET | `/api/dashboard/calendar?width={width}&height={height}` | `DashboardImageRequest` | image/png | None |
| Combined Image | GET | `/api/dashboard/combined?width={width}&height={height}` | `DashboardImageRequest` | image/png | None |

Note: Standard HTTP status codes must be returned when the underlying service raises domain errors (e.g., `404 Not Found` or `400 Bad Request`).

### 4. Dependencies
* [Core Image Module](../core/image.md)
