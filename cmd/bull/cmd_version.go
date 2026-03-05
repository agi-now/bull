package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version and build info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bull %s\n", Version)
			fmt.Printf("built:  %s\n", BuildTime)
			fmt.Printf("go:     %s\n", runtime.Version())
			fmt.Printf("os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}
