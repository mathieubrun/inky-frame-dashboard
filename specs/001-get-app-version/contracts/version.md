# API Contract: Get App Version

## GET /version

Retrieve the application's current version.

**Response**: `200 OK`
**Content-Type**: `application/json`

### Body
```json
{
  "version": "1.0.0"
}
```

## CLI Interface

### Subcommand: `version`
**Command**: `inky version`
**Output**: `1.0.0` (Plain text)
**Exit Codes**:
- `0`: Success
- `1`: Failure

### Global Flag: `--version`
**Command**: `inky --version`
**Output**: `1.0.0` (Plain text)
**Exit Codes**:
- `0`: Success
- `1`: Failure
