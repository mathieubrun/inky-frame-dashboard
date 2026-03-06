# Feature Specification: Google Agenda Display

**Feature Branch**: `006-google-agenda-display`  
**Created**: 2026-03-06  
**Status**: Draft  
**Input**: User description: "i want to display information from google agendas. the agendas are configured on server side. agendas can be queried by api or cli. for the dashboard image, weather information is displayed on the left, agenda information on the right."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Retrieve Agenda Events (Priority: P1)

As a user or automated system, I want to query upcoming events from configured Google Agendas via the API or CLI so that I can verify the data is being correctly synchronized.

**Why this priority**: This is the foundation of the feature. Data must be retrievable before it can be rendered into an image.

**Independent Test**: Can be fully tested by running a CLI command or API request and verifying that a list of upcoming events is returned in a readable format.

**Acceptance Scenarios**:

1. **Given** one or more Google Agendas are configured, **When** a retrieval request is made, **Then** the system returns a list of events scheduled for the next 24 hours.
2. **Given** an invalid or unauthorized agenda configuration, **When** a retrieval is attempted, **Then** the system returns a clear error message.

---

### User Story 2 - Combined Dashboard Image (Priority: P2)

As a dashboard device, I want to fetch a single image containing both weather and agenda information so that I can display a comprehensive overview of my day.

**Why this priority**: This fulfills the primary visual requirement of the feature.

**Independent Test**: Can be tested by requesting the dashboard image and verifying that weather info is on the left half and agenda info is on the right half.

**Acceptance Scenarios**:

1. **Given** active weather and agenda data, **When** the dashboard image is requested, **Then** an 800x480 PNG is returned.
2. **Given** the generated image, **When** viewed, **Then** the left 50% displays weather data and the right 50% displays a list of upcoming calendar events.

---

### User Story 3 - Server-Side Agenda Configuration (Priority: P3)

As an administrator, I want to configure which Google Agendas are tracked by editing a configuration file on the server so that I don't have to manage credentials on the display device.

**Why this priority**: Centralizes security and configuration, following the principle of logic offloading.

**Independent Test**: Can be tested by adding a new agenda ID to the server configuration and verifying it appears in the API/CLI output.

**Acceptance Scenarios**:

1. **Given** access to the server, **When** a new Google Agenda ID is added to the configuration, **Then** the system includes that agenda's events in subsequent queries.

---

### Edge Cases

- **What happens when an agenda has no upcoming events?** The right side of the dashboard should display a "No upcoming events" message.
- **How are overlapping events handled?** Events should be displayed in chronological order; overlapping events are listed sequentially.
- **What if Google API is unreachable?** The system should display a "Calendar Unavailable" message on the right side, preserving the weather display on the left.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support server-side configuration of Google Agenda IDs and credentials.
- **FR-002**: System MUST authenticate with Google APIs using a Service Account.
- **FR-003**: System MUST provide a CLI command to list upcoming agenda events.
- **FR-004**: System MUST provide an API endpoint to retrieve agenda data in JSON format.
- **FR-005**: System MUST generate a combined PNG image (800x480) with a vertical split layout.
- **FR-006**: The left panel MUST display weather information (Temperature, Condition, Location).
- **FR-007**: The right panel MUST display up to 8 upcoming events.
- **FR-008**: Each event entry MUST include its Title and Start Time.

### Key Entities *(include if feature involves data)*

- **Agenda Event**: Represents a single calendar entry (Summary, StartTime, EndTime, Location).
- **Agenda Configuration**: Server-side settings mapping an ID to a specific Google Calendar and its access credentials.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Agenda data is retrieved from Google APIs in under 2 seconds.
- **SC-002**: Combined dashboard images are generated and served in under 3 seconds.
- **SC-003**: The right panel displays at least 5 upcoming events if they exist.
- **SC-004**: 100% of event titles are legible on an Inky Frame 7.3" display.
