# Battery API Specification

### 1. Purpose
The HTTP endpoints for reporting battery metrics, fetching the latest battery status, and downloading historical battery metrics are defined by this specification.

### 2. Data Models

#### BatteryReportRequest
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| voltage | float | Yes | None | The measured battery voltage to be reported (Used in Request). |

#### BatteryReportResponse
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| voltage | float | Yes | None | The measured battery voltage. |
| percentage | float | Yes | None | The calculated battery percentage. |
| timestamp | datetime | Yes | None | The timestamp of the report. |

#### BatteryStatusResponse
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| latest | `BatteryReportResponse` | Yes | None | The most recent battery report (Used in Status Response). |
| is_low | boolean | Yes | None | Flag indicating if the battery percentage is less than `20%` (Used in Status Response). |
| last_reported | datetime | Yes | None | The timestamp of the latest report (Used in Status Response). |

#### BatteryHistoryRequest
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| limit | integer | No | 10 | The maximum number of historical reports to return. |
| offset | integer | No | 0 | The pagination offset for historical reports. |

### 3. Endpoints
| Endpoint Name | Verb | URL with query parameters | Request Model | Response Type | Response Model |
|---|---|---|---|---|---|
| Post Report | POST | `/api/battery` | `BatteryReportRequest` | JSON | `BatteryReportResponse` |
| Get Status | GET | `/api/battery/status` | None | JSON | `BatteryStatusResponse` |
| Get History | GET | `/api/battery/history?limit={limit}&offset={offset}` | `BatteryHistoryRequest` | JSON | Array of `BatteryReportResponse` |

### 4. Dependencies
* [Core Battery Service](../core/battery.md)
