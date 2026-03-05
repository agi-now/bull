# Showcase: 电商多源数据分析流水线

多引擎协作：SQL 做数据清洗和聚合，Search 建商品搜索索引，KV 存中间状态。

---

## 用户 Prompt

> 我有三个文件：orders.csv（订单）、products.csv（商品）、customers.csv（客户）。
> 帮我分析：
> 1. 每个城市的总销售额和订单数
> 2. 销售额 top 10 的商品
> 3. 把商品信息建个搜索索引，方便后面按名称搜
> 4. 记录一下分析进度，下次我可以接着来

---

## Agent 执行过程

### Step 1 — 导入数据（SQL）

```bash
bull sql import shop orders orders.csv
bull sql import shop products products.csv
bull sql import shop customers customers.csv
```

确认导入情况：

```bash
bull sql count shop orders
bull sql count shop products
bull sql count shop customers
```

### Step 2 — 记录流水线状态（KV）

```bash
bull kv put pipeline:shop step1_import '{"status":"done","tables":["orders","products","customers"]}'
```

### Step 3 — 按城市汇总销售额（SQL JOIN + GROUP BY）

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

### Step 4 — Top 10 商品（SQL JOIN + ORDER BY）

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

### Step 5 — 导出商品数据并建搜索索引（SQL + Search）

先导出为 NDJSON：

```bash
bull sql query shop "SELECT id AS _id, name, category, price FROM products" --format json > /tmp/products.ndjson
```

> 注意：需要将 JSON 数组转为 NDJSON（每行一个对象），agent 可用脚本处理。

创建索引并批量导入：

```bash
bull search create products
bull search bulk products /tmp/products.ndjson
```

验证搜索功能：

```bash
bull search query products "category:电子" --field name --field price --format json
```

### Step 6 — 更新流水线状态（KV）

```bash
bull kv put pipeline:shop step2_analysis '{"status":"done","results":["city_sales","top10_products"]}'
bull kv put pipeline:shop step3_search '{"status":"done","index":"products","doc_count":500}'
bull kv put pipeline:shop current_step '3'
```

### Step 7 — 查看完整进度

```bash
bull kv list pipeline:shop --format json
```

---

## 涉及引擎

| 引擎 | 用途 |
|------|------|
| SQL | 数据导入、JOIN 查询、聚合分析 |
| Search | 商品全文搜索索引 |
| KV | 流水线状态持久化 |
