# 🚰 Tap

Reusable site programs and web extraction powered by
[agent-browser](https://github.com/vercel-labs/agent-browser).

Tap owns website knowledge: script discovery, metadata, arguments, environment
expansion, headers, and readable extraction. agent-browser owns Chrome,
sessions, profiles, tabs, CDP, interaction, and network tooling.

## Install

Install the matching native agent-browser binary from
[GitHub Releases](https://github.com/vercel-labs/agent-browser/releases/latest),
then install Chrome:

```bash
install -m 0755 agent-browser-<platform>-<arch> ~/.local/bin/agent-browser
agent-browser install
```

Release assets include `darwin-arm64`, `darwin-x64`, `linux-arm64`,
`linux-x64`, musl Linux variants, and `win32-x64.exe`. Homebrew and npm remain
fallbacks:

```bash
brew install agent-browser       # macOS
# or: npm install -g agent-browser
```

Then install Tap:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
# or
go install github.com/vaayne/tap/cmd/tap@latest
```

Verify both with `tap doctor`.

## Quick start

### Discover, inspect, execute

```bash
tap site list
tap site search github
tap site info exa/search
tap site exa/search query="agent-browser" count=5
```

Scripts auto-sync from [tap.vaayne.com](https://tap.vaayne.com) into
`~/.cache/tap/sites/`. Local overrides at
`~/.config/tap/sites/{site}/{script}.js` take precedence.

### Extract readable content

```bash
tap fetch https://example.com/article
tap fetch --json https://example.com/article
```

With no URL, Tap extracts the current agent-browser tab without navigating:

```bash
agent-browser open https://example.com/article
tap fetch
```

### Continue with arbitrary interaction

Tap does not wrap browser automation commands. Use agent-browser directly:

```bash
agent-browser snapshot -i
agent-browser click @e3
agent-browser network requests --filter api
```

## Sessions

Tap never creates, names, persists, or closes browser sessions. It inherits
agent-browser's environment unchanged:

```bash
export AGENT_BROWSER_SESSION=my-task

agent-browser open https://github.com
tap site github/notifications
tap fetch
agent-browser snapshot -i
```

Without `AGENT_BROWSER_SESSION`, agent-browser's default session is used.

## Site format

The script name comes from its path. For example,
`sites/exa/search.js` is `exa/search`:

```javascript
/* @meta
{
  "description": "Search the web with Exa",
  "domain": "mcp.exa.ai",
  "args": {
    "query": {"required": true},
    "count": {}
  },
  "headers": {
    "X-API-Key": "${EXA_API_KEY}"
  }
}
*/

async function(args) {
  // fetch(...) runs in agent-browser; metadata headers are merged into requests.
}
```

Environment variables are inferred from `${VAR}` references. A header is
omitted when one of its referenced variables is unset. Resolved headers are
applied before navigation and cleared after script execution.

## Command map

```text
tap
├── site       discover, inspect, sync, and execute site programs
├── fetch      extract a URL or the current agent-browser tab
└── doctor     check the agent-browser runtime dependency
```

`upgrade`, `skill`, `docs`, and `completion` remain maintenance commands.

## Agent skill

Tap ships the `tap-web` skill:

```bash
npx skills add vaayne/tap
# or
tap skill install
```

Its escalation order is `tap site` → `tap fetch` → `agent-browser read` →
agent-browser snapshot/interaction.

## Go library

```bash
go get github.com/vaayne/tap
```

```go
client, err := tap.New(ctx, tap.WithSitesDir("./sites"))
if err != nil {
    log.Fatal(err)
}

result, err := client.RunScript(ctx, "exa/search", map[string]string{
    "query": "agent-browser",
})
content, err := client.Fetch(ctx, "", &fetch.Options{Markdown: true})
```

The library uses the same inherited agent-browser session and does not close it.

## Docs

- [docs/cli.md](docs/cli.md) — generated CLI reference
- [skills/tap-web/references/script-development.md](skills/tap-web/references/script-development.md) — site script development
- [web/README.md](web/README.md) — registry web app

## License

MIT
