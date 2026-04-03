# AGENTS.md

## Project

`web/` is the Tap web app: a TanStack Start + React frontend deployed to Cloudflare Workers with Cloudflare D1 as storage.

It powers script browsing at `tap.vaayne.com` and exposes API endpoints used for search, script retrieval, sync manifests, usage reporting, and batch updates from automation.

## Stack

- React 19
- TanStack Start / TanStack Router
- Cloudflare Workers
- Cloudflare D1
- Tailwind CSS v4
- Vite Plus
- Vitest
- pnpm

## Commands

Run from `web/` unless noted otherwise:

```bash
pnpm install
pnpm dev        # start local dev server
pnpm build      # production build
pnpm preview    # preview build locally
pnpm test       # run tests
pnpm cf-typegen # regenerate wrangler/cloudflare types
pnpm deploy     # deploy to Cloudflare Workers
```

## Architecture

```text
src/routes/           → file-based app routes and API endpoints
src/lib/db.ts         → D1 queries and data mapping
src/lib/server-fns.ts → app-facing server functions for loaders
src/lib/auth.ts       → X-Tap-Secret validation for batch sync
migrations/           → D1 schema
wrangler.jsonc        → Worker, domain, and D1 bindings
```

Key API endpoints:

- `GET /api/search`
- `GET /api/scripts/:name`
- `GET /api/scripts/:name/content`
- `GET /api/sync`
- `POST /api/usage`
- `POST /api/batch` (protected by `X-Tap-Secret`)

## Data model

D1 schema currently includes:

- `scripts` — script metadata, content, hashes, timestamps
- `usage_events` — raw usage events for popularity counts

If you change schema, update:

1. `migrations/`
2. `src/lib/types.ts`
3. `src/lib/db.ts`
4. any affected API routes/UI/docs

## Code style

- Prefer `type` aliases over `interface`.
- Keep route handlers thin; move shared logic into `src/lib/`.
- Reuse existing UI primitives in `src/components/ui/`.
- Match existing import style and naming conventions.
- Keep components focused and small.

## Guardrails

- Do not break the CLI sync contract without updating the Go CLI and docs.
- `POST /api/batch` must remain authenticated and payload-limited.
- Keep public response shapes stable unless the caller updates too.
- When changing API behavior or schema, update docs in both repo root and `web/`.

## Testing and verification

Before finishing web changes, prefer to run:

```bash
cd web
pnpm build
pnpm test
```

If API, binding, or D1 config changes, also review:

- `web/wrangler.jsonc`
- `web/migrations/`
- root `README.md`
- `skills/tap-web/`

## Commits

Use emoji-prefixed Conventional Commits, for example:

- `✨ feat: add web script filters`
- `🐛 fix: handle missing batch auth secret`
- `📝 docs: add web app documentation`
- `♻️ refactor: simplify D1 query helpers`
