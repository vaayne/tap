# cdp

A CLI tool that runs JavaScript scripts to fetch data from websites. Scripts execute with a two-tier strategy: first via QuickJS (fast, no browser), then falling back to a real browser via Chrome DevTools Protocol when needed.

## Install

```bash
go install
```

Requires Go 1.22+ and Google Chrome (or Chromium) for browser fallback.

## Quick Start

```bash
# List all available scripts
cdp list

# Run a script
cdp v2ex/hot

# Run a script with arguments
cdp twitter/search query=claude
cdp bilibili/search keyword=编程 order=click
cdp boss/search query=golang city=上海
```

## How It Works

Every script goes through a two-tier execution strategy:

```
┌──────────────┐     ┌───────────────┐     ┌───────────────┐
│  Parse script │──>│  Try QuickJS   │──X──>│ Fall back to  │
│  (@meta + fn) │    │  (Go fetch())  │     │ CDP browser   │
└──────────────┘     └───────────────┘     └───────────────┘
                           │                      │
                        Success                Success
                           │                      │
                           v                      v
                        JSON output           JSON output
```

1. **QuickJS** — Runs the script in an embedded JS engine ([fastschema/qjs](https://github.com/fastschema/qjs)) with a Go-backed `fetch()` polyfill via `net/http`. Fast, no browser overhead. Works for scripts that call simple APIs.
2. **CDP Browser** — If QuickJS fails (missing Web APIs like `URLSearchParams`, CORS, needs cookies/auth), falls back to a real Chrome browser. Navigates to the script's domain first for proper origin and cookie access.

```
# QuickJS succeeds (fast path)
Running: v2ex/hot — 获取 V2EX 最热主题
{"count": 10, "topics": [...]}

# QuickJS fails, falls back to browser
Running: bilibili/search — Search Bilibili videos by keyword
QuickJS failed: ReferenceError: URLSearchParams is not defined, falling back to browser
{"count": 20, "videos": [...]}
```

## Browser Modes

### Local Chrome (default)

When no `CDP_WS_URL` is set, `cdp` launches a local headless Chrome. A persistent profile directory is used by default so cookies and storage survive across runs:

```
~/.cache/cdp/chrome-profile-<username>
```

Override with:

```bash
export CDP_PROFILE_DIR=/path/to/profile
```

### Remote Browser

Connect to a remote CDP endpoint (e.g., [Lightpanda](https://lightpanda.io), [Browserless](https://browserless.io)):

```bash
export CDP_WS_URL=wss://your-remote-browser/ws
cdp v2ex/hot
```

## Configuration

All config is via environment variables or `.env` file:

| Variable | Description | Default |
|---|---|---|
| `CDP_WS_URL` | Remote CDP WebSocket URL. If set, local Chrome is not launched. | _(unset — uses local Chrome)_ |
| `CDP_PROFILE_DIR` | Chrome user data directory for persistent cookies/storage. | `~/.cache/cdp/chrome-profile-$USER` |

## Scripts

Scripts live in `sites/` organized by site name. Each script has a `@meta` header and an async function body:

```javascript
/* @meta
{
  "name": "twitter/search",
  "description": "Search tweets",
  "domain": "x.com",
  "args": {
    "query": {"required": true, "description": "Search query"},
    "count": {"required": false, "description": "Number of results (default 20)"}
  },
  "readOnly": true
}
*/

async function(args) {
  const resp = await fetch('https://api.example.com/search?q=' + args.query);
  return await resp.json();
}
```

### Available Sites

Over 100 scripts across 30+ sites including:

- **Search**: Google, Bing, Baidu, DuckDuckGo
- **Social**: Twitter, Weibo, Reddit, 小红书, 即刻
- **Video**: YouTube, Bilibili
- **News**: Hacker News, BBC, Reuters, 今日头条, 36氪
- **Dev**: GitHub, Stack Overflow, Dev.to, npm, PyPI
- **Finance**: 雪球, 东方财富, Yahoo Finance
- **Knowledge**: Wikipedia, 知乎, Douban, arXiv

Run `cdp list` for the full list.

## Project Structure

```
├── main.go        # CLI entry point, arg parsing, orchestration
├── browser.go     # Browser context (local Chrome / remote CDP)
├── qjsrunner.go   # QuickJS runner with Go fetch() polyfill
├── meta.go        # Script parser (@meta + function body)
├── registry.go    # Script discovery and indexing
├── sites/         # JavaScript scripts organized by site
│   ├── twitter/
│   ├── github/
│   ├── bilibili/
│   └── ...
└── .env           # Environment variables (gitignored)
```

## License

MIT
