# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

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

[0.1.3]: https://github.com/vaayne/tap/compare/v0.1.2...v0.1.3
[0.1.1]: https://github.com/vaayne/tap/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/vaayne/tap/releases/tag/v0.1.0
