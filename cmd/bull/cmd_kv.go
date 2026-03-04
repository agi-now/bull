package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bull-cli/bull/internal/kv"
	"github.com/spf13/cobra"
)

func kvCmd() *cobra.Command {
	var bucket string

	cmd := &cobra.Command{
		Use:   "kv",
		Short: "KV store operations (bbolt)",
	}

	put := &cobra.Command{
		Use:   "put <db> <key> <value>",
		Short: "Put a key-value pair",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return kv.Put(args[0], bucket, args[1], args[2])
		},
	}

	get := &cobra.Command{
		Use:   "get <db> <key>",
		Short: "Get value by key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := kv.Get(args[0], bucket, args[1])
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}

	del := &cobra.Command{
		Use:   "del <db> <key>",
		Short: "Delete a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return kv.Del(args[0], bucket, args[1])
		},
	}

	var prefix string
	var listFormat string
	list := &cobra.Command{
		Use:   "list <db>",
		Short: "List keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := kv.List(args[0], bucket, prefix)
			if err != nil {
				return err
			}
			switch listFormat {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(pairs)
			default:
				for _, p := range pairs {
					fmt.Printf("%s\t%s\n", p.Key, p.Value)
				}
			}
			return nil
		},
	}
	list.Flags().StringVar(&prefix, "prefix", "", "key prefix filter")
	list.Flags().StringVar(&listFormat, "format", "tsv", "output format: tsv|json")

	dbs := &cobra.Command{
		Use:   "dbs",
		Short: "List all KV databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := kv.ListDBs()
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	buckets := &cobra.Command{
		Use:   "buckets <db>",
		Short: "List buckets in a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := kv.ListBuckets(args[0])
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	count := &cobra.Command{
		Use:   "count <db>",
		Short: "Count keys in a bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := kv.Count(args[0], bucket)
			if err != nil {
				return err
			}
			fmt.Println(n)
			return nil
		},
	}

	exists := &cobra.Command{
		Use:   "exists <db> <key>",
		Short: "Check if a key exists",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := kv.Exists(args[0], bucket, args[1])
			if err != nil {
				return err
			}
			if ok {
				fmt.Println("true")
			} else {
				fmt.Println("false")
			}
			return nil
		},
	}

	incr := &cobra.Command{
		Use:   "incr <db> <key> [delta]",
		Short: "Increment a numeric key (default delta=1)",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var delta int64 = 1
			if len(args) == 3 {
				fmt.Sscanf(args[2], "%d", &delta)
			}
			val, err := kv.Incr(args[0], bucket, args[1], delta)
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}

	decr := &cobra.Command{
		Use:   "decr <db> <key> [delta]",
		Short: "Decrement a numeric key (default delta=1)",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var delta int64 = 1
			if len(args) == 3 {
				fmt.Sscanf(args[2], "%d", &delta)
			}
			val, err := kv.Incr(args[0], bucket, args[1], -delta)
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}

	exportJSON := &cobra.Command{
		Use:   "export <db>",
		Short: "Export bucket data as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := kv.ExportJSON(args[0], bucket)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(pairs)
		},
	}

	var importFile string
	importJSON := &cobra.Command{
		Use:   "import <db> -f <file.json>",
		Short: "Import key-value pairs from JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(importFile)
			if err != nil {
				return err
			}
			var pairs []kv.KVPair
			if err := json.Unmarshal(data, &pairs); err != nil {
				return err
			}
			if err := kv.ImportJSON(args[0], bucket, pairs); err != nil {
				return err
			}
			fmt.Printf("imported %d pairs\n", len(pairs))
			return nil
		},
	}
	importJSON.Flags().StringVarP(&importFile, "file", "f", "", "JSON file path")
	importJSON.MarkFlagRequired("file")

	mget := &cobra.Command{
		Use:   "mget <db> <key1> <key2> ...",
		Short: "Get multiple keys at once",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := kv.MGet(args[0], bucket, args[1:])
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(pairs)
		},
	}

	mput := &cobra.Command{
		Use:   "mput <db> <key1> <val1> <key2> <val2> ...",
		Short: "Put multiple key-value pairs at once",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			kvArgs := args[1:]
			if len(kvArgs)%2 != 0 {
				return fmt.Errorf("expected even number of key-value arguments, got %d", len(kvArgs))
			}
			var pairs []kv.KVPair
			for i := 0; i < len(kvArgs); i += 2 {
				pairs = append(pairs, kv.KVPair{Key: kvArgs[i], Value: kvArgs[i+1]})
			}
			if err := kv.MPut(args[0], bucket, pairs); err != nil {
				return err
			}
			fmt.Printf("put %d pairs\n", len(pairs))
			return nil
		},
	}

	var startKey, endKey string
	var scanFormat string
	scan := &cobra.Command{
		Use:   "scan <db>",
		Short: "Range scan keys [start, end)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := kv.Scan(args[0], bucket, startKey, endKey)
			if err != nil {
				return err
			}
			switch scanFormat {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(pairs)
			default:
				for _, p := range pairs {
					fmt.Printf("%s\t%s\n", p.Key, p.Value)
				}
			}
			return nil
		},
	}
	scan.Flags().StringVar(&startKey, "start", "", "start key (inclusive)")
	scan.Flags().StringVar(&endKey, "end", "", "end key (exclusive)")
	scan.Flags().StringVar(&scanFormat, "format", "tsv", "output format: tsv|json")

	dropDB := &cobra.Command{
		Use:   "drop <db>",
		Short: "Delete an entire KV database file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := kv.DropDB(args[0]); err != nil {
				return err
			}
			fmt.Printf("dropped %s\n", args[0])
			return nil
		},
	}

	dropBucket := &cobra.Command{
		Use:   "drop-bucket <db> <bucket>",
		Short: "Delete a bucket from a database",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := kv.DropBucket(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("dropped bucket %s\n", args[1])
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&bucket, "bucket", "", "bucket name (default: \"default\")")
	cmd.AddCommand(put, get, del, list, dbs, buckets, count, exists, incr, decr, exportJSON, importJSON, scan, dropDB, dropBucket, mget, mput)
	return cmd
}
