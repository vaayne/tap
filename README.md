# 🚰 Tap

Tap into any website from your terminal.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
```

Or with Go:

```bash
go install github.com/vaayne/tap/cmd/tap@latest
```

Upgrade later with:

```bash
tap upgrade
```

Browser features use agent-browser. Install with `tap doctor --install`.

## Quick start

### Structured data from site scripts

```bash
tap site list
tap site hackernews/top
tap site run github/repo repo=vaayne/tap
tap site search weather
```

Scripts auto-sync from [tap.vaayne.com](https://tap.vaayne.com) into the local cache. Local overrides at `~/.config/tap/sites/{site}/{script}.js` take precedence.

### Readable content from any URL

```bash
tap fetch https://example.com/article
tap fetch --json https://example.com/article
tap fetch -b https://example.com/app --wait-selector '.content'
```

### Use a visible browser for auth when needed

```bash
tap attach chrome
tap browser open https://github.com/login --show
tap site -b github/notifications
tap fetch -b https://github.com/notifications
```

### Reuse your existing Chrome

agent-browser manages Chrome automatically. To attach to an existing Chrome with DevTools enabled:

```bash
tap attach chrome
tap attach status
tap browser open https://example.com
tap browser snapshot --interactive
tap browser click '#submit'
tap browser click @e1
tap browser text
```

You can also attach explicitly:

```bash
tap attach chrome --browser-url http://127.0.0.1:9222
tap attach chrome --browser-url http://127.0.0.1:9222
```

### Browser workflow

```bash
tap browser open https://news.ycombinator.com
tap browser open https://github.com --new-tab
tap browser tabs
tap browser switch tab-2
tap browser screenshot --output github.png
tap browser status
```

### Attached Chrome workflow

```bash
tap attach chrome
tap attach status --json
tap browser evaluate 'document.title'
tap browser screenshot
```

If the attached state becomes stale, rerun `tap attach chrome`.

## Command map

```text
tap
├── site        structured extraction from known sites
├── fetch       clean readable content from arbitrary URLs
├── browser     open pages and automate the current browser context
├── attach      connect tap to an existing Chrome browser
├── status      show the active browser context and current tab
├── doctor      dependency and environment checks
├── upgrade     update tap
└── completion  generate shell completion scripts
```

## Shell completion

Tap can generate shell completion scripts for bash, zsh, fish, and pwsh.

```bash
# Bash
source <(tap completion bash)

# Zsh
source <(tap completion zsh)

# Fish
mkdir -p ~/.config/fish/completions
tap completion fish > ~/.config/fish/completions/tap.fish

# PowerShell / pwsh
tap completion pwsh > ~/.config/powershell/tap.ps1
```

Persistent install paths commonly used by package managers and dotfiles:

```bash
# Bash
mkdir -p ~/.local/share/bash-completion/completions
tap completion bash > ~/.local/share/bash-completion/completions/tap

# Zsh
mkdir -p ~/.zfunc
tap completion zsh > ~/.zfunc/_tap
```

## Common browser-backed flags

These show up on the relevant commands instead of only in global help:

| Flag | Description |
| --- | --- |
| `--browser`, `-b` | Force browser execution and reuse the resolved browser context |
| `--show` | Run the browser visibly |
| `--wait` | Wait a fixed duration after navigation |
| `--wait-selector` | Wait for a CSS selector |
| `--wait-js` | Wait for a JS expression |
| `--timeout` | Set execution timeout |
| `--browser-url` | One-shot agent-browser/DevTools connection override |
| `--profile-dir` | One-shot profile override |
| `--lightpanda`, `--lp` | Compatibility alias for fast headless browser mode |

Compatibility aliases still work:
- `--ws-url` -> `--browser-url`
- `--delay` -> `--wait`
- `--no-headless` -> `--show`

## Advanced browser commands

The browser command still includes lower-level tools when needed:

```bash
tap browser evaluate ...
tap browser snapshot
tap browser forms
tap browser cookies ...
tap browser network ...
tap browser set ...
tap browser storage ...
tap browser state ...
tap browser auth ...
tap browser get ...
tap browser vitals
tap browser diff ...
```

## Browser backend

Tap uses agent-browser as the single browser backend. It manages Chrome for full browser automation, auth, screenshots, and network workflows. Install or update the backend with:

```bash
tap doctor --install
```

## Docs

- [docs/browser.md](docs/browser.md) — browser UX and reference
- [docs/network.md](docs/network.md) — network interception reference
- [web/README.md](web/README.md) — web app docs

## Agent skill

Tap ships with a built-in agent skill that gives coding agents web access, site scripts, and browser automation.

```bash
npx skills add vaayne/tap
```

## Go library

```bash
go get github.com/vaayne/tap
```

```go
client, _ := tap.New()
defer client.Close()

result, _ := client.RunScript(ctx, "hackernews/top", nil)
content, _ := client.Fetch(ctx, "https://example.com", &fetch.Options{Markdown: true})
fmt.Println(content.Markdown)
```

## License

MIT
