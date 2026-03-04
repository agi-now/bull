package sql

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bull-cli/bull/internal/config"
	_ "modernc.org/sqlite"
)

func dbPath(name string) string {
	if name == ":memory:" {
		return ":memory:"
	}
	return filepath.Join(config.SQLDir(), name+".db")
}

func OpenDB(name string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath(name))
}

func Exec(dbName, sqlStr string) (int64, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type QueryResult struct {
	Columns []string
	Rows    [][]string
}

func Query(dbName, sqlStr string) (*QueryResult, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := &QueryResult{Columns: cols}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			if v == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		result.Rows = append(result.Rows, row)
	}
	return result, rows.Err()
}

func (r *QueryResult) FormatTable(w io.Writer) {
	widths := make([]int, len(r.Columns))
	for i, c := range r.Columns {
		widths[i] = len(c)
	}
	for _, row := range r.Rows {
		for i, v := range row {
			if len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}
	sep := "+"
	for _, w := range widths {
		sep += strings.Repeat("-", w+2) + "+"
	}
	fmt.Fprintln(w, sep)
	fmt.Fprint(w, "|")
	for i, c := range r.Columns {
		fmt.Fprintf(w, " %-*s |", widths[i], c)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, sep)
	for _, row := range r.Rows {
		fmt.Fprint(w, "|")
		for i, v := range row {
			fmt.Fprintf(w, " %-*s |", widths[i], v)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, sep)
}

func (r *QueryResult) FormatCSV(w io.Writer) {
	cw := csv.NewWriter(w)
	cw.Write(r.Columns)
	for _, row := range r.Rows {
		cw.Write(row)
	}
	cw.Flush()
}

func (r *QueryResult) FormatJSON(w io.Writer) {
	var out []map[string]string
	for _, row := range r.Rows {
		m := make(map[string]string)
		for i, c := range r.Columns {
			m[c] = row[i]
		}
		out = append(out, m)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

func ImportCSV(dbName, table, csvFile string) (int, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		return 0, err
	}

	db, err := OpenDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	colDefs := make([]string, len(headers))
	for i, h := range headers {
		colDefs[i] = fmt.Sprintf(`"%s" TEXT`, h)
	}
	createSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, table, strings.Join(colDefs, ", "))
	if _, err := db.Exec(createSQL); err != nil {
		return 0, err
	}

	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	insertSQL := fmt.Sprintf(`INSERT INTO "%s" VALUES (%s)`, table, strings.Join(placeholders, ", "))

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return count, err
		}
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}
		if _, err := stmt.Exec(args...); err != nil {
			tx.Rollback()
			return count, err
		}
		count++
	}
	return count, tx.Commit()
}

func Tables(dbName string) ([]string, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func Schema(dbName, table string) (string, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var ddl string
	err = db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&ddl)
	return ddl, err
}

func CountRows(dbName, table string) (int64, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var n int64
	err = db.QueryRow(fmt.Sprintf(`SELECT count(*) FROM "%s"`, table)).Scan(&n)
	return n, err
}

func ExecFile(dbName, sqlFile string) error {
	data, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}
	db, err := OpenDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(data))
	return err
}

func ImportJSON(dbName, table, jsonFile string) (int, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return 0, err
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}

	var columns []string
	for k := range rows[0] {
		columns = append(columns, k)
	}

	db, err := OpenDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	colDefs := make([]string, len(columns))
	for i, c := range columns {
		colDefs[i] = fmt.Sprintf(`"%s" TEXT`, c)
	}
	createSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, table, strings.Join(colDefs, ", "))
	if _, err := db.Exec(createSQL); err != nil {
		return 0, err
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	insertSQL := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table,
		strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, row := range rows {
		vals := make([]interface{}, len(columns))
		for i, c := range columns {
			vals[i] = fmt.Sprintf("%v", row[c])
		}
		if _, err := stmt.Exec(vals...); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	return len(rows), tx.Commit()
}

type ColumnInfo struct {
	CID       int    `json:"cid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	NotNull   bool   `json:"notnull"`
	Default   string `json:"default"`
	PK        bool   `json:"pk"`
}

func Describe(dbName, table string) ([]ColumnInfo, error) {
	db, err := OpenDB(dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(\"%s\")", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		var notnull, pk int
		var dflt *string
		if err := rows.Scan(&c.CID, &c.Name, &c.Type, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		c.NotNull = notnull == 1
		c.PK = pk == 1
		if dflt != nil {
			c.Default = *dflt
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

func ImportNDJSON(dbName, table, ndjsonFile string) (int, error) {
	f, err := os.Open(ndjsonFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	db, err := OpenDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var columns []string
	var stmt *sql.Stmt
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	count := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var row map[string]interface{}
		if err := json.Unmarshal(line, &row); err != nil {
			tx.Rollback()
			return count, fmt.Errorf("line %d: %w", count+1, err)
		}
		if stmt == nil {
			for k := range row {
				columns = append(columns, k)
			}
			colDefs := make([]string, len(columns))
			for i, c := range columns {
				colDefs[i] = fmt.Sprintf(`"%s" TEXT`, c)
			}
			createSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s)`, table, strings.Join(colDefs, ", "))
			if _, err := tx.Exec(createSQL); err != nil {
				tx.Rollback()
				return 0, err
			}
			placeholders := make([]string, len(columns))
			for i := range placeholders {
				placeholders[i] = "?"
			}
			insertSQL := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table,
				strings.Join(columns, ", "), strings.Join(placeholders, ", "))
			stmt, err = tx.Prepare(insertSQL)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
		vals := make([]interface{}, len(columns))
		for i, c := range columns {
			vals[i] = fmt.Sprintf("%v", row[c])
		}
		if _, err := stmt.Exec(vals...); err != nil {
			tx.Rollback()
			return count, err
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		tx.Rollback()
		return count, err
	}
	return count, tx.Commit()
}

func DropDB(dbName string) error {
	return os.Remove(dbPath(dbName))
}

func ExportCSV(dbName, sqlStr string, w io.Writer) error {
	result, err := Query(dbName, sqlStr)
	if err != nil {
		return err
	}
	result.FormatCSV(w)
	return nil
}

func ListDBs() ([]string, error) {
	pattern := filepath.Join(config.SQLDir(), "*.db")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		name := filepath.Base(m)
		names = append(names, name[:len(name)-3])
	}
	return names, nil
}
