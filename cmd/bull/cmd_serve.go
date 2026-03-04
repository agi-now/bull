package main

import (
	"fmt"

	"github.com/bull-cli/bull/internal/server"
	"github.com/spf13/cobra"
)

func serveCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := server.New(Version, BuildTime)
			addr := fmt.Sprintf(":%d", port)
			return srv.ListenAndServe(addr)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 2880, "listen port")
	return cmd
}
