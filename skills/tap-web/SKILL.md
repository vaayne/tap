---
name: tap-web
description: >
  Access websites, search the web, and extract clean content using the `tap` CLI.
  Use when the user asks to search the web, read a webpage, fetch article content,
  get trending topics, look up social media posts, check stock prices, search videos,
  or retrieve structured data from any supported site. Triggers on: "search for",
  "what's trending on", "fetch this page", "read this URL", "get content from",
  "look up on Twitter/Weibo/Reddit/YouTube/etc.", "check stock", "translate",
  or any web access task.
---

# tap-web

Use the `tap` CLI to access websites from the terminal. Two main commands:

## `tap fetch` — Extract clean content from any URL

```bash
tap fetch <url>              # Clean markdown output
tap fetch --json <url>       # JSON with full metadata
```

Use for reading articles, documentation, blog posts, or any webpage.

## `tap site` — Run site-specific scripts for structured data

**Always discover available scripts first. Do not guess script names or parameters.**

### Step 1: Find the right script

```bash
tap site list                # List all available scripts
tap site search <keyword>    # Search scripts by keyword (e.g., "twitter", "news", "stock")
```

### Step 2: Check script details

```bash
tap site info <script>       # Show description, parameters (* = required), and usage
```

### Step 3: Run the script

Follow the exact usage shown by `tap site info`. Use `key=value` syntax for all parameters.

```bash
tap site <site/action> [key=value ...]
```

### Output format

```bash
tap site <script> -f json     # Raw JSON (default for piping)
tap site <script> -f pretty   # Human-readable (default for terminal)
tap site <script> -f raw      # Unformatted raw output
```

Pipe to `jq` for filtering:

```bash
tap site <script> -f json | jq '.'
```

## Login & interactive browser

Some sites require login or CAPTCHA solving before scripts work.

### `tap login` — Log in once, run scripts many times

```bash
tap login https://github.com/login          # Opens visible browser, press Enter when done
tap site -b github/notifications             # Subsequent runs use saved cookies

tap login https://twitter.com/i/flow/login
tap site twitter/search query="AI agents"
```

Cookies persist in the Chrome profile directory (`~/.cache/tap/chrome-profile-$USER`).
Use `--profile-dir` to manage separate profiles (e.g., work vs personal).

### `--pause` — One-off interaction before script execution

```bash
tap site --pause twitter/search query=claude
# Browser opens → solve CAPTCHA / interact → press Enter → script runs → browser closes
```

`--pause` implies `--no-headless` (visible browser) and `-b` (browser mode).

## Browser backends

When HTTP isn't enough (JS-rendered pages), use a browser backend:

- **`--lp`** — Lightpanda: fast headless browser, no cookies/auth. Prefer this when you just need JS rendering.
- **`-b`** — Chrome: full browser with cookies, login state, `--pause` support. Use when auth is needed, **or when Lightpanda has compatibility issues with a site**.

```bash
tap fetch https://spa-site.com --lp          # Fast JS rendering
tap site hackernews/top --lp                 # Structured data via Lightpanda
tap site -b github/notifications             # Needs saved cookies → Chrome
```

## Global options

| Flag | Description |
|---|---|
| `--lightpanda, --lp` | Use Lightpanda headless browser (fast, no cookies) |
| `--browser, -b` | Force Chrome browser (cookies, login, interactive) |
| `--no-headless` | Show Chrome window (debug auth issues) |
| `--pause` | Pause for manual interaction; implies `--no-headless` and `-b` |
| `--timeout, -t` | Execution timeout (e.g., `30s`, `2m`) |
| `--quiet, -q` | Suppress log output |
| `--verbose` | Enable verbose logging |

## Tips

- Prefer `--lp` over `-b` when you just need JS rendering without auth.
- For auth-required sites, use `tap login <url>` first, then `-b`. Cookies persist in the Chrome profile.
- Use `--pause` for one-off CAPTCHA solving before a script runs.
- Prefer `tap site` over `tap fetch` when a site script exists for better structured output.
- If QuickJS execution fails, tap automatically falls back to browser mode.
