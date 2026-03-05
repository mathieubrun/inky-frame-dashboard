package api

import (
	"fmt"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"net/http"
)

// Server represents the HTTP server.
type Server struct {
	Config *config.Config
}

// NewServer creates a new server instance.
func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/version", VersionHandler)
	mux.HandleFunc("/weather/swiss", s.WeatherHandler)

	core.InfoLogger.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, mux)
}
