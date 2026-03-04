package kv

import (
	"os"
	"testing"

	"github.com/agi-now/bull/internal/config"
)

func setup(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	config.DataDir = dir
}

func TestPutGetDel(t *testing.T) {
	setup(t)
	if err := Put("test", "", "k1", "v1"); err != nil {
		t.Fatal(err)
	}
	val, err := Get("test", "", "k1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "v1" {
		t.Fatalf("expected v1, got %s", val)
	}
	if err := Del("test", "", "k1"); err != nil {
		t.Fatal(err)
	}
	_, err = Get("test", "", "k1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestBucket(t *testing.T) {
	setup(t)
	if err := Put("test", "mybkt", "a", "1"); err != nil {
		t.Fatal(err)
	}
	if err := Put("test", "mybkt2", "b", "2"); err != nil {
		t.Fatal(err)
	}
	buckets, err := ListBuckets("test")
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestListWithPrefix(t *testing.T) {
	setup(t)
	Put("test", "", "user:1", "alice")
	Put("test", "", "user:2", "bob")
	Put("test", "", "item:1", "phone")

	pairs, err := List("test", "", "user:")
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs with prefix user:, got %d", len(pairs))
	}
}

func TestScan(t *testing.T) {
	setup(t)
	Put("test", "", "a", "1")
	Put("test", "", "b", "2")
	Put("test", "", "c", "3")
	Put("test", "", "d", "4")

	pairs, err := Scan("test", "", "b", "c")
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs in [b,c], got %d", len(pairs))
	}
	if pairs[0].Key != "b" || pairs[1].Key != "c" {
		t.Fatalf("unexpected keys: %v", pairs)
	}
}

func TestExists(t *testing.T) {
	setup(t)
	Put("test", "", "k1", "v1")

	ok, err := Exists("test", "", "k1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected k1 to exist")
	}
	ok, err = Exists("test", "", "nope")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected nope to not exist")
	}
}

func TestCount(t *testing.T) {
	setup(t)
	Put("test", "", "a", "1")
	Put("test", "", "b", "2")

	n, err := Count("test", "")
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected count 2, got %d", n)
	}
}

func TestIncr(t *testing.T) {
	setup(t)
	val, err := Incr("test", "", "counter", 1)
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
	val, err = Incr("test", "", "counter", 5)
	if err != nil {
		t.Fatal(err)
	}
	if val != 6 {
		t.Fatalf("expected 6, got %d", val)
	}
	val, err = Incr("test", "", "counter", -2)
	if err != nil {
		t.Fatal(err)
	}
	if val != 4 {
		t.Fatalf("expected 4, got %d", val)
	}
}

func TestMGetMPut(t *testing.T) {
	setup(t)
	pairs := []KVPair{
		{Key: "a", Value: "1"},
		{Key: "b", Value: "2"},
		{Key: "c", Value: "3"},
	}
	if err := MPut("test", "", pairs); err != nil {
		t.Fatal(err)
	}
	result, err := MGet("test", "", []string{"a", "c", "missing"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	if result[0].Value != "1" {
		t.Fatalf("expected 1 for a, got %s", result[0].Value)
	}
	if result[1].Value != "3" {
		t.Fatalf("expected 3 for c, got %s", result[1].Value)
	}
	if result[2].Value != "" {
		t.Fatalf("expected empty for missing, got %s", result[2].Value)
	}
}

func TestExportImportJSON(t *testing.T) {
	setup(t)
	Put("test", "", "x", "10")
	Put("test", "", "y", "20")

	exported, err := ExportJSON("test", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(exported) != 2 {
		t.Fatalf("expected 2 exported pairs, got %d", len(exported))
	}

	if err := ImportJSON("test2", "", exported); err != nil {
		t.Fatal(err)
	}
	val, err := Get("test2", "", "x")
	if err != nil {
		t.Fatal(err)
	}
	if val != "10" {
		t.Fatalf("expected 10, got %s", val)
	}
}

func TestDropDB(t *testing.T) {
	setup(t)
	Put("test", "", "k", "v")
	if err := DropDB("test"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dbPath("test")); !os.IsNotExist(err) {
		t.Fatal("expected db file to be removed")
	}
}

func TestDropBucket(t *testing.T) {
	setup(t)
	Put("test", "bkt1", "k", "v")
	Put("test", "bkt2", "k", "v")

	if err := DropBucket("test", "bkt1"); err != nil {
		t.Fatal(err)
	}
	buckets, _ := ListBuckets("test")
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket after drop, got %d", len(buckets))
	}
}

func TestListDBs(t *testing.T) {
	setup(t)
	Put("db1", "", "k", "v")
	Put("db2", "", "k", "v")

	names, err := ListDBs()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 dbs, got %d", len(names))
	}
}
