# Network Interception

Tap forwards network commands to agent-browser for the current managed or attached browser session.

## Commands

```bash
# List captured requests
tap browser network requests
tap browser network requests --filter "*/api/*"

# Route or block matching requests
tap browser network route "*.ads.*" --abort
tap browser network route "*/api/mock" --body '{"ok":true}'

# Remove routes
tap browser network unroute
tap browser network unroute "*.ads.*"

# HAR capture
tap browser network har start /tmp/session.har
tap browser network har stop
```

All arguments after the tap subcommand are passed through to agent-browser. Use `--session <name>` with these pass-through commands to target a named tap session.

## Common workflows

### Inspect API calls

```bash
tap browser open https://example.com --show
tap browser network requests --filter api
```

### Block noisy requests

```bash
tap browser network route "*.doubleclick.net*" --abort
tap browser reload
```

### Clear interception rules

```bash
tap browser network unroute
```

## Notes

- Request history is collected by the browser backend for the active session.
- Route rules apply while the agent-browser session is active.
- For exact flag support, run `agent-browser network <action> --help` or `agent-browser skills get core --full`.
