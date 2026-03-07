package battery

import (
	"fmt"
	"time"
)

// Processor handles high-level battery report logic.
type Processor struct {
	storage *Storage
}

// NewProcessor creates a new processor.
func NewProcessor(storage *Storage) *Processor {
	return &Processor{
		storage: storage,
	}
}

// AddReport validates and persists a new battery report.
func (p *Processor) AddReport(voltage float64) (BatteryReport, error) {
	if voltage < 0 {
		return BatteryReport{}, fmt.Errorf("voltage cannot be negative")
	}

	report := BatteryReport{
		Timestamp: time.Now().UTC(),
		Voltage:   voltage,
	}

	if err := p.storage.Append(report); err != nil {
		return BatteryReport{}, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

// GetLatest returns the latest battery report.
func (p *Processor) GetLatest() (BatteryReport, error) {
	return p.storage.GetLatest()
}

// GetHistoryRaw returns the raw history.
func (p *Processor) GetHistoryRaw() ([]byte, error) {
	return p.storage.ReadRaw()
}

// ClearHistory clears all battery history.
func (p *Processor) ClearHistory() error {
	return p.storage.Clear()
}
