package weather

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	"github.com/fogleman/gg"
)

// Inky Palette colors
var (
	InkyBlack  = color.RGBA{0, 0, 0, 255}
	InkyWhite  = color.RGBA{255, 255, 255, 255}
	InkyRed    = color.RGBA{255, 0, 0, 255}
	InkyGreen  = color.RGBA{0, 255, 0, 255}
	InkyBlue   = color.RGBA{0, 0, 255, 255}
	InkyYellow = color.RGBA{255, 255, 0, 255}
)

var Spectra6Palette = color.Palette{
	InkyBlack, InkyWhite, InkyRed, InkyGreen, InkyBlue, InkyYellow,
}

func convertToPaletted(src image.Image, p color.Palette) *image.Paletted {
	bounds := src.Bounds()
	dst := image.NewPaletted(bounds, p)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

func drawWeatherIcon(dc *gg.Context, x, y, size float64, condition string) {
	dc.Push()
	defer dc.Pop()

	condition = strings.ToLower(condition)
	switch {
	case strings.Contains(condition, "sun") || strings.Contains(condition, "clear"):
		// Sun
		dc.SetColor(InkyYellow)
		dc.DrawCircle(x, y, size*0.4)
		dc.Fill()
	case strings.Contains(condition, "cloud"):
		// Cloud
		dc.SetColor(InkyBlue)
		dc.DrawEllipse(x, y, size*0.4, size*0.25)
		dc.DrawEllipse(x-size*0.2, y-size*0.1, size*0.25, size*0.2)
		dc.DrawEllipse(x+size*0.2, y-size*0.1, size*0.25, size*0.2)
		dc.Fill()
	case strings.Contains(condition, "rain"):
		// Rain
		dc.SetColor(InkyBlue)
		dc.DrawEllipse(x, y-size*0.1, size*0.4, size*0.25)
		dc.Fill()
		// Rain drops
		dc.SetColor(InkyBlue)
		dc.SetLineWidth(2)
		dc.DrawLine(x-size*0.2, y+size*0.1, x-size*0.25, y+size*0.3)
		dc.DrawLine(x, y+size*0.1, x-0.05*size, y+size*0.3)
		dc.DrawLine(x+size*0.2, y+size*0.1, x+size*0.15, y+size*0.3)
		dc.Stroke()
	default:
		// Unknown - Question mark
		dc.SetColor(InkyBlack)
		dc.DrawCircle(x, y, size*0.4)
		dc.Stroke()
	}
}

// WeatherImageRenderer handles the generation of weather images.
type WeatherImageRenderer struct {
	// icons maps weather condition codes or names to their image representation.
	icons map[string]image.Image
	// fontPath is the path to the TrueType font to use.
	fontPath string
}

// NewWeatherImageRenderer creates a new WeatherImageRenderer.
func NewWeatherImageRenderer(fontPath string) *WeatherImageRenderer {
	return &WeatherImageRenderer{
		icons:    make(map[string]image.Image),
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
	drawWeatherIcon(dc, float64(req.Width)/2, 160, 150, forecast.Current.Condition)

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
		img = convertToPaletted(img, Spectra6Palette)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
