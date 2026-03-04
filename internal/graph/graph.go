package graph

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bull-cli/bull/internal/config"
	gr "github.com/dominikbraun/graph"
)

type graphData struct {
	Directed bool                `json:"directed"`
	Vertices []vertexData        `json:"vertices"`
	Edges    []edgeData          `json:"edges"`
}

type vertexData struct {
	ID    string            `json:"id"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

type edgeData struct {
	From   string            `json:"from"`
	To     string            `json:"to"`
	Weight int               `json:"weight,omitempty"`
	Attrs  map[string]string `json:"attrs,omitempty"`
}

func dbPath(name string) string {
	return filepath.Join(config.GraphDir(), name+".json")
}

func Load(name string, directed bool) (gr.Graph[string, string], error) {
	g := newGraph(directed)
	path := dbPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return g, nil
		}
		return nil, err
	}
	var gd graphData
	if err := json.Unmarshal(data, &gd); err != nil {
		return nil, err
	}
	for _, v := range gd.Vertices {
		attrs := make([]func(*gr.VertexProperties), 0)
		for k, val := range v.Attrs {
			attrs = append(attrs, gr.VertexAttribute(k, val))
		}
		g.AddVertex(v.ID, attrs...)
	}
	for _, e := range gd.Edges {
		opts := []func(*gr.EdgeProperties){gr.EdgeWeight(e.Weight)}
		for k, val := range e.Attrs {
			opts = append(opts, gr.EdgeAttribute(k, val))
		}
		g.AddEdge(e.From, e.To, opts...)
	}
	return g, nil
}

func Save(name string, g gr.Graph[string, string], directed bool) error {
	gd := graphData{Directed: directed}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return err
	}
	for v := range adjMap {
		_, props, _ := g.VertexWithProperties(v)
		vd := vertexData{ID: v}
		if props.Attributes != nil && len(props.Attributes) > 0 {
			vd.Attrs = props.Attributes
		}
		gd.Vertices = append(gd.Vertices, vd)
	}
	edges, _ := g.Edges()
	for _, e := range edges {
		ed := edgeData{
			From:   e.Source,
			To:     e.Target,
			Weight: e.Properties.Weight,
		}
		if e.Properties.Attributes != nil && len(e.Properties.Attributes) > 0 {
			ed.Attrs = e.Properties.Attributes
		}
		gd.Edges = append(gd.Edges, ed)
	}
	data, err := json.MarshalIndent(gd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbPath(name), data, 0644)
}

func newGraph(directed bool) gr.Graph[string, string] {
	if directed {
		return gr.New(gr.StringHash, gr.Directed(), gr.Weighted())
	}
	return gr.New(gr.StringHash, gr.Weighted())
}

func AddVertex(name string, directed bool, id string, attrs map[string]string) error {
	g, err := Load(name, directed)
	if err != nil {
		return err
	}
	opts := make([]func(*gr.VertexProperties), 0)
	for k, v := range attrs {
		opts = append(opts, gr.VertexAttribute(k, v))
	}
	if err := g.AddVertex(id, opts...); err != nil {
		return err
	}
	return Save(name, g, directed)
}

func AddEdge(name string, directed bool, from, to string, weight int, attrs map[string]string) error {
	g, err := Load(name, directed)
	if err != nil {
		return err
	}
	opts := []func(*gr.EdgeProperties){gr.EdgeWeight(weight)}
	for k, v := range attrs {
		opts = append(opts, gr.EdgeAttribute(k, v))
	}
	if err := g.AddEdge(from, to, opts...); err != nil {
		return err
	}
	return Save(name, g, directed)
}

func ShortestPath(name string, directed bool, from, to string) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	return gr.ShortestPath(g, from, to)
}

func DFS(name string, directed bool, start string) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	var result []string
	err = gr.DFS(g, start, func(v string) bool {
		result = append(result, v)
		return false
	})
	return result, err
}

func BFS(name string, directed bool, start string) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	var result []string
	err = gr.BFS(g, start, func(v string) bool {
		result = append(result, v)
		return false
	})
	return result, err
}

func Vertices(name string, directed bool) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	var result []string
	for v := range adjMap {
		result = append(result, v)
	}
	return result, nil
}

type EdgeInfo struct {
	From   string
	To     string
	Weight int
}

func EdgeList(name string, directed bool) ([]EdgeInfo, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	edges, err := g.Edges()
	if err != nil {
		return nil, err
	}
	var result []EdgeInfo
	for _, e := range edges {
		result = append(result, EdgeInfo{From: e.Source, To: e.Target, Weight: e.Properties.Weight})
	}
	return result, nil
}

func RemoveVertex(name string, directed bool, id string) error {
	g, err := Load(name, directed)
	if err != nil {
		return err
	}
	if err := g.RemoveVertex(id); err != nil {
		return err
	}
	return Save(name, g, directed)
}

func RemoveEdge(name string, directed bool, from, to string) error {
	g, err := Load(name, directed)
	if err != nil {
		return err
	}
	if err := g.RemoveEdge(from, to); err != nil {
		return err
	}
	return Save(name, g, directed)
}

func Neighbors(name string, directed bool, id string) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	edgesMap, ok := adjMap[id]
	if !ok {
		return nil, fmt.Errorf("vertex %q not found", id)
	}
	var result []string
	for target := range edgesMap {
		result = append(result, target)
	}
	return result, nil
}

func Degree(name string, directed bool, id string) (int, error) {
	g, err := Load(name, directed)
	if err != nil {
		return 0, err
	}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return 0, err
	}
	edgesMap, ok := adjMap[id]
	if !ok {
		return 0, fmt.Errorf("vertex %q not found", id)
	}
	return len(edgesMap), nil
}

func HasPath(name string, directed bool, from, to string) (bool, error) {
	path, err := ShortestPath(name, directed, from, to)
	if err != nil {
		return false, nil
	}
	return len(path) > 0, nil
}

type GraphStats struct {
	VertexCount int `json:"vertex_count"`
	EdgeCount   int `json:"edge_count"`
}

func Stats(name string, directed bool) (*GraphStats, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	edges, _ := g.Edges()
	return &GraphStats{
		VertexCount: len(adjMap),
		EdgeCount:   len(edges),
	}, nil
}

func VertexAttrs(name string, directed bool, id string) (map[string]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	_, props, err := g.VertexWithProperties(id)
	if err != nil {
		return nil, err
	}
	return props.Attributes, nil
}

func ConnectedComponents(name string, directed bool) ([][]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	if directed {
		return gr.StronglyConnectedComponents(g)
	}
	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}
	visited := make(map[string]bool)
	var components [][]string
	for v := range adjMap {
		if visited[v] {
			continue
		}
		var component []string
		queue := []string{v}
		visited[v] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			component = append(component, cur)
			for neighbor := range adjMap[cur] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
		components = append(components, component)
	}
	return components, nil
}

func TopologicalSort(name string, directed bool) ([]string, error) {
	g, err := Load(name, directed)
	if err != nil {
		return nil, err
	}
	return gr.TopologicalSort(g)
}

func HasCycle(name string, directed bool) (bool, error) {
	g, err := Load(name, directed)
	if err != nil {
		return false, err
	}
	_, err = gr.TopologicalSort(g)
	if err != nil {
		return true, nil
	}
	return false, nil
}

func ListDBs() ([]string, error) {
	pattern := filepath.Join(config.GraphDir(), "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		name := filepath.Base(m)
		names = append(names, name[:len(name)-5])
	}
	return names, nil
}

func DropDB(name string) error {
	return os.Remove(dbPath(name))
}

func ImportCSV(name string, directed bool, csvFile string) (int, int, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	g, err := Load(name, directed)
	if err != nil {
		return 0, 0, err
	}

	scanner := bufio.NewScanner(f)
	vCount, eCount := 0, 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		weight := 0
		if len(parts) >= 3 {
			fmt.Sscanf(strings.TrimSpace(parts[2]), "%d", &weight)
		}

		if err := g.AddVertex(from); err != nil && err != gr.ErrVertexAlreadyExists {
			return vCount, eCount, err
		} else if err == nil {
			vCount++
		}
		if err := g.AddVertex(to); err != nil && err != gr.ErrVertexAlreadyExists {
			return vCount, eCount, err
		} else if err == nil {
			vCount++
		}
		opts := []func(*gr.EdgeProperties){gr.EdgeWeight(weight)}
		if err := g.AddEdge(from, to, opts...); err != nil && err != gr.ErrEdgeAlreadyExists {
			return vCount, eCount, err
		} else if err == nil {
			eCount++
		}
	}
	if err := scanner.Err(); err != nil {
		return vCount, eCount, err
	}
	return vCount, eCount, Save(name, g, directed)
}

func ExportJSON(name string, directed bool) ([]byte, error) {
	path := dbPath(name)
	return os.ReadFile(path)
}
