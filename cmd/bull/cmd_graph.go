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

	delVertex := &cobra.Command{
		Use:   "del-vertex <db> <id>",
		Short: "Remove a vertex and its edges",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bgraph.RemoveVertex(args[0], directed(), args[1])
		},
	}

	delEdge := &cobra.Command{
		Use:   "del-edge <db> <from> <to>",
		Short: "Remove an edge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bgraph.RemoveEdge(args[0], directed(), args[1], args[2])
		},
	}

	neighbors := &cobra.Command{
		Use:   "neighbors <db> <vertex>",
		Short: "List neighbors of a vertex",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.Neighbors(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			for _, v := range result {
				fmt.Println(v)
			}
			return nil
		},
	}

	degree := &cobra.Command{
		Use:   "degree <db> <vertex>",
		Short: "Get degree of a vertex",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := bgraph.Degree(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			fmt.Println(n)
			return nil
		},
	}

	hasPath := &cobra.Command{
		Use:   "has-path <db> <from> <to>",
		Short: "Check if a path exists between two vertices",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := bgraph.HasPath(args[0], directed(), args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Println(ok)
			return nil
		},
	}

	stats := &cobra.Command{
		Use:   "stats <db>",
		Short: "Show graph statistics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := bgraph.Stats(args[0], directed())
			if err != nil {
				return err
			}
			fmt.Printf("vertices: %d\nedges: %d\n", s.VertexCount, s.EdgeCount)
			return nil
		},
	}

	vertexAttrs := &cobra.Command{
		Use:   "attrs <db> <vertex>",
		Short: "Show vertex attributes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			attrs, err := bgraph.VertexAttrs(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			for k, v := range attrs {
				fmt.Printf("%s=%s\n", k, v)
			}
			return nil
		},
	}

	components := &cobra.Command{
		Use:   "components <db>",
		Short: "Find connected components (SCC for directed)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := bgraph.ConnectedComponents(args[0], directed())
			if err != nil {
				return err
			}
			for i, c := range cc {
				fmt.Printf("Component %d: %s\n", i+1, strings.Join(c, ", "))
			}
			return nil
		},
	}

	toposort := &cobra.Command{
		Use:   "toposort <db>",
		Short: "Topological sort (directed graphs)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := bgraph.TopologicalSort(args[0], directed())
			if err != nil {
				return err
			}
			for _, v := range result {
				fmt.Println(v)
			}
			return nil
		},
	}

	hasCycle := &cobra.Command{
		Use:   "has-cycle <db>",
		Short: "Check if the graph has a cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := bgraph.HasCycle(args[0], directed())
			if err != nil {
				return err
			}
			fmt.Println(ok)
			return nil
		},
	}

	importCSV := &cobra.Command{
		Use:   "import-csv <db> <file.csv>",
		Short: "Import edges from CSV (from,to[,weight])",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, e, err := bgraph.ImportCSV(args[0], directed(), args[1])
			if err != nil {
				return err
			}
			fmt.Printf("added %d vertices, %d edges\n", v, e)
			return nil
		},
	}

	exportJSON := &cobra.Command{
		Use:   "export <db>",
		Short: "Export graph as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := bgraph.ExportJSON(args[0], directed())
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		},
	}

	dropDB := &cobra.Command{
		Use:   "drop <db>",
		Short: "Delete a graph",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bgraph.DropDB(args[0]); err != nil {
				return err
			}
			fmt.Printf("dropped %s\n", args[0])
			return nil
		},
	}

	allCmds := []*cobra.Command{addVertex, addEdge, shortestPath, dfs, bfs, vertices, edges,
		delVertex, delEdge, neighbors, degree, hasPath, stats, vertexAttrs, components, toposort, hasCycle, importCSV, exportJSON, dropDB}
	for _, c := range allCmds {
		c.Flags().BoolVar(&undirected, "undirected", false, "use undirected graph")
	}

	cmd.AddCommand(allCmds...)
	cmd.AddCommand(dbs)
	return cmd
}
