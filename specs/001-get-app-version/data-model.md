# Data Model: Get App Version

## Entities

### VersionInfo
Represents the application's version information.

| Field | Type | Description |
|-------|------|-------------|
| version | string | Semantic version string (X.Y.Z) |

## Constraints
- Version MUST follow SemVer 2.0.0.
- Version MUST be raw numbers only (no 'v' prefix).
