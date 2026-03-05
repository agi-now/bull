# Showcase: Server Log Analysis & Error Tracking

Multi-engine collaboration: SQL imports structured logs for aggregation, Search enables full-text log search, TS tracks error rates over time, KV stores alert thresholds and state.

---

## User Prompt

> I have server access logs in access.ndjson (each line: timestamp, level, service, message, status_code, latency_ms).
> Help me:
> 1. Import and find the top error-producing services
> 2. Search logs for specific error messages
> 3. Track error rate over time
> 4. Set up alert thresholds I can check later

---

## Agent Execution

### Step 1 — Import Logs (SQL)

```bash
bull sql import-ndjson logs access access.ndjson
```

Verify:

```bash
bull sql count logs access
bull sql describe logs access
```

### Step 2 — Error Analysis (SQL)

Top services by error count:

```bash
bull sql query logs "
  SELECT service, COUNT(*) AS error_count
  FROM access
  WHERE level = 'ERROR'
  GROUP BY service
  ORDER BY error_count DESC
  LIMIT 10
" --format json
```

Error distribution by status code:

```bash
bull sql query logs "
  SELECT status_code, COUNT(*) AS count
  FROM access
  WHERE CAST(status_code AS INTEGER) >= 400
  GROUP BY status_code
  ORDER BY count DESC
" --format json
```

Slow requests (latency > 1000ms):

```bash
bull sql query logs "
  SELECT service, message, latency_ms, timestamp
  FROM access
  WHERE CAST(latency_ms AS REAL) > 1000
  ORDER BY CAST(latency_ms AS REAL) DESC
  LIMIT 20
" --format json
```

### Step 3 — Full-Text Log Search (Search)

Build a search index for log messages:

```bash
bull search create server_logs
bull search bulk server_logs access.ndjson
```

Search for specific patterns:

```bash
bull search query server_logs "connection timeout" --field service --field message --limit 10 --format json
bull search query server_logs "out of memory" --field service --field message --limit 10 --format json
```

### Step 4 — Error Rate Metrics (TS)

Record hourly error counts per service:

```bash
bull ts write log_metrics error_count 23 --label service=api-gateway
bull ts write log_metrics error_count 5 --label service=user-svc
bull ts write log_metrics error_count 42 --label service=payment-svc
```

Track P99 latency:

```bash
bull ts write log_metrics p99_latency 850 --label service=api-gateway
bull ts write log_metrics p99_latency 1200 --label service=payment-svc
```

Query trends:

```bash
bull ts query log_metrics error_count --label service=payment-svc --format json
bull ts latest log_metrics p99_latency --label service=api-gateway --format json
```

### Step 5 — Alert Thresholds (KV)

```bash
bull kv put alerts:config error_rate_threshold '50'
bull kv put alerts:config latency_p99_threshold '1000'
bull kv put alerts:state payment-svc '{"status":"critical","error_count":42,"last_check":"2026-03-05T12:00:00Z"}'
bull kv put alerts:state api-gateway '{"status":"ok","error_count":23,"last_check":"2026-03-05T12:00:00Z"}'
```

Check alert state:

```bash
bull kv list alerts:state --format json
```

---

## Engines Used

| Engine | Purpose |
|--------|---------|
| SQL | Log import, error aggregation, slow request analysis |
| Search | Full-text log message search |
| TS | Error rate and latency metrics over time |
| KV | Alert thresholds and service status |
