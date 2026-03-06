package weather

import (
	"encoding/json"
	"fmt"
	"inky-frame-dashboard/internal/core"
	"net/http"
	"net/url"
	"time"
)

// OpenMeteoProvider implements the Provider interface using Open-Meteo APIs.
type OpenMeteoProvider struct {
	httpClient   *http.Client
	geocodingURL string
	forecastURL  string
}

// NewOpenMeteoProvider creates a new OpenMeteoProvider.
func NewOpenMeteoProvider() *OpenMeteoProvider {
	return &OpenMeteoProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		geocodingURL: "https://geocoding-api.open-meteo.com/v1/search",
		forecastURL:  "https://api.open-meteo.com/v1/forecast",
	}
}

// GetForecast resolves the city to coordinates and fetches the weather forecast.
func (p *OpenMeteoProvider) GetForecast(city string) (*WeatherForecast, error) {
	core.InfoLogger.Printf("Fetching weather forecast for city: %s", city)
	location, err := p.geocode(city)
	if err != nil {
		core.ErrorLogger.Printf("Geocoding failed for city %s: %v", city, err)
		return nil, fmt.Errorf("geocoding error: %w", err)
	}

	core.InfoLogger.Printf("Geocoding successful for %s: Lat=%.4f, Lon=%.4f", city, location.Latitude, location.Longitude)

	forecast, err := p.fetchWeather(location)
	if err != nil {
		core.ErrorLogger.Printf("Weather fetch failed for %s: %v", city, err)
		return nil, fmt.Errorf("weather fetch error: %w", err)
	}

	core.InfoLogger.Printf("Successfully fetched weather for %s", city)
	return forecast, nil
}

type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

func (p *OpenMeteoProvider) geocode(city string) (*Location, error) {
	apiURL := fmt.Sprintf("%s?name=%s&count=1&language=en&format=json", p.geocodingURL, url.QueryEscape(city))
	core.InfoLogger.Printf("Calling geocoding API: %s", apiURL)
	
	resp, err := p.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var data geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}

	res := data.Results[0]
	return &Location{
		City:      res.Name,
		Latitude:  res.Latitude,
		Longitude: res.Longitude,
		Country:   res.Country,
	}, nil
}

type weatherResponse struct {
	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2m            []float64 `json:"temperature_2m"`
		Weathercode              []int     `json:"weathercode"`
		Windspeed10m             []float64 `json:"windspeed_10m"`
		Winddirection10m         []float64 `json:"winddirection_10m"`
		Precipitation            []float64 `json:"precipitation"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
	} `json:"hourly"`
}

func weatherCodeToCondition(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1, 2, 3:
		return "Mainly clear, partly cloudy, and overcast"
	case 45, 48:
		return "Fog and depositing rime fog"
	case 51, 53, 55:
		return "Drizzle: Light, moderate, and dense intensity"
	case 56, 57:
		return "Freezing Drizzle: Light and dense intensity"
	case 61, 63, 65:
		return "Rain: Slight, moderate and heavy intensity"
	case 66, 67:
		return "Freezing Rain: Light and heavy intensity"
	case 71, 73, 75:
		return "Snow fall: Slight, moderate, and heavy intensity"
	case 77:
		return "Snow grains"
	case 80, 81, 82:
		return "Rain showers: Slight, moderate, and violent"
	case 85, 86:
		return "Snow showers slight and heavy"
	case 95:
		return "Thunderstorm: Slight or moderate"
	case 96, 99:
		return "Thunderstorm with slight and heavy hail"
	default:
		return "Unknown"
	}
}

func (p *OpenMeteoProvider) fetchWeather(location *Location) (*WeatherForecast, error) {
	// Use best_match models (automatically selects ICON-CH for Switzerland)
	apiURL := fmt.Sprintf(
		"%s?latitude=%.4f&longitude=%.4f&hourly=temperature_2m,weathercode,windspeed_10m,winddirection_10m,precipitation,precipitation_probability&forecast_days=1",
		p.forecastURL, location.Latitude, location.Longitude,
	)
	core.InfoLogger.Printf("Calling weather forecast API: %s", apiURL)

	resp, err := p.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	hourly := make([]WeatherRecord, len(data.Hourly.Time))
	for i := range data.Hourly.Time {
		t, _ := time.Parse("2006-01-02T15:04", data.Hourly.Time[i])
		hourly[i] = WeatherRecord{
			Timestamp:         t,
			Temperature:       data.Hourly.Temperature2m[i],
			Condition:         weatherCodeToCondition(data.Hourly.Weathercode[i]),
			WindSpeed:         data.Hourly.Windspeed10m[i],
			WindDirection:     data.Hourly.Winddirection10m[i],
			Precipitation:     data.Hourly.Precipitation[i],
			PrecipitationProb: data.Hourly.PrecipitationProbability[i],
		}
	}

	// Find the "current" record (closest to now)
	now := time.Now()
	var current WeatherRecord
	if len(hourly) > 0 {
		current = hourly[0]
		minDiff := now.Sub(hourly[0].Timestamp)
		if minDiff < 0 {
			minDiff = -minDiff
		}
		for _, r := range hourly {
			diff := now.Sub(r.Timestamp)
			if diff < 0 {
				diff = -diff
			}
			if diff < minDiff {
				minDiff = diff
				current = r
			}
		}
	}

	return &WeatherForecast{
		Location:  *location,
		Current:   current,
		Hourly:    hourly,
		FetchedAt: now,
	}, nil
}

// Ensure OpenMeteoProvider implements Provider.
var _ Provider = (*OpenMeteoProvider)(nil)
