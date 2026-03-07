package agenda

import (
	"fmt"
	"time"
)

// MockCalendarProvider implements the CalendarProvider interface for testing.
type MockCalendarProvider struct{}

// NewMockCalendarProvider creates a new MockCalendarProvider.
func NewMockCalendarProvider() *MockCalendarProvider {
	return &MockCalendarProvider{}
}

// GetAgenda returns a set of mock events.
func (p *MockCalendarProvider) GetAgenda(calendarID string, count int) (*AgendaForecast, error) {
	now := time.Now()

	// Create some mock events
	events := make([]AgendaEvent, 0, count)

	titles := []string{
		"Morning Standup",
		"Project Review",
		"Lunch with Team",
		"Product Demo",
		"Focus Time",
		"Gym Session",
		"Dinner with Family",
		"Evening Walk",
	}

	locations := []string{
		"Meeting Room A",
		"Conference Hall",
		"Downtown Cafe",
		"Boardroom",
		"Home Office",
		"City Gym",
		"Trattoria",
		"Park",
	}

	for i := 0; i < count; i++ {
		title := titles[i%len(titles)]
		if i >= len(titles) {
			title = fmt.Sprintf("%s %d", title, i/len(titles)+1)
		}

		startTime := now.Add(time.Duration(i+1) * time.Hour)
		// Round to nearest 15 mins for cleaner look
		startTime = startTime.Truncate(15 * time.Minute)

		events = append(events, AgendaEvent{
			Summary:   title,
			StartTime: startTime,
			EndTime:   startTime.Add(1 * time.Hour),
			Location:  locations[i%len(locations)],
		})
	}

	return &AgendaForecast{
		Events:    events,
		FetchedAt: now,
	}, nil
}

// Ensure MockCalendarProvider implements CalendarProvider.
var _ CalendarProvider = (*MockCalendarProvider)(nil)
