# 🚰 Tap

Tap into any website from your terminal.

Tap is a Go library, CLI, and web app for running JavaScript scripts against real websites. It uses QuickJS first for speed, falls back to Chrome CDP when needed, and can extract clean content from any URL via [go-defuddle](https://github.com/vaayne/go-defuddle).

## Repository

This monorepo contains:

- `cmd/tap/` — Go CLI
- `web/` — web app for browsing scripts ([web/README.md](web/README.md))
- `script/` — script registry and parser
- `engine/` — QuickJS and browser engines
- `transport/` — shared HTTP and CDP transport
- `tap.go` — library entry point

## Install

### CLI

```bash
go install github.com/vaayne/tap/cmd/tap@latest
```

### Library

```bash
go get github.com/vaayne/tap
```

Requires Go 1.22+ and Google Chrome or Chromium for browser fallback.

## Web

Browse scripts at [tap.vaayne.com](https://tap.vaayne.com).

## CLI usage

### Site scripts

Scripts are downloaded from [tap.vaayne.com](https://tap.vaayne.com) and cached in `$XDG_CACHE_HOME/tap/sites/` (default: `~/.cache/tap/sites/`). The cache refreshes every 24 hours.

```bash
# List scripts
tap site list

# Run a script
tap site v2ex/hot
tap site twitter/search query=claude
tap site bilibili/search keyword=编程 order=click

# Pipe to jq
tap site hackernews/top | jq '.stories[:3]'

# Search scripts online
tap site search bilibili

# Manually sync the cache
tap site sync

# Use a local override at ~/.config/tap/sites/{site}/{script}.js
tap site github/repo vaayne/tap

# Skip the remote cache entirely
tap --local-only site list
tap --local-only site github/repo vaayne/tap
```

### Fetch content

```bash
# Extract clean markdown
tap fetch https://example.com/article

# Output JSON with metadata
tap fetch --json https://example.com/article
```

### Browser automation

Tap also supports persistent browser sessions:

```bash
tap browser session new work
tap browser tab new main --url https://example.com
tap browser navigate https://httpbin.org/html
tap browser evaluate 'document.title'
tap browser screenshot
tap browser session close work
```

See [docs/browser.md](docs/browser.md) for details.

## Library usage

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

    result, err := client.RunScript(context.Background(), "v2ex/hot", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)

    content, err := client.Fetch(context.Background(), "https://example.com", &fetch.Options{
        Markdown: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(content.Markdown)
}
```

## How it works

Both `tap site` and `tap fetch` share the same transport layer:

```text
                    ┌─────────────────────────────────┐
                    │       Shared Transport Layer     │
                    │  Level 1: HTTP  │  Level 2: CDP  │
                    └────────┬────────┴───────┬────────┘
                             │                │
              ┌──────────────┴──┐    ┌────────┴──────────┐
              │    tap site     │    │    tap fetch      │
              │  QuickJS → CDP  │    │  HTTP → CDP       │
              │  → structured   │    │  → defuddle       │
              │    JSON         │    │  → markdown/HTML  │
              └─────────────────┘    └───────────────────┘
```

- **Transport** — shared HTTP client and headless Chrome via CDP
- **Site scripts** — QuickJS first, browser fallback for cookies, DOM, or auth
- **Fetch** — direct HTTP first, browser fallback for JS-rendered pages

## Configuration

Config can be set with environment variables, `.env`, or CLI flags:

| Variable | Flag | Description | Default |
|---|---|---|---|
| `TAP_SITES_DIR` | `--sites-dir` | Directory containing site scripts | `~/.config/tap/sites` |
| `TAP_WS_URL` | `--ws-url` | Remote CDP WebSocket URL | _(local Chrome)_ |
| `TAP_BROWSER` | `--browser`, `-b` | Force browser execution | `false` |
| `TAP_PROFILE_DIR` | `--profile-dir` | Chrome profile for persistent cookies | `~/.cache/tap/chrome-profile-$USER` |
| | `--pause` | Pause after navigation for manual interaction | `false` |
| | `--delay` | Wait a fixed duration after navigation | `0s` |
| | `--wait-selector` | Wait until a CSS selector becomes visible | `""` |
| | `--wait-js` | Wait until a JavaScript expression becomes truthy | `""` |
| | `--no-headless` | Run browser in visible mode | `false` |

## Browser modes

**Local Chrome** is the default and uses a persistent profile so cookies survive across runs.

**Remote browser** uses a CDP WebSocket endpoint:

```bash
export TAP_WS_URL=wss://your-remote-browser/ws
tap site v2ex/hot
```

## Login and interactive browser

Use `tap login` to log in once and reuse saved cookies later:

```bash
tap login https://github.com/login
tap site -b github/notifications
```

Interactive wait modes are also supported:

```bash
tap site --pause twitter/search query=claude
tap fetch https://example.com --delay 5s
tap fetch https://example.com --wait-selector '.redeem-code'
tap fetch https://example.com --wait-js 'document.body.innerText.includes("Code")'
```

`--pause`, `--delay`, `--wait-selector`, and `--wait-js` imply visible browser mode and browser execution.

## Writing scripts

Scripts live in the separate [tap-scripts](https://github.com/vaayne/tap-scripts) repository.

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
  const resp = await fetch('https://api.example.com?q=' + args.query)
  return await resp.json()
}
```

## Available sites

100+ scripts across 30+ sites, including:

- Search: Google, Bing, Baidu, DuckDuckGo
- Social: Twitter, Weibo, Reddit, 小红书, 即刻
- Video: YouTube, Bilibili
- News: Hacker News, BBC, Reuters, 今日头条, 36氪
- Dev: GitHub, Stack Overflow, Dev.to, npm, PyPI
- Finance: 雪球, 东方财富, Yahoo Finance
- Knowledge: Wikipedia, 知乎, Douban, arXiv

Run `tap site list` for the full list.

## Docs

- [docs/browser.md](docs/browser.md) — persistent browser sessions
- [docs/network.md](docs/network.md) — network interception
- [web/README.md](web/README.md) — web app docs

## Roadmap

- [x] Site scripts with QuickJS + browser fallback
- [x] `tap fetch <url>`
- [x] `tap login <url>`
- [x] `--pause`
- [x] `tap browser`
- [x] `tap browser forms` / `tap browser fill`
- [ ] `tap pdf <url>`

## License

MIT
