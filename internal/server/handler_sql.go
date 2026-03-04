package server

import (
	"net/http"
	"strconv"

	bsql "github.com/agi-now/bull/internal/sql"
)

func (s *Server) registerSQL() {
	s.mux.HandleFunc("GET /api/sql/dbs", s.sqlDbs)
	s.mux.HandleFunc("POST /api/sql/{db}/exec", s.sqlExec)
	s.mux.HandleFunc("POST /api/sql/{db}/query", s.sqlQuery)
	s.mux.HandleFunc("GET /api/sql/{db}/tables", s.sqlTables)
	s.mux.HandleFunc("GET /api/sql/{db}/schema/{table}", s.sqlSchema)
	s.mux.HandleFunc("GET /api/sql/{db}/describe/{table}", s.sqlDescribe)
	s.mux.HandleFunc("GET /api/sql/{db}/count/{table}", s.sqlCount)
	s.mux.HandleFunc("DELETE /api/sql/{db}", s.sqlDrop)
}

func (s *Server) sqlDbs(w http.ResponseWriter, r *http.Request) {
	names, err := bsql.ListDBs()
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, names)
}

func (s *Server) sqlExec(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		SQL string `json:"sql"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	affected, err := bsql.Exec(db, req.SQL)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, map[string]int64{"rows_affected": affected})
}

func (s *Server) sqlQuery(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	var req struct {
		SQL   string `json:"sql"`
		Limit int    `json:"limit"`
	}
	if err := bind(r, &req); err != nil {
		fail(w, 400, err.Error())
		return
	}
	sqlStr := req.SQL
	if req.Limit > 0 {
		sqlStr = "SELECT * FROM (" + sqlStr + ") LIMIT " + strconv.Itoa(req.Limit)
	}
	result, err := bsql.Query(db, sqlStr)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	rows := make([]map[string]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		m := make(map[string]string, len(result.Columns))
		for i, c := range result.Columns {
			m[c] = row[i]
		}
		rows = append(rows, m)
	}
	ok(w, map[string]interface{}{
		"columns": result.Columns,
		"rows":    rows,
		"count":   len(result.Rows),
	})
}

func (s *Server) sqlTables(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	tables, err := bsql.Tables(db)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, tables)
}

func (s *Server) sqlSchema(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	ddl, err := bsql.Schema(db, table)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, ddl)
}

func (s *Server) sqlDescribe(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	cols, err := bsql.Describe(db, table)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, cols)
}

func (s *Server) sqlCount(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	table := r.PathValue("table")
	n, err := bsql.CountRows(db, table)
	if err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, n)
}

func (s *Server) sqlDrop(w http.ResponseWriter, r *http.Request) {
	db := r.PathValue("db")
	if err := bsql.DropDB(db); err != nil {
		fail(w, 500, err.Error())
		return
	}
	ok(w, nil)
}

