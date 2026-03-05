export type Lang = 'en' | 'zh'

const dict = {
  en: {
    nav_features: 'Features',
    nav_engines: 'Engines',
    nav_commands: 'Commands',
    nav_skills: 'AI Skills',
    nav_get_started: 'Get Started',

    hero_badge: '~8 MB · 5 Engines · 72+ Commands · Zero Dependencies',
    hero_title_1: 'One Binary.',
    hero_title_2: 'Five Engines.',
    hero_title_3: 'Infinite Possibilities.',
    hero_desc: 'Bull packs KV, SQL, Graph, Full-Text Search, and Time-Series into a single static Go binary. Purpose-built for AI Agent skill extensions — download once, use everywhere.',
    hero_cta: 'Get Started',
    hero_cta2: 'View on GitHub',

    feat_title: 'Why Bull?',
    feat_desc: 'Everything an AI Agent needs for local data processing, in one portable binary.',
    feat_1_title: 'Single Binary',
    feat_1_desc: 'One ~8 MB executable. No runtime, no containers, no package managers. Just copy and run.',
    feat_2_title: 'Five Engines',
    feat_2_desc: 'KV store, SQL database, graph engine, full-text search, and time-series — all embedded.',
    feat_3_title: 'Pure Go',
    feat_3_desc: 'Zero CGo dependencies. Fully static compilation. Cross-compile to any OS/arch in seconds.',
    feat_4_title: 'AI-Native',
    feat_4_desc: 'Ships with YAML skill definitions. AI agents read them to pick the right engine and commands.',
    feat_5_title: 'JSON Output',
    feat_5_desc: 'Every query command supports --format json for machine-readable structured output.',
    feat_6_title: 'Persistent',
    feat_6_desc: 'All data survives restarts. Each engine writes to its own directory under --data-dir.',

    engines_title: 'Five Engines, One Toolkit',
    engines_desc: 'Each engine is carefully selected for minimal footprint and maximum capability.',
    engine_kv_name: 'KV Store',
    engine_kv_lib: 'bbolt · ~1 MB',
    engine_kv_desc: 'B+tree key-value storage with buckets, batch operations, atomic counters, range scans, and JSON import/export.',
    engine_sql_name: 'SQL Database',
    engine_sql_lib: 'SQLite (pure Go) · ~7 MB',
    engine_sql_desc: 'Full SQLite — CSV/JSON/NDJSON import, multi-format export, interactive shell, PRAGMA introspection.',
    engine_graph_name: 'Graph Engine',
    engine_graph_lib: 'dominikbraun/graph · ~1 MB',
    engine_graph_desc: 'Weighted directed & undirected graphs — shortest path, DFS/BFS, topological sort, cycle detection, connected components.',
    engine_search_name: 'Full-Text Search',
    engine_search_lib: 'SQLite FTS5',
    engine_search_desc: 'Index JSON documents. Query with scoring, field return, pagination. Bulk ingest from NDJSON.',
    engine_ts_name: 'Time-Series',
    engine_ts_lib: 'tstorage · ~0.5 MB',
    engine_ts_desc: 'Write labeled metrics with timestamps. Range queries, latest-point lookup, counting, CSV export.',

    cmd_title: '72 Commands at Your Fingertips',
    cmd_desc: 'Every engine exposes a consistent CLI interface. Compose them freely in scripts and agent workflows.',

    skills_title: 'AI Agent Ready',
    skills_desc: 'Bull ships with machine-readable YAML skill definitions that tell AI agents exactly what to do.',
    skills_card_1_title: '21 Decision Scenarios',
    skills_card_1_desc: 'A decision matrix maps real-world tasks to the right engine and command sequence.',
    skills_card_2_title: '45 Example Prompts',
    skills_card_2_desc: 'Pre-built prompts demonstrate how to trigger each capability from natural language.',
    skills_card_3_title: 'Structured Output',
    skills_card_3_desc: 'Use --format json everywhere. Agents parse results directly — no regex scraping.',

    get_started_title: 'Get Started in 30 Seconds',

    platforms_title: 'Run Everywhere',
    platforms_desc: 'Pure Go. Static binary. No CGo. Cross-compile to any target.',

    footer_desc: 'All-in-One Embedded Engine Toolkit for AI Agents.',
    footer_license: 'MIT License',
  },

  zh: {
    nav_features: '特性',
    nav_engines: '引擎',
    nav_commands: '命令',
    nav_skills: 'AI 技能',
    nav_get_started: '快速开始',

    hero_badge: '~8 MB · 5 大引擎 · 72+ 条命令 · 零外部依赖',
    hero_title_1: '一个二进制。',
    hero_title_2: '五大引擎。',
    hero_title_3: '无限可能。',
    hero_desc: 'Bull 将 KV、SQL、图引擎、全文搜索和时序存储打包进单个静态 Go 二进制。专为 AI Agent 技能扩展打造——下载即用，随处运行。',
    hero_cta: '快速开始',
    hero_cta2: '查看 GitHub',

    feat_title: '为什么选择 Bull？',
    feat_desc: 'AI Agent 本地数据处理所需的一切，浓缩在一个便携二进制中。',
    feat_1_title: '单个二进制',
    feat_1_desc: '约 8 MB 的可执行文件。无运行时、无容器、无包管理器。复制即用。',
    feat_2_title: '五大引擎',
    feat_2_desc: 'KV 存储、SQL 数据库、图引擎、全文搜索、时序存储——全部内嵌。',
    feat_3_title: '纯 Go 实现',
    feat_3_desc: '零 CGo 依赖，完全静态编译。秒级交叉编译到任意系统架构。',
    feat_4_title: 'AI 原生',
    feat_4_desc: '内置 YAML 技能定义文件，AI Agent 读取后即知该用哪个引擎、调用什么命令。',
    feat_5_title: 'JSON 输出',
    feat_5_desc: '所有查询命令均支持 --format json，为机器解析提供结构化输出。',
    feat_6_title: '持久存储',
    feat_6_desc: '所有数据重启后保留。每个引擎写入 --data-dir 下各自的目录。',

    engines_title: '五大引擎，一套工具',
    engines_desc: '每个引擎都经过精心挑选，体积最小化、能力最大化。',
    engine_kv_name: 'KV 存储',
    engine_kv_lib: 'bbolt · ~1 MB',
    engine_kv_desc: 'B+tree 键值存储，支持桶管理、批量操作、原子计数器、范围扫描和 JSON 导入导出。',
    engine_sql_name: 'SQL 数据库',
    engine_sql_lib: 'SQLite（纯 Go）· ~7 MB',
    engine_sql_desc: '完整 SQLite——CSV/JSON/NDJSON 导入、多格式导出、交互式 Shell、PRAGMA 自省。',
    engine_graph_name: '图引擎',
    engine_graph_lib: 'dominikbraun/graph · ~1 MB',
    engine_graph_desc: '有向/无向加权图——最短路径、DFS/BFS、拓扑排序、环检测、连通分量。',
    engine_search_name: '全文搜索',
    engine_search_lib: 'SQLite FTS5',
    engine_search_desc: '索引 JSON 文档，支持评分查询、字段返回、分页。NDJSON 批量导入。',
    engine_ts_name: '时序存储',
    engine_ts_lib: 'tstorage · ~0.5 MB',
    engine_ts_desc: '写入带标签的时间戳指标，范围查询、最新值查询、计数、CSV 导出。',

    cmd_title: '72 条命令，触手可及',
    cmd_desc: '每个引擎都暴露一致的 CLI 接口。在脚本和 Agent 工作流中自由组合。',

    skills_title: 'AI Agent 就绪',
    skills_desc: 'Bull 内置机器可读的 YAML 技能定义文件，精确告诉 AI Agent 该做什么。',
    skills_card_1_title: '21 个决策场景',
    skills_card_1_desc: '决策矩阵将真实任务映射到正确的引擎和命令序列。',
    skills_card_2_title: '45 个示例 Prompt',
    skills_card_2_desc: '预置提示词演示如何用自然语言触发每项能力。',
    skills_card_3_title: '结构化输出',
    skills_card_3_desc: '到处使用 --format json。Agent 直接解析结果——无需正则提取。',

    get_started_title: '30 秒快速上手',

    platforms_title: '随处运行',
    platforms_desc: '纯 Go 编译，静态二进制，无 CGo，秒级交叉编译到任意目标。',

    footer_desc: '面向 AI Agent 的全能嵌入式引擎工具箱。',
    footer_license: 'MIT 开源协议',
  },
} as const

export type DictKey = keyof typeof dict.en

export function t(lang: Lang, key: DictKey): string {
  return dict[lang][key] ?? key
}
