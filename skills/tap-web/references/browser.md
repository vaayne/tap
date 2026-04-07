# Browser Sessions & Tabs

Manage long-lived browser instances that survive across CLI invocations.

## Sessions

```bash
tap browser session new <name>                # Launch headless Chrome
tap browser session new <name> --no-headless  # Visible browser
tap browser session new <name> --ws-url <url> # Remote CDP endpoint
tap browser session list
tap browser session info [name]
tap browser session close [name]
```

## Tabs

```bash
tap browser tab new <name> [--url <url>]
tap browser tab list
tap browser tab select <name>
tap browser tab close [name]
```

## Actions

All accept `--session` and `--tab` to override defaults.

### Page actions

```bash
tap browser navigate <url>
tap browser back
tap browser forward
tap browser reload
tap browser evaluate <javascript>
tap browser screenshot [--output <path>]
tap browser text [selector]                    # Clean readable text (defuddle)
tap browser text ".main-content"               # Scoped to selector
tap browser text -f json                       # JSON with title, markdown, wordCount
tap browser pdf [--output <path>] [--landscape] [--background] [--scale 1.0]
```

### Human-like interaction

These dispatch real CDP events (mouse moves, key presses) — indistinguishable from a real user.

```bash
tap browser click <selector>                    # Full mouse event chain
tap browser type <selector> <text>              # Per-keystroke typing
tap browser hover <selector>                    # mouseMoved to element center
tap browser scroll <selector>                   # Scroll element into view
tap browser scroll --x 0 --y 1000              # Scroll to pixel position
tap browser select <selector> <value>           # Pick <select> option
tap browser wait <selector> [--timeout 30s]     # Wait for element visible
```

| Command | vs JS eval | When to use |
|---|---|---|
| `click` | Real mousedown/mouseup chain | Sites listening on mousedown, hover menus |
| `type` | Individual keyDown/keyUp events | Anti-bot detection, per-keystroke validation |
| `hover` | Real mouseMoved to coordinates | CSS :hover states, mouseenter listeners |
| `scroll` | Triggers IntersectionObserver | Lazy-loaded content, infinite scroll |
| `select` | Fires focus/input/change events | Native `<select>` elements |
| `wait` | CDP visibility polling | Wait for dynamic content before acting |

### Keyboard & dialogs

```bash
tap browser keypress Enter                  # Submit form
tap browser keypress Tab                    # Move focus
tap browser keypress Escape                 # Dismiss/close
tap browser keypress Ctrl+a                 # Select all
tap browser keypress Ctrl+c                 # Copy
tap browser dialog                          # Accept pending alert/confirm/prompt
tap browser dialog --accept=false           # Dismiss dialog
tap browser dialog --text "answer"          # Fill prompt and accept
```

Keys: Enter, Tab, Escape, Backspace, Delete, Space, ArrowUp/Down/Left/Right, Home, End, PageUp, PageDown, Ctrl, Alt, Shift, Meta. Modifiers with `+`: `Ctrl+a`, `Shift+Tab`.

### Cookies

```bash
tap browser cookies get                     # List all (includes httpOnly)
tap browser cookies get -f json             # JSON output
tap browser cookies set <name> <value>      # Set cookie
tap browser cookies set <name> <value> --domain .example.com --path /
tap browser cookies clear                   # Delete all cookies
```

### Forms

```bash
tap browser forms
tap browser fill <sel> <val> [<sel> <val>...]
tap browser fill <sel> <val> --submit <sel>
```

`tap browser forms` returns JSON with each element's `selector`, `type`, `name`, `placeholder`, `label`, `value`, `role`.

`tap browser fill` uses React-compatible native setters — works with React, Vue, Angular, vanilla HTML.

## Electron Apps

Connect tap to any Electron (or CEF-based) desktop app via CDP.

```bash
tap electron ps                              # List running processes with --remote-debugging-port
tap electron launch <binary> [app-args...]   # Launch with --remote-debugging-port=0
tap electron launch <binary> --session <name>
tap electron attach --port <port>            # Attach to running app
tap electron attach --port <port> --session <name>
tap electron discover [--session <name>]     # Adopt live windows as tracked tabs
```

After `discover`, all `tap browser` commands work against the session. macOS `.app` bundles: pass the inner binary (`/Applications/MyApp.app/Contents/MacOS/MyApp`).

**Note:** Cherry Studio and similar apps use `<i role="button">` for the send icon rather than `<button>` — use `tap browser evaluate` to inspect the DOM when `click` targets need investigation.

## Resolution

When `--session`/`--tab` omitted, tap resolves automatically:

- **Session**: flag → `default` (auto-created if missing)
- **Tab**: flag → selected tab → the only live tab

## Session Strategy

**The `default` session is always used when `--session` is omitted.** If it doesn't exist, tap auto-creates a headless `default` session. The profile directory is derived from the session name, so same name = same profile = same cookies. Only use `--session <name>` for isolation (parallel subagents, different accounts).

### Quick start

```bash
tap browser tab new main --url <start-url>  # Auto-creates default session, opens a tab
```

### Recover a stale session

If `default` is unresponsive (Chrome crashed, PID gone), close and recreate. **Note: closing deletes the profile directory (including cookies).** Use `tap login` for persistent auth that survives session close.

```bash
tap browser session close default
tap browser session new default
# Re-login if needed: tap login <url>
```

### When to create a named session

- **Parallel subagents** — each agent needs its own Chrome (profile lock). Give each a unique name.
- **Account isolation** — different login states that must not interfere.

```bash
tap browser session new research-a
tap browser session new research-b
```

### When to use `-b` vs persistent sessions

| Need | Use |
|---|---|
| Single script with auth | `tap site -b <script>` — ephemeral, no session needed |
| Multi-step workflow | `tap browser` with `default` session |
| Network interception + scripting | `tap browser` — intercept rules need a persistent session |

## Examples

```bash
# Reuse default session for a workflow
tap browser session new default
tap browser tab new docs --url https://go.dev/doc
tap browser evaluate 'document.title'

# Open a second tab in the same session
tap browser tab new api --url https://pkg.go.dev
tap browser tab select docs
tap browser evaluate 'document.querySelectorAll("a").length'

# Form filling (visible browser for manual CAPTCHA)
tap browser session new login --no-headless
tap browser tab new page --url https://example.com/login
tap browser forms
tap browser fill "#email" "me@example.com" "#password" "secret" --submit "button[type=submit]"
tap browser session close login
```
