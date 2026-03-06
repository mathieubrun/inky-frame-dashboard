# Tasks: Google Agenda Display

**Feature**: Google Agenda Display
**Status**: Immediately Executable
**Story Dependencies**: US1 (P1) → US2 (P2) → US3 (P3)

## Phase 1: Setup

- [x] T001 Add `google.golang.org/api/calendar/v3` and `google.golang.org/api/option` to `go.mod`
- [x] T002 Create directory structure for agenda and dashboard packages in `internal/core/agenda/` and `internal/core/dashboard/`

## Phase 2: Foundational

- [x] T003 Define `AgendaEvent` and `AgendaForecast` types in `internal/core/agenda/types.go`
- [x] T004 Define `CalendarProvider` interface in `internal/core/agenda/provider.go`
- [x] T005 Update `internal/config/config.go` to include Google API credentials and Agenda ID settings
- [x] T006 [P] Implement agenda file-based caching in `internal/core/agenda/cache.go`

## Phase 3: User Story 1 - Retrieve Agenda Events (Priority: P1)

**Goal**: Retrieve upcoming events from Google Agenda via API or CLI.
**Independent Test**: `inky agenda list --mock` returns a list of mock events.

- [x] T007 [P] [US1] Implement `MockCalendarProvider` in `internal/core/agenda/mock.go`
- [x] T008 [US1] Implement `GoogleCalendarProvider` using Service Account auth in `internal/core/agenda/google.go`
- [x] T009 [US1] Create unit tests for Google Calendar retrieval in `internal/core/agenda/google_test.go` using `httptest`
- [x] T010 [US1] Create HTTP handler for `/api/v1/agenda` in `internal/api/agenda.go`
- [x] T011 [US1] Implement `inky agenda list` command in `internal/cli/agenda.go`
- [x] T012 [US1] Create Bruno test for agenda endpoint in `bruno/Get Agenda.bru`

## Phase 4: User Story 2 - Combined Dashboard Image (Priority: P2)

**Goal**: Generate an 800x480 PNG with weather on the left and agenda on the right.
**Independent Test**: `inky dashboard image --mock --output test.png` generates a valid split-screen image.

- [x] T013 [US2] Implement combined dashboard rendering logic (split-screen) in `internal/core/dashboard/render.go`
- [x] T014 [US2] Create unit tests for combined rendering in `internal/core/dashboard/render_test.go`
- [x] T015 [US2] Create HTTP handler for `/api/v1/dashboard/image` in `internal/api/dashboard.go`
- [x] T016 [US2] Implement `inky dashboard image` command in `internal/cli/dashboard.go`
- [x] T017 [US2] Create Bruno test for dashboard image endpoint in `bruno/Get Dashboard Image.bru`

## Phase 5: User Story 3 - Server-Side Agenda Configuration (Priority: P3)

**Goal**: Configure Google Agendas via server-side settings.
**Independent Test**: Application starts and validates credentials path from `internal/config/config.go`.

- [x] T018 [US3] Ensure Google Service Account JSON path can be configured via environment variables in `internal/config/config.go`
- [x] T019 [US3] Add validation for Google credentials path during application startup in `internal/core/agenda/google.go`

## Phase 6: Polish & Cross-cutting Concerns

- [x] T020 Ensure >80% test coverage for new agenda and dashboard logic in `internal/core/`
- [x] T021 Perform end-to-end integration test with mock providers enabled using the CLI

## Implementation Strategy

- **MVP**: Complete US1 using the `MockCalendarProvider` to provide immediate CLI/API data access.
- **Incremental Delivery**: US2 adds the visual dashboard requirement; US3 ensures production-ready configuration.
- **Parallel Opportunities**: Agenda caching (T006) and Mock provider (T007) can be developed simultaneously.

## Dependency Graph

```text
Phase 1 (Setup)
  └── Phase 2 (Foundational)
        ├── Phase 3 (US1: Retrieval)
        │     └── Phase 4 (US2: Rendering)
        └── Phase 5 (US3: Config)
              └── Phase 6 (Polish)
```
