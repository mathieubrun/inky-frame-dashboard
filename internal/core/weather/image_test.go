package weather

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"
	"time"

	"inky-frame-dashboard/internal/core"

	"github.com/fogleman/gg"
)

func TestWeatherImageRenderer_Render(t *testing.T) {
	forecast := &WeatherForecast{
		Location: Location{
			City:    "Zurich",
			Country: "CH",
		},
		Current: WeatherRecord{
			Temperature: 22.5,
			Condition:   "Clear sky",
			Timestamp:   time.Now(),
		},
		FetchedAt: time.Now(),
	}

	req := &ImageRequest{
		Location: "Zurich",
		Width:    800,
		Height:   480,
	}

	renderer := NewWeatherImageRenderer("") // No font for now
	data, err := renderer.Render(forecast, req)
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

func TestWeatherImageRenderer_Palette(t *testing.T) {
	forecast := &WeatherForecast{
		Location: Location{City: "Zurich"},
		Current:  WeatherRecord{Temperature: 22.5, Condition: "Clear sky"},
	}

	req := &ImageRequest{
		Location: "Zurich",
		Width:    100,
		Height:   100,
		Palette:  "spectra6",
	}

	renderer := NewWeatherImageRenderer("")
	data, err := renderer.Render(forecast, req)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode rendered image: %v", err)
	}

	paletted, ok := img.(*image.Paletted)
	if !ok {
		t.Fatal("Expected paletted image when spectra6 is requested")
	}

	if len(paletted.Palette) != len(core.Spectra6Palette) {
		t.Errorf("Expected palette size %d, got %d", len(core.Spectra6Palette), len(paletted.Palette))
	}
}

func TestDrawWeatherIcon(t *testing.T) {
	dc := gg.NewContext(100, 100)
	
	conditions := []string{"Sun", "Clear", "Cloudy", "Rain", "Unknown"}
	for _, cond := range conditions {
		t.Run(cond, func(t *testing.T) {
			core.DrawWeatherIcon(dc, 50, 50, 50, cond)
			// If it doesn't panic, it's at least not crashing.
		})
	}
}
