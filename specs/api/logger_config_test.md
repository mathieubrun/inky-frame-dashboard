# Logger Configuration Test Specification

## 1. Response Header Injection
### Given
A running API application configured with the `RequestIdMiddleware` is provided.
### When
A standard HTTP GET request is executed to any valid endpoint.
### Then
The HTTP response contains the header `X-Request-ID`, and its value is verified to be exactly 8 characters long and consisting only of valid hexadecimal characters.

## 2. Request ID Uniqueness
### Given
A running API application configured with the `RequestIdMiddleware` is provided.
### When
Multiple concurrent or sequential HTTP requests are executed to the application.
### Then
A unique `X-Request-ID` value is contained in every HTTP response, confirming that a new identifier is generated per request.

## 3. Context Log Emission
### Given
An API application with the `RequestIdMiddleware` and `RequestIdFilter` enabled is provided, and a mock logger is attached to capture emitted log records.
### When
A standard HTTP GET request is executed so that a log output is triggered.
### Then
A `request_id` attribute that matches the `X-Request-ID` from the response header is contained in the captured log record.
