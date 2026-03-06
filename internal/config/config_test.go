package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad_DotEnv(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a .env file in the temp dir
	envContent := "PORT=9090\nWEATHER_MOCK=true\n"
	err = os.WriteFile(filepath.Join(tempDir, ".env"), []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to create .env file: %v", err)
	}

	// Change working directory to temp dir for the test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Reset viper for the test
	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}

	if cfg.WeatherMock != true {
		t.Errorf("expected WeatherMock to be true")
	}
}
