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

	var limit, offset int
	var format string
	var fields []string
	query := &cobra.Command{
		Use:   "query <index> <query>",
		Short: "Search the index",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := search.QueryIndexWithFields(args[0], args[1], limit, offset, fields)
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
					fmt.Printf("  %s (score: %.4f)", h.ID, h.Score)
					if len(h.Fields) > 0 {
						for k, v := range h.Fields {
							fmt.Printf(" %s=%s", k, v)
						}
					}
					fmt.Println()
				}
			}
			return nil
		},
	}
	query.Flags().IntVar(&limit, "limit", 10, "max results")
	query.Flags().IntVar(&offset, "offset", 0, "skip first N results")
	query.Flags().StringVar(&format, "format", "table", "output format: table|json")
	query.Flags().StringArrayVar(&fields, "field", nil, "fields to return (default: all)")

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

	bulk := &cobra.Command{
		Use:   "bulk <index> <ndjson-file>",
		Short: "Bulk index documents from NDJSON file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := search.BulkIndex(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("indexed %d documents\n", count)
			return nil
		},
	}

	deleteDoc := &cobra.Command{
		Use:   "delete <index> <docID>",
		Short: "Delete a document by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return search.DeleteDoc(args[0], args[1])
		},
	}

	getDoc := &cobra.Command{
		Use:   "get <index> <docID>",
		Short: "Get a document by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doc, err := search.GetDoc(args[0], args[1])
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(doc)
		},
	}

	updateDoc := &cobra.Command{
		Use:   "update <index> <docID> <json>",
		Short: "Update (re-index) a document",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return search.UpdateDoc(args[0], args[1], args[2])
		},
	}

	dropIdx := &cobra.Command{
		Use:   "drop <index>",
		Short: "Delete a search index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := search.DropIndex(args[0]); err != nil {
				return err
			}
			fmt.Printf("dropped %s\n", args[0])
			return nil
		},
	}

	cmd.AddCommand(create, index, query, info, dbs, bulk, deleteDoc, getDoc, updateDoc, dropIdx)
	return cmd
}
