# Tasks: Battery Refresh Optimization

**Input**: Design documents from `/specs/008-battery-refresh-optimization/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Comprehensive tests are REQUIRED for all new features and bug fixes as per the Project Constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Source Code**: `cmd/inky/`, `internal/api/`, `internal/cli/`, `internal/core/`, `internal/config/`
- **Tests**: Next to source files as `*_test.go` (e.g., `internal/core/[file]_test.go`)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and base utilities

- [X] T001 [P] Implement MD5 hashing utility for image data in `internal/core/render.go`
- [X] T002 [P] Unit test for MD5 utility in `internal/core/render_test.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core logic for weather scheduling and cache invalidation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Implement "Last 04:00 AM" freshness logic in `internal/core/weather/cache.go`
- [X] T004 [P] Unit tests for "Last 04:00 AM" logic in `internal/core/weather/cache_test.go`
- [X] T005 Update `internal/core/weather/cache.go` to use the new freshness check instead of fixed TTL

**Checkpoint**: Foundation ready - weather data now follows the daily refresh schedule.

---

## Phase 3: User Story 1 - Conditional Screen Refresh (Priority: P1) 🎯 MVP

**Goal**: Skip physical screen refresh if the dashboard image hasn't changed using ETag/304.

**Independent Test**: Request the same image twice with `If-None-Match` header and verify 304 Not Modified.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T006 [P] [US1] Unit test for ETag generation and 304 handling in `internal/api/dashboard_test.go`
- [X] T007 [P] [US1] Create Bruno file `bruno/Get Dashboard Image with ETag.bru` with ETag assertion
- [X] T008 [P] [US1] Verify coverage >80% for ETag logic (`go test ./internal/api/... -cover`)

### Implementation for User Story 1

- [X] T009 [US1] Integrate MD5 hashing into `DashboardImageHandler` in `internal/api/dashboard.go`
- [X] T010 [US1] Implement `If-None-Match` check and `StatusNotModified` response in `internal/api/dashboard.go`
- [X] T011 [US1] Update `firmware/main.py` to read/write ETag from `/etag.txt`
- [X] T012 [US1] Update `firmware/main.py` `fetch_image` function to send `If-None-Match` header
- [X] T013 [US1] Update `firmware/main.py` `main` loop to skip `render_image` and `display.update` on 304

**Checkpoint**: User Story 1 is functional. Devices skip refreshes for identical content.

---

## Phase 4: User Story 2 - Daily Weather Update (Priority: P2)

**Goal**: Ensure weather is updated once a day at 04:00 server-side.

**Independent Test**: Simulate a request before and after 04:00 AM and verify weather fetch occurs only after 04:00 AM.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T014 [P] [US2] Unit test for weather fetch timing in `internal/core/weather/cache_test.go` (mocking time)
- [X] T015 [P] [US2] Verify coverage >80% for daily refresh logic
- [X] T016 [US2] Refactor `DashboardImageHandler` in `internal/api/dashboard.go` to ensure weather provider uses the daily refresh logic
- [X] T017 [US2] Add logging for weather cache invalidation at 04:00 AM

**Checkpoint**: User Story 2 is functional. Weather refreshes follow the 04:00 AM schedule.

---

## Phase 5: User Story 3 - Calendar-Driven Refresh (Priority: P3)

**Goal**: Dashboard refreshes when Google Calendar changes.

**Independent Test**: Modify a calendar event and verify the Inky Frame receives a 200 OK (new ETag) instead of 304.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T018 [P] [US3] Unit test for agenda-driven image change in `internal/api/dashboard_test.go`
- [X] T019 [P] [US3] Verify coverage >80% for live calendar check logic
- [X] T020 [US3] Update `DashboardImageHandler` in `internal/api/dashboard.go` to perform synchronous calendar fetch
- [X] T021 [US3] Ensure ETag is calculated AFTER calendar data is merged into the image

**Checkpoint**: User Story 3 is functional. Calendar changes trigger screen updates on next poll.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and documentation

- [X] T022 Update `firmware/env.template.py` with default `SLEEP_MINUTES = 60` per spec FR-006
- [X] T023 [P] Update `README.md` with ETag/304 optimization details
- [X] T024 [P] Final verification of project-wide test coverage (>80%)
- [X] T025 Run `quickstart.md` validation steps

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Phase 1 completion.
- **User Story 1 (P1)**: Depends on Phase 2 completion.
- **User Story 2 (P2)**: Depends on Phase 2 completion.
- **User Story 3 (P3)**: Depends on User Story 1 completion.
- **Polish (Phase 6)**: Depends on all user stories being complete.

### User Story Dependencies

- **US1** is the primary driver for energy savings.
- **US2** relies on the foundational weather cache logic.
- **US3** relies on the ETag mechanism established in US1 to detect image differences.

---

## Parallel Opportunities

- T001 and T003 can be worked on in parallel.
- All test tasks marked [P] can run in parallel within their respective stories.
- Firmware updates (T011-T013) can happen in parallel with API handler updates (T009-T010).

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup and Foundational logic for MD5.
2. Implement ETag/304 in API.
3. Implement ETag handling in Firmware.
4. **VALIDATE**: Verify 304 response and skipped refresh.

### Incremental Delivery

1. Weather 04:00 AM logic → Ensures data freshness schedule.
2. Calendar sync check → Ensures responsiveness to events.
3. Firmware 60-minute interval → Maximizes battery life.

---

## Notes

- The 60-minute wake interval is a client-side setting (`SLEEP_MINUTES` in `env.py`).
- All server-side time logic should use local time or a consistent timezone as per user preference (defaulting to Zurich per handlers).
