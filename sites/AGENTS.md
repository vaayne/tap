# Adding Site Scripts

Scripts live in `sites/<site>/<action>.js` and are embedded into the binary at build time (`embed.go`).

## Structure

Copy an existing script as a starting point:
- `exa/search.js` — API key auth, SSE response parsing
- `twitter/getxapi-tweet-detail.js` — simple authenticated GET
- `twitter/post-tweet.js` — authenticated POST, non-read-only

## Meta fields

See field-level docs on the `Meta` struct in `script/parser.go`.

## Key rules

- `env` + `headers` handle auth — **never re-read env vars inside the function body** (`args.MY_API_KEY` is always undefined). Meta headers are resolved via `tap.go:165` and injected into every `fetch()` at `engine/quickjs.go:131`, before any JS-level headers.
- Return `{error: 'message'}` on failure; any other JSON value is success

## After adding

```bash
mise run build           # re-embeds sites/**/*.js
tap site <site>/<action> [key=value ...]
```
