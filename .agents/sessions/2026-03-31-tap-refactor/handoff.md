# Handoff

<!-- Append a new phase section after each phase completes. -->

## Phase 1: Project scaffolding & script package

**Status:** complete

**Tasks completed:**

- 1.1: Init module as `github.com/vaayne/tap`, created `cmd/tap/main.go` stub, added urfave/cli and go-defuddle deps
- 1.2: Extracted script parser to `script/parser.go` with exported types (Meta, ArgDef, Script, Parse)
- 1.3: Extracted registry to `script/registry.go` with NewRegistry, Get, List
- 1.4: Added unit tests for parser (5 tests) and registry (3 tests), all passing

**Files changed:**

- `go.mod` — renamed module to `github.com/vaayne/tap`, added deps
- `cmd/tap/main.go` — CLI stub
- `script/parser.go` — script @meta parser
- `script/registry.go` — script directory scanner + index
- `script/parser_test.go` — parser unit tests
- `script/registry_test.go` — registry unit tests

**Commits:**

- `699acd3` — ✨ feat: init module github.com/vaayne/tap with cmd/tap stub
- `287f826` — ♻️ refactor: extract script parser to script/parser.go
- `2a036dc` — ♻️ refactor: extract script registry to script/registry.go
- `a6af7ba` — ✅ test: add unit tests for script parser and registry

**Decisions & context for next phase:**

- Old root-level `meta.go` and `registry.go` still exist — will be removed in Phase 6
- `script.Script` type is the canonical type going forward
- Old `main.go` still references old types — will be replaced in Phase 5

## Phase 2: Engine package

**Status:** complete

**Tasks completed:**

- 2.1: Defined `Engine` interface with `Run()`, `Name()`, `Close()` + `RunScript()` orchestrator that tries engines in order
- 2.2: Moved QuickJS runner to `engine/quickjs.go` — Go-backed fetch() polyfill, async eval
- 2.3: Moved CDP browser to `engine/browser.go` — `BrowserConfig` struct for WSURL/ProfileDir, local/remote context creation
- 2.4: Added 3 unit tests with mock engines: first succeeds, fallback, all fail

**Files changed:**

- `engine/engine.go` — Engine interface + RunScript orchestrator
- `engine/quickjs.go` — QuickJS engine implementation
- `engine/browser.go` — CDP browser engine implementation
- `engine/engine_test.go` — orchestrator fallback tests

**Commits:**

- `7cd996c` — ✨ feat: add engine interface and RunScript orchestrator
- `174c108` — ♻️ refactor: move QuickJS runner to engine/quickjs.go
- `1c2619f` — ♻️ refactor: move CDP browser to engine/browser.go
- `10b5448` — ✅ test: add engine fallback orchestrator tests

**Decisions & context for next phase:**

- `BrowserConfig` holds WSURL + ProfileDir — will be wired from Client options in Phase 4
- Default profile dir is `~/.cache/tap/chrome-profile-$USER`
- Engine interface uses `context.Context` for cancellation support

## Phase 3: Fetch package with go-defuddle

**Status:** complete

**Tasks completed:**

- 3.1: Created `fetch/fetch.go` with Fetcher, Result type, Options, HTTP fetch + defuddle parsing
- 3.2: Added tests — real URL fetch (example.com) and invalid URL error handling

**Files changed:**

- `fetch/fetch.go` — Fetcher with New/Close/Fetch, Result type, HTML fetching
- `fetch/fetch_test.go` — integration test with real URL + error case

**Commits:**

- `e744992` — ✨ feat: add fetch package with go-defuddle content extraction
- `88bcd5e` — ✅ test: add fetch package tests

**Decisions & context for next phase:**

- Fetcher holds a defuddle.Parser (expensive to create, reuse across calls)
- Default Options enables Markdown output
- Client.Fetch() in Phase 4 will wrap this package

## Phase 4: Top-level client API

**Status:** complete

**Tasks completed:**

- 4.1+4.2: Created `tap.go` with Client struct (New, Close, RunScript, Fetch, ListScripts, GetScript) and `options.go` with functional options (WithSitesDir, WithWSURL, WithProfileDir)
- 4.3: Client.RunScript wires registry + engine orchestrator with arg validation
- 4.4: Client.Fetch wraps fetch package
- 4.5: Integration tests — 5 tests covering New, ListScripts, RunScript errors, Fetch

**Files changed:**

- `tap.go` — Client struct with all methods
- `options.go` — functional options
- `tap_test.go` — integration tests
- `.old/` — moved old root-level files out (was blocking package conflict)

**Commits:**

- `f4f4639` — ✨ feat: add top-level Client API with options
- `1fac912` — ✅ test: add Client integration tests

**Decisions & context for next phase:**

- Old files moved to `.old/` (will be deleted in Phase 6)
- Client wires QuickJS → Browser engine chain automatically
- CLI in Phase 5 just creates a Client with options from env/flags

## Phase 5: CLI with urfave/cli

**Status:** complete

**Tasks completed:**

- 5.1: Set up urfave/cli v3 app with global flags (sites-dir, ws-url, profile-dir) + env var sources
- 5.2: `tap site list` — lists all scripts with arg hints
- 5.3: `tap site <name> [key=value ...]` — runs scripts with QuickJS→Browser fallback
- 5.4: `tap fetch <url>` — markdown by default, `--json` for full metadata
- 5.5: Loads .env via godotenv, env vars: TAP_SITES_DIR, TAP_WS_URL, TAP_PROFILE_DIR

**Files changed:**

- `cmd/tap/main.go` — full CLI implementation
- `.gitignore` — added `tap` binary

**Commits:**

- `f270c98` — ✨ feat: add CLI with urfave/cli — site and fetch subcommands

**Decisions & context for next phase:**

- All commands smoke-tested and working
- Ready for cleanup: remove `.old/`, update README, add LICENSE

## Phase 6: Cleanup & docs

**Status:** complete

**Tasks completed:**

- 6.1: Removed old root-level Go files (.old/)
- 6.2: Updated README with library + CLI usage, architecture diagram, roadmap
- 6.3: Added .env.example, fixed .gitignore to allow it
- 6.4: Added MIT LICENSE
- 6.5: Full build + all tests passing + smoke test
- 6.6: Pushed to github.com/vaayne/tap

**Commits:**

- `e80cf67` — 🗑️ chore: remove old root-level Go files
- `3a1b049` — 📝 docs: update README for tap library + CLI
- `c316cde` — 📝 docs: add .env.example and MIT LICENSE
