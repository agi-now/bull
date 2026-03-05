<p align="center">
  <img src="https://img.shields.io/badge/体积-~8MB-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/引擎-5个-green?style=flat-square" />
  <img src="https://img.shields.io/badge/命令-72+-orange?style=flat-square" />
  <img src="https://img.shields.io/badge/CGo-无依赖-brightgreen?style=flat-square" />
  <img src="https://img.shields.io/badge/平台-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square" />
</p>

<p align="center"><b>Bull</b> — 全能嵌入式引擎工具箱<br/>数据不进 LLM 上下文，全部本地处理。Save your tokens and money.</p>

<p align="center">
  <a href="README.md">English</a>
</p>

---

五个数据引擎，一个静态二进制，零外部依赖。

**Bull** 将 KV 存储、SQL 数据库、图引擎、全文搜索和时序存储打包进一个约 8 MB 的 Go 可执行文件。它专为 **AI Agent 技能扩展** 设计——将二进制丢入任何受限环境，即刻获得本地数据处理能力，无需安装任何数据库服务。

## 要解决的问题

AI Agent 处理大规模数据（CSV、日志、文档等）时，传统做法是将数据加载到对话上下文中。这会导致**巨额 token 消耗**和**响应缓慢**。一个 10 MB 的 CSV 文件光读取就要消耗几十万 token，更别说分析。

Bull 的做法是**将数据处理卸载到本地引擎**。Agent 将数据导入 Bull，在本地执行查询和聚合，只将精简的结果返回到对话中。原始数据不进入 token 流。

```
不用 Bull：用户 → [10MB CSV 作为 token] → LLM → 回答         （几十万 token）
使用 Bull：用户 → Agent → bull sql import + query → LLM → 回答（几百 token）
```

## 为什么选择 Bull?

- **大幅降低 token 开销** — 数据留在本地，只有查询结果进入对话。处理百万行数据无需消耗 token。
- **单文件部署，零依赖** — 不需要安装数据库服务，不需要配置运行时环境。复制一个文件即可运行。
- **5 引擎，72+ 命令** — KV、SQL、图、全文搜索、时序 — 覆盖 AI Agent 绝大多数数据处理场景。
- **CLI 优先** — 每个引擎都通过简洁的命令行接口操作，方便任何语言和框架集成。
- **AI Agent 原生支持** — 内置 `skills/` 机器可读技能定义，AI Agent 可自主决策使用哪个引擎、调用哪个命令。
- **纯 Go，静态编译** — 无 CGo、无动态库。一条命令交叉编译到 Linux / macOS / Windows。

## 速览

```
bull kv put config host 10.0.0.1          # 持久化键值对
bull sql query db "SELECT * FROM t"       # 完整 SQL（SQLite）
bull graph shortest-path g A B            # 图算法
bull search query idx "error timeout"     # 全文搜索
bull ts latest mon cpu --format json      # 时序指标
```

## 示例：60 秒完成故障排查

AI Agent 收到服务器日志和服务依赖图，在本地完成导入、分析、搜索、追踪——原始数据零 token 消耗。

```bash
# 1. 导入 5 万条访问日志到 SQL
bull sql import-ndjson incident access access.ndjson        # imported 50000 rows

# 2. 找到报错最多的服务
bull sql query incident "SELECT service, COUNT(*) c FROM access WHERE level='ERROR' GROUP BY service ORDER BY c DESC LIMIT 5" --format json

# 3. 建搜索索引，定位根因
bull search create logs
bull search bulk logs access.ndjson
bull search query logs "connection refused port 5432" --field service --field message --limit 5 --format json

# 4. 导入服务依赖图，追踪爆炸半径
bull graph import-csv incident deps.csv
bull graph shortest-path incident api-gateway db-primary
bull graph bfs incident db-primary                          # 所有受影响的下游服务

# 5. 记录事件时间线
bull ts write incident_metrics error_rate 127 --label service=db-primary
bull ts write incident_metrics error_rate 3 --label service=db-primary  # 修复后

# 6. 保存结论供复盘
bull kv put incident:2026-03-05 root_cause '{"service":"db-primary","issue":"connection pool exhausted"}'
bull kv put incident:2026-03-05 blast_radius '["api-gateway","user-svc","payment-svc"]'
```

5 个引擎，12 条命令，一个二进制——Agent 处理了 5 万条日志，没有一行数据进入 LLM 上下文。

## 下载

从 [**GitHub Releases**](https://github.com/agi-now/bull/releases/latest) 获取对应平台的最新二进制：

| 平台 | 文件名 |
|------|--------|
| Linux amd64 | `bull-linux-amd64` |
| Linux arm64 | `bull-linux-arm64` |
| macOS amd64 | `bull-darwin-amd64` |
| macOS arm64 | `bull-darwin-arm64` |
| Windows amd64 | `bull-windows-amd64.exe` |

下载后重命名为 `bull`（Windows 为 `bull.exe`），放入 PATH 即可。

## 引擎一览

| 引擎 | 底层库 | 能力 |
|------|--------|------|
| **KV** | [bbolt](https://github.com/etcd-io/bbolt) | B+tree 键值存储——桶管理、批量操作、原子计数器、范围扫描、JSON 导入导出 |
| **SQL** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | 纯 Go 的完整 SQLite——CSV/JSON/NDJSON 导入、多格式查询输出、交互式 Shell |
| **Graph** | [dominikbraun/graph](https://github.com/dominikbraun/graph) | 有向/无向加权图——最短路径、DFS/BFS、拓扑排序、环检测、连通分量 |
| **Search** | SQLite FTS5 | 全文索引——评分查询、字段返回、分页、NDJSON 批量索引 |
| **TS** | [tstorage](https://github.com/nakabonne/tstorage) | 时序存储——带标签的指标、范围查询、最新值查询、CSV 导出 |

## 从源码编译

```bash
git clone https://github.com/agi-now/bull.git && cd bull
./build.sh                # Linux / macOS
.\build.ps1               # Windows
```

## 命令总览

```
bull ─┬─ kv ─────┬─ put / get / del          单键操作
      │          ├─ mget / mput              批量操作
      │          ├─ list / scan              范围查询 (--format tsv|json)
      │          ├─ exists / count           检查
      │          ├─ incr / decr              原子计数器
      │          ├─ buckets                  列出桶
      │          ├─ export / import          JSON 导入导出
      │          ├─ drop / drop-bucket       清理
      │          └─ dbs                      列出数据库
      │
      ├─ sql ────┬─ exec / query             DDL/DML 和 SELECT (--format, --limit)
      │          ├─ exec-file                执行 .sql 文件
      │          ├─ tables / schema / describe / count
      │          ├─ import / import-json / import-ndjson
      │          ├─ export                   导出 CSV 或 JSON (--format)
      │          ├─ shell                    交互式 REPL
      │          ├─ drop                     清理
      │          └─ dbs                      列出数据库
      │
      ├─ graph ──┬─ add-vertex / add-edge    构建图 (--attr, --weight)
      │          ├─ del-vertex / del-edge    删除
      │          ├─ vertices / edges         列出
      │          ├─ neighbors / degree / attrs
      │          ├─ shortest-path / has-path 路径查找
      │          ├─ dfs / bfs               遍历
      │          ├─ components / toposort / has-cycle
      │          ├─ stats                    统计
      │          ├─ import-csv / export      导入导出
      │          ├─ drop                     清理
      │          └─ dbs                      列出图
      │
      ├─ search ─┬─ create                   创建索引
      │          ├─ index / bulk             添加文档
      │          ├─ query                    搜索 (--field, --limit, --offset)
      │          ├─ get / update / delete    CRUD
      │          ├─ info                     索引统计
      │          ├─ drop                     清理
      │          └─ dbs                      列出索引
      │
      ├─ ts ─────┬─ write / bulk             写入
      │          ├─ query                    范围查询 (--from, --to, --label)
      │          ├─ latest / count           检查
      │          ├─ export                   CSV 导出
      │          ├─ drop                     清理
      │          └─ dbs                      列出数据库
      │
      ├─ version                             构建信息
      └─ info                                数据目录总览
```

各命令的详细用法见 [`skills/bull/references/`](skills/bull/references/)。

## AI Agent 技能集成

`skills/` 目录包含机器可读的技能定义文件。AI Agent 读取这些文件来决定**用哪个引擎**、**调用哪些命令**。使用 `--format json` 获取结构化输出以便机器解析。

技能文件会在 AI Agent 安装 Bull 时自动下载。如果自动下载失败，可手动将本仓库的 `skills/` 目录拷贝到你的项目中。

## 数据持久化

所有数据持久化在 `--data-dir`（默认 `./bull_data`）下，按引擎分目录存储：

| 引擎 | 路径 | 格式 |
|------|------|------|
| KV | `bull_data/kv/<name>.db` | bbolt B+tree |
| SQL | `bull_data/sql/<name>.db` | SQLite |
| Graph | `bull_data/graph/<name>.json` | JSON |
| Search | `bull_data/search/<name>.db` | SQLite FTS5 |
| TS | `bull_data/ts/<name>/` | tstorage WAL |

## 开源协议

MIT
