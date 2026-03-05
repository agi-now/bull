# Showcase: E-Commerce Multi-Source Data Analysis Pipeline

Multi-engine collaboration: SQL handles data import and aggregation, Search builds a product search index, KV tracks pipeline state.

---

## User Prompt

> I have three files: orders.csv (orders), products.csv (products), customers.csv (customers).
> Help me analyze:
> 1. Total sales and order count per city
> 2. Top 10 products by revenue
> 3. Build a search index for products so I can search by name later
> 4. Track analysis progress so I can resume next time

---

## Agent Execution

### Step 1 — Import Data (SQL)

```bash
bull sql import shop orders orders.csv
bull sql import shop products products.csv
bull sql import shop customers customers.csv
```

Verify imports:

```bash
bull sql count shop orders
bull sql count shop products
bull sql count shop customers
```

### Step 2 — Track Pipeline State (KV)

```bash
bull kv put pipeline:shop step1_import '{"status":"done","tables":["orders","products","customers"]}'
```

### Step 3 — Sales by City (SQL JOIN + GROUP BY)

```bash
bull sql query shop "
  SELECT c.city,
         COUNT(o.id) AS order_count,
         SUM(CAST(o.amount AS REAL)) AS total_sales
  FROM orders o
  JOIN customers c ON o.customer_id = c.id
  GROUP BY c.city
  ORDER BY total_sales DESC
" --format json
```

### Step 4 — Top 10 Products (SQL JOIN + ORDER BY)

```bash
bull sql query shop "
  SELECT p.name,
         p.category,
         COUNT(o.id) AS sold,
         SUM(CAST(o.amount AS REAL)) AS revenue
  FROM orders o
  JOIN products p ON o.product_id = p.id
  GROUP BY p.id
  ORDER BY revenue DESC
  LIMIT 10
" --format json
```

### Step 5 — Export Products and Build Search Index (SQL + Search)

Export as NDJSON:

```bash
bull sql query shop "SELECT id AS _id, name, category, price FROM products" --format json > /tmp/products.ndjson
```

> Note: The agent may need to convert the JSON array to NDJSON (one object per line) via a script.

Create index and bulk import:

```bash
bull search create products
bull search bulk products /tmp/products.ndjson
```

Verify search:

```bash
bull search query products "electronics" --field name --field price --format json
```

### Step 6 — Update Pipeline State (KV)

```bash
bull kv put pipeline:shop step2_analysis '{"status":"done","results":["city_sales","top10_products"]}'
bull kv put pipeline:shop step3_search '{"status":"done","index":"products","doc_count":500}'
bull kv put pipeline:shop current_step '3'
```

### Step 7 — View Full Progress

```bash
bull kv list pipeline:shop --format json
```

---

## Engines Used

| Engine | Purpose |
|--------|---------|
| SQL | Data import, JOIN queries, aggregation analysis |
| Search | Product full-text search index |
| KV | Pipeline state persistence |
