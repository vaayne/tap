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

**Before every commit**, run `mise run lint && mise run test` and fix any failures. Do not commit code that fails lint or tests.

## Architecture

```
tap.go / options.go   → Client API + functional options
transport/            → Shared HTTP + CDP browser layer
browser/              → Persistent sessions, tabs, network interception
engine/               → QuickJS + browser fallback
fetch/                → URL → clean content (go-defuddle)
cmd/tap/              → CLI (site, fetch, sync, browser)
web/                  → Web UI — TanStack Start + Cloudflare Workers + D1
```

## Scripts

Site scripts live in [tap-scripts](https://github.com/vaayne/tap-scripts) (separate repo).
They auto-sync to D1 on push via GitHub Actions → `POST /api/batch`.
The CLI caches them locally in `~/.cache/tap/sites/` (refreshed every 24 h).

Local overrides: drop a `.js` file at `~/.config/tap/sites/{site}/{script}.js`.
It takes precedence over the cache automatically. Use `--local-only` to skip the cache entirely.

## Testing

Run `go test ./... -timeout 60s -race` before pushing. CI: lint + test on ubuntu/macos.

## Release

1. Ensure CI passes (`mise run lint && mise run test`).
2. Update `CHANGELOG.md` ([Keep a Changelog](https://keepachangelog.com/) format).
3. Commit: `📝 docs: update CHANGELOG for vx.y.z`.
4. Tag and push: `git tag vx.y.z && git push origin main --tags`.
5. GoReleaser publishes automatically.

## Documentation

User-facing docs live in three places. Keep them in sync when making changes:

| Location | Purpose | Update when |
|---|---|---|
| `README.md` | Project overview, quick start, links to docs | New features, commands, or doc files |
| `docs/` | Full reference docs (`browser.md`, `network.md`) | Command changes, new flags, behavior changes |
| `web/` | Web UI source — TanStack Start + Cloudflare Workers | API route changes, D1 schema changes |
| `skills/tap-web/` | Agent skill (`SKILL.md` + `references/`) | Command changes, new capabilities |

## Skills

Agent skills live in `skills/`. See `skills/tap-web/` for web access.

Skill structure uses progressive disclosure:
- `SKILL.md` — lean quick-reference (loaded on trigger)
- `references/` — detailed docs loaded on demand by the agent
