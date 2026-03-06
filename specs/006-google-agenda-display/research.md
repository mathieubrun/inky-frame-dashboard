# Research: Google Agenda Display

This document outlines the technical decisions and research findings for integrating Google Calendar events into the Inky Frame Dashboard.

## Google Calendar API Integration

**Decision**: Use the official `google.golang.org/api/calendar/v3` Go library.

**Authentication**: 
- Use **Service Account** authentication as requested.
- The JSON credentials file will be loaded via `option.WithCredentialsFile` or `option.WithCredentialsJSON`.
- Users must share their target calendars with the service account's email address.

**Mocking for Testing**:
- Use `httptest.NewServer` to mock the Google API responses.
- Override the API endpoint in tests using `option.WithEndpoint(ts.URL)` and `option.WithoutAuthentication()`.

## Dashboard Layout (Combined Image)

**Decision**: Implement a split-screen rendering strategy in `internal/core/weather/image.go` (or a new rendering service).

**Layout Details**:
- **Canvas Size**: 800x480 pixels.
- **Left Panel (0-399px)**: Weather data (Temperature, Icon, Condition).
- **Right Panel (400-799px)**: Upcoming Agenda events (Up to 8 entries).
- **Separation**: A vertical divider line at x=400.

**Rendering with `gg`**:
- Use `dc.Push()` and `dc.Translate(400, 0)` to render the agenda panel independently of the weather panel.
- Or simply use an `offsetX` variable for all drawing commands in the right panel.

## Data Fetching & Caching

**Decision**: 
- Implement a `CalendarProvider` interface similar to the `WeatherProvider`.
- Use a dedicated cache for agenda events to respect API quotas and ensure fast image generation.
- The `GetDashboardImage` endpoint will trigger parallel fetches for both weather and calendar data.

## Configuration

**Decision**: Use `viper` to manage server-side agenda configuration.
- `agenda_id`: The Google Calendar ID.
- `google_credentials_path`: Path to the service account JSON file.
- `agenda_mock`: Boolean flag to use mock data for testing/development.
