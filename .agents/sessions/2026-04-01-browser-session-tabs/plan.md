# Plan: Browser Sessions and Tabs

## Overview

Add a new persistent browser automation workflow to `tap` that is separate from the existing one-shot `site`, `fetch`, and `login` flows. The new workflow must support first-class browser sessions and tabs from day one, allowing navigation, JavaScript evaluation, and screenshots against persistent browser state across separate CLI invocations.

### Goals

- Add persistent browser automation commands that survive across separate CLI runs.
- Support both sessions and tabs as explicit concepts from the first release.
- Keep existing `site`, `fetch`, and `login` behavior unchanged.
- Build on the current `chromedp` transport layer with a clear path to future browser automation features.

### Success Criteria

- [ ] Users can create, list, inspect, and close persistent browser sessions.
- [ ] Users can create, list, select, and close tabs within a session.
- [ ] Users can run `navigate`, `evaluate`, and `screenshot` against a chosen session and tab across separate CLI invocations.
- [ ] Session/tab metadata survives process restarts and is reconciled when browser targets disappear.
- [ ] Local-browser and remote-CDP configurations behave predictably and report actionable errors.
- [ ] Automated tests cover metadata management, target resolution, and CLI behavior for the new commands.

### Out of Scope

- Multi-window orchestration beyond a session’s tracked tabs.
- Form automation, PDF export, network interception, and recorder-style workflows.
- Synchronizing sessions across machines.
- Extending `site` or `fetch` to reuse the new persistent session model in this change.

## Technical Approach

Introduce a dedicated persistent browser session manager alongside the existing one-shot browser helpers. The manager will store session and tab metadata on disk, launch or reconnect to independently managed browser processes, resolve CDP targets, and expose high-level operations for session lifecycle, tab lifecycle, navigation, evaluation, and screenshots. For local persistent sessions, `tap` cannot rely on the current one-shot `chromedp.NewExecAllocator` lifecycle; instead it must own browser processes separately and use `chromedp` only to attach and control them.

For local Chrome, `tap` will launch a persistent browser process with a dedicated remote debugging port and user data directory per session, then save its connection metadata for later invocations. The runtime contract must define crash/restart behavior explicitly: if a managed browser disappears, the session record remains, live target IDs are invalidated, tracked tabs are marked stale or recreated according to the reconciliation rules, and selected-tab state is updated deterministically. For remote CDP mode, `tap` will create a session by binding a session record to an explicit endpoint (`tap browser session new <name> --ws-url <url>`), persist that endpoint in session metadata, and ignore global `--ws-url` defaults for already-created remote sessions so reconnect behavior is deterministic. Connection, auth, and TLS failures must be validated at session creation time and surfaced with actionable errors; unsupported lifecycle operations will fail clearly rather than partially working.

The CLI should introduce a dedicated `browser` namespace to reflect the larger scope and avoid overloading top-level commands as more session/tab operations are added.

### Components

- **`browser/` package**: New focused package for persisted session metadata, browser process management, tab bookkeeping, target lookup, and automation operations.
- **Metadata store**: Disk-backed JSON state under a durable app-state directory (for example XDG state/data locations rather than a cache directory) containing sessions, tabs, selected session, selected tab, browser endpoint, process info, ownership markers, and timestamps. Persistence must use atomic write-then-rename semantics plus interprocess locking so separate CLI invocations cannot corrupt state.
- **Runtime locking**: Session-mutating commands must also take a session-scoped lock around `load -> reconcile -> maybe relaunch/reconnect -> mutate targets -> persist` so concurrent CLI invocations cannot race browser relaunch, tab creation, or selection updates.
- **Transport extensions**: Reusable low-level helpers in `transport/` for attaching to specific browser targets, creating tabs, navigating existing tabs, evaluating JS on a target, capturing screenshots, and enumerating/closing targets.
- **Client API additions**: New `tap.Client` methods for browser session and tab operations so the library stays aligned with CLI capabilities.
- **CLI commands**: New `tap browser ...` command tree for session management, tab management, navigation, evaluation, and screenshots.
- **Tests**: Unit tests for the metadata store and command parsing; integration-oriented tests around target resolution and stale-session reconciliation where practical.

### Key Design Decisions

- **CLI namespace**: Use `tap browser ...` rather than bare top-level commands because sessions and tabs require additional management commands and the namespace scales better.
- **Session model**: A session represents one persistent browser instance managed by `tap`, with its own profile directory and metadata.
- **State/config contract**: The browser state root must be configurable (for tests and advanced users). Session creation freezes the session’s core runtime configuration into metadata: local vs remote mode, remote endpoint if applicable, profile directory, and browser-launch settings required for deterministic reconnect behavior. Later commands may override only per-operation concerns that do not change session identity.
- **Tab model**: A tab is an explicitly named tracked CDP target within a session. Commands may target a tab by name or fall back to the session’s selected tab.
- **Session resolution behavior**: Commands that operate on a session will resolve in this order: explicit `--session`, globally selected/default session, error with guidance if no session is available or the reference is ambiguous. The CLI will include explicit session-selection UX (`tap browser session select <name>`) so the default session is user-controlled and persisted. Session names are globally unique; tab names are unique within a session.
- **Tab resolution behavior**: Commands that need a tab will resolve in this order: explicit `--tab`, session’s selected tab, error with guidance if no tab is available.
- **Metadata authority**: Stored metadata is the source of user-facing names and defaults, but each command must reconcile it against live CDP state before acting.
- **Tracked vs untracked tabs**: `tap` manages only tracked tabs recorded in session metadata. `tab list` and `session info` will show tracked tabs with an explicit status (`live`, `stale`, `closed`). Untracked live browser tabs are ignored by default in v1 rather than auto-adopted, preventing surprising name/selection changes.
- **Restored-tab behavior**: If a browser profile restores tabs after restart, `tap` will not auto-adopt them. A tracked tab whose target disappears becomes `stale`; future enhancements can add explicit import/adopt commands if needed.
- **Crash/restart reconciliation**: If a local managed browser process dies, `tap` will reconnect or relaunch the session browser, preserve session metadata, mark tabs without matching live targets as stale, clear invalid selected-tab references, and require the user to recreate or reselect tabs if automatic restoration is not possible.
- **Session lifecycle contract**: `session close` will mean: detach from or terminate the managed browser for that session, remove live endpoint/process bindings from metadata, clear selected-session references that point to it, and keep durable session history only if explicitly needed for diagnostics; by default the session record and managed profile directory are removed for local sessions. Remote sessions remove `tap` metadata only and never attempt to terminate the remote browser.
- **Tab lifecycle contract**: `tab close` removes the tracked target and metadata for that tab; if it was selected, the session will select the next remaining tab by deterministic creation order, or leave no selected tab if none remain.
- **Remote browser support**: Allowed, but constrained: Phase 1 will define a capability matrix covering `session new`, `session close`, `tab new`, `tab close`, `navigate`, `evaluate`, and `screenshot` against remote CDP. If a remote flow cannot satisfy the lifecycle contract cleanly, the command should fail with a clear explanation rather than partially working.
- **Process ownership safety**: Local session metadata must include enough launch markers to verify ownership before killing or reusing a stored PID/endpoint, preventing `session close` from terminating an unrelated process after PID reuse.
- **Backward compatibility**: Existing transport methods remain intact for `site`, `fetch`, and `login`; the new persistent path is additive.

## Implementation Phases

### Phase 1: Design session/tab model and metadata store

1. Define persistent browser UX and command surface for `browser session`, `browser tab`, `browser navigate`, `browser evaluate`, and `browser screenshot`, including explicit session resolution rules, `session select`, and selected/default session behavior (files: `cmd/tap/main.go`, new `cmd/tap/browser.go`, docs in `README.md`).
2. Define and document the exact lifecycle contracts for `session close`, `tab close`, stale sessions, tracked vs untracked tabs, and stale/restored tabs so validation, runtime behavior, and help text share the same semantics (files: planning/docs updates in `README.md`, CLI help text, new browser package docs/comments as needed).
3. Define and document the local-vs-remote capability matrix up front, including remote session creation syntax, config precedence (`session metadata` over global flag/env after creation), persisted endpoint behavior, and expected connection/auth/TLS failures (files: planning/docs updates in `README.md`, CLI help text, new browser package docs/comments as needed).
4. Add a new `browser/` package with types for session records, tab records, durable metadata persistence, selection rules, atomic store updates, interprocess locking, and stale-state reconciliation (files: new `browser/*.go`).
5. Add unit tests for metadata CRUD, default session/tab resolution, lock/write behavior, stale tab/session handling, and validation rules (files: new `browser/*_test.go`).

### Phase 2: Implement persistent browser runtime and transport support

1. Extend `transport/` with target-aware CDP helpers for browser attach/reconnect, target enumeration, tab creation/closure, navigation on existing targets, JS evaluation on a target, and screenshots (files: `transport/transport.go` or split transport helpers if needed).
2. Implement local browser process launch/reconnect logic for persistent sessions, including profile directory allocation, debugging-port management, PID/endpoint recording, launch/ownership markers, and liveness checks (files: new `browser/runtime.go`, `browser/process.go`, and related tests/mocks as appropriate).
3. Implement session manager operations that combine metadata + transport runtime, including reconciliation when tabs or browser processes disappear, relaunch semantics after crashes, and safe handling of stale PIDs/endpoints (files: new `browser/manager.go`, tests in `browser/*_test.go`).

### Phase 3: Expose library API and CLI commands

1. Add `tap.Client` methods for browser session/tab lifecycle and browser actions while preserving current client initialization patterns, including options for browser state root and session-config persistence (files: `tap.go`, `options.go` if additional browser-store options are needed, new tests if applicable).
2. Add CLI command tree for:
   - `tap browser session new|list|info|select|close`
   - `tap browser tab new|list|select|close`
   - `tap browser navigate`
   - `tap browser evaluate`
   - `tap browser screenshot`
   (files: new `cmd/tap/browser*.go`, `cmd/tap/main.go`, shared CLI helpers as needed).
3. Add screenshot output handling (`--output` with sensible default), evaluation formatting integration, and clear UX around selected/default tab behavior (files: `cmd/tap/output.go`, new browser CLI files).
4. Update `README.md` and command help text to document the new persistent browser workflow and the separation from one-shot commands.

### Phase 4: Validation and hardening

1. Add/expand tests covering CLI validation, session/tab resolution, stale metadata recovery, and error messaging for remote/local runtime edge cases (files: `cmd/tap/*_test.go`, `browser/*_test.go`, `transport/*_test.go` where feasible).
2. Run formatting, linting, and full test suite; fix issues found during validation (commands: `mise run fmt`, `mise run lint`, `go test ./... -timeout 60s -race`).
3. Update session docs (`handoff.md`) and prepare implementation notes for follow-up enhancements such as top-level aliases or richer tab inspection.

## Testing Strategy

- Unit test metadata store create/read/update/delete flows for sessions and tabs.
- Unit test selected-session and selected-tab resolution with explicit flags, `session select`, implicit defaults, missing defaults, and ambiguous references.
- Unit test metadata store locking and atomic write behavior under repeated/competing command access.
- Unit test stale metadata reconciliation when a saved target no longer exists.
- Unit test crash/restart semantics: dead browser process, relaunched browser, stale endpoints, cleared selected tabs, and stale tab markers.
- Unit test session validation rules for duplicate names, invalid names, missing sessions, and missing tabs.
- Unit test process ownership checks so stale PID metadata cannot terminate an unrelated process.
- Unit test screenshot path resolution and evaluate output formatting decisions.
- Integration-style tests for target-aware transport helpers using mocked or controlled CDP interactions where practical.
- Automated cross-invocation integration tests that run separate `tap` processes against the same durable state directory and verify session creation, session selection, tab creation/selection, navigate/evaluate/screenshot, and crash/relaunch reconciliation.
- CLI tests for help text, argument validation, and command dispatch.
- Manual validation on local Chrome for end-to-end flow:
  - create session
  - select session
  - create/select tab
  - navigate
  - evaluate
  - screenshot
  - close tab/session
- Manual validation in remote CDP mode for supported lifecycle flows and clear failures for unsupported ones.

## Risks

| Risk | Impact | Mitigation |
| ---- | ------ | ---------- |
| Persistent browser runtime is more complex than the current one-shot model | High | Isolate logic in a dedicated `browser/` package with clean metadata and runtime boundaries; keep one-shot transport behavior untouched |
| chromedp target/session APIs may require lower-level CDP usage for tab attachment | High | Encapsulate target-aware logic in transport helpers and add focused tests before wiring CLI behavior |
| Stale metadata after browser crashes or user-closing tabs can make UX unreliable | High | Reconcile metadata against live targets on every command, define deterministic crash/restart semantics, and return actionable recovery errors |
| Remote CDP lifecycle may not support the same guarantees as locally launched Chrome | Medium | Define a capability matrix in Phase 1, validate capabilities early, document constraints, and fail clearly when invariants cannot be met |
| Concurrent CLI invocations could corrupt session state | Medium | Use interprocess locking and atomic writes in the metadata store, with tests covering repeated access patterns |
| PID reuse could make `session close` terminate the wrong process | Medium | Record launch ownership markers and verify them before kill/reconnect operations |
| Command surface may become confusing with sessions and tabs added at once | Medium | Use a dedicated `browser` namespace, consistent flags, and explicit help/examples in README and CLI descriptions |
| Screenshot/evaluate output behavior may be inconsistent with existing commands | Low | Reuse existing output helpers where possible and define predictable defaults in CLI docs |

## Open Questions

- [x] Use a `browser` namespace instead of top-level commands to leave room for session/tab management.
- [x] Support sessions and tabs from day one, with selected-tab fallback when `--tab` is omitted.
- [x] Treat `site`/`fetch` as one-shot workflows and keep the new browser workflow separate.
- [x] Assume one browser instance per local session with its own profile directory and persisted metadata.
- [x] Assume `screenshot` will support `--output` plus a generated default file path.
- [x] Assume `evaluate` will reuse existing formatting conventions rather than inventing a separate output mode.

## Review Feedback

Round 1 reviewer feedback addressed:
- Added explicit session resolution rules alongside tab resolution.
- Added atomic write + interprocess locking requirements for the metadata store.
- Added crash/restart reconciliation semantics for local managed browsers.
- Added a concrete remote-mode capability-matrix deliverable in Phase 1.
- Added process ownership verification requirements for local session teardown.

Round 2 reviewer feedback addressed:
- Added explicit `session select` UX for persistent default-session behavior.
- Moved persistent metadata expectation from cache storage to durable app state/data storage.
- Defined explicit `session close` and `tab close` lifecycle semantics for local and remote sessions.

Round 3 reviewer feedback addressed:
- Defined concrete remote session creation/configuration behavior and config precedence.
- Defined tracked-vs-untracked and stale/restored-tab reconciliation rules for `tab list` and session inspection.
- Added automated cross-invocation integration coverage for the persistence contract.

Additional follow-up incorporated before user review:
- Added explicit `session select` to the Phase 3 CLI command plan.
- Added session-scoped runtime locking beyond atomic metadata writes.
- Added a configurable durable browser state root and a session creation-time config contract.
- Clarified that local persistent sessions require independently managed Chrome processes, with `chromedp` used for attach/control only.

## Final Status

All phases (1-4) complete.
