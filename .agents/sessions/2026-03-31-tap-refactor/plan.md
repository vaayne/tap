# Plan: Tap — Web Interaction Toolkit

## Overview

Refactor the existing `cdp` project into **Tap** (`github.com/vaayne/tap`), a Go library + CLI toolkit for interacting with web pages. Library-first design with a CLI built on urfave/cli.

### Goals

- Restructure as an importable Go library with clean public API
- CLI via urfave/cli with dedicated subcommands (`site`, `fetch`)
- Integrate go-defuddle for `tap fetch` (clean content extraction)
- Three-tier execution engine: Go net/http → QuickJS → CDP Browser
- Rename all references from `cdp` to `tap`

### Success Criteria

- [ ] `go get github.com/vaayne/tap` works — library importable
- [ ] `tap site list` lists all scripts
- [ ] `tap site v2ex/hot` runs a script (QuickJS → CDP fallback)
- [ ] `tap fetch <url>` returns clean markdown content via defuddle
- [ ] `tap fetch --json <url>` returns full metadata as JSON
- [ ] All existing scripts work unchanged
- [ ] Env vars: `TAP_WS_URL`, `TAP_PROFILE_DIR`

### Out of Scope

- `tap screenshot`, `tap pdf`, `tap eval`, `tap fill` (roadmap only)
- New site scripts
- Tier 1 pure Go fetch for site scripts (future optimization)

## Technical Approach

### Architecture

```
github.com/vaayne/tap/              # module root
├── tap.go                          # Client struct — unified library API
├── options.go                      # Client options (WithWSURL, WithProfileDir, etc.)
├── engine/
│   ├── engine.go                   # Engine interface + RunScript orchestrator
│   ├── quickjs.go                  # QuickJS runner + fetch polyfill
│   └── browser.go                  # CDP browser context + runner
├── fetch/
│   └── fetch.go                    # Fetch URL → clean HTML/Markdown via defuddle
├── script/
│   ├── parser.go                   # @meta + body parser
│   └── registry.go                 # scan sites/ dir, index by name
├── cmd/tap/
│   └── main.go                     # urfave/cli: site, fetch subcommands
├── sites/                          # community JS scripts (unchanged)
└── go.mod                          # module: github.com/vaayne/tap
```

### Key Design Decisions

1. **Library API** — top-level `Client` struct with `Fetch()` and `RunScript()`. Sub-packages also importable directly.
2. **Engine interface** — `engine.Engine` with `Run(script, args)` method. QuickJS and Browser implement it. Orchestrator tries in order.
3. **go-defuddle** — used inside `fetch/` package. Fetches HTML via Go net/http, then parses with defuddle.
4. **Script registry** — accepts a directory path. CLI defaults to embedded or local `sites/`. Library users pass their own.
5. **Config** — functional options pattern for Client. CLI reads env vars / flags and maps to options.

### Public API

```go
// Top-level client
client, err := tap.New(
    tap.WithSitesDir("./sites"),
    tap.WithWSURL("wss://..."),          // optional
    tap.WithProfileDir("~/.cache/tap"),   // optional
)
defer client.Close()

// Run a site script
result, err := client.RunScript(ctx, "v2ex/hot", map[string]string{})

// Fetch clean content from URL
content, err := client.Fetch(ctx, "https://example.com/article")
content.Title    // "Article Title"
content.Markdown // clean markdown
content.HTML     // clean HTML
```

## Implementation Phases

### Phase 1: Project scaffolding & script package

Set up new module structure, move script parsing + registry into `script/` package.

1. Init new module `github.com/vaayne/tap`, set up `cmd/tap/` (files: `go.mod`, `cmd/tap/main.go`)
2. Move meta parsing to `script/parser.go` (files: `script/parser.go`)
3. Move registry to `script/registry.go` (files: `script/registry.go`)
4. Unit test parser + registry (files: `script/parser_test.go`, `script/registry_test.go`)

### Phase 2: Engine package

Move QuickJS runner and browser runner into `engine/` with a common interface.

1. Define Engine interface + orchestrator in `engine/engine.go` (files: `engine/engine.go`)
2. Move QuickJS runner to `engine/quickjs.go` (files: `engine/quickjs.go`)
3. Move CDP browser to `engine/browser.go` (files: `engine/browser.go`)
4. Unit test engine fallback logic (files: `engine/engine_test.go`)

### Phase 3: Fetch package with go-defuddle

Implement `tap fetch` — HTTP fetch + defuddle content extraction.

1. Create `fetch/fetch.go` — fetch URL, parse with defuddle, return Result (files: `fetch/fetch.go`)
2. Define Result type with Title, Markdown, HTML, metadata (files: `fetch/fetch.go`)
3. Unit test with sample HTML (files: `fetch/fetch_test.go`)

### Phase 4: Top-level client API

Create the unified `Client` with functional options.

1. Create `tap.go` — Client struct, New(), Close() (files: `tap.go`)
2. Create `options.go` — WithSitesDir, WithWSURL, WithProfileDir (files: `options.go`)
3. Implement Client.RunScript() — wires script registry + engine (files: `tap.go`)
4. Implement Client.Fetch() — wires fetch package (files: `tap.go`)
5. Integration test (files: `tap_test.go`)

### Phase 5: CLI with urfave/cli

Build the CLI binary with subcommands.

1. Set up urfave/cli app with global flags (files: `cmd/tap/main.go`)
2. `tap site list` command (files: `cmd/tap/main.go`)
3. `tap site <name> [args]` command (files: `cmd/tap/main.go`)
4. `tap fetch <url>` command with `--json`/`--markdown` flags (files: `cmd/tap/main.go`)
5. Load .env, wire env vars to client options (files: `cmd/tap/main.go`)

### Phase 6: Cleanup & docs

Remove old files, update README, push to vaayne/tap.

1. Remove old root-level Go files (`main.go`, `meta.go`, `registry.go`, `browser.go`, `qjsrunner.go`)
2. Update README.md for new structure + library usage
3. Update .gitignore, .env.example
4. Add LICENSE (MIT)
5. Final build + smoke test all commands
6. Push to github.com/vaayne/tap

## Testing Strategy

- `script/parser_test.go` — parse @meta from sample scripts, edge cases (no meta, malformed)
- `script/registry_test.go` — scan test fixtures dir, verify indexing
- `engine/engine_test.go` — test fallback: mock QuickJS fail → browser fallback
- `fetch/fetch_test.go` — test defuddle parsing with sample HTML
- `tap_test.go` — integration test: RunScript + Fetch via Client
- Manual smoke test: `tap site v2ex/hot`, `tap site bilibili/search keyword=编程`, `tap fetch https://example.com`

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| go-defuddle QuickJS conflicts with engine QuickJS | Medium | They use separate runtimes, should be isolated |
| Large binary size (wazero + chromedp + defuddle) | Low | Acceptable for a CLI tool |
| Some site scripts may break during refactor | Medium | Don't touch sites/ contents, only move wrapper code |
| urfave/cli learning curve | Low | Well-documented, simple API |

## Open Questions

_None — all resolved into assumptions:_
- Assumption: `sites/` ships with the repo, not embedded in binary (users can add custom scripts)
- Assumption: `tap fetch` defaults to markdown output, `--json` for full metadata
- Assumption: No Tier 1 pure Go fetch for site scripts in this phase

## Review Feedback

_(Updated during plan review rounds)_

## Final Status

_(Updated after implementation completes)_
