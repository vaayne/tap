# Tasks: Browser Sessions and Tabs

- [x] Phase 1 — Design session/tab model and metadata store
  - [x] Define final CLI command tree and argument/flag defaults, including `session select`
  - [x] Add `browser/` package types for sessions, tabs, and durable persisted state
  - [x] Implement metadata store CRUD + selected-session/selected-tab resolution
  - [x] Implement session-scoped runtime locking and configurable state-root support
  - [x] Add metadata-focused unit tests

- [ ] Phase 2 — Implement persistent browser runtime and transport support
  - [ ] Add target-aware CDP helpers in `transport/`
  - [ ] Implement local browser process launch/reconnect support
  - [ ] Implement session manager reconciliation logic and crash/restart handling
  - [ ] Add runtime and reconciliation tests

- [ ] Phase 3 — Expose library API and CLI commands
  - [ ] Add `tap.Client` browser session/tab methods and state-root/session-config options
  - [ ] Add `tap browser session ...` commands, including `session select`
  - [ ] Add `tap browser tab ...` commands
  - [ ] Add `tap browser navigate`
  - [ ] Add `tap browser evaluate`
  - [ ] Add `tap browser screenshot`
  - [ ] Update help text and README examples

- [ ] Phase 4 — Validation and hardening
  - [ ] Add CLI validation/error-path tests
  - [ ] Run `mise run fmt`
  - [ ] Run `mise run lint`
  - [x] Run `go test ./... -timeout 60s -race`
  - [x] Update session handoff notes
