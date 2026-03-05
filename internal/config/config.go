package config

import (
	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	Port int
}

// Load loads the configuration from environment variables and flags.
func Load() (*Config, error) {
	viper.SetDefault("port", 8080)
	viper.AutomaticEnv()

	return &Config{
		Port: viper.GetInt("port"),
	}, nil
}
