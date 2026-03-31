# cdp

A CLI tool that runs JavaScript scripts in a real browser via Chrome DevTools Protocol (CDP). Scripts execute with full browser capabilities — cookies, `fetch`, DOM — making it easy to scrape or interact with websites that require authentication.

## Install

```bash
go install
```

Requires Go 1.22+ and Google Chrome (or Chromium) installed locally.

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

## How It Works

1. Reads the script file from `sites/`
2. Parses the `@meta` JSON header and the async function body
3. Connects to a browser (local Chrome or remote CDP endpoint)
4. Evaluates the function in the browser context with `await`
5. Returns the result as JSON

Since scripts run inside a real browser, they have access to cookies, sessions, and all Web APIs — no need to manage auth tokens manually.

## Project Structure

```
├── main.go        # CLI entry point
├── browser.go     # Browser context (local/remote)
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
