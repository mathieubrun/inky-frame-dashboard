package api

import (
	"fmt"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/battery"
	"inky-frame-dashboard/internal/core/dashboard"
	"inky-frame-dashboard/internal/core/weather"
	"net/http"
)

// Server represents the HTTP server.
type Server struct {
	Config            *config.Config
	WeatherRenderer   *weather.WeatherImageRenderer
	DashboardRenderer *dashboard.DashboardRenderer
	ImageCache        *weather.WeatherImageCache
	BatteryProcessor  *battery.Processor
}

// NewServer creates a new server instance.
func NewServer(cfg *config.Config) *Server {
	batteryStorage := battery.NewStorage(cfg.BatteryCSVPath)
	return &Server{
		Config:            cfg,
		WeatherRenderer:   weather.NewWeatherImageRenderer(cfg.FontPath),
		DashboardRenderer: dashboard.NewDashboardRenderer(cfg.FontPath),
		ImageCache:        weather.NewWeatherImageCache(cfg.WeatherImageCacheDir, cfg.WeatherImageCacheTTL),
		BatteryProcessor:  battery.NewProcessor(batteryStorage),
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/version", VersionHandler)
	mux.HandleFunc("/weather/swiss", s.WeatherHandler)
	mux.HandleFunc("/weather/image", s.WeatherImageHandler)
	mux.HandleFunc("/agenda", s.AgendaHandler)
	mux.HandleFunc("/dashboard/image", s.DashboardImageHandler)

	// Battery routes
	mux.HandleFunc("/api/v1/battery", s.BatteryReportHandler)
	mux.HandleFunc("/api/v1/battery/status", s.BatteryStatusHandler)
	mux.HandleFunc("/api/v1/battery/history", s.BatteryHistoryHandler)

	core.InfoLogger.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, mux)
}
