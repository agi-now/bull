package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/agi-now/bull/internal/graph"
	"github.com/agi-now/bull/internal/kv"
	"github.com/agi-now/bull/internal/search"
	bsql "github.com/agi-now/bull/internal/sql"
	"github.com/agi-now/bull/internal/ts"
)

type Server struct {
	Version   string
	BuildTime string
	mux       *http.ServeMux
}

func New(version, buildTime string) *Server {
	s := &Server{Version: version, BuildTime: buildTime, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	log.Printf("bull http server listening on %s", addr)
	return http.ListenAndServe(addr, s)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/version", s.handleVersion)
	s.mux.HandleFunc("GET /api/info", s.handleInfo)
	s.registerKV()
	s.registerSQL()
	s.registerGraph()
	s.registerSearch()
	s.registerTS()
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	ok(w, map[string]string{
		"version":   s.Version,
		"build":     s.BuildTime,
		"go":        runtime.Version(),
		"os_arch":   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	kvDBs, _ := kv.ListDBs()
	sqlDBs, _ := bsql.ListDBs()
	graphDBs, _ := graph.ListDBs()
	searchDBs, _ := search.ListDBs()
	tsDBs, _ := ts.ListDBs()

	ok(w, map[string]interface{}{
		"kv":     kvDBs,
		"sql":    sqlDBs,
		"graph":  graphDBs,
		"search": searchDBs,
		"ts":     tsDBs,
	})
}

// --- JSON helpers ---

func ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "data": data})
}

func fail(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": msg})
}

func bind(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
