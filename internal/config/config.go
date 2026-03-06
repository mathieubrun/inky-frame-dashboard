package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	Port                 int
	WeatherCacheDir      string
	WeatherCacheTTL      time.Duration
	WeatherImageCacheDir string
	WeatherImageCacheTTL time.Duration
	WeatherMock          bool
	FontPath             string
	GoogleCredentials    string
	AgendaID             string
	AgendaMock           bool
	AgendaCacheDir       string
	AgendaCacheTTL       time.Duration
}

// Load loads the configuration from environment variables and flags.
func Load() (*Config, error) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("weather_cache_dir", "./.inky/cache")
	viper.SetDefault("weather_cache_ttl", 1*time.Hour)
	viper.SetDefault("weather_image_cache_dir", "./.inky/image_cache")
	viper.SetDefault("weather_image_cache_ttl", 15*time.Minute)
	viper.SetDefault("weather_mock", false)
	viper.SetDefault("font_path", "")
	viper.SetDefault("google_credentials", "")
	viper.SetDefault("agenda_id", "primary")
	viper.SetDefault("agenda_mock", false)
	viper.SetDefault("agenda_cache_dir", "./.inky/cache/agenda")
	viper.SetDefault("agenda_cache_ttl", 15*time.Minute)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./.inky")
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	return &Config{
		Port:                 viper.GetInt("port"),
		WeatherCacheDir:      viper.GetString("weather_cache_dir"),
		WeatherCacheTTL:      viper.GetDuration("weather_cache_ttl"),
		WeatherImageCacheDir: viper.GetString("weather_image_cache_dir"),
		WeatherImageCacheTTL: viper.GetDuration("weather_image_cache_ttl"),
		WeatherMock:          viper.GetBool("weather_mock"),
		FontPath:             viper.GetString("font_path"),
		GoogleCredentials:    viper.GetString("google_credentials"),
		AgendaID:             viper.GetString("agenda_id"),
		AgendaMock:           viper.GetBool("agenda_mock"),
		AgendaCacheDir:       viper.GetString("agenda_cache_dir"),
		AgendaCacheTTL:       viper.GetDuration("agenda_cache_ttl"),
	}, nil
}
