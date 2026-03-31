# Tasks: Tap Refactor

## Phase 1: Project scaffolding & script package

- [x] 1.1 — Init module `github.com/vaayne/tap`, create `cmd/tap/main.go` stub (`go.mod`, `cmd/tap/main.go`)
- [x] 1.2 — Move meta parsing to `script/parser.go` (`script/parser.go`)
- [x] 1.3 — Move registry to `script/registry.go` (`script/registry.go`)
- [x] 1.4 — Unit test parser + registry (`script/parser_test.go`, `script/registry_test.go`)

## Phase 2: Engine package

- [ ] 2.1 — Define Engine interface + orchestrator (`engine/engine.go`)
- [ ] 2.2 — Move QuickJS runner to `engine/quickjs.go` (`engine/quickjs.go`)
- [ ] 2.3 — Move CDP browser to `engine/browser.go` (`engine/browser.go`)
- [ ] 2.4 — Unit test engine fallback logic (`engine/engine_test.go`)

## Phase 3: Fetch package with go-defuddle

- [ ] 3.1 — Create fetch package with Result type + Fetch function (`fetch/fetch.go`)
- [ ] 3.2 — Unit test with sample HTML (`fetch/fetch_test.go`)

## Phase 4: Top-level client API

- [ ] 4.1 — Create Client struct with New(), Close() (`tap.go`)
- [ ] 4.2 — Create functional options (`options.go`)
- [ ] 4.3 — Implement Client.RunScript() (`tap.go`)
- [ ] 4.4 — Implement Client.Fetch() (`tap.go`)
- [ ] 4.5 — Integration test (`tap_test.go`)

## Phase 5: CLI with urfave/cli

- [ ] 5.1 — Set up urfave/cli app with global flags (`cmd/tap/main.go`)
- [ ] 5.2 — `tap site list` command (`cmd/tap/main.go`)
- [ ] 5.3 — `tap site <name> [args]` command (`cmd/tap/main.go`)
- [ ] 5.4 — `tap fetch <url>` command with `--json`/`--markdown` flags (`cmd/tap/main.go`)
- [ ] 5.5 — Load .env, wire env vars to client options (`cmd/tap/main.go`)

## Phase 6: Cleanup & docs

- [ ] 6.1 — Remove old root-level Go files
- [ ] 6.2 — Update README.md
- [ ] 6.3 — Update .gitignore, add .env.example
- [ ] 6.4 — Add LICENSE (MIT)
- [ ] 6.5 — Final build + smoke test
- [ ] 6.6 — Push to github.com/vaayne/tap
