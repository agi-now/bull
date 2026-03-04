package server

import (
	"net/http"

	bgraph "github.com/bull-cli/bull/internal/graph"
)

func (s *Server) registerGraph() {
	s.mux.HandleFunc("GET /api/graph/dbs", s.graphDbs)
	s.mux.HandleFunc("POST /api/graph/{db}/add-vertex", s.graphAddVertex)
	s.mux.HandleFunc("POST /api/graph/{db}/add-edge", s.graphAddEdge)
	s.mux.HandleFunc("POST /api/graph/{db}/del-vertex", s.graphDelVertex)
	s.mux.HandleFunc("POST /api/graph/{db}/del-edge", s.graphDelEdge)
	s.mux.HandleFunc("POST /api/graph/{db}/vertices", s.graphVertices)
	s.mux.HandleFunc("POST /api/graph/{db}/edges", s.graphEdges)
	s.mux.HandleFunc("POST /api/graph/{db}/neighbors", s.graphNeighbors)
	s.mux.HandleFunc("POST /api/graph/{db}/degree", s.graphDegree)
	s.mux.HandleFunc("POST /api/graph/{db}/attrs", s.graphAttrs)
	s.mux.HandleFunc("POST /api/graph/{db}/shortest-path", s.graphShortestPath)
	s.mux.HandleFunc("POST /api/graph/{db}/has-path", s.graphHasPath)
	s.mux.HandleFunc("POST /api/graph/{db}/dfs", s.graphDFS)
	s.mux.HandleFunc("POST /api/graph/{db}/bfs", s.graphBFS)
	s.mux.HandleFunc("POST /api/graph/{db}/stats", s.graphStats)
	s.mux.HandleFunc("POST /api/graph/{db}/components", s.graphComponents)
	s.mux.HandleFunc("POST /api/graph/{db}/toposort", s.graphToposort)
	s.mux.HandleFunc("POST /api/graph/{db}/has-cycle", s.graphHasCycle)
	s.mux.HandleFunc("POST /api/graph/{db}/export", s.graphExport)
	s.mux.HandleFunc("DELETE /api/graph/{db}", s.graphDrop)
}

type graphReq struct {
	ID         string            `json:"id"`
	From       string            `json:"from"`
	To         string            `json:"to"`
	Start      string            `json:"start"`
	Weight     int               `json:"weight"`
	Attrs      map[string]string `json:"attrs"`
	Undirected bool              `json:"undirected"`
}

func (g *graphReq) directed() bool { return !g.Undirected }

func (s *Server) graphDbs(w http.ResponseWriter, r *http.Request) {
	names, err := bgraph.ListDBs()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) graphAddVertex(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := bgraph.AddVertex(db, req.directed(), req.ID, req.Attrs); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) graphAddEdge(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := bgraph.AddEdge(db, req.directed(), req.From, req.To, req.Weight, req.Attrs); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) graphDelVertex(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := bgraph.RemoveVertex(db, req.directed(), req.ID); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) graphDelEdge(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := bgraph.RemoveEdge(db, req.directed(), req.From, req.To); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) graphVertices(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	verts, err := bgraph.Vertices(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, verts)
}

func (s *Server) graphEdges(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	edges, err := bgraph.EdgeList(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, edges)
}

func (s *Server) graphNeighbors(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	result, err := bgraph.Neighbors(db, req.directed(), req.ID)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, result)
}

func (s *Server) graphDegree(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	n, err := bgraph.Degree(db, req.directed(), req.ID)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, n)
}

func (s *Server) graphAttrs(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	attrs, err := bgraph.VertexAttrs(db, req.directed(), req.ID)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, attrs)
}

func (s *Server) graphShortestPath(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	path, err := bgraph.ShortestPath(db, req.directed(), req.From, req.To)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, path)
}

func (s *Server) graphHasPath(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	has, err := bgraph.HasPath(db, req.directed(), req.From, req.To)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, has)
}

func (s *Server) graphDFS(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	result, err := bgraph.DFS(db, req.directed(), req.Start)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, result)
}

func (s *Server) graphBFS(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	result, err := bgraph.BFS(db, req.directed(), req.Start)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, result)
}

func (s *Server) graphStats(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	stats, err := bgraph.Stats(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, stats)
}

func (s *Server) graphComponents(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	cc, err := bgraph.ConnectedComponents(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, cc)
}

func (s *Server) graphToposort(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	order, err := bgraph.TopologicalSort(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, order)
}

func (s *Server) graphHasCycle(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	has, err := bgraph.HasCycle(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, has)
}

func (s *Server) graphExport(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req graphReq
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	data, err := bgraph.ExportJSON(db, req.directed())
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) graphDrop(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	if err := bgraph.DropDB(db); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}
