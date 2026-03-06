package dashboard

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/agenda"
	"inky-frame-dashboard/internal/core/weather"

	"github.com/fogleman/gg"
)

// DashboardRenderer handles the generation of combined dashboard images.
type DashboardRenderer struct {
	fontPath string
}

// NewDashboardRenderer creates a new DashboardRenderer.
func NewDashboardRenderer(fontPath string) *DashboardRenderer {
	return &DashboardRenderer{
		fontPath: fontPath,
	}
}

// Render generates a split-screen dashboard image.
func (r *DashboardRenderer) Render(wForecast *weather.WeatherForecast, aForecast *agenda.AgendaForecast, palette string) ([]byte, error) {
	width := 800
	height := 480
	dc := gg.NewContext(width, height)

	// Draw background
	dc.SetRGB(1, 1, 1) // White
	dc.Clear()

	// --- Left Panel: Weather ---
	r.renderWeatherPanel(dc, wForecast)

	// --- Divider ---
	dc.SetColor(core.InkyBlack)
	dc.SetLineWidth(2)
	dc.DrawLine(400, 20, 400, 460)
	dc.Stroke()

	// --- Right Panel: Agenda ---
	r.renderAgendaPanel(dc, aForecast)

	img := dc.Image()
	if palette == "spectra6" {
		img = core.ConvertToPaletted(img, core.Spectra6Palette)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *DashboardRenderer) renderWeatherPanel(dc *gg.Context, forecast *weather.WeatherForecast) {
	centerX := 200.0
	
	// City
	dc.SetColor(core.InkyBlack)
	if r.fontPath != "" {
		_ = dc.LoadFontFace(r.fontPath, 32)
	}
	dc.DrawStringAnchored(forecast.Location.City, centerX, 60, 0.5, 0.5)

	// Icon
	core.DrawWeatherIcon(dc, centerX, 180, 160, forecast.Current.Condition)

	// Temperature
	dc.SetColor(core.InkyRed)
	if r.fontPath != "" {
		_ = dc.LoadFontFace(r.fontPath, 64)
	}
	dc.DrawStringAnchored(fmt.Sprintf("%.1f°C", forecast.Current.Temperature), centerX, 320, 0.5, 0.5)

	// Condition
	dc.SetColor(core.InkyBlack)
	if r.fontPath != "" {
		_ = dc.LoadFontFace(r.fontPath, 24)
	}
	dc.DrawStringAnchored(forecast.Current.Condition, centerX, 380, 0.5, 0.5)
}

func (r *DashboardRenderer) renderAgendaPanel(dc *gg.Context, forecast *agenda.AgendaForecast) {
	startX := 420.0
	startY := 60.0
	
	// Title
	dc.SetColor(core.InkyBlack)
	if r.fontPath != "" {
		_ = dc.LoadFontFace(r.fontPath, 28)
	}
	dc.DrawString(fmt.Sprintf("Upcoming (%d)", len(forecast.Events)), startX, startY)
	
	dc.SetLineWidth(1)
	dc.DrawLine(startX, startY+10, 780, startY+10)
	dc.Stroke()

	// Events
	eventY := startY + 50.0
	if r.fontPath != "" {
		_ = dc.LoadFontFace(r.fontPath, 18)
	}

	for i, e := range forecast.Events {
		if i >= 8 {
			break
		}
		
		// Time
		dc.SetColor(core.InkyBlue)
		timeStr := e.StartTime.Format("15:04")
		if e.StartTime.Day() != time.Now().Day() {
			timeStr = e.StartTime.Format("02 Jan 15:04")
		}
		dc.DrawString(timeStr, startX, eventY)
		
		// Summary
		dc.SetColor(core.InkyBlack)
		// Truncate summary if too long
		summary := e.Summary
		if len(summary) > 25 {
			summary = summary[:22] + "..."
		}
		dc.DrawString(summary, startX+100, eventY)
		
		eventY += 45.0
	}

	if len(forecast.Events) == 0 {
		dc.DrawString("No upcoming events", startX, eventY)
	}
}
