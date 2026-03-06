package agenda

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"inky-frame-dashboard/internal/core"
)

// CachedProvider wraps another CalendarProvider and caches results in files.
type CachedProvider struct {
	provider CalendarProvider
	cacheDir string
	ttl      time.Duration
}

// NewCachedProvider creates a new CachedProvider.
func NewCachedProvider(provider CalendarProvider, cacheDir string, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		provider: provider,
		cacheDir: cacheDir,
		ttl:      ttl,
	}
}

// GetAgenda returns agenda data, using the cache if it's available and fresh.
func (p *CachedProvider) GetAgenda(calendarID string, count int) (*AgendaForecast, error) {
	filename := p.getCacheFilename(calendarID, count)

	// Try to load from cache
	forecast, err := p.loadFromCache(filename)
	if err == nil {
		// Check if the cache is fresh
		if time.Since(forecast.FetchedAt) < p.ttl {
			core.InfoLogger.Printf("Agenda cache hit for calendar: %s", calendarID)
			return forecast, nil
		}
		core.InfoLogger.Printf("Agenda cache entry for calendar: %s is stale", calendarID)
	} else {
		core.InfoLogger.Printf("Agenda cache miss for calendar: %s", calendarID)
	}

	// Fetch from the underlying provider
	core.InfoLogger.Printf("Fetching agenda data for calendar: %s from provider", calendarID)
	forecast, err = p.provider.GetAgenda(calendarID, count)
	if err != nil {
		core.ErrorLogger.Printf("Failed to fetch agenda data for calendar: %s: %v", calendarID, err)
		return nil, err
	}

	// Save to cache
	p.saveToCache(filename, forecast)

	return forecast, nil
}

func (p *CachedProvider) getCacheFilename(calendarID string, count int) string {
	normalized := strings.ToLower(strings.TrimSpace(calendarID))
	normalized = strings.ReplaceAll(normalized, "@", "_at_")
	normalized = strings.ReplaceAll(normalized, ".", "_")
	return filepath.Join(p.cacheDir, fmt.Sprintf("agenda_%s_%d.json", normalized, count))
}

func (p *CachedProvider) loadFromCache(filename string) (*AgendaForecast, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var forecast AgendaForecast
	if err := json.Unmarshal(data, &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}

func (p *CachedProvider) saveToCache(filename string, forecast *AgendaForecast) {
	// Ensure the cache directory exists
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		core.ErrorLogger.Printf("Failed to create agenda cache directory %s: %v", filepath.Dir(filename), err)
		return
	}

	data, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		core.ErrorLogger.Printf("Failed to marshal agenda data: %v", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		core.ErrorLogger.Printf("Failed to write agenda cache file %s: %v", filename, err)
	} else {
		core.InfoLogger.Printf("Saved agenda data to cache")
	}
}

// Ensure CachedProvider implements CalendarProvider.
var _ CalendarProvider = (*CachedProvider)(nil)
