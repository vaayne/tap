# 🚰 Tap

Reusable site programs, browser workflows, and web extraction powered by
[agent-browser](https://github.com/vercel-labs/agent-browser).

Tap owns host-side workflow execution and website knowledge: script discovery,
metadata, arguments, environment expansion, headers, and readable extraction.
agent-browser owns Chrome, sessions, profiles, tabs, CDP, interaction, and
network tooling.

## Install

The default installer bootstraps Tap, a pinned native agent-browser binary, and
its Chrome runtime:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
```

Every downloaded binary is pinned and SHA-256 verified. Existing
`agent-browser` installations are preserved. Installer options:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --full
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --skip-chrome
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --without-agent-browser
```

Full archives are published for macOS, glibc/musl Linux, and Windows x64. They
contain both CLIs and the agent-browser license for offline transfer; use an
existing Chrome installation or run `agent-browser install` when online.
Tap resolves the runtime in this order: `TAP_AGENT_BROWSER`, a sibling binary
from a full bundle, then `agent-browser` on `PATH`.

Manual Tap-only installation remains available, but requires an independently
installed agent-browser:

```bash
go install github.com/vaayne/tap/cmd/tap@latest
```

Verify both with `tap doctor`.

See [Installation](docs/install.md) for the platform matrix, checksum
verification, offline-transfer caveat, runtime resolution, and upgrade paths.

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

### Run a browser workflow

`tap run` executes host-side JavaScript from a file or stdin. Browser commands
still run through agent-browser in the inherited session:

```bash
tap run <<'JS'
await browser.open("https://example.com")
const page = await browser.snapshot("-i")
console.log(page.snapshot)
JS
```

Use `browser.cmd(...args)` for any supported browser command,
`browser.eval(js)` for page-side JavaScript, or the `open` and `snapshot`
shortcuts. Command results are the agent-browser JSON `data` payload.

### Continue with arbitrary interaction

For one-off interaction, use agent-browser directly:

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
  // fetch(...) runs in agent-browser. Metadata headers are merged only into
  // requests targeting the declared domain; browser CORS/CSP still applies.
}
```

Environment variables are inferred from `${VAR}` references. A header is
omitted when one of its referenced variables is unset. Resolved headers are
injected into same-domain script fetches and are never installed as
browser-wide navigation headers.

## Command map

```text
tap
├── site       discover, inspect, sync, and execute site programs
├── fetch      extract a URL or the current agent-browser tab
├── run        execute a host-side JavaScript browser workflow
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

Its escalation order is `tap site` → `tap fetch` → `tap run` for programmable
workflows → direct agent-browser interaction.

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

- [docs/install.md](docs/install.md) — installation, full bundles, and upgrades
- [docs/cli.md](docs/cli.md) — generated CLI reference
- [skills/tap-web/references/script-development.md](skills/tap-web/references/script-development.md) — site script development
- [web/README.md](web/README.md) — registry web app

## License

MIT
