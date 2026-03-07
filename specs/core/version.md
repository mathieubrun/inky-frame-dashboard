# Version Module Specification

### 1. Purpose
The current application version string is retrieved and returned by this module.

### 2. Functional Requirements
* A mechanism to fetch the application version must be exposed by the module.
* The application version must be defined exactly once as a single source of truth (e.g., extracted from `pyproject.toml` or a dedicated application constant).
