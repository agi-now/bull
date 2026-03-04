# Graph Engine Reference

Embedded graph data structure and algorithms (pure Go). Supports directed/undirected weighted graphs with JSON persistence.

Default: directed. Use `--undirected` to switch.

## Commands

### add-vertex
```
bull graph add-vertex <db> <id> [--attr key=value ...] [--undirected]
```
Add a vertex with optional attributes.
```bash
bull graph add-vertex deps svc-auth --attr type=service --attr team=platform
```

### add-edge
```
bull graph add-edge <db> <from> <to> [--weight N] [--attr key=value ...] [--undirected]
```
Add a weighted edge. Both vertices must exist.
```bash
bull graph add-edge deps svc-auth svc-db --weight 5
```

### del-vertex
```
bull graph del-vertex <db> <id> [--undirected]
```
Remove a vertex (must have no edges).

### del-edge
```
bull graph del-edge <db> <from> <to> [--undirected]
```
Remove an edge.

### vertices
```
bull graph vertices <db> [--undirected]
```
List all vertex IDs. One per line.

### edges
```
bull graph edges <db> [--undirected]
```
List all edges as `from -> to (weight: N)`.

### neighbors
```
bull graph neighbors <db> <vertex> [--undirected]
```
List direct neighbors. One per line.

### degree
```
bull graph degree <db> <vertex> [--undirected]
```
Get degree of a vertex. Single integer.

### attrs
```
bull graph attrs <db> <vertex> [--undirected]
```
Show vertex attributes as `key=value` per line.

### shortest-path
```
bull graph shortest-path <db> <from> <to> [--undirected]
```
Find shortest weighted path (Dijkstra). Prints `A -> B -> C`.
```bash
bull graph shortest-path deps svc-auth svc-cache
```

### has-path
```
bull graph has-path <db> <from> <to> [--undirected]
```
Check reachability. Prints `true` or `false`.

### dfs
```
bull graph dfs <db> <start> [--undirected]
```
Depth-first traversal. One vertex per line.

### bfs
```
bull graph bfs <db> <start> [--undirected]
```
Breadth-first traversal. One vertex per line.

### stats
```
bull graph stats <db> [--undirected]
```
Show vertex and edge counts.

### components
```
bull graph components <db> [--undirected]
```
Find connected components (SCC for directed, BFS for undirected).

### toposort
```
bull graph toposort <db>
```
Topological sort (DAG only). Errors if cycle detected.

### has-cycle
```
bull graph has-cycle <db>
```
Check for cycles. Prints `true` or `false`.

### import-csv
```
bull graph import-csv <db> <file.csv> [--undirected]
```
Bulk import edges from CSV (`from,to[,weight]`). Auto-creates vertices.
```bash
bull graph import-csv social friends.csv --undirected
```

### export
```
bull graph export <db> [--undirected]
```
Export full graph as JSON.

### drop
```
bull graph drop <db>
```
Delete a graph file.

### dbs
```
bull graph dbs
```
List all graph files.

## HTTP API Endpoints

| Method | Path | Body |
|--------|------|------|
| GET | `/api/graph/dbs` | — |
| POST | `/api/graph/{db}/add-vertex` | `{"id","attrs?","undirected?"}` |
| POST | `/api/graph/{db}/add-edge` | `{"from","to","weight?","attrs?","undirected?"}` |
| POST | `/api/graph/{db}/del-vertex` | `{"id","undirected?"}` |
| POST | `/api/graph/{db}/del-edge` | `{"from","to","undirected?"}` |
| POST | `/api/graph/{db}/vertices` | `{"undirected?"}` |
| POST | `/api/graph/{db}/edges` | `{"undirected?"}` |
| POST | `/api/graph/{db}/neighbors` | `{"id","undirected?"}` |
| POST | `/api/graph/{db}/degree` | `{"id","undirected?"}` |
| POST | `/api/graph/{db}/attrs` | `{"id","undirected?"}` |
| POST | `/api/graph/{db}/shortest-path` | `{"from","to","undirected?"}` |
| POST | `/api/graph/{db}/has-path` | `{"from","to","undirected?"}` |
| POST | `/api/graph/{db}/dfs` | `{"start","undirected?"}` |
| POST | `/api/graph/{db}/bfs` | `{"start","undirected?"}` |
| POST | `/api/graph/{db}/stats` | `{"undirected?"}` |
| POST | `/api/graph/{db}/components` | `{"undirected?"}` |
| POST | `/api/graph/{db}/toposort` | `{"undirected?"}` |
| POST | `/api/graph/{db}/has-cycle` | `{"undirected?"}` |
| POST | `/api/graph/{db}/export` | `{"undirected?"}` |
| DELETE | `/api/graph/{db}` | — |
