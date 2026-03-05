# Showcase: 内容检索与监控指标追踪

多引擎协作：Search 索引文档做全文检索，TS 记录查询和访问指标，KV 做配置和热门缓存，SQL 做结构化统计。

---

## 用户 Prompt

> 我有一批技术博客文章（articles.ndjson，每行一个 JSON，包含 title、author、tags、body、published_at 字段）。
> 帮我：
> 1. 建个全文搜索索引
> 2. 搜一下关于 "distributed systems" 的文章
> 3. 每次搜索模拟记一次指标（搜索次数、耗时）
> 4. 统计每个 author 发了多少篇文章
> 5. 把搜索最多的关键词缓存起来

---

## Agent 执行过程

### Step 1 — 建索引并导入文档（Search）

```bash
bull search create articles
bull search bulk articles articles.ndjson
```

确认索引状态：

```bash
bull search info articles
```

### Step 2 — 全文搜索（Search）

```bash
bull search query articles "distributed systems" --field title --field author --field tags --limit 10 --format json
```

按字段精确搜索：

```bash
bull search query articles "tags:consensus" --field title --field author --limit 5 --format json
```

### Step 3 — 记录搜索指标（TS）

每次搜索后记录指标：

```bash
bull ts write content_metrics search_count 1 --label query="distributed systems"
bull ts write content_metrics search_latency_ms 45 --label query="distributed systems"
```

模拟多次搜索后查看趋势：

```bash
bull ts query content_metrics search_count --format json
bull ts latest content_metrics search_latency_ms --format json
```

### Step 4 — 结构化统计（SQL）

将文章数据导入 SQL 做聚合：

```bash
bull sql import-ndjson content articles articles.ndjson
```

按作者统计：

```bash
bull sql query content "
  SELECT author, COUNT(*) AS article_count
  FROM articles
  GROUP BY author
  ORDER BY article_count DESC
" --format json
```

按月份统计发文趋势：

```bash
bull sql query content "
  SELECT SUBSTR(published_at, 1, 7) AS month, COUNT(*) AS count
  FROM articles
  GROUP BY month
  ORDER BY month
" --format json
```

### Step 5 — 缓存热门关键词（KV）

```bash
bull kv put content:cache hot_keywords '["distributed systems","consensus","raft","kubernetes"]'
bull kv put content:cache last_updated '2026-03-05T10:00:00Z'
```

下次查询时先读缓存：

```bash
bull kv get content:cache hot_keywords
```

用计数器追踪每个关键词被搜索的次数：

```bash
bull kv incr content:counters "query:distributed systems"
bull kv incr content:counters "query:kubernetes"
bull kv incr content:counters "query:kubernetes"
bull kv list content:counters --prefix "query:" --format json
```

---

## 涉及引擎

| 引擎 | 用途 |
|------|------|
| Search | 全文索引、关键词搜索、字段过滤 |
| TS | 搜索次数和延迟指标记录与查询 |
| SQL | 文章元数据聚合统计 |
| KV | 热门关键词缓存、搜索计数器 |
