package cli

import (
	"fmt"
	"inky-frame-dashboard/internal/config"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application version",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
