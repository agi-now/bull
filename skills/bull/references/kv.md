# KV Engine Reference

Embedded key-value store powered by bbolt (B+tree). Fast, persistent, bucket-aware.

## Commands

### put
```
bull kv put <db> <key> <value> [--bucket <name>]
```
Store a key-value pair. Bucket defaults to "default".
```bash
bull kv put cache session_token abc123 --bucket sessions
```

### get
```
bull kv get <db> <key> [--bucket <name>]
```
Retrieve the value for a key. Prints the raw value. Exit code 1 if not found.

### del
```
bull kv del <db> <key> [--bucket <name>]
```
Delete a key.

### mget
```
bull kv mget <db> <key1> <key2> ... [--bucket <name>]
```
Get multiple keys in one call. Returns JSON array `[{"key":"k","value":"v"},...]`.
```bash
bull kv mget cache token1 token2 token3
```

### mput
```
bull kv mput <db> <key1> <val1> <key2> <val2> ... [--bucket <name>]
```
Put multiple key-value pairs in one atomic batch.
```bash
bull kv mput cache k1 v1 k2 v2 k3 v3
```

### list
```
bull kv list <db> [--bucket <name>] [--prefix <prefix>] [--format tsv|json]
```
List all key-value pairs, optionally filtered by prefix.
- `--format tsv` (default): `<key>\t<value>` per line
- `--format json`: JSON array `[{"key":"k","value":"v"},...]`
```bash
bull kv list myapp --prefix user: --format json
```

### scan
```
bull kv scan <db> [--bucket <name>] [--start <key>] [--end <key>] [--format tsv|json]
```
Range scan keys within [start, end) bounds.
```bash
bull kv scan myapp --start "a" --end "m" --format json
```

### exists
```
bull kv exists <db> <key> [--bucket <name>]
```
Check if a key exists. Prints `true` or `false`.

### count
```
bull kv count <db> [--bucket <name>]
```
Count keys in a bucket. Prints a single integer.

### incr
```
bull kv incr <db> <key> [delta] [--bucket <name>]
```
Atomically increment a numeric value. Default delta=1. Prints the new value.
```bash
bull kv incr counters page_views
```

### decr
```
bull kv decr <db> <key> [delta] [--bucket <name>]
```
Atomically decrement. Default delta=1. Prints the new value.
```bash
bull kv decr counters page_views 5
```

### buckets
```
bull kv buckets <db>
```
List all bucket names. One per line.

### export
```
bull kv export <db> [--bucket <name>]
```
Export all key-value pairs as JSON array.

### import
```
bull kv import <db> -f <file.json> [--bucket <name>]
```
Import key-value pairs from a JSON file. Format: `[{"key":"k","value":"v"},...]`.

### drop
```
bull kv drop <db>
```
Delete an entire database file.

### drop-bucket
```
bull kv drop-bucket <db> <bucket>
```
Delete a bucket from a database.

### dbs
```
bull kv dbs
```
List all KV databases. One name per line.
