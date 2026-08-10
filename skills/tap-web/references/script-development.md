# Site Script Development

Site scripts live at `{site}/{action}.js`; the relative path is their name.

## Workflow

1. Use agent-browser to inspect the site or its network API.
2. Test the request with `agent-browser eval --stdin`.
3. Write the script under `~/.config/tap/sites/{site}/{action}.js`.
4. Run `tap --local-only site {site}/{action} key=value`.

```bash
agent-browser open https://example.com
agent-browser network requests --filter api

cat <<'JS' | agent-browser eval --stdin
fetch('/api/items?q=test', {credentials: 'include'}).then(r => r.json())
JS
```

## Metadata

```javascript
/* @meta
{
  "description": "Search example.com",
  "domain": "example.com",
  "args": {
    "query": {"required": true, "description": "Search query"}
  },
  "headers": {
    "Authorization": "Bearer ${EXAMPLE_TOKEN}"
  },
  "readOnly": true
}
*/

async function(args) {
  const response = await fetch(`/api/search?q=${encodeURIComponent(args.query)}`);
  if (!response.ok) return {error: `HTTP ${response.status}`};
  return response.json();
}
```

`name`, `runtime`, and `env` are not metadata fields. Environment variables are
inferred from `${VAR}` references in headers. Unresolved headers are omitted.

Metadata headers are applied before domain navigation, merged into every script
`fetch()` call, then cleared so credentials do not linger in the shared session.

## Errors

Return a plain object with `error` and an optional `hint`:

```javascript
return {error: 'Missing argument: query'};
return {error: 'HTTP 401', hint: 'Authenticate in the current agent-browser session'};
```

Scripts are contributed upstream to [bb-sites](https://github.com/epiral/bb-sites).
