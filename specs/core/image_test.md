# Image Generation Service Test Specification

## 1. Render Weather Image (Success)
### Given
A valid city name, `width`, and `height` are provided. The Weather Module dependency is explicitly mocked to return a successful weather forecast domain model.
### When
The weather image generation function is called.
### Then
The weather dependency is invoked correctly, and a valid image object matching the specified `width` and `height` dimensions is successfully returned.

## 2. Render Calendar Image (Success)
### Given
A `width` and `height` are provided. The Calendar Module dependency is explicitly mocked to return a chronological list of upcoming events.
### When
The calendar image generation function is called.
### Then
The calendar dependency is invoked correctly, and a valid image object matching the specified `width` and `height` dimensions is successfully returned.

## 3. Render Combined Image (Success)
### Given
A valid city name, `width`, and `height` are provided. Both the Weather Module and Calendar Module dependencies are explicitly mocked to return their respective domain models successfully.
### When
The combined image generation function is called.
### Then
Both upstream module dependencies are invoked correctly, and a single composite image object matching the specified `width` and `height` dimensions is successfully returned.

## 4. Upstream Dependency Failure Handling (Parameterized)
### Given
A parameterized set of simulated exceptions (e.g., `CityNotFoundError`, `AuthenticationError`, `UpstreamAPIError`) is configured to be raised by either the mocked Weather Module or Calendar Module dependencies.
### When
Any of the image generation functions are called.
### Then
The underlying domain exceptions raised by the dependent modules are safely propagated upwards.
