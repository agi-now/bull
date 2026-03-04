# HTTP API Reference

Start the server:
```
bull serve [--port 2880]
```

All responses are JSON:
```json
{"ok": true, "data": ...}
{"ok": false, "error": "..."}
```

## Global Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/version` | Version, build time, Go version, OS/arch |
| GET | `/api/info` | Summary of all databases across engines |

## KV Endpoints (16)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/kv/dbs` | List databases |
| POST | `/api/kv/{db}/put` | Store key-value |
| POST | `/api/kv/{db}/get` | Get value |
| POST | `/api/kv/{db}/del` | Delete key |
| POST | `/api/kv/{db}/mget` | Batch get |
| POST | `/api/kv/{db}/mput` | Batch put |
| POST | `/api/kv/{db}/list` | List keys |
| POST | `/api/kv/{db}/scan` | Range scan |
| POST | `/api/kv/{db}/exists` | Check key |
| POST | `/api/kv/{db}/count` | Count keys |
| POST | `/api/kv/{db}/incr` | Increment |
| GET | `/api/kv/{db}/buckets` | List buckets |
| POST | `/api/kv/{db}/export` | Export JSON |
| POST | `/api/kv/{db}/import` | Import JSON |
| DELETE | `/api/kv/{db}` | Drop database |
| POST | `/api/kv/{db}/drop-bucket` | Drop bucket |

## SQL Endpoints (8)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/sql/dbs` | List databases |
| POST | `/api/sql/{db}/exec` | Execute SQL |
| POST | `/api/sql/{db}/query` | Query SQL |
| GET | `/api/sql/{db}/tables` | List tables |
| GET | `/api/sql/{db}/schema/{table}` | Table DDL |
| GET | `/api/sql/{db}/describe/{table}` | Column info |
| GET | `/api/sql/{db}/count/{table}` | Row count |
| DELETE | `/api/sql/{db}` | Drop database |

## Graph Endpoints (20)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/graph/dbs` | List graphs |
| POST | `/api/graph/{db}/add-vertex` | Add vertex |
| POST | `/api/graph/{db}/add-edge` | Add edge |
| POST | `/api/graph/{db}/del-vertex` | Remove vertex |
| POST | `/api/graph/{db}/del-edge` | Remove edge |
| POST | `/api/graph/{db}/vertices` | List vertices |
| POST | `/api/graph/{db}/edges` | List edges |
| POST | `/api/graph/{db}/neighbors` | Neighbors |
| POST | `/api/graph/{db}/degree` | Degree |
| POST | `/api/graph/{db}/attrs` | Vertex attrs |
| POST | `/api/graph/{db}/shortest-path` | Shortest path |
| POST | `/api/graph/{db}/has-path` | Reachability |
| POST | `/api/graph/{db}/dfs` | DFS traversal |
| POST | `/api/graph/{db}/bfs` | BFS traversal |
| POST | `/api/graph/{db}/stats` | Graph stats |
| POST | `/api/graph/{db}/components` | Components |
| POST | `/api/graph/{db}/toposort` | Topo sort |
| POST | `/api/graph/{db}/has-cycle` | Cycle check |
| POST | `/api/graph/{db}/export` | Export JSON |
| DELETE | `/api/graph/{db}` | Drop graph |

## Search Endpoints (9)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/search/dbs` | List indexes |
| POST | `/api/search/{idx}/create` | Create index |
| POST | `/api/search/{idx}/index` | Index document |
| POST | `/api/search/{idx}/query` | Search |
| GET | `/api/search/{idx}/get/{docID}` | Get document |
| POST | `/api/search/{idx}/update` | Update document |
| DELETE | `/api/search/{idx}/doc/{docID}` | Delete document |
| GET | `/api/search/{idx}/info` | Index info |
| DELETE | `/api/search/{idx}` | Drop index |

## TS Endpoints (7)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/ts/dbs` | List databases |
| POST | `/api/ts/{db}/write` | Write point |
| POST | `/api/ts/{db}/query` | Query range |
| POST | `/api/ts/{db}/latest` | Latest point |
| POST | `/api/ts/{db}/count` | Count points |
| POST | `/api/ts/{db}/export` | Export CSV |
| DELETE | `/api/ts/{db}` | Drop database |
