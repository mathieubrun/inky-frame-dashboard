# Feature Specification: Get App Version

**Feature Branch**: `001-get-app-version`  
**Created**: 2026-03-05  
**Status**: Draft  
**Input**: User description: "I want to be able to get the inky-dashboard app version"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Retrieve Version via CLI (Priority: P1)

As a developer or system administrator, I want to retrieve the application version using a CLI command so that I can quickly verify which version is installed on a local machine.

**Why this priority**: Essential for local development, troubleshooting, and deployment verification.

**Independent Test**: Running the version command in the terminal should display the correct semantic version.

**Acceptance Scenarios**:

1. **Given** the application is installed, **When** I run `inky --version` or `inky version`, **Then** the current semantic version is printed to the standard output.
2. **Given** the application is installed, **When** I run the version command, **Then** the output format is concise and easily readable.

---

### User Story 2 - Retrieve Version via API (Priority: P1)

As a remote user or monitoring system, I want to retrieve the application version via an API endpoint so that I can remotely verify the version of a running service.

**Why this priority**: Critical for remote management and ensuring compatibility between different system components.

**Independent Test**: Calling the version endpoint using an HTTP client should return a JSON response with the version.

**Acceptance Scenarios**:

1. **Given** the API service is running, **When** I make a GET request to `/version`, **Then** I receive a `200 OK` response with a JSON body containing the version string.
2. **Given** the API service is running, **When** I make a GET request to `/version`, **Then** the response includes the version string.

---

### Edge Cases

- **Service Unavailable**: If the API is down, the CLI command should still work for the local binary.
- **Malformed Version**: If the version cannot be determined (e.g., development build), the system should return a sensible fallback like `0.0.0-dev`.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a CLI subcommand `version` to display the application version.
- **FR-002**: System MUST support a global flag `--version` to display the application version.
- **FR-003**: System MUST provide an HTTP GET endpoint `/version` to retrieve the application version.
- **FR-004**: The version string MUST follow semantic versioning (e.g., X.Y.Z).
- **FR-005**: The CLI output MUST be ONLY plain text.
- **FR-006**: The API response MUST be in JSON format.
- **FR-007**: The version information MUST be consistent between the CLI and the API.
- **FR-008**: The version response MUST include ONLY the semantic version string.
- **FR-009**: The API MUST listen on a port configurable via flags or environment variables (default: 8080).
- **FR-010**: The application version MUST be hardcoded in a dedicated configuration file (e.g., `internal/config/version.go`) for easy maintenance.

### Key Entities *(include if feature involves data)*

- **Version Info**: A structure containing the major, minor, and patch numbers (semantic version).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can retrieve the version via CLI in under 100ms.
- **SC-002**: API `/version` endpoint responds in under 50ms.
- **SC-003**: Version information is 100% consistent across all access methods (CLI and API).
- **SC-004**: 100% of automated health checks can successfully parse the version from the JSON response.
