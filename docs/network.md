# Network Interception

Tap can capture, inspect, and intercept network requests on the current browser tab using Chrome DevTools Protocol (CDP). This covers two use cases:

- **Passive capture** (Network domain): observe requests, wait for specific API responses, stream network logs
- **Active interception** (Fetch domain): block, mock, or modify requests

All commands operate on the resolved current tab. Advanced `--session` and `--tab` overrides still exist but are hidden from the common UX.

## Quick Start

```bash
# Open a page in the default browser context
tap browser open https://example.com

# Wait for a specific API request to complete
tap browser network wait --url-pattern "*/api/*"

# Stream all network activity as NDJSON
tap browser network log

# Block ad requests
tap browser network intercept --block --url-pattern "*.ads.*"

# Clear interception rules
tap browser network clear
```

## Commands

### `tap browser network wait`

Block until a network request matching filters completes, then print the captured entry.

```bash
tap browser network wait [flags]
```

| Flag | Description | Default |
|---|---|---|
| `--url-pattern` | Glob pattern (`*` matches any chars including `/`) | match all |
| `--method` | HTTP method(s), comma-separated | match all |
| `--resource-type` | Resource type(s), comma-separated (e.g. `XHR,Fetch,Document`) | match all |
| `--timeout` | Maximum time to wait | `30s` |
| `--body` | Include the response body in output | `false` |
| `--format` | Output format: `pretty`, `json`, `raw` | `pretty` |

```bash
# Wait for any XHR/Fetch request
tap browser network wait --resource-type XHR,Fetch

# Wait for a specific API endpoint and capture the response body
tap browser network wait --url-pattern "*/api/search*" --body --timeout 10s

# Wait for a POST request
tap browser network wait --url-pattern "*/graphql" --method POST --body --format json
```

The output is a `NetworkEntry` object:

```json
{
  "requestId": "1234.56",
  "url": "https://example.com/api/search?q=hello",
  "method": "GET",
  "status": 200,
  "mimeType": "application/json",
  "resourceType": "XHR",
  "reqHeaders": { "Accept": "application/json" },
  "resHeaders": { "Content-Type": "application/json" },
  "body": "{\"results\": [...]}",
  "timestamp": "2026-04-02T12:00:00Z"
}
```

If a request fails before receiving a response (DNS failure, connection refused), the entry has `status: 0` and a populated `error` field.

### `tap browser network body`

Fetch the response body for a completed request by its request ID.

```bash
tap browser network body <requestId> [flags]
```

| Flag | Description | Default |
|---|---|---|
| `--format` | `raw` (binary to stdout) or `json` (base64-encoded) | `raw` |

```bash
# Get the request ID from a previous wait/log
tap browser network wait --url-pattern "*/api/data" --format json | jq -r .requestId
# => 1234.56

# Fetch the body
tap browser network body "1234.56"
tap browser network body "1234.56" --format json
```

### `tap browser network log`

Stream completed network requests as newline-delimited JSON (NDJSON).

```bash
tap browser network log [flags]
```

| Flag | Description | Default |
|---|---|---|
| `--url-pattern` | Glob pattern to filter URLs | match all |
| `--method` | HTTP method(s), comma-separated | match all |
| `--resource-type` | Resource type(s), comma-separated | match all |
| `--timeout` | Duration to capture (0 = until interrupted) | `0` |

```bash
# Stream all network requests
tap browser network log

# Only XHR/Fetch requests for 30 seconds
tap browser network log --resource-type XHR,Fetch --timeout 30s

# Filter to API calls, pipe through jq
tap browser network log --url-pattern "*/api/*" | jq '.url, .status'
```

Each line is a JSON `NetworkEntry` object. The stream runs until `--timeout` expires or the process is interrupted (Ctrl-C).

### `tap browser network intercept`

Set Fetch domain interception rules. Rules are **replace-all**: each invocation replaces any previously set rules.

```bash
tap browser network intercept [flags]
```

| Flag | Description | Default |
|---|---|---|
| `--url-pattern` | Glob pattern to match request URLs | match all |
| `--method` | HTTP method(s), comma-separated | match all |
| `--resource-type` | Resource type(s), comma-separated | match all |
| `--block` | Block matching requests | `false` |
| `--header` | Add/override a request header (repeatable, `"Key: Value"`) | — |
| `--respond` | Mock response body | — |
| `--status` | Mock response HTTP status code | `200` |
| `--content-type` | Mock response Content-Type | `application/json` |

The `--block` and `--respond` flags are mutually exclusive.

The process stays alive while interception is active. Press Ctrl-C to stop.

```bash
# Block ad/tracking requests
tap browser network intercept --block --url-pattern "*.doubleclick.net*"
tap browser network intercept --block --url-pattern "*.google-analytics.com*"

# Mock an API endpoint
tap browser network intercept \
  --url-pattern "*/api/user/profile" \
  --respond '{"name":"test","plan":"premium"}' \
  --status 200

# Add an auth header to all API requests
tap browser network intercept \
  --url-pattern "*/api/*" \
  --header "Authorization: Bearer tok_abc123"

# Combine header injection with content-type
tap browser network intercept \
  --url-pattern "*/api/*" \
  --header "Authorization: Bearer tok_abc123" \
  --header "X-Custom: value"
```

### `tap browser network clear`

Disable the Fetch domain and remove all interception rules.

```bash
tap browser network clear
```

This only clears active interception rules. Passive capture (`wait`, `log`) is unaffected.

## URL Pattern Syntax

URL patterns use glob-style matching where `*` matches **any characters including `/`**. This differs from `filepath.Match` which treats `/` as a boundary.

| Pattern | Matches | Doesn't Match |
|---|---|---|
| `*/api/*` | `https://example.com/api/v1/users` | `https://example.com/other/v1` |
| `*.ads.*` | `https://tracker.ads.example.com/pixel` | `https://tracker-ads-example.com/pixel` |
| `*/search?q=*` | `https://example.com/search?q=hello` | `https://example.com/search` |
| `*` | anything | — |

An empty pattern (or omitted `--url-pattern`) matches all URLs.

## Resource Types

CDP resource types correspond to how Chrome classifies each request:

| Type | Description |
|---|---|
| `Document` | HTML pages |
| `Stylesheet` | CSS files |
| `Image` | Images |
| `Media` | Audio/video |
| `Font` | Web fonts |
| `Script` | JavaScript files |
| `XHR` | XMLHttpRequest |
| `Fetch` | Fetch API requests |
| `WebSocket` | WebSocket connections |
| `Other` | Everything else |

Use `--resource-type XHR,Fetch` to capture only API calls.

## Use Cases

### Capture API responses behind SPAs

Many modern sites (Twitter/X, Instagram, Reddit) load data via internal API calls. The JSON from these APIs is far cleaner than scraping the DOM.

```bash
tap browser open https://x.com/elonmusk --show

# Wait for the tweets API endpoint
tap browser network wait --url-pattern "*/UserTweets*" --body --format json
```

### Debug page network activity

Record what a page loads to find hidden API endpoints or audit third-party requests.

```bash
tap browser network log --timeout 30s | jq '{url: .url, status: .status, type: .resourceType}'
```

### Wait for API response after a user action

Fill a search form and capture the resulting API call.

```bash
# In one terminal: start waiting for the search API
tap browser network wait --url-pattern "*/api/search*" --body

# In another terminal: trigger the search
tap browser fill "#search" "golang" --submit "#submit-btn"
```

### Speed up page loads by blocking resources

Block ads, trackers, and heavy assets.

```bash
tap browser network intercept --block --url-pattern "*.doubleclick.net*"

# Navigate — the blocked requests will fail immediately
tap browser navigate https://news-site.com/article
```

### Inject authentication headers

Add auth tokens to API requests without modifying the page.

```bash
tap browser network intercept \
  --url-pattern "*/api/*" \
  --header "Authorization: Bearer tok_abc123"
```

### Mock API responses for testing

Return custom responses to test frontend behavior.

```bash
tap browser network intercept \
  --url-pattern "*/api/feature-flags" \
  --respond '{"dark_mode": true, "beta_features": true}' \
  --status 200

tap browser navigate https://app.example.com
tap browser screenshot --output with-feature-flags.png
```

## Limitations

- **Chrome only.** Network interception requires Chrome/Chromium. Lightpanda does not support the Network or Fetch CDP domains. Run `tap doctor` to verify Chrome is installed.
- **Ephemeral.** Network logs are not persisted. They exist only for the lifetime of the `wait`/`log` process.
- **`body` fetch may fail** for redirected or cached requests. This is a CDP limitation.
- **`intercept` is process-bound.** The interception goroutine runs in the `tap browser network intercept` process. When the process exits, interception stops.
- **`log` channel buffer.** The log stream buffers up to 256 entries. If the consumer can't keep up, newer entries are dropped silently.
- **No WebSocket frame capture.** Only HTTP/HTTPS requests are captured.

## Go Library

The network primitives are available as Go functions in the `browser` package:

```go
import "github.com/vaayne/tap/browser"

// Wait for a matching request (single-shot)
entry, err := browser.WaitForRequest(ctx, debugURL, targetID, browser.NetworkFilter{
    URLPattern: "*/api/*",
    Methods:    []string{"GET"},
}, true) // includeBody

// Fetch response body by request ID
body, err := browser.GetResponseBody(ctx, debugURL, targetID, "1234.56")

// Stream network activity
ch, cancel, err := browser.EnableNetworkLog(ctx, debugURL, targetID, browser.NetworkFilter{})
defer cancel()
for entry := range ch {
    fmt.Println(entry.URL, entry.Status)
}

// Set interception rules
cancel, err := browser.SetInterceptRules(ctx, debugURL, targetID, []browser.InterceptRule{
    {Filter: browser.NetworkFilter{URLPattern: "*.ads.*"}, Block: true},
})
defer cancel()

// Clear interception
err = browser.ClearIntercept(ctx, debugURL, targetID)
```

Or use the `Manager` for session/tab resolution:

```go
mgr := browser.NewManager(store)

entry, err := mgr.NetworkWait(ctx, "session", "tab", filter, true)
body, err := mgr.NetworkGetBody(ctx, "session", "tab", "1234.56")
ch, cancel, err := mgr.NetworkLog(ctx, "session", "tab", filter)
err = mgr.NetworkIntercept(ctx, "session", "tab", rules)
err = mgr.NetworkClearIntercept(ctx, "session", "tab")
```
