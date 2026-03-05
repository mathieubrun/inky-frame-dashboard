# Research: Get App Version

## Decision: Version Storage & Injection
- **Decision**: Hardcode the version string in `internal/config/version.go`.
- **Rationale**: User preference for simplicity and ease of maintenance.
- **Alternatives considered**: Injected via `ldflags` at build time (rejected per user choice).

## Decision: CLI Framework Integration
- **Decision**: Use `spf13/cobra` to implement both the `version` subcommand and the `--version` flag on the root command.
- **Rationale**: Industry standard for Go CLIs, already part of the project's technical stack.

## Decision: API Implementation
- **Decision**: Use standard library `net/http` for the `/version` endpoint.
- **Rationale**: Minimal overhead for a simple metadata endpoint, aligns with project principles.

## Best Practices: Versioning
- Follow strict Semantic Versioning 2.0.0.
- Ensure the version string is a single source of truth used by both CLI and API handlers.
- Standard access logging for the API endpoint as requested.
