# Tasks: Battery Level Monitoring

**Input**: Design documents from `/specs/007-battery-level-monitoring/`
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

**Purpose**: Project initialization and basic structure

- [ ] T001 Update `internal/config/config.go` to include `BatteryCSVPath` with default value `.inky/battery.csv`
- [ ] T002 Update `internal/config/config.go` to load `BatteryCSVPath` from environment and flags
- [ ] T003 Create directory `internal/core/battery/` for core logic
- [ ] T004 [P] Create `bruno/` entries for battery endpoints

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure for thread-safe CSV management

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 [P] Define `BatteryReport` struct in `internal/core/battery/types.go`
- [ ] T006 Implement thread-safe CSV append logic with `sync.Mutex` in `internal/core/battery/storage.go`
- [ ] T007 [P] Implement CSV read logic in `internal/core/battery/storage.go`
- [ ] T008 [P] Unit tests for CSV storage logic in `internal/core/battery/storage_test.go` (verify Mutex and format)

**Checkpoint**: Foundation ready - battery storage and retrieval logic is verified.

---

## Phase 3: User Story 1 - Device Battery Reporting (Priority: P1) 🎯 MVP

**Goal**: Enable Inky Frame to report battery voltage to the server for persistence.

**Independent Test**: Send a POST request to `/api/v1/battery` and verify the value appears in the CSV file.

### Tests for User Story 1 (REQUIRED) ⚠️

- [ ] T009 [P] [US1] Unit test for report validation logic in `internal/core/battery/battery_test.go`
- [ ] T010 [P] [US1] Integration test for POST handler in `internal/api/battery_test.go`
- [ ] T011 [P] [US1] Create Bruno file `bruno/Report Battery.bru` with POST example
- [ ] T012 [P] [US1] Verify coverage >80% for US1 logic (`go test ./internal/core/battery/... ./internal/api/... -cover`)

### Implementation for User Story 1

- [ ] T013 [P] [US1] Implement report validation and processing in `internal/core/battery/battery.go`
- [ ] T014 [US1] Implement `BatteryReportHandler` (POST) in `internal/api/battery.go`
- [ ] T015 [US1] Register POST `/api/v1/battery` route in `internal/api/server.go`
- [ ] T016 [US1] Add logging for received battery reports in `internal/api/battery.go`

**Checkpoint**: User Story 1 is functional. Devices can report battery levels.

---

## Phase 4: User Story 2 - Viewing Current Battery Status (Priority: P2)

**Goal**: Allow users/apps to query the latest reported battery level.

**Independent Test**: Query `/api/v1/battery/status` and receive the most recent report as JSON.

### Tests for User Story 2 (REQUIRED) ⚠️

- [ ] T017 [P] [US2] Unit test for "get latest" logic in `internal/core/battery/storage_test.go`
- [ ] T018 [P] [US2] Integration test for GET status handler in `internal/api/battery_test.go`
- [ ] T019 [P] [US2] Create Bruno file `bruno/Get Battery Status.bru`
- [ ] T020 [P] [US2] Verify coverage >80% for US2 logic

### Implementation for User Story 2

- [ ] T021 [P] [US2] Implement `GetLatestReport` logic in `internal/core/battery/storage.go`
- [ ] T022 [US2] Implement `BatteryStatusHandler` (GET) in `internal/api/battery.go` returning JSON
- [ ] T023 [US2] Register GET `/api/v1/battery/status` route in `internal/api/server.go`

**Checkpoint**: User Story 2 is functional. Latest battery status is accessible via API.

---

## Phase 5: User Story 3 - Battery Level History (Priority: P3)

**Goal**: Retrieve the full history of battery reports for analysis.

**Independent Test**: Query `/api/v1/battery/history` and receive raw CSV text. Run `inky battery history` from CLI.

### Tests for User Story 3 (REQUIRED) ⚠️

- [ ] T024 [P] [US3] Unit test for "get history" logic in `internal/core/battery/storage_test.go`
- [ ] T025 [P] [US3] Integration test for GET history handler in `internal/api/battery_test.go`
- [ ] T026 [P] [US3] Functional test for CLI command in `internal/cli/battery_test.go`
- [ ] T027 [P] [US3] Create Bruno file `bruno/Get Battery History.bru`
- [ ] T028 [P] [US3] Verify coverage >80% for US3 logic

### Implementation for User Story 3

- [ ] T029 [P] [US3] Implement `GetHistoryRaw` logic in `internal/core/battery/storage.go`
- [ ] T030 [US3] Implement `BatteryHistoryHandler` (GET) in `internal/api/battery.go` returning raw text
- [ ] T031 [US3] Register GET `/api/v1/battery/history` route in `internal/api/server.go`
- [ ] T032 [US3] Implement `battery` command and `history` subcommand in `internal/cli/battery.go`
- [ ] T033 [US3] Implement `clear` subcommand in `internal/cli/battery.go` to reset CSV

**Checkpoint**: User Story 3 is functional. History is accessible via API and CLI.

---

## Phase 6: Firmware Integration (Inky Frame)

**Purpose**: Update the MicroPython code to actually report the battery level.

- [ ] T034 [US1] Add `report_battery` function to `firmware/main.py`
- [ ] T035 [US1] Integrate `report_battery` call before dashboard update in `firmware/main.py`
- [ ] T036 [US1] Add configuration for API URL in `firmware/env.template.py`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements and final validation

- [ ] T037 [P] Update `README.md` with battery monitoring API details
- [ ] T038 Final code cleanup and `gofmt` verification
- [ ] T039 [P] Final verification of project-wide test coverage (>80%)
- [ ] T040 Run `quickstart.md` validation steps

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup (T001-T004). BLOCKS all user stories.
- **User Story 1 (P1)**: Depends on Foundational (Phase 2).
- **User Story 2 (P2)**: Depends on Foundational (Phase 2). Can be done in parallel with US1.
- **User Story 3 (P3)**: Depends on Foundational (Phase 2). Can be done in parallel with US1/US2.
- **Firmware (Phase 6)**: Depends on US1 (P1) API implementation.
- **Polish (Phase 7)**: Depends on all user stories.

### Parallel Opportunities

- T004, T005, T007 can be done in parallel.
- All test tasks (T009-T012, T017-T020, T024-T028) can run in parallel within their stories.
- API Handlers for US1, US2, and US3 are largely independent once core logic is ready.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup and Foundational phases.
2. Implement US1 (Reporting) to start collecting data.
3. Validate data is saved correctly in `.inky/battery.csv`.

### Incremental Delivery

1. Foundation → Data model and storage ready.
2. US1 → Device can send data.
3. US2 → User can see "now".
4. US3 → User can see "history".
5. Firmware → Close the loop.

---

## Notes

- All battery operations must respect the `sync.Mutex` in `internal/core/battery/storage.go`.
- API handlers should use the logic from `internal/core/battery`.
- CLI commands should also use the logic from `internal/core/battery` for consistency.
