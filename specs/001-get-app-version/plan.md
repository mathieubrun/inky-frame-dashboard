# Implementation Plan: Get App Version

**Branch**: `003-get-app-version` | **Date**: 2026-03-05 | **Spec**: [specs/003-get-app-version/spec.md](spec.md)
**Input**: Feature specification from `/specs/003-get-app-version/spec.md`

## Summary
Implement a unified way to retrieve the application version through a CLI subcommand (`inky version`) and an HTTP API endpoint (`GET /version`). The version will be a hardcoded semantic string in `internal/config/version.go`. Following user refinement, the global flag `--version` is removed.

## Technical Context

**Language/Version**: Go (Golang) 1.22+ (Managed by `go mod`)
**Primary Dependencies**: standard library `net/http`, `spf13/cobra`, `spf13/viper`
**Storage**: N/A
**Testing**: standard library `testing`
**Target Platform**: Raspberry Pi Pico W (Inky Frame) + Linux Server (API)
**Project Type**: CLI + web-service
**Performance Goals**: CLI response < 100ms, API response < 50ms
**Constraints**: Single binary release, raw numbers only (X.Y.Z), hardcoded in `internal/config/version.go`, port configurable (default 8080), public API access, standard access logging. Global flag `--version` is explicitly excluded.
**Scale/Scope**: Minimal metadata retrieval.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Logic Offloading**: ✅ Metadata retrieval is lightweight.
- **II. Energy-First Lifecycle**: ✅ Fast responses minimize wake time for clients.
- **III. Data Integrity & Freshness**: ✅ Hardcoded version ensures consistency.
- **IV. Resource-Conscious Image Delivery**: ✅ N/A (JSON/Text response).
- **V. API-First Development**: ✅ Defined /version contract in spec.
- **VI. Tooling Consistency**: ✅ Using standard Go tools.
- **VII. Modular & Unified Architecture**: ✅ Unified single binary with CLI subcommand parity. Refined to subcommand only per user request.
- **VIII. Flexible Configuration**: ✅ Configurable port via flags/env.

## Project Structure

### Documentation (this feature)

```text
specs/003-get-app-version/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── version.md
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
cmd/
└── inky/
    └── main.go          # CLI entry point updated

internal/
├── config/
│   └── version.go       # Hardcoded version string
├── cli/
│   └── version.go       # version subcommand implementation
└── api/
    └── version.go       # GET /version handler
```

**Structure Decision**: Version string is centralized in `internal/config/version.go`. The CLI uses `cobra` subcommands. The API uses a standard `http.Handler`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
