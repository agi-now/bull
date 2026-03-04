---
name: bull
description: "All-in-one embedded engine toolkit (~18MB Go binary) providing local KV store, SQL database, graph analysis, full-text search, and time-series storage. TRIGGER when: user needs local data storage, CSV/JSON analysis, graph traversal, full-text search, metrics recording, or persistent state between agent steps. Supports both CLI and HTTP API. Zero external dependencies."
license: Apache-2.0
compatibility: "Requires the bull binary in PATH. Supports linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64."
metadata:
  author: bull-cli
  version: "1.0"
---

# Bull — All-in-One Embedded Engine Toolkit

Bull packages five embedded engines into a single static binary. No servers to install, no network dependencies. Just download and run.

## Engines

| Engine | Powered By | What It Does |
|--------|-----------|--------------|
| **kv** | bbolt (B+tree) | Persistent key-value storage with buckets, counters, batch ops |
| **sql** | SQLite (pure Go) | Full SQL — import CSV/JSON, query, join, aggregate, export |
| **graph** | dominikbraun/graph | Directed/undirected graphs — shortest path, DFS/BFS, cycle detection, toposort |
| **search** | bleve | Full-text search — index JSON documents, query with scoring and field return |
| **ts** | tstorage | Time-series — write metrics with labels, range query, export CSV |

## Quick Decision: Which Engine to Use

**Need to store/retrieve a value by key?** → `bull kv`
- Config, session, cache, counters, state between agent steps
- `bull kv put <db> <key> <value>` / `bull kv get <db> <key>`

**Need to analyze structured data with SQL?** → `bull sql`
- CSV/JSON import, GROUP BY, JOIN, aggregation, export
- `bull sql import <db> <table> data.csv` → `bull sql query <db> "SELECT ..." --format json`

**Need to model relationships or dependencies?** → `bull graph`
- Service dependencies, social networks, task ordering, reachability
- `bull graph add-vertex`, `bull graph add-edge`, `bull graph shortest-path`

**Need to search text content?** → `bull search`
- Articles, logs, documents — keyword search with scoring
- `bull search create <idx>` → `bull search bulk <idx> data.ndjson` → `bull search query <idx> "keyword"`

**Need to record metrics over time?** → `bull ts`
- CPU, memory, request counts, latency, any numeric value with timestamps
- `bull ts write <db> <metric> <value> --label host=server1`

## Global Flags

All commands accept:
- `--data-dir <path>` — Override data directory (default: `./data`)

## Output Conventions

- Use `--format json` wherever available for machine-readable output
- Mutation commands (put, del, add-vertex) produce no stdout on success, exit code 0
- Errors go to stderr with non-zero exit code

## HTTP API Mode

Start the HTTP server:
```
bull serve --port 2880
```

All endpoints return unified JSON:
```json
{"ok": true, "data": ...}
{"ok": false, "error": "..."}
```

API routes mirror CLI commands under `/api/{engine}/{db}/{action}`. See [references/http.md](references/http.md) for the full endpoint list.

## Common Workflows

### Import CSV and analyze with SQL
```bash
bull sql import analytics users users.csv
bull sql query analytics "SELECT city, COUNT(*) as cnt FROM users GROUP BY city ORDER BY cnt DESC" --format json
```

### Build a dependency graph and find paths
```bash
bull graph add-vertex deps auth-svc --attr type=service
bull graph add-vertex deps user-svc --attr type=service
bull graph add-vertex deps db-svc --attr type=database
bull graph add-edge deps auth-svc user-svc
bull graph add-edge deps user-svc db-svc
bull graph shortest-path deps auth-svc db-svc
```

### Index documents and search
```bash
bull search create articles
bull search bulk articles articles.ndjson
bull search query articles "machine learning" --limit 5 --format json
```

### Record and query metrics
```bash
bull ts write monitoring cpu_usage 72.5 --label host=web-01
bull ts query monitoring cpu_usage --label host=web-01 --format json
```

### Persist agent state between steps
```bash
bull kv put pipeline step1 '{"status":"done","rows":1500}'
bull kv put pipeline step2 '{"status":"done","rows":1200}'
bull kv list pipeline --format json
```

### Compare two datasets with SQL JOIN
```bash
bull sql import compare old_prices old.csv
bull sql import compare new_prices new.csv
bull sql query compare "SELECT n.name, o.price as old, n.price as new FROM new_prices n JOIN old_prices o ON n.id=o.id WHERE CAST(n.price AS REAL) > CAST(o.price AS REAL)" --format json
```

### Detect cycles in a build graph
```bash
bull graph has-cycle pipeline
bull graph toposort pipeline
```

## Detailed Command References

For the complete command list of each engine, read the corresponding reference file:

- [references/kv.md](references/kv.md) — 17 commands: put, get, del, mget, mput, list, scan, exists, count, incr, decr, buckets, export, import, drop, drop-bucket, dbs
- [references/sql.md](references/sql.md) — 15 commands: exec, query, exec-file, tables, schema, describe, count, import, import-json, import-ndjson, export, shell, drop, dbs
- [references/graph.md](references/graph.md) — 21 commands: add-vertex, add-edge, del-vertex, del-edge, vertices, edges, neighbors, degree, attrs, shortest-path, has-path, dfs, bfs, stats, components, toposort, has-cycle, import-csv, export, drop, dbs
- [references/search.md](references/search.md) — 11 commands: create, index, bulk, query, get, update, delete, info, drop, dbs
- [references/ts.md](references/ts.md) — 8 commands: write, bulk, query, latest, count, export, drop, dbs
- [references/http.md](references/http.md) — HTTP API endpoint reference

## Edge Cases

- KV bucket defaults to `"default"` if `--bucket` is omitted
- Graph defaults to directed mode; use `--undirected` for undirected graphs
- SQL import auto-creates tables with TEXT columns; cast in queries for numeric comparison
- Search index must be created (`bull search create`) before indexing documents
- TS timestamps are Unix seconds; `--from`/`--to` default to last 1 hour if omitted
- All data persists under `--data-dir` (default `./data`) with per-engine subdirectories
