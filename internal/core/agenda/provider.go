package agenda

// CalendarProvider defines the interface for fetching calendar events.
type CalendarProvider interface {
	GetAgenda(calendarID string, count int) (*AgendaForecast, error)
}
