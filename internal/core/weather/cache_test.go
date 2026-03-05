package weather

import (
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

	// 3. Third fetch with expired TTL
	cache.ttl = -1 * time.Second // Force expiry
	f3, err := cache.GetForecast(city)
	if err != nil {
		t.Fatalf("third fetch failed: %v", err)
	}

	if f3.FetchedAt.Equal(f1.FetchedAt) {
		t.Errorf("expected new data after expiry, but got cached data")
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
