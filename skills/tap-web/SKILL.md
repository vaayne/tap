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
  "mock API", "network log", "wait for request", "save as PDF", "export PDF",
  or any web access task.
---

# tap-web

**Before accessing any site, check `$XDG_CONFIG_HOME/tap/site-notes/{domain}.md` for saved knowledge. Update after learning.** See [references/site-notes.md](references/site-notes.md).

## Quick reference

```bash
# Content extraction
tap fetch <url>                          # Clean markdown from any URL
tap fetch --json <url>                   # JSON with metadata

# Site scripts — always discover first, never guess
tap site list                            # List all local scripts
tap site search <kw>                     # Search remote catalog
tap site info <script>                   # Show script details
tap site sync                            # Force re-sync from remote
tap site <site/action> [key=value]       # Run a script
tap site <site/action> -f json           # Run with JSON output
tap site -b <script>                     # With auth (browser cookies)

# Login
tap login <url>                          # Opens browser, press Enter when done
```

| Flag | Use when |
|---|---|
| `--lp` | JS rendering without auth (implies `-b`, prefer over `-b` alone when no auth needed) |
| `-b` | Auth needed — run `tap login <url>` first to save cookies |
| `--pause` | One-off CAPTCHA (implies `-b --no-headless`). **Requires interactive terminal — agents use `--delay` or `--wait-selector` instead** |
| `--delay <dur>` | Wait fixed duration after navigation (e.g., `--delay 3s`) |
| `--wait-selector <sel>` | Wait until CSS selector visible before continuing |
| `--wait-js <expr>` | Wait until JS expression truthy before continuing |
| `--timeout <dur>` | Global execution timeout (e.g., `--timeout 30s`, `-t 2m`) |

Prefer `tap site` over `tap fetch` when a script exists — structured data, less tokens.

## Browser sessions

**Always reuse the `default` session.** Only create named sessions for isolation (parallel subagents, different accounts). Stale recovery: close + recreate with same name. Note: closing a session deletes its profile directory (including cookies) — use `tap login` for persistent auth that survives session close.

See [references/browser.md](references/browser.md) for full commands, session strategy, and recovery.

```bash
tap browser session list | new <name> | info [name] | close [name]
tap browser tab new <name> --url <url>
# Page
tap browser navigate <url>              # Go to URL
tap browser back | forward | reload     # History navigation
tap browser text [selector]             # Clean text via defuddle (token-efficient)
tap browser evaluate <js>               # Run JavaScript
tap browser screenshot | pdf            # Capture output
# Interaction (real CDP events, human-like)
tap browser click | hover <sel>         # Mouse actions
tap browser type <sel> <text>           # Per-keystroke typing
tap browser scroll <sel>                # Scroll into view
tap browser select <sel> <value>        # Pick <select> option
tap browser fill <sel> <val> [--submit] # Set form values
tap browser wait <sel> [--timeout 30s]  # Wait for element visible
tap browser keypress Enter              # Keyboard (Enter, Tab, Escape, Ctrl+a...)
tap browser dialog [--accept=false]     # Handle alert/confirm/prompt
# State
tap browser forms                       # Discover form fields
tap browser cookies get | set | clear   # Cookie management (includes httpOnly)
# Network — see references/network.md
tap browser network wait --url-pattern "*/api/*" --body
tap browser network log --resource-type XHR,Fetch
tap browser network intercept --block --url-pattern "*.ads.*"
```

## Reading pages efficiently

**Never use `evaluate "document.documentElement.outerHTML"`.** Pick the cheapest path:

1. **API exists?** → `network wait --body` (~1-5k tokens)
2. **Readable content?** → `text [selector]` (~2-10k tokens, defuddle strips boilerplate)
3. **Form fields?** → `forms` (~0.5-2k tokens)
4. **Specific data?** → `evaluate` with targeted selector (~0.5-5k tokens)
5. **Visual layout?** → `screenshot` (fixed cost)

## References

- [browser.md](references/browser.md) — Sessions, tabs, actions, session strategy
- [network.md](references/network.md) — Network capture & interception
- [script-development.md](references/script-development.md) — Writing site scripts
- [site-notes.md](references/site-notes.md) — Per-site knowledge system
