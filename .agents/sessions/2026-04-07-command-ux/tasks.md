# Tasks: CLI Command UX Redesign

## Phase 1: Command model and help scaffolding

### Tasks

- [x] Add `tap status` command
  - Show default browser context status
  - Show current tab info
  - Support `--json` flag
  - Files: `cmd/tap/status.go`

- [x] Add `tap attach` command family
  - `tap attach chrome` - auto-discover and attach to user Chrome
  - `tap attach chrome --browser-url <url>` - explicit URL
  - `tap attach chrome --port-file <path>` - explicit port file
  - `tap attach electron --port <port>` - attach to running Electron
  - `tap attach electron --launch <binary>` - launch and attach
  - `tap attach status` - show attachment info
  - `tap attach clear` - detach from external target
  - Support `--json` for status
  - Files: `cmd/tap/attach.go`, `cmd/tap/attach_chrome.go`, `cmd/tap/attach_electron.go`

- [x] Add simplified `tap browser` commands
  - `tap browser open <url>` - open/navigate current tab
  - `tap browser tabs` - list tabs (alias for `tab list`)
  - `tap browser switch <tab-id>` - switch tab (alias for `tab select`)
  - `tap browser close-tab [tab-id]` - close tab (alias for `tab close`)
  - `tap browser status` - show browser context and tab status (alias for `session info`)
  - Support `--json` for tabs and status
  - Files: `cmd/tap/browser_simple.go`

- [x] Update help text and descriptions
  - Rewrite top-level help to emphasize task-first workflows
  - Rewrite `tap browser --help` to show common workflows first
  - Keep advanced commands but demote in help presentation
  - Files: `cmd/tap/main.go`, `cmd/tap/browser.go`

- [x] Add flag aliases
  - `--show` alias for `--no-headless`
  - `--wait` alias for `--delay`
  - `--browser-url` alias for `--ws-url`
  - `--devtools-url` alias for `--ws-url`
  - Files: `cmd/tap/main.go`, command-specific flag definitions

## Phase 2: Default-context behavior (pending)

- [x] Implement active-context state model
- [x] Add default-context resolution logic
- [x] Ensure no silent fallback on stale contexts
- [x] Add `--json` to all status commands

## Phase 3: Attach workflow integration (pending)

- [x] Wire `attach chrome` to discovery/proxy internals
- [x] Wire `attach electron` to existing attach/launch
- [x] Auto-adopt windows/tabs where safe

## Phase 4: Docs and migration (pending)

- [x] Rewrite README with new workflows
- [x] Rewrite docs/browser.md
- [x] Add migration guide
- [x] Mark advanced commands

## Status

- [x] Plan created and approved
- [x] Phase 1 complete
- [x] Phase 2 complete
- [x] Phase 3 complete
- [x] Phase 4 complete
