# Tasks: Swiss Weather Forecast Integration

**Feature**: Swiss Weather Forecast Integration
**Plan**: [/specs/003-swiss-weather-forecast/plan.md](/specs/003-swiss-weather-forecast/plan.md)

## Implementation Strategy

We follow an incremental delivery approach, starting with the core weather engine and foundational providers (Mock and Cache), followed by the CLI implementation (US1) and the API endpoint (US2). Each phase delivers a complete, testable increment.

1.  **Foundational Phase**: Establish the core data models and provider interfaces. Implement the Mock and Cache providers to enable development and offline testing.
2.  **Weather Engine (Open-Meteo)**: Implement the real data fetcher using Open-Meteo with MeteoSwiss ICON-CH models and geocoding.
3.  **CLI Increment (US1)**: Deliver the `inky weather` command with formatted output.
4.  **API Increment (US2)**: Deliver the `/weather/swiss` endpoint with JSON response and Bruno integration tests.

## Phase 1: Setup

- [ ] T001 Create directory structure for weather core logic in `internal/core/weather/`
- [ ] T002 Add weather-related configuration parameters (Cache Dir, TTL, Mock toggle) to `internal/config/config.go`

## Phase 2: Foundational Weather Engine

- [ ] T003 Define weather data structures (`WeatherRecord`, `WeatherForecast`, `Location`) in `internal/core/weather/types.go`
- [ ] T004 Define the `Provider` interface in `internal/core/weather/provider.go`
- [ ] T005 [P] Implement `MockProvider` in `internal/core/weather/mock.go` returning randomized data
- [ ] T006 [P] Implement file-based caching logic in `internal/core/weather/cache.go` for `WeatherForecast` persistence
- [ ] T007 Implement `OpenMeteoProvider` in `internal/core/weather/openmeteo.go` with geocoding and ICON-CH model fetching

## Phase 3: User Story 1 - Retrieve Swiss Weather via CLI (Priority: P1)

**Goal**: Users can fetch weather via `inky weather --city "City"`
**Test Criteria**: `inky weather --city "Zurich"` displays a table with temp, wind, and rain details.

- [ ] T008 [US1] Create the `weather` command definition in `internal/cli/weather.go`
- [ ] T009 [US1] Implement CLI output formatting (table and JSON) in `internal/cli/weather.go`
- [ ] T010 [US1] Integrate `internal/core/weather` providers with the `weather` command in `internal/cli/weather.go`
- [ ] T011 [US1] Verify CLI functionality using both `--mock` and live providers

## Phase 4: User Story 2 - Retrieve Swiss Weather via API (Priority: P1)

**Goal**: Dashboard can fetch weather via `/weather/swiss?city=City`
**Test Criteria**: `GET /weather/swiss?city=Geneva` returns `200 OK` with JSON weather data.

- [ ] T012 [US2] Implement the weather handler function in `internal/api/weather.go`
- [ ] T013 [US2] Register the `/weather/swiss` route in `internal/api/server.go`
- [ ] T014 [P] [US2] Create Bruno test file at `bruno/Get Swiss Weather.bru` with response assertions
- [ ] T015 [US2] Verify API response matches the JSON contract defined in `contracts/api.md`

## Phase 5: Polish & Cross-cutting

- [ ] T016 Ensure `golangci-lint` passes for all new files
- [ ] T017 [P] Add unit tests for caching logic in `internal/core/weather/cache_test.go`
- [ ] T018 [P] Add unit tests for geocoding logic in `internal/core/weather/openmeteo_test.go`

## Dependencies

1.  **Phase 2** must be complete before **Phase 3** or **Phase 4** can be fully verified with real data.
2.  **T003 (Types)** and **T004 (Interface)** are prerequisites for all other weather tasks.
3.  **US1** and **US2** can be implemented in parallel once Phase 2 is complete.

## Parallel Execution Examples

### Team A (Core & CLI)
- T003 -> T004 -> T007 -> T008 -> T010

### Team B (Mocks & API)
- T005 -> T012 -> T013 -> T014
