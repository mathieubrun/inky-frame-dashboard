package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	Port            int
	WeatherCacheDir string
	WeatherCacheTTL time.Duration
	WeatherMock     bool
}

// Load loads the configuration from environment variables and flags.
func Load() (*Config, error) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("weather_cache_dir", "./.inky/cache")
	viper.SetDefault("weather_cache_ttl", 1*time.Hour)
	viper.SetDefault("weather_mock", false)

	viper.AutomaticEnv()

	return &Config{
		Port:            viper.GetInt("port"),
		WeatherCacheDir: viper.GetString("weather_cache_dir"),
		WeatherCacheTTL: viper.GetDuration("weather_cache_ttl"),
		WeatherMock:     viper.GetBool("weather_mock"),
	}, nil
}
