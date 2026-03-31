---
name: tap-web
description: >
  Access websites, search the web, and extract clean content using the `tap` CLI.
  Use when the user asks to search the web, read a webpage, fetch article content,
  get trending topics, look up social media posts, check stock prices, search videos,
  or retrieve structured data from any supported site. Triggers on: "search for",
  "what's trending on", "fetch this page", "read this URL", "get content from",
  "look up on Twitter/Weibo/Reddit/YouTube/etc.", "check stock", "translate",
  or any web access task.
---

# tap-web

Use the `tap` CLI to access websites from the terminal. Two main commands:

## `tap fetch` — Extract clean content from any URL

```bash
tap fetch <url>              # Clean markdown output
tap fetch --json <url>       # JSON with full metadata
```

Use for reading articles, documentation, blog posts, or any webpage.

## `tap site` — Run site-specific scripts for structured data

```bash
tap site <site/action> [key=value ...]
```

### Discover scripts

```bash
tap site list                # List all available scripts
tap site info <script>       # Show script details and parameters
tap site search <keyword>    # Search scripts by name/description
```

### Common examples

```bash
# Search engines
tap site google/search query="climate change"
tap site bing/search query="rust programming"
tap site duckduckgo/search query="best laptops 2026"
tap site baidu/search query="人工智能"

# Social media
tap site twitter/search query="AI agents" count=10
tap site twitter/tweets screen_name=elonmusk count=5
tap site weibo/hot count=10
tap site xiaohongshu/search keyword="咖啡推荐"
tap site reddit/search query="coding agents"

# News & trending
tap site hackernews/top count=10
tap site toutiao/hot count=10
tap site bbc/news count=5
tap site reuters/search query="technology"

# Video
tap site youtube/search query="SwiftUI tutorial" max=5
tap site bilibili/search keyword="编程" count=5

# Developer
tap site github/repo repo=vaayne/tap
tap site github/issues repo=vaayne/tap
tap site stackoverflow/search query="golang concurrency"
tap site npm/search query="react hooks"

# Finance
tap site xueqiu/stock symbol=AAPL
tap site yahoo-finance/quote symbol=TSLA
tap site eastmoney/stock query="贵州茅台"

# Knowledge
tap site zhihu/search keyword="量子计算" count=5
tap site wikipedia/search query="machine learning"

# Translation
tap site youdao/translate query="hello world"
```

### Parameter conventions

- Parameters marked with `*` are required (shown in `tap site list` output).
- Use `key=value` syntax for all parameters.
- Use `count` or `max` to limit results.

### Output format

```bash
tap site <script> -f json     # Raw JSON (default for piping)
tap site <script> -f pretty   # Human-readable (default for terminal)
tap site <script> -f raw      # Unformatted raw output
```

Pipe to `jq` for filtering:

```bash
tap site hackernews/top count=5 -f json | jq '.'
```

## Global options

| Flag | Description |
|---|---|
| `--browser, -b` | Force browser mode (skip QuickJS) — use for sites needing login/cookies |
| `--no-headless` | Show browser window (debug auth issues) |
| `--timeout, -t` | Execution timeout (e.g., `30s`, `2m`) |
| `--quiet, -q` | Suppress log output |
| `--verbose` | Enable verbose logging |

## Tips

- For sites requiring authentication (Twitter, Xiaohongshu, Bilibili, etc.), run once with `--no-headless` to log in manually. Cookies persist in the Chrome profile.
- Use `tap fetch` for arbitrary URLs; use `tap site` for structured data from known sites.
- When a site script exists for the target, prefer `tap site` over `tap fetch` for better structured output.
- If QuickJS execution fails, tap automatically falls back to browser mode.
