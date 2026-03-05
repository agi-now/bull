package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agi-now/bull/internal/config"
	"github.com/spf13/cobra"
)

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show summary of all databases across engines",
		RunE: func(cmd *cobra.Command, args []string) error {
			engines := []struct {
				name string
				dir  string
			}{
				{"kv", config.KVDir()},
				{"sql", config.SQLDir()},
				{"graph", config.GraphDir()},
				{"search", config.SearchDir()},
				{"ts", config.TSDir()},
			}

			fmt.Printf("data-dir: %s\n\n", config.DataDir)

			for _, e := range engines {
				entries, err := os.ReadDir(e.dir)
				if err != nil {
					fmt.Printf("[%s] (empty)\n", e.name)
					continue
				}
				var names []string
				for _, entry := range entries {
					name := entry.Name()
					name = strings.TrimSuffix(name, filepath.Ext(name))
					names = append(names, name)
				}
				if len(names) == 0 {
					fmt.Printf("[%s] (empty)\n", e.name)
				} else {
					fmt.Printf("[%s] %s\n", e.name, strings.Join(names, ", "))
				}
			}
			return nil
		},
	}
}
