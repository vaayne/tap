---
name: tap-web
description: >
  Access websites, search the web, and extract clean content using the `tap` CLI.
  Also manage persistent browser sessions for long-lived automation workflows.
  Use when the user asks to search the web, read a webpage, fetch article content,
  get trending topics, look up social media posts, check stock prices, search videos,
  retrieve structured data from any supported site, manage browser sessions/tabs,
  capture or intercept network requests, block or mock API calls.
  Triggers on: "search for", "what's trending on", "fetch this page", "read this URL",
  "get content from", "look up on Twitter/Weibo/Reddit/YouTube/etc.", "check stock",
  "translate", "browser session", "open a tab", "navigate to", "take a screenshot",
  "evaluate javascript", "persistent browser", "fill form", "form fields",
  "discover inputs", "intercept requests", "capture network", "block requests",
  "mock API", "network log", "wait for request", or any web access task.
---

# tap-web

Use the `tap` CLI to access websites from the terminal.

## `tap fetch` — Clean content from any URL

```bash
tap fetch <url>              # Clean markdown
tap fetch --json <url>       # JSON with metadata
```

## `tap site` — Structured data via site scripts

**Always discover scripts first. Do not guess names or parameters.**

```bash
tap site list                         # List all scripts
tap site search <keyword>             # Search by keyword
tap site info <script>                # Show params and usage
tap site <site/action> [key=value]    # Run a script
tap site <script> -f json | jq '.'   # JSON output
```

To write or contribute a new script, see [references/script-development.md](references/script-development.md).

## Login & browser backends

```bash
tap login https://github.com/login   # Opens browser, press Enter when done
tap site -b github/notifications     # Use saved cookies
tap site --pause twitter/search query=claude  # One-off CAPTCHA
```

| Flag | Use when |
|---|---|
| `--lp` | JS rendering without auth (fast, no cookies) |
| `-b` | Auth needed (cookies, login state) |
| `--pause` | One-off CAPTCHA/interaction (implies `-b --no-headless`) |

## Session strategy

**Always reuse the `default` session.** Only create named sessions when you need isolation (different accounts, parallel work).

```bash
# Ensure default session exists (skip if already running)
tap browser session list              # Check first
tap browser session new default       # Create only if missing
```

**When to use which:**

| Need | Approach |
|---|---|
| Single script with auth | `tap site -b <script>` |
| Multi-step workflow (navigate, fill, extract) | `tap browser` with `default` session |
| Parallel browser tasks (subagents) | Each subagent creates its own named session |

**Stale session recovery:** If `default` is unresponsive, close and recreate — same name = same profile directory = cookies preserved.

```bash
tap browser session close default
tap browser session new default
```

See [references/browser.md](references/browser.md) for full session/tab/action commands and recovery details.

## `tap browser` — Persistent sessions, tabs, and network

Manage long-lived browser instances across CLI invocations.

```bash
tap browser session new <name>        # Launch Chrome
tap browser tab new <name> --url <url>
tap browser navigate <url>
tap browser evaluate <js>
tap browser screenshot
tap browser forms
tap browser fill <sel> <val> [--submit <sel>]
tap browser click <sel>               # Real mouse click
tap browser type <sel> <text>         # Per-keystroke typing
tap browser hover <sel>               # Trigger :hover / mouseenter
tap browser scroll <sel>              # Scroll element into view
tap browser scroll --x 0 --y 1000    # Scroll to position
tap browser select <sel> <value>      # Pick <select> option
tap browser wait <sel> [--timeout 30s]  # Wait for element visible
tap browser back                      # History back
tap browser forward                   # History forward
tap browser reload                    # Reload page
tap browser session close [name]
```

### Network capture & interception

Capture API responses and intercept requests on tracked tabs via CDP.
See [references/network.md](references/network.md) for all flags, patterns, and examples.

```bash
tap browser network wait --url-pattern "*/api/*" --body   # Capture API response
tap browser network log --resource-type XHR,Fetch         # Stream as NDJSON
tap browser network body <requestId>                      # Fetch body by ID
tap browser network intercept --block --url-pattern "*.ads.*"  # Block requests
tap browser network intercept --url-pattern "*/api/*" --respond '{}' --status 200
tap browser network intercept --url-pattern "*/api/*" --header "Authorization: Bearer tok"
tap browser network clear                                 # Remove rules
```

## Tips

- Reuse the `default` session — don't create new sessions unless you need isolation.
- Prefer `--lp` over `-b` when you just need JS rendering without auth.
- Use `tap login` first for auth-required sites, then `-b`.
- Prefer `tap site` over `tap fetch` when a script exists.
- Use `tap browser` for multi-step workflows needing state across invocations.
- Use `tap browser network wait` to capture clean API JSON instead of scraping DOM.
- Use `tap --local-only site ...` to skip remote cache during script development.
