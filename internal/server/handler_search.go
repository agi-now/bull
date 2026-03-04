package server

import (
	"encoding/json"
	"net/http"

	"github.com/agi-now/bull/internal/search"
)

func (s *Server) registerSearch() {
	s.mux.HandleFunc("GET /api/search/dbs", s.searchDbs)
	s.mux.HandleFunc("POST /api/search/{idx}/create", s.searchCreate)
	s.mux.HandleFunc("POST /api/search/{idx}/index", s.searchIndex)
	s.mux.HandleFunc("POST /api/search/{idx}/query", s.searchQuery)
	s.mux.HandleFunc("GET /api/search/{idx}/get/{docID}", s.searchGetDoc)
	s.mux.HandleFunc("POST /api/search/{idx}/update", s.searchUpdate)
	s.mux.HandleFunc("DELETE /api/search/{idx}/doc/{docID}", s.searchDeleteDoc)
	s.mux.HandleFunc("GET /api/search/{idx}/info", s.searchInfo)
	s.mux.HandleFunc("DELETE /api/search/{idx}", s.searchDrop)
}

func (s *Server) searchDbs(w http.ResponseWriter, r *http.Request) {
	names, err := search.ListDBs()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) searchCreate(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	if err := search.Create(idx); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) searchIndex(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	var req struct {
		ID  string      `json:"id"`
		Doc interface{} `json:"doc"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	jsonBytes, _ := json.Marshal(req.Doc)
	if err := search.Index(idx, req.ID, string(jsonBytes)); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) searchQuery(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	var req struct {
		Query  string   `json:"query"`
		Limit  int      `json:"limit"`
		Offset int      `json:"offset"`
		Fields []string `json:"fields"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	result, err := search.QueryIndexWithFields(idx, req.Query, req.Limit, req.Offset, req.Fields)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, result)
}

func (s *Server) searchGetDoc(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	docID := r.PathValue("docID")
	doc, err := search.GetDoc(idx, docID)
	if err != nil {
		fail(w, 404, err.Error())
		return
	}
	ok(w, doc)
}

func (s *Server) searchUpdate(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	var req struct {
		ID  string      `json:"id"`
		Doc interface{} `json:"doc"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	jsonBytes, _ := json.Marshal(req.Doc)
	if err := search.UpdateDoc(idx, req.ID, string(jsonBytes)); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) searchDeleteDoc(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	docID := r.PathValue("docID")
	if err := search.DeleteDoc(idx, docID); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) searchInfo(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	count, err := search.Info(idx)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, map[string]interface{}{"index": idx, "documents": count})
}

func (s *Server) searchDrop(w http.ResponseWriter, r *http.Request) {
	idx := r.PathValue("idx")
	if err := search.DropIndex(idx); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

