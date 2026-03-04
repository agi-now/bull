package main

import (
	"fmt"
	"strings"

	bgraph "github.com/bull-cli/bull/internal/graph"
	"github.com/spf13/cobra"
)

func graphCmd() *cobra.Command {
	var undirected bool

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Graph operations",
	}

	directed := func() bool { return !undirected }

	var attrArgs []string

	parseAttrs := func(args []string) map[string]string {
		m := make(map[string]string)
		for _, a := range args {
			parts := strings.SplitN(a, "=", 2)
			if len(parts) == 2 {
				m[parts[0]] = parts[1]
			}
		}
		return m
	}

	addVertex := &cobra.Command{
		Use:   "add-vertex <db> <id>",
		Short: "Add a vertex",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bgraph.AddVertex(args[0], directed(), args[1], parseAttrs(attrArgs))
		},
	}
	addVertex.Flags().StringArrayVar(&attrArgs, "attr", nil, "vertex attribute (key=value)")

	var weight int
	var edgeAttrArgs []string
	addEdge := &cobra.Command{
		Use:   "add-edge <db> <from> <to>",
		Short: "Add an edge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bgraph.AddEdge(args[0], directed(), args[1], args[2], weight, parseAttrs(edgeAttrArgs))
		},
	}
	addEdge.Flags().IntVar(&weight, "weight", 0, "edge weight")
	addEdge.Flags().StringArrayVar(&edgeAttrArgs, "attr", nil, "edge attribute (key=value)")

	shortestPath := &cobra.Command{
		Use:   "shortest-path <db> <from> <to>",
		Short: "Find shortest path",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := bgraph.ShortestPath(args[0], directed(), args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(path, " -> "))
			return nil
		},
	}

	dfs := &cobra.Command{
		Use:   "dfs <db> <start>",
		Short: "DFS traversal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.DFS(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			for _, v := range result {
				fmt.Println(v)
			}
			return nil
		},
	}

	bfs := &cobra.Command{
		Use:   "bfs <db> <start>",
		Short: "BFS traversal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.BFS(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			for _, v := range result {
				fmt.Println(v)
			}
			return nil
		},
	}

	vertices := &cobra.Command{
		Use:   "vertices <db>",
		Short: "List all vertices",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.Vertices(args[0], directed())
			if err != nil {
				return err
			}
			for _, v := range result {
				fmt.Println(v)
			}
			return nil
		},
	}

	edges := &cobra.Command{
		Use:   "edges <db>",
		Short: "List all edges",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.EdgeList(args[0], directed())
			if err != nil {
				return err
			}
			for _, e := range result {
				if e.Weight != 0 {
					fmt.Printf("%s -> %s (weight: %d)\n", e.From, e.To, e.Weight)
				} else {
					fmt.Printf("%s -> %s\n", e.From, e.To)
				}
			}
			return nil
		},
	}

	dbs := &cobra.Command{
		Use:   "dbs",
		Short: "List all graphs",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := bgraph.ListDBs()
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}

	for _, c := range []*cobra.Command{addVertex, addEdge, shortestPath, dfs, bfs, vertices, edges} {
		c.Flags().BoolVar(&undirected, "undirected", false, "use undirected graph")
	}

	cmd.AddCommand(addVertex, addEdge, shortestPath, dfs, bfs, vertices, edges, dbs)
	return cmd
}
