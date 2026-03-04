package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bull-cli/bull/internal/config"
)

func setup(t *testing.T) {
	t.Helper()
	config.DataDir = t.TempDir()
}

func createAndIndex(t *testing.T) {
	t.Helper()
	if err := Create("idx"); err != nil {
		t.Fatal(err)
	}
	if err := Index("idx", "doc1", `{"title":"Hello World","body":"This is a test document"}`); err != nil {
		t.Fatal(err)
	}
	if err := Index("idx", "doc2", `{"title":"Goodbye World","body":"Another test document"}`); err != nil {
		t.Fatal(err)
	}
}

func TestCreateAndInfo(t *testing.T) {
	setup(t)
	if err := Create("idx"); err != nil {
		t.Fatal(err)
	}
	count, err := Info("idx")
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 docs, got %d", count)
	}
}

func TestIndexAndQuery(t *testing.T) {
	setup(t)
	createAndIndex(t)

	result, err := QueryIndex("idx", "Hello", 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Fatalf("expected 1 hit for Hello, got %d", result.Total)
	}
	if result.Hits[0].ID != "doc1" {
		t.Fatalf("expected doc1, got %s", result.Hits[0].ID)
	}
}

func TestQueryWithFields(t *testing.T) {
	setup(t)
	createAndIndex(t)

	result, err := QueryIndexWithFields("idx", "World", 10, 0, []string{"title"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 hits for World, got %d", result.Total)
	}
	for _, h := range result.Hits {
		if h.Fields == nil {
			t.Fatal("expected fields in result")
		}
		if _, ok := h.Fields["title"]; !ok {
			t.Fatal("expected title field in results")
		}
	}
}

func TestQueryWithOffset(t *testing.T) {
	setup(t)
	createAndIndex(t)

	result, err := QueryIndexWithFields("idx", "World", 1, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Hits) != 1 {
		t.Fatalf("expected 1 hit with limit=1, got %d", len(result.Hits))
	}

	result2, err := QueryIndexWithFields("idx", "World", 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result2.Hits) != 1 {
		t.Fatalf("expected 1 hit at offset=1, got %d", len(result2.Hits))
	}
	if result.Hits[0].ID == result2.Hits[0].ID {
		t.Fatal("offset=0 and offset=1 should return different docs")
	}
}

func TestGetDoc(t *testing.T) {
	setup(t)
	createAndIndex(t)

	doc, err := GetDoc("idx", "doc1")
	if err != nil {
		t.Fatal(err)
	}
	if doc["_id"] != "doc1" {
		t.Fatalf("expected _id=doc1, got %v", doc["_id"])
	}
}

func TestUpdateDoc(t *testing.T) {
	setup(t)
	createAndIndex(t)

	if err := UpdateDoc("idx", "doc1", `{"title":"Updated Title","body":"new body"}`); err != nil {
		t.Fatal(err)
	}
	result, err := QueryIndex("idx", "Updated", 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Fatalf("expected 1 hit for Updated, got %d", result.Total)
	}
}

func TestDeleteDoc(t *testing.T) {
	setup(t)
	createAndIndex(t)

	if err := DeleteDoc("idx", "doc1"); err != nil {
		t.Fatal(err)
	}
	count, _ := Info("idx")
	if count != 1 {
		t.Fatalf("expected 1 doc after delete, got %d", count)
	}
}

func TestBulkIndex(t *testing.T) {
	setup(t)
	Create("idx")

	ndjson := `{"_id":"a1","title":"First"}
{"_id":"a2","title":"Second"}
{"_id":"a3","title":"Third"}
`
	path := filepath.Join(t.TempDir(), "docs.ndjson")
	os.WriteFile(path, []byte(ndjson), 0644)

	n, err := BulkIndex("idx", path)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 indexed, got %d", n)
	}
	count, _ := Info("idx")
	if count != 3 {
		t.Fatalf("expected 3 docs, got %d", count)
	}
}

func TestDropIndex(t *testing.T) {
	setup(t)
	Create("idx")
	if err := DropIndex("idx"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(indexPath("idx")); !os.IsNotExist(err) {
		t.Fatal("expected index dir to be removed")
	}
}

func TestListDBs(t *testing.T) {
	setup(t)
	Create("idx1")
	Create("idx2")

	names, err := ListDBs()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 indexes, got %d", len(names))
	}
}
