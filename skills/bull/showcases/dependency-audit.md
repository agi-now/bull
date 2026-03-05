# Showcase: 微服务架构依赖审计

多引擎协作：Graph 建模服务拓扑并做图分析，SQL 存服务元数据做关联查询，KV 缓存审计结论。

---

## 用户 Prompt

> 我们有一组微服务，依赖关系在 services.csv 里（格式：from,to,weight）。
> 另外 service_info.csv 记录了每个服务的 team、language、deploy_env。
> 帮我做个架构审计：
> 1. 有没有循环依赖
> 2. 找出关键路径（auth-svc 到 db-svc 的最短路径）
> 3. 哪些服务是孤立的（没有任何依赖关系）
> 4. 按 team 统计各组负责了多少服务
> 5. 给出部署顺序

---

## Agent 执行过程

### Step 1 — 导入依赖图（Graph）

```bash
bull graph import-csv infra services.csv
```

查看概况：

```bash
bull graph stats infra
```

### Step 2 — 导入服务元数据（SQL）

```bash
bull sql import infra services service_info.csv
```

### Step 3 — 循环依赖检测（Graph）

```bash
bull graph has-cycle infra
```

如果输出 `true`，通过 SCC 定位循环所在的组件：

```bash
bull graph components infra
```

### Step 4 — 关键路径分析（Graph）

```bash
bull graph shortest-path infra auth-svc db-svc
```

输出形如 `auth-svc -> user-svc -> db-svc`。

也可以确认两个服务之间是否可达：

```bash
bull graph has-path infra gateway db-svc
```

### Step 5 — 孤立服务发现（Graph + SQL）

先拿到图中所有节点：

```bash
bull graph vertices infra
```

对每个节点查入度+出度，degree 为 0 的就是孤立服务：

```bash
bull graph degree infra monitoring-svc
```

或者用更高效的做法——从 SQL 侧对比：

```bash
bull sql query infra "
  SELECT s.name, s.team
  FROM services s
  WHERE s.name NOT IN (
    SELECT DISTINCT from_col FROM edges
    UNION
    SELECT DISTINCT to_col FROM edges
  )
" --format json
```

> 注意：这种方法需要先将边数据也导入 SQL，或者 agent 将 graph vertices 的结果与 SQL 做交叉比对。

### Step 6 — 按 team 统计服务数（SQL）

```bash
bull sql query infra "
  SELECT team, COUNT(*) AS svc_count, GROUP_CONCAT(name) AS services
  FROM services
  GROUP BY team
  ORDER BY svc_count DESC
" --format json
```

### Step 7 — 生成部署顺序（Graph）

```bash
bull graph toposort infra
```

输出从无依赖到有依赖的顺序，即安全的部署顺序。

### Step 8 — 缓存审计结论（KV）

```bash
bull kv put audit:infra has_cycle 'true'
bull kv put audit:infra critical_path 'auth-svc -> user-svc -> db-svc'
bull kv put audit:infra isolated '["monitoring-svc"]'
bull kv put audit:infra deploy_order '["db-svc","user-svc","auth-svc","gateway"]'
bull kv put audit:infra timestamp '2026-03-05T10:30:00Z'
```

---

## 涉及引擎

| 引擎 | 用途 |
|------|------|
| Graph | 依赖建模、循环检测、最短路径、拓扑排序、连通分量 |
| SQL | 服务元数据查询、按 team 聚合 |
| KV | 审计结论持久化 |
