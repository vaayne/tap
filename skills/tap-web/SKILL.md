---
name: tap-web
description: >
  Access websites, search the web, and extract clean content using the `tap` CLI.
  Also manage persistent browser sessions for long-lived automation workflows,
  and debug or automate Electron apps via CDP.
  Use when the user asks to search the web, read a webpage, fetch article content,
  get trending topics, look up social media posts, check stock prices, search videos,
  retrieve structured data from any supported site, manage browser sessions/tabs,
  capture or intercept network requests, block or mock API calls,
  debug an Electron app, automate an Electron app, connect to a desktop app via CDP.
  Triggers on: "search for", "what's trending on", "fetch this page", "read this URL",
  "get content from", "look up on Twitter/Weibo/Reddit/YouTube/etc.", "check stock",
  "translate", "browser session", "open a tab", "navigate to", "take a screenshot",
  "evaluate javascript", "persistent browser", "fill form", "form fields",
  "discover inputs", "intercept requests", "capture network", "block requests",
  "mock API", "network log", "wait for request", "save as PDF", "export PDF",
  "debug electron", "automate electron", "connect to electron app",
  or any web access task.
---

# tap-web

## Install

```bash
# Install tap CLI
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
# Install this skill
npx skills add vaayne/tap
# Upgrade tap
tap upgrade
```

**Before accessing any site, check `$XDG_CONFIG_HOME/tap/site-notes/{domain}.md` for saved knowledge. Update after learning.** See [references/site-notes.md](references/site-notes.md).

## Pick the right tool

tap has three tiers — always start from the top and escalate only when needed:

| Tier | Tool | What it's for | When to use |
|---|---|---|---|
| **1. Script** | `tap site` | Structured data from known sites | A script exists for the site (check `tap site list` / `search` first). Fastest, least tokens, returns clean JSON. |
| **2. Fetch** | `tap fetch` | Clean article/page content | Reading a URL's main content as Markdown. No interaction needed. Works for blogs, docs, news. |
| **3. Browser** | `tap browser` | Full browser automation | Multi-step workflows, form filling, JS-heavy SPAs, network interception, screenshots. Most powerful but most expensive. |

**Decision flow:**
1. Is there a site script? → `tap site` (always check first — never guess script names)
2. Just need to read a page? → `tap fetch <url>` (add `--lp` if JS-rendered, `-b` if auth needed)
3. Need to interact, navigate multiple pages, or capture network? → `tap browser`

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
| `--lp` | JS rendering without auth (implies `-b`, macOS/Linux only, not all sites work — fall back to `-b` if rendering fails) |
| `-b` | Auth needed — run `tap login <url>` first to save cookies |
| `--pause` | One-off CAPTCHA (implies `-b --no-headless`). **Requires interactive terminal — agents use `--delay` or `--wait-selector` instead** |
| `--delay <dur>` | Wait fixed duration after navigation (e.g., `--delay 3s`, implies `-b`) |
| `--wait-selector <sel>` | Wait until CSS selector visible before continuing (implies `-b`) |
| `--wait-js <expr>` | Wait until JS expression truthy before continuing (implies `-b`) |
| `--timeout <dur>` | Global execution timeout (e.g., `--timeout 30s`, `-t 2m`) |

## Browser sessions

**Always reuse the persisted default browser context.** Resolution order is:
1. explicit `--session`
2. persisted default context (`tap attach ...` or a previously-created `default` session)
3. managed local `default` session auto-created only when no persisted context exists

This means `tap attach chrome` becomes sticky across later `tap browser ...` commands. If the attached context becomes unreachable, tap marks it stale and fails explicitly — it does **not** silently fall back to a managed local browser.

Only create named sessions for isolation (parallel subagents, different accounts). Stale recovery: re-attach or close + recreate the same session. Note: closing a managed local session deletes its profile directory (including cookies) — use `tap login` for persistent auth that survives session close.

See [references/browser.md](references/browser.md) for full commands, session strategy, status output, and recovery.

```bash
tap status [--json]                     # Show resolved default context + current tab
tap attach status [--json]              # Show attachment/default-context state
tap browser status [--json]             # Show current browser context + tab state
tap browser session list [--json]
tap browser session new <name>
tap browser session info [name] [--json]
tap browser session close [name]
tap browser tab new <name> --url <url>
# Page
tap browser open <url>                  # Resolve current context, create/navigate tab
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

## Electron apps

Connect tap to any Electron (or CEF-based) desktop app via CDP — no plugins needed.

```bash
# Launch app with debugging + auto-create session
tap electron launch "/Applications/MyApp.app/Contents/MacOS/MyApp" [--session <name>]

# Attach to already-running app (started with --remote-debugging-port=PORT)
tap electron attach --port <port> [--session <name>]

# Adopt open windows as tracked tabs (run after attach or launch)
tap electron discover [--session <name>]

# List running processes with debug ports
tap electron ps
```

After `discover`, use `tap browser` commands with `--session <name>` normally — screenshot, evaluate, click, type, network intercept, etc.

See [references/browser.md](references/browser.md) for full browser command reference.

## Reading pages efficiently

**Never use `evaluate "document.documentElement.outerHTML"`.** Pick the cheapest path:

1. **API exists?** → `network wait --body` (~1-5k tokens)
2. **Readable content?** → `text [selector]` (~2-10k tokens, defuddle strips boilerplate)
3. **Form fields?** → `forms` (~0.5-2k tokens)
4. **Specific data?** → `evaluate` with targeted selector (~0.5-5k tokens)
5. **Visual layout?** → `screenshot` (fixed cost)

## Doctor — browser dependencies

```bash
tap doctor              # Check tap, Chrome, Lightpanda status
tap doctor --install    # Install or update Lightpanda to latest nightly
```

**Lightpanda limitations:** macOS/Linux only. Not all sites work (advanced Web APIs, heavy JS frameworks may fail). No network interception, no persistent sessions/cookies. Use Chrome (`-b`) when `--lp` doesn't work.

## References

- [browser.md](references/browser.md) — Sessions, tabs, actions, session strategy
- [network.md](references/network.md) — Network capture & interception
- [script-development.md](references/script-development.md) — Writing site scripts
- [site-notes.md](references/site-notes.md) — Per-site knowledge system
