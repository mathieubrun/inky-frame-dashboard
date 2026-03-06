package agenda

import "time"

// AgendaEvent represents a single event from a Google Calendar.
type AgendaEvent struct {
	Summary   string    `json:"summary"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Location  string    `json:"location,omitempty"`
}

// AgendaForecast represents a collection of upcoming events for a specific period.
type AgendaForecast struct {
	Events    []AgendaEvent `json:"events"`
	FetchedAt time.Time     `json:"fetched_at"`
}
