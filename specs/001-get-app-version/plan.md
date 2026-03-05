# Implementation Plan: Get App Version

**Branch**: `003-get-app-version` | **Date**: 2026-03-05 | **Spec**: [/specs/003-get-app-version/spec.md](spec.md)
**Input**: Feature specification from `/specs/003-get-app-version/spec.md`

## Summary
Implement a unified way to retrieve the application version through both a CLI command (`inky version` and `inky --version`) and an HTTP API endpoint (`GET /version`). The version will be a hardcoded semantic string in `internal/config/version.go`.

## Technical Context

**Language/Version**: Go (Golang) 1.22+ (Managed by `go mod`)
**Primary Dependencies**: standard library `net/http`, `spf13/cobra`, `spf13/viper`
**Storage**: N/A
**Testing**: standard library `testing`
**Target Platform**: Raspberry Pi Pico W (Inky Frame) + Linux Server (API)
**Project Type**: CLI + web-service
**Performance Goals**: CLI response < 100ms, API response < 50ms
**Constraints**: Single binary release, raw numbers only (X.Y.Z), hardcoded in `internal/config/version.go`, port configurable (default 8080), public API access, standard access logging.
**Scale/Scope**: Minimal metadata retrieval.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Logic Offloading**: ✅ N/A (this is a metadata feature)
- **II. Energy-First Lifecycle**: ✅ Fast responses minimize wake time.
- **III. Data Integrity & Freshness**: ✅ N/A (static version)
- **IV. Resource-Conscious Image Delivery**: ✅ N/A (JSON/Text response)
- **V. API-First Development**: ✅ Defined /version contract.
- **VI. Tooling Consistency**: ✅ Using standard Go tools.
- **VII. Modular & Unified Architecture**: ✅ Unified single binary with CLI and API parity.
- **VIII. Flexible Configuration**: ✅ Configurable port via flags/env.

## Project Structure

### Documentation (this feature)

```text
specs/003-get-app-version/
├── plan.md              # This file
├── research.md          # Research on version injection and cobra integration
├── data-model.md        # VersionInfo entity
├── quickstart.md        # How to test CLI/API version retrieval
├── contracts/           # API and CLI interface definitions
│   └── version.md
└── tasks.md             # Implementation tasks (Phase 2 output)
```

### Source Code (repository root)

```text
cmd/
└── inky/
    └── main.go          # Root command configuration

internal/
├── config/
│   └── version.go       # Hardcoded version string
├── cli/
│   └── version.go       # version subcommand and root flag implementation
└── api/
    └── version.go       # GET /version handler
```

**Structure Decision**: Version string is centralized in `internal/config/version.go` to serve as the single source of truth for both the CLI and API layers.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | | |
