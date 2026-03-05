package cli

import (
	"inky-frame-dashboard/internal/api"
	"inky-frame-dashboard/internal/config"
	"inky-frame-dashboard/internal/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			core.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
		}

		server := api.NewServer(cfg)
		if err := server.Start(); err != nil {
			core.ErrorLogger.Fatalf("Failed to start server: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	if err := viper.BindPFlag("port", serveCmd.Flags().Lookup("port")); err != nil {
		core.ErrorLogger.Fatalf("Failed to bind port flag: %v", err)
	}

	rootCmd.AddCommand(serveCmd)
}
