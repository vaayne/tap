# 🚰 Tap

Tap into any website from your terminal.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
```

Or with a custom directory:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --dir /usr/local/bin
```

Or install with Go:

```bash
go install github.com/vaayne/tap/cmd/tap@latest
```

To upgrade to the latest version:

```bash
tap upgrade
```

Requires Google Chrome for browser features.

## Site scripts — structured data from 100+ sites

Run community scripts to get structured JSON from popular websites.

```bash
tap site hackernews/top | jq '.stories[:3]'
tap site twitter/search query=claude
tap site bilibili/search keyword=编程 order=click
tap site github/repo vaayne/tap
```

Scripts auto-sync from [tap.vaayne.com](https://tap.vaayne.com) and cache locally. Run `tap site list` for the full catalog, or `tap site search <query>` to find scripts.

Local overrides: drop a `.js` file at `~/.config/tap/sites/{site}/{script}.js`. Use `--local-only` to skip the remote cache entirely. Write your own — see the [script development guide](https://github.com/vaayne/tap-scripts).

## Fetch — clean content from any URL

Extract readable content as markdown or JSON, powered by [go-defuddle](https://github.com/vaayne/go-defuddle).

```bash
tap fetch https://example.com/article
tap fetch --json https://example.com/article
```

Falls back to a real browser for JS-rendered pages automatically.

## Browser — persistent sessions and automation

Manage long-lived Chrome sessions, automate interactions, and intercept network requests.

```bash
tap browser session new work
tap browser tab new main --url https://example.com
tap browser navigate https://httpbin.org/html
tap browser click '#submit'
tap browser text
tap browser screenshot
tap browser network wait --url '*/api/data'
tap browser session close work
```

Save a PDF, handle dialogs, fill forms, press keys, and more. See [docs/browser.md](docs/browser.md) and [docs/network.md](docs/network.md) for the full reference.

### Login and interactive mode

Log in once and reuse saved cookies:

```bash
tap login https://github.com/login
tap site -b github/notifications
```

Wait for dynamic content:

```bash
tap fetch https://example.com --wait-selector '.content'
tap fetch https://example.com --wait-js 'document.body.innerText.includes("ready")'
```

## Agent skill

Tap ships with a built-in [agent skill](https://agentskills.io) that gives AI coding agents web access, site scripts, and browser automation. Works with Claude Code, Cursor, Copilot, and [45+ other agents](https://github.com/vercel-labs/skills).

```bash
npx skills add vaayne/tap
```

Once installed, your agent will automatically use `tap` when you ask it to search the web, read a URL, run site scripts, or automate a browser.

## Agent / LLM usage

Tap is designed as a tool-use target for AI agents. Every `tap site` script returns deterministic structured JSON, and every command is a single CLI call with no interactive prompts — ideal for function calling and pipelines.

## Go library

```bash
go get github.com/vaayne/tap
```

```go
client, _ := tap.New()
defer client.Close()

// Run a site script
result, _ := client.RunScript(ctx, "hackernews/top", nil)

// Extract clean content
content, _ := client.Fetch(ctx, "https://example.com", &fetch.Options{Markdown: true})
fmt.Println(content.Markdown)
```

## Configuration

Essential flags and environment variables:

| Flag | Env | Description |
|---|---|---|
| `--browser`, `-b` | `TAP_BROWSER` | Force browser execution |
| `--sites-dir` | `TAP_SITES_DIR` | Script directory (default `~/.config/tap/sites`) |
| `--ws-url` | `TAP_WS_URL` | Remote CDP WebSocket URL |
| `--local-only` | | Skip remote cache, use local scripts only |

Browser flags (imply `--browser`):

| Flag | Description |
|---|---|
| `--pause` | Pause after navigation for manual interaction |
| `--delay` | Wait a fixed duration after navigation |
| `--wait-selector` | Wait until a CSS selector is visible |
| `--wait-js` | Wait until a JS expression is truthy |
| `--no-headless` | Run browser in visible mode |
| `--profile-dir` | Chrome profile directory |

## Links

- [tap.vaayne.com](https://tap.vaayne.com) — browse scripts online
- [tap-scripts](https://github.com/vaayne/tap-scripts) — script repository
- [docs/browser.md](docs/browser.md) — browser sessions reference
- [docs/network.md](docs/network.md) — network interception reference
- [web/README.md](web/README.md) — web app docs

## License

MIT
