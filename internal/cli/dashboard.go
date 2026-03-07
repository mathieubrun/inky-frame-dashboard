package cli

import (
	"fmt"

	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/agenda"
	"inky-frame-dashboard/internal/core/dashboard"
	"inky-frame-dashboard/internal/core/weather"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Operations for the full dashboard",
}

var dashboardImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate a combined weather and agenda image",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		location := viper.GetString("dashboard_location")
		output := viper.GetString("dashboard_output")
		palette := viper.GetString("dashboard_palette")
		useMock := viper.GetBool("dashboard_mock")

		// --- 1. Fetch Weather ---
		var wProvider weather.Provider
		if useMock {
			wProvider = weather.NewMockProvider()
		} else {
			wProvider = weather.NewOpenMeteoProvider()
		}
		wProvider = weather.NewCachedProvider(wProvider, cfg.WeatherCacheDir, cfg.WeatherCacheTTL)

		wForecast, err := wProvider.GetForecast(location)
		if err != nil {
			core.ErrorLogger.Printf("Dashboard: Failed to fetch weather: %v", err)
		}

		// --- 2. Fetch Agenda ---
		var aProvider agenda.CalendarProvider
		if useMock {
			aProvider = agenda.NewMockCalendarProvider()
		} else {
			aProvider = agenda.NewGoogleCalendarProvider(cfg.GoogleCredentials)
		}
		aProvider = agenda.NewCachedProvider(aProvider, cfg.AgendaCacheDir, cfg.AgendaCacheTTL)

		aForecast, err := aProvider.GetAgenda(cfg.AgendaID, 8)
		if err != nil {
			core.ErrorLogger.Printf("Dashboard: Failed to fetch agenda: %v", err)
			aForecast = &agenda.AgendaForecast{Events: []agenda.AgendaEvent{}}
		}

		// --- 3. Render Combined Image ---
		renderer := dashboard.NewDashboardRenderer(cfg.FontPath)
		data, err := renderer.Render(wForecast, aForecast, palette)
		if err != nil {
			core.ErrorLogger.Fatalf("Dashboard: Failed to render image: %v", err)
		}

		// Write to output file
		if err := core.WriteFile(output, data); err != nil {
			core.ErrorLogger.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Dashboard image saved to %s\n", output)
	},
}

func init() {
	dashboardImageCmd.Flags().StringP("location", "l", "Zurich", "Weather location")
	dashboardImageCmd.Flags().StringP("output", "o", "dashboard.png", "Output file path")
	dashboardImageCmd.Flags().StringP("palette", "p", "spectra6", "Color palette (spectra6, grayscale)")
	dashboardImageCmd.PersistentFlags().Bool("mock", false, "Use mock data")

	_ = viper.BindPFlag("dashboard_location", dashboardImageCmd.Flags().Lookup("location"))
	_ = viper.BindPFlag("dashboard_output", dashboardImageCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("dashboard_palette", dashboardImageCmd.Flags().Lookup("palette"))
	_ = viper.BindPFlag("dashboard_mock", dashboardImageCmd.PersistentFlags().Lookup("mock"))

	dashboardCmd.AddCommand(dashboardImageCmd)
	rootCmd.AddCommand(dashboardCmd)
}
