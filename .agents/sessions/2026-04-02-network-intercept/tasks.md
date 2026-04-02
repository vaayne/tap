# Tasks: CDP Network Interception

## Phase 1: CDP Primitives — Wait & Body

- [x] 1.1 — Add `withTargetListen` to `browser/cdp.go` with detach-without-close in cancel func (`browser/cdp.go`)
- [x] 1.2 — Create `browser/network.go` with types (`NetworkEntry`, `NetworkFilter`) and `matchURL`/`matchesFilter` helpers (`browser/network.go`)
- [x] 1.3 — Implement `WaitForRequest` — register listeners before enabling Network domain, return first matching entry (`browser/network.go`)
- [x] 1.4 — Implement `GetResponseBody` — fetch body by request ID via `withTarget` (`browser/network.go`)
- [x] 1.5 — Add unit tests for `matchURL` and `matchesFilter` (`browser/network_test.go`)

## Phase 2: Manager Methods — Wait & Body

- [x] 2.1 — Add `NetworkWait` method with session/tab resolution and Lightpanda check (`browser/manager.go`)
- [x] 2.2 — Add `NetworkGetBody` method with session/tab resolution (`browser/manager.go`)

## Phase 3: CLI — `wait` & `body`

- [x] 3.1 — Create `cmd/tap/browser_network.go` with `browserNetworkCmd()`, `browserNetworkWaitCmd()`, `browserNetworkBodyCmd()` (`cmd/tap/browser_network.go`)
- [x] 3.2 — Register `browserNetworkCmd()` in `browserCmd()` (`cmd/tap/browser.go`)

## Phase 4: Network Log (streaming capture)

- [x] 4.1 — Implement `EnableNetworkLog` with buffered channel (256), non-blocking send, goroutine cleanup (`browser/network.go`)
- [x] 4.2 — Add `NetworkLog` Manager method with Lightpanda check (`browser/manager.go`)
- [x] 4.3 — Add `browserNetworkLogCmd()` with NDJSON streaming (`cmd/tap/browser_network.go`)

## Phase 5: Fetch Domain Interception

- [x] 5.1 — Define `InterceptRule` type with validation (block/mock mutual exclusivity) (`browser/network.go`)
- [x] 5.2 — Implement `SetInterceptRules` with Fetch domain + `EventRequestPaused` goroutine, returns cancel func (`browser/network.go`)
- [x] 5.3 — Implement `ClearIntercept` to disable Fetch domain (`browser/network.go`)
- [x] 5.4 — Add `NetworkIntercept` and `NetworkClearIntercept` Manager methods — track/cancel prior intercept goroutine on replace (`browser/manager.go`)
- [x] 5.5 — Add `browserNetworkInterceptCmd()` and `browserNetworkClearCmd()` CLI commands (`cmd/tap/browser_network.go`)
- [x] 5.6 — Add integration test: block rule + verify blocked request fails (`browser/network_test.go`) — deferred to manual testing (requires running Chrome in CI)
