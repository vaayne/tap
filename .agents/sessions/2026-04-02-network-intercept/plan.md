# Plan: CDP Network Interception

## Overview

Add CDP Network/Fetch domain support to `tap browser` so users can passively capture network requests, wait for specific API responses, fetch response bodies, and actively intercept (block/mock/modify) requests.

### Goals

- Capture network activity on tracked browser tabs (passive, via `Network` domain)
- Wait for a specific request matching a URL pattern and return its response (single-shot)
- Fetch response body by request ID
- Actively intercept requests to block, mock, or modify headers (via `Fetch` domain)
- Expose all features via CLI under `tap browser network`

### Success Criteria

- [ ] `tap browser network wait --url-pattern "*/api/*" --timeout 30s` blocks until a matching request completes and prints the response (including `requestId`)
- [ ] `tap browser network log` streams captured requests as NDJSON (each entry includes `requestId`)
- [ ] `tap browser network body <requestId>` fetches and prints a response body
- [ ] `tap browser network intercept --block "*.ads.*"` blocks matching requests
- [ ] `tap browser network intercept --url-pattern "*/api/*" --respond '{"mock":true}' --status 200 --content-type "application/json"` mocks responses
- [ ] `tap browser network clear` removes all Fetch domain intercept rules
- [ ] All new code has tests (including Fetch domain integration test)
- [ ] Existing tests continue to pass

### Out of Scope

- HAR export format (entries are JSON, not HAR)
- WebSocket frame capture
- Persisting network logs to the metadata store (ephemeral only)
- Network capture during `tap site` / `tap fetch` (transport layer) — only `tap browser` tabs

## Technical Approach

Two CDP domains serve different purposes:

- **`Network` domain** (passive): observe requests/responses without interference. Used for `wait`, `log`, `body`.
- **`Fetch` domain** (active): intercept and modify/block/mock requests. Used for `intercept`, `clear`.

Both require a long-lived CDP session that stays attached to the target while events are being listened to.

### `withTargetListen` — Long-lived CDP Helper

The existing `withTarget` helper detaches immediately after running actions. Network capture needs a long-lived session. New helper contract:

```go
// withTargetListen connects to a target and returns contexts for long-lived
// event listening. The caller is responsible for calling cancel when done.
//
// Usage pattern:
//   1. Call withTargetListen to get taskCtx
//   2. Register event listeners with chromedp.ListenTarget(taskCtx, ...)
//   3. Enable the domain with chromedp.Run(taskCtx, ...)
//   4. Wait/process events
//   5. Call cancel — this clears TargetID before closing (detach-without-close)
func withTargetListen(ctx context.Context, debugURL string, targetID string) (
    taskCtx context.Context, cancel func(), err error,
)
```

The cancel func uses the same detach-without-close trick as `withTarget`: clears `c.Target.TargetID` before calling chromedp cancel so the tab survives.

**Event listener ordering**: `chromedp.ListenTarget` MUST be called before `chromedp.Run(taskCtx, network.Enable())` to avoid missing initial events. The CDP primitives enforce this ordering internally.

### URL Matching

Use `strings.Contains` for substring matching and `filepath.Match` for glob patterns. Since `filepath.Match` doesn't support `**` or bare `*` matching `/`, we'll treat `*` in URL patterns as "match any characters including `/`" by implementing a simple `matchURL(pattern, url string) bool` helper that replaces `*` with a regex `.*`. This aligns with CDP's `urlPattern` semantics.

### Lightpanda Gating

Browser-type detection is done at the Manager level. The Manager already has access to session metadata which includes the browser mode. We'll check `session.Mode` and return a clear error like `"network interception requires Chrome (not supported by Lightpanda)"` from Manager methods. No new detection helper needed — just a check in each Manager method.

### Components

- **`browser/network.go`** — CDP-level primitives: types, `WaitForRequest`, `GetResponseBody`, `EnableNetworkLog`, `SetInterceptRules`, `ClearIntercept`
- **`browser/cdp.go`** — New `withTargetListen` helper for long-lived CDP sessions
- **`browser/manager.go`** — Manager methods wrapping the primitives with session/tab resolution
- **`cmd/tap/browser_network.go`** — CLI subcommands under `tap browser network`
- **`cmd/tap/browser.go`** — Register the new `network` subcommand group
- **`browser/network_test.go`** — Unit tests for matching/filtering logic and integration tests

## Implementation Phases

### Phase 1: CDP Primitives — Wait & Body

Core building blocks: types, the long-lived target helper, `WaitForRequest`, and `GetResponseBody`.

1. Add `withTargetListen` to `browser/cdp.go` — returns `(context.Context, func(), error)`, uses detach-without-close in cancel func (files: `browser/cdp.go`)
2. Create `browser/network.go` with types and `matchURL` helper (files: `browser/network.go`):
   - `NetworkEntry` — includes `RequestID`, `URL`, `Method`, `Status`, `MIMEType`, `ReqHeaders`, `ResHeaders`, `Body`, `BodyBase64`, `Error`, `Timestamp`
   - `NetworkFilter` — `URLPattern`, `Methods`, `ResourceTypes`
   - `matchURL(pattern, url string) bool` — `*` matches any chars including `/`
   - `matchesFilter(entry, filter) bool` — combines URL, method, resource type matching
3. Implement `WaitForRequest` — registers `chromedp.ListenTarget` FIRST, then enables `Network` domain, listens for `ResponseReceived` + `LoadingFinished`/`LoadingFailed`, returns first matching entry. Respects context cancellation (caller sets timeout via context). Optionally fetches response body based on `includeBody` param. (files: `browser/network.go`)
4. Implement `GetResponseBody` — calls `network.GetResponseBody` on a target via `withTarget` (files: `browser/network.go`)
5. Add unit tests for `matchURL` and `matchesFilter` (files: `browser/network_test.go`)

### Phase 2: Manager Methods — Wait & Body

Wire CDP primitives through Manager with session/tab resolution.

1. Add `NetworkWait` method to Manager — resolves session/tab, checks for Lightpanda (returns error), delegates to `WaitForRequest` (files: `browser/manager.go`)
2. Add `NetworkGetBody` method to Manager — resolves session/tab, delegates to `GetResponseBody` (files: `browser/manager.go`)

### Phase 3: CLI — `wait` & `body`

Expose via CLI, register in browser command tree.

1. Create `cmd/tap/browser_network.go` with `browserNetworkCmd()` group, `browserNetworkWaitCmd()`, `browserNetworkBodyCmd()` (files: `cmd/tap/browser_network.go`)
   - `wait` flags: `--url-pattern`, `--method`, `--resource-type`, `--timeout` (default 30s), `--body` (include response body), `--format`
   - `body` args: `<requestId>`, flags: `--format`
2. Register `browserNetworkCmd()` in `browserCmd()` (files: `cmd/tap/browser.go`)

### Phase 4: Network Log (streaming capture)

Add streaming network log support.

1. Implement `EnableNetworkLog` — returns `(<-chan NetworkEntry, func(), error)` (files: `browser/network.go`)
   - **Channel**: buffered (256 entries). If buffer fills, newest entries are dropped (non-blocking send).
   - **Cancel func**: stops the internal goroutine, which then closes the channel.
   - **Goroutine lifecycle**: goroutine exits when cancel is called OR parent context is done. Always closes channel on exit.
   - Caller must either drain the channel or call cancel to avoid goroutine leak.
2. Add `NetworkLog` Manager method — resolves session/tab, checks Lightpanda, delegates (files: `browser/manager.go`)
3. Add `browserNetworkLogCmd()` — streams NDJSON to stdout. Flags: `--url-pattern`, `--method`, `--resource-type`, `--timeout` (0 = indefinite), `--body` (files: `cmd/tap/browser_network.go`)

### Phase 5: Fetch Domain Interception

Add active request interception (block/mock/headers).

1. Define `InterceptRule` type (files: `browser/network.go`):
   - `Filter NetworkFilter` — URL pattern, methods, resource types
   - `Block bool` — fail the request (mutually exclusive with `Mock*`)
   - `AddHeaders map[string]string` — inject/override request headers
   - `MockStatus int` — mock response status code (requires `MockBody`)
   - `MockBody string` — mock response body
   - `MockHeaders map[string]string` — mock response headers (default `Content-Type: application/json` if not set)
   - Validation: `Block` and `MockBody` are mutually exclusive. Error if both set.
   - **Rule semantics**: `SetInterceptRules` is **replace-all** (not additive). Pass the complete set of desired rules each time. Pass `nil` to clear.
2. Implement `SetInterceptRules` — enables `Fetch` domain with `RequestPattern` entries derived from rules, handles `EventRequestPaused` in a goroutine to apply block/mock/header logic, returns cancel func for cleanup (files: `browser/network.go`)
3. Implement `ClearIntercept` — disables `Fetch` domain on the target (files: `browser/network.go`)
4. Add `NetworkIntercept` and `NetworkClearIntercept` Manager methods — with Lightpanda check (files: `browser/manager.go`)
5. Add `browserNetworkInterceptCmd()` and `browserNetworkClearCmd()` CLI commands (files: `cmd/tap/browser_network.go`)
   - `intercept` flags: `--url-pattern`, `--method`, `--resource-type`, `--block`, `--header "Key: Value"` (repeatable), `--respond <body>`, `--status <code>`, `--content-type <mime>`
   - `clear` — no extra flags, just `--session`/`--tab`
6. Add integration test: launch headless Chrome, set block rule, navigate to a page that loads a resource, verify the blocked request fails (files: `browser/network_test.go`)

## Testing Strategy

- **Unit tests**: `matchURL` with various glob patterns including multi-segment paths, edge cases (empty pattern, exact match)
- **Unit tests**: `matchesFilter` with various filter combinations (URL + method, URL + resource type, empty filter)
- **Unit tests**: `InterceptRule` validation (block + mock mutual exclusivity)
- **Integration test (Phase 1)**: launch headless Chrome, navigate to a test page, use `WaitForRequest` to capture a known request, verify entry fields
- **Integration test (Phase 5)**: launch headless Chrome, set block rule, navigate, verify blocked request
- **Manual tests**: `tap browser network wait` against a real site with known API patterns
- Existing `go test ./... -timeout 60s -race` must pass throughout

## Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Long-lived CDP session leaks if caller doesn't cancel | Medium | Cancel func in `withTargetListen`; defer in Manager methods; document in godoc |
| `Network.getResponseBody` fails for redirected/cached requests | Low | Return error gracefully; document limitation |
| `Fetch` domain + `Network` domain conflict | Low | They operate independently per CDP spec; test co-existence |
| Lightpanda doesn't support Network/Fetch domains | Medium | Check session mode in Manager methods; return clear error |
| Channel buffer overflow in `EnableNetworkLog` | Low | Non-blocking send with buffer of 256; document drop behavior |
| `filepath.Match` doesn't match `*` across `/` | Medium | Custom `matchURL` using regex with `.*` for `*` wildcards |

## Open Questions

_None — all resolved to assumptions:_

- **Assumption**: Lightpanda does not support Network/Fetch domains. Manager methods return a clear error.
- **Assumption**: URL matching uses custom `matchURL` where `*` matches any characters including `/` (like CDP `urlPattern`).
- **Assumption**: Network log is ephemeral (process-lifetime only), not persisted to the store.
- **Assumption**: `SetInterceptRules` is replace-all semantics. The CLI sends the complete rule set each call.
- **Assumption**: `Block` and `MockBody` are mutually exclusive in `InterceptRule`.
- **Assumption**: `network clear` only clears Fetch domain rules (interception). Passive Network domain capture is stopped by its own cancel func.

## Review Feedback

### Round 1 (Codex Reviewer)

**Verdict: CHANGES NEEDED** — 8 actionable items integrated:

1. ✅ Specified `withTargetListen` contract — signature, detach-without-close, event listener ordering
2. ✅ Added `--timeout` flag to `network wait` (default 30s)
3. ✅ Specified `EnableNetworkLog` channel semantics — buffer 256, cancel closes channel, goroutine lifecycle
4. ✅ Clarified `intercept` rule semantics — replace-all, `Block`/`MockBody` mutually exclusive, added `--status`/`--content-type` flags
5. ✅ Added `--method`/`--resource-type` filters to `intercept`
6. ✅ `NetworkEntry` includes `RequestID` field (already in types, now explicitly noted in success criteria)
7. ✅ Added Fetch domain integration test to Phase 5
8. ✅ Replaced `path.Match` with custom `matchURL` using regex for multi-segment wildcard support

**Minor nits addressed:**
- Removed `NetworkTiming` type (not needed for v1, can add later)
- Clarified that `network clear` only clears Fetch domain (not Network domain)

## Final Status

_(Updated after implementation completes)_
