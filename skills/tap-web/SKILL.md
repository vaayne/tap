---
name: tap-web
description: >
  Access websites, search the web, and extract clean content using the `tap` CLI.
  Also manage persistent browser sessions for long-lived automation workflows.
  Use when the user asks to search the web, read a webpage, fetch article content,
  get trending topics, look up social media posts, check stock prices, search videos,
  retrieve structured data from any supported site, or manage browser sessions/tabs
  for persistent automation. Triggers on: "search for", "what's trending on",
  "fetch this page", "read this URL", "get content from",
  "look up on Twitter/Weibo/Reddit/YouTube/etc.", "check stock", "translate",
  "browser session", "open a tab", "navigate to", "take a screenshot",
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
tap site <script> -f json | jq '.'   # JSON output, pipe to jq
```

## Login & browser backends

```bash
tap login https://github.com/login   # Opens browser, press Enter when done
tap site -b github/notifications     # Use saved cookies

tap site --pause twitter/search query=claude  # One-off CAPTCHA/interaction
```

| Flag | Description |
|---|---|
| `--lp` | Lightpanda: fast headless, no cookies. Prefer for JS rendering. |
| `-b` | Chrome: cookies, login, `--pause`. Use when auth needed. |
| `--no-headless` | Show Chrome window |
| `--pause` | Pause for interaction; implies `--no-headless` and `-b` |
| `-t` | Timeout (e.g., `30s`, `2m`) |

## `tap browser` — Persistent sessions and tabs

For full reference, read `docs/browser.md` in the tap repo.

```bash
# Sessions
tap browser session new <name>              # Launch headless Chrome
tap browser session new <name> --no-headless
tap browser session new <name> --ws-url <url>
tap browser session list
tap browser session close [name]

# Tabs
tap browser tab new <name> [--url <url>]
tap browser tab list
tap browser tab close [name]

# Actions (use --session/--tab to override defaults)
tap browser navigate <url>
tap browser evaluate <javascript>
tap browser screenshot [--output <path>]
tap browser forms
tap browser fill <sel> <val> [--submit <sel>]
```

### Network capture & interception

For full reference, read `docs/network.md` in the tap repo.

```bash
# Passive capture (Network domain)
tap browser network wait --url-pattern "*/api/*" --body --timeout 30s
tap browser network log --resource-type XHR,Fetch
tap browser network body <requestId>

# Active interception (Fetch domain, process stays alive)
tap browser network intercept --block --url-pattern "*.ads.*"
tap browser network intercept --url-pattern "*/api/*" --respond '{"ok":true}' --status 200
tap browser network intercept --url-pattern "*/api/*" --header "Authorization: Bearer token"
tap browser network clear
```

Key flags: `--url-pattern` (glob, `*` matches across `/`), `--method`, `--resource-type`, `--timeout`, `--body`, `--block`, `--respond`, `--status`, `--header`.

## Tips

- Prefer `--lp` over `-b` when you just need JS rendering without auth.
- Use `tap login` first for auth-required sites, then `-b`.
- Prefer `tap site` over `tap fetch` when a script exists.
- Use `tap browser` for multi-step workflows needing state across invocations.
- Use `tap browser network wait` to capture clean API JSON instead of scraping DOM.
