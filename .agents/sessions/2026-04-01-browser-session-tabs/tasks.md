# Tasks: Browser Sessions and Tabs

- [x] Phase 1 — Design session/tab model and metadata store
  - [x] Define final CLI command tree and argument/flag defaults, including `session select`
  - [x] Add `browser/` package types for sessions, tabs, and durable persisted state
  - [x] Implement metadata store CRUD + selected-session/selected-tab resolution
  - [x] Implement session-scoped runtime locking and configurable state-root support
  - [x] Add metadata-focused unit tests

- [x] Phase 2 — Implement persistent browser runtime and transport support
  - [x] Add target-aware CDP helpers in `transport/`
  - [x] Implement local browser process launch/reconnect support
  - [x] Implement session manager reconciliation logic and crash/restart handling
  - [x] Add runtime and reconciliation tests

- [x] Phase 3 — Expose library API and CLI commands
  - [x] Add `tap browser session ...` commands, including `session select`
  - [x] Add `tap browser tab ...` commands
  - [x] Add `tap browser navigate`
  - [x] Add `tap browser evaluate`
  - [x] Add `tap browser screenshot`

- [x] Phase 4 — Validation and hardening
  - [x] Run `mise run fmt`
  - [x] Run `mise run lint`
  - [x] Update session handoff notes
