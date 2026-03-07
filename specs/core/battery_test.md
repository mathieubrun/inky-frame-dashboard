# Battery Service Test Specification

## 1. Percentage Calculation (Parameterized)
### Given
A parameterized set of voltage inputs and expected percentages is provided, specifically including edge cases:
  * `2.5V` (Below minimum) -> Expect `0.0%`
  * `3.0V` (Exact minimum) -> Expect `0.0%`
  * `3.75V` (Midpoint) -> Expect `50.0%`
  * `4.5V` (Exact maximum) -> Expect `100.0%`
  * `5.0V` (Above maximum) -> Expect `100.0%`
### When
`calculate_percentage(voltage)` is called.
### Then
The expected percentage float bounded between 0.0 and 100.0 is returned.

## 2. Add New Battery Report
### Given
A `BatteryReportRequest` and a mocked `BatteryPersistence` layer are provided, with the Configuration module mocked to provide a temporary `data_dir`.
### When
`add_report(request)` is called.
### Then
`calculate_percentage` is executed, the `save_report` method on the persistence layer is called with the calculated values, and the `BatteryReport` is returned.

## 3. Retrieve Latest Status (Parameterized)
### Given
A parameterized set of returns from the mocked `BatteryPersistence.get_latest_report()` is provided, and the Configuration module is mocked:
  * Returns `None` (empty history).
  * Returns a report with percentage `19.9%` (low battery threshold).
  * Returns a report with percentage `20.0%` (not low battery).
### When
`get_latest_status()` is called.
### Then
* `None` is returned if persistence returns `None`.
* `BatteryStatusResponse` with `is_low=True` is returned if percentage < 20.0.
* `BatteryStatusResponse` with `is_low=False` is returned if percentage >= 20.0.

## 4. Retrieve History (Parameterized)
### Given
A mocked `BatteryPersistence` layer containing historical report data is provided, with the Configuration module mocked, alongside parameterized pagination inputs (e.g., various combinations of `limit` and `offset` bounds, including edge cases like `limit=0` or `offset` exceeding history size).
### When
`get_history(limit, offset)` is called.
### Then
The pagination parameters are successfully delegated to the persistence layer, and the corresponding paginable array of `BatteryReport` entries is returned by the service.
