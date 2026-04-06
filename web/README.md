# Web

Web UI for Tap script discovery and distribution.

This app serves the public script catalog at `tap.vaayne.com`, exposes API routes for script lookup and sync, and stores script metadata plus usage events in Cloudflare D1.

## Stack

- TanStack Start + React 19
- Cloudflare Workers
- Cloudflare D1
- Tailwind CSS v4
- Vite Plus / Vitest
- pnpm

## What lives here

- Browse and search available site scripts
- View script details and raw script content
- Expose sync endpoints for the CLI and `bb-sites` automation
- Record usage events for popularity ranking

## Structure

```text
web/
├── src/
│   ├── components/      # UI components
│   ├── lib/             # D1 access, auth, shared types, server fns
│   ├── routes/          # Pages and API routes
│   │   └── api/         # /api/search, /api/scripts, /api/sync, /api/batch, /api/usage
│   ├── router.tsx       # TanStack router setup
│   └── styles.css       # Global styles
├── migrations/          # D1 schema migrations
├── wrangler.jsonc       # Cloudflare Worker + D1 config
├── vite.config.ts       # Vite Plus + TanStack Start config
└── package.json
```

## Local development

From the repo root:

```bash
cd web
pnpm install
pnpm dev
```

Common commands:

```bash
pnpm dev        # local dev server
pnpm build      # production build
pnpm preview    # preview built app
pnpm test       # vitest
pnpm cf-typegen # regenerate Cloudflare types
pnpm deploy     # build and deploy with wrangler
```

## Data model

D1 tables are defined in `migrations/0001_init.sql`.

- `scripts` — canonical script metadata and source
- `usage_events` — raw usage reports used for popularity ranking

Apply local/remote D1 migrations with Wrangler as needed.

## API routes

### Public

- `GET /api/search?q=&sort=&site=` — search and filter scripts
- `GET /api/scripts/:name` — fetch full script detail
- `GET /api/scripts/:name/content` — fetch raw JavaScript source
- `GET /api/sync` — sync manifest `{ name, hash, updatedAt }`
- `POST /api/usage` — report a script usage event

### Protected

- `POST /api/batch` — replace the script dataset in D1
  - Requires `X-Tap-Secret` header
  - Validates JSON payload and enforces a 500 KB limit
  - Used by automation syncing from `bb-sites`

## Configuration

Important runtime config lives in `web/wrangler.jsonc`.

Current bindings/config:

- `DB` — D1 database binding
- `TAP_SCRIPTS_SECRET` — shared secret for `POST /api/batch`
- custom domain route for `tap.vaayne.com`

For local development, provide any required Worker secrets with Wrangler.

## Implementation notes

- Prefer server functions in `src/lib/server-fns.ts` for app-facing data loading.
- Keep DB access in `src/lib/db.ts`.
- API routes should stay thin and delegate validation/business logic to `lib/` helpers.
- Preserve compatibility with the Go CLI sync flow and the separate `bb-sites` repo.

## When updating this app

Also consider whether related docs need updates in:

- `README.md`
- `docs/`
- `skills/tap-web/`

Keep behavior and docs in sync when API routes, sync semantics, or D1 schema change.
