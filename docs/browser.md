# Persistent Browser Sessions

Tap provides persistent browser automation via `tap browser`. Sessions and tabs survive across CLI invocations, letting you navigate pages, evaluate JavaScript, and capture screenshots against long-lived browser state.

## Quick Start

```bash
# Create a session (launches a managed headless Chrome)
tap browser session new work

# Open a tab and navigate
tap browser tab new main --url https://example.com
tap browser navigate https://httpbin.org/html

# Interact with the page
tap browser evaluate 'document.querySelector("h1").textContent'
tap browser screenshot --output page.png

# Clean up
tap browser tab close main
tap browser session close work
```

## Commands

### Sessions

```bash
tap browser session new <name>              # Create a local session
tap browser session new <name> --no-headless  # Visible browser
tap browser session new <name> --ws-url <url> # Remote CDP session
tap browser session list                    # List all sessions
tap browser session info [name]             # Show session details
tap browser session select <name>           # Set the default session
tap browser session close [name]            # Close and remove a session
```

### Tabs

```bash
tap browser tab new <name> [--url <url>]    # Create a tracked tab
tap browser tab list                        # List tracked tabs
tap browser tab select <name>               # Set the default tab
tap browser tab close [name]                # Close and remove a tab
```

### Actions

```bash
tap browser navigate <url>                  # Navigate the selected tab
tap browser evaluate <javascript>           # Run JS and print the result
tap browser screenshot [--output <path>]    # Capture a full-page PNG
tap browser forms                           # Discover fillable form elements
tap browser fill <sel> <val> [<sel> <val>]  # Fill form fields
```

All action commands accept `--session <name>` and `--tab <name>` to override the defaults.

## Resolution Rules

When `--session` or `--tab` is omitted, tap resolves them automatically:

**Session resolution order:**
1. `--session` flag
2. The selected session from `tap browser session select`
3. The only available session, when exactly one exists

**Tab resolution order:**
1. `--tab` flag
2. The selected tab within the resolved session
3. The only live tracked tab, when exactly one exists

If the resolution is ambiguous, tap fails with guidance instead of guessing.

## Local vs Remote Sessions

| Operation | Local managed browser | Remote CDP session |
|---|---|---|
| `session new` | Launches Chrome with a dedicated profile and debug endpoint | Persists the `--ws-url` and validates the connection at creation time |
| `session close` | Verifies ownership, stops Chrome, removes metadata and profile | Removes tap metadata only; never kills the remote browser |
| `tab new/close` | Supported | Supported when the remote endpoint allows target management |
| `navigate/evaluate/screenshot` | Operate on the resolved tracked tab | Operate through the saved endpoint, ignoring later global `--ws-url` overrides |

### Remote sessions

```bash
# Connect to a remote browser
tap browser session new remote-box --ws-url wss://remote:9222/devtools/browser/abc

# Use it like a local session
tap browser tab new t1 --url https://example.com
tap browser evaluate 'document.title'
```

Remote sessions freeze the `--ws-url` from creation time. Later global `--ws-url` flags or `TAP_WS_URL` environment changes are ignored for reconnects.

## Lifecycle Rules

- A **session** is one persistent browser instance (local or remote).
- A **tab** is a named tracked CDP target within a session.
- Only tracked tabs are part of tap metadata. Untracked live browser tabs are ignored.
- If a tracked target disappears, tap marks the tab **stale** and clears invalid selected-tab state.
- `tab close` promotes the next remaining live tracked tab by creation order, or leaves no selected tab if none remain.
- `session close` for local sessions verifies browser-process ownership before terminating Chrome and deleting its profile directory.

## Configuration

| Variable | Flag | Description | Default |
|---|---|---|---|
| `TAP_BROWSER_STATE_DIR` | `--state-root` | Browser metadata directory | `$XDG_CACHE_HOME/tap/browser` or `~/.cache/tap/browser` |

## Output Formatting

`tap browser evaluate` supports the same output formats as `tap site`:

```bash
tap browser evaluate 'document.title'                    # Pretty (default)
tap browser evaluate --format json 'document.title'      # JSON
tap browser evaluate --format raw 'document.title'       # Raw value
```

## Examples

### Multi-tab workflow

```bash
tap browser session new research

tap browser tab new docs --url https://go.dev/doc
tap browser tab new api --url https://pkg.go.dev

# Switch between tabs
tap browser tab select docs
tap browser evaluate 'document.title'

tap browser tab select api
tap browser evaluate 'document.title'

tap browser session close research
```

### Screenshot automation

```bash
tap browser session new screenshots
tap browser tab new page --url https://example.com

# Navigate and capture
tap browser navigate https://example.com/page1
tap browser screenshot --output page1.png

tap browser navigate https://example.com/page2
tap browser screenshot --output page2.png

tap browser session close screenshots
```

### Non-headless debugging

```bash
# Launch a visible browser for debugging
tap browser session new debug --no-headless
tap browser tab new test --url https://example.com

# Interact manually in the visible browser, then:
tap browser evaluate 'document.querySelector(".result").textContent'
tap browser session close debug
```

### Form discovery and filling

```bash
tap browser session new forms-demo --no-headless
tap browser tab new page --url https://example.com/login

# Discover fillable elements on the page
tap browser forms
# Returns JSON with selector, type, name, placeholder, label, role for each element

# Fill fields using the reported selectors
tap browser fill "#username" "myuser" "#password" "secret"

# Fill and submit in one command
tap browser fill "#username" "myuser" "#password" "secret" --submit "button[type=submit]"

tap browser session close forms-demo
```

`tap browser forms` reports all fillable elements (inputs, textareas, selects, buttons) with their best CSS selector, type, label, placeholder, current value, and role (`text`, `toggle`, `select`, `submit`). Use the selectors directly with `tap browser fill`.

`tap browser fill` uses React-compatible native value setters with proper `input`/`change` event dispatch, so it works with React, Vue, Angular, and vanilla HTML forms.

## Browser Backends

| | Chrome | Lightpanda (`--lp`) |
|---|---|---|
| **Platforms** | macOS, Linux, Windows | macOS, Linux |
| **Install** | Manual | `tap doctor --install` |
| **Update** | — | `tap doctor --update` |
| **Sessions & cookies** | Yes | No |
| **Network interception** | Yes | No |
| **Site compatibility** | All | Partial — nightly builds, not all sites render correctly |
| **Best for** | Auth, full automation | Fast headless JS rendering without auth |
