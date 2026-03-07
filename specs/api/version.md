# Version API Specification

### 1. Purpose
The HTTP endpoints for retrieving the application version string are defined by this specification.

### 2. Data Models

#### VersionResponse
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|
| version | string | Yes | None | The semantic version string of the application. |

### 3. Endpoints
| Endpoint Name | Verb | URL with query parameters | Request Model | Response Type | Response Model |
|---|---|---|---|---|---|
| Get Version | GET | `/api/version` | None | JSON | `VersionResponse` |

### 4. Dependencies
* [Core Version Module](../core/version.md)
