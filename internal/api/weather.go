package api

import (
	"encoding/json"
	"fmt"
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

// WeatherImageHandler returns a pre-rendered weather image for a location.
func (s *Server) WeatherImageHandler(w http.ResponseWriter, r *http.Request) {
	core.InfoLogger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing location parameter"})
		return
	}

	// Parse dimensions
	width := 800
	height := 480
	if wStr := r.URL.Query().Get("width"); wStr != "" {
		fmt.Sscanf(wStr, "%d", &width)
	}
	if hStr := r.URL.Query().Get("height"); hStr != "" {
		fmt.Sscanf(hStr, "%d", &height)
	}
	palette := r.URL.Query().Get("palette")
	if palette == "" {
		palette = "spectra6"
	}

	cacheKey := s.ImageCache.GenerateKey(location, width, height, palette)

	// Try image cache first
	if data, err := s.ImageCache.GetImage(cacheKey); err == nil {
		core.InfoLogger.Printf("Image cache hit for %s", cacheKey)
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(data)
		return
	}

	// Cache miss, need weather data
	mockParam := r.URL.Query().Get("mock")
	useMock := s.Config.WeatherMock || strings.ToLower(mockParam) == "true"

	var provider weather.Provider
	if useMock {
		provider = weather.NewMockProvider()
	} else {
		provider = weather.NewOpenMeteoProvider()
	}
	provider = weather.NewCachedProvider(provider, s.Config.WeatherCacheDir, s.Config.WeatherCacheTTL)

	forecast, err := provider.GetForecast(location)
	if err != nil {
		core.ErrorLogger.Printf("Failed to retrieve weather for %s for image generation: %v", location, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "weather service unavailable"})
		return
	}

	// Render image
	req := &weather.ImageRequest{
		Location: location,
		Width:    width,
		Height:   height,
		Palette:  palette,
	}

	data, err := s.WeatherRenderer.Render(forecast, req)
	if err != nil {
		core.ErrorLogger.Printf("Failed to render weather image for %s: %v", location, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate image"})
		return
	}

	// Save to cache
	if err := s.ImageCache.SaveImage(cacheKey, data); err != nil {
		core.ErrorLogger.Printf("Failed to save image to cache for %s: %v", cacheKey, err)
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(data)
}
