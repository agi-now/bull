export type Lang = 'en' | 'zh'

const dict = {
  en: {
    hero_title: 'Bull',
    hero_sub: 'Micro data environment in a single ~8 MB Go binary. Five embedded engines — KV, SQL, Graph, Search, Time-Series — ready the moment you download it.',
    hero_github: 'GitHub',
    hero_get_started: 'Get Started',

    engines_title: 'Engines',
    engines_desc: 'Each engine is embedded — no external servers, no dependencies.',

    commands_title: 'Commands',
    commands_desc: 'All engines share a consistent CLI interface. Every query supports --format json.',

    skills_title: 'Agent Skills',
    skills_desc: 'Bull ships with SKILL.md following the Agent Skills specification. Agents discover capabilities at startup and activate on demand.',
    skills_p1: 'Metadata (~100 tokens) loaded at startup for task matching.',
    skills_p2: 'Full instructions loaded only when activated.',
    skills_p3: 'Per-engine references in references/ loaded on demand.',

    api_title: 'HTTP API',
    api_desc: 'Run bull serve to expose all engines as a JSON REST API.',

    install_title: 'Install',

    platforms_title: 'Platforms',

    footer: 'Bull — MIT License',
  },
  zh: {
    hero_title: 'Bull',
    hero_sub: '单个 ~8 MB Go 二进制的微型数据环境。五大内嵌引擎 — KV、SQL、图、搜索、时序 — 下载即用。',
    hero_github: 'GitHub',
    hero_get_started: '快速开始',

    engines_title: '引擎',
    engines_desc: '每个引擎都内嵌运行 — 无外部服务器，无依赖。',

    commands_title: '命令',
    commands_desc: '所有引擎共享一致的 CLI 接口。每个查询支持 --format json。',

    skills_title: 'Agent Skills',
    skills_desc: 'Bull 内置 SKILL.md，遵循 Agent Skills 规范。Agent 启动时发现能力，按需激活。',
    skills_p1: '元数据（~100 token）启动时加载，用于任务匹配。',
    skills_p2: '完整指令仅在激活时加载。',
    skills_p3: '各引擎参考文档在 references/ 中，按需读取。',

    api_title: 'HTTP API',
    api_desc: '运行 bull serve 将所有引擎暴露为 JSON REST API。',

    install_title: '安装',

    platforms_title: '平台',

    footer: 'Bull — MIT 开源协议',
  },
} as const

export type DictKey = keyof typeof dict.en

export function t(lang: Lang, key: DictKey): string {
  return dict[lang][key] ?? key
}
