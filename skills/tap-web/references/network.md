# Network Capture & Interception

Capture, inspect, and intercept network requests on tracked browser tabs via CDP.

- **Network domain** (passive): `wait`, `log`, `body`
- **Fetch domain** (active): `intercept`, `clear`

## Commands

### wait — Block until a matching request completes

```bash
tap browser network wait --url-pattern "*/api/*" --body --timeout 30s
tap browser network wait --resource-type XHR,Fetch --method POST
tap browser network wait --url-pattern "*/graphql" --body --format json
```

| Flag | Default | Description |
|---|---|---|
| `--url-pattern` | all | Glob (`*` matches any chars including `/`) |
| `--method` | all | Comma-separated (e.g. `GET,POST`) |
| `--resource-type` | all | Comma-separated (e.g. `XHR,Fetch,Document`) |
| `--timeout` | `30s` | Max wait time |
| `--body` | false | Include response body |
| `--format` | `pretty` | `pretty`, `json`, `raw` |

Output is a JSON entry with `requestId`, `url`, `method`, `status`, `mimeType`, `resourceType`, `reqHeaders`, `resHeaders`, `body`, `error`, `timestamp`.

### body — Fetch response body by request ID

```bash
tap browser network body <requestId>              # Raw to stdout
tap browser network body <requestId> --format json # Base64-encoded
```

### log — Stream requests as NDJSON

```bash
tap browser network log                                      # All requests
tap browser network log --url-pattern "*/api/*" --timeout 30s
tap browser network log --resource-type XHR,Fetch | jq '.url'
```

Runs until `--timeout` or Ctrl-C. Same filter flags as `wait` (except `--body`).

**Always pass `--timeout` — default is infinite.** Without it, `log` runs forever and agents will hang.

### intercept — Block, mock, or modify requests

Process stays alive while active. Ctrl-C to stop. Rules are **replace-all** per call.

```bash
# Block requests
tap browser network intercept --block --url-pattern "*.ads.*"

# Mock a response
tap browser network intercept \
  --url-pattern "*/api/user" \
  --respond '{"name":"test"}' --status 200 --content-type "application/json"

# Inject headers
tap browser network intercept \
  --url-pattern "*/api/*" \
  --header "Authorization: Bearer tok_abc"
```

| Flag | Description |
|---|---|
| `--block` | Fail matching requests (exclusive with `--respond`) |
| `--respond` | Mock response body (exclusive with `--block`) |
| `--status` | Mock status code (default `200`) |
| `--content-type` | Mock Content-Type (default `application/json`) |
| `--header` | Add request header, repeatable (`"Key: Value"`) |

### clear — Remove interception rules

```bash
tap browser network clear
```

Only clears Fetch domain rules. Passive capture (`wait`/`log`) is unaffected.

## URL Pattern Syntax

`*` matches **any characters including `/`** (unlike filepath globs).

| Pattern | Matches |
|---|---|
| `*/api/*` | `https://example.com/api/v1/users` |
| `*.ads.*` | `https://tracker.ads.example.com/pixel` |
| `*/search?q=*` | `https://example.com/search?q=hello` |

## Resource Types

`Document`, `Stylesheet`, `Image`, `Media`, `Font`, `Script`, `XHR`, `Fetch`, `WebSocket`, `Other`

Use `--resource-type XHR,Fetch` to capture only API calls.

## Examples

```bash
# Capture SPA API response instead of scraping DOM
tap browser session new scrape
tap browser tab new page --url https://example.com/dashboard
tap browser network wait --url-pattern "*/api/data*" --body --format json
tap browser session close scrape

# Wait for API after form submit
tap browser network wait --url-pattern "*/api/search*" --body &
tap browser fill "#search" "golang" --submit "#submit-btn"

# Block ads while browsing
tap browser network intercept --block --url-pattern "*.doubleclick.net*"
tap browser navigate https://news-site.com/article
```

## Limitations

- Chrome only (not Lightpanda)
- Ephemeral — logs not persisted
- `body` fetch may fail for redirected/cached requests
- `intercept` stops when the process exits
- `log` drops entries if buffer (256) fills
