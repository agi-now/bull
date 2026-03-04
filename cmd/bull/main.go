package main

import (
	"fmt"
	"os"

	"github.com/bull-cli/bull/internal/config"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "bull",
	Short: "Bull - All-in-One embedded engine toolkit",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.DataDir, "data-dir", "./data", "data directory path")
	rootCmd.AddCommand(
		kvCmd(),
		sqlCmd(),
		graphCmd(),
		searchCmd(),
		tsCmd(),
		serveCmd(),
		versionCmd(),
		infoCmd(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
