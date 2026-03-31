# AGENTS.md

## Project

Tap is a Go CLI and library that runs JavaScript scripts against real websites and extracts clean content from any URL. It uses QuickJS for fast execution with headless Chrome (CDP) fallback.

## Stack

- **Language**: Go 1.26+
- **CLI framework**: urfave/cli v3
- **JS engine**: QuickJS (via fastschema/qjs)
- **Browser**: Chrome DevTools Protocol (chromedp)
- **Content extraction**: go-defuddle
- **Task runner**: mise

## Commands

```bash
mise run build        # Build binary to bin/tap
mise run test         # Run tests
mise run lint         # Run golangci-lint
mise run fmt          # Format code
mise run tidy         # Tidy go modules
```

## Code Style

- Follow standard Go conventions (`gofmt`, `golangci-lint`).
- Always check error returns, including `resp.Body.Close()`.
- Use functional options pattern for configuration (see `options.go`).
- Keep packages focused: `transport/`, `engine/`, `fetch/`, `cmd/tap/`.

## Commit Convention

Use emoji-prefixed Conventional Commits:

- `✨ feat:` — new feature
- `🐛 fix:` — bug fix
- `♻️ refactor:` — code restructuring
- `📝 docs:` — documentation
- `🔥 chore:` — maintenance

## Architecture

```
tap.go              → Client API, unified entry point
options.go          → Functional options (WithSitesDir, WithWSURL, ...)
transport/          → Shared HTTP + CDP browser layer
engine/             → QuickJS engine + browser fallback orchestrator
fetch/              → URL → clean content via go-defuddle
cmd/tap/            → CLI binary (site, fetch subcommands + sync)
```

## Testing

- Run `go test ./... -timeout 60s -race` before pushing.
- CI runs lint + tests on ubuntu and macos.

## Skills

Agent skills live in `skills/`. See `skills/tap-web/` for the web access skill.
