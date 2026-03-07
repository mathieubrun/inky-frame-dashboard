package api

import (
	"encoding/json"
	"inky-frame-dashboard/internal/core"
	"net/http"
)

// BatteryHandler handles battery-related requests.
func (s *Server) BatteryReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Voltage float64 `json:"voltage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := s.BatteryProcessor.AddReport(req.Voltage)
	if err != nil {
		core.ErrorLogger.Printf("Failed to add battery report: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	core.InfoLogger.Printf("Battery report received: %.2fV", report.Voltage)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// BatteryStatusHandler returns the latest battery status.
func (s *Server) BatteryStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := s.BatteryProcessor.GetLatest()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// BatteryHistoryHandler returns the raw CSV battery history.
func (s *Server) BatteryHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := s.BatteryProcessor.GetHistoryRaw()
	if err != nil {
		core.ErrorLogger.Printf("Failed to get battery history: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}
