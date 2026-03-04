package graph

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/bull-cli/bull/internal/config"
)

func setup(t *testing.T) {
	t.Helper()
	config.DataDir = t.TempDir()
}

func buildDAG(t *testing.T) {
	t.Helper()
	AddVertex("g", true, "A", nil)
	AddVertex("g", true, "B", nil)
	AddVertex("g", true, "C", nil)
	AddVertex("g", true, "D", nil)
	AddEdge("g", true, "A", "B", 1, nil)
	AddEdge("g", true, "A", "C", 2, nil)
	AddEdge("g", true, "B", "D", 1, nil)
	AddEdge("g", true, "C", "D", 1, nil)
}

func TestAddVertexAndEdge(t *testing.T) {
	setup(t)
	if err := AddVertex("g", true, "A", map[string]string{"type": "svc"}); err != nil {
		t.Fatal(err)
	}
	if err := AddVertex("g", true, "B", nil); err != nil {
		t.Fatal(err)
	}
	if err := AddEdge("g", true, "A", "B", 5, nil); err != nil {
		t.Fatal(err)
	}

	verts, err := Vertices("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(verts) != 2 {
		t.Fatalf("expected 2 vertices, got %d", len(verts))
	}

	edges, err := EdgeList("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Weight != 5 {
		t.Fatalf("expected weight 5, got %d", edges[0].Weight)
	}
}

func TestRemoveEdge(t *testing.T) {
	setup(t)
	buildDAG(t)

	if err := RemoveEdge("g", true, "A", "B"); err != nil {
		t.Fatal(err)
	}
	edges, _ := EdgeList("g", true)
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges after removing one, got %d", len(edges))
	}
}

func TestRemoveVertex(t *testing.T) {
	setup(t)
	AddVertex("g", true, "X", nil)
	AddVertex("g", true, "Y", nil)

	if err := RemoveVertex("g", true, "X"); err != nil {
		t.Fatal(err)
	}
	verts, _ := Vertices("g", true)
	if len(verts) != 1 {
		t.Fatalf("expected 1 vertex, got %d", len(verts))
	}
}

func TestShortestPath(t *testing.T) {
	setup(t)
	buildDAG(t)

	path, err := ShortestPath("g", true, "A", "D")
	if err != nil {
		t.Fatal(err)
	}
	if len(path) < 2 {
		t.Fatalf("expected a path from A to D, got %v", path)
	}
	if path[0] != "A" || path[len(path)-1] != "D" {
		t.Fatalf("path should start with A and end with D, got %v", path)
	}
}

func TestDFSBFS(t *testing.T) {
	setup(t)
	buildDAG(t)

	dfsResult, err := DFS("g", true, "A")
	if err != nil {
		t.Fatal(err)
	}
	if len(dfsResult) != 4 {
		t.Fatalf("DFS expected 4 vertices, got %d", len(dfsResult))
	}

	bfsResult, err := BFS("g", true, "A")
	if err != nil {
		t.Fatal(err)
	}
	if len(bfsResult) != 4 {
		t.Fatalf("BFS expected 4 vertices, got %d", len(bfsResult))
	}
}

func TestNeighborsAndDegree(t *testing.T) {
	setup(t)
	buildDAG(t)

	neighbors, err := Neighbors("g", true, "A")
	if err != nil {
		t.Fatal(err)
	}
	if len(neighbors) != 2 {
		t.Fatalf("expected 2 neighbors for A, got %d", len(neighbors))
	}

	deg, err := Degree("g", true, "A")
	if err != nil {
		t.Fatal(err)
	}
	if deg != 2 {
		t.Fatalf("expected degree 2, got %d", deg)
	}
}

func TestHasPath(t *testing.T) {
	setup(t)
	buildDAG(t)

	ok, err := HasPath("g", true, "A", "D")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected path from A to D")
	}

	ok, _ = HasPath("g", true, "D", "A")
	if ok {
		t.Fatal("expected no path from D to A in directed graph")
	}
}

func TestVertexAttrs(t *testing.T) {
	setup(t)
	AddVertex("g", true, "X", map[string]string{"color": "red", "size": "large"})

	attrs, err := VertexAttrs("g", true, "X")
	if err != nil {
		t.Fatal(err)
	}
	if attrs["color"] != "red" || attrs["size"] != "large" {
		t.Fatalf("unexpected attrs: %v", attrs)
	}
}

func TestStats(t *testing.T) {
	setup(t)
	buildDAG(t)

	s, err := Stats("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if s.VertexCount != 4 {
		t.Fatalf("expected 4 vertices, got %d", s.VertexCount)
	}
	if s.EdgeCount != 4 {
		t.Fatalf("expected 4 edges, got %d", s.EdgeCount)
	}
}

func TestTopologicalSort(t *testing.T) {
	setup(t)
	buildDAG(t)

	order, err := TopologicalSort("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 vertices in toposort, got %d", len(order))
	}
	posA, posD := -1, -1
	for i, v := range order {
		if v == "A" {
			posA = i
		}
		if v == "D" {
			posD = i
		}
	}
	if posA > posD {
		t.Fatal("A should come before D in topological order")
	}
}

func TestHasCycle(t *testing.T) {
	setup(t)
	buildDAG(t)

	has, err := HasCycle("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Fatal("DAG should not have a cycle")
	}

	AddEdge("g2", true, "X", "Y", 0, nil)
	AddVertex("g2", true, "X", nil)
	AddVertex("g2", true, "Y", nil)
	// Build a cycle: manually create circular edges
	AddVertex("cyc", true, "A", nil)
	AddVertex("cyc", true, "B", nil)
	AddEdge("cyc", true, "A", "B", 0, nil)
	AddEdge("cyc", true, "B", "A", 0, nil)

	has, err = HasCycle("cyc", true)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal("graph with A->B->A should have a cycle")
	}
}

func TestConnectedComponents(t *testing.T) {
	setup(t)
	AddVertex("g", false, "A", nil)
	AddVertex("g", false, "B", nil)
	AddVertex("g", false, "C", nil)
	AddVertex("g", false, "D", nil)
	AddEdge("g", false, "A", "B", 0, nil)
	AddEdge("g", false, "C", "D", 0, nil)

	cc, err := ConnectedComponents("g", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(cc) != 2 {
		t.Fatalf("expected 2 connected components, got %d", len(cc))
	}
}

func TestImportCSVAndExport(t *testing.T) {
	setup(t)
	csvContent := "# comment\nA,B,1\nB,C,2\nC,D,3\n"
	csvPath := filepath.Join(t.TempDir(), "edges.csv")
	os.WriteFile(csvPath, []byte(csvContent), 0644)

	v, e, err := ImportCSV("g", true, csvPath)
	if err != nil {
		t.Fatal(err)
	}
	if v != 4 {
		t.Fatalf("expected 4 vertices, got %d", v)
	}
	if e != 3 {
		t.Fatalf("expected 3 edges, got %d", e)
	}

	data, err := ExportJSON("g", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("exported JSON should not be empty")
	}
}

func TestDropDBAndListDBs(t *testing.T) {
	setup(t)
	AddVertex("g1", true, "A", nil)
	AddVertex("g2", true, "B", nil)

	names, err := ListDBs()
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(names)
	if len(names) != 2 {
		t.Fatalf("expected 2 graphs, got %d", len(names))
	}

	if err := DropDB("g1"); err != nil {
		t.Fatal(err)
	}
	names, _ = ListDBs()
	if len(names) != 1 {
		t.Fatalf("expected 1 graph after drop, got %d", len(names))
	}
}
