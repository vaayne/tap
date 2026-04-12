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
