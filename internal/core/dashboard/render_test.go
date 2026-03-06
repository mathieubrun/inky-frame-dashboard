package dashboard

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"
	"time"

	"inky-frame-dashboard/internal/core/agenda"
	"inky-frame-dashboard/internal/core/weather"
)

func TestDashboardRenderer_Render(t *testing.T) {
	wForecast := &weather.WeatherForecast{
		Location: weather.Location{City: "Zurich"},
		Current:  weather.WeatherRecord{Temperature: 22.5, Condition: "Clear sky"},
		FetchedAt: time.Now(),
	}

	aForecast := &agenda.AgendaForecast{
		Events: []agenda.AgendaEvent{
			{Summary: "Test Event", StartTime: time.Now()},
		},
		FetchedAt: time.Now(),
	}

	renderer := NewDashboardRenderer("")
	data, err := renderer.Render(wForecast, aForecast, "spectra6")
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Render returned empty data")
	}

	// Verify it's a valid image
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode rendered image: %v", err)
	}

	if format != "png" {
		t.Errorf("Expected format png, got %s", format)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 800 || bounds.Dy() != 480 {
		t.Errorf("Expected dimensions 800x480, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
