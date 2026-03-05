package weather

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
			return forecast, nil
		}
	}

	// Fetch from the underlying provider
	forecast, err = p.provider.GetForecast(city)
	if err != nil {
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
		return
	}

	data, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(filename, data, 0644)
}

// Ensure CachedProvider implements Provider.
var _ Provider = (*CachedProvider)(nil)
