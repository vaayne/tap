# Browser Automation

Tap now exposes a task-first browser UX for common workflows while keeping the older session/tab/proxy commands for advanced use.

## Common workflows

### Reuse your existing Chrome

Chrome must already expose a DevTools endpoint.

```bash
tap attach chrome
tap attach status
tap fetch -b https://example.com/private
tap site -b github/notifications
tap browser open https://example.com
tap browser snapshot --interactive
tap browser click '#submit'
tap browser click @e1
```

Supported attach inputs:

```bash
tap attach chrome
tap attach chrome --browser-url http://127.0.0.1:9222
tap attach chrome --port-file ~/Library/Application\ Support/Google/Chrome/DevToolsActivePort
```

If the attached browser becomes unreachable, tap marks the default context as **stale** and fails explicitly. It does **not** silently switch to another browser context.

### Use a visible browser for auth when needed

```bash
tap attach chrome
tap browser open https://github.com/login --show
tap site -b github/notifications
tap fetch -b https://github.com/notifications
tap browser open https://github.com
```

Auth state lives in the resolved browser context:
- attached Chrome/Electron context, if one is active
- otherwise tap's managed default browser profile

### Browser open / tabs / switch

```bash
tap browser open https://news.ycombinator.com
tap browser open https://github.com --new-tab
tap browser tabs
tap browser switch tab-2
tap browser screenshot --output github.png
tap browser status
```

Default behavior:
- `tap browser open <url>` navigates the current tab
- `--new-tab` creates another tracked tab
- `tap browser tabs` shows stable tab IDs like `tab-1`, `tab-2`
- `tap browser switch <tab-id>` switches by exact ID
- `tap browser close-tab` closes the current tab

### Attach to Electron

```bash
tap attach electron --port 9333
tap browser tabs --json
tap browser evaluate 'document.title'
tap browser screenshot
```

Or launch and attach in one step:

```bash
tap attach electron --launch /Applications/MyApp.app/Contents/MacOS/MyApp
```

`--port` must be a **browser CDP port**, not a generic Node inspector port.

## Status commands

These commands provide machine-readable output for automation:

```bash
tap status --json
tap attach status --json
tap browser status --json
tap browser tabs --json
```

## Command map

### Task-first commands

```bash
tap site list
tap site search <query>
tap site info <script>
tap site sync
tap site run <script> [key=value ...]
tap site <script> [key=value ...]     # compatibility shorthand

tap fetch <url>

tap attach chrome
tap attach electron --port <port>
tap attach status
tap attach clear

tap browser open <url>
tap browser tabs
tap browser switch <tab-id>
tap browser close-tab [tab-id]
tap browser status
```

### Additional browser tools

```bash
tap browser evaluate ...
tap browser snapshot
tap browser forms
tap browser cookies ...
tap browser network ...
```

Use these when you need lower-level page inspection, cookie management, or network tooling.

### Query element and page properties (`get` / `is`)

All selectors accept a CSS selector or a snapshot ref (`@eN`).

```bash
tap browser get text "h1"
tap browser get html "#content"
tap browser get value "input[name=q]"
tap browser get attr "a.logo" href
tap browser get title
tap browser get url
tap browser get count "li.item"          # CSS selector only — refs address one element
tap browser get box "#sidebar"           # prints JSON {x,y,width,height}
tap browser get styles "button.primary"  # prints JSON computed-style map

tap browser is visible "#modal"
tap browser is enabled "button[type=submit]"
tap browser is checked "input[type=checkbox]"
tap browser is visible @e3              # ref form
```

`is` commands always exit 0 and print `true` or `false`.

### Extra interaction commands

```bash
tap browser dblclick "td.editable"       # real double-click (clickCount=2)
tap browser focus "input[name=email]"    # move keyboard focus
tap browser check "input[name=agree]"    # idempotent — clicks only if unchecked
tap browser uncheck "input[name=newsletter]"
tap browser scrollintoview "#footer"     # scrollIntoViewIfNeeded
tap browser upload "input[type=file]" /tmp/report.pdf
tap browser drag ".card" ".dropzone"     # move→press→interpolate→release
```

Low-level mouse events:

```bash
tap browser mouse move 640 400           # absolute coordinates
tap browser mouse down [left|right|middle]
tap browser mouse up   [left|right|middle]
tap browser mouse wheel 300              # dy (positive = down); optional dx
```

Low-level keyboard events:

```bash
tap browser keyboard type "hello"        # per-character keyDown/char/keyUp
tap browser keyboard insert "paste me"  # Input.insertText, no key events
tap browser keydown Shift               # hold key
tap browser keyup   Shift               # release key
```

### Wait modes

`tap browser wait` accepts exactly one condition per invocation:

```bash
tap browser wait "#login-form"               # element visible (default)
tap browser wait ".spinner" --state hidden   # state: visible|hidden|attached|detached
tap browser wait 2000                        # plain ms
tap browser wait 1.5s                        # Go duration
tap browser wait --text "Welcome back"       # body text contains substring
tap browser wait --url "**/dashboard"        # location.href glob (* and **)
tap browser wait --load networkidle          # load | domcontentloaded | networkidle
tap browser wait --fn "window.__ready"       # poll until JS expression is truthy
```

Default `--timeout` is 30 s. `tap browser open` now honours `--wait-selector` by calling `tap browser wait` internally before returning.

### Semantic locator (`find`)

`find` lets you target elements without knowing CSS selectors.

Locator kinds: `role`, `text`, `label`, `placeholder`, `alt`, `title`, `testid`, `first`, `last`, `nth`.

Actions: `click`, `fill <value>`, `type <value>`, `hover`, `focus`, `check`, `uncheck`, `text`.

`fill` and `type` require a `<value>` argument. `text` prints the element's trimmed textContent.

```bash
tap browser find role button click --name "Submit"    # --name filters by accessible name
tap browser find text "Sign in" click
tap browser find text "Submit" click --exact          # --exact requires exact match
tap browser find label "Email" fill "me@example.com"
tap browser find placeholder "Search…" type "golang"
tap browser find alt "company logo" click
tap browser find title "Close dialog" click
tap browser find testid "login-btn" click
tap browser find first "li.item" text
tap browser find last "tr" text
tap browser find nth 2 "li.item" click                # 1-based index
```

### Web storage (`storage`)

```bash
tap browser storage local                  # dump all localStorage as JSON
tap browser storage local myKey            # print one value
tap browser storage local set myKey value  # set an entry
tap browser storage local clear            # clear all entries

tap browser storage session [...]          # same four forms for sessionStorage
```

### Auth state (`state`)

Save and load browser auth state in Playwright `storageState` format (cookies + localStorage).

```bash
tap browser state save auth.json
tap browser state load auth.json
```

**Limitations:**
- `save`: captures localStorage only for the **current tab's origin**; all cookies from the entire browser context are included.
- `load`: cookies are applied globally; localStorage is restored only for origins matching the current page — other origins in the file are skipped with a warning.
- The state file is written with `0600` permissions.

### Emulation overrides (`set`)

Settings are persisted per tab and **re-applied automatically on every invocation** — no need to repeat them.

```bash
tap browser set viewport 1280 720          # width height [scale]
tap browser set device "iPhone 14"         # sets viewport + UA + touch in one shot
tap browser set geo 37.7749 -122.4194      # lat lng
tap browser set offline on                 # on|off
tap browser set headers '{"Authorization":"Bearer token123"}'  # JSON object
tap browser set media dark                 # dark|light (prefers-color-scheme)
tap browser set useragent "MyBot/1.0"
tap browser set clear                      # remove all persisted overrides for the tab
```

## Browser-backed flags

Common commands expose browser-related flags directly:

```bash
--browser, -b          Use browser execution and reuse the resolved context
--show                 Run visibly
--wait <duration>      Wait after navigation
--wait-selector <css>  Wait for an element
--wait-js <expr>       Wait for a JS expression
--timeout <duration>   Set execution timeout
--browser-url <url>    One-shot DevTools override
--profile-dir <path>   One-shot profile override
```

Compatibility aliases remain accepted:
- `--ws-url` -> `--browser-url`
- `--delay` -> `--wait`
- `--no-headless` -> `--show`
