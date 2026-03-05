<p align="center">
  <img src="https://img.shields.io/badge/体积-~20MB-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/引擎-5个-green?style=flat-square" />
  <img src="https://img.shields.io/badge/命令-72+-orange?style=flat-square" />
  <img src="https://img.shields.io/badge/CGo-无依赖-brightgreen?style=flat-square" />
  <img src="https://img.shields.io/badge/平台-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square" />
</p>

<p align="center"><b>Bull</b> — 全能嵌入式引擎工具箱</p>

<p align="center">
  <a href="README.md">English</a>
</p>

---

五个数据引擎，一个静态二进制，零外部依赖。

**Bull** 将 KV 存储、SQL 数据库、图引擎、全文搜索和时序存储打包进一个约 20 MB 的 Go 可执行文件。它专为 **AI Agent 技能扩展** 设计——将二进制丢入任何受限环境，即刻获得本地数据处理能力，无需安装任何数据库服务。

## 为什么选择 Bull?

- **单文件部署，零依赖** — 不需要安装数据库服务，不需要配置运行时环境。复制一个文件即可运行。
- **5 引擎，72+ 命令** — KV、SQL、图、全文搜索、时序 — 覆盖 AI Agent 绝大多数数据处理场景。
- **CLI + HTTP API** — 每个引擎都同时提供命令行和 RESTful HTTP 接口（`bull serve`），方便任何语言和框架集成。
- **AI Agent 原生支持** — 内置 `skills/` 机器可读技能定义，AI Agent 可自主决策使用哪个引擎、调用哪个命令。
- **纯 Go，静态编译** — 无 CGo、无 Wasm、无动态库。一条命令交叉编译到 Linux / macOS / Windows。

## 速览

```
bull kv put config host 10.0.0.1          # 持久化键值对
bull sql query db "SELECT * FROM t"       # 完整 SQL（SQLite）
bull graph shortest-path g A B            # 图算法
bull search query idx "error timeout"     # 全文搜索
bull ts latest mon cpu --format json      # 时序指标
bull serve -p 2880                        # 启动 HTTP API 服务
```

## 引擎一览

| 引擎 | 底层库 | 能力 |
|------|--------|------|
| **KV** | [bbolt](https://github.com/etcd-io/bbolt) | B+tree 键值存储——桶管理、批量操作、原子计数器、范围扫描、JSON 导入导出 |
| **SQL** | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | 纯 Go 的完整 SQLite——CSV/JSON/NDJSON 导入、多格式查询输出、交互式 Shell |
| **Graph** | [dominikbraun/graph](https://github.com/dominikbraun/graph) | 有向/无向加权图——最短路径、DFS/BFS、拓扑排序、环检测、连通分量 |
| **Search** | [bleve](https://github.com/blevesearch/bleve) | 全文索引——评分查询、字段返回、分页、NDJSON 批量索引 |
| **TS** | [tstorage](https://github.com/nakabonne/tstorage) | 时序存储——带标签的指标、范围查询、最新值查询、CSV 导出 |

所有引擎均为**纯 Go 实现**——无 CGo、无 Wasm、无动态库。二进制完全静态编译。

## 安装

```bash
# 克隆并构建
git clone https://github.com/agi-now/bull.git && cd bull
go build -ldflags="-s -w" -o bull ./cmd/bull/

# 或使用构建脚本（自动注入版本号和构建时间）
./build.sh                # Linux / macOS
.\build.ps1               # Windows
```

交叉编译：

```bash
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o bull           ./cmd/bull/
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o bull           ./cmd/bull/
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bull.exe       ./cmd/bull/
```

## 使用方式

所有数据持久化在 `--data-dir`（默认 `./data`）下，按引擎分目录存储。

### KV 存储

```bash
# 单键读写
bull kv put cache token abc123
bull kv get cache token                        # abc123

# 批量操作
bull kv mput cache k1 v1 k2 v2 k3 v3          # put 3 pairs
bull kv mget cache k1 k2 k3                   # [{"key":"k1","value":"v1"}, ...]

# 原子计数器
bull kv incr stats page_views                  # 1
bull kv incr stats page_views 10               # 11

# 扫描与过滤
bull kv list cache --prefix k --format json
bull kv scan cache --start a --end z --format json

# 导入 / 导出
bull kv export cache                           # JSON 输出到 stdout
bull kv import cache -f backup.json
```

### SQL 数据库

```bash
# 导入数据
bull sql import analytics users users.csv                    # imported 1500 rows
bull sql import-ndjson analytics logs access.ndjson           # imported 8234 rows

# 查询
bull sql query analytics "SELECT city, COUNT(*) c FROM users GROUP BY city ORDER BY c DESC" --format json --limit 10

# 表结构检查
bull sql tables analytics
bull sql describe analytics users
bull sql schema analytics users

# 导出
bull sql export analytics "SELECT * FROM users" -o out.csv --format csv
bull sql export analytics "SELECT * FROM users WHERE age>30" --format json

# 交互式 Shell
bull sql shell analytics
```

### 图引擎

```bash
# 构建图
bull graph add-vertex deps auth --attr type=service
bull graph add-vertex deps db --attr type=database
bull graph add-edge deps auth db --weight 1

# 图算法
bull graph shortest-path deps auth db          # auth -> db
bull graph toposort deps                       # 拓扑排序
bull graph has-cycle deps                      # false
bull graph components deps                     # 连通分量

# 遍历
bull graph dfs deps auth
bull graph bfs deps auth
bull graph neighbors deps auth

# 批量导入
bull graph import-csv social edges.csv --undirected
bull graph stats social --undirected            # vertices: 42, edges: 87
```

### 全文搜索

```bash
# 建索引
bull search create articles
bull search bulk articles docs.ndjson           # indexed 500 documents

# 分页查询
bull search query articles "machine learning" --limit 10 --offset 0 --format json
bull search query articles "title:hello" --field title --field body --format json

# CRUD
bull search index articles doc1 '{"title":"Hello","body":"World"}'
bull search update articles doc1 '{"title":"Updated","body":"New content"}'
bull search get articles doc1
bull search delete articles doc1
```

### 时序存储

```bash
# 写入
bull ts write mon cpu 72.5 --label host=web-01
bull ts write mon cpu 68.3 --label host=web-01

# 读取
bull ts latest mon cpu --label host=web-01 --format json
bull ts query mon cpu --from 1700000000 --to 1700003600 --format json
bull ts count mon cpu

# 导出
bull ts export mon cpu -o metrics.csv
```

### HTTP API

所有引擎同样支持 HTTP 访问。启动服务后，通过 RESTful 接口操作任意引擎：

```bash
bull serve -p 2880                            # 在 2880 端口启动

curl localhost:2880/api/version
curl -X POST localhost:2880/api/kv/mydb/get -d '{"key":"mykey"}'
curl -X POST localhost:2880/api/sql/mydb/query -d '{"sql":"SELECT 1"}'
```

### 全局命令

```bash
bull version                                   # 版本、构建时间、Go 版本、系统架构
bull info                                      # 所有引擎的数据库汇总
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
      ├─ serve                               HTTP API 服务 (--port)
      ├─ version                             构建信息
      └─ info                                数据目录总览
```

**72+ 条命令**，横跨 5 个引擎 + HTTP API 服务 + 2 个全局命令。

## AI Agent 技能集成

`skills/` 目录包含机器可读的技能定义文件。AI Agent 读取这些文件来决定**用哪个引擎**、**调用哪些命令**。使用 `--format json` 获取结构化输出以便机器解析。

## 项目结构

```
bull/
├── cmd/bull/              CLI 入口（cobra）
│   ├── main.go            根命令、version、info
│   ├── cmd_kv.go          KV 子命令
│   ├── cmd_sql.go         SQL 子命令
│   ├── cmd_graph.go       Graph 子命令
│   ├── cmd_search.go      Search 子命令
│   ├── cmd_ts.go          TS 子命令
│   └── cmd_serve.go       HTTP API 服务
├── internal/
│   ├── config/            数据目录配置
│   ├── kv/                bbolt 封装
│   ├── sql/               SQLite 封装
│   ├── graph/             图算法封装
│   ├── search/            bleve 封装
│   ├── ts/                tstorage 封装
│   └── server/            HTTP API 处理器
├── skills/                AI Agent 技能定义
├── web/                   Web 前端（Vite + React）
├── build.sh / build.ps1   带版本注入的构建脚本
└── data/                  运行时存储（已 gitignore）
```

## 数据持久化

| 引擎 | 路径 | 格式 |
|------|------|------|
| KV | `data/kv/<name>.db` | bbolt B+tree |
| SQL | `data/sql/<name>.db` | SQLite |
| Graph | `data/graph/<name>.json` | JSON |
| Search | `data/search/<name>.bleve/` | bleve 索引 |
| TS | `data/ts/<name>/` | tstorage WAL |

## 平台支持

| 系统 | 架构 |
|------|------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## 开源协议

MIT
