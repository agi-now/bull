# Time-Series Engine Reference

Embedded time-series storage engine. Write and query timestamped metric data with labels.

Timestamps are Unix seconds. Query defaults to last 1 hour if `--from`/`--to` omitted.

## Commands

### write
```
bull ts write <db> <metric> <value> [--time <unix>] [--label key=value ...]
```
Write a single data point. Timestamp defaults to now.
```bash
bull ts write monitoring cpu_usage 75.5 --label host=server1 --label region=us-east
```

### bulk
```
bull ts bulk <db> <ndjson-file>
```
Bulk write from NDJSON. Format: `{"metric":"name","value":N,"timestamp":T,"labels":{"k":"v"}}`.

### query
```
bull ts query <db> <metric> [--from <unix>] [--to <unix>] [--label key=value ...] [--format table|csv|json]
```
Query data points in a time range.
```bash
bull ts query monitoring cpu_usage --from 1700000000 --to 1700003600 --label host=server1 --format json
```

### latest
```
bull ts latest <db> <metric> [--label key=value ...] [--format table|json]
```
Get the most recent data point.
```bash
bull ts latest monitoring cpu_usage --label host=server1 --format json
```

### count
```
bull ts count <db> <metric> [--from <unix>] [--to <unix>] [--label key=value ...]
```
Count data points in a time range. Single integer.
```bash
bull ts count monitoring cpu_usage --from 1700000000 --to 1700003600 --label host=server1
```

### export
```
bull ts export <db> <metric> [--from <unix>] [--to <unix>] [--label key=value ...] [-o <file>]
```
Export as CSV (`timestamp,value`).
```bash
bull ts export monitoring cpu_usage --label host=server1 -o cpu_report.csv
```

### drop
```
bull ts drop <db>
```
Delete a time-series database.

### dbs
```
bull ts dbs
```
List all time-series databases.

## HTTP API Endpoints

Fields ending with `?` are optional. CLI-only commands not available via HTTP: `bulk`.

| Method | Path | Body |
|--------|------|------|
| GET | `/api/ts/dbs` | — |
| POST | `/api/ts/{db}/write` | `{"metric","value","timestamp?","labels?"}` |
| POST | `/api/ts/{db}/query` | `{"metric","from?","to?","labels?"}` |
| POST | `/api/ts/{db}/latest` | `{"metric","labels?"}` |
| POST | `/api/ts/{db}/count` | `{"metric","from?","to?","labels?"}` |
| POST | `/api/ts/{db}/export` | `{"metric","from?","to?","labels?"}` (returns CSV) |
| DELETE | `/api/ts/{db}` | — |
