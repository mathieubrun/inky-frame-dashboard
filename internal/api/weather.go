package api

import (
	"encoding/json"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/weather"
	"net/http"
	"strings"
)

// WeatherResponse represents the JSON response for the weather API.
type WeatherResponse struct {
	*weather.WeatherForecast
	Source string `json:"source"`
}

// WeatherHandler returns the weather forecast for a Swiss city.
func (s *Server) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	core.InfoLogger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing city parameter"})
		return
	}

	mockParam := r.URL.Query().Get("mock")
	useMock := s.Config.WeatherMock || strings.ToLower(mockParam) == "true"

	var provider weather.Provider
	source := "MeteoSwiss (ICON-CH via Open-Meteo)"
	if useMock {
		provider = weather.NewMockProvider()
		source = "Mock Provider"
	} else {
		provider = weather.NewOpenMeteoProvider()
	}

	// Wrap with cache
	provider = weather.NewCachedProvider(provider, s.Config.WeatherCacheDir, s.Config.WeatherCacheTTL)

	forecast, err := provider.GetForecast(city)
	if err != nil {
		core.ErrorLogger.Printf("Failed to retrieve weather for %s: %v", city, err)
		
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "city not found") {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "city not found: " + city})
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "weather service unavailable"})
		}
		return
	}

	resp := WeatherResponse{
		WeatherForecast: forecast,
		Source:          source,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		core.ErrorLogger.Printf("Failed to encode weather response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
