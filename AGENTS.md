# AGENTS.md

## Project

Go CLI and library for running JS scripts against websites (QuickJS + Chrome CDP fallback) and extracting clean content from URLs via go-defuddle.

## Stack

Go 1.26+, urfave/cli v3, QuickJS (fastschema/qjs), chromedp, go-defuddle, mise.

## Commands

```bash
mise run build    # Build to bin/tap
mise run test     # Run tests
mise run lint     # golangci-lint
mise run fmt      # gofmt
mise run tidy     # go mod tidy
```

## Code Style

- Standard Go conventions, `gofmt`, `golangci-lint`.
- Always check error returns (including `resp.Body.Close()`).
- Functional options pattern (see `options.go`).
- Focused packages: `transport/`, `engine/`, `fetch/`, `cmd/tap/`.

## Commits

Emoji-prefixed Conventional Commits: `✨ feat:`, `🐛 fix:`, `♻️ refactor:`, `📝 docs:`, `🔥 chore:`.

## Architecture

```
tap.go / options.go   → Client API + functional options
transport/            → Shared HTTP + CDP browser layer
engine/               → QuickJS + browser fallback
fetch/                → URL → clean content (go-defuddle)
cmd/tap/              → CLI (site, fetch, sync)
```

## Testing

Run `go test ./... -timeout 60s -race` before pushing. CI: lint + test on ubuntu/macos.

## Release

1. Ensure CI passes (`mise run lint && mise run test`).
2. Update `CHANGELOG.md` ([Keep a Changelog](https://keepachangelog.com/) format).
3. Commit: `📝 docs: update CHANGELOG for vx.y.z`.
4. Tag and push: `git tag vx.y.z && git push origin main --tags`.
5. GoReleaser publishes automatically.

## Skills

Agent skills live in `skills/`. See `skills/tap-web/` for web access.
