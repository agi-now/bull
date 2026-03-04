package sql

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agi-now/bull/internal/config"
)

func setup(t *testing.T) {
	t.Helper()
	config.DataDir = t.TempDir()
}

func TestExecAndQuery(t *testing.T) {
	setup(t)
	_, err := Exec("test", "CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := Exec("test", "INSERT INTO t VALUES (1, 'alice'), (2, 'bob')")
	if err != nil {
		t.Fatal(err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 rows affected, got %d", affected)
	}
	result, err := Query("test", "SELECT * FROM t ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(result.Columns))
	}
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
	if result.Rows[0][1] != "alice" {
		t.Fatalf("expected alice, got %s", result.Rows[0][1])
	}
}

func TestFormatTable(t *testing.T) {
	result := &QueryResult{
		Columns: []string{"id", "name"},
		Rows:    [][]string{{"1", "alice"}},
	}
	var buf bytes.Buffer
	result.FormatTable(&buf)
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Fatal("table format should contain alice")
	}
	if !strings.Contains(out, "+") {
		t.Fatal("table format should have borders")
	}
}

func TestFormatCSV(t *testing.T) {
	result := &QueryResult{
		Columns: []string{"id", "name"},
		Rows:    [][]string{{"1", "alice"}},
	}
	var buf bytes.Buffer
	result.FormatCSV(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 CSV lines, got %d", len(lines))
	}
}

func TestFormatJSON(t *testing.T) {
	result := &QueryResult{
		Columns: []string{"id", "name"},
		Rows:    [][]string{{"1", "alice"}},
	}
	var buf bytes.Buffer
	result.FormatJSON(&buf)
	var arr []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatal(err)
	}
	if len(arr) != 1 || arr[0]["name"] != "alice" {
		t.Fatalf("unexpected JSON: %v", arr)
	}
}

func TestTablesSchemaCount(t *testing.T) {
	setup(t)
	Exec("test", "CREATE TABLE users (id INTEGER, name TEXT)")
	Exec("test", "INSERT INTO users VALUES (1,'a'),(2,'b'),(3,'c')")

	tables, err := Tables("test")
	if err != nil {
		t.Fatal(err)
	}
	if len(tables) != 1 || tables[0] != "users" {
		t.Fatalf("expected [users], got %v", tables)
	}

	ddl, err := Schema("test", "users")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(ddl, "CREATE TABLE") {
		t.Fatalf("expected CREATE TABLE in schema, got %s", ddl)
	}

	n, err := CountRows("test", "users")
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 rows, got %d", n)
	}
}

func TestDescribe(t *testing.T) {
	setup(t)
	Exec("test", "CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT NOT NULL, price REAL)")

	cols, err := Describe("test", "items")
	if err != nil {
		t.Fatal(err)
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
	if cols[0].Name != "id" || !cols[0].PK {
		t.Fatal("expected id to be PK")
	}
	if cols[1].Name != "name" || !cols[1].NotNull {
		t.Fatal("expected name to be NOT NULL")
	}
}

func TestImportCSV(t *testing.T) {
	setup(t)
	csvContent := "name,age\nalice,30\nbob,25\n"
	csvPath := filepath.Join(t.TempDir(), "test.csv")
	os.WriteFile(csvPath, []byte(csvContent), 0644)

	n, err := ImportCSV("test", "people", csvPath)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 rows imported, got %d", n)
	}
	cnt, _ := CountRows("test", "people")
	if cnt != 2 {
		t.Fatalf("expected 2 rows in table, got %d", cnt)
	}
}

func TestImportJSON(t *testing.T) {
	setup(t)
	jsonContent := `[{"name":"alice","age":"30"},{"name":"bob","age":"25"}]`
	jsonPath := filepath.Join(t.TempDir(), "test.json")
	os.WriteFile(jsonPath, []byte(jsonContent), 0644)

	n, err := ImportJSON("test", "people", jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 rows, got %d", n)
	}
}

func TestImportNDJSON(t *testing.T) {
	setup(t)
	ndjson := "{\"name\":\"alice\",\"city\":\"NYC\"}\n{\"name\":\"bob\",\"city\":\"LA\"}\n{\"name\":\"charlie\",\"city\":\"SF\"}\n"
	ndjsonPath := filepath.Join(t.TempDir(), "test.ndjson")
	os.WriteFile(ndjsonPath, []byte(ndjson), 0644)

	n, err := ImportNDJSON("test", "logs", ndjsonPath)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 rows, got %d", n)
	}
}

func TestExecFile(t *testing.T) {
	setup(t)
	sqlContent := "CREATE TABLE t (x TEXT); INSERT INTO t VALUES ('hello');"
	sqlPath := filepath.Join(t.TempDir(), "init.sql")
	os.WriteFile(sqlPath, []byte(sqlContent), 0644)

	if err := ExecFile("test", sqlPath); err != nil {
		t.Fatal(err)
	}
	cnt, _ := CountRows("test", "t")
	if cnt != 1 {
		t.Fatalf("expected 1 row, got %d", cnt)
	}
}

func TestExportCSV(t *testing.T) {
	setup(t)
	Exec("test", "CREATE TABLE t (a TEXT, b TEXT)")
	Exec("test", "INSERT INTO t VALUES ('1','x'),('2','y')")

	var buf bytes.Buffer
	if err := ExportCSV("test", "SELECT * FROM t", &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 CSV lines (header+2 rows), got %d", len(lines))
	}
}

func TestDropDB(t *testing.T) {
	setup(t)
	Exec("test", "CREATE TABLE t (x TEXT)")
	if err := DropDB("test"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dbPath("test")); !os.IsNotExist(err) {
		t.Fatal("expected db file removed")
	}
}

func TestListDBs(t *testing.T) {
	setup(t)
	Exec("db1", "CREATE TABLE t (x TEXT)")
	Exec("db2", "CREATE TABLE t (x TEXT)")

	names, err := ListDBs()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 dbs, got %d", len(names))
	}
}
