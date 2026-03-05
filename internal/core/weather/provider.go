package weather

// Provider defines the interface for fetching weather data.
type Provider interface {
	GetForecast(city string) (*WeatherForecast, error)
}
