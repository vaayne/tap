# Handoff

<!-- Append a new phase section after each phase completes. -->

## Implementation Notes

- **Reviewer implementation note (Round 2)**: `NetworkIntercept` in the Manager must cancel any prior intercept goroutine before setting new rules, to prevent goroutine leaks on consecutive `intercept` calls. Tracked in Task 5.4.

## Phase 1: CDP Primitives — Wait & Body (Complete)

**Commits:**
- `aa11306` — ✨ feat: add withTargetListen for long-lived CDP sessions
- `a2a9d9b` — ✨ feat: add network types, URL matching, WaitForRequest, and GetResponseBody
- `10536ca` — ✅ test: add unit tests for matchURL and matchesFilter

**Files created/modified:**
- `browser/cdp.go` — Added `withTargetListen(ctx, debugURL, targetID)` → `(context.Context, func(), error)`. Uses detach-without-close trick (clears `TargetID` before cancel). Runs empty `chromedp.Run` to force CDP attachment before returning.
- `browser/network.go` (new) — Contains:
  - Types: `NetworkEntry` (with `RequestID`, `URL`, `Method`, `Status`, `MIMEType`, `ResourceType`, `ReqHeaders`, `ResHeaders`, `Body`, `Error`, `Timestamp`), `NetworkFilter` (`URLPattern`, `Methods`, `ResourceTypes`)
  - `matchURL(pattern, url)` — glob-to-regex conversion where `*` matches any char including `/`. Uses `regexp.QuoteMeta` + replace `\*` → `.*`.
  - `matchesFilter(entry, filter)` — combines URL, method (case-insensitive), resource type (case-insensitive) matching. Empty filter matches all.
  - `WaitForRequest(ctx, debugURL, targetID, filter, includeBody)` — Uses `withTargetListen`, registers `chromedp.ListenTarget` BEFORE `network.Enable()`. Tracks 4 events: `RequestWillBeSent`, `ResponseReceived`, `LoadingFinished`, `LoadingFailed`. Returns first matching entry after response+loading completes. Optional body fetch via `getResponseBodyOn`. Context cancellation = timeout.
  - `GetResponseBody(ctx, debugURL, targetID, requestID)` — Uses existing `withTarget` pattern. Returns `[]byte`.
  - `headersToMap` — converts CDP `network.Headers` (map[string]any) to `map[string]string`.
- `browser/network_test.go` (new) — 33 test cases across `TestMatchURL` (16 cases) and `TestMatchesFilter` (17 cases).

**API note:** CDP's `network.GetResponseBody` returns `[]byte` (not string+base64 flag), so `GetResponseBody` and `NetworkEntry.Body` use `[]byte`/`string` directly. No `DecodeBody` helper needed.

**All tests pass:** `go test ./... -timeout 60s -race` ✓

**Ready for Phase 2:** Manager methods (`NetworkWait`, `NetworkGetBody`) can wrap `WaitForRequest` and `GetResponseBody` using the `resolveTarget` pattern from existing Manager methods.
