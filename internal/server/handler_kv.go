package server

import (
	"net/http"

	"github.com/agi-now/bull/internal/kv"
)

func (s *Server) registerKV() {
	s.mux.HandleFunc("GET /api/kv/dbs", s.kvDbs)
	s.mux.HandleFunc("POST /api/kv/{db}/put", s.kvPut)
	s.mux.HandleFunc("POST /api/kv/{db}/get", s.kvGet)
	s.mux.HandleFunc("POST /api/kv/{db}/del", s.kvDel)
	s.mux.HandleFunc("POST /api/kv/{db}/mget", s.kvMGet)
	s.mux.HandleFunc("POST /api/kv/{db}/mput", s.kvMPut)
	s.mux.HandleFunc("POST /api/kv/{db}/list", s.kvList)
	s.mux.HandleFunc("POST /api/kv/{db}/scan", s.kvScan)
	s.mux.HandleFunc("POST /api/kv/{db}/exists", s.kvExists)
	s.mux.HandleFunc("POST /api/kv/{db}/count", s.kvCount)
	s.mux.HandleFunc("POST /api/kv/{db}/incr", s.kvIncr)
	s.mux.HandleFunc("GET /api/kv/{db}/buckets", s.kvBuckets)
	s.mux.HandleFunc("POST /api/kv/{db}/export", s.kvExport)
	s.mux.HandleFunc("POST /api/kv/{db}/import", s.kvImport)
	s.mux.HandleFunc("DELETE /api/kv/{db}", s.kvDrop)
	s.mux.HandleFunc("POST /api/kv/{db}/drop-bucket", s.kvDropBucket)
}

func (s *Server) kvDbs(w http.ResponseWriter, r *http.Request) {
	names, err := kv.ListDBs()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) kvPut(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := kv.Put(db, req.Bucket, req.Key, req.Value); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) kvGet(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Key    string `json:"key"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	val, err := kv.Get(db, req.Bucket, req.Key)
	if err != nil {
		fail(w, 404, err.Error())
		return
	}
	ok(w, val)
}

func (s *Server) kvDel(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Key    string `json:"key"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := kv.Del(db, req.Bucket, req.Key); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) kvMGet(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Keys   []string `json:"keys"`
		Bucket string   `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	pairs, err := kv.MGet(db, req.Bucket, req.Keys)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, pairs)
}

func (s *Server) kvMPut(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Pairs  []kv.KVPair `json:"pairs"`
		Bucket string      `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := kv.MPut(db, req.Bucket, req.Pairs); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, len(req.Pairs))
}

func (s *Server) kvList(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Prefix string `json:"prefix"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	pairs, err := kv.List(db, req.Bucket, req.Prefix)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, pairs)
}

func (s *Server) kvScan(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Start  string `json:"start"`
		End    string `json:"end"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	pairs, err := kv.Scan(db, req.Bucket, req.Start, req.End)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, pairs)
}

func (s *Server) kvExists(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Key    string `json:"key"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	found, err := kv.Exists(db, req.Bucket, req.Key)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, found)
}

func (s *Server) kvCount(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	n, err := kv.Count(db, req.Bucket)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, n)
}

func (s *Server) kvIncr(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Key    string `json:"key"`
		Delta  int64  `json:"delta"`
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if req.Delta == 0 {
		req.Delta = 1
	}
	val, err := kv.Incr(db, req.Bucket, req.Key, req.Delta)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, val)
}

func (s *Server) kvBuckets(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	names, err := kv.ListBuckets(db)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) kvExport(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	pairs, err := kv.ExportJSON(db, req.Bucket)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, pairs)
}

func (s *Server) kvImport(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Pairs  []kv.KVPair `json:"pairs"`
		Bucket string      `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := kv.ImportJSON(db, req.Bucket, req.Pairs); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, len(req.Pairs))
}

func (s *Server) kvDrop(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	if err := kv.DropDB(db); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

func (s *Server) kvDropBucket(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		Bucket string `json:"bucket"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	if err := kv.DropBucket(db, req.Bucket); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}
