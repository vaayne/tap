# Handoff

<!-- Append a new phase section after each phase completes. -->

## Phase 1: Design session/tab model and metadata store

**Status:** complete

**Tasks completed:**

- 1.1: Defined the `tap browser` command tree, flag defaults, selection rules, and lifecycle semantics in CLI help text plus README planning docs.
- 1.2: Added a focused `browser/` package with durable session and tab record types, validation rules, capability metadata, and stale-target reconciliation behavior.
- 1.3: Implemented a disk-backed metadata store with atomic JSON writes, configurable durable state-root resolution, and file-lock based store/session serialization.
- 1.4: Added metadata-focused tests for CRUD, session/tab resolution, stale-target reconciliation, environment-based state-root selection, and session lock exclusivity.
- 1.5: Ran `go test ./... -timeout 60s -race` after the Phase 1 commit series.

**Files changed:**

- `browser/doc.go` — documented the persistent browser model and the Phase 1 local/remote contract.
- `browser/state.go` — added session/tab types, validation, default capability matrix, selection logic, and reconciliation helpers.
- `browser/store.go` — added durable state loading/saving, atomic writes, state-root discovery, and session-scoped update helpers.
- `browser/lock.go` — added advisory file-lock helpers for store and per-session serialization.
- `browser/state_test.go` — covered validation, session/tab selection, and stale-target reconciliation.
- `browser/store_test.go` — covered persisted metadata updates and state-root behavior.
- `browser/lock_test.go` — covered session-lock exclusivity.
- `cmd/tap/browser.go` — added the top-level `tap browser` namespace and shared browser command help text.
- `cmd/tap/browser_session.go` — defined the Phase 1 session-management command surface and lifecycle help text.
- `cmd/tap/browser_tab.go` — defined the Phase 1 tab-management command surface and selection behavior.
- `cmd/tap/browser_action.go` — defined the Phase 1 navigate/evaluate/screenshot command surface.
- `cmd/tap/main.go` — registered the new `browser` namespace.
- `README.md` — documented the persistent browser command tree, lifecycle rules, and local-vs-remote capability contract.

**Commits:**

- `3b3c317` — `✨ feat: add browser session state model`
- `0dd35f2` — `✨ feat: add browser state store and locks`
- `c5d471f` — `✅ test: cover browser metadata store behavior`
- `5775b1a` — `📝 docs: define browser CLI surface`

**Decisions & context for next phase:**

- The persistent browser workflow is intentionally namespaced under `tap browser` to keep it separate from the existing one-shot `site`, `fetch`, and `login` flows.
- Session resolution falls back to a selected session first, then a single available session; tab resolution falls back to a selected tab first, then a single live tracked tab.
- The store uses a durable state directory (`TAP_BROWSER_STATE_DIR` or `--state-root`) instead of cache storage, with advisory file locks on both global state and per-session operations.
- Remote session metadata freezes the creation-time `--ws-url` so Phase 2 reconnect logic can ignore later global endpoint overrides.
- Selected tabs are never allowed to remain stale or closed in durable state; reconciliation clears them so Phase 2 runtime commands can fail predictably instead of guessing.

## Phase 2.1: Target-aware CDP helpers in transport

**Status:** complete

**Tasks completed:**

- Added `transport/cdp.go` with stateless, package-level CDP helper functions that connect to an already-running browser via its debug WebSocket URL and operate on specific targets.
- Implemented `TargetInfo` struct and seven public functions: `ListTargets`, `CreateTarget`, `CloseTarget`, `NavigateTarget`, `EvalTarget`, `ScreenshotTarget`.
- Implemented two internal helpers (`withBrowser`, `withTarget`) that manage remote allocator and context lifecycle for browser-level vs target-level operations.
- Each function is self-contained: connects, performs its work, and disconnects.

**Files changed:**

- `transport/cdp.go` — new file with target-aware CDP helpers.

**Commits:**

- `706147c` — `✨ feat: add target-aware CDP helpers in transport`

**Decisions & context for next phase:**

- Functions are stateless and accept `debugURL` directly rather than being methods on `Transport`, since the session manager (Phase 2.3) will supply the debug URL from session metadata.
- `ListTargets` filters to `type == "page"` only; other target types (service workers, iframes) are excluded.
- `CloseTarget` in the current cdproto version returns only an error (no boolean success flag), so we rely on the error alone.
- `EvalTarget` uses `WithReturnByValue(true).WithAwaitPromise(true)` to match the existing `BrowseEval` semantics in `transport.go`.
- `ScreenshotTarget` uses quality 90 for full-page PNG capture.

## Phase 2.2: Local browser process management

**Status:** complete

**Tasks completed:**

- 2.2: Added browser/process.go with Chrome binary discovery, process launch with auto-port and ownership tokens, liveness checking, and safe termination.

**Files changed:**

- `browser/process.go` — Chrome discovery, LaunchBrowser, CheckProcess, KillProcess, debug URL parsing

**Commits:**

- `c353635` — `✨ feat: add browser process launch and lifecycle helpers`

**Decisions & context for next phase:**

- `findChrome` uses `exec.LookPath` first (cross-platform) then falls back to well-known absolute paths per OS (`runtime.GOOS`). No build tags needed — runtime detection is sufficient for darwin/linux.
- `LaunchBrowser` uses `--remote-debugging-port=0` so the OS auto-assigns a free port; the actual debug URL is parsed from Chrome's stderr (`DevTools listening on ...`) with a 10-second timeout.
- `SysProcAttr.Setpgid = true` ensures Chrome runs in its own process group and survives tap process exit.
- `cmd.Wait()` is intentionally never called — Chrome is long-lived and `cmd.Process.Pid` is the only handle stored.
- `KillProcess` sends SIGTERM first, polls with signal-0 every 100ms, and escalates to SIGKILL after 5 seconds.
- Ownership token is 16 random bytes hex-encoded (32 chars). Phase 2.3 will use it to verify that the tap instance that launched a browser is the one managing it.
- `debugURLToHTTP` handles both `ws://` and `wss://` schemes for forward compatibility with TLS-enabled debug endpoints.

## Phase 2.3: Session manager

**Status:** complete

**Tasks completed:**

- 2.3: Added browser/manager.go with session lifecycle, tab lifecycle, browser actions, and reconciliation.

**Files changed:**

- `browser/manager.go` — Manager type + all session/tab/action/reconciliation methods

**Commits:**

- `aac731a` — `✨ feat: add browser session manager`

**Decisions & context for next phase:**

- `Manager` is a thin coordination layer: it delegates process management to `browser/process.go`, CDP operations to `transport/cdp.go`, and metadata persistence to `browser/store.go`.
- `CreateSession` for local mode launches the browser first, then persists metadata; on metadata save failure, it best-effort kills the just-launched process to avoid orphans.
- `CreateSession` for remote mode validates endpoint reachability via HTTP GET `/json/version` before persisting, and stores the WSURL in both `Remote.WSURL` and `Process.DebugURL` so `resolveDebugURL` works uniformly.
- `CloseSession` uses `WithSessionLock` + manual Load/Save rather than `UpdateSession` because it needs to perform process kill and profile cleanup between load and save.
- Read-only operations (`ListSessions`, `GetSession`, `ListTabs`) use `store.Load()` without locking — they tolerate momentarily stale reads.
- Mutating tab/action operations use `store.UpdateSession` which combines session-scoped locking with atomic store mutation.
- `Evaluate` and `Screenshot` capture their return values via closure variables and use `UpdateSession` for locking, even though they don't strictly mutate state — this ensures consistent reads under concurrent access.
- `Navigate` updates `tab.URL` and `tab.UpdatedAt` in the same locked transaction that performs the CDP navigation.
- Session and tab name parameters that are empty strings flow through to `ResolveSession("")` / `ResolveTab("")` which implement the fallback-to-selected-or-only logic.

### Review fixes applied (Phase 2.1–2.3)

- `NavigateTarget` error wrapping added (reviewer Phase 2.1).
- `LaunchBrowser` switched from `exec.CommandContext` to `exec.Command` to prevent context cancellation from killing detached Chrome; `cmd.Wait()` added on error path, `cmd.Process.Release()` on success; ESRCH handled in `KillProcess` SIGTERM path; `CheckProcess` drains response body (reviewer Phase 2.2).
- `Evaluate`/`Screenshot` decoupled from state locks via `resolveTarget` helper; `Navigate` uses three-phase resolve→CDP→persist; `CloseSession` uses two-phase read→kill→delete with best-effort `KillProcess` (reviewer Phase 2.3).

## Phase 2.4: Runtime and reconciliation tests

**Status:** complete

**Tasks completed:**

- 2.4: Added unit tests for transport/cdp.go (TargetInfo struct), browser/process.go (debugURLToHTTP, parseDebugURL, findChrome, KillProcess edge cases, CheckProcess nil), and browser/manager.go (CreateSession validation, remote endpoint rejection, ListSessions, SelectSession, GetSession resolution, resolveDebugURL, requireLiveTab).

**Files changed:**

- `transport/cdp_test.go` — TargetInfo field validation
- `browser/process_test.go` — debugURLToHTTP table tests, parseDebugURL parsing/timeout, findChrome smoke, KillProcess/CheckProcess edge cases
- `browser/manager_test.go` — manager lifecycle validation, session resolution, helper coverage

**Commits:**

- `86c5db5` — `✅ test: add unit tests for CDP helpers, process lifecycle, and manager`

**Decisions & context for next phase:**

- Integration tests requiring a real Chrome instance are deferred to Phase 4 validation.
- Tests cover metadata-level and helper-level logic that can run without a browser.

## Phase 3: Expose library API and CLI commands

**Status:** complete

**Tasks completed:**

- Wired all browser CLI commands to the session manager
- session new/list/info/select/close
- tab new/list/select/close
- navigate, evaluate, screenshot

**Files changed:**

- `cmd/tap/browser.go` — added newBrowserManager helper, removed Phase 1 placeholders
- `cmd/tap/browser_session.go` — wired session commands
- `cmd/tap/browser_tab.go` — wired tab commands
- `cmd/tap/browser_action.go` — wired navigate/evaluate/screenshot

**Commits:**

- `51634c1` — `✨ feat: wire browser CLI commands to session manager`

**Decisions & context for next phase:**

- Every CLI action starts with `configureLogging(cmd)` and creates a fresh `Manager` via `newBrowserManager(cmd)` — no shared state between invocations.
- Status messages go to stderr; data output (evaluate results, tab/session tables) goes to stdout, matching existing tap conventions.
- Session/tab name arguments pass empty string to the Manager when omitted, relying on the Manager's `ResolveSession`/`ResolveTab` fallback logic (selected → only-one).
- `session list` determines the selected session by calling `GetSession(ctx, "")` which resolves to the selected/only session; errors are silently ignored (no selection marker shown if resolution fails).
- `screenshot` generates a deterministic filename from session name, tab name, and Unix timestamp when `--output` is omitted.
- `session new` added `--ws-url` and `--no-headless` flags; headless defaults to true (local mode launches headless unless `--no-headless` is set).

### Review fixes applied (Phase 3)

- Added `resolveSessionName` helper in `browser/manager.go` to resolve empty session names before they reach `WithSessionLock`/`UpdateSession` (which reject empty names via `ValidateSessionName`). Applied to CloseSession, CreateTab, CloseTab, SelectTab, Navigate, Evaluate, Screenshot, and Reconcile.
- Fixed errcheck lint issues on `fmt.Fprintf`/`fmt.Fprintln` calls in CLI output code.

## Phase 4: Validation and hardening

**Status:** complete

**Tasks completed:**

- Ran `mise run fmt` — applied gofmt to all files (minor alignment fixes in pre-existing lightpanda.go, sync.go)
- Ran `mise run lint` — fixed 7 errcheck issues in CLI and test files, 0 issues remaining
- Committed import cycle fix after rebase: moved CDP helpers from `transport/cdp.go` to `browser/cdp.go` to resolve cycle introduced by `browser/lightpanda.go` from main

**Files changed:**

- `browser/cdp.go` — moved from `transport/cdp.go` (package change)
- `browser/cdp_test.go` — moved from `transport/cdp_test.go`
- `browser/manager.go` — removed `transport` import, added `resolveSessionName` helper
- `browser/process_test.go` — fixed errcheck on `w.Close()`
- `cmd/tap/browser_session.go` — fixed errcheck on `fmt.Fprintf`/`fmt.Fprintln`
- `cmd/tap/browser_tab.go` — fixed errcheck on `fmt.Fprintf`/`fmt.Fprintln`
- `browser/lightpanda.go`, `cmd/tap/sync.go`, `cmd/tap/browser.go` — gofmt alignment

**Commits:**

- `9cdacbe` — `♻️ refactor: move CDP helpers from transport to browser package`
- `4b52c3f` — `🐛 fix: resolve optional session names before locking`
- `d04fac3` — `🐛 fix: satisfy errcheck linter for fmt write calls`
- `c1d6aee` — `🔥 chore: apply gofmt to pre-existing files`
