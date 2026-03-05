# SQL Engine Reference

Embedded SQLite database (pure Go, no CGo). Full SQL support with CSV/JSON import and export.

## Commands

### exec
```
bull sql exec <db> <sql>
```
Execute DDL/DML (CREATE, INSERT, UPDATE, DELETE). Prints rows affected.
```bash
bull sql exec mydb "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)"
```

### query
```
bull sql query <db> <sql> [--format table|csv|json] [--limit N]
```
Run a SELECT query.
- `--format table` (default): ASCII table with borders
- `--format csv`: CSV with header row
- `--format json`: JSON array of objects
- `--limit N`: Wraps query in sub-select with LIMIT
```bash
bull sql query mydb "SELECT * FROM users WHERE age > 18" --format json --limit 100
```

### exec-file
```
bull sql exec-file <db> <file.sql>
```
Execute all SQL statements from a file.

### tables
```
bull sql tables <db>
```
List all table names. One per line.

### schema
```
bull sql schema <db> <table>
```
Show the CREATE TABLE DDL statement.

### describe
```
bull sql describe <db> <table>
```
Show column info (CID, name, type, notnull, default, pk) via PRAGMA table_info.

### count
```
bull sql count <db> <table>
```
Count rows. Prints a single integer.

### import (CSV)
```
bull sql import <db> <table> <file.csv>
```
Import CSV into a table. Auto-creates table with TEXT columns if needed.
```bash
bull sql import analytics users users.csv
```

### import-json
```
bull sql import-json <db> <table> <file.json>
```
Import a JSON array of objects. Auto-creates table.

### import-ndjson
```
bull sql import-ndjson <db> <table> <file.ndjson>
```
Import newline-delimited JSON (one object per line). Auto-creates table.

### export
```
bull sql export <db> <sql> [-o <file>] [--format csv|json]
```
Export query results to stdout or file.

### shell
```
bull sql shell <db>
```
Interactive SQL REPL. Meta-commands: `.tables`, `.schema <table>`, `.count <table>`, `.quit`.

### drop
```
bull sql drop <db>
```
Delete a database file.

### dbs
```
bull sql dbs
```
List all SQL databases.
