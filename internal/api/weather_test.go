package api

import (
	"inky-frame-dashboard/internal/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestWeatherImageHandler(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "api_weather_image_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	cfg := &config.Config{
		WeatherCacheDir:      tempDir,
		WeatherCacheTTL:      1 * time.Hour,
		WeatherImageCacheDir: tempDir,
		WeatherImageCacheTTL: 1 * time.Hour,
		WeatherMock:          true,
	}

	s := NewServer(cfg)

	req, err := http.NewRequest("GET", "/weather/image?location=Zurich&mock=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.WeatherImageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "image/png" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "image/png")
	}

	if len(rr.Body.Bytes()) == 0 {
		t.Fatal("handler returned empty body")
	}
}

func TestWeatherImageHandler_MissingLocation(t *testing.T) {
	cfg := &config.Config{}
	s := NewServer(cfg)

	req, err := http.NewRequest("GET", "/weather/image", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.WeatherImageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
