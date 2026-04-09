# Current Command UX

Tap now exposes a smaller public CLI centered on the common jobs:

```text
tap
├── site
├── fetch
├── browser
├── attach
├── status
├── doctor
└── upgrade
```

## Common flows

### Read structured site data

```bash
tap site list
tap site hackernews/top
tap site run github/repo repo=vaayne/tap
```

### Read clean content from a URL

```bash
tap fetch https://example.com/article
tap fetch -b https://example.com/app --wait-selector '.content'
```

### Reuse an existing browser

```bash
tap attach chrome
tap browser open https://example.com
tap browser click '#submit'
```

### Use a visible browser for auth

```bash
tap attach chrome
tap browser open https://github.com/login --show
tap site -b github/notifications
```

### Multi-tab browser flow

```bash
tap browser open https://news.ycombinator.com
tap browser open https://github.com --new-tab
tap browser tabs
tap browser switch tab-2
```

## Browser-related flags

Use these on `site`, `fetch`, and relevant `browser` commands:

- `--browser`, `-b`
- `--show`
- `--wait`
- `--wait-selector`
- `--wait-js`
- `--timeout`
- `--browser-url`
- `--profile-dir`
- `--lightpanda`
