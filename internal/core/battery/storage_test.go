package battery

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStorage_AppendAndRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)

	now := time.Now().UTC().Truncate(time.Second)
	report := BatteryReport{
		Timestamp: now,
		Voltage:   3.75,
	}

	if err := storage.Append(report); err != nil {
		t.Fatalf("failed to append report: %v", err)
	}

	data, err := storage.ReadRaw()
	if err != nil {
		t.Fatalf("failed to read raw data: %v", err)
	}

	expectedHeader := "Timestamp,Voltage\n"
	expectedRow := now.Format(time.RFC3339) + ",3.75\n"
	expectedContent := expectedHeader + expectedRow

	if string(data) != expectedContent {
		t.Errorf("unexpected content.\nGot:\n%s\nWant:\n%s", string(data), expectedContent)
	}
}

func TestStorage_GetLatest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)

	t1 := time.Now().UTC().Truncate(time.Second).Add(-1 * time.Hour)
	t2 := time.Now().UTC().Truncate(time.Second)

	storage.Append(BatteryReport{Timestamp: t1, Voltage: 3.80})
	storage.Append(BatteryReport{Timestamp: t2, Voltage: 3.75})

	latest, err := storage.GetLatest()
	if err != nil {
		t.Fatalf("failed to get latest: %v", err)
	}

	if latest.Voltage != 3.75 {
		t.Errorf("expected voltage 3.75, got %f", latest.Voltage)
	}

	if !latest.Timestamp.Equal(t2) {
		t.Errorf("expected timestamp %v, got %v", t2, latest.Timestamp)
	}
}

func TestStorage_GetLatest_OnlyHeader(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-test-header")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)

	// Create file with only header
	os.MkdirAll(tmpDir, 0755)
	os.WriteFile(csvPath, []byte("Timestamp,Voltage\n"), 0644)

	_, err = storage.GetLatest()
	if err == nil {
		t.Error("expected error for header-only CSV, got nil")
	}
}

func TestStorage_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	storage := NewStorage(csvPath)

	storage.Append(BatteryReport{Timestamp: time.Now(), Voltage: 3.75})

	if err := storage.Clear(); err != nil {
		t.Fatalf("failed to clear storage: %v", err)
	}

	if _, err := os.Stat(csvPath); !os.IsNotExist(err) {
		t.Errorf("file still exists after clear")
	}
}
