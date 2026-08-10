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

- `headers` handles auth. Environment variables are inferred from `${VAR}`
  references and expanded before execution; unresolved headers are omitted.
- `name`, `runtime`, and `env` are not metadata fields. The file path determines
  the script name and agent-browser is the only runtime.
- Never read secrets from `args`; metadata headers are merged into script
  `fetch()` calls by Tap's browser-side wrapper.
- Return `{error: 'message'}` on failure; any other JSON value is success

## After adding

```bash
mise run build           # re-embeds sites/**/*.js
tap site <site>/<action> [key=value ...]
```
