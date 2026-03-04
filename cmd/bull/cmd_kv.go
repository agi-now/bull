package main

import (
	"fmt"

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
	list := &cobra.Command{
		Use:   "list <db>",
		Short: "List keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := kv.List(args[0], bucket, prefix)
			if err != nil {
				return err
			}
			for _, p := range pairs {
				fmt.Printf("%s\t%s\n", p.Key, p.Value)
			}
			return nil
		},
	}
	list.Flags().StringVar(&prefix, "prefix", "", "key prefix filter")

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

	for _, c := range []*cobra.Command{put, get, del, list} {
		c.Flags().StringVar(&bucket, "bucket", "", "bucket name (default: \"default\")")
	}

	cmd.AddCommand(put, get, del, list, dbs)
	return cmd
}
