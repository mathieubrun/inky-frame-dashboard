# Data Model: Weather Data Image Generation

This document defines the entities and data structures used for generating and serving weather images.

## Entities

### Image Request Parameters
Represents the inputs required to generate a specific weather image.

| Field | Type | Description |
| :--- | :--- | :--- |
| `location` | `string` | City name or postcode (e.g., "Zurich", "8001"). |
| `width` | `int` | Target image width in pixels (Default: 800). |
| `height` | `int` | Target image height in pixels (Default: 480). |
| `palette` | `string` | The target display palette (e.g., "spectra6"). |

### Render Template
Defines the layout and styling of weather data on the image canvas.

| Field | Type | Description |
| :--- | :--- | :--- |
| `Font` | `string` | The font family used for text rendering. |
| `Icons` | `map[string]image.Image` | Mapping of weather conditions to pre-loaded icons. |
| `Background` | `color.Color` | The background color of the canvas. |

### Weather Cache Entry
Represents a cached weather image stored on disk.

| Field | Type | Description |
| :--- | :--- | :--- |
| `Key` | `string` | Unique identifier (e.g., `zurich_800x480_spectra6`). |
| `FilePath` | `string` | Path to the generated PNG file. |
| `CreatedAt` | `time.Time` | Timestamp when the image was generated. |
| `ExpiresAt` | `time.Time` | Timestamp when the cache entry expires. |

## Relationships

- A **Weather Image** is generated from a **Weather Forecast** (internal core entity).
- A **Weather Forecast** is retrieved based on a **Location Identifier** (internal core entity).
- A **Weather Image** is stored in the **Weather Cache** and identified by a unique key derived from request parameters.

## Validation Rules

- **Location**: MUST be a non-empty string.
- **Resolution**: Width and Height MUST be positive integers.
- **Palette**: MUST be one of the supported values (`spectra6`).
