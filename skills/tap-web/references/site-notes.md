# Site Notes — Remember How Sites Work

Maintain per-site knowledge so you don't rediscover access patterns on every session.

## Location

```
$XDG_CONFIG_HOME/tap/site-notes/{domain}.md    # if XDG_CONFIG_HOME set
~/.config/tap/site-notes/{domain}.md           # fallback
```

One file per domain. Read before interacting with a site. Update after learning something new.

## Workflow

1. **Before accessing a site**: check if `site-notes/{domain}.md` exists — read it first
2. **First visit**: create the file with whatever you learn
3. **Something stops working**: update the file — mark the old approach as broken with a date, log the new approach
4. **Found a better path**: replace the old approach (e.g., discovered API endpoint that replaces DOM scraping)

## What to log

| Category | Examples |
|---|---|
| **Access pattern** | Needs `-b`? WAF behavior? Cloudflare? |
| **Working endpoints** | API URLs, query params, response shape |
| **Broken paths** | What doesn't work and why (with date) |
| **Selectors** | CSS selectors that reliably extract data |
| **Auth** | How to log in, cookie expiry, CSRF token location |
| **Gotchas** | Rate limits, anti-bot detection, content-type traps |
| **Token-efficient path** | Best extraction method for this site |

## File format

```markdown
# {domain}
Last updated: {date}

## Access pattern
- API requires browser cookies (`-b` flag)
- WAF returns 200 + HTML instead of JSON — check content-type

## Working endpoints
- `GET /api/v1/data?q=foo` — returns JSON array
- Best extraction: `tap browser network wait --url-pattern "*/api/v1/*" --body`

## Broken
- `GET /old/api` — returns 403 since 2026-04-01

## Selectors
- `.result-item h3 a` — search result links (verified 2026-04-03)

## Auth
- `tap login https://example.com/login`
- Cookies expire after ~24h
```

## Key principles

- **Be specific**: log exact URLs, params, and response shapes — not vague descriptions
- **Date your entries**: especially broken paths, so you know when to retry
- **Log the cheapest path**: record which extraction method works best (API > text > evaluate)
- **Update, don't append**: keep the file current, remove outdated info
