# Configuration Module Specification

### 1. Purpose
The global configuration and environment variables for the application are managed by this module.

### 2. Functional Requirements
* A centralized `Settings` object that parses environment variables must be provided by the module.
* The persistence data directory (`data_dir`) must be automatically created upon startup if it does not exist.
* The `Settings` object must include a `data_dir` field of type `Path`, which is not strictly required and defaults to `.inky`, representing the path to the application data directory.
