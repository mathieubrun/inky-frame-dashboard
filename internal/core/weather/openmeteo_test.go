package weather

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "Zurich" {
			_, _ = fmt.Fprintln(w, `{"results":[{"name":"Zurich","latitude":47.3717,"longitude":8.5422,"country":"Switzerland"}]}`)
		} else if name == "8001" {
			_, _ = fmt.Fprintln(w, `{"results":[{"name":"Zurich","latitude":47.3717,"longitude":8.5422,"country":"Switzerland"}]}`)
		} else {
			_, _ = fmt.Fprintln(w, `{"results":[]}`)
		}
	}))
	defer ts.Close()

	p := NewOpenMeteoProvider()
	p.geocodingURL = ts.URL

	tests := []struct {
		city    string
		wantErr bool
		name    string
	}{
		{"Zurich", false, "Zurich"},
		{"8001", false, "Zurich"},
		{"Unknown", true, ""},
	}

	for _, tt := range tests {
		loc, err := p.geocode(tt.city)
		if (err != nil) != tt.wantErr {
			t.Errorf("geocode(%q) error = %v, wantErr %v", tt.city, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && loc.City != tt.name {
			t.Errorf("geocode(%q) = %v, want name %v", tt.city, loc.City, tt.name)
		}
	}
}

func TestWeatherCodeToCondition(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "Clear sky"},
		{1, "Mainly clear, partly cloudy, and overcast"},
		{95, "Thunderstorm: Slight or moderate"},
		{99, "Thunderstorm with slight and heavy hail"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		got := weatherCodeToCondition(tt.code)
		if got != tt.expected {
			t.Errorf("weatherCodeToCondition(%d) = %q; want %q", tt.code, got, tt.expected)
		}
	}
}

func TestGetForecast(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/search" {
			_, _ = fmt.Fprintln(w, `{"results":[{"name":"Zurich","latitude":47.3717,"longitude":8.5422,"country":"Switzerland"}]}`)
		} else {
			_, _ = fmt.Fprintln(w, `{"hourly":{"time":["2026-03-05T00:00"],"temperature_2m":[10.0],"weathercode":[0],"windspeed_10m":[5.0],"winddirection_10m":[180.0],"precipitation":[0.0],"precipitation_probability":[0.0]}}`)
		}
	}))
	defer ts.Close()

	p := NewOpenMeteoProvider()
	p.geocodingURL = ts.URL + "/v1/search"
	p.forecastURL = ts.URL + "/v1/forecast"

	forecast, err := p.GetForecast("Zurich")
	if err != nil {
		t.Fatalf("GetForecast failed: %v", err)
	}

	if forecast.Location.City != "Zurich" {
		t.Errorf("Expected city Zurich, got %s", forecast.Location.City)
	}
	if len(forecast.Hourly) != 1 {
		t.Errorf("Expected 1 hourly record, got %d", len(forecast.Hourly))
	}
}
