package weather

import (
	"math/rand"
	"time"
)

// MockProvider implements the Provider interface for testing and development.
type MockProvider struct{}

// NewMockProvider creates a new MockProvider.
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// GetForecast returns randomized weather data for any city.
func (p *MockProvider) GetForecast(city string) (*WeatherForecast, error) {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	now := time.Now().Truncate(time.Hour)

	location := Location{
		City:      city,
		Latitude:  47.3769,
		Longitude: 8.5417,
		Country:   "CH",
	}

	current := WeatherRecord{
		Timestamp:         now,
		Temperature:       10.0 + r.Float64()*10.0,
		WindSpeed:         5.0 + r.Float64()*15.0,
		WindDirection:     r.Float64() * 360.0,
		Precipitation:     r.Float64() * 5.0,
		PrecipitationProb: r.Float64() * 100.0,
	}

	hourly := make([]WeatherRecord, 24)
	for i := 0; i < 24; i++ {
		ts := now.Add(time.Duration(i) * time.Hour)
		hourly[i] = WeatherRecord{
			Timestamp:         ts,
			Temperature:       current.Temperature + (r.Float64()*4.0 - 2.0),
			WindSpeed:         current.WindSpeed + (r.Float64()*5.0 - 2.5),
			WindDirection:     current.WindDirection + (r.Float64()*20.0 - 10.0),
			Precipitation:     r.Float64() * 2.0,
			PrecipitationProb: r.Float64() * 100.0,
		}
	}

	return &WeatherForecast{
		Location:  location,
		Current:   current,
		Hourly:    hourly,
		FetchedAt: time.Now(),
	}, nil
}

// Ensure MockProvider implements Provider.
var _ Provider = (*MockProvider)(nil)
