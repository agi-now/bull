# Showcase: Project Task Scheduling & Dependency Tracking

Multi-engine collaboration: Graph models task dependencies and generates execution order, SQL stores task metadata for queries, KV tracks completion state, TS records progress metrics.

---

## User Prompt

> I'm managing a software release with tasks in tasks.csv (columns: id, name, assignee, estimate_hours, priority).
> Task dependencies are in deps.csv (columns: from, to — meaning "from" must finish before "to" starts).
> Help me:
> 1. Find a valid execution order
> 2. Check for circular dependencies
> 3. Find the critical path from "design" to "release"
> 4. Show workload per assignee
> 5. Track task completion

---

## Agent Execution

### Step 1 — Import Task Dependencies (Graph)

```bash
bull graph import-csv release deps.csv
```

Overview:

```bash
bull graph stats release
```

### Step 2 — Import Task Metadata (SQL)

```bash
bull sql import release tasks tasks.csv
```

### Step 3 — Cycle Detection (Graph)

```bash
bull graph has-cycle release
```

### Step 4 — Execution Order (Graph)

```bash
bull graph toposort release
```

Output gives a valid ordering that respects all dependencies.

### Step 5 — Critical Path (Graph)

```bash
bull graph shortest-path release design release
```

Check what blocks a specific task:

```bash
bull graph dfs release design
```

### Step 6 — Workload Per Assignee (SQL)

```bash
bull sql query release "
  SELECT assignee,
         COUNT(*) AS task_count,
         SUM(CAST(estimate_hours AS REAL)) AS total_hours
  FROM tasks
  GROUP BY assignee
  ORDER BY total_hours DESC
" --format json
```

High-priority tasks not yet assigned:

```bash
bull sql query release "
  SELECT id, name, priority, estimate_hours
  FROM tasks
  WHERE priority = 'high' AND (assignee IS NULL OR assignee = '')
" --format json
```

### Step 7 — Track Completion (KV)

Mark tasks as done:

```bash
bull kv put release:status design '{"done":true,"finished_at":"2026-03-01T09:00:00Z"}'
bull kv put release:status backend '{"done":true,"finished_at":"2026-03-03T17:00:00Z"}'
bull kv put release:status testing '{"done":false}'
```

View all statuses:

```bash
bull kv list release:status --format json
```

### Step 8 — Progress Metrics (TS)

```bash
bull ts write release_metrics tasks_completed 5
bull ts write release_metrics tasks_remaining 12
bull ts write release_metrics hours_burned 40
```

Check latest progress:

```bash
bull ts latest release_metrics tasks_completed --format json
bull ts latest release_metrics tasks_remaining --format json
```

---

## Engines Used

| Engine | Purpose |
|--------|---------|
| Graph | Task dependency modeling, cycle detection, topological sort, critical path |
| SQL | Task metadata queries, workload aggregation |
| KV | Task completion state tracking |
| TS | Progress metrics over time |
