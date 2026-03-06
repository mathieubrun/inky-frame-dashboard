# Feature Specification: Inky Frame Dashboard Display

**Feature Branch**: `005-inky-display-image`  
**Created**: 2026-03-05  
**Status**: Draft  
**Input**: User description: "I want to display the image given by the dashboard on my inky frame 7.3"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Automatic Dashboard Update (Priority: P1)

As a user, I want my Inky Frame 7.3 to wake up automatically at a set interval, connect to Wi-Fi, fetch the latest weather image from the dashboard API, and display it so that I always have an up-to-date information board without manual intervention.

**Why this priority**: This is the core purpose of the device. Without automatic updates, the Inky Frame is just a static picture.

**Independent Test**: Can be fully tested by configuring the device with the dashboard API URL and observing it perform a successful update and image refresh on the display.

**Acceptance Scenarios**:

1. **Given** the device is powered and configured, **When** the scheduled update time is reached, **Then** it must successfully connect to the configured Wi-Fi and fetch the dashboard image.
2. **Given** a successful image download, **When** the rendering process completes, **Then** the Inky Frame 7.3 display must update with the new content and then return to deep sleep.

---

### User Story 2 - Visual Confirmation of Connection (Priority: P2)

As a user, I want to see a visual indicator (like an icon or brief text) when the device is connecting to Wi-Fi or fetching data, so that I know it is actively working and not frozen.

**Why this priority**: Improves user experience and provides feedback during the slow process of Wi-Fi connection and image rendering on e-ink.

**Independent Test**: Observe the display (or an onboard LED) during the connection phase to verify a status change is visible.

**Acceptance Scenarios**:

1. **Given** the device wakes up, **When** it starts the Wi-Fi connection, **Then** it provides a clear visual signal that it is "In Progress".

---

### User Story 3 - Error Feedback (Priority: P3)

As a user, I want a clear message or error screen if the Wi-Fi connection fails or the API is unreachable, so that I can troubleshoot the issue.

**Why this priority**: Essential for maintenance. Without feedback, the user won't know if the batteries died, Wi-Fi changed, or the server is down.

**Independent Test**: Can be tested by disconnecting the Wi-Fi router or disabling the API and verifying that the device displays an error message on the e-ink screen before sleeping.

**Acceptance Scenarios**:

1. **Given** the API is unreachable, **When** the fetch fails, **Then** the device displays an "Update Failed" message including the time of failure.

---

### Edge Cases

- **What happens when the battery is critically low?** The system should attempt one final display update showing a "Low Battery" warning before remaining in deep sleep until charged.
- **How does the system handle an image larger than the display memory?** The device should only request/process images matching its specific resolution (800x480).
- **What happens if the Wi-Fi signal is weak?** The system MUST have a strict timeout to prevent battery drain while attempting to connect indefinitely.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Device MUST support configuration for Wi-Fi SSID, Password, and API Endpoint URL.
- **FR-002**: Device MUST wake up from deep sleep at a configurable interval (Default: 30 minutes).
- **FR-003**: Device MUST fetch the image from the specified endpoint using a standard HTTP GET request.
- **FR-004**: System MUST render the downloaded PNG data to the 800x480 e-ink display.
- **FR-005**: Device MUST return to deep sleep immediately after a successful or failed update attempt to conserve power.
- **FR-006**: System MUST support a timeout for both Wi-Fi connection and HTTP request to prevent battery depletion.

### Key Entities *(include if feature involves data)*

- **Device Configuration**: Set of parameters (SSID, Password, API URL, Sleep Interval) stored locally on the Inky Frame.
- **Dashboard Image**: Binary PNG data (800x480) fetched from the Go API.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Device completes a full update cycle (wake, fetch, render, sleep) in under 60 seconds.
- **SC-002**: Device achieves at least 3 months of battery life with a 30-minute update interval using standard AA batteries.
- **SC-003**: 100% of successful API fetches result in a clear, correctly dithered image on the display.
- **SC-004**: Device correctly enters deep sleep mode after every attempt (Pass/Fail).
