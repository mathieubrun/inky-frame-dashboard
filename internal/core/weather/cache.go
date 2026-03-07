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

// WeatherImageCache handles caching of generated weather images.
type WeatherImageCache struct {
	cacheDir string
	ttl      time.Duration
}

// NewWeatherImageCache creates a new WeatherImageCache.
func NewWeatherImageCache(cacheDir string, ttl time.Duration) *WeatherImageCache {
	return &WeatherImageCache{
		cacheDir: cacheDir,
		ttl:      ttl,
	}
}

// GenerateKey creates a unique key for a weather image based on request parameters.
func (c *WeatherImageCache) GenerateKey(city string, width, height int, palette string) string {
	normalized := strings.ToLower(strings.TrimSpace(city))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return fmt.Sprintf("%s_%dx%d_%s", normalized, width, height, palette)
}

// GetImage returns a cached image if it exists and is fresh.
func (c *WeatherImageCache) GetImage(key string) ([]byte, error) {
	filename := filepath.Join(c.cacheDir, key+".png")

	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if time.Since(info.ModTime()) > c.ttl {
		return nil, fmt.Errorf("cache entry for image %s is stale", key)
	}

	return os.ReadFile(filename)
}

// SaveImage saves an image to the cache.
func (c *WeatherImageCache) SaveImage(key string, data []byte) error {
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(c.cacheDir, key+".png")
	return os.WriteFile(filename, data, 0644)
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
