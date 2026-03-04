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
	var queryLimit int

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
			sqlStr := args[1]
			if queryLimit > 0 {
				sqlStr = fmt.Sprintf("SELECT * FROM (%s) LIMIT %d", sqlStr, queryLimit)
			}
			result, err := bsql.Query(args[0], sqlStr)
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
	query.Flags().IntVar(&queryLimit, "limit", 0, "limit number of rows (0 = no limit)")

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
				if line == ".tables" {
					names, err := bsql.Tables(dbName)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
					} else {
						for _, n := range names {
							fmt.Println(n)
						}
					}
					continue
				}
				if strings.HasPrefix(line, ".schema") {
					parts := strings.Fields(line)
					if len(parts) < 2 {
						fmt.Fprintln(os.Stderr, "usage: .schema <table>")
					} else {
						ddl, err := bsql.Schema(dbName, parts[1])
						if err != nil {
							fmt.Fprintf(os.Stderr, "error: %v\n", err)
						} else {
							fmt.Println(ddl)
						}
					}
					continue
				}
				if strings.HasPrefix(line, ".count") {
					parts := strings.Fields(line)
					if len(parts) < 2 {
						fmt.Fprintln(os.Stderr, "usage: .count <table>")
					} else {
						n, err := bsql.CountRows(dbName, parts[1])
						if err != nil {
							fmt.Fprintf(os.Stderr, "error: %v\n", err)
						} else {
							fmt.Println(n)
						}
					}
					continue
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
			n, err := bsql.ImportCSV(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Printf("imported %d rows\n", n)
			return nil
		},
	}

	var output string
	var exportFormat string
	export := &cobra.Command{
		Use:   "export <db> <sql>",
		Short: "Export query result to file",
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
			switch exportFormat {
			case "json":
				result.FormatJSON(w)
			default:
				result.FormatCSV(w)
			}
			return nil
		},
	}
	export.Flags().StringVarP(&output, "output", "o", "", "output file (default: stdout)")
	export.Flags().StringVar(&exportFormat, "format", "csv", "output format: csv|json")

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

	tables := &cobra.Command{
		Use:   "tables <db>",
		Short: "List all tables",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := bsql.Tables(args[0])
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	schema := &cobra.Command{
		Use:   "schema <db> <table>",
		Short: "Show CREATE TABLE statement",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ddl, err := bsql.Schema(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(ddl)
			return nil
		},
	}

	countRows := &cobra.Command{
		Use:   "count <db> <table>",
		Short: "Count rows in a table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := bsql.CountRows(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(n)
			return nil
		},
	}

	execFile := &cobra.Command{
		Use:   "exec-file <db> <file.sql>",
		Short: "Execute SQL from a file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bsql.ExecFile(args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("ok")
			return nil
		},
	}

	importJSONCmd := &cobra.Command{
		Use:   "import-json <db> <table> <file.json>",
		Short: "Import JSON array into table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := bsql.ImportJSON(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Printf("imported %d rows\n", n)
			return nil
		},
	}

	describe := &cobra.Command{
		Use:   "describe <db> <table>",
		Short: "Show column info (PRAGMA table_info)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cols, err := bsql.Describe(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("%-4s %-20s %-10s %-8s %-10s %-4s\n", "CID", "NAME", "TYPE", "NOTNULL", "DEFAULT", "PK")
			for _, c := range cols {
				nn := ""
				if c.NotNull {
					nn = "YES"
				}
				pk := ""
				if c.PK {
					pk = "YES"
				}
				fmt.Printf("%-4d %-20s %-10s %-8s %-10s %-4s\n", c.CID, c.Name, c.Type, nn, c.Default, pk)
			}
			return nil
		},
	}

	importNDJSON := &cobra.Command{
		Use:   "import-ndjson <db> <table> <file.ndjson>",
		Short: "Import NDJSON (one JSON object per line) into table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := bsql.ImportNDJSON(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Printf("imported %d rows\n", n)
			return nil
		},
	}

	dropDB := &cobra.Command{
		Use:   "drop <db>",
		Short: "Delete a SQL database file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bsql.DropDB(args[0]); err != nil {
				return err
			}
			fmt.Printf("dropped %s\n", args[0])
			return nil
		},
	}

	cmd.AddCommand(exec, query, shell, importCmd, export, dbs, tables, schema, countRows, execFile, importJSONCmd, describe, importNDJSON, dropDB)
	return cmd
}
