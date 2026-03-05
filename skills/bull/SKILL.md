---
name: bull
description: "Micro data environment in a single ~18MB binary — KV store, SQL database, graph analysis, full-text search, and time-series storage, all out of the box. No servers to install, no dependencies to manage. TRIGGER when: user needs local data storage, CSV/JSON analysis, graph traversal, full-text search, metrics recording, persistent state between agent steps, or building data pipelines. CLI-driven, scriptable, pipe-friendly."
license: Apache-2.0
compatibility: "Requires the bull binary in PATH. Supports linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64."
repository: https://github.com/agi-now/bull
install:
  binary_pattern: "bull-{os}-{arch}"
  platforms:
    - linux/amd64
    - linux/arm64
    - darwin/amd64
    - darwin/arm64
    - windows/amd64
  steps:
    - "Download from GitHub Releases: https://github.com/agi-now/bull/releases/latest"
    - "Move to PATH and chmod +x"
  quick: "curl -fsSL https://github.com/agi-now/bull/releases/latest/download/bull-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o /usr/local/bin/bull && chmod +x /usr/local/bin/bull"
metadata:
  author: agi-now
  version: "1.0"
---

# Bull — Micro Data Environment, One Binary

Bull is a self-contained data environment in a single static binary. Five embedded engines — KV, SQL, Graph, Search, Time-Series — ready to use the moment you download it. No databases to install, no servers to configure, no dependencies to chase. Just `bull` and your data.

Use CLI to build data pipelines step by step — scriptable, pipe-friendly, and perfect for agent workflows that persist state, analyze, and query — all locally.

## Installation

**One-liner (Linux / macOS):**
```bash
curl -fsSL https://github.com/agi-now/bull/releases/latest/download/bull-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o /usr/local/bin/bull && chmod +x /usr/local/bin/bull
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/agi-now/bull/releases/latest/download/bull-windows-amd64.exe" -OutFile "$env:LOCALAPPDATA\bull.exe"
```

**Manual:** Download the binary for your platform from [GitHub Releases](https://github.com/agi-now/bull/releases/latest), rename to `bull` (or `bull.exe`), and place it in your PATH.

**Verify:**
```bash
bull version
```

## Engines

| Engine | Powered By | What It Does |
|--------|-----------|--------------|
| **kv** | bbolt (B+tree) | Persistent key-value storage with buckets, counters, batch ops |
| **sql** | SQLite (pure Go) | Full SQL — import CSV/JSON, query, join, aggregate, export |
| **graph** | dominikbraun/graph | Directed/undirected graphs — shortest path, DFS/BFS, cycle detection, toposort |
| **search** | bleve | Full-text search — index JSON documents, query with scoring and field return |
| **ts** | tstorage | Time-series — write metrics with labels, range query, export CSV |

## Agent Strategy: Offload to Bull, Save Tokens

When working with data, follow this principle: **understand the data first, preprocess with Python, then hand off to bull for storage and analysis.**

1. **Inspect** — Read a small sample of the data source (head, schema, dtypes) to understand structure. Do NOT load entire datasets into conversation context.
2. **Preprocess with Python** — Use Python scripts for cleaning, type conversion, column renaming, filtering, reshaping, or format conversion (e.g. Excel/Parquet → CSV/NDJSON). Write the result to a file.
3. **Import into bull** — Use `bull sql import`, `bull search bulk`, `bull graph import-csv`, or `bull ts bulk` to load the preprocessed file.
4. **Analyze with bull** — Run queries, aggregations, searches, and graph algorithms via bull CLI. Return `--format json` results directly — no need to parse large outputs in Python.

This keeps token usage minimal: Python handles the heavy byte-level work offline, bull handles persistent storage and query execution, and only compact JSON results enter the conversation.

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
- Errors go to stderr with non-zero exit code. Examples:
  ```
  Error: key "session_token" not found in bucket "default"
  Error: index "articles" does not exist
  Error: vertex "svc-unknown" not found in graph "deps"
  ```

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

## Edge Cases

- KV bucket defaults to `"default"` if `--bucket` is omitted
- Graph defaults to directed mode; use `--undirected` for undirected graphs
- SQL import auto-creates tables with TEXT columns; cast in queries for numeric comparison
- Search index must be created (`bull search create`) before indexing documents
- TS timestamps are Unix seconds; `--from`/`--to` default to last 1 hour if omitted
- All data persists under `--data-dir` (default `./data`) with per-engine subdirectories
