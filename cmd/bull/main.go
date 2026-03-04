package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bull-cli/bull/internal/config"
	"github.com/bull-cli/bull/internal/graph"
	"github.com/bull-cli/bull/internal/kv"
	"github.com/bull-cli/bull/internal/search"
	bsql "github.com/bull-cli/bull/internal/sql"
	"github.com/bull-cli/bull/internal/ts"
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
	rootCmd.AddCommand(kvCmd(), sqlCmd(), graphCmd(), searchCmd(), tsCmd(), serveCmd(), versionCmd(), infoCmd())
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bull %s\n", Version)
			fmt.Printf("build: %s\n", BuildTime)
			fmt.Printf("go: %s\n", runtime.Version())
			fmt.Printf("os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show data directory summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("data-dir: %s\n\n", config.DataDir)

			kvDBs, _ := kv.ListDBs()
			fmt.Printf("[KV]     %d databases\n", len(kvDBs))
			for _, n := range kvDBs {
				fmt.Printf("  - %s\n", n)
			}

			sqlDBs, _ := bsql.ListDBs()
			fmt.Printf("[SQL]    %d databases\n", len(sqlDBs))
			for _, n := range sqlDBs {
				fmt.Printf("  - %s\n", n)
			}

			graphDBs, _ := graph.ListDBs()
			fmt.Printf("[Graph]  %d databases\n", len(graphDBs))
			for _, n := range graphDBs {
				fmt.Printf("  - %s\n", n)
			}

			searchDBs, _ := search.ListDBs()
			fmt.Printf("[Search] %d indexes\n", len(searchDBs))
			for _, n := range searchDBs {
				fmt.Printf("  - %s\n", n)
			}

			tsDBs, _ := ts.ListDBs()
			fmt.Printf("[TS]     %d databases\n", len(tsDBs))
			for _, n := range tsDBs {
				fmt.Printf("  - %s\n", n)
			}
			return nil
		},
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
