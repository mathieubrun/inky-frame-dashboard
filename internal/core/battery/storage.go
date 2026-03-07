package battery

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Storage handles thread-safe persistence of battery reports to a CSV file.
type Storage struct {
	csvPath string
	mu      sync.Mutex
}

// NewStorage creates a new storage instance.
func NewStorage(csvPath string) *Storage {
	return &Storage{
		csvPath: csvPath,
	}
}

// Append adds a new battery report to the CSV file.
func (s *Storage) Append(report BatteryReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(s.csvPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists to determine if we need to write a header
	fileExists := true
	if _, err := os.Stat(s.csvPath); os.IsNotExist(err) {
		fileExists = false
	}

	f, err := os.OpenFile(s.csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write header if new file
	if !fileExists {
		if err := writer.Write([]string{"Timestamp", "Voltage"}); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	// Write report
	row := []string{
		report.Timestamp.Format(time.RFC3339),
		strconv.FormatFloat(report.Voltage, 'f', 2, 64),
	}
	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write row: %w", err)
	}

	return nil
}

// ReadRaw returns the raw content of the CSV file.
func (s *Storage) ReadRaw() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return just the header if file doesn't exist
			return []byte("Timestamp,Voltage\n"), nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// GetLatest returns the most recent battery report.
func (s *Storage) GetLatest() (BatteryReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			return BatteryReport{}, fmt.Errorf("no battery history available")
		}
		return BatteryReport{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return BatteryReport{}, fmt.Errorf("failed to read csv: %w", err)
	}

	if len(records) <= 1 {
		return BatteryReport{}, fmt.Errorf("no battery history available")
	}

	// Get the last record
	lastRecord := records[len(records)-1]

	timestamp, err := time.Parse(time.RFC3339, lastRecord[0])
	if err != nil {
		return BatteryReport{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	voltage, err := strconv.ParseFloat(lastRecord[1], 64)
	if err != nil {
		return BatteryReport{}, fmt.Errorf("failed to parse voltage: %w", err)
	}

	return BatteryReport{
		Timestamp: timestamp,
		Voltage:   voltage,
	}, nil
}

// Clear removes all battery history.
func (s *Storage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.Remove(s.csvPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}
