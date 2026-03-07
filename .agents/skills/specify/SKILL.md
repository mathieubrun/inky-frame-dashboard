---
name: specify
description: Interactively help the user author a new, fully compliant feature specification and its testing sidecar.
---

# Skill: specify

This skill provides a guided, interactive workflow to help the user create comprehensive specification documentation for a new feature. It ensures that the generated specification strictly adheres to the formatting, philosophy, and testing requirements defined in `specs/SPECIFICATIONS.md`.

---

## 1. Context & Verification

Before drafting any specifications, the agent must use the `view_file` tool to read the `specs/SPECIFICATIONS.md` document in its entirety to understand the required formatting, philosophy, templates, and testing rules.

---

## 2. Step-by-Step Execution Guide

When the `specify` skill is invoked, execute the following steps:

1. Requirements Gathering & Clarification: Ask the user what new feature they want to specify. If the initial description is vague or misses key architectural details, dependencies, edge cases, or domain boundaries, you **must explicitly ask clarification questions** to gather all necessary information before writing the spec.
2. Determine the Specification Type: Based on the gathered requirements, decide whether the feature requires a standard functional specification or an API specification. Determine its correct directory location within `specs/` (e.g., `specs/core/` or `specs/api/`).
3. Draft the Functional/API Specification: Write the draft for the main specification file (e.g., `feature.md`), strictly following the exact section headers required by `specs/SPECIFICATIONS.md` (e.g., Purpose, Functional Requirements/Data Models, Dependencies).
4. Draft the Testing Sidecar: Write the corresponding draft for the testing sidecar (e.g., `feature_test.md`). Ensure it explicitly follows the "One Section Per Test" rule and includes explicit `Given`, `When`, and `Then` sub-sections, along with parameterization, invalid payload handling, and dependency mocking details where applicable.
5. Review and Iterate: Present the drafted markdown specifications to the user (e.g., as markdown artifacts) for review. Request their feedback and make any necessary adjustments.
6. Finalize: Once the user explicitly approves the drafts, write both the specification and its testing sidecar to their correct locations inside the `specs/` directory.
