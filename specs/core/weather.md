# Weather Service Specification

### 1. Purpose
The business logic for retrieving weather forecasts by city name via the Open-Meteo API (using the MeteoSwiss ICON-CH models) is managed by the weather service module.

### 2. Functional Requirements
* Geocoding: The provided city name must be translated into geographic coordinates (latitude and longitude) using the Open-Meteo Geocoding API.
* Forecast Retrieval: Weather forecast data must be retrieved from the Open-Meteo API using the resolved geographic coordinates.
* Caching: Parsed weather forecast domain models must be temporarily stored in an in-memory cache for a configurable duration managed by the Configuration Module. Only successful responses must be cached.
* Error Handling: Explicit domain exceptions must be raised if the city cannot be located or if the upstream APIs fail to return a valid response.
* Data Transformation: Raw upstream JSON responses must be parsed and returned as a standardized internal domain model representing current and forecasted weather conditions.

### 3. Dependencies
* [Configuration Module](config.md)
