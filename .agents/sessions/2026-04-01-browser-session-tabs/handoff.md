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
