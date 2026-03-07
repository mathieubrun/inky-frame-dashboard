# Weather Service Test Specification

## 1. Retrieve Forecast by City Name (Success)
### Given
A valid city name is provided, and the Open-Meteo Geocoding API and Forecast API dependencies are explicitly mocked to return successful, well-formed JSON responses. The Configuration module is mocked to provide API base URLs.
### When
The forecast retrieval function is called with the city name.
### Then
The geocoding and forecast methods are executed with correct parameters (including the `icon_ch_seamless` model flag), and a structured domain model containing the weather data is successfully returned.

## 2. In-Memory Caching (Success)
### Given
A valid city name is provided, and the Configuration module is mocked to define a specific cache duration. The Open-Meteo API dependencies are mocked to return successful responses.
### When
The forecast retrieval function is called multiple times sequentially with the same city name within the configured cache duration.
### Then
The underlying geocoding and forecast API dependencies are invoked exactly once. The subsequent function calls successfully return the cached parsed domain model.

## 3. Cache Expiration
### Given
A valid city name is provided, and the Configuration module is mocked to define a specific cache duration. The Open-Meteo API dependencies are mocked to return successful responses.
### When
The forecast retrieval function is called, the system time or cache state is advanced beyond the configured cache duration, and the function is called again.
### Then
The underlying geocoding and forecast API dependencies are invoked twice, and a fresh parsed domain model is returned by the second call.

## 4. Unknown City Name Handling (Parameterized)
### Given
A parameterized set of invalid or unknown city names (e.g., an empty string, a non-existent city name, special characters) is provided. The Geocoding API dependency is mocked to return an empty result list or a `404 Not Found`.
### When
The forecast retrieval function is called.
### Then
An explicit domain exception (e.g., `CityNotFoundError`) is raised, the Forecast API dependency is never called, and the error is explicitly not cached.

## 5. Upstream API Failure Handling
### Given
A valid city name is provided, and the Geocoding API returns valid coordinates, but the Forecast API dependency is mocked to return an HTTP error (e.g., `500 Internal Server Error` or network timeout).
### When
The forecast retrieval function is called.
### Then
An explicit domain exception indicating a third-party API failure is safely raised and propagated, and the error is explicitly not cached.
