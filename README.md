# 🚰 Tap

Tap into any website from your terminal.

A Go library and CLI toolkit that runs JavaScript scripts against real websites — fast via QuickJS, with full browser fallback when needed. Also extracts clean content from any URL via [go-defuddle](https://github.com/vaayne/go-defuddle).

## Install

### CLI

```bash
go install github.com/vaayne/tap/cmd/tap@latest
```

### Library

```bash
go get github.com/vaayne/tap
```

Requires Go 1.22+ and Google Chrome (or Chromium) for browser fallback.

## CLI Usage

### Site Scripts

```bash
# List all available scripts
tap site list

# Run a script (QuickJS first, browser fallback)
tap site v2ex/hot
tap site twitter/search query=claude
tap site bilibili/search keyword=编程 order=click

# Pipe to jq
tap site hackernews/top | jq '.stories[:3]'
```

### Fetch Content

```bash
# Extract clean markdown from any URL
tap fetch https://example.com/article

# Output as JSON with full metadata
tap fetch --json https://example.com/article
```

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/vaayne/tap"
    "github.com/vaayne/tap/fetch"
)

func main() {
    client, err := tap.New(
        tap.WithSitesDir("./sites"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Run a site script
    result, err := client.RunScript(context.Background(), "v2ex/hot", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)

    // Fetch clean content from a URL
    content, err := client.Fetch(context.Background(), "https://example.com", &fetch.Options{
        Markdown: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(content.Markdown)
}
```

## How It Works

Every site script runs through a two-tier execution engine:

```
┌─────────────┐     ┌─────────────┐
│   QuickJS   │──X──>│   Chrome    │
│  Go fetch() │     │   via CDP   │
└─────────────┘     └─────────────┘
  Fast, no browser    Full browser
  Works for APIs      Cookies, DOM, auth
```

- **QuickJS** — Runs scripts in an embedded JS engine ([fastschema/qjs](https://github.com/fastschema/qjs)) with a Go-backed `fetch()` polyfill. Fast, no browser overhead.
- **Chrome (CDP)** — If QuickJS fails (missing Web APIs, CORS, needs cookies), falls back to headless Chrome via Chrome DevTools Protocol. Navigates to the script's domain for proper origin and cookie access.

For `tap fetch`, content is fetched via Go HTTP and parsed with [go-defuddle](https://github.com/vaayne/go-defuddle) to extract clean HTML/Markdown.

## Configuration

All config via environment variables, `.env` file, or CLI flags:

| Variable | Flag | Description | Default |
|---|---|---|---|
| `TAP_SITES_DIR` | `--sites-dir` | Directory containing site scripts | `./sites` |
| `TAP_WS_URL` | `--ws-url` | Remote CDP WebSocket URL | _(local Chrome)_ |
| `TAP_PROFILE_DIR` | `--profile-dir` | Chrome profile for persistent cookies | `~/.cache/tap/chrome-profile-$USER` |

### Browser Modes

**Local Chrome (default)** — launches headless Chrome with a persistent profile so cookies survive across runs.

**Remote browser** — connect to a remote CDP endpoint:
```bash
export TAP_WS_URL=wss://your-remote-browser/ws
tap site v2ex/hot
```

## Writing Scripts

Scripts live in `sites/` organized by site name:

```javascript
/* @meta
{
  "name": "site/action",
  "description": "What this script does",
  "domain": "example.com",
  "args": {
    "query": {"required": true, "description": "Search query"}
  }
}
*/

async function(args) {
  const resp = await fetch('https://api.example.com?q=' + args.query);
  return await resp.json();
}
```

### Available Sites

100+ scripts across 30+ sites:

- **Search**: Google, Bing, Baidu, DuckDuckGo
- **Social**: Twitter, Weibo, Reddit, 小红书, 即刻
- **Video**: YouTube, Bilibili
- **News**: Hacker News, BBC, Reuters, 今日头条, 36氪
- **Dev**: GitHub, Stack Overflow, Dev.to, npm, PyPI
- **Finance**: 雪球, 东方财富, Yahoo Finance
- **Knowledge**: Wikipedia, 知乎, Douban, arXiv

Run `tap site list` for the full list.

## Project Structure

```
github.com/vaayne/tap/
├── tap.go              # Client API — unified entry point
├── options.go          # Functional options (WithSitesDir, WithWSURL, ...)
├── engine/
│   ├── engine.go       # Engine interface + fallback orchestrator
│   ├── quickjs.go      # QuickJS engine with Go fetch() polyfill
│   └── browser.go      # Chrome CDP engine (local/remote)
├── fetch/
│   └── fetch.go        # URL → clean content via go-defuddle
├── script/
│   ├── parser.go       # Script @meta parser
│   └── registry.go     # Script directory scanner + index
├── cmd/tap/
│   └── main.go         # CLI binary (urfave/cli)
└── sites/              # 100+ community site scripts
```

## Roadmap

- [x] Site scripts with QuickJS + browser fallback
- [x] `tap fetch <url>` — clean content extraction
- [ ] `tap screenshot <url>` — page screenshots
- [ ] `tap pdf <url>` — save as PDF
- [ ] `tap eval <js> --url <url>` — run arbitrary JS on a page
- [ ] `tap fill <script>` — form automation

## License

MIT
