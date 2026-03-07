# Specification Writing Guidelines

The strict standards for writing and maintaining technical specification files are described by this document.

---

## 1. Core Philosophy

The following principles must be adhered to by all specifications:

* Language Agnostic: The architecture, requirements, and contracts should be described by specifications without being tightly coupled to a specific programming language. (Exception: Implementation guidelines, such as `IMPLEMENTATION.md`, are explicitly exempt from this rule).
* Voice: The use of the first person is strictly prohibited.
* Concise and Straight to the Point: Filler words, introductions, or essays should be avoided. The requirements must be stated directly.
* Limited Formatting: Markdown formatting (bold, italics, etc.) should be limited only to emphasize important or key points. Visual clutter must be avoided.
* Avoid Superlatives: Subjective superlatives (e.g., "the best", "fastest", "amazing", and so on) must not be used unless an explicit, measurable requirement is described by them. Adjectives are allowed (e.g, "clean", "structured", and so on)
* Testing Sidecar: Every functional specification template must be accompanied by a sidecar `_test` testing specification template in the same directory (e.g., a corresponding `feature_test.md` must be provided for `feature.md`).
* Explicit Dependencies: All dependencies on other modules or features required by a module or feature must be explicitly listed so that architectural coupling is made transparent.

---

## 2. Spec Template

This exact section structure must be followed by standard functional specification files:

### 1. Purpose
A single, clear sentence explaining what is accomplished by the module or feature.

### 2. Functional Requirements
A concise, bulleted list of strict requirements and rules using unambiguous language.

### 3. Dependencies (if applicable)
A bulleted list of markdown links of all modules or features required by the specified feature.


---

## 3. API Spec Template

For specifications dedicated specifically to API endpoints and routers, this exact structure must be followed by the file:

### 1. Purpose
A single sentence by which the domain and goal of the API boundary are described.

### 2. Data Models
All request, response, and internal data structures must be documented using a Markdown table with exactly the following columns. Each data model must be specified in its own separate table, and the word "Request" or "Response" must be included in the name of each data model:
| Field Name | Type | Required | Default Value | Description |
|---|---|---|---|---|

### 3. Endpoints
All API endpoints must be documented using a Markdown table with exactly the following columns. A request model is mandatory if the request takes a payload or parameters:
| Endpoint Name | Verb | URL with query parameters | Request Model | Response Type | Response Model |
|---|---|---|---|---|---|

### 4. Dependencies (if applicable)
A bulleted list of markdown links of all modules or features required by the API boundary

---

## 4. Testing Specification Template

This strict structure must be followed by all `_test.md` sidecar files:

* One Section Per Test: Its own heading/section must be assigned to each distinct test case. For API specifications, each individual API endpoint must be documented in its own dedicated test section.
* Given / When / Then: Explicitly separated `Given`, `When`, and `Then` subsections by which the state, execution, and expected outcome are described must be included in each test section.
* Parameterization & Edge Cases: Tests should be parameterized whenever possible so that multiple code paths can be efficiently covered. Boundary values and edge cases must be explicitly included in parameterized data sets so that system resilience is ensured.
* Invalid Payloads: Test cases for handling invalid payloads and malformed requests must be explicitly included in every API endpoint test section.
* Dependency Mocking: All explicitly listed dependencies from the associated functional specification must be addressed, and how they are mocked or substituted in the test environment must be detailed.
