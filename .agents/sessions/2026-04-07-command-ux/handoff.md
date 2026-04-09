# Handoff: CLI Command UX Redesign

## 2026-04-07

### Completed

- Finished the task-first top-level help rewrite in `cmd/tap/main.go`
- Hid browser-specific global flags from top-level help and surfaced them on relevant commands
- Added browser-context resolution for `site`, `fetch`, and `login` so they reuse:
  1. explicit `--browser-url` / `--profile-dir`
  2. persisted default context
  3. managed default profile for login flows
- Added `tap login status` and `tap login clear`
- Added `tap site run <script>` while keeping `tap site <script>` compatibility
- Improved `tap browser tabs` output to match the new common UX better
- Rewrote `README.md`
- Rewrote `docs/browser.md`
- Added `docs/command-ux-migration.md`

### Validation

- Ran `gofmt -w` on modified Go files
- Ran `go test ./...`
- Checked `tap --help`, `tap login --help`, `tap site --help`, `tap fetch --help`

### Notes

- The persisted default-context record is still lighter than the original aspirational JSON shape in the plan, but behavior now matches the intended UX contract: explicit overrides win, attached contexts are sticky, and stale attached contexts fail explicitly rather than silently falling back.
