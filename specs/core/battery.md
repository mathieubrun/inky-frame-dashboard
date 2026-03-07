# Battery Service Specification

### 1. Purpose
The business logic for processing battery reports, calculating statuses, and integrating with the persistence layer is managed by the battery service module.

### 2. Functional Requirements
* Report Insertion: A new battery report containing the battery voltage must be accepted by the battery service, the corresponding battery percentage and timestamp must be automatically computed, and the persistence layer must be called so that the report is saved in a CSV file.
* Percentage Formula: The battery percentage must be calculated based on 3 x 1.5V AA batteries, using the voltage range 3.0V (0%) to 4.5V (100%): `percentage = (voltage - 3.0) / (4.5 - 3.0) * 100`, bounded strictly between `0` and `100`.
* Latest Status Query: The most recent battery report must be retrieved from the persistence layer by the service, and a determination of whether the battery is low (defined as a battery percentage strictly below `20%`) must be made.
* History Retrieval: Historical data must be retrieved and returned as a paginable array of battery report entries by the battery service.

### 3. Dependencies
* [Configuration Module](config.md)
