---
description: "Task list for Get App Version implementation"
---

# Tasks: Get App Version

**Input**: Design documents from `/specs/001-get-app-version/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Source Code**: `cmd/inky/`, `internal/api/`, `internal/cli/`, `internal/core/`, `internal/config/`
- **Tests**: Next to source files as `*_test.go` (e.g., `internal/core/[file]_test.go`)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan (`cmd/inky`, `internal/{api,cli,core,config}`)
- [ ] T002 Initialize Go project using `go mod init` and add dependencies (spf13/cobra, spf13/viper)
- [ ] T003 [P] Configure `golangci-lint` for linting
- [ ] T004 [P] Setup base test utilities for `go test`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 [P] Create `internal/config/version.go` with hardcoded SemVer string `1.0.0`
- [ ] T006 [P] Implement shared `VersionInfo` entity in `internal/core/version.go` (if needed for shared logic)
- [ ] T007 [P] Setup base Cobra command root in `internal/cli/root.go` and `cmd/inky/main.go`
- [ ] T008 Configure error handling and logging infrastructure

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Retrieve Version via CLI Subcommand (Priority: P1) 🎯 MVP

**Goal**: Users can retrieve the application version using `inky version`.

**Independent Test**: Running `./inky version` prints `1.0.0` and returns exit code 0.

### Implementation for User Story 1

- [ ] T009 [P] [US1] Create unit test for CLI version output in `internal/cli/version_test.go`
- [ ] T010 [US1] Implement `version` subcommand in `internal/cli/version.go` using `internal/config/version.go`
- [ ] T011 [US1] Register `version` subcommand in Root command within `internal/cli/root.go`
- [ ] T012 [US1] Verify `inky version` prints raw numbers and returns exit code 0

**Checkpoint**: At this point, User Story 1 is fully functional and testable independently.

---

## Phase 4: User Story 2 - Retrieve Version via API (Priority: P1)

**Goal**: Users can retrieve the application version via HTTP `GET /version`.

**Independent Test**: `curl http://localhost:8080/version` returns `{"version": "1.0.0"}` with 200 OK.

### Implementation for User Story 2

- [ ] T013 [P] [US2] Create integration test for `/version` endpoint in `internal/api/version_test.go`
- [ ] T014 [US2] Implement `VersionHandler` in `internal/api/version.go` returning JSON object with version
- [ ] T015 [US2] Setup HTTP server in `internal/api/server.go` and register `/version` route
- [ ] T016 [US2] Implement `serve` subcommand in `internal/cli/serve.go` to start the API server
- [ ] T017 [US2] Configure port retrieval via `viper` (flags/env) in `internal/config/config.go`
- [ ] T018 [US2] Add standard access logging for `/version` endpoint

**Checkpoint**: At this point, User Stories 1 AND 2 work independently.

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T019 [P] Documentation updates in `README.md` regarding version retrieval
- [ ] T020 Code cleanup and refactoring
- [ ] T021 [P] Ensure no global `--version` flag exists in `spf13/cobra` default configuration
- [ ] T022 Run `quickstart.md` validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User Story 1 (CLI) and User Story 2 (API) can proceed in parallel once foundation is ready.
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Parallel Opportunities

- Foundation: T005, T006, T007 can run in parallel.
- CLI (Phase 3) and API (Phase 4) can run in parallel.
- Within CLI: T009 can be started first.
- Within API: T013 can be started first.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test `inky version` independently.

### Incremental Delivery

1. Foundation ready.
2. Add CLI version retrieval (US1).
3. Add API version retrieval (US2).
4. Each story adds value without breaking the other.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Verify SemVer format: raw numbers only (X.Y.Z)
- Access logging is mandatory for API endpoint
