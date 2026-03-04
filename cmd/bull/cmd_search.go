package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bull-cli/bull/internal/search"
	"github.com/spf13/cobra"
)

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Full-text search operations (bleve)",
	}

	create := &cobra.Command{
		Use:   "create <index>",
		Short: "Create a new search index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return search.Create(args[0])
		},
	}

	index := &cobra.Command{
		Use:   "index <index> <docID> <json>",
		Short: "Index a document",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return search.Index(args[0], args[1], args[2])
		},
	}

	var limit int
	var format string
	query := &cobra.Command{
		Use:   "query <index> <query>",
		Short: "Search the index",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := search.QueryIndex(args[0], args[1], limit)
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				enc.Encode(result)
			default:
				fmt.Printf("Total: %d\n", result.Total)
				for _, h := range result.Hits {
					fmt.Printf("  %s (score: %.4f)\n", h.ID, h.Score)
				}
			}
			return nil
		},
	}
	query.Flags().IntVar(&limit, "limit", 10, "max results")
	query.Flags().StringVar(&format, "format", "table", "output format: table|json")

	info := &cobra.Command{
		Use:   "info <index>",
		Short: "Show index info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := search.Info(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Index: %s\nDocuments: %d\n", args[0], count)
			return nil
		},
	}

	dbs := &cobra.Command{
		Use:   "dbs",
		Short: "List all search indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := search.ListDBs()
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	cmd.AddCommand(create, index, query, info, dbs)
	return cmd
}
