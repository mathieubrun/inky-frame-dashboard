package agenda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/calendar/v3"
)

func TestGoogleCalendarProvider_GetAgenda(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock events response
		resp := calendar.Events{
			Items: []*calendar.Event{
				{
					Summary: "Test Event",
					Start:   &calendar.EventDateTime{DateTime: "2026-03-06T10:00:00Z"},
					End:     &calendar.EventDateTime{DateTime: "2026-03-06T11:00:00Z"},
					Location: "Test Location",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	p := NewGoogleCalendarProvider("")
	p.apiURL = ts.URL

	forecast, err := p.GetAgenda("primary", 10)
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}

	if len(forecast.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(forecast.Events))
	}

	if forecast.Events[0].Summary != "Test Event" {
		t.Errorf("Expected summary 'Test Event', got %s", forecast.Events[0].Summary)
	}
}

func TestGoogleCalendarProvider_Validate(t *testing.T) {
	// 1. Unconfigured
	p1 := NewGoogleCalendarProvider("")
	if err := p1.Validate(); err == nil {
		t.Error("Expected error for unconfigured credentials, got nil")
	}

	// 2. Missing file
	p2 := NewGoogleCalendarProvider("non-existent.json")
	if err := p2.Validate(); err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

func TestParseGoogleTime(t *testing.T) {
	tests := []struct {
		name    string
		input   *calendar.EventDateTime
		wantErr bool
	}{
		{
			name: "DateTime",
			input: &calendar.EventDateTime{DateTime: "2026-03-06T10:00:00Z"},
			wantErr: false,
		},
		{
			name: "DateOnly",
			input: &calendar.EventDateTime{Date: "2026-03-06"},
			wantErr: false,
		},
		{
			name: "Nil",
			input: nil,
			wantErr: true,
		},
		{
			name: "Empty",
			input: &calendar.EventDateTime{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseGoogleTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGoogleTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
