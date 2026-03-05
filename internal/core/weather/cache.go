package weather

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"inky-frame-dashboard/internal/core"
)

// CachedProvider wraps another Provider and caches results in files.
type CachedProvider struct {
	provider Provider
	cacheDir string
	ttl      time.Duration
}

// NewCachedProvider creates a new CachedProvider.
func NewCachedProvider(provider Provider, cacheDir string, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		provider: provider,
		cacheDir: cacheDir,
		ttl:      ttl,
	}
}

// GetForecast returns weather data, using the cache if it's available and fresh.
func (p *CachedProvider) GetForecast(city string) (*WeatherForecast, error) {
	filename := p.getCacheFilename(city)

	// Try to load from cache
	forecast, err := p.loadFromCache(filename)
	if err == nil {
		// Check if the cache is fresh
		if time.Since(forecast.FetchedAt) < p.ttl {
			core.InfoLogger.Printf("Cache hit for city: %s", city)
			return forecast, nil
		}
		core.InfoLogger.Printf("Cache entry for city: %s is stale", city)
	} else {
		core.InfoLogger.Printf("Cache miss for city: %s", city)
	}

	// Fetch from the underlying provider
	core.InfoLogger.Printf("Fetching weather data for city: %s from provider", city)
	forecast, err = p.provider.GetForecast(city)
	if err != nil {
		core.ErrorLogger.Printf("Failed to fetch weather data for city: %s: %v", city, err)
		return nil, err
	}

	// Save to cache
	p.saveToCache(filename, forecast)

	return forecast, nil
}

func (p *CachedProvider) getCacheFilename(city string) string {
	normalized := strings.ToLower(strings.TrimSpace(city))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return filepath.Join(p.cacheDir, fmt.Sprintf("weather_%s.json", normalized))
}

func (p *CachedProvider) loadFromCache(filename string) (*WeatherForecast, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var forecast WeatherForecast
	if err := json.Unmarshal(data, &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}

func (p *CachedProvider) saveToCache(filename string, forecast *WeatherForecast) {
	// Ensure the cache directory exists
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		core.ErrorLogger.Printf("Failed to create cache directory %s: %v", filepath.Dir(filename), err)
		return
	}

	data, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		core.ErrorLogger.Printf("Failed to marshal weather data for city %s: %v", forecast.Location.City, err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		core.ErrorLogger.Printf("Failed to write cache file %s: %v", filename, err)
	} else {
		core.InfoLogger.Printf("Saved weather data to cache for city: %s", forecast.Location.City)
	}
}

// Ensure CachedProvider implements Provider.
var _ Provider = (*CachedProvider)(nil)
