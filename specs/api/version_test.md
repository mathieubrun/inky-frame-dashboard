# Version API Test Specification

## 1. GET /api/version
### Given
The HTTP API is run with the version router, and the Core Version Module is mocked to return a valid version string. Valid requests and malformed requests (e.g., POST instead of GET) are prepared.
### When
An HTTP GET request is executed to `/api/version`.
### Then
For a valid request, a status code of `200 OK` is returned, matching the `VersionResponse` model. For a malformed request, a `405 Method Not Allowed` or `422 Unprocessable Entity` is returned.
