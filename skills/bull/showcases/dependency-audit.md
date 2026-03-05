# Showcase: Microservice Dependency Audit

Multi-engine collaboration: Graph models the service topology and runs graph analysis, SQL stores service metadata for relational queries, KV caches audit conclusions.

---

## User Prompt

> We have a set of microservices with dependencies in services.csv (format: from,to,weight).
> Another file service_info.csv records each service's team, language, and deploy_env.
> Help me do an architecture audit:
> 1. Are there any circular dependencies?
> 2. Find the critical path (shortest path from auth-svc to db-svc)
> 3. Which services are isolated (no dependencies at all)?
> 4. Count how many services each team owns
> 5. Generate a safe deployment order

---

## Agent Execution

### Step 1 — Import Dependency Graph (Graph)

```bash
bull graph import-csv infra services.csv
```

Check overview:

```bash
bull graph stats infra
```

### Step 2 — Import Service Metadata (SQL)

```bash
bull sql import infra services service_info.csv
```

### Step 3 — Cycle Detection (Graph)

```bash
bull graph has-cycle infra
```

If output is `true`, locate the cycle via connected components:

```bash
bull graph components infra
```

### Step 4 — Critical Path Analysis (Graph)

```bash
bull graph shortest-path infra auth-svc db-svc
```

Output looks like `auth-svc -> user-svc -> db-svc`.

Check reachability between two services:

```bash
bull graph has-path infra gateway db-svc
```

### Step 5 — Isolated Service Discovery (Graph + SQL)

Get all vertices in the graph:

```bash
bull graph vertices infra
```

Check degree for each vertex — degree 0 means isolated:

```bash
bull graph degree infra monitoring-svc
```

### Step 6 — Services Per Team (SQL)

```bash
bull sql query infra "
  SELECT team, COUNT(*) AS svc_count, GROUP_CONCAT(name) AS services
  FROM services
  GROUP BY team
  ORDER BY svc_count DESC
" --format json
```

### Step 7 — Deployment Order (Graph)

```bash
bull graph toposort infra
```

Output is ordered from no-dependency to most-dependent — a safe deployment sequence.

### Step 8 — Cache Audit Results (KV)

```bash
bull kv put audit:infra has_cycle 'true'
bull kv put audit:infra critical_path 'auth-svc -> user-svc -> db-svc'
bull kv put audit:infra isolated '["monitoring-svc"]'
bull kv put audit:infra deploy_order '["db-svc","user-svc","auth-svc","gateway"]'
bull kv put audit:infra timestamp '2026-03-05T10:30:00Z'
```

---

## Engines Used

| Engine | Purpose |
|--------|---------|
| Graph | Dependency modeling, cycle detection, shortest path, topological sort, connected components |
| SQL | Service metadata queries, team aggregation |
| KV | Audit result persistence |
