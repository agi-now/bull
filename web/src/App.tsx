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
      {/* NAV */}
      <nav>
        <div className="container">
          <span className="nav-logo">Bull</span>
          <ul className="nav-links">
            <li><a href="#features">{T('nav_features')}</a></li>
            <li><a href="#engines">{T('nav_engines')}</a></li>
            <li><a href="#commands">{T('nav_commands')}</a></li>
            <li><a href="#skills">{T('nav_skills')}</a></li>
            <li><a href="#start">{T('nav_get_started')}</a></li>
          </ul>
          <div className="nav-controls">
            <button className="toggle-btn" onClick={() => setLang(l => l === 'en' ? 'zh' : 'en')}>
              {lang === 'en' ? '中文' : 'EN'}
            </button>
            <button className="toggle-btn" onClick={() => setDark(d => !d)}>
              {dark ? '☀️' : '🌙'}
            </button>
          </div>
        </div>
      </nav>

      {/* HERO */}
      <section className="hero">
        <div className="container">
          <span className="hero-badge">{T('hero_badge')}</span>
          <h1>
            {T('hero_title_1')}<br />
            <span className="gradient">{T('hero_title_2')}</span><br />
            {T('hero_title_3')}
          </h1>
          <p>{T('hero_desc')}</p>
          <div className="hero-actions">
            <a href="#start" className="btn-primary">{T('hero_cta')}</a>
            <a href="https://github.com/agi-now/bull" target="_blank" rel="noopener" className="btn-secondary">{T('hero_cta2')}</a>
          </div>

          <div className="hero-code">
            <pre>{heroCode}</pre>
          </div>

          <div className="stats-row">
            <div className="stat-item">
              <div className="stat-num">~18 MB</div>
              <div className="stat-label">{lang === 'en' ? 'Binary Size' : '二进制体积'}</div>
            </div>
            <div className="stat-item">
              <div className="stat-num">5</div>
              <div className="stat-label">{lang === 'en' ? 'Engines' : '数据引擎'}</div>
            </div>
            <div className="stat-item">
              <div className="stat-num">72</div>
              <div className="stat-label">{lang === 'en' ? 'Commands' : '条命令'}</div>
            </div>
            <div className="stat-item">
              <div className="stat-num">0</div>
              <div className="stat-label">{lang === 'en' ? 'External Deps' : '外部依赖'}</div>
            </div>
            <div className="stat-item">
              <div className="stat-num">45</div>
              <div className="stat-label">{lang === 'en' ? 'AI Skill Prompts' : '技能 Prompt'}</div>
            </div>
          </div>
        </div>
      </section>

      {/* FEATURES */}
      <section className="alt" id="features">
        <div className="container">
          <div className="section-header">
            <h2>{T('feat_title')}</h2>
            <p>{T('feat_desc')}</p>
          </div>
          <div className="feat-grid">
            {[
              { icon: '📦', title: T('feat_1_title'), desc: T('feat_1_desc') },
              { icon: '🔧', title: T('feat_2_title'), desc: T('feat_2_desc') },
              { icon: '🦫', title: T('feat_3_title'), desc: T('feat_3_desc') },
              { icon: '🤖', title: T('feat_4_title'), desc: T('feat_4_desc') },
              { icon: '📄', title: T('feat_5_title'), desc: T('feat_5_desc') },
              { icon: '💾', title: T('feat_6_title'), desc: T('feat_6_desc') },
            ].map((f, i) => (
              <div className="feat-card" key={i}>
                <div className="feat-icon">{f.icon}</div>
                <h3>{f.title}</h3>
                <p>{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ENGINES */}
      <section id="engines">
        <div className="container">
          <div className="section-header">
            <h2>{T('engines_title')}</h2>
            <p>{T('engines_desc')}</p>
          </div>
          <div className="engine-grid">
            {[
              { name: T('engine_kv_name'), lib: T('engine_kv_lib'), desc: T('engine_kv_desc'), color: 'var(--engine-kv)', cmds: 17 },
              { name: T('engine_sql_name'), lib: T('engine_sql_lib'), desc: T('engine_sql_desc'), color: 'var(--engine-sql)', cmds: 15 },
              { name: T('engine_graph_name'), lib: T('engine_graph_lib'), desc: T('engine_graph_desc'), color: 'var(--engine-graph)', cmds: 21 },
              { name: T('engine_search_name'), lib: T('engine_search_lib'), desc: T('engine_search_desc'), color: 'var(--engine-search)', cmds: 11 },
              { name: T('engine_ts_name'), lib: T('engine_ts_lib'), desc: T('engine_ts_desc'), color: 'var(--engine-ts)', cmds: 8 },
            ].map((e, i) => (
              <div className="engine-card" key={i} style={{ '--engine-color': e.color } as React.CSSProperties}>
                <h3>{e.name}</h3>
                <div className="engine-lib">{e.lib}</div>
                <p>{e.desc}</p>
                <span className="engine-cmd-count">{e.cmds} {lang === 'en' ? 'commands' : '条命令'}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* COMMANDS */}
      <section className="alt" id="commands">
        <div className="container">
          <div className="section-header">
            <h2>{T('cmd_title')}</h2>
            <p>{T('cmd_desc')}</p>
          </div>
          <div className="cmd-block">
            <pre>{cmdTree}</pre>
          </div>
        </div>
      </section>

      {/* SKILLS */}
      <section id="skills">
        <div className="container">
          <div className="section-header">
            <h2>{T('skills_title')}</h2>
            <p>{T('skills_desc')}</p>
          </div>
          <div className="skills-grid">
            {[
              { num: '21', title: T('skills_card_1_title'), desc: T('skills_card_1_desc') },
              { num: '45', title: T('skills_card_2_title'), desc: T('skills_card_2_desc') },
              { num: 'JSON', title: T('skills_card_3_title'), desc: T('skills_card_3_desc') },
            ].map((s, i) => (
              <div className="skill-card" key={i}>
                <div className="skill-number">{s.num}</div>
                <h3>{s.title}</h3>
                <p>{s.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* GET STARTED */}
      <section className="alt" id="start">
        <div className="container">
          <div className="section-header">
            <h2>{T('get_started_title')}</h2>
          </div>
          <div className="get-started-code">
            <pre>
              <span className="comment">{lang === 'en' ? '# Clone & build' : '# 克隆并构建'}</span>{'\n'}
              <span className="cmd">git clone</span> https://github.com/agi-now/bull.git{'\n'}
              <span className="cmd">cd</span> bull{'\n'}
              <span className="cmd">go build</span> -ldflags=<span className="str">"-s -w"</span> -o bull ./cmd/bull/{'\n'}
              {'\n'}
              <span className="comment">{lang === 'en' ? '# Start using' : '# 开始使用'}</span>{'\n'}
              <span className="cmd">./bull</span> kv put mydb hello world{'\n'}
              <span className="cmd">./bull</span> sql query mydb <span className="str">"SELECT * FROM t"</span> --format json{'\n'}
              <span className="cmd">./bull</span> info
            </pre>
          </div>
        </div>
      </section>

      {/* PLATFORMS */}
      <section id="platforms">
        <div className="container">
          <div className="section-header">
            <h2>{T('platforms_title')}</h2>
            <p>{T('platforms_desc')}</p>
          </div>
          <div className="platform-badges">
            {['Linux amd64', 'Linux arm64', 'macOS amd64', 'macOS arm64', 'Windows amd64'].map(p => (
              <div className="platform-badge" key={p}>{p}</div>
            ))}
          </div>
        </div>
      </section>

      {/* FOOTER */}
      <footer>
        <div className="container">
          <div className="nav-logo">Bull</div>
          <p>{T('footer_desc')}</p>
          <p style={{ marginTop: 4 }}>{T('footer_license')}</p>
        </div>
      </footer>
    </>
  )
}

const heroCode = <>{
  [
    <span key="c1" className="comment"># KV — persistent key-value</span>,
    '\n',
    <span key="h1" className="cmd">bull kv put</span>, ' config host ', <span key="s1" className="str">10.0.0.1</span>,
    '\n',
    <span key="h2" className="cmd">bull kv mget</span>, ' config host port name',
    '\n\n',
    <span key="c2" className="comment"># SQL — full SQLite queries</span>,
    '\n',
    <span key="h3" className="cmd">bull sql import</span>, ' db users ', <span key="a1" className="arg">data.csv</span>,
    '\n',
    <span key="h4" className="cmd">bull sql query</span>, ' db ', <span key="s2" className="str">"SELECT city, COUNT(*) FROM users GROUP BY city"</span>, <span key="f1" className="flag"> --format json</span>,
    '\n\n',
    <span key="c3" className="comment"># Graph — algorithms on the fly</span>,
    '\n',
    <span key="h5" className="cmd">bull graph shortest-path</span>, ' deps auth cache',
    '\n',
    <span key="h6" className="cmd">bull graph toposort</span>, ' pipeline',
    '\n\n',
    <span key="c4" className="comment"># Search — full-text with scoring</span>,
    '\n',
    <span key="h7" className="cmd">bull search query</span>, ' articles ', <span key="s3" className="str">"machine learning"</span>, <span key="f2" className="flag"> --format json</span>,
    '\n\n',
    <span key="c5" className="comment"># Time-Series — metrics & export</span>,
    '\n',
    <span key="h8" className="cmd">bull ts latest</span>, ' monitoring cpu ', <span key="f3" className="flag">--format json</span>,
  ]
}</>

const cmdTree = `bull ─┬─ kv ─────┬─ put / get / del / mget / mput
      │          ├─ list / scan          (--format tsv|json)
      │          ├─ exists / count / incr / decr
      │          ├─ buckets / export / import
      │          └─ drop / drop-bucket / dbs
      │
      ├─ sql ────┬─ exec / query         (--format, --limit)
      │          ├─ exec-file / shell
      │          ├─ tables / schema / describe / count
      │          ├─ import / import-json / import-ndjson
      │          ├─ export               (--format csv|json)
      │          └─ drop / dbs
      │
      ├─ graph ──┬─ add-vertex / add-edge / del-vertex / del-edge
      │          ├─ vertices / edges / neighbors / degree / attrs
      │          ├─ shortest-path / has-path / dfs / bfs
      │          ├─ components / toposort / has-cycle
      │          ├─ stats / import-csv / export
      │          └─ drop / dbs
      │
      ├─ search ─┬─ create / index / bulk
      │          ├─ query                (--field, --limit, --offset)
      │          ├─ get / update / delete
      │          ├─ info / drop / dbs
      │
      ├─ ts ─────┬─ write / bulk
      │          ├─ query / latest / count / export
      │          └─ drop / dbs
      │
      ├─ version
      └─ info`

export default App
