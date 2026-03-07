---
name: verify
description: Verify implementation, code completeness, or specification compliance based on user argument.
---

# Skill: verify

This skill acts as a unified verification workflow to ensure the integrity and compliance of the project's codebase and specifications.

---

## 1. Execution Mode Routing

When the `verify` skill is invoked, first check if a specific target mode (argument) was provided. If no argument is provided, or if the argument is unrecognized, you **must explicitly ask the user** to choose one of the following execution modes before proceeding:

- `implementation`: Verifies that the implementation matches `specs/ARCHITECTURE.md`, `specs/IMPLEMENTATION.md`, and the feature-specific specification and testing sidecar.
- `code`: Verifies that the specifications do not miss an implemented code feature.
- `specification`: Verifies that the specifications conform strictly to `specs/SPECIFICATIONS.md`.

Once the target mode is established, proceed exclusively to the corresponding section below.

---

## 2. Mode: implementation

This mode validates that written code aligns with architectural boundaries, tooling standards, and feature requirements.

1. Read Core Specs: Use the `view_file` tool to read `specs/ARCHITECTURE.md` and `specs/IMPLEMENTATION.md` to establish the baseline rules.
2. Locate Target Files: Find the specific feature's code in `src/`, its tests in `tests/`, and its specification files in `specs/`.
3. Audit: Systematically compare the Python code and test suites against the principles in the core specs and the feature's specific requirements. Verify that the tests adequately cover the scenarios defined in the feature's `_test.md` sidecar.
4. Report: Output a report containing:
   * Status: `PASSED`, `WARNING`, or `FAILED`.
   * Violations Found: Specific, bulleted references to code/test violations or missing coverage.
   * Suggested Remediation: Concrete refactoring recommendations.
5. Approval: Request user approval before making any code or test changes.

---

## 3. Mode: code

This mode validates that existing code is completely documented by corresponding specification files.

1. Locate Target Code: Find the target source files inside `src/`.
2. Locate Specs: Determine the expected location of the specification file in the `specs/` directory based on the directory mirroring rule.
3. Audit: Compare the code's behavior, API contracts, and logic against the content of the corresponding specification file. Note any undocumented features, properties, or missing API contracts.
4. Report: Output a report containing:
   * Status: `PASSED`, `WARNING`, or `FAILED`.
   * Violations Found: Undocumented code features or discrepancies.
   * Suggested Remediation: Recommendations for creating or updating the spec.
5. Approval: Request user approval before modifying specification files.

---

## 4. Mode: specification

This mode validates that specification files are well-written, concise, and properly formatted.

1. Read Core Specs: Use the `view_file` tool to read `specs/SPECIFICATIONS.md`.
2. Locate Target Specs: Find the target `.md` specification files inside the `specs/` directory.
3. Audit: Compare the target specification file's content against the detailed guidelines in `specs/SPECIFICATIONS.md`. Check for first-person voice, excessive formatting, superlatives, explicit dependencies, and the presence of a testing sidecar.
4. Report: Output a report containing:
   * Status: `PASSED`, `WARNING`, or `FAILED`.
   * Violations Found: Specific formatting or structural violations.
   * Suggested Remediation: Text rewrite suggestions or direct markdown diffs.
5. Approval: Request user approval before altering the specification files.
