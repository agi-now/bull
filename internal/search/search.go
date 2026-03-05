package search

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/agi-now/bull/internal/config"
	_ "modernc.org/sqlite"
)

func dbPath(name string) string {
	return filepath.Join(config.SearchDir(), name+".db")
}

func openDB(name string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath(name))
}

func Create(name string) error {
	db, err := openDB(name)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS docs USING fts5(doc_id, content, detail=full)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS docs_raw (doc_id TEXT PRIMARY KEY, data TEXT NOT NULL)`)
	return err
}

func ensureTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS docs USING fts5(doc_id, content, detail=full)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS docs_raw (doc_id TEXT PRIMARY KEY, data TEXT NOT NULL)`)
	return err
}

func flattenJSON(data map[string]interface{}) string {
	var parts []string
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%v", data[k]))
	}
	return strings.Join(parts, " ")
}

func Index(name, docID, jsonStr string) error {
	db, err := openDB(name)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureTables(db); err != nil {
		return err
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return err
	}
	content := flattenJSON(doc)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Remove old entry if exists
	tx.Exec(`DELETE FROM docs WHERE doc_id = ?`, docID)
	tx.Exec(`DELETE FROM docs_raw WHERE doc_id = ?`, docID)

	if _, err := tx.Exec(`INSERT INTO docs (doc_id, content) VALUES (?, ?)`, docID, content); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`INSERT INTO docs_raw (doc_id, data) VALUES (?, ?)`, docID, jsonStr); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
	return QueryIndexWithFields(name, queryStr, limit, 0, nil)
}

func QueryIndexWithFields(name, queryStr string, limit, offset int, fields []string) (*SearchResult, error) {
	db, err := openDB(name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	countRow := db.QueryRow(`SELECT count(*) FROM docs WHERE docs MATCH ?`, queryStr)
	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`SELECT d.doc_id, rank, r.data FROM docs d JOIN docs_raw r ON d.doc_id = r.doc_id WHERE docs MATCH ? ORDER BY rank LIMIT ? OFFSET ?`,
		queryStr, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &SearchResult{Total: total}
	for rows.Next() {
		var id, rawData string
		var rank float64
		if err := rows.Scan(&id, &rank, &rawData); err != nil {
			return nil, err
		}
		h := SearchHit{
			ID:    id,
			Score: -rank, // FTS5 rank is negative (lower = better), invert for consistency
		}
		if len(fields) > 0 {
			var doc map[string]interface{}
			if json.Unmarshal([]byte(rawData), &doc) == nil {
				h.Fields = make(map[string]string)
				for _, f := range fields {
					if v, ok := doc[f]; ok {
						h.Fields[f] = fmt.Sprintf("%v", v)
					}
				}
			}
		} else {
			var doc map[string]interface{}
			if json.Unmarshal([]byte(rawData), &doc) == nil {
				h.Fields = make(map[string]string)
				for k, v := range doc {
					h.Fields[k] = fmt.Sprintf("%v", v)
				}
			}
		}
		result.Hits = append(result.Hits, h)
	}
	return result, rows.Err()
}

func DeleteDoc(name, docID string) error {
	db, err := openDB(name)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	tx.Exec(`DELETE FROM docs WHERE doc_id = ?`, docID)
	tx.Exec(`DELETE FROM docs_raw WHERE doc_id = ?`, docID)
	return tx.Commit()
}

func GetDoc(name, docID string) (map[string]interface{}, error) {
	db, err := openDB(name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rawData string
	err = db.QueryRow(`SELECT data FROM docs_raw WHERE doc_id = ?`, docID).Scan(&rawData)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("document %q not found", docID)
	}
	if err != nil {
		return nil, err
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(rawData), &doc); err != nil {
		return nil, err
	}
	doc["_id"] = docID
	return doc, nil
}

func Info(name string) (uint64, error) {
	db, err := openDB(name)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var count uint64
	err = db.QueryRow(`SELECT count(*) FROM docs_raw`).Scan(&count)
	return count, err
}

func BulkIndex(name, ndjsonFile string) (int, error) {
	db, err := openDB(name)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := ensureTables(db); err != nil {
		return 0, err
	}

	f, err := os.Open(ndjsonFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmtFTS, err := tx.Prepare(`INSERT INTO docs (doc_id, content) VALUES (?, ?)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	stmtRaw, err := tx.Prepare(`INSERT OR REPLACE INTO docs_raw (doc_id, data) VALUES (?, ?)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	count := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(line, &doc); err != nil {
			tx.Rollback()
			return count, fmt.Errorf("line %d: %w", count+1, err)
		}
		docID := fmt.Sprintf("%d", count+1)
		if id, ok := doc["_id"]; ok {
			docID = fmt.Sprintf("%v", id)
			delete(doc, "_id")
		} else if id, ok := doc["id"]; ok {
			docID = fmt.Sprintf("%v", id)
		}

		content := flattenJSON(doc)
		rawBytes, _ := json.Marshal(doc)

		if _, err := stmtFTS.Exec(docID, content); err != nil {
			tx.Rollback()
			return count, err
		}
		if _, err := stmtRaw.Exec(docID, string(rawBytes)); err != nil {
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

func UpdateDoc(name, docID, jsonStr string) error {
	return Index(name, docID, jsonStr)
}

func DropIndex(name string) error {
	return os.Remove(dbPath(name))
}

func ListDBs() ([]string, error) {
	pattern := filepath.Join(config.SearchDir(), "*.db")
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
