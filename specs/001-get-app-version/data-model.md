# Data Model: Get App Version

## Entities

### VersionInfo
Represents the application version metadata.

| Field | Type | Description |
|-------|------|-------------|
| version | string | Semantic version string (e.g., "1.0.0") |

## Persistence
- Static hardcoded value in `internal/config/version.go`.
