# Calendar Service Test Specification

## 1. Retrieve Upcoming Events (Success)
### Given
The Google Calendar API dependency is explicitly mocked to return a successful, well-formed JSON response containing a list of future events. The Configuration module is mocked to provide necessary API credentials and base URLs.
### When
The upcoming events retrieval function is called.
### Then
The upstream API method is executed with correct parameters (e.g., filtering for upcoming dates), and a structured domain model containing the sorted chronological list of events is successfully returned.

## 2. In-Memory Caching (Success)
### Given
The Configuration module is mocked to define a specific cache duration. The Google Calendar API dependency is explicitly mocked to return a successful JSON response containing a list of events.
### When
The upcoming events retrieval function is called multiple times sequentially within the configured cache duration.
### Then
The underlying Google Calendar API dependency is invoked exactly once. The subsequent function calls successfully return the cached parsed domain model.

## 3. Cache Expiration
### Given
The Configuration module is mocked to define a specific cache duration. The Google Calendar API dependency is explicitly mocked to return a successful JSON response containing a list of events.
### When
The upcoming events retrieval function is called, the system time or cache state is advanced beyond the configured cache duration, and the function is called again.
### Then
The underlying Google Calendar API dependency is invoked twice, and a fresh parsed domain model is returned by the second call.

## 4. Empty Upcoming Events Handling
### Given
The Google Calendar API dependency is explicitly mocked to return a successful JSON response containing no events.
### When
The upcoming events retrieval function is called.
### Then
An empty list is successfully returned without raising any exceptions.

## 5. Upstream API Failure Handling (Parameterized)
### Given
A parameterized set of simulated failures from the mocked Google Calendar API dependency is provided (e.g., `401 Unauthorized` for authentication failure, `500 Internal Server Error` for upstream outage, or network timeouts).
### When
The upcoming events retrieval function is called.
### Then
An explicit domain exception indicating a third-party API or authentication failure is safely raised and propagated, and the error is explicitly not cached.
