package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	bsql "github.com/bull-cli/bull/internal/sql"
	"github.com/spf13/cobra"
)

func sqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql",
		Short: "SQL operations (SQLite)",
	}

	var format string

	exec := &cobra.Command{
		Use:   "exec <db> <sql>",
		Short: "Execute SQL statement",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			affected, err := bsql.Exec(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("%d rows affected\n", affected)
			return nil
		},
	}

	query := &cobra.Command{
		Use:   "query <db> <sql>",
		Short: "Query and display results",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bsql.Query(args[0], args[1])
			if err != nil {
				return err
			}
			switch format {
			case "csv":
				result.FormatCSV(os.Stdout)
			case "json":
				result.FormatJSON(os.Stdout)
			default:
				result.FormatTable(os.Stdout)
			}
			return nil
		},
	}
	query.Flags().StringVar(&format, "format", "table", "output format: table|csv|json")

	shell := &cobra.Command{
		Use:   "shell <db>",
		Short: "Interactive SQL shell",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbName := args[0]
			fmt.Printf("bull sql shell [%s] - type .quit to exit\n", dbName)
			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print("sql> ")
				if !scanner.Scan() {
					break
				}
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				if line == ".quit" || line == ".exit" {
					break
				}
				upper := strings.ToUpper(line)
				if strings.HasPrefix(upper, "SELECT") || strings.HasPrefix(upper, "PRAGMA") || strings.HasPrefix(upper, "EXPLAIN") {
					result, err := bsql.Query(dbName, line)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						continue
					}
					result.FormatTable(os.Stdout)
				} else {
					affected, err := bsql.Exec(dbName, line)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						continue
					}
					fmt.Printf("%d rows affected\n", affected)
				}
			}
			return nil
		},
	}

	importCmd := &cobra.Command{
		Use:   "import <db> <table> <file.csv>",
		Short: "Import CSV into table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bsql.ImportCSV(args[0], args[1], args[2])
		},
	}

	var output string
	export := &cobra.Command{
		Use:   "export <db> <sql>",
		Short: "Export query result to CSV",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bsql.Query(args[0], args[1])
			if err != nil {
				return err
			}
			w := os.Stdout
			if output != "" {
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				w = f
			}
			result.FormatCSV(w)
			return nil
		},
	}
	export.Flags().StringVarP(&output, "output", "o", "", "output file (default: stdout)")

	dbs := &cobra.Command{
		Use:   "dbs",
		Short: "List all SQL databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := bsql.ListDBs()
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	cmd.AddCommand(exec, query, shell, importCmd, export, dbs)
	return cmd
}
