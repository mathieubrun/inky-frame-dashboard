# Feature Specification: Battery Refresh Optimization

**Feature Branch**: `008-battery-refresh-optimization`
**Created**: 2026-03-07
**Status**: Draft
## Input: User description: "in order to preserve battery, the number of screen refreshes must be limited. Only refresh if the dashboard image has changed. The weather data is refreshed every day at 04:00. Calendar data is refreshed when calendar is updated."

## Clarifications
### Session 2026-03-07
- Q: How does the server detect if the calendar is updated? → A: When the inky frame makes a request, the server checks google calendar.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Conditional Screen Refresh (Priority: P1)

As an Inky Frame device, I want to skip a full screen refresh if the dashboard image content hasn't changed since the last update, so that I can significantly extend my battery life.

**Why this priority**: Core requirement for battery preservation. E-ink refreshes are the most power-consuming operation on the device.

**Independent Test**: Can be tested by requesting the dashboard image twice with the same data and verifying that the device (or a simulator) does not trigger a refresh the second time.

**Acceptance Scenarios**:

1. **Given** the server has a generated dashboard image, **When** the device requests the image and it is identical to the one previously displayed, **Then** the device should not perform a physical screen refresh.
2. **Given** the server has a new dashboard image (due to data changes), **When** the device requests the image, **Then** the device MUST perform a full screen refresh.

---

### User Story 2 - Daily Weather Update (Priority: P2)

As a user, I want my weather information to be updated once a day at 04:00 so that I have the forecast for the day when I wake up, while minimizing refreshes during the day.

**Why this priority**: High value for battery saving by limiting periodic weather polls to a single daily event.

**Independent Test**: Monitor server logs or device wake-ups to ensure weather-driven refreshes only occur at or shortly after 04:00.

**Acceptance Scenarios**:

1. **Given** it is before 04:00, **When** the device wakes up for a routine check, **Then** it should not trigger a weather update or screen refresh based on weather changes.
2. **Given** it is 04:00 or later (first wake-up of the day), **When** the device checks for updates, **Then** the weather data should be refreshed and the screen updated if the forecast has changed.

---

### User Story 3 - Calendar-Driven Refresh (Priority: P3)

As a user, I want my dashboard to refresh whenever my calendar is updated so that I always see my current appointments without unnecessary periodic polling.

**Why this priority**: Ensures data freshness for dynamic content while remaining "lazy" about refreshes when no changes occur.

**Independent Test**: Update a calendar event and verify the dashboard refreshes on the next device check-in.

**Acceptance Scenarios**:

1. **Given** the calendar data has not changed, **When** the device checks for updates, **Then** no refresh should occur.
2. **Given** a new event has been added to the calendar, **When** the device checks for updates, **Then** the dashboard should be refreshed to show the new event.

---

### Edge Cases

- **Partial Image Changes**: How does the system handle subtle changes (e.g., a clock or timestamp) that might trigger "accidental" refreshes? (Assumption: Server-side rendering will handle this by only changing the image if significant data changes).
- **Communication Failure**: What happens if the device cannot reach the server to check for changes? (Default: Stay asleep and retry later to save battery).
- **Manual Override**: Should the user be able to force a refresh if they suspect the "lazy" logic is stuck?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a mechanism to detect if the dashboard image has changed since the last delivery to the device (e.g., via ETag or Hash).
- **FR-002**: The device MUST only perform a physical screen refresh if the server indicates new content is available.
- **FR-003**: The server MUST refresh weather data precisely once per day at 04:00 local time.
- **FR-004**: The server MUST fetch fresh calendar data from Google Calendar during every device request to determine if an update is needed.
- **FR-005**: The server MUST maintain a record of the last "significant" image hash sent to the device.
- **FR-006**: The device MUST strictly enforce a 60-minute wake interval to check for these changes.

### Key Entities *(include if feature involves data)*

- **Dashboard State**: Represents the current version/content of the dashboard.
  - Image Hash (MD5/SHA1 of the rendered image)
  - Last Update Timestamp
  - Source Data Version (Weather/Calendar)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Reduce the number of screen refreshes by at least 70% for a typical user with infrequent calendar changes.
- **SC-002**: Weather data is never older than 24 hours plus the wake interval.
- **SC-003**: Calendar updates are reflected on the display within the next scheduled wake-up interval (e.g., within 30 minutes).
- **SC-004**: The check-for-change request from the device to the server MUST complete in under 500ms to minimize Wi-Fi on-time.
