# Architecture Principles

The high-level, language-agnostic architecture principles for the project are outlined by this document. These principles are enforced so that the codebase is kept maintainable, testable, and traceable as the system evolves.

---

## 1. High Level Overview

The system is structured using a strict layered architecture so that a unidirectional dependency flow is ensured.

* Remote Client (Inky Frame): The physical hardware client. Waking up, reporting any essential physical metrics (like battery voltage), fetching the pre-compiled dashboard image via HTTP, and updating the e-ink screen before returning to deep sleep are its sole responsibilities. No presentation or business logic is contained within it.
* API Layer: The presentation boundary of the server. The HTTP endpoints are defined, requests from the Remote Client or other consumers are captured, schemas are validated, and operations are delegated to the business logic layer by this layer. No core business logic itself is contained within it.
* Business Logic Layer: The core of the application. All core domain entities, state computations, application rules, and workflows are encapsulated by this layer. It is completely decoupled from the physical screen hardware, specific API frameworks, and databases.
* Infrastructure Layer: The infrastructure boundary. Reading from and writing to local storage, external databases, or third-party APIs are handled by this layer. Communication with the Business Logic Layer is performed exclusively through standard interfaces, so that the core domain is kept completely unaware of how data is physically persisted.

---

## 2. Common Concerns

Cross-cutting concerns must be solved consistently across all modules so that stability, observability, and ease of configuration are ensured.

* Error Handling: System failures must be caught at boundary layers and converted into standard representations. Underlying domain or persistence errors should be gracefully translated into appropriate HTTP responses by the API layer, so that the system execution loop is never unexpectedly crashed.
* Logging: A central, language-agnostic logging abstraction must be used for all execution tracing, warnings, and error reports. Contextual metadata (such as transaction IDs or component names) should be included in log entries and clear severity levels (`DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`) must be enforced by them.
* Configuration: All operational variables, environment-specific credentials, and feature flags must be loaded into a central, read-only configuration model at startup. Runtime side effects are prevented and development, staging, and production environments are cleanly separated by this single source of truth.

---

## 3. Testing Strategy

High software quality is maintained through a multi-tiered, automation-first testing pyramid.

* End-to-End (E2E) Verification (Top of Pyramid): The complete user flow and the final rendering accuracy of the dashboard images are verified by these tests. The actual hardware client requests are simulated by them so that the correct behavior of the system as a whole is ensured.
* Integration Tests (Middle of Pyramid): The correct interaction between architectural boundaries is focused upon by these tests. The validation of API request schemas and the testing of the real Persistence Layer against local files or mock databases are included.
* Unit Tests (Base of Pyramid): High-coverage testing of the Business Logic Layer is focused upon by these tests. These tests are fast, strictly stateless, and any network, file system, or infrastructure dependencies are entirely mocked out by them.
