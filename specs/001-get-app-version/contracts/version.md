# Contracts: Get App Version

## API Contract

### GET /version
- **Request**: `GET /version`
- **Headers**: `Accept: application/json`
- **Response**: `200 OK`
- **Body**:
  ```json
  {
    "version": "1.0.0"
  }
  ```

## CLI Contract

### `inky version`
- **Command**: `inky version`
- **Output**: `1.0.0` (followed by newline)
- **Exit Code**: `0` on success.
