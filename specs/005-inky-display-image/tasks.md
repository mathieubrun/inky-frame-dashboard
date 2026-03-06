# Tasks: Inky Frame Dashboard Display

**Feature**: Inky Frame Dashboard Display
**Status**: Immediately Executable
**Story Dependencies**: US1 (P1) → US2 (P2) → US3 (P3)

## Phase 1: Setup

- [ ] T001 [P] Create initial MicroPython environment structure with `main.py` and `.env.py` template
- [ ] T002 Define constants and hardware pin mappings in `main.py` per research (VSYS=29, Hold=2, ReadEn=25)

## Phase 2: Foundational

- [ ] T003 Implement battery voltage monitoring using ADC(29) in `main.py`
- [ ] T004 Implement Wi-Fi connection handler with timeout and error retry logic in `main.py`
- [ ] T005 [P] Implement hardware-managed deep sleep using `inky_frame.sleep_for()` in `main.py`

## Phase 3: User Story 1 - Automatic Dashboard Update (Priority: P1)

**Goal**: Fetch and display the weather image automatically on a 30-minute interval.
**Independent Test**: Device wakes from sleep, connects to Wi-Fi, downloads PNG, renders to e-ink, and returns to sleep.

- [ ] T006 [US1] Implement HTTP GET request using `urequests` to fetch dashboard PNG in `main.py`
- [ ] T007 [US1] Implement PNG decoding and rendering using `pngdec` and `picographics` in `main.py`
- [ ] T008 [US1] Integrate full update cycle (Wake -> Fetch -> Render -> Sleep) in `main.py`
- [ ] T009 [US1] Implement memory management with `gc.collect()` before and after rendering in `main.py`

## Phase 4: User Story 2 - Visual Confirmation of Connection (Priority: P2)

**Goal**: Provide visual feedback during the update process.
**Independent Test**: Onboard LEDs blink or change state during Wi-Fi connection and image download.

- [ ] T010 [US2] Implement LED status signaling for "Connecting", "Downloading", and "Rendering" phases in `main.py`

## Phase 5: User Story 3 - Error Feedback (Priority: P3)

**Goal**: Display clear error messages on the e-ink screen if an update fails.
**Independent Test**: Disconnecting Wi-Fi results in an "Update Failed" screen being rendered before sleep.

- [ ] T011 [US3] Implement basic text-based error screen rendering in `main.py`
- [ ] T012 [US3] Add error handling catch-all to ensure device sleeps even if an exception occurs in `main.py`

## Phase 6: Polish & Cross-cutting Concerns

- [ ] T013 Implement low battery overlay (icon or text) in top-right corner if voltage < threshold in `main.py`
- [ ] T014 Perform final end-to-end integration test with the Go API server

## Implementation Strategy

- **MVP**: Complete US1 to achieve the core functionality of a periodic information board.
- **Incremental Delivery**: US2 adds operational feedback; US3 adds robustness and troubleshooting capability.
- **Parallel Opportunities**: US2 (LED signaling) can be developed independently once the main state machine is defined.
