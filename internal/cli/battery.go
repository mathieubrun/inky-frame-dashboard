package cli

import (
	"fmt"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"
	"inky-frame-dashboard/internal/core/battery"

	"github.com/spf13/cobra"
)

var batteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Manage and view battery reports",
}

var batteryHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Print the full battery history as raw CSV",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		storage := battery.NewStorage(cfg.BatteryCSVPath)
		processor := battery.NewProcessor(storage)

		data, err := processor.GetHistoryRaw()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to retrieve battery history: %v", err)
		}

		fmt.Fprint(cmd.OutOrStdout(), string(data))
	},
}

var batteryClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all battery history",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		storage := battery.NewStorage(cfg.BatteryCSVPath)
		processor := battery.NewProcessor(storage)

		if err := processor.ClearHistory(); err != nil {
			core.ErrorLogger.Fatalf("Failed to clear battery history: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Battery history cleared successfully.")
	},
}

func init() {
	batteryCmd.AddCommand(batteryHistoryCmd)
	batteryCmd.AddCommand(batteryClearCmd)
	rootCmd.AddCommand(batteryCmd)
}
