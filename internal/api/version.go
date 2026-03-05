package api

import (
	"encoding/json"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"net/http"
)

// VersionHandler returns the application version in JSON format.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	core.InfoLogger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	versionInfo := core.VersionInfo{
		Version: config.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
		core.ErrorLogger.Printf("Failed to encode version info: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
