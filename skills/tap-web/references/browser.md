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

## Examples

```bash
# Multi-tab workflow
tap browser session new research
tap browser tab new docs --url https://go.dev/doc
tap browser tab new api --url https://pkg.go.dev
tap browser tab select docs
tap browser evaluate 'document.title'
tap browser session close research

# Form filling
tap browser session new login --no-headless
tap browser tab new page --url https://example.com/login
tap browser forms
tap browser fill "#email" "me@example.com" "#password" "secret" --submit "button[type=submit]"
tap browser session close login
```
