# Search Engine Reference

Embedded full-text search powered by bleve. Index JSON documents, query with scoring, field return, and pagination.

An index must be created before indexing documents. Query syntax: `field:value`, `+required`, `-excluded`, `"exact phrase"`, `term*` wildcards.

## Commands

### create
```
bull search create <index>
```
Create a new empty search index.

### index
```
bull search index <index> <docID> <json>
```
Index a single JSON document.
```bash
bull search index articles doc1 '{"title":"Hello","body":"World"}'
```

### bulk
```
bull search bulk <index> <ndjson-file>
```
Bulk index from NDJSON. Uses `_id` or `id` field as doc ID; auto-increments otherwise.

### query
```
bull search query <index> <query> [--field name ...] [--limit N] [--offset N] [--format table|json]
```
Search with pagination. Returns IDs, scores, and optionally field values.
```bash
bull search query articles "title:kubernetes" --field title --field author --limit 5 --format json
```

### get
```
bull search get <index> <docID>
```
Retrieve a document by ID with all stored fields.

### update
```
bull search update <index> <docID> <json>
```
Re-index (replace) a document.

### delete
```
bull search delete <index> <docID>
```
Delete a document by ID.

### info
```
bull search info <index>
```
Show index name and document count.

### drop
```
bull search drop <index>
```
Delete an entire search index.

### dbs
```
bull search dbs
```
List all search indexes.

## HTTP API Endpoints

| Method | Path | Body |
|--------|------|------|
| GET | `/api/search/dbs` | — |
| POST | `/api/search/{idx}/create` | — |
| POST | `/api/search/{idx}/index` | `{"id","doc"}` |
| POST | `/api/search/{idx}/query` | `{"query","limit?","offset?","fields?"}` |
| GET | `/api/search/{idx}/get/{docID}` | — |
| POST | `/api/search/{idx}/update` | `{"id","doc"}` |
| DELETE | `/api/search/{idx}/doc/{docID}` | — |
| GET | `/api/search/{idx}/info` | — |
| DELETE | `/api/search/{idx}` | — |
