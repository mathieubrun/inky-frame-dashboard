package api

import (
	"fmt"
	"inky-frame-dashboard/internal/core"
	"net/http"
)

// Server represents the HTTP server.
type Server struct {
	Addr string
}

// NewServer creates a new server instance.
func NewServer(port int) *Server {
	return &Server{
		Addr: fmt.Sprintf(":%d", port),
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", VersionHandler)

	core.InfoLogger.Printf("Starting server on %s", s.Addr)
	return http.ListenAndServe(s.Addr, mux)
}
