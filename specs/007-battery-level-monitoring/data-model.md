# Data Model: Battery Monitoring

## Entities

### BatteryReport
Represents a single battery level measurement.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Timestamp | Time (RFC3339) | When the report was received by the server. | Non-empty |
| Voltage | float64 | Battery voltage in Volts (e.g., 3.75). | MUST be > 0.0 |

## Storage Strategy

### Local CSV File
- **Path**: `.inky/battery.csv`
- **Format**: `Timestamp,Voltage`
- **Concurrency**: `sync.Mutex` protected writes in `internal/core/battery`.
- **Initialization**: Create file with header `Timestamp,Voltage` if not exists.

## State Transitions
1. **Report Received**: `BatteryReport` created with current server time.
2. **Persistence**: `BatteryReport` appended to CSV file.
3. **Retrieval**: CSV file read as raw text for history endpoint.
