package api

import (
	"bytes"
	"encoding/json"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core/battery"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestBatteryHandlers(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "battery-api-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	csvPath := filepath.Join(tmpDir, "battery.csv")
	cfg := &config.Config{BatteryCSVPath: csvPath}
	server := NewServer(cfg)

	// Test POST /api/v1/battery
	t.Run("ReportBattery", func(t *testing.T) {
		reqBody := []byte(`{"voltage": 3.75}`)
		req, err := http.NewRequest("POST", "/api/v1/battery", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.BatteryReportHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var report battery.BatteryReport
		if err := json.NewDecoder(rr.Body).Decode(&report); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if report.Voltage != 3.75 {
			t.Errorf("expected voltage 3.75, got %f", report.Voltage)
		}
	})

	// Test GET /api/v1/battery/status
	t.Run("GetStatus", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/battery/status", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.BatteryStatusHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var report battery.BatteryReport
		if err := json.NewDecoder(rr.Body).Decode(&report); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if report.Voltage != 3.75 {
			t.Errorf("expected voltage 3.75, got %f", report.Voltage)
		}
	})

	// Test GET /api/v1/battery/history
	t.Run("GetHistory", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/battery/history", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.BatteryHistoryHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Errorf("wrong content type: got %v", rr.Header().Get("Content-Type"))
		}
	})
}
