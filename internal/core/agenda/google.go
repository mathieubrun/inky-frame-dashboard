package agenda

import (
	"context"
	"fmt"
	"os"
	"time"

	"inky-frame-dashboard/internal/core"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// GoogleCalendarProvider implements the CalendarProvider interface using Google Calendar API.
type GoogleCalendarProvider struct {
	credsPath string
	apiURL    string
}

// NewGoogleCalendarProvider creates a new GoogleCalendarProvider.
func NewGoogleCalendarProvider(credsPath string) *GoogleCalendarProvider {
	return &GoogleCalendarProvider{
		credsPath: credsPath,
	}
}

// Validate checks if the credentials file exists.
func (p *GoogleCalendarProvider) Validate() error {
	if p.credsPath == "" {
		return fmt.Errorf("google credentials path is not configured")
	}
	
	// Check if file exists
	if _, err := os.Stat(p.credsPath); err != nil {
		return fmt.Errorf("google credentials file not found at %s: %w", p.credsPath, err)
	}
	
	return nil
}

// GetAgenda fetches upcoming events from Google Calendar.
func (p *GoogleCalendarProvider) GetAgenda(calendarID string, count int) (*AgendaForecast, error) {
	ctx := context.Background()

	var opts []option.ClientOption
	if p.credsPath != "" {
		opts = append(opts, option.WithCredentialsFile(p.credsPath))
	}
	if p.apiURL != "" {
		opts = append(opts, option.WithEndpoint(p.apiURL))
		opts = append(opts, option.WithoutAuthentication())
	}
	opts = append(opts, option.WithScopes(calendar.CalendarReadonlyScope))

	srv, err := calendar.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	now := time.Now()
	tMin := now.Format(time.RFC3339)

	events, err := srv.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(tMin).
		MaxResults(int64(count)).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	agendaEvents := make([]AgendaEvent, 0, len(events.Items))
	for _, item := range events.Items {
		start, err := parseGoogleTime(item.Start)
		if err != nil {
			core.ErrorLogger.Printf("Failed to parse start time for event %s: %v", item.Summary, err)
			continue
		}

		end, err := parseGoogleTime(item.End)
		if err != nil {
			core.ErrorLogger.Printf("Failed to parse end time for event %s: %v", item.Summary, err)
			continue
		}

		agendaEvents = append(agendaEvents, AgendaEvent{
			Summary:   item.Summary,
			StartTime: start,
			EndTime:   end,
			Location:  item.Location,
		})
	}

	return &AgendaForecast{
		Events:    agendaEvents,
		FetchedAt: now,
	}, nil
}

func parseGoogleTime(gt *calendar.EventDateTime) (time.Time, error) {
	if gt == nil {
		return time.Time{}, fmt.Errorf("nil event date time")
	}

	if gt.DateTime != "" {
		return time.Parse(time.RFC3339, gt.DateTime)
	}

	if gt.Date != "" {
		return time.Parse("2006-01-02", gt.Date)
	}

	return time.Time{}, fmt.Errorf("no date or date-time provided")
}

// Ensure GoogleCalendarProvider implements CalendarProvider.
var _ CalendarProvider = (*GoogleCalendarProvider)(nil)
