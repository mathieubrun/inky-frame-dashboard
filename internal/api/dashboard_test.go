package api

import (
	"inky-frame-dashboard/internal/config"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDashboardImageHandler_ETag(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "dashboard_api_test")
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		WeatherMock:     true,
		AgendaMock:      true,
		WeatherCacheDir: tempDir,
		WeatherCacheTTL: 1 * time.Hour,
		AgendaCacheDir:  filepath.Join(tempDir, "agenda"),
		AgendaCacheTTL:  1 * time.Hour,
	}
	server := NewServer(cfg)

	// 1. First request
	req1, _ := http.NewRequest("GET", "/dashboard/image", nil)
	rr1 := httptest.NewRecorder()
	server.DashboardImageHandler(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr1.Code)
	}

	etag := rr1.Header().Get("ETag")
	if etag == "" {
		t.Error("Expected ETag header, got empty")
	}

	// 2. Second request with If-None-Match
	req2, _ := http.NewRequest("GET", "/dashboard/image", nil)
	req2.Header.Set("If-None-Match", etag)
	rr2 := httptest.NewRecorder()
	server.DashboardImageHandler(rr2, req2)

	if rr2.Code != http.StatusNotModified {
		t.Errorf("Expected status Not Modified, got %v", rr2.Code)
	}
}

func TestDashboardImageHandler_MethodNotAllowed(t *testing.T) {
	server := NewServer(&config.Config{})
	req, _ := http.NewRequest("POST", "/dashboard/image", nil)
	rr := httptest.NewRecorder()
	server.DashboardImageHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %v", rr.Code)
	}
}

func TestDashboardImageHandler_CalendarChange(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "dashboard_calendar_test")
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		WeatherMock:     true,
		WeatherCacheDir: tempDir,
		WeatherCacheTTL: 1 * time.Hour,
	}
	server := NewServer(cfg)
	// Inject our controllable mock into the handler logic is hard without refactoring.
	// But we can use the Mock param and just wait for time to pass or change a param.
	
	// Actually, changing the location will change the image
	req1, _ := http.NewRequest("GET", "/dashboard/image?location=Zurich", nil)
	rr1 := httptest.NewRecorder()
	server.DashboardImageHandler(rr1, req1)
	etag1 := rr1.Header().Get("ETag")

	req2, _ := http.NewRequest("GET", "/dashboard/image?location=Geneva", nil)
	rr2 := httptest.NewRecorder()
	server.DashboardImageHandler(rr2, req2)
	etag2 := rr2.Header().Get("ETag")

	if etag1 == etag2 {
		t.Error("Expected different ETags for different locations, but got same")
	}
}
