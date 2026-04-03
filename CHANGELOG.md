# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.3.0] - 2026-04-03

### Added

- **`web/`** — Web UI source moved into the tap monorepo (`tap/web/`); full TanStack Start + Cloudflare Workers + D1 app
- **`POST /api/batch`** — new API endpoint that accepts a full script payload from tap-scripts, validates authentication via `X-Tap-Secret`, and atomically replaces all D1 scripts in a single batch
- **`tap-scripts` repository** — standalone repo with 106 site scripts, `validate.ts`, `deploy.ts`, and a GitHub Action that auto-syncs to D1 on every push to main
- **Local script overrides** — drop a `.js` file at `~/.config/tap/sites/{site}/{script}.js` to shadow the cached version; a warning is printed when the local copy is used
- **`--local-only` flag** — skips the cache and auto-sync; only scripts in `~/.config/tap/sites/` are visible
- `WithLocalOverrideDir` functional option on `tap.Client` for library users
- `script.NewRegistryWithOverride` for constructing a registry with an override layer
- `tap.Client.ListScriptsLocalOnly()` returns only locally-overridden scripts

### Changed

- `script.Registry` now loads the override directory after the main cache, so local scripts always win
- `AGENTS.md` updated with `web/` in architecture diagram, scripts section, and documentation table
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

- Replace sites submodule with remote sync from tap-sites
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

[Unreleased]: https://github.com/vaayne/tap/compare/v0.1.6...HEAD
[0.1.6]: https://github.com/vaayne/tap/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/vaayne/tap/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/vaayne/tap/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/vaayne/tap/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/vaayne/tap/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/vaayne/tap/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/vaayne/tap/releases/tag/v0.1.0
