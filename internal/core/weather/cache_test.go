package weather

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCachedProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "weather_cache_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	mock := NewMockProvider()
	ttl := 1 * time.Hour
	cache := NewCachedProvider(mock, tempDir, ttl)

	city := "Zurich"

	// 1. First fetch (should hit mock and save to cache)
	f1, err := cache.GetForecast(city)
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}

	// Verify file exists
	filename := cache.getCacheFilename(city)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("cache file not created: %s", filename)
	}

	// 2. Second fetch (should hit cache)
	f2, err := cache.GetForecast(city)
	if err != nil {
		t.Fatalf("second fetch failed: %v", err)
	}

	// Since mock returns random data, if f1 and f2 are identical, it hit the cache
	if !f1.FetchedAt.Equal(f2.FetchedAt) {
		t.Errorf("expected cached data, but got new data (FetchedAt: %v != %v)", f1.FetchedAt, f2.FetchedAt)
	}

	// 3. Third fetch with mocked "stale" time (after 04:00 AM)
	// If f1 was fetched at T, and we are now at T+24h, it should be stale
	cache.now = func() time.Time {
		return f1.FetchedAt.Add(24 * time.Hour)
	}
	f3, err := cache.GetForecast(city)
	if err != nil {
		t.Fatalf("third fetch failed: %v", err)
	}

	if f3.FetchedAt.Equal(f1.FetchedAt) {
		t.Errorf("expected new data after 24h, but got cached data")
	}
}

func TestGetCacheFilename(t *testing.T) {
	tempDir := "/tmp"
	cache := NewCachedProvider(nil, tempDir, 1*time.Hour)

	tests := []struct {
		city     string
		expected string
	}{
		{"Zurich", filepath.Join(tempDir, "weather_zurich.json")},
		{" Geneva ", filepath.Join(tempDir, "weather_geneva.json")},
		{"St. Gallen", filepath.Join(tempDir, "weather_st._gallen.json")},
	}

	for _, tt := range tests {
		got := cache.getCacheFilename(tt.city)
		if got != tt.expected {
			t.Errorf("getCacheFilename(%q) = %q; want %q", tt.city, got, tt.expected)
		}
	}
}

func TestWeatherImageCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "weather_image_cache_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	c := NewWeatherImageCache(tempDir, 1*time.Minute)
	key := c.GenerateKey("Zurich", 800, 480, "spectra6")
	data := []byte("fake-png-data")

	// Test Save
	if err := c.SaveImage(key, data); err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}

	// Test Get
	got, err := c.GetImage(key)
	if err != nil {
		t.Fatalf("GetImage failed: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("expected %s, got %s", data, got)
	}

	// Test Stale
	c.ttl = -1 * time.Second
	_, err = c.GetImage(key)
	if err == nil {
		t.Error("expected error for stale cache entry, got nil")
	}

	// Test Missing
	_, err = c.GetImage("missing")
	if err == nil {
		t.Error("expected error for missing cache entry, got nil")
	}
}

func TestIsWeatherFresh(t *testing.T) {
	// Case 1: Currently 10:00 AM, data from 05:00 AM (Today) -> Fresh
	now1 := time.Date(2026, 3, 7, 10, 0, 0, 0, time.UTC)
	fetched1 := time.Date(2026, 3, 7, 5, 0, 0, 0, time.UTC)
	if !IsWeatherFresh(fetched1, now1) {
		t.Errorf("Expected data from %v to be fresh at %v", fetched1, now1)
	}

	// Case 2: Currently 10:00 AM, data from 03:00 AM (Today) -> Stale
	fetched2 := time.Date(2026, 3, 7, 3, 0, 0, 0, time.UTC)
	if IsWeatherFresh(fetched2, now1) {
		t.Errorf("Expected data from %v to be stale at %v", fetched2, now1)
	}

	// Case 3: Currently 02:00 AM, data from 23:00 (Yesterday) -> Fresh
	now2 := time.Date(2026, 3, 7, 2, 0, 0, 0, time.UTC)
	fetched3 := time.Date(2026, 3, 6, 23, 0, 0, 0, time.UTC)
	if !IsWeatherFresh(fetched3, now2) {
		t.Errorf("Expected data from %v to be fresh at %v", fetched3, now2)
	}

	// Case 4: Currently 02:00 AM, data from 03:00 (Yesterday) -> Stale
	fetched4 := time.Date(2026, 3, 6, 3, 0, 0, 0, time.UTC)
	if IsWeatherFresh(fetched4, now2) {
		t.Errorf("Expected data from %v to be stale at %v", fetched4, now2)
	}
}

func TestCachedProvider_Error(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "weather_cache_err_test")
	defer os.RemoveAll(tempDir)

	// Provider that always fails
	mock := &errorProvider{}
	cache := NewCachedProvider(mock, tempDir, 1*time.Hour)

	_, err := cache.GetForecast("London")
	if err == nil {
		t.Error("Expected error from GetForecast when provider fails, got nil")
	}
}

type errorProvider struct{}

func (p *errorProvider) GetForecast(city string) (*WeatherForecast, error) {
	return nil, fmt.Errorf("provider error")
}
