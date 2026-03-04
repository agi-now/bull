package search

import (
	"encoding/json"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
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

func Info(name string) (uint64, error) {
	idx, err := openIndex(name)
	if err != nil {
		return 0, err
	}
	defer idx.Close()
	return idx.DocCount()
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
