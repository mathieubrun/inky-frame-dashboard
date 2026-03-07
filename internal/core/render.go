package core

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/fogleman/gg"
)

// CalculateMD5 returns the MD5 hash of the provided data as a hex string.
func CalculateMD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

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

// ConvertToPaletted converts an image to a paletted one using the provided palette.
func ConvertToPaletted(src image.Image, p color.Palette) *image.Paletted {
	bounds := src.Bounds()
	dst := image.NewPaletted(bounds, p)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// DrawWeatherIcon draws a simple weather icon based on the condition string.
func DrawWeatherIcon(dc *gg.Context, x, y, size float64, condition string) {
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
