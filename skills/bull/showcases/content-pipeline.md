# Showcase: Content Search & Metrics Tracking

Multi-engine collaboration: Search indexes documents for full-text retrieval, TS records query and access metrics, KV handles config and hot caches, SQL performs structured aggregation.

---

## User Prompt

> I have a batch of tech blog articles (articles.ndjson, one JSON per line with title, author, tags, body, published_at fields).
> Help me:
> 1. Build a full-text search index
> 2. Search for articles about "distributed systems"
> 3. Record a metric for each search (count, latency)
> 4. Count how many articles each author has published
> 5. Cache the most searched keywords

---

## Agent Execution

### Step 1 — Create Index and Import Documents (Search)

```bash
bull search create articles
bull search bulk articles articles.ndjson
```

Check index status:

```bash
bull search info articles
```

### Step 2 — Full-Text Search (Search)

```bash
bull search query articles "distributed systems" --field title --field author --field tags --limit 10 --format json
```

Search by specific field:

```bash
bull search query articles "tags:consensus" --field title --field author --limit 5 --format json
```

### Step 3 — Record Search Metrics (TS)

Record metrics after each search:

```bash
bull ts write content_metrics search_count 1 --label query="distributed systems"
bull ts write content_metrics search_latency_ms 45 --label query="distributed systems"
```

View trends after multiple searches:

```bash
bull ts query content_metrics search_count --format json
bull ts latest content_metrics search_latency_ms --format json
```

### Step 4 — Structured Aggregation (SQL)

Import article data into SQL for aggregation:

```bash
bull sql import-ndjson content articles articles.ndjson
```

Count by author:

```bash
bull sql query content "
  SELECT author, COUNT(*) AS article_count
  FROM articles
  GROUP BY author
  ORDER BY article_count DESC
" --format json
```

Monthly publishing trend:

```bash
bull sql query content "
  SELECT SUBSTR(published_at, 1, 7) AS month, COUNT(*) AS count
  FROM articles
  GROUP BY month
  ORDER BY month
" --format json
```

### Step 5 — Cache Hot Keywords (KV)

```bash
bull kv put content:cache hot_keywords '["distributed systems","consensus","raft","kubernetes"]'
bull kv put content:cache last_updated '2026-03-05T10:00:00Z'
```

Read cache before next query:

```bash
bull kv get content:cache hot_keywords
```

Track search count per keyword with counters:

```bash
bull kv incr content:counters "query:distributed systems"
bull kv incr content:counters "query:kubernetes"
bull kv incr content:counters "query:kubernetes"
bull kv list content:counters --prefix "query:" --format json
```

---

## Engines Used

| Engine | Purpose |
|--------|---------|
| Search | Full-text indexing, keyword search, field filtering |
| TS | Search count and latency metric recording and querying |
| SQL | Article metadata aggregation |
| KV | Hot keyword cache, search counters |
