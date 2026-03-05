# Research: Get App Version

## Decision: Version Storage
- **Decision**: Centralize version in `internal/config/version.go`.
- **Rationale**: User request for easy maintenance.
- **Alternatives considered**: Build tags (rejected), separate VERSION file (rejected).

## Decision: CLI Command Implementation
- **Decision**: Use `spf13/cobra` subcommands.
- **Rationale**: Matches project stack.
- **Implementation**: `versionCmd` added to the root command in `cmd/inky/`.

## Decision: API Handler
- **Decision**: Standard `http.HandlerFunc`.
- **Rationale**: Lightweight, no external router needed for this simple requirement.
- **Payload**: `{"version": "X.Y.Z"}`.

## Best Practices
- Ensure `http.ResponseWriter` sets `Content-Type: application/json`.
- CLI output should go to `Stdout`.
- Use standard access logging for the API endpoint as per `FR-012`.
