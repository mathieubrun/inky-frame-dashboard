# Feature Specification: Get App Version

**Feature Branch**: `001-get-app-version`  
**Created**: 2026-03-05  
**Status**: Draft  
**Input**: User description: "I want to be able to get the inky-dashboard app version"

## Clarifications

### Session 2026-03-05
- Q: Should the JSON response be a simple object with a "version" key, or a raw JSON string? → A: JSON Object: `{"version": "1.0.0"}`
- Q: Should the version string include a lowercase 'v' prefix (e.g., `v1.0.0`) or be the raw numbers only (e.g., `1.0.0`)? → A: Raw numbers only (e.g., `1.0.0`)
- Q: Should the `/version` API endpoint be publicly accessible, or should it require authentication? → A: Publicly accessible (no auth required)
- Q: Should successful requests to the `/version` API endpoint be recorded in the application logs? → A: Log every request (standard behavior)
- Q: Which exit codes should the CLI `version` command return to indicate its execution status? → A: 0 on success, 1 on failure

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Retrieve Version via CLI (Priority: P1)

As a developer or system administrator, I want to retrieve the application version using a CLI command so that I can quickly verify which version is installed on a local machine.

**Why this priority**: Essential for local development, troubleshooting, and deployment verification.

**Independent Test**: Running the version command in the terminal should display the correct semantic version.

**Acceptance Scenarios**:

1. **Given** the application is installed, **When** I run `inky --version` or `inky version`, **Then** the current semantic version (raw numbers only) is printed to the standard output.
2. **Given** the application is installed, **When** I run the version command, **Then** the output format is concise and easily readable.
3. **Given** the command executes, **When** it completes successfully, **Then** it returns an exit code of 0.

---

### User Story 2 - Retrieve Version via API (Priority: P1)

As a remote user or monitoring system, I want to retrieve the application version via an API endpoint so that I can remotely verify the version of a running service.

**Why this priority**: Critical for remote management and ensuring compatibility between different system components.

**Independent Test**: Calling the version endpoint using an HTTP client should return a JSON response with the version.

**Acceptance Scenarios**:

1. **Given** the API service is running, **When** I make a GET request to `/version`, **Then** I receive a `200 OK` response with a JSON body containing the version string (raw numbers only).
2. **Given** the API service is running, **When** I make a GET request to `/version`, **Then** the response includes the version string within a "version" field.

---

### Edge Cases

- **Service Unavailable**: If the API is down, the CLI command should still work for the local binary.
- **Malformed Version**: If the version cannot be determined (e.g., development build), the system should return a sensible fallback like `0.0.0-dev`.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a CLI subcommand `version` to display the application version.
- **FR-002**: System MUST support a global flag `--version` to display the application version.
- **FR-003**: System MUST provide an HTTP GET endpoint `/version` to retrieve the application version.
- **FR-004**: The version string MUST follow semantic versioning (e.g., X.Y.Z) and MUST be raw numbers only (no 'v' prefix).
- **FR-005**: The CLI output MUST be ONLY plain text.
- **FR-006**: The API response MUST be a JSON object with a `version` key.
- **FR-007**: The version information MUST be consistent between the CLI and the API.
- **FR-008**: The version response MUST include ONLY the semantic version string as the value.
- **FR-009**: The API MUST listen on a port configurable via flags or environment variables (default: 8080).
- **FR-010**: The application version MUST be hardcoded in a dedicated configuration file (e.g., `internal/config/version.go`) for easy maintenance.
- **FR-011**: The `/version` API endpoint MUST be publicly accessible without authentication.
- **FR-012**: The application MUST log all requests to the `/version` endpoint following standard access logging practices.
- **FR-013**: The CLI version command MUST return exit code 0 on success and 1 on failure.

### Key Entities *(include if feature involves data)*

- **Version Info**: A structure containing the major, minor, and patch numbers (semantic version).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can retrieve the version via CLI in under 100ms.
- **SC-002**: API `/version` endpoint responds in under 50ms.
- **SC-003**: Version information is 100% consistent across all access methods (CLI and API).
- **SC-004**: 100% of automated health checks can successfully parse the version from the JSON response.
