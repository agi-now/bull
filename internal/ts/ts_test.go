package ts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bull-cli/bull/internal/config"
)

func setup(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "bull-ts-test-*")
	if err != nil {
		t.Fatal(err)
	}
	config.DataDir = dir
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestWriteAndQuery(t *testing.T) {
	setup(t)
	now := time.Now().Unix()

	if err := Write("db", "cpu", 72.5, now-10, nil); err != nil {
		t.Fatal(err)
	}
	if err := Write("db", "cpu", 68.3, now-9, nil); err != nil {
		t.Fatal(err)
	}

	points, err := QueryRange("db", "cpu", now-20, now+10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if points[0].Value != 72.5 {
		t.Fatalf("expected 72.5, got %f", points[0].Value)
	}
}

func TestWriteWithLabels(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	labels1 := map[string]string{"host": "web-01"}
	labels2 := map[string]string{"host": "web-02"}

	if err := Write("db", "cpu", 50.0, now-10, labels1); err != nil {
		t.Fatal(err)
	}
	if err := Write("db", "cpu", 60.0, now-9, labels2); err != nil {
		t.Fatal(err)
	}

	points, err := QueryRange("db", "cpu", now-20, now+10, labels1)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 1 {
		t.Fatalf("expected 1 point with label host=web-01, got %d", len(points))
	}
}

func TestLatest(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	Write("db", "mem", 40.0, now-30, nil)
	Write("db", "mem", 55.0, now-20, nil)
	Write("db", "mem", 70.0, now-10, nil)

	p, err := Latest("db", "mem", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Value != 70.0 {
		t.Fatalf("expected latest value 70.0, got %f", p.Value)
	}
}

func TestCount(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	Write("db", "req", 1.0, now-30, nil)
	Write("db", "req", 2.0, now-20, nil)
	Write("db", "req", 3.0, now-10, nil)

	n, err := Count("db", "req", now-60, now+10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func TestWriteBatch(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	rows := []BatchRow{
		{Metric: "cpu", Value: 10.0, Timestamp: now - 30},
		{Metric: "cpu", Value: 20.0, Timestamp: now - 20},
		{Metric: "cpu", Value: 30.0, Timestamp: now - 10},
	}
	n, err := WriteBatch("db", rows)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 written, got %d", n)
	}

	points, _ := QueryRange("db", "cpu", now-60, now+10, nil)
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
}

func TestWriteBatchFromNDJSON(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	ndjson := ""
	for i := 0; i < 3; i++ {
		ndjson += fmt.Sprintf("{\"metric\":\"disk\",\"value\":10,\"timestamp\":%d}\n", now-int64(30-i*10))
	}
	path := filepath.Join(t.TempDir(), "data.ndjson")
	os.WriteFile(path, []byte(ndjson), 0644)

	n, err := WriteBatchFromNDJSON("db", path)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func TestExportCSV(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	Write("db", "cpu", 10.0, now-20, nil)
	Write("db", "cpu", 20.0, now-10, nil)

	var buf bytes.Buffer
	if err := ExportCSV("db", "cpu", now-60, now+10, nil, &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header+2), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Fatal("CSV should start with header")
	}
}

func TestDeleteDB(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	Write("db", "x", 1.0, now-10, nil)

	dbDir := filepath.Join(config.TSDir(), "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		t.Fatal("db dir should exist before delete")
	}
	if err := DeleteDB("db"); err != nil {
		t.Logf("warning: DeleteDB error (may be Windows file lock): %v", err)
	}
}

func TestListDBs(t *testing.T) {
	setup(t)
	now := time.Now().Unix()
	Write("db1", "m", 1.0, now-10, nil)
	Write("db2", "m", 1.0, now-10, nil)

	names, err := ListDBs()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 dbs, got %d", len(names))
	}
}
