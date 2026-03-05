# Tasks: Weather Data Image Generation

**Feature**: Weather Data Image Generation
**Status**: Immediately Executable
**Story Dependencies**: US1 (P1) → US2 (P2) → US3 (P3)

## Phase 1: Setup

- [ ] T001 Add `github.com/fogleman/gg` and `golang.org/x/image` to `go.mod`
- [ ] T002 Create directory for embedded weather icons at `internal/core/weather/icons/`

## Phase 2: Foundational

- [ ] T003 Define `ImageRequest` and `RenderTemplate` structures in `internal/core/weather/types.go`
- [ ] T004 [P] Implement file-based image caching logic in `internal/core/weather/cache.go`
- [ ] T005 Create base rendering service structure in `internal/core/weather/image.go`

## Phase 3: User Story 1 - Retrieve Weather Dashboard Image (Priority: P1)

**Goal**: Fetch a pre-rendered PNG image of the current weather forecast via an API endpoint.
**Independent Test**: `curl -o test.png "http://localhost:8080/api/v1/weather/image?location=Zurich"` returns a valid PNG.

- [ ] T006 [US1] Create unit tests for core weather image rendering in `internal/core/weather/image_test.go`
- [ ] T007 [US1] Implement weather image rendering logic using `gg` library in `internal/core/weather/image.go`
- [ ] T008 [US1] Create HTTP handler for weather image endpoint in `internal/api/weather.go`
- [ ] T009 [US1] Register `/api/v1/weather/image` route in `internal/api/server.go`
- [ ] T010 [US1] Create Bruno test configuration at `bruno/Get Weather Image.bru`

## Phase 4: User Story 2 - Localized Weather Image (Priority: P2)

**Goal**: Support location lookup via city name or postcode for image generation.
**Independent Test**: `inky weather image "8001" --output test.png` generates a local image for the specific postcode.

- [ ] T011 [US2] Update `OpenMeteoProvider.geocode` to support postcode strings in `internal/core/weather/openmeteo.go`
- [ ] T012 [US2] Add unit tests for postcode-to-coordinate resolution in `internal/core/weather/openmeteo_test.go`
- [ ] T013 [US2] Implement `inky weather image` command in `internal/cli/weather.go` using shared rendering logic

## Phase 5: User Story 3 - Visual Consistency and Clarity (Priority: P3)

**Goal**: Optimize image layout and colors for the 6-color Inky e-ink display.
**Independent Test**: Generated images use only Black, White, Red, Green, Blue, and Yellow pigments.

- [ ] T014 [US3] Implement Spectra 6 color palette mapping and dithering in `internal/core/weather/image.go`
- [ ] T015 [US3] Optimize text contrast and font sizes for 800x480 resolution in `internal/core/weather/image.go`
- [ ] T016 [US3] Add and embed optimized weather icons (PNG) in `internal/core/weather/icons/`

## Phase 6: Polish & Cross-cutting Concerns

- [ ] T017 Ensure >80% test coverage for all new rendering and caching logic in `internal/core/weather/`
- [ ] T018 Perform final integration test between API, CLI, and Open-Meteo provider

## Implementation Strategy

- **MVP**: Complete US1 to provide the basic image endpoint.
- **Incremental Delivery**: US2 adds flexibility for users; US3 ensures the product is optimized for the target hardware.
- **Parallel Opportunities**: Cache implementation (T004) can be done in parallel with renderer setup (T005).
