# API Contract: Battery Refresh Optimization

## Base Path: `/`

### GET `/dashboard/image`
Returns a combined weather and agenda image with ETag support.

- **Request Headers**:
  - `If-None-Match`: (Optional) The last ETag received from the server.
- **Query Parameters**:
  - `location`: (Optional) City for weather.
  - `calendar_id`: (Optional) Google Calendar ID.
  - `palette`: (Optional) E-ink color palette.
- **Responses**:
  - `200 OK`: Success, new image returned.
    - **Headers**: `ETag: "<md5_hash>"`
    - **Body**: Binary PNG image data.
  - `304 Not Modified`: Image content has not changed since the last ETag.
    - **Body**: Empty.
  - `400 Bad Request`: Invalid parameters.
  - `500 Internal Server Error`: Generation failure.
