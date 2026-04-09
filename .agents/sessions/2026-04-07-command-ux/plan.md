# Plan: CLI Command UX Redesign

## Overview

Redesign tap's public CLI surface so common jobs are task-first and discoverable while preserving current advanced browser/CDP capabilities behind compatibility commands and advanced subcommands.

### Goals

- Make the default CLI model centered around user jobs instead of CDP/session internals.
- Simplify browser usage around one default browser context and one current tab.
- Make "use my existing Chrome", "log in once and reuse", and "attach to Electron" first-class flows.
- Preserve current advanced capabilities (`session`, `tab`, `proxy`, network tooling) for power users and automation.
- Ship the redesign incrementally with backward compatibility.

### Success Criteria

- [ ] A new top-level command map is coherent across `site`, `fetch`, `login`, `browser`, and `electron`.
- [ ] Common workflows take fewer steps and avoid explicit session/tab setup.
- [ ] Advanced capabilities remain available without breaking existing scripts immediately.
- [ ] Help text and docs clearly distinguish common vs advanced flows.
- [ ] Migration path is explicit, with aliases and deprecation sequencing.
- [ ] Default-context resolution is deterministic and never silently switches browser/account state.

### Out of Scope

- A full transport/browser-engine rewrite.
- Immediate removal of current `browser session`, `browser tab`, `browser proxy`, or `electron` commands.
- Deep changes to site-script runtime behavior unrelated to CLI UX.

## Technical Approach

Adopt a three-layer command model:

1. **Task-first layer** for common use
   - user-facing commands focused on actions and workflows
   - one default browser context and one current tab mental model
2. **Compatibility layer**
   - existing commands kept working
   - aliased to the new mental model where possible
3. **Advanced layer**
   - explicit session/tab/proxy/network/CDP operations remain available
   - documented as advanced

### Product Model

#### Core user concepts

- **Default browser context**: the browser state used by `login`, browser-backed `fetch/site`, and plain `tap browser` commands unless overridden.
- **Current tab**: the default working page inside the resolved browser context.
- **Attached browser/app**: an existing Chrome/Electron target adopted by tap for reuse.

#### Internal concepts kept but demoted

- named sessions
- tracked tabs
- CDP endpoints / WebSocket URLs
- local proxy plumbing
- remote/local mode distinctions

## Persisted State Contract

This is the key new UX contract and must be implemented explicitly.

### New persisted metadata

Store a small default-context record alongside the existing browser store under the same state root and lock discipline.

Suggested shape:

```json
{
  "version": 1,
  "activeContext": {
    "kind": "managed|attached-chrome|attached-electron",
    "sessionName": "default",
    "source": {
      "type": "managed-profile|proxy|browser-url|electron-port|electron-launch",
      "value": "..."
    },
    "status": "ready|stale",
    "savedAt": "..."
  }
}
```

Use the existing browser state/store lock so concurrent CLI invocations cannot partially update active-context state.

### Persistence rules

- **Persistent state changes** are made only by explicit stateful commands:
  - `tap attach ...`
  - `tap attach clear`
  - `tap login clear`
  - future explicit `tap browser use ...` style commands if added
- **One-shot overrides** do **not** mutate the persisted default:
  - `tap fetch/site/browser --browser-url ...`
  - `tap fetch/site/browser --profile-dir ...`
  - `tap browser --session ...`
  - `tap browser --tab ...`
- `tap login <url>` writes auth into the resolved browser context but does not change which context is active unless `login` itself is given a future explicit `--make-default` option.

### Resolution precedence

For any browser-backed command:

1. explicit command override (`--session`, `--context`, `--browser-url`, `--profile-dir`)
2. persisted active context
3. managed local default session/profile

### Invalidation rules

- If an attached target becomes unreachable, mark persisted active context as `stale`.
- `tap status` / `tap attach status` / `tap login status` must surface `stale` explicitly.
- Commands must **not** silently fall back from a stale attached context to a managed context, because that can cross account boundaries.
- Recovery is explicit: `tap attach clear`, `tap attach chrome`, or an explicit one-shot override.

## Tab Identity Contract

The common tab layer needs stable identity independent of fluctuating titles/URLs.

### Stable tab IDs

- Every tab exposed in the common UX gets a stable user-facing ID within its context: `tab-1`, `tab-2`, ...
- `tap browser open --new-tab` auto-generates the next ID.
- A future `--name <id>` may allow explicit naming, but exact ID matching is the default.

### `tap browser tabs` output

Human output columns:
- `ID`
- `TITLE`
- `URL`
- `CURRENT`
- `STATUS`

Machine output:
- `tap browser tabs --json`

### `switch` resolution

- `tap browser switch <id>` accepts exact tab ID.
- Optional 1-based numeric index may be supported in human mode only, but exact ID is preferred.
- Title/URL fuzzy matching is out of scope for the first redesign; it introduces ambiguity.

### Current-tab transitions

- `open <url>` without `--new-tab` navigates the current tab.
- If there is no current tab, `open` creates `tab-1`.
- `close-tab` promotes the next live tab by creation order.
- If no live tabs remain, current tab becomes unset and later commands fail with action-oriented guidance.

## Supported Attach Contract

### `tap attach chrome`

This command does **not** magically control arbitrary running Chrome. Supported cases are:

1. Chrome/Chromium already running with remote debugging enabled and discoverable via `DevToolsActivePort`
2. explicit DevTools endpoint via `--browser-url`
3. explicit `--port-file`

If none apply, `tap attach chrome` fails with actionable guidance. It does **not** relaunch the user's existing browser automatically in the first implementation.

### `tap attach electron`

Supported cases are:

1. `--port <port>` where `http://127.0.0.1:<port>/json/version` exposes a browser CDP endpoint
2. `--launch <binary>` which launches an app with a Chrome DevTools-compatible endpoint

Important: `--port` means **CDP DevTools port**, not generic Node inspector. The docs/examples must stop implying that any inspector port works.

### Attach persistence

- `tap attach ...` persists the attached target as the active default context.
- `tap attach clear` only removes tap's saved attachment metadata; it never deletes cookies or kills the external browser/app.
- `tap attach status --json` reports the exact source, resolved endpoint, status, and target type.

## Proposed Command Taxonomy

### Top-level commands

```text
tap
├── site        # structured extraction from known sites
├── fetch       # clean readable content from arbitrary URLs
├── login       # log into a site in the default browser context
├── browser     # open pages and automate the current browser context
├── attach      # connect tap to existing Chrome/Electron/browser targets
├── status      # show the active browser/auth context and current tab
├── doctor      # dependency and environment checks
├── upgrade     # update tap
└── help
```

`tap status` is a convenience summary over `attach status` + browser context resolution.

### 1) `tap site`

Purpose stays the same: structured data from known sites.

#### Proposed commands

```text
tap site list
tap site search <query>
tap site info <script>
tap site sync
tap site run <script> [key=value ...]
```

#### UX notes

- `tap site <script> ...` remains supported as the primary compatibility shorthand.
- Docs can keep `tap site <script>` as the main example if retraining users seems unnecessary.
- `run` exists mainly for discoverability/help-tree consistency.

#### Proposed flags

Command-scoped:
- `-f, --format <pretty|json|raw>`
- `-b, --browser` -- run through browser path
- `--show` alias for visible browser execution when applicable
- `--wait <duration>` friendly alias for fixed delay
- `--wait-selector <css>`
- `--wait-js <expr>`
- `--timeout <duration>`

Advanced/compatibility:
- `--profile-dir <path>`
- `--browser-url <url>` alias of `--ws-url`
- hidden compatibility: `--ws-url`

### 2) `tap fetch`

Purpose stays the same: readable content extraction.

#### Proposed commands

```text
tap fetch <url>
```

Keep it single-command.

#### Proposed flags

Core:
- `--json`
- `-b, --browser`
- `--show`
- `--wait <duration>` alias of current `--delay`
- `--wait-selector <css>`
- `--wait-js <expr>`
- `--timeout <duration>`

Backend selection:
- `--lightpanda, --lp`

Advanced:
- `--browser-url <url>` alias of `--ws-url`
- `--profile-dir <path>`

#### UX notes

- Keep the "simple by default" single command.
- Rename `--delay` in the UX to `--wait`; keep `--delay` as compatibility alias.
- `-b` means "use the browser path and its auth/browser context", not just "skip QuickJS".

### 3) `tap login`

Purpose: establish reusable auth/browser state in the default browser context.

#### Proposed commands

```text
tap login <url>
tap login status
tap login clear
```

#### Proposed flags

Core:
- `--show` (default true for login if using a managed browser)
- `--timeout <duration>`
- `--json` for `status`

Advanced:
- `--profile-dir <path>`
- `--browser-url <url>`

#### Clear semantics

- `tap login clear` deletes tap-managed saved auth/profile state for the managed default context only.
- It does **not** clear cookies from an attached external Chrome/Electron target.
- If the current default context is attached, `login clear` reports that attached state is external and suggests `tap attach clear` instead.

### 4) `tap attach`

New top-level family. This is the key simplification for "use my existing Chrome" and Electron.

#### Proposed commands

```text
tap attach chrome
tap attach chrome --browser-url <devtools-url>
tap attach chrome --port-file <DevToolsActivePort>
tap attach electron --port <port>
tap attach electron --launch <binary> [-- <app-args...>]
tap attach status
tap attach clear
```

#### UX notes

- `tap attach chrome` should auto-discover the user's Chrome/Chromium **only when** remote debugging is already enabled and discoverable.
- Internally it may start/reuse the built-in proxy, but that is not exposed as the default workflow.
- `tap attach electron --port` attaches to a running Electron app with a valid CDP endpoint.
- `tap attach electron --launch` launches and attaches in one step.
- `status` shows the current attached target and what later commands will reuse.
- `clear` detaches tap from the current external target without mutating the external app/browser.

#### Proposed flags

Chrome attach:
- `--browser-url <devtools-url>` alias `--url`
- `--port-file <path>`
- hidden compatibility: `--user-chrome`, `--devtools-port-file`

Electron attach:
- `--port <cdp-port>`
- `--launch <binary>`
- `--name <context-name>` optional for advanced named contexts

Advanced:
- `--listen <addr>` for internal proxy bind address
- hidden compatibility: `--upstream`

### 5) `tap browser`

Reframe around current browser context + current tab.

#### Proposed commands

```text
tap browser open <url>
tap browser tabs
tap browser switch <tab-id>
tap browser close-tab [tab-id]
tap browser status

tap browser back
tap browser forward
tap browser reload

tap browser text [selector]
tap browser evaluate <js>
tap browser screenshot [--output <path>]
tap browser pdf [--output <path>]

tap browser click <selector>
tap browser type <selector> <text>
tap browser fill <selector> <value> [<selector> <value> ...]
tap browser hover <selector>
tap browser scroll [selector]
tap browser select <selector> <value>
tap browser wait <selector>
tap browser keypress <key>
tap browser dialog

tap browser forms
tap browser cookies get|set|clear
tap browser network wait|log|body|intercept|clear
```

#### Advanced subcommands retained

```text
tap browser session ...   # advanced
tap browser tab ...       # advanced
tap browser proxy ...     # advanced / hidden from primary docs
```

#### Core behavior defaults

- `open <url>` uses the resolved browser context.
- If no current tab exists, `open` creates one automatically.
- If a current tab exists, `open` navigates it by default.
- Multi-tab workflows use `tabs` and `switch`; explicit named tracked tabs remain advanced.
- `status` shows browser mode, attachment source, current tab, and context status in user terms.

#### Proposed flags

Common browser-action flags:
- `--tab <id>` advanced override
- `--context <name>` with `--session` kept as compatibility alias
- `--new-tab` for `open`
- `--show` for visible local browser where applicable
- `--timeout <duration>` where relevant
- `--json` for `status` and `tabs`

Output flags:
- `-f, --format <pretty|json|raw>` for `text`, `evaluate`, `forms`, `cookies`, `network wait/body`
- `-o, --output <path>` for `screenshot`, `pdf`

Advanced/compatibility:
- `--session` kept supported
- hidden compatibility for current browser action flags

### 6) `tap electron`

Keep as compatibility / advanced namespace during migration.

#### Proposed position

- `tap attach electron ...` becomes the primary documented path.
- `tap electron ps|attach|launch|discover` remains supported.
- `tap electron` help should explicitly say: prefer `tap attach electron` for new usage.

## Flag Model

### Keep as global

Only truly global, cross-command concerns should remain global:
- `--verbose`
- `--quiet`
- `--no-color`

### Move out of global help for user-facing commands

These should become command-scoped or advanced-scoped in docs/help:
- `--browser`, `-b`
- `--show` / `--no-headless`
- `--wait` / `--delay`
- `--wait-selector`
- `--wait-js`
- `--profile-dir`
- `--browser-url` / `--ws-url`
- `--lightpanda`

### Naming changes

- `--no-headless` -> primary alias `--show`
- `--delay` -> primary alias `--wait`
- `--ws-url` -> primary alias `--browser-url` or `--devtools-url`
- keep old names as compatibility aliases

## Behavioral Defaults

### Default browser context

`tap` should expose one clear default browser context used by:
- `tap login`
- `tap fetch -b`
- `tap site -b`
- `tap browser ...`

If the active context is attached Chrome/Electron, commands should say so in `status` and related help text.

### Current tab

Most browser commands operate on the current tab.

Resolution order:
1. explicit `--tab`
2. current tab
3. only live tab
4. error with action-oriented guidance

### Attach behavior

- `tap attach chrome` makes the attached Chrome the default browser context.
- `tap attach electron ...` makes the attached app/browser the default browser context unless `--name` creates a named one.
- `tap attach clear` resets to managed/local default behavior.

### Auth behavior

- `login` clearly reports where auth was saved and which future commands will reuse it.
- browser-backed `fetch/site` reuse the default browser context without requiring users to understand `profile-dir` or `ws-url`.
- explicit `--profile-dir` / `--browser-url` are one-shot overrides unless a future explicit persist flag is provided.

## Machine-readable Status Contract

To preserve automation while improving human UX, add `--json` to:

- `tap status`
- `tap attach status`
- `tap login status`
- `tap browser status`
- `tap browser tabs`

This lets scripts migrate onto stable machine-readable output while human-readable output evolves.

## Example Flows

### Flow 1: Use my existing Chrome

```bash
# Chrome must already expose DevTools
tap attach chrome
tap fetch -b https://example.com/private-page
tap site -b github/notifications
tap browser open https://example.com
tap browser click '#submit'
```

Expected UX:
- one attach step
- no explicit proxy or session creation
- same browser/account reused across commands
- clear failure if Chrome is not debuggable yet

### Flow 2: Log in once and reuse

```bash
tap login https://github.com/login
tap login status --json
tap site -b github/notifications
tap fetch -b https://github.com/notifications
tap browser open https://github.com
```

Expected UX:
- `login status` makes reuse visible
- browser-backed commands share the same default auth context

### Flow 3: Open page, click, extract

```bash
tap browser open https://example.com
tap browser click 'button.submit'
tap browser wait '.result'
tap browser text '.result'
```

Expected UX:
- no session/tab setup ceremony
- current tab model is obvious

### Flow 4: Persistent automation session

```bash
tap browser open https://news.ycombinator.com
tap browser open https://github.com --new-tab
tap browser tabs
tap browser switch tab-2
tap browser screenshot --output github.png
tap browser status
```

Expected UX:
- common tab management without exposing tracked-tab internals
- stable tab IDs make switching scriptable and unambiguous

### Flow 5: Attach to Electron app

```bash
tap attach electron --port 9333
tap browser tabs --json
tap browser evaluate 'document.title'
tap browser screenshot
```

Alternative launch flow:

```bash
tap attach electron --launch /Applications/MyApp.app/Contents/MacOS/MyApp
tap browser text
```

Expected UX:
- attach and use browser commands immediately
- window discovery happens automatically when possible
- docs avoid conflating Node inspector ports with browser CDP ports

### Flow 6: Advanced remote/proxy workflow

```bash
# advanced only
tap browser proxy --upstream http://127.0.0.1:9222
tap browser session new remote --browser-url http://127.0.0.1:9401
tap browser tab new page --url https://example.com --session remote
tap browser network wait --url-pattern '*/api/*' --body --session remote --tab page
```

Expected UX:
- old power-user path remains possible
- advanced docs only

## Migration Strategy

### Backward-compatible additions first

Phase the redesign without breaking scripts:

1. Add new command aliases and workflow-first help.
2. Keep existing commands working unchanged.
3. Update docs/examples/help to prefer new commands.
4. Add deprecation notes only after the new model is proven.

### Compatibility mapping

- `tap site <script> ...` -> supported, docs may also show `tap site run <script>`
- `tap browser tab list` -> alias `tap browser tabs`
- `tap browser tab select <name>` -> alias `tap browser switch <id>` where possible
- `tap browser session info` -> alias `tap browser status`
- `tap browser tab close` -> alias `tap browser close-tab`
- `tap electron attach|launch` -> retained, docs steer to `tap attach electron`
- `tap browser proxy` -> retained as advanced
- `--no-headless` / `--delay` / `--ws-url` -> retained as hidden or secondary aliases

### Compatibility guarantees

For all existing scripted paths during the migration window:
- existing command spellings remain valid
- exit codes remain stable
- existing machine-readable outputs remain stable where already defined
- new human-oriented aliases may render different text, but `--json` becomes the supported machine contract for new status/list commands

## Implementation Phases

### Phase 1: Command model and help scaffolding

1. Add new command aliases and entrypoints: `status`, `attach`, `browser open`, `browser tabs`, `browser switch`, `browser close-tab`, `browser status`.
2. Update top-level and browser help text to present task-first flows.
3. Keep existing commands intact and wired to the same underlying operations.

### Phase 2: Default-context behavior

1. Introduce a clear persisted default browser context abstraction that can represent managed browser vs attached target.
2. Make `browser open` and related commands resolve through that abstraction.
3. Add `attach status` and `login status` output that explains reuse.

### Phase 3: Attach workflow integration

1. Implement `tap attach chrome` as a wrapper over discovery/proxy/session plumbing.
2. Implement `tap attach electron` wrappers over current Electron attach/launch paths.
3. Reduce explicit `discover` needs by auto-adopting windows/tabs when safe.

### Phase 4: Docs, migration, and cleanup

1. Rewrite README and docs around primary flows.
2. Mark advanced commands clearly.
3. Add deprecation guidance for renamed flags/paths once stable.

## Testing Strategy

- Help output tests for new top-level command structure and aliases.
- Behavioral tests for `browser open` default tab resolution.
- Backward compatibility tests ensuring existing `browser session/tab` flows still work.
- Attach flow tests for Chrome and Electron wrappers where deterministic.
- Docs/examples validation against actual command behavior.

## Risks

| Risk | Impact | Mitigation |
| --- | --- | --- |
| New UX layer creates ambiguous state with old session model | High | Keep explicit advanced overrides and define one clear default-context resolution path |
| Attach abstraction hides too much for power users | Medium | Preserve advanced commands and surface `status` clearly |
| Global flag restructuring causes CLI churn | Medium | Add aliases first, deprecate later |
| Electron and Chrome attach behaviors diverge in edge cases | Medium | Use shared attach abstractions but keep target-specific fallbacks |
| Current code structure may not support clean aliasing without refactor | Medium | Phase command scaffolding before deeper behavior changes |

## Open Questions

- Should `tap site run` become the primary documented form, or should docs keep `tap site <script>` as the main UX and use `run` only as an optional explicit form?
- Should `tap attach chrome` persist across shells by default, or only for the current local state directory/session? (recommended: yes, via state root)
- Should `tap browser open` navigate the current tab by default, with `--new-tab` optional? (recommended: yes)
- Should `login status` report only managed auth state, or also attached-browser reuse state? (recommended: both)

## Review Feedback

(To be updated after plan review rounds)

## Final Status

Approved and ready for implementation.
