# Browser Sessions & Tabs

Manage long-lived browser instances that survive across CLI invocations.

## Sessions

```bash
tap browser session new <name>                # Launch headless Chrome
tap browser session new <name> --no-headless  # Visible browser
tap browser session new <name> --ws-url <url> # Remote CDP endpoint
tap browser session list
tap browser session info [name]
tap browser session select <name>             # Set default
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

```bash
tap browser navigate <url>
tap browser evaluate <javascript>
tap browser screenshot [--output <path>]
tap browser forms
tap browser fill <sel> <val> [<sel> <val>...]
tap browser fill <sel> <val> --submit <sel>
```

## Resolution

When `--session`/`--tab` omitted, tap resolves automatically:

- **Session**: flag → selected → the only session
- **Tab**: flag → selected tab → the only live tab

## Forms

`tap browser forms` returns JSON with each element's `selector`, `type`, `name`, `placeholder`, `label`, `value`, `role`.

`tap browser fill` uses React-compatible native setters — works with React, Vue, Angular, vanilla HTML.

## Session Strategy

**Always reuse the `default` session.** The profile directory is derived from the session name, so same name = same profile = same cookies. Only create named sessions for isolation (parallel subagents, different accounts).

### Bootstrap the default session

```bash
tap browser session list                    # Check if default exists
tap browser session new default             # Create if missing
tap browser tab new main --url <start-url>  # Open a tab
```

### Recover a stale session

If `default` is unresponsive (Chrome crashed, PID gone), close and recreate. Cookies are preserved because the profile directory survives:

```bash
tap browser session close default
tap browser session new default
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
