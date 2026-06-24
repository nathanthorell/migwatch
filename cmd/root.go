package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	envFile    string
	envFilter  string
)

var rootCmd = &cobra.Command{
	Use:   "migwatch",
	Short: "Visualize database migration state across environments",
	RunE:  runSummary,
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show summary of migration state per environment (default)",
	RunE:  runSummary,
}

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "Show full migration history table per environment",
	RunE:  runFull,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default: ./migwatch.toml or ~/.config/migwatch/migwatch.toml)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", "", "env file (default: ./.env)")
	rootCmd.PersistentFlags().StringVarP(&envFilter, "env", "e", "", "filter to a single environment")

	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(fullCmd)
}
