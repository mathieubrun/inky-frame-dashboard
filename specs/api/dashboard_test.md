# Dashboard API Test Specification

## 1. Retrieve Weather Image (Success)
### Given
A valid city name, width, and height query parameters are provided. The Core Image Module dependency is mocked to return a valid binary image stream (e.g., a synthetic Pillow `Image`).
### When
A `GET` request is made to `/api/dashboard/weather`.
### Then
The endpoint responds with a `200 OK` status and the `Content-Type: image/png` HTTP header.

## 2. Retrieve Calendar Image (Success)
### Given
Valid width and height query parameters are provided. The Core Image Module dependency is mocked to return a valid binary image stream.
### When
A `GET` request is made to `/api/dashboard/calendar`.
### Then
The endpoint responds with a `200 OK` status and the `Content-Type: image/png` HTTP header.

## 3. Retrieve Combined Image (Success)
### Given
A valid city name, width, and height query parameters are provided. The Core Image Module dependency is mocked to return a valid binary image stream.
### When
A `GET` request is made to `/api/dashboard/combined`.
### Then
The endpoint responds with a `200 OK` status and the `Content-Type: image/png` HTTP header.

## 4. Invalid Payload/Method Handling (Parameterized)
### Given
A parameterized set of invalid requests (e.g., an unsupported `POST` request, missing `city_name`, or negative `width`/`height` dimensions) is provided.
### When
The request is made to the respective endpoint.
### Then
A `405 Method Not Allowed` or `422 Unprocessable Entity` is correctly returned.

## 5. Domain Exception Mapping (Parameterized)
### Given
The Core Image Module dependency is mocked to raise specific domain exceptions (e.g., `CityNotFoundError`, `AuthenticationError`, `UpstreamAPIError`).
### When
A `GET` request is made to `/api/dashboard/combined`.
### Then
The API exception handlers correctly intercept and map the internal domain errors to appropriate HTTP responses (e.g., `400 Bad Request`).
