package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/agenda"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "Retrieve upcoming calendar events",
}

var agendaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List upcoming events",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		calendarID := viper.GetString("agenda_id")
		count := viper.GetInt("agenda_count")
		useMock := viper.GetBool("agenda_mock")
		useJSON := viper.GetBool("agenda_json")

		var provider agenda.CalendarProvider
		if useMock {
			provider = agenda.NewMockCalendarProvider()
		} else {
			provider = agenda.NewGoogleCalendarProvider(cfg.GoogleCredentials)
		}

		// Wrap with cache
		provider = agenda.NewCachedProvider(provider, cfg.AgendaCacheDir, cfg.AgendaCacheTTL)

		forecast, err := provider.GetAgenda(calendarID, count)
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to retrieve agenda: %v", err)
		}

		if useJSON {
			printAgendaJSON(forecast)
		} else {
			printAgendaTable(forecast)
		}
	},
}

func printAgendaJSON(forecast *agenda.AgendaForecast) {
	data, err := json.MarshalIndent(forecast, "", "  ")
	if err != nil {
		core.ErrorLogger.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(data))
}

func printAgendaTable(forecast *agenda.AgendaForecast) {
	fmt.Printf("Upcoming events (Fetched: %s)\n", forecast.FetchedAt.Format(time.RFC1123))
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("%-20s | %-30s | %s\n", "Time", "Summary", "Location")
	fmt.Println("----------------------------------------------------------------------")

	for _, e := range forecast.Events {
		timeStr := e.StartTime.Format("02 Jan 15:04")
		fmt.Printf("%-20s | %-30s | %s\n", timeStr, e.Summary, e.Location)
	}
}

func init() {
	agendaListCmd.Flags().StringP("calendar-id", "c", "", "Google Calendar ID")
	agendaListCmd.Flags().IntP("count", "n", 10, "Number of events to list")
	agendaListCmd.PersistentFlags().Bool("mock", false, "Use mock agenda data")
	agendaListCmd.Flags().Bool("json", false, "Output in JSON format")

	_ = viper.BindPFlag("agenda_id", agendaListCmd.Flags().Lookup("calendar-id"))
	_ = viper.BindPFlag("agenda_count", agendaListCmd.Flags().Lookup("count"))
	_ = viper.BindPFlag("agenda_mock", agendaListCmd.PersistentFlags().Lookup("mock"))
	_ = viper.BindPFlag("agenda_json", agendaListCmd.Flags().Lookup("json"))

	agendaCmd.AddCommand(agendaListCmd)
	rootCmd.AddCommand(agendaCmd)
}
