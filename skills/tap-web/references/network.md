# Network Capture & Interception

Capture, inspect, and intercept network requests on the current browser tab via agent-browser.

## Commands

```bash
# List captured requests
tap browser network requests
tap browser network requests --filter "*/api/*"

# Route/block requests
tap browser network route "*.ads.*" --abort
tap browser network route "*/api/mock" --body '{"ok":true}'

# Remove routes
tap browser network unroute
tap browser network unroute "*.ads.*"

# HAR capture
tap browser network har start /tmp/session.har
tap browser network har stop
```

## Quick workflow

```bash
tap browser open https://example.com --show
tap browser network requests --filter "*/api/*"
```

## Notes

- Prefer `network requests --filter api` over scraping DOM when the site has clean JSON APIs.
- Network interception runs through agent-browser against managed or attached Chrome sessions.
- Pass-through commands accept agent-browser flags directly; use `--session <name>` to target a named tap session.

## Example: capture SPA data

```bash
tap browser open https://x.com/elonmusk --show
tap browser network requests --filter "UserTweets"
```

## Example: capture after a form action

```bash
tap browser open https://example.com/search --show
tap browser fill "#search" "golang"
tap browser keypress Enter
tap browser network requests --filter "api/search"
```
