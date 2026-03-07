package battery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessor_AddReport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-proc-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)
	processor := NewProcessor(storage)

	// Valid report
	report, err := processor.AddReport(3.75)
	if err != nil {
		t.Fatalf("failed to add valid report: %v", err)
	}
	if report.Voltage != 3.75 {
		t.Errorf("expected voltage 3.75, got %f", report.Voltage)
	}

	// Invalid report (negative voltage)
	_, err = processor.AddReport(-1.0)
	if err == nil {
		t.Error("expected error for negative voltage, got nil")
	}
}

func TestProcessor_GetLatestAndHistory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-proc-history-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)
	processor := NewProcessor(storage)

	// GetLatest when empty
	_, err = processor.GetLatest()
	if err == nil {
		t.Error("expected error for empty history, got nil")
	}

	// Add data
	processor.AddReport(3.75)
	processor.AddReport(3.70)

	// GetLatest
	latest, err := processor.GetLatest()
	if err != nil {
		t.Fatalf("failed to get latest: %v", err)
	}
	if latest.Voltage != 3.70 {
		t.Errorf("expected voltage 3.70, got %f", latest.Voltage)
	}

	// GetHistoryRaw
	history, err := processor.GetHistoryRaw()
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}
	if len(history) == 0 {
		t.Error("expected history data, got empty")
	}

	// ClearHistory
	if err := processor.ClearHistory(); err != nil {
		t.Fatalf("failed to clear history: %v", err)
	}
	_, err = processor.GetLatest()
	if err == nil {
		t.Error("expected error for cleared history, got nil")
	}
}
