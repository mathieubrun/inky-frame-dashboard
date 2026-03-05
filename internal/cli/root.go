package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inky",
	Short: "Inky Frame Dashboard CLI",
	Long:  `A command line interface for managing the Inky Frame Dashboard.`,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Disable default --version flag to avoid conflict with the version subcommand if it were there,
    // though cobra's default --version is only added if version field is set in rootCmd.
    // However, the plan specifically says: "Global flag --version is explicitly excluded."
    rootCmd.Flags().BoolP("version", "v", false, "display version")
    _ = rootCmd.Flags().MarkHidden("version")
    // Let's just make sure we don't use it.
}
