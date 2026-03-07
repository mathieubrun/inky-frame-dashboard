package api

import (
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/agenda"
	"inky-frame-dashboard/internal/core/weather"
	"net/http"
	"strings"
)

// DashboardImageHandler returns a combined weather and agenda image.
func (s *Server) DashboardImageHandler(w http.ResponseWriter, r *http.Request) {
	core.InfoLogger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		location = "Zurich" // Default
	}

	calendarID := r.URL.Query().Get("calendar_id")
	if calendarID == "" {
		calendarID = s.Config.AgendaID
	}

	palette := r.URL.Query().Get("palette")
	if palette == "" {
		palette = "spectra6"
	}

	mockParam := r.URL.Query().Get("mock")
	useMock := s.Config.WeatherMock || s.Config.AgendaMock || strings.ToLower(mockParam) == "true"

	// --- 1. Fetch Weather ---
	var wProvider weather.Provider
	if useMock {
		wProvider = weather.NewMockProvider()
	} else {
		wProvider = weather.NewOpenMeteoProvider()
	}
	wProvider = weather.NewCachedProvider(wProvider, s.Config.WeatherCacheDir, s.Config.WeatherCacheTTL)

	wForecast, err := wProvider.GetForecast(location)
	if err != nil {
		core.ErrorLogger.Printf("Dashboard: Failed to fetch weather: %v", err)
		// We proceed with empty weather if possible, or fail
	}

	// --- 2. Fetch Agenda ---
	var aProvider agenda.CalendarProvider
	if useMock {
		aProvider = agenda.NewMockCalendarProvider()
	} else {
		aProvider = agenda.NewGoogleCalendarProvider(s.Config.GoogleCredentials)
	}
	aProvider = agenda.NewCachedProvider(aProvider, s.Config.AgendaCacheDir, s.Config.AgendaCacheTTL)

	aForecast, err := aProvider.GetAgenda(calendarID, 8)
	if err != nil {
		core.ErrorLogger.Printf("Dashboard: Failed to fetch agenda: %v", err)
		// We proceed with empty agenda if possible
		aForecast = &agenda.AgendaForecast{Events: []agenda.AgendaEvent{}}
	}

	// --- 3. Render Combined Image ---
	data, err := s.DashboardRenderer.Render(wForecast, aForecast, palette)
	if err != nil {
		core.ErrorLogger.Printf("Dashboard: Failed to render image: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// --- 4. ETag / 304 Logic ---
	etag := "\"" + core.CalculateMD5(data) + "\""
	w.Header().Set("ETag", etag)

	if r.Header.Get("If-None-Match") == etag {
		core.InfoLogger.Printf("Dashboard: Image not modified (ETag: %s)", etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(data)
}
