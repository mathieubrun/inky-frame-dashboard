package cli

import (
	"bytes"
	"inky-frame-dashboard/internal/core/battery"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBatteryCLI(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-cli-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	os.Setenv("BATTERY_CSV_PATH", csvPath)
	defer os.Unsetenv("BATTERY_CSV_PATH")

	storage := battery.NewStorage(csvPath)
	storage.Append(battery.BatteryReport{Voltage: 3.75})

	// Test battery history
	t.Run("History", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"battery", "history"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("failed to execute history command: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "3.75") {
			t.Errorf("expected 3.75 in output, got %s", output)
		}
	})

	// Test battery clear
	t.Run("Clear", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"battery", "clear"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("failed to execute clear command: %v", err)
		}

		if _, err := os.Stat(csvPath); !os.IsNotExist(err) {
			t.Errorf("file still exists after clear command")
		}
	})
}
