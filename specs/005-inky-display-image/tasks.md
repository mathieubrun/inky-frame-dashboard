# Tasks: Inky Frame Dashboard Display

**Feature**: Inky Frame Dashboard Display
**Status**: Immediately Executable
**Story Dependencies**: US1 (P1) → US2 (P2) → US3 (P3)

## Phase 1: Setup

- [ ] T001 Create `firmware/` directory for MicroPython code
- [ ] T002 [P] Create `firmware/.env.py` template with configuration variables
- [ ] T003 [P] Create initial `firmware/main.py` with hardware constants and imports

## Phase 2: Foundational

- [ ] T004 Implement `get_battery_voltage()` using ADC(29) in `firmware/main.py`
- [ ] T005 Implement `connect_wifi()` with SSID/Password from `env` in `firmware/main.py`
- [ ] T006 Implement hardware-managed deep sleep using `inky_frame.sleep_for()` in `firmware/main.py`

## Phase 3: User Story 1 - Automatic Dashboard Update (Priority: P1)

**Goal**: Wake every 30 minutes, fetch weather image, and render to e-ink.
**Independent Test**: Device wakes, connects, downloads PNG, and updates screen correctly.

- [ ] T007 [US1] Implement `fetch_image()` using `urequests` in `firmware/main.py`
- [ ] T008 [US1] Implement `render_image()` using `pngdec` and `picographics` in `firmware/main.py`
- [ ] T009 [US1] Orchestrate update cycle (Wake -> Battery -> WiFi -> Fetch -> Render -> Sleep) in `firmware/main.py`
- [ ] T010 [US1] Implement memory safety with `gc.collect()` before/after rendering in `firmware/main.py`

## Phase 4: User Story 2 - Visual Confirmation of Connection (Priority: P2)

**Goal**: Provide visual feedback during active update phases.
**Independent Test**: Onboard LEDs signal "WiFi Connecting" and "Downloading" states.

- [ ] T011 [US2] Implement LED signaling for WiFi and Download phases in `firmware/main.py`

## Phase 5: User Story 3 - Error Feedback (Priority: P3)

**Goal**: Display troubleshooting info on screen if update fails.
**Independent Test**: Turning off server results in "Fetch Failed" screen on Inky Frame.

- [ ] T012 [US3] Implement `draw_error_screen()` using `picographics` text in `firmware/main.py`
- [ ] T013 [US3] Add try/except block around update logic to ensure sleep on error in `firmware/main.py`

## Phase 6: Polish & Cross-cutting Concerns

- [ ] T014 Implement low battery overlay (red rectangle/text) if voltage < threshold in `firmware/main.py`
- [ ] T015 Verify final end-to-end integration between `firmware/` and Go API server

## Implementation Strategy

- **MVP**: Complete US1 to have a functional, periodic dashboard.
- **Incremental Delivery**: US2 adds UX feedback; US3 adds reliability and diagnostics.
- **Parallel Opportunities**: LED signaling (T011) and Error screen (T012) can be worked on once the main flow is stable.
