package weather

import (
	"bytes"
	"fmt"
	"image/png"

	"inky-frame-dashboard/internal/core"

	"github.com/fogleman/gg"
)

// WeatherImageRenderer handles the generation of weather images.
type WeatherImageRenderer struct {
	// fontPath is the path to the TrueType font to use.
	fontPath string
}

// NewWeatherImageRenderer creates a new WeatherImageRenderer.
func NewWeatherImageRenderer(fontPath string) *WeatherImageRenderer {
	return &WeatherImageRenderer{
		fontPath: fontPath,
	}
}

// Render generates a weather image from a forecast and a request.
func (r *WeatherImageRenderer) Render(forecast *WeatherForecast, req *ImageRequest) ([]byte, error) {
	// Set default dimensions if not specified
	if req.Width <= 0 {
		req.Width = 800
	}
	if req.Height <= 0 {
		req.Height = 480
	}

	dc := gg.NewContext(req.Width, req.Height)

	// Draw background
	dc.SetRGB(1, 1, 1) // White
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0) // Black
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, 36); err != nil {
			return nil, fmt.Errorf("failed to load font: %w", err)
		}
	}
	dc.DrawStringAnchored(fmt.Sprintf("Weather for %s", forecast.Location.City), float64(req.Width)/2, 40, 0.5, 0.5)

	// Draw Icon
	core.DrawWeatherIcon(dc, float64(req.Width)/2, 160, 150, forecast.Current.Condition)

	// Temperature
	dc.SetRGB(1, 0, 0) // Red (testing colors)
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, 72); err != nil {
			return nil, fmt.Errorf("failed to load font: %w", err)
		}
	}
	dc.DrawStringAnchored(fmt.Sprintf("%.1f°C", forecast.Current.Temperature), float64(req.Width)/2, 300, 0.5, 0.5)

	// Condition Text
	dc.SetRGB(0, 0, 0) // Black
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, 24); err != nil {
			return nil, fmt.Errorf("failed to load font: %w", err)
		}
	}
	dc.DrawStringAnchored(forecast.Current.Condition, float64(req.Width)/2, 360, 0.5, 0.5)

	// Success indication (Bottom right)
	dc.SetRGB(0, 0, 0)
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, 14); err != nil {
			return nil, fmt.Errorf("failed to load font: %w", err)
		}
	}
	dc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", forecast.FetchedAt.Format("15:04:05")), float64(req.Width)-10, float64(req.Height)-10, 1, 1)

	img := dc.Image()
	if req.Palette == "spectra6" {
		img = core.ConvertToPaletted(img, core.Spectra6Palette)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
