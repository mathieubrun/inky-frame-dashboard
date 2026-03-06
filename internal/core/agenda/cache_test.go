package agenda

import (
	"os"
	"testing"
	"time"
)

func TestCachedProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agenda_cache_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mock := NewMockCalendarProvider()
	cache := NewCachedProvider(mock, tempDir, 1*time.Minute)

	calendarID := "test@example.com"
	count := 5

	// 1. Initial fetch (Cache miss)
	forecast, err := cache.GetAgenda(calendarID, count)
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}
	if len(forecast.Events) != count {
		t.Errorf("Expected %d events, got %d", count, len(forecast.Events))
	}

	// 2. Second fetch (Cache hit)
	forecast2, err := cache.GetAgenda(calendarID, count)
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}
	if !forecast2.FetchedAt.Equal(forecast.FetchedAt) {
		t.Errorf("Expected cached forecast, but FetchTime changed")
	}

	// 3. Stale cache
	cache.ttl = -1 * time.Second
	forecast3, err := cache.GetAgenda(calendarID, count)
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}
	if forecast3.FetchedAt.Equal(forecast.FetchedAt) {
		t.Errorf("Expected new fetch for stale cache, but got cached one")
	}
}
