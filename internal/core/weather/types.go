package weather

import "time"

// WeatherRecord represents the weather state at a specific point in time.
type WeatherRecord struct {
	Timestamp         time.Time `json:"timestamp"`
	Temperature       float64   `json:"temperature"`
	Condition         string    `json:"condition"`
	WindSpeed         float64   `json:"wind_speed"`
	WindDirection     float64   `json:"wind_direction"`
	Precipitation     float64   `json:"precipitation"`
	PrecipitationProb float64   `json:"precipitation_prob"`
}

// Location represents a geographic point.
type Location struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

// ImageRequest represents the parameters for generating a weather image.
type ImageRequest struct {
	Location string `json:"location"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Palette  string `json:"palette"`
}

// WeatherForecast represents the current weather and the 24-hour hourly forecast.
type WeatherForecast struct {
	Location  Location        `json:"location"`
	Current   WeatherRecord   `json:"current"`
	Hourly    []WeatherRecord `json:"hourly"`
	FetchedAt time.Time       `json:"fetched_at"`
}
