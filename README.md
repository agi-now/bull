<p align="center">
  <img src="https://img.shields.io/badge/binary-~20MB-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/engines-5-green?style=flat-square" />
  <img src="https://img.shields.io/badge/commands-72+-orange?style=flat-square" />
  <img src="https://img.shields.io/badge/CGo-none-brightgreen?style=flat-square" />
  <img src="https://img.shields.io/badge/platforms-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square" />
</p>

<p align="center"><b>Bull</b> — All-in-One Embedded Engine Toolkit</p>

<p align="center">
  <a href="README_zh.md">中文文档</a>
</p>

---

Five data engines. One static binary. Zero external dependencies.

**Bull** packs a KV store, SQL database, graph engine, full-text search, and time-series storage into a single ~20 MB Go executable. It is purpose-built for **AI Agent skill extensions** — drop the binary into any sandboxed environment and instantly unlock local data processing capabilities that would normally require installing multiple database servers.

## Why Bull?

- **Single binary, zero dependencies** — no database servers to install, no runtime to configure. Copy one file and you're done.
- **5 engines, 72+ commands** — KV, SQL, Graph, Full-text Search, Time-Series — covers the vast majority of data processing scenarios an AI Agent would encounter.
- **CLI first** — every engine is accessible via straightforward command-line interface, making it easy to integrate with any language or framework.
- **AI-Agent native** — ships with machine-readable skill definitions in `skills/`, enabling AI Agents to autonomously decide which engine and command to use.
- **Pure Go, statically compiled** — no CGo, no Wasm, no shared libraries. Cross-compile for Linux / macOS / Windows with a single command.

## At a Glance

```
bull kv put config host 10.0.0.1          # persistent key-value
bull sql query db "SELECT * FROM t"       # full SQL (SQLite)
bull graph shortest-path g A B            # graph algorithms
bull search query idx "error timeout"     # full-text search
bull ts latest mon cpu --format json      # time-series metrics
```

## Engines

| Engine | Powered By | What It Does |
|--------|-----------|--------------|
| **KV** | [bbolt](https://github.com/etcd-io/bbolt) | B+tree KV store — buckets, batch ops, atomic counters, range scans, JSON import/export |
| **SQL** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Full SQLite in pure Go — CSV/JSON/NDJSON import, multi-format query output, interactive shell |
| **Graph** | [dominikbraun/graph](https://github.com/dominikbraun/graph) | Directed & undirected weighted graphs — shortest path, DFS/BFS, topological sort, cycle detection, connected components |
| **Search** | [bleve](https://github.com/blevesearch/bleve) | Full-text indexing — scored queries, field return, pagination, bulk NDJSON indexing |
| **TS** | [tstorage](https://github.com/nakabonne/tstorage) | Time-series storage — labeled metrics, range queries, latest-point lookup, CSV export |

All engines are **pure Go** — no CGo, no Wasm, no shared libraries. The binary is fully statically compiled.

## Install

```bash
# Clone and build
git clone https://github.com/agi-now/bull.git && cd bull
go build -ldflags="-s -w" -o bull ./cmd/bull/

# Or use build scripts (inject version + build time automatically)
./build.sh                # Linux / macOS
.\build.ps1               # Windows
```

Cross-compile for any target:

```bash
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o bull           ./cmd/bull/
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o bull           ./cmd/bull/
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bull.exe       ./cmd/bull/
```

## Usage

All data persists under `--data-dir` (default `./data`), organized by engine.

### KV

```bash
# Single read/write
bull kv put cache token abc123
bull kv get cache token                        # abc123

# Batch operations
bull kv mput cache k1 v1 k2 v2 k3 v3          # put 3 pairs
bull kv mget cache k1 k2 k3                   # [{"key":"k1","value":"v1"}, ...]

# Counters
bull kv incr stats page_views                  # 1
bull kv incr stats page_views 10               # 11

# Scan & filter
bull kv list cache --prefix k --format json
bull kv scan cache --start a --end z --format json

# Import / Export
bull kv export cache                           # JSON to stdout
bull kv import cache -f backup.json
```

### SQL

```bash
# Import data
bull sql import analytics users users.csv                    # imported 1500 rows
bull sql import-ndjson analytics logs access.ndjson           # imported 8234 rows

# Query
bull sql query analytics "SELECT city, COUNT(*) c FROM users GROUP BY city ORDER BY c DESC" --format json --limit 10

# Inspect
bull sql tables analytics
bull sql describe analytics users
bull sql schema analytics users

# Export
bull sql export analytics "SELECT * FROM users" -o out.csv --format csv
bull sql export analytics "SELECT * FROM users WHERE age>30" --format json

# Interactive
bull sql shell analytics
```

### Graph

```bash
# Build
bull graph add-vertex deps auth --attr type=service
bull graph add-vertex deps db --attr type=database
bull graph add-edge deps auth db --weight 1

# Algorithms
bull graph shortest-path deps auth db          # auth -> db
bull graph toposort deps                       # topological order
bull graph has-cycle deps                      # false
bull graph components deps                     # connected components

# Traversal
bull graph dfs deps auth
bull graph bfs deps auth
bull graph neighbors deps auth

# Bulk
bull graph import-csv social edges.csv --undirected
bull graph stats social --undirected            # vertices: 42, edges: 87
```

### Search

```bash
# Setup
bull search create articles
bull search bulk articles docs.ndjson           # indexed 500 documents

# Query with pagination
bull search query articles "machine learning" --limit 10 --offset 0 --format json
bull search query articles "title:hello" --field title --field body --format json

# CRUD
bull search index articles doc1 '{"title":"Hello","body":"World"}'
bull search update articles doc1 '{"title":"Updated","body":"New content"}'
bull search get articles doc1
bull search delete articles doc1
```

### Time-Series

```bash
# Write
bull ts write mon cpu 72.5 --label host=web-01
bull ts write mon cpu 68.3 --label host=web-01

# Read
bull ts latest mon cpu --label host=web-01 --format json
bull ts query mon cpu --from 1700000000 --to 1700003600 --format json
bull ts count mon cpu

# Export
bull ts export mon cpu -o metrics.csv
```

### Global

```bash
bull version                                   # version, build time, go, os/arch
bull info                                      # summary of all databases across engines
```

## Command Reference

```
bull ─┬─ kv ─────┬─ put / get / del          single key ops
      │          ├─ mget / mput              batch ops
      │          ├─ list / scan              range queries (--format tsv|json)
      │          ├─ exists / count           inspection
      │          ├─ incr / decr              atomic counters
      │          ├─ buckets                  list buckets
      │          ├─ export / import          JSON I/O
      │          ├─ drop / drop-bucket       cleanup
      │          └─ dbs                      list databases
      │
      ├─ sql ────┬─ exec / query             DDL/DML and SELECT (--format, --limit)
      │          ├─ exec-file                run .sql file
      │          ├─ tables / schema / describe / count
      │          ├─ import / import-json / import-ndjson
      │          ├─ export                   CSV or JSON (--format)
      │          ├─ shell                    interactive REPL
      │          ├─ drop                     cleanup
      │          └─ dbs                      list databases
      │
      ├─ graph ──┬─ add-vertex / add-edge    build graph (--attr, --weight)
      │          ├─ del-vertex / del-edge    remove
      │          ├─ vertices / edges         list
      │          ├─ neighbors / degree / attrs
      │          ├─ shortest-path / has-path path finding
      │          ├─ dfs / bfs               traversal
      │          ├─ components / toposort / has-cycle
      │          ├─ stats                    counts
      │          ├─ import-csv / export      I/O
      │          ├─ drop                     cleanup
      │          └─ dbs                      list graphs
      │
      ├─ search ─┬─ create                   new index
      │          ├─ index / bulk             add documents
      │          ├─ query                    search (--field, --limit, --offset)
      │          ├─ get / update / delete    CRUD
      │          ├─ info                     statistics
      │          ├─ drop                     cleanup
      │          └─ dbs                      list indexes
      │
      ├─ ts ─────┬─ write / bulk             ingest
      │          ├─ query                    range query (--from, --to, --label)
      │          ├─ latest / count           inspect
      │          ├─ export                   CSV output
      │          ├─ drop                     cleanup
      │          └─ dbs                      list databases
      │
      ├─ version                             build info
      └─ info                                data directory overview
```

**72+ commands** across 5 engines + 2 global utilities.

## AI Agent Skills

The `skills/` directory contains machine-readable skill definitions for each engine. An AI agent reads these files to decide **which engine to use** and **which commands to invoke** for any given task. Use `--format json` for structured output.

## Project Layout

```
bull/
├── cmd/bull/              CLI entry (cobra)
│   ├── main.go            root, version, info
│   ├── cmd_kv.go          KV subcommands
│   ├── cmd_sql.go         SQL subcommands
│   ├── cmd_graph.go       Graph subcommands
│   ├── cmd_search.go      Search subcommands
│   └── cmd_ts.go          TS subcommands
├── internal/
│   ├── config/            data directory config
│   ├── kv/                bbolt wrapper
│   ├── sql/               SQLite wrapper
│   ├── graph/             graph algorithms
│   ├── search/            bleve wrapper
│   └── ts/                tstorage wrapper
├── skills/                AI Agent skill definitions
├── build.sh / build.ps1   build with version injection
└── data/                  runtime storage (gitignored)
```

## Data Storage

| Engine | Location | Format |
|--------|----------|--------|
| KV | `data/kv/<name>.db` | bbolt B+tree |
| SQL | `data/sql/<name>.db` | SQLite |
| Graph | `data/graph/<name>.json` | JSON |
| Search | `data/search/<name>.bleve/` | bleve index |
| TS | `data/ts/<name>/` | tstorage WAL |

## Platform Support

| OS | Architecture |
|----|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## License

MIT
