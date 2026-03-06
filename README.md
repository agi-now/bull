<p align="center">
  <img src="https://img.shields.io/badge/binary-~8MB-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/engines-5-green?style=flat-square" />
  <img src="https://img.shields.io/badge/commands-72+-orange?style=flat-square" />
  <img src="https://img.shields.io/badge/CGo-none-brightgreen?style=flat-square" />
  <img src="https://img.shields.io/badge/platforms-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square" />
</p>

<p align="center"><b>Bull</b> — All-in-One Embedded Engine Toolkit<br/>Data never enters LLM context — all processing stays local. Save your tokens and money.</p>

<p align="center">
  <a href="README_zh.md">中文文档</a>
</p>

---

Five data engines. One static binary. Zero external dependencies.

**Bull** packs a KV store, SQL database, graph engine, full-text search, and time-series storage into a single ~8 MB Go executable. It is purpose-built for **AI Agent skill extensions** — drop the binary into any sandboxed environment and instantly unlock local data processing capabilities that would normally require installing multiple database servers.

## The Problem

When AI Agents process large datasets — CSV files, logs, documents — the conventional approach is to load data into conversation context. This leads to **massive token consumption** and **slow response times**. A 10 MB CSV file easily burns through hundreds of thousands of tokens just to read, let alone analyze.

Bull solves this by **offloading data processing to local engines**. Instead of feeding raw data into the LLM, the Agent imports it into Bull, runs queries and aggregations locally, and only returns compact results to the conversation. The data never enters the token stream.

```
Without Bull:  User → [10MB CSV as tokens] → LLM → answer     (hundreds of thousands of tokens)
With Bull:     User → Agent → bull sql import + query → LLM → answer   (a few hundred tokens)
```

## Why Bull?

- **Slash token costs** — data stays local, only query results enter the conversation. Process millions of rows without burning tokens on raw data.
- **Single binary, zero dependencies** — no database servers to install, no runtime to configure. Copy one file and you're done. Works in environments where Python or other runtimes are not available (minimal containers, CI pipelines, restricted sandboxes).
- **5 engines, 72+ commands** — KV, SQL, Graph, Full-text Search, Time-Series — covers the vast majority of data processing scenarios an AI Agent would encounter.
- **CLI first** — every command is a deterministic shell call. No script to write, no indentation errors, no dependency conflicts. More reliable than LLM-generated Python scripts.
- **AI-Agent native** — ships with machine-readable skill definitions in `skills/`, enabling AI Agents to autonomously decide which engine and command to use.
- **Pure Go, statically compiled** — no CGo, no shared libraries. Cross-compile for Linux / macOS / Windows with a single command.

## At a Glance

```
bull kv put config host 10.0.0.1          # persistent key-value
bull sql query db "SELECT * FROM t"       # full SQL (SQLite)
bull graph shortest-path g A B            # graph algorithms
bull search query idx "error timeout"     # full-text search
bull ts latest mon cpu --format json      # time-series metrics
```

## Example: Incident Investigation in 60 Seconds

An AI Agent receives server logs and a service dependency map. It imports, analyzes, searches, and tracks — all locally, zero tokens wasted on raw data.

```bash
# 1. Import 50k access logs into SQL
bull sql import-ndjson incident access access.ndjson        # imported 50000 rows

# 2. Find the top error-producing services
bull sql query incident "SELECT service, COUNT(*) c FROM access WHERE level='ERROR' GROUP BY service ORDER BY c DESC LIMIT 5" --format json

# 3. Build a search index and find the root cause
bull search create logs
bull search bulk logs access.ndjson
bull search query logs "connection refused port 5432" --field service --field message --limit 5 --format json

# 4. Import service dependency graph and trace the blast radius
bull graph import-csv incident deps.csv
bull graph shortest-path incident api-gateway db-primary
bull graph bfs incident db-primary                          # all affected downstream services

# 5. Record the incident timeline
bull ts write incident_metrics error_rate 127 --label service=db-primary
bull ts write incident_metrics error_rate 3 --label service=db-primary  # after fix

# 6. Save conclusions for the post-mortem
bull kv put incident:2026-03-05 root_cause '{"service":"db-primary","issue":"connection pool exhausted"}'
bull kv put incident:2026-03-05 blast_radius '["api-gateway","user-svc","payment-svc"]'
```

5 engines, 12 commands, one binary — the Agent processed 50k logs without a single row entering the LLM context.

## Download

Get the latest binary for your platform from [**GitHub Releases**](https://github.com/agi-now/bull/releases/latest).

| Platform | Binary |
|----------|--------|
| Linux amd64 | `bull-linux-amd64` |
| Linux arm64 | `bull-linux-arm64` |
| macOS amd64 | `bull-darwin-amd64` |
| macOS arm64 | `bull-darwin-arm64` |
| Windows amd64 | `bull-windows-amd64.exe` |

Download, rename to `bull` (or `bull.exe`), and place it in your PATH.

## Engines

| Engine | Powered By | What It Does |
|--------|-----------|--------------|
| **KV** | [bbolt](https://github.com/etcd-io/bbolt) | B+tree KV store — buckets, batch ops, atomic counters, range scans, JSON import/export |
| **SQL** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Full SQLite in pure Go — CSV/JSON/NDJSON import, multi-format query output, interactive shell |
| **Graph** | [dominikbraun/graph](https://github.com/dominikbraun/graph) | Directed & undirected weighted graphs — shortest path, DFS/BFS, topological sort, cycle detection, connected components |
| **Search** | SQLite FTS5 | Full-text indexing — scored queries, field return, pagination, bulk NDJSON indexing |
| **TS** | [tstorage](https://github.com/nakabonne/tstorage) | Time-series storage — labeled metrics, range queries, latest-point lookup, CSV export |

## Build from Source

```bash
git clone https://github.com/agi-now/bull.git && cd bull
./build.sh                # Linux / macOS
.\build.ps1               # Windows
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

For detailed usage of each command, see [`skills/bull/references/`](skills/bull/references/).

## AI Agent Skills

The `skills/` directory contains machine-readable skill definitions for each engine. An AI agent reads these files to decide **which engine to use** and **which commands to invoke** for any given task. Use `--format json` for structured output.

Skill files are automatically downloaded by the AI Agent when it installs Bull. If automatic download fails, manually copy the `skills/` directory from this repository into your project.

## Data Storage

All data persists under `--data-dir` (default `./bull_data`), organized by engine:

| Engine | Location | Format |
|--------|----------|--------|
| KV | `bull_data/kv/<name>.db` | bbolt B+tree |
| SQL | `bull_data/sql/<name>.db` | SQLite |
| Graph | `bull_data/graph/<name>.json` | JSON |
| Search | `bull_data/search/<name>.db` | SQLite FTS5 |
| TS | `bull_data/ts/<name>/` | tstorage WAL |

## License

MIT
