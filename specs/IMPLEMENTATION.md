# Python Implementation & Best Practices

This document defines the Python-specific implementation guidelines, library ecosystems, and testing patterns to maintain a high-quality and consistent codebase.

---

## 1. Libraries

Use a modern Python toolchain as the standard.

* uv: Use as the primary dependency manager, environment manager, and workflow tool. Execute all scripts, tests, and formatting commands via `uv run`.
* Ruff: Use as the code linter and auto-formatter. Enforce strict linting rules and format consistency across the entire codebase using `uv run ruff check` and `uv run ruff format`.
* pyright: Use as the static type checker to verify and enforce strict type safety across the entire codebase.
* FastAPI: Use as the core web framework for all HTTP endpoints. It provides automatic OpenAPI generation and declarative routing.
* Pydantic: Use Pydantic v2 for all data validation, serialization, and settings management. Use modern V2 API methods (e.g., `model_dump()`, `model_validate()`) exclusively. Represent every data input and output across the API and core domain logic as a Pydantic model. Use `pydantic-settings` to load environment-specific variables into a central configuration model.
* pytest: Use as the primary testing framework. All test execution must be driven by `pytest`.
* dishka: Use as the dependency injection framework to manage the lifecycle and resolution of dependencies across the core and API layers.
* cachetools: Use for robust, in-memory TTL and LRU caching instead of building custom caching logic.

---

## 2. Directory Structure

Organize the codebase strictly into the following top-level layers to enforce architectural boundaries:

* `api/`: Contains all presentation layer components, including FastAPI routers, HTTP request/response Pydantic models, and API-specific exception handlers. No business logic may reside here.
* `core/`: Contains all business logic, domain entities, rules, and workflows. This layer must remain completely decoupled from the web framework and physical infrastructure.

---

## 3. Code Standards

Maintain a strict and organized architecture using standard Python features.

* Classes: Encapsulate all business logic within classes. Never expose standalone functional APIs from core modules. The presentation layer must instantiate these classes and invoke their methods.
* Separation of Concerns: Should a feature require multiple dependent classes (e.g., persistence, database access, external API access) to facilitate separation of concerns and testing, they must each be separated into their respective files, inside a folder having the same name as the feature. Their associated tests must also be separated into corresponding test files.
* Type Annotations: Python type annotations (type hints) are strictly mandatory. Explicitly annotate every function argument, class attribute, return type, and test fixture.
* Dependency Injection: Use `dishka` for managing dependency injection. Use constructor injection (`__init__`) for injecting dependencies into business logic classes. In the API layer, use `dishka`'s integration to inject dependencies into FastAPI route handlers instead of `Depends()`.
* Docstrings: The generation of docstrings is strictly forbidden across the entire codebase to avoid redundant noise. Code should be self-documenting through clear variable, class, and method names.
* External APIs: Pydantic models must be created for handling responses and requests of external APIs. Use the official Python module for the external API if it exists.


---

## 4. Cross-Cutting Concerns

Implement common architectural concerns consistently across the Python codebase.

* Error Handling: Catch domain-specific exceptions at the API boundary. Create a base `DomainError` exception class that all core domain exceptions must inherit from. Exceptions should be declared in their respective feature modules, except for those that can be (and are) reused across multiple features. Use FastAPI exception handlers (`@app.exception_handler`) to map these internal Python domain exceptions cleanly to standard HTTP error responses.
* Logging: Use the standard Python `logging` module. Configure a central logger that enforces clear severity levels (`DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`) and includes contextual metadata (like component names).
* Configuration: Load all configurations at startup into a read-only Pydantic settings model.

---

## 5. Testing Standards

A robust test suite is critical for preventing regressions. Write tests that are clean, structured, and readable.

* Test File Naming: Mirror the naming of the specification components being verified in the test files.
* Test Classes: Encapsulate all tests within test classes.
* System Under Test (`sut`): Always name the main component being tested `sut`. Instantiate and inject it into test methods via a `@pytest.fixture`.
* Fixture Mock Injection: When the `sut` relies on external dependencies, mock those dependencies, create them as fixtures, and cleanly inject them into the `sut` fixture.
* Explicit Given/When/Then Pattern: Strictly follow the Given, When, Then (GWT) pattern in every test method. Use explicit `# Given`, `# When`, and `# Then` comments to separate setup, execution, and validation phases.

---

## 6. FastAPI Standards

Follow consistent organizational and documentation rules in the API layer.

* One Router Per Feature: Organize API endpoints into clean, domain-specific routers. Implement each distinct feature specification in its own separate FastAPI router module.
* OpenAPI Spec Present: Keep FastAPI's automatic OpenAPI specification generation active. Correctly document endpoint contracts (including status codes, request models, and response models).

---

## 7. Bruno API Collections

Maintain a comprehensive and up-to-date suite of local API requests for reproducible testing.

* Bruno API Collections: Create Bruno request files (`.bru`) so that all HTTP endpoints can be reproducibly tested locally.
* Use a folder structure within the `bruno/` directory that strictly mirrors the API endpoint paths.
* All request URLs must use the `{{base_url}}` environment parameter instead of hardcoded hostnames.
* Bruno API collections must be completely aligned with the defined API specifications and their respective parameters.
