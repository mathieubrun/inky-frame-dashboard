package battery

import "time"

// BatteryReport represents a single battery level measurement.
type BatteryReport struct {
	Timestamp time.Time `json:"timestamp"`
	Voltage   float64   `json:"voltage"`
}
