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

Scripts are automatically downloaded from [tap.vaayne.com](https://tap.vaayne.com) and cached in `$XDG_CACHE_HOME/tap/sites/` (default `~/.cache/tap/sites/`). The cache auto-refreshes every 24 hours.

```bash
# List all available scripts
tap site list

# Run a script (QuickJS first, browser fallback)
tap site v2ex/hot
tap site twitter/search query=claude
tap site bilibili/search keyword=编程 order=click

# Pipe to jq
tap site hackernews/top | jq '.stories[:3]'

# Search scripts online
tap site search bilibili

# Manually sync/update scripts
tap site sync
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

Both `tap site` and `tap fetch` share a common transport layer with two-tier network access:

```
                    ┌─────────────────────────────────┐
                    │       Shared Transport Layer     │
                    │  Level 1: HTTP  │  Level 2: CDP  │
                    └────────┬────────┴───────┬────────┘
                             │                │
              ┌──────────────┴──┐    ┌────────┴──────────┐
              │    tap site     │    │    tap fetch       │
              │  QuickJS → CDP  │    │  HTTP → CDP        │
              │  → structured   │    │  → defuddle        │
              │    JSON         │    │  → markdown/HTML   │
              └─────────────────┘    └───────────────────┘
```

**Transport layer** — Shared HTTP client and headless Chrome (CDP) browser, configured once and used by all consumers.

**Site scripts** — Predefined recipes that know the optimal path to fetch structured data. Tries [QuickJS](https://github.com/fastschema/qjs) (fast, Go-backed `fetch()`) first, falls back to Chrome via CDP for pages needing cookies, DOM, or auth.

**Fetch** — Generic content extraction from any URL. Tries direct HTTP first, falls back to browser for JS-rendered pages. Parses with [go-defuddle](https://github.com/vaayne/go-defuddle) to extract clean HTML/Markdown.

## Configuration

All config via environment variables, `.env` file, or CLI flags:

| Variable | Flag | Description | Default |
|---|---|---|---|
| `TAP_SITES_DIR` | `--sites-dir` | Directory containing site scripts | `~/.config/tap/sites` |
| `TAP_WS_URL` | `--ws-url` | Remote CDP WebSocket URL | _(local Chrome)_ |
| `TAP_BROWSER` | `--browser`, `-b` | Force browser execution, skip QuickJS | `false` |
| `TAP_PROFILE_DIR` | `--profile-dir` | Chrome profile for persistent cookies | `~/.cache/tap/chrome-profile-$USER` |
| | `--pause` | Pause after navigation for manual interaction (TTY only) | `false` |
| | `--delay` | Wait a fixed duration after navigation | `0s` |
| | `--wait-selector` | Wait until a CSS selector becomes visible | `""` |
| | `--wait-js` | Wait until a JavaScript expression becomes truthy | `""` |
| | `--no-headless` | Run browser in visible mode | `false` |

### Browser Modes

**Local Chrome (default)** — launches headless Chrome with a persistent profile so cookies survive across runs.

**Remote browser** — connect to a remote CDP endpoint:
```bash
export TAP_WS_URL=wss://your-remote-browser/ws
tap site v2ex/hot
```

### Login & Interactive Browser

Some sites require login or CAPTCHA solving. Tap provides two ways to interact with the browser before running scripts:

**`tap login`** — open a browser to log in, cookies are saved for future runs:
```bash
# Log in once, run scripts many times
tap login https://github.com/login
tap site -b github/notifications

# Use a specific profile for work accounts
tap login --profile-dir ~/.tap/work https://internal.corp.com
```

**Browser wait modes** — interact with or wait on the page before extraction/script execution:
```bash
# Browser opens → you solve CAPTCHA → press Enter → script executes
tap site --pause twitter/search query=claude

# Wait 5 seconds after navigation
tap fetch https://example.com --delay 5s

# Continue once a selector becomes visible
tap fetch https://example.com --wait-selector '.redeem-code'

# Continue once a JS expression becomes truthy
tap fetch https://example.com --wait-js 'document.body.innerText.includes("Code")'
```

`--pause`, `--delay`, `--wait-selector`, and `--wait-js` each imply `--no-headless` (visible browser) and `-b` (browser mode). `--pause` requires an interactive terminal. Cookies are always saved to the Chrome profile directory.

### Persistent Browser Sessions

Tap's persistent browser workflow lives under `tap browser ...`. Phase 1 defines the durable state model, command surface, and lifecycle contract; runtime commands land in later phases.

Planned command tree:

```bash
tap browser session new <name> [--profile-dir <dir>] [--no-headless]
tap browser session new <name> --ws-url <remote-cdp-url>
tap browser session list
tap browser session info [name]
tap browser session select <name>
tap browser session close [name]

tap browser tab new <name> [--session <session>] [--url about:blank]
tap browser tab list [--session <session>]
tap browser tab select <name> [--session <session>]
tap browser tab close [name] [--session <session>]

tap browser navigate <url> [--session <session>] [--tab <tab>]
tap browser evaluate <javascript> [--session <session>] [--tab <tab>] [--format pretty|json|raw]
tap browser screenshot [--session <session>] [--tab <tab>] [--output <path>]
```

Resolution rules:

- Session resolution order: `--session`, then the selected session from `tap browser session select`, then the only available session when exactly one exists.
- Tab resolution order: `--tab`, then the selected tab within the resolved session, then the only live tracked tab when exactly one exists.
- If a session or tab reference is still ambiguous, tap fails with guidance instead of guessing.

Lifecycle rules:

- A session is one persistent browser instance, either local or remote.
- A tab is a named tracked CDP target within a session.
- Only tracked tabs are part of tap metadata. Untracked live browser tabs are ignored by default.
- If a tracked target disappears, tap marks the tab `stale`, clears invalid selected-tab state, and requires the user to recreate or reselect a usable tab.
- `tap browser tab close` removes the tracked tab and promotes the next remaining live tracked tab by creation order, or leaves no selected tab if none remain.
- `tap browser session close` removes the resolved session. For local sessions, tap first verifies browser-process ownership before terminating the managed browser and deleting its managed profile directory. For remote sessions, tap removes metadata only and never kills the remote browser.

Local-vs-remote contract:

| Operation | Local managed browser | Remote CDP session |
|---|---|---|
| `session new` | Launches a managed Chrome with a dedicated profile and debugging endpoint | Persists the explicit `--ws-url` and validates connection/auth/TLS at creation time |
| `session close` | Verifies ownership, stops the managed browser, then removes metadata and managed profile state | Removes tap metadata only |
| `tab new` / `tab close` | Supported against the managed browser | Supported when the remote endpoint allows target management; failures must be returned clearly |
| `navigate` / `evaluate` / `screenshot` | Operate on the resolved tracked tab | Operate through the persisted session endpoint, ignoring later global `--ws-url` overrides |

The browser metadata store uses a durable state directory instead of cache storage. Override it with `TAP_BROWSER_STATE_DIR` or `tap browser --state-root <dir> ...`.

## Writing Scripts

Scripts live in the [tap-sites](https://github.com/vaayne/tap-sites) catalog, organized by site name:

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
├── transport/
│   └── transport.go    # Shared network layer (HTTP + CDP browser)
├── engine/
│   ├── engine.go       # Engine interface + fallback orchestrator
│   ├── quickjs.go      # QuickJS engine with Go fetch() polyfill
│   └── browser.go      # Chrome CDP engine (delegates to transport)
├── fetch/
│   └── fetch.go        # URL → clean content via go-defuddle (HTTP → browser fallback)
├── script/
│   ├── parser.go       # Script @meta parser
│   └── registry.go     # Script directory scanner + index
└── cmd/tap/
    ├── main.go         # CLI binary (urfave/cli)
    └── sync.go         # Remote script sync + search
```

## Roadmap

- [x] Site scripts with QuickJS + browser fallback
- [x] `tap fetch <url>` — clean content extraction
- [x] `tap login <url>` — interactive browser login with cookie persistence
- [x] `--pause` flag — manual interaction before script execution
- [ ] `tap screenshot <url>` — page screenshots
- [ ] `tap pdf <url>` — save as PDF
- [ ] `tap eval <js> --url <url>` — run arbitrary JS on a page
- [ ] `tap fill <script>` — form automation

## License

MIT
