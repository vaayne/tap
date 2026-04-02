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

## `tap browser` — Persistent browser sessions and tabs

Manage long-lived browser instances that survive across CLI invocations. Use for multi-step automation, multi-tab workflows, or interactive debugging.

### Session lifecycle

```bash
tap browser session new <name>                # Launch a managed headless Chrome
tap browser session new <name> --no-headless  # Visible browser
tap browser session new <name> --ws-url <url> # Connect to remote CDP endpoint
tap browser session list                      # List all sessions
tap browser session info [name]               # Show session details and tabs
tap browser session select <name>             # Set default session
tap browser session close [name]              # Stop browser and remove metadata
```

### Tab lifecycle

```bash
tap browser tab new <name> [--url <url>]      # Create a tracked tab
tap browser tab list                          # List tracked tabs
tap browser tab select <name>                 # Set default tab
tap browser tab close [name]                  # Close and remove a tab
```

### Browser actions

```bash
tap browser navigate <url>                    # Navigate the selected tab
tap browser evaluate <javascript>             # Run JS and print the result
tap browser screenshot [--output <path>]      # Capture a full-page PNG
tap browser forms                             # Discover fillable form elements
tap browser fill <sel> <val> [<sel> <val>...] # Fill form fields
tap browser fill <sel> <val> --submit <sel>   # Fill and click submit
```

All action commands accept `--session <name>` and `--tab <name>` to override defaults.

### Resolution rules

When `--session` or `--tab` is omitted, tap resolves automatically:

- **Session**: `--session` flag → selected session → the only session
- **Tab**: `--tab` flag → selected tab → the only live tracked tab

If ambiguous, tap fails with guidance instead of guessing.

### Example: form discovery and filling

```bash
tap browser session new login-demo --no-headless
tap browser tab new page --url https://example.com/login

# Discover fillable elements (inputs, textareas, selects, buttons)
tap browser forms
# Returns JSON: selector, type, name, placeholder, label, value, role per element

# Fill fields using the reported selectors
tap browser fill "#username" "myuser" "#password" "secret"

# Fill and submit in one command
tap browser fill "#email" "me@example.com" --submit "button[type=submit]"

tap browser session close login-demo
```

`tap browser fill` uses React-compatible native value setters with proper `input`/`change` event dispatch — works with React, Vue, Angular, and vanilla HTML forms.

### Network interception

Capture, inspect, and intercept network requests on tracked browser tabs using CDP Network and Fetch domains.

```bash
# Wait for a specific API request to complete
tap browser network wait --url-pattern "*/api/*" --body --timeout 30s

# Stream all network activity as NDJSON
tap browser network log --resource-type XHR,Fetch

# Fetch response body by request ID (from wait/log output)
tap browser network body "1234.56"

# Block ad/tracking requests (process stays alive, Ctrl-C to stop)
tap browser network intercept --block --url-pattern "*.ads.*"

# Mock an API response
tap browser network intercept \
  --url-pattern "*/api/user" \
  --respond '{"name":"test"}' --status 200

# Add auth headers to API requests
tap browser network intercept \
  --url-pattern "*/api/*" \
  --header "Authorization: Bearer tok_abc123"

# Remove all interception rules
tap browser network clear
```

**Network wait flags:**

| Flag | Description | Default |
|---|---|---|
| `--url-pattern` | Glob pattern (`*` matches any chars including `/`) | match all |
| `--method` | HTTP method(s), comma-separated | match all |
| `--resource-type` | Resource type(s): `XHR`, `Fetch`, `Document`, `Script`, etc. | match all |
| `--timeout` | Maximum time to wait | `30s` |
| `--body` | Include response body in output | `false` |
| `--format` | `pretty`, `json`, `raw` | `pretty` |

**Network intercept flags:**

| Flag | Description |
|---|---|
| `--block` | Block matching requests (mutually exclusive with `--respond`) |
| `--respond` | Mock response body (mutually exclusive with `--block`) |
| `--status` | Mock response status code (default `200`) |
| `--content-type` | Mock response Content-Type (default `application/json`) |
| `--header` | Add/override request header (repeatable, `"Key: Value"`) |

**Common use case — capture SPA API responses:**

```bash
tap browser session new scrape
tap browser tab new page --url https://example.com/dashboard

# Wait for the data API call and capture its JSON response
tap browser network wait --url-pattern "*/api/data*" --body --format json

tap browser session close scrape
```

### Example: multi-tab workflow

```bash
tap browser session new research
tap browser tab new docs --url https://go.dev/doc
tap browser tab new api --url https://pkg.go.dev

tap browser tab select docs
tap browser evaluate 'document.title'

tap browser tab select api
tap browser evaluate 'document.title'

tap browser session close research
```

## Tips

- Prefer `--lp` over `-b` when you just need JS rendering without auth.
- For auth-required sites, use `tap login <url>` first, then `-b`. Cookies persist in the Chrome profile.
- Use `--pause` for one-off CAPTCHA solving before a script runs.
- Prefer `tap site` over `tap fetch` when a site script exists for better structured output.
- If QuickJS execution fails, tap automatically falls back to browser mode.
- Use `tap browser` for multi-step workflows that need state across invocations.
- Prefer `tap browser` over `-b` when you need to keep a browser alive between commands.
