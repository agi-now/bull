package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/agi-now/bull/internal/ts"
	"github.com/spf13/cobra"
)

func tsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ts",
		Short: "Time-series operations (tstorage)",
	}

	var timestamp int64
	var labelArgs []string

	parseLabels := func(args []string) map[string]string {
		m := make(map[string]string)
		for _, a := range args {
			parts := strings.SplitN(a, "=", 2)
			if len(parts) == 2 {
				m[parts[0]] = parts[1]
			}
		}
		return m
	}

	write := &cobra.Command{
		Use:   "write <db> <metric> <value>",
		Short: "Write a data point",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return fmt.Errorf("invalid value: %w", err)
			}
			return ts.Write(args[0], args[1], val, timestamp, parseLabels(labelArgs))
		},
	}
	write.Flags().Int64Var(&timestamp, "time", 0, "unix timestamp (default: now)")
	write.Flags().StringArrayVar(&labelArgs, "label", nil, "label (key=value)")

	var from, to int64
	var queryLabelArgs []string
	var format string
	query := &cobra.Command{
		Use:   "query <db> <metric>",
		Short: "Query data points",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now().Unix()
			if from == 0 {
				from = now - 3600
			}
			if to == 0 {
				to = now
			}
			points, err := ts.QueryRange(args[0], args[1], from, to, parseLabels(queryLabelArgs))
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				enc.Encode(points)
			case "csv":
				fmt.Println("timestamp,value")
				for _, p := range points {
					fmt.Printf("%d,%f\n", p.Timestamp, p.Value)
				}
			default:
				fmt.Printf("%-20s %s\n", "TIMESTAMP", "VALUE")
				for _, p := range points {
					fmt.Printf("%-20d %f\n", p.Timestamp, p.Value)
				}
			}
			return nil
		},
	}
	query.Flags().Int64Var(&from, "from", 0, "start unix timestamp (default: now - 1h)")
	query.Flags().Int64Var(&to, "to", 0, "end unix timestamp (default: now)")
	query.Flags().StringArrayVar(&queryLabelArgs, "label", nil, "label filter (key=value)")
	query.Flags().StringVar(&format, "format", "table", "output format: table|csv|json")

	var latestLabelArgs []string
	latest := &cobra.Command{
		Use:   "latest <db> <metric>",
		Short: "Get the latest data point",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := ts.Latest(args[0], args[1], parseLabels(latestLabelArgs))
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				enc.Encode(p)
			default:
				fmt.Printf("timestamp: %d\nvalue: %f\n", p.Timestamp, p.Value)
			}
			return nil
		},
	}
	latest.Flags().StringArrayVar(&latestLabelArgs, "label", nil, "label filter (key=value)")
	latest.Flags().StringVar(&format, "format", "table", "output format: table|json")

	var countLabelArgs []string
	var countFrom, countTo int64
	countCmd := &cobra.Command{
		Use:   "count <db> <metric>",
		Short: "Count data points in range",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now().Unix()
			if countFrom == 0 {
				countFrom = now - 3600
			}
			if countTo == 0 {
				countTo = now
			}
			n, err := ts.Count(args[0], args[1], countFrom, countTo, parseLabels(countLabelArgs))
			if err != nil {
				return err
			}
			fmt.Println(n)
			return nil
		},
	}
	countCmd.Flags().Int64Var(&countFrom, "from", 0, "start unix timestamp (default: now - 1h)")
	countCmd.Flags().Int64Var(&countTo, "to", 0, "end unix timestamp (default: now)")
	countCmd.Flags().StringArrayVar(&countLabelArgs, "label", nil, "label filter (key=value)")

	var exportLabelArgs []string
	var exportFrom, exportTo int64
	var exportOutput string
	exportCmd := &cobra.Command{
		Use:   "export <db> <metric>",
		Short: "Export data points as CSV",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now().Unix()
			if exportFrom == 0 {
				exportFrom = now - 3600
			}
			if exportTo == 0 {
				exportTo = now
			}
			w := os.Stdout
			if exportOutput != "" {
				f, err := os.Create(exportOutput)
				if err != nil {
					return err
				}
				defer f.Close()
				w = f
			}
			return ts.ExportCSV(args[0], args[1], exportFrom, exportTo, parseLabels(exportLabelArgs), w)
		},
	}
	exportCmd.Flags().Int64Var(&exportFrom, "from", 0, "start unix timestamp (default: now - 1h)")
	exportCmd.Flags().Int64Var(&exportTo, "to", 0, "end unix timestamp (default: now)")
	exportCmd.Flags().StringArrayVar(&exportLabelArgs, "label", nil, "label filter (key=value)")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "output file (default: stdout)")

	dbs := &cobra.Command{
		Use:   "dbs",
		Short: "List all TS databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := ts.ListDBs()
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
		Use:   "bulk <db> <ndjson-file>",
		Short: "Bulk write data points from NDJSON file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := ts.WriteBatchFromNDJSON(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("wrote %d data points\n", count)
			return nil
		},
	}

	deleteDB := &cobra.Command{
		Use:   "drop <db>",
		Short: "Delete a TS database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ts.DropDB(args[0]); err != nil {
				return err
			}
			fmt.Printf("dropped %s\n", args[0])
			return nil
		},
	}

	cmd.AddCommand(write, query, latest, countCmd, exportCmd, dbs, bulk, deleteDB)
	return cmd
}
