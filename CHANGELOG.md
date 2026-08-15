# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- `tap run` executes host-side JavaScript workflows from a file or stdin, with
  a thin `browser.cmd()`/`browser.eval()` bridge to the inherited agent-browser
  session.

### Fixed

- Structured agent-browser errors remain available to workflow `try/catch`
  blocks when the subprocess exits nonzero.

## [1.0.1] - 2026-08-11

### Changed

- Metadata headers are injected only into fetches targeting the script's
  declared domain. Cross-origin fetches remain governed by browser CORS/CSP and
  never receive Tap-configured credentials.

### Fixed

- Structured agent-browser batch failures are now preserved instead of being
  reduced to an opaque `exit status 1` error.

## [1.0.0] - 2026-08-11

### Added

- The default installer now bootstraps a pinned, checksummed agent-browser
  binary and Chrome runtime. Releases also publish optional full offline
  bundles containing both Tap and agent-browser.

### Changed

- **Tap 1.0 runtime boundary** — agent-browser is now the sole browser runtime;
  Tap inherits `AGENT_BROWSER_SESSION` and no longer owns Chrome, CDP, profiles,
  tabs, or session state.
- Site names now derive from `{site}/{action}.js`; metadata no longer declares
  `name`, `runtime`, or `env`.
- `tap fetch` now runs embedded Defuddle in agent-browser and accepts no URL to
  extract the active tab.

### Removed

- QuickJS, chromedp, Lightpanda, the internal transport/browser runtimes, and
  their dependencies.
- `tap browser`, `tap attach`, and `tap status`; use agent-browser directly.

### Fixed

- `tap fetch` no longer depends on site discovery or synchronization, so it
  remains usable when the remote registry is unavailable.
- Web builds now use a locked, compatible Cloudflare toolchain and deploy from
  GitHub Actions with reproducible dependencies.

## [0.4.10] - 2026-06-12

### Fixed

- **Stale Chrome singleton locks no longer block launch** — leftover `SingletonLock` files from dead Chrome processes are detected and cleared automatically before launch; a lock held by a live Chrome fails fast with an actionable error instead of raw Chrome stderr
- **Engine failures are no longer swallowed** — when script execution falls back across engines, the final error lists every attempted engine with its own cause (e.g. the QuickJS error is shown alongside the browser error), and single-engine runs no longer claim "all engines failed"

## [0.4.9] - 2026-06-12

### Added

- **`tap browser get` / `tap browser is`** — read-only element queries: `text`, `html`, `value`, `attr`, `title`, `url`, `count`, `box` (JSON), `styles` (JSON); boolean state checks `visible`, `enabled`, `checked`; all accept CSS selectors or snapshot refs (`@eN`)
- **Extra interaction commands** — `dblclick`, `focus`, `check` (idempotent), `uncheck` (idempotent), `scrollintoview`, `upload` (file input), `drag` (mouse move→press→release); low-level `mouse move/down/up/wheel` and `keyboard type/insert`; `keydown`/`keyup` for modifier keys
- **Enhanced `wait` modes** — plain duration (`2000` ms or `1.5s`), `--text` substring, `--url` glob, `--load load|domcontentloaded|networkidle`, `--fn` JS poll, `--state visible|hidden|attached|detached` for element waits; `tap browser open --wait-selector` now correctly blocks until the element appears
- **Semantic locator (`find`)** — locate elements by `role`, `text`, `label`, `placeholder`, `alt`, `title`, `testid`, `first`, `last`, `nth` and perform `click/fill/type/hover/focus/check/uncheck/text` actions without needing a CSS selector; `--name` filter for role, `--exact` for text
- **Web storage (`storage local|session`)** — read/write/clear `localStorage` and `sessionStorage` of the current tab
- **Auth state save/load (`state save|load`)** — export and import cookies + current-origin localStorage in Playwright `storageState` format (`0600` perms); non-matching origins are skipped with a warning on load
- **Emulation overrides (`set`)** — persist `viewport`, `device` preset, `geo`, `offline`, `headers`, `media` (color-scheme), and `useragent` per tab; settings are re-applied automatically on every subsequent invocation; `set clear` removes all overrides
- **Categorized `tap browser --help` output** — browser subcommands are grouped by task for faster command discovery
- **Automatic `tap-web` skill refresh after `tap upgrade`** — installed skill files are refreshed when the CLI upgrades
- **Discoverable `tap docs` command** — documentation generation is now visible from CLI help

### Changed

- Rewrote `tap-web` SKILL.md as a lean decision guide, moving command syntax to `--help` and reference docs

## [0.4.8] - 2026-05-18

### Added

- **FxEmbed Twitter scripts** — added full FxTwitter API coverage for posts, threads, conversations, profiles, search, trends, and typeahead
- **Twitter article Markdown output** — `twitter/fxembed-status` can return long-form X articles as Markdown with `format=markdown -f raw`

### Changed

- Prefer FxEmbed scripts first for X/Twitter content in the embedded `tap-web` skill
- Raw output now prints string results directly instead of JSON-escaping them

## [0.4.7] - 2026-05-12

### Added

- **WeChat article script** (`weixin/article`) — fetches WeChat MP articles via HTTP using WeChat Mobile User-Agent, no browser required; extracts title, author, content, images, and metadata

### Fixed

- **QuickJS fetch race condition** — the goroutine resolving a fetch promise was calling back into the wazero WASM module while `QJS_Eval` still held it, causing memory corruption and 100% CPU spin on large responses (e.g. 3 MB pages); fetch now runs synchronously within the WASM host callback

## [0.4.6] - 2026-05-12

### Added

- **Runtime-based engine routing** — site scripts can declare `runtime` (http/browser/lightpanda) to control which engine executes them
- **Environment variables and headers for site scripts** — `EnvDef`/`Env` fields in script Meta with validation and header resolution; resolved headers auto-injected into QuickJS fetch and Browser CDP requests
- **Twitter site scripts** — `twitter/getxapi-tweet-detail`, `twitter/getxapi-article`, `twitter/post-tweet`
- **Embedded sites sync to D1** — batch sync now accepts multiple script directories with override-by-name priority
- Full metadata display in `tap site info` and `tap site list` (runtime, env, headers)

### Changed

- Narrowed embed pattern to `*/*.js` to exclude markdown files from binary
- Added field-level godoc to `Meta` struct
- Added `sites/CLAUDE.md` documenting script structure and conventions

### Fixed

- Removed redundant authorization header from `fetch()` in scripts (meta headers are injected automatically)
- `forceBrowser` now returns only browser engines

## [0.4.5] - 2026-05-11

### Fixed

- `tap browser attach` status no longer incorrectly shows local sessions as Chrome attachment candidates

## [0.4.4] - 2026-04-13

### Added

- Shell completion support for bash, zsh, fish, and PowerShell via `tap completion`
- Command coverage and regression tests for shell completion generation

### Fixed

- Check out the tap repo before the bb-sites sync workflow runs
- Recreate the D1 staging table during batch sync so repeated syncs stay consistent
- Tolerate a missing D1 capabilities migration during batch sync

## [0.4.3] - 2026-04-13

### Changed

- Updated `github.com/vaayne/go-defuddle` to `v0.1.2`
- Restored the embedded `tap-web` skill version-check note so installed skills can be verified against `tap --version`

## [0.4.2] - 2026-04-12

### Added

- Browser snapshot refs (`@eN`) for stable element targeting across `tap browser` actions
- Ref-based `tap browser` commands for click, type, fill submit, and select workflows on dynamic pages

### Changed

- Shortened the embedded `tap-web` skill description for cleaner quick-reference guidance

## [0.4.1] - 2026-04-09

### Added

- Embedded tap-web skill in binary — `tap skill install` extracts bundled skill to local config

## [0.4.0] - 2026-04-09

### Added

- **Electron app debugging via CDP** — `tap electron` commands to connect to, inspect, and control Electron apps
- **Built-in Chrome proxy** — `tap chrome` command manages a dedicated Chrome process with configurable port and headless mode
- **Attach to user's existing Chrome** — `tap attach` connects to an already-running Chrome instance via CDP WebSocket
- **Simplified CLI commands** — streamlined `tap open`/`tap do` shorthand for common browser operations
- **Persisted default browser context** — automatic reuse of the "default" session without explicit session management

### Changed

- Pinned Go version to 1.26 for reproducible builds
- Refactored CLI surface around attach and browser commands for cleaner UX
- Updated documentation to cover Chrome proxy and browser attachment workflows

### Fixed

- Accept valid CDP WebSocket endpoints during attach (relaxed URL validation)

## [0.3.3] - 2026-04-07

### Added

- Capabilities field and daily bb-sites sync for site scripts

### Changed

- Renamed wrangler worker from `tap-web` to `tap`
- Replaced tap-scripts references with epiral/bb-sites in documentation

## [0.3.2] - 2026-04-06

### Added

- uTLS fingerprinting to bypass bot detection — HTTP transport now mimics real browser TLS fingerprints
- `tap doctor` command — check, install, and update browser dependencies (Chrome detection, Lightpanda download/update)
- Lightpanda download metadata tracking (`.meta.json`) to support update workflows
- Chrome detection utility with path and version reporting

### Changed

- Browser session resolution simplified to always use "default" session
- Documentation updated with Lightpanda platform support (macOS/Linux only), site compatibility limitations, and browser backend comparison

### Fixed

- Added h2 transport to fix uTLS ALPN negotiation on CI

## [0.3.1] - 2026-04-04

### Added

- Install script for easy one-line installation
- `tap upgrade` command for self-upgrading the CLI

### Changed

- README rewritten with use-case-driven structure
- Agent skill section added to README
- CLAUDE.md now requires lint and test pass before every commit

### Fixed

- Resolved errcheck lint errors in upgrade command

## [0.3.0] - 2026-04-04

### Added

- **`web/`** — Web UI source moved into the tap monorepo (`tap/web/`); full TanStack Start + Cloudflare Workers + D1 app
- **`POST /api/batch`** — new API endpoint that accepts a full script payload from bb-sites, validates authentication via `X-Tap-Secret`, and atomically replaces all D1 scripts in a single batch
- **`bb-sites` repository** — standalone repo with 106 site scripts, `validate.ts`, `deploy.ts`, and a GitHub Action that auto-syncs to D1 on every push to main
- **Local script overrides** — drop a `.js` file at `~/.config/tap/sites/{site}/{script}.js` to shadow the cached version; a warning is printed when the local copy is used
- **`--local-only` flag** — skips the cache and auto-sync; only scripts in `~/.config/tap/sites/` are visible
- `WithLocalOverrideDir` functional option on `tap.Client` for library users
- `script.NewRegistryWithOverride` for constructing a registry with an override layer
- `tap.Client.ListScriptsLocalOnly()` returns only locally-overridden scripts
- `tap browser text` — token-efficient plain-text page reading
- `tap browser pdf` — save page as PDF
- `tap browser keypress`, `dialog`, and `cookies` commands
- 9 human-like browser interaction commands (click, type, hover, scroll, select, screenshot, etc.)

### Changed

- `script.Registry` now loads the override directory after the main cache, so local scripts always win
- `AGENTS.md` updated with `web/` in architecture diagram, scripts section, and documentation table
- CLI help text improved for self-contained discoverability
- `tap-web` skill condensed to concise quick-reference card with decision-making guide
## [0.2.0] - 2026-04-02

### Added

- CDP network interception for browser tabs — capture, inspect, and intercept network requests on tracked tabs
- `tap browser network wait` — block until a matching request completes, print the entry
- `tap browser network body` — fetch response body by request ID
- `tap browser network log` — stream completed requests as NDJSON
- `tap browser network intercept` — block, mock, or modify matching requests via Fetch domain
- `tap browser network clear` — remove all interception rules
- URL pattern matching with glob syntax (`*` matches any character including `/`)
- `withTargetListen` CDP helper for long-lived event listening sessions
- `docs/network.md` full reference documentation
- `skills/tap-web/references/network.md` agent skill reference

### Changed

- Restructured tap-web skill with progressive disclosure (SKILL.md → 92 lines, details in `references/`)
- Added documentation sync rule to AGENTS.md

## [0.1.8] - 2026-04-02

### Added

- Full Windows support for browser session management (`tap browser` commands)
- Windows process management using `OpenProcess` / `GetExitCodeProcess` for liveness checks and `taskkill` for two-phase graceful+forced termination
- Windows Chrome discovery via `Program Files`, `LocalAppData`, and registry (`App Paths\chrome.exe`)
- Real file locking on Windows via `LockFileEx` / `UnlockFileEx` with retry semantics

### Changed

- Split `browser/process.go` into portable shared code and platform-specific files (`process_unix.go`, `process_windows.go`)
- Made `browser/manager.go` fully portable (removed `//go:build !windows` constraint)
- Profile directory cleanup now retries on failure to handle transient file locks

## [0.1.7] - 2026-04-02

### Added

- Persistent browser sessions and tabs — manage long-lived Chrome sessions with named sessions and multiple tabs (`tap browser session`, `tap browser tab`)
- Browser forms and fill commands — discover form fields and fill them programmatically (`tap browser action form`, `tap browser action fill`)
- Browser session state management with file-based store and process lifecycle tracking
- Browser lock mechanism for safe concurrent access to sessions
- `docs/browser.md` documentation for the new browser subsystem

## [0.1.6] - 2026-04-02

### Added

- Lightpanda headless browser as an alternative CDP backend (`--lightpanda` / `--lp` flag), with automatic binary download from GitHub nightly releases
- Browser wait modes for headed debugging
- `--lightpanda` support in `tap login` command

### Changed

- `transport.New` and `tap.New` now accept `context.Context` for cancelable browser startup
- Transport is closed on fetcher initialization failure to avoid leaking child processes
- Browser package uses standard `log` instead of `slog` for consistent `--quiet`/`--verbose` output

## [0.1.5] - 2026-04-01

### Added

- `tap login <url>` command — open a visible browser for login, CAPTCHA, or manual interaction; cookies persist in Chrome profile
- `--pause` flag — pause after browser navigation to allow user interaction before script execution (implies `--no-headless` and `-b`)
- `transport.PauseFunc` type and `BrowseInteractive`, `BrowseHTMLWithPause`, `BrowseEvalWithPause` transport methods
- `WithPause` functional option for the Go library

## [0.1.4] - 2026-03-31

### Fixed

- Browser engine: preserve native `fetch` before page scripts override it (fixes GitHub and similar sites)
- Browser engine: disable CORS enforcement for cross-origin API calls in scripts

## [0.1.3] - 2026-03-31

### Added

- `tap-web` agent skill for CLI web access
- `AGENTS.md` (symlinked as `CLAUDE.md`) for coding agent context

### Changed

- Replace sites submodule with remote sync from bb-sites
- Use `XDG_CACHE_HOME` for scripts cache instead of `XDG_CONFIG_HOME`

### Fixed

- Check `resp.Body.Close()` error returns to satisfy errcheck linter

## [0.1.2] - 2026-03-31

### Fixed

- Recover from QuickJS WASM panics in the engine fallback chain so runtime panics become errors and browser fallback can continue

## [0.1.1] - 2026-03-31

### Added

- mise.toml for dependency management and task running
- CLI grouped list view, search command, browser mode, and better error messages

### Changed

- Split `cmd/tap/main.go` into focused modules

### Fixed

- Resolve all golangci-lint errcheck and unused issues

## [0.1.0] - 2026-03-31

### Added

- Initial release
- Go library and CLI for running JavaScript scripts against websites
- QuickJS engine with browser (Chrome CDP) fallback
- Site script registry with submodule support
- `tap site` command — list, run, and search site scripts
- `tap fetch` command — extract clean markdown/JSON from any URL
- Shared transport layer for HTTP and browser access
- Content extraction via [go-defuddle](https://github.com/vaayne/go-defuddle)
- CI workflows for lint, test, and release
- GoReleaser configuration for cross-platform builds

[Unreleased]: https://github.com/vaayne/tap/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/vaayne/tap/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/vaayne/tap/compare/v0.4.10...v1.0.0
[0.4.10]: https://github.com/vaayne/tap/compare/v0.4.9...v0.4.10
[0.4.9]: https://github.com/vaayne/tap/compare/v0.4.8...v0.4.9
[0.4.8]: https://github.com/vaayne/tap/compare/v0.4.7...v0.4.8
[0.4.7]: https://github.com/vaayne/tap/compare/v0.4.6...v0.4.7
[0.4.6]: https://github.com/vaayne/tap/compare/v0.4.5...v0.4.6
[0.4.5]: https://github.com/vaayne/tap/compare/v0.4.4...v0.4.5
[0.4.4]: https://github.com/vaayne/tap/compare/v0.4.3...v0.4.4
[0.4.3]: https://github.com/vaayne/tap/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/vaayne/tap/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/vaayne/tap/compare/v0.4.0...v0.4.1
[0.3.3]: https://github.com/vaayne/tap/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/vaayne/tap/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/vaayne/tap/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/vaayne/tap/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/vaayne/tap/compare/v0.1.8...v0.2.0
[0.1.8]: https://github.com/vaayne/tap/compare/v0.1.7...v0.1.8
[0.1.7]: https://github.com/vaayne/tap/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/vaayne/tap/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/vaayne/tap/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/vaayne/tap/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/vaayne/tap/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/vaayne/tap/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/vaayne/tap/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/vaayne/tap/releases/tag/v0.1.0
