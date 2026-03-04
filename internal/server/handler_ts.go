package server

import (
	"bytes"
	"net/http"

	"github.com/agi-now/bull/internal/ts"
)

func (s *Server) registerTS() {
	s.mux.HandleFunc("GET /api/ts/dbs", s.tsDbs)
	s.mux.HandleFunc("POST /api/ts/{db}/write", s.tsWrite)
	s.mux.HandleFunc("POST /api/ts/{db}/query", s.tsQuery)
	s.mux.HandleFunc("POST /api/ts/{db}/latest", s.tsLatest)
	s.mux.HandleFunc("POST /api/ts/{db}/count", s.tsCount)
	s.mux.HandleFunc("POST /api/ts/{db}/export", s.tsExport)
	s.mux.HandleFunc("DELETE /api/ts/{db}", s.tsDrop)
}

func (s *Server) tsDbs(w http.ResponseWriter, r *http.Request) {
	names, err := ts.ListDBs()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) tsWrite(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Metric    string            `json:"metric"`
		Value     float64           `json:"value"`
		Timestamp int64             `json:"timestamp"`
		Labels    map[string]string `json:"labels"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := ts.Write(db, req.Metric, req.Value, req.Timestamp, req.Labels); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) tsQuery(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Metric string            `json:"metric"`
		From   int64             `json:"from"`
		To     int64             `json:"to"`
		Labels map[string]string `json:"labels"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	points, err := ts.QueryRange(db, req.Metric, req.From, req.To, req.Labels)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, points)
}

func (s *Server) tsLatest(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Metric string            `json:"metric"`
		Labels map[string]string `json:"labels"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	dp, err := ts.Latest(db, req.Metric, req.Labels)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, dp)
}

func (s *Server) tsCount(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Metric string            `json:"metric"`
		From   int64             `json:"from"`
		To     int64             `json:"to"`
		Labels map[string]string `json:"labels"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	n, err := ts.Count(db, req.Metric, req.From, req.To, req.Labels)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, n)
}

func (s *Server) tsExport(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Metric string            `json:"metric"`
		From   int64             `json:"from"`
		To     int64             `json:"to"`
		Labels map[string]string `json:"labels"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	var buf bytes.Buffer
	if err := ts.ExportCSV(db, req.Metric, req.From, req.To, req.Labels, &buf); err != nil {
		fail(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Write(buf.Bytes())
}

func (s *Server) tsDrop(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	if err := ts.DropDB(db); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}
