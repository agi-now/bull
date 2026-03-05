import { useState, useEffect } from 'react'
import { type Lang, t } from './i18n'

function App() {
  const [lang, setLang] = useState<Lang>('en')
  const [dark, setDark] = useState(() => window.matchMedia('(prefers-color-scheme: dark)').matches)

  useEffect(() => {
    document.documentElement.className = dark ? 'dark' : ''
  }, [dark])

  const T = (key: Parameters<typeof t>[1]) => t(lang, key)

  return (
    <>
      <nav>
        <div className="container">
          <span className="nav-logo">Bull</span>
          <div className="nav-right">
            <a href="https://github.com/agi-now/bull" target="_blank" rel="noopener">GitHub</a>
            <button className="toggle-btn" onClick={() => setLang(l => l === 'en' ? 'zh' : 'en')}>
              {lang === 'en' ? '中文' : 'EN'}
            </button>
            <button className="toggle-btn" onClick={() => setDark(d => !d)}>
              {dark ? '☀' : '☾'}
            </button>
          </div>
        </div>
      </nav>

      <div className="hero">
        <div className="container">
          <h1>{T('hero_title')}</h1>
          <p>{T('hero_sub')}</p>
          <div className="hero-links">
            <a href="https://github.com/agi-now/bull" target="_blank" rel="noopener">{T('hero_github')}</a>
            <a href="#install">{T('hero_get_started')}</a>
          </div>
          <pre>
{`$ bull kv put config host 10.0.0.1
$ bull sql import db users data.csv
$ bull sql query db "SELECT city, COUNT(*) FROM users GROUP BY city" --format json
$ bull graph shortest-path deps auth-svc cache-svc
$ bull search query articles "machine learning" --format json
$ bull serve --addr :9090`}
          </pre>
        </div>
      </div>

      <section id="engines">
        <div className="container">
          <h2>{T('engines_title')}</h2>
          <p>{T('engines_desc')}</p>
          <table>
            <thead>
              <tr>
                <th>Engine</th>
                <th>Powered by</th>
                <th>{lang === 'en' ? 'Capabilities' : '能力'}</th>
                <th>{lang === 'en' ? 'Cmds' : '命令数'}</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>kv</code></td>
                <td>bbolt</td>
                <td>{lang === 'en' ? 'Key-value, buckets, counters, batch, scan' : '键值、桶、计数器、批量、扫描'}</td>
                <td>17</td>
              </tr>
              <tr>
                <td><code>sql</code></td>
                <td>SQLite (pure Go)</td>
                <td>{lang === 'en' ? 'Full SQL, CSV/JSON import, export, shell' : '完整 SQL、CSV/JSON 导入、导出、Shell'}</td>
                <td>15</td>
              </tr>
              <tr>
                <td><code>graph</code></td>
                <td>dominikbraun/graph</td>
                <td>{lang === 'en' ? 'Shortest path, DFS/BFS, toposort, cycles, components' : '最短路径、DFS/BFS、拓扑排序、环检测、连通分量'}</td>
                <td>21</td>
              </tr>
              <tr>
                <td><code>search</code></td>
                <td>SQLite FTS5</td>
                <td>{lang === 'en' ? 'Full-text search, scoring, pagination, bulk' : '全文搜索、评分、分页、批量导入'}</td>
                <td>10</td>
              </tr>
              <tr>
                <td><code>ts</code></td>
                <td>tstorage</td>
                <td>{lang === 'en' ? 'Metrics, labels, range query, CSV export' : '指标、标签、范围查询、CSV 导出'}</td>
                <td>8</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section id="commands">
        <div className="container">
          <h2>{T('commands_title')}</h2>
          <p>{T('commands_desc')}</p>
          <pre className="cmd-tree">{cmdTree}</pre>
        </div>
      </section>

      <section id="skills">
        <div className="container">
          <h2>{T('skills_title')}</h2>
          <p>{T('skills_desc')}</p>
          <pre>
{`skills/bull/
├── SKILL.md              # ${lang === 'en' ? 'Frontmatter + instructions' : 'Frontmatter + 指令'}
└── references/
    ├── kv.md             # 17 commands
    ├── sql.md            # 15 commands
    ├── graph.md          # 21 commands
    ├── search.md         # 10 commands
    └── ts.md             #  8 commands`}
          </pre>
          <p style={{ marginTop: 16 }}>
            <strong>{lang === 'en' ? 'Progressive disclosure' : '渐进式加载'}:</strong>
          </p>
          <ul style={{ paddingLeft: 20, color: 'var(--muted)', fontSize: 15 }}>
            <li>{T('skills_p1')}</li>
            <li>{T('skills_p2')}</li>
            <li>{T('skills_p3')}</li>
          </ul>
        </div>
      </section>

      <section id="api">
        <div className="container">
          <h2>{T('api_title')}</h2>
          <p>{T('api_desc')}</p>
          <pre>
{`$ bull serve --addr :9090

# KV
POST /kv/put     {"db":"mydb","key":"k","value":"v"}
POST /kv/get     {"db":"mydb","key":"k"}

# SQL
POST /sql/query  {"db":"mydb","sql":"SELECT ...","format":"json"}

# Graph
POST /graph/shortest-path  {"db":"deps","from":"a","to":"b"}

# Search
POST /search/query  {"index":"articles","q":"keyword","limit":10}

# Time-Series
POST /ts/query  {"db":"mon","metric":"cpu","format":"json"}`}
          </pre>
        </div>
      </section>

      <section id="install">
        <div className="container">
          <h2>{T('install_title')}</h2>
          <pre>
{`${lang === 'en' ? '# Build from source' : '# 从源码构建'}
git clone https://github.com/agi-now/bull.git
cd bull
go build -ldflags="-s -w" -o bull ./cmd/bull/

${lang === 'en' ? '# Verify' : '# 验证'}
./bull version`}
          </pre>

          <h3 style={{ marginTop: 24, fontSize: 16, fontWeight: 600 }}>{T('platforms_title')}</h3>
          <div className="tags">
            {['linux/amd64', 'linux/arm64', 'darwin/amd64', 'darwin/arm64', 'windows/amd64'].map(p => (
              <span className="tag" key={p}>{p}</span>
            ))}
          </div>
        </div>
      </section>

      <footer>
        <div className="container">
          <p>{T('footer')}</p>
        </div>
      </footer>
    </>
  )
}

const cmdTree = `bull
├── kv        put get del mget mput list scan exists count
│             incr decr buckets export import drop drop-bucket dbs
├── sql       exec query exec-file tables schema describe count
│             import import-json import-ndjson export shell drop dbs
├── graph     add-vertex add-edge del-vertex del-edge vertices edges
│             neighbors degree attrs shortest-path has-path dfs bfs
│             stats components toposort has-cycle import-csv export drop dbs
├── search    create index bulk query get update delete info drop dbs
├── ts        write bulk query latest count export drop dbs
├── serve     --addr (HTTP JSON API)
├── version
└── info`

export default App
