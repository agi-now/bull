package search

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	index "github.com/blevesearch/bleve_index_api"
	"github.com/bull-cli/bull/internal/config"
)

func indexPath(name string) string {
	return filepath.Join(config.SearchDir(), name+".bleve")
}

func Create(name string) error {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.New(indexPath(name), mapping)
	if err != nil {
		return err
	}
	return idx.Close()
}

func openIndex(name string) (bleve.Index, error) {
	return bleve.Open(indexPath(name))
}

func Index(name, docID, jsonStr string) error {
	idx, err := openIndex(name)
	if err != nil {
		return err
	}
	defer idx.Close()

	var doc interface{}
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return err
	}
	return idx.Index(docID, doc)
}

type SearchHit struct {
	ID     string            `json:"id"`
	Score  float64           `json:"score"`
	Fields map[string]string `json:"fields,omitempty"`
}

type SearchResult struct {
	Total int         `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

func QueryIndex(name, queryStr string, limit int) (*SearchResult, error) {
	idx, err := openIndex(name)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	q := bleve.NewQueryStringQuery(queryStr)
	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	res, err := idx.Search(req)
	if err != nil {
		return nil, err
	}

	result := &SearchResult{Total: int(res.Total)}
	for _, hit := range res.Hits {
		result.Hits = append(result.Hits, SearchHit{
			ID:    hit.ID,
			Score: hit.Score,
		})
	}
	return result, nil
}

func DeleteDoc(name, docID string) error {
	idx, err := openIndex(name)
	if err != nil {
		return err
	}
	defer idx.Close()
	return idx.Delete(docID)
}

func QueryIndexWithFields(name, queryStr string, limit int, fields []string) (*SearchResult, error) {
	idx, err := openIndex(name)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	q := bleve.NewQueryStringQuery(queryStr)
	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	if len(fields) > 0 {
		req.Fields = fields
	} else {
		req.Fields = []string{"*"}
	}
	res, err := idx.Search(req)
	if err != nil {
		return nil, err
	}

	result := &SearchResult{Total: int(res.Total)}
	for _, hit := range res.Hits {
		h := SearchHit{
			ID:    hit.ID,
			Score: hit.Score,
		}
		if hit.Fields != nil {
			h.Fields = make(map[string]string)
			for k, v := range hit.Fields {
				h.Fields[k] = fmt.Sprintf("%v", v)
			}
		}
		result.Hits = append(result.Hits, h)
	}
	return result, nil
}

func GetDoc(name, docID string) (map[string]interface{}, error) {
	idx, err := openIndex(name)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	doc, err := idx.Document(docID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("document %q not found", docID)
	}
	result := make(map[string]interface{})
	result["_id"] = docID
	doc.VisitFields(func(field index.Field) {
		result[field.Name()] = string(field.Value())
	})
	return result, nil
}

func Info(name string) (uint64, error) {
	idx, err := openIndex(name)
	if err != nil {
		return 0, err
	}
	defer idx.Close()
	return idx.DocCount()
}

func BulkIndex(name, ndjsonFile string) (int, error) {
	idx, err := openIndex(name)
	if err != nil {
		return 0, err
	}
	defer idx.Close()

	f, err := os.Open(ndjsonFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	batch := idx.NewBatch()
	count := 0
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(line, &doc); err != nil {
			return count, fmt.Errorf("line %d: %w", count+1, err)
		}
		docID := fmt.Sprintf("%d", count+1)
		if id, ok := doc["_id"]; ok {
			docID = fmt.Sprintf("%v", id)
			delete(doc, "_id")
		} else if id, ok := doc["id"]; ok {
			docID = fmt.Sprintf("%v", id)
		}
		if err := batch.Index(docID, doc); err != nil {
			return count, err
		}
		count++
		if count%1000 == 0 {
			if err := idx.Batch(batch); err != nil {
				return count, err
			}
			batch = idx.NewBatch()
		}
	}
	if batch.Size() > 0 {
		if err := idx.Batch(batch); err != nil {
			return count, err
		}
	}
	return count, scanner.Err()
}

func DropIndex(name string) error {
	return os.RemoveAll(indexPath(name))
}

func ListDBs() ([]string, error) {
	pattern := filepath.Join(config.SearchDir(), "*.bleve")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		name := filepath.Base(m)
		names = append(names, name[:len(name)-6])
	}
	return names, nil
}
