package api

import (
	"encoding/json"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/agenda"
	"net/http"
	"strconv"
	"strings"
)

// AgendaHandler returns upcoming calendar events in JSON format.
func (s *Server) AgendaHandler(w http.ResponseWriter, r *http.Request) {
	core.InfoLogger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	calendarID := r.URL.Query().Get("calendar_id")
	if calendarID == "" {
		calendarID = s.Config.AgendaID
	}

	count := 10
	if cStr := r.URL.Query().Get("count"); cStr != "" {
		if c, err := strconv.Atoi(cStr); err == nil {
			count = c
		}
	}

	mockParam := r.URL.Query().Get("mock")
	useMock := s.Config.AgendaMock || strings.ToLower(mockParam) == "true"

	var provider agenda.CalendarProvider
	if useMock {
		provider = agenda.NewMockCalendarProvider()
	} else {
		provider = agenda.NewGoogleCalendarProvider(s.Config.GoogleCredentials)
	}

	// Wrap with cache
	provider = agenda.NewCachedProvider(provider, s.Config.AgendaCacheDir, s.Config.AgendaCacheTTL)

	forecast, err := provider.GetAgenda(calendarID, count)
	if err != nil {
		core.ErrorLogger.Printf("Failed to retrieve agenda for %s: %v", calendarID, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "agenda service unavailable"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(forecast); err != nil {
		core.ErrorLogger.Printf("Failed to encode agenda response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
