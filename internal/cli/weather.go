package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/weather"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Retrieve weather forecast data for a Swiss city",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		city := viper.GetString("weather_city")
		if city == "" {
			core.ErrorLogger.Fatal("Error: --city flag or WEATHER_CITY environment variable is required")
		}

		useMock := viper.GetBool("weather_mock")
		useJSON := viper.GetBool("weather_json")

		var provider weather.Provider
		if useMock {
			provider = weather.NewMockProvider()
		} else {
			provider = weather.NewOpenMeteoProvider()
		}

		// Wrap with cache
		provider = weather.NewCachedProvider(provider, cfg.WeatherCacheDir, cfg.WeatherCacheTTL)

		forecast, err := provider.GetForecast(city)
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to retrieve weather: %v", err)
		}

		if useJSON {
			printJSON(forecast)
		} else {
			printTable(forecast)
		}
	},
}

func printJSON(forecast *weather.WeatherForecast) {
	data, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		core.ErrorLogger.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(data))
}

func printTable(forecast *weather.WeatherForecast) {
	fmt.Printf("Weather for %s, %s (Source: MeteoSwiss via Open-Meteo)\n", forecast.Location.City, forecast.Location.Country)
	fmt.Println("------------------------------------")
	
	cur := forecast.Current
	fmt.Printf("Current: %.1f°C | Wind: %.1f km/h (%.0f°) | Rain: %.1f mm (%.0f%%)\n",
		cur.Temperature, cur.WindSpeed, cur.WindDirection, cur.Precipitation, cur.PrecipitationProb)

	fmt.Println("\nNext 24h Forecast:")
	// Show a subset of forecast (e.g., every 3 hours)
	for i, r := range forecast.Hourly {
		if i%3 == 0 {
			fmt.Printf("- %s: %.1f°C | Rain: %.1f mm (%.0f%%)\n",
				r.Timestamp.Format("15:04"), r.Temperature, r.Precipitation, r.PrecipitationProb)
		}
	}

	fmt.Printf("\n(Fetched: %s)\n", forecast.FetchedAt.Format(time.RFC1123))
}

func init() {
	weatherCmd.Flags().StringP("city", "c", "", "Swiss city name")
	weatherCmd.Flags().Bool("mock", false, "Use mock data")
	weatherCmd.Flags().Bool("json", false, "Output in JSON format")

	_ = viper.BindPFlag("weather_city", weatherCmd.Flags().Lookup("city"))
	_ = viper.BindPFlag("weather_mock", weatherCmd.Flags().Lookup("mock"))
	_ = viper.BindPFlag("weather_json", weatherCmd.Flags().Lookup("json"))

	rootCmd.AddCommand(weatherCmd)
}
