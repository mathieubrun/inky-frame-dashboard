---
name: implement
description: Execute the coding and testing implementation of a feature based on its specification.
---

# Skill: implement

This skill guides the agent through the rigorous process of implementing a feature. It ensures that the newly written code perfectly aligns with the project's architecture, implementation standards, and the feature's dedicated specification documents.

---

## 1. Context & Verification

Before writing any code, the agent must use the `view_file` tool to read the following documents to establish the correct execution context:
1. `specs/ARCHITECTURE.md`
2. `specs/IMPLEMENTATION.md`
3. The specific feature specification (e.g., `specs/core/feature.md` or `specs/api/feature_api.md`)
4. The corresponding testing sidecar specification (e.g., `specs/core/feature_test.md`)

Do not proceed with coding until these documents have been fully read and understood. If the required feature specification is missing, instruct the user to use the `/specify` skill first.

---

## 2. Step-by-Step Execution Guide

When the `implement` skill is invoked, execute the following steps:

1. Analyze Specifications: Ensure the feature's functional requirements and testing strategy are thoroughly understood from the read specification files.
2. Implement the Business Logic & API: Write the actual Python implementation in the appropriate `src/` directories. You must ensure strict compliance with `specs/IMPLEMENTATION.md` (e.g., encapsulating business logic in classes, using Pydantic V2 exclusively, employing dependency injection via Dishka, forbidding docstrings, and enforcing explicit type annotations).
3. Implement the Test Suite: Write the test cases in the `tests/` directory exactly as prescribed by the testing sidecar specification. Strictly follow the `Given`, `When`, `Then` comment pattern, implement the `sut` fixture, and correctly mock out all external dependencies.
4. Validate Tooling: You MUST always run the project's formatting and linting tools (`uv run ruff format` and `uv run ruff check --fix`) using the `run_command` tool after making code changes. Fix any remaining linting errors before proceeding.
5. Execute the Test Suite: You MUST always run the automated tests (`uv run pytest`) using the `run_command` tool after implementation to verify your work. Iterate and debug until all tests pass successfully.
6. Finalize: Present a summary of the implemented feature and the passing test results to the user.
