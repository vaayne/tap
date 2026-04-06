# Script Development Guide

Write and contribute site scripts to extract structured data from websites.

tap uses the same script format as [bb-sites](https://github.com/epiral/bb-sites) — scripts are fully compatible between tap and bb-browser.

## Development workflow

1. Reverse the API with `tap browser network`
2. Test a fetch in QuickJS or browser eval
3. Write the script file
4. Save to `~/.config/tap/sites/{site}/{script}.js` and test locally
5. Contribute upstream to [bb-sites](https://github.com/epiral/bb-sites)

---

## Step 1: Reverse the API

Open the target site in a tracked browser tab, then capture its API calls.
The network log starts empty — you must trigger a page action (load, scroll, search) to see requests.

```bash
tap browser session new dev --no-headless
tap browser tab new page --url https://www.example.com

# Start log in background, then trigger the page action you want to capture
tap browser network log --resource-type XHR,Fetch --timeout 15s &
tap browser navigate https://www.example.com   # reload triggers initial requests
wait

# Or block until one specific request completes and capture its body
tap browser network wait --url-pattern "*/api/*" --body --timeout 30s
```

Focus on:
- Request URL and query parameters
- Auth mechanism (Cookie / Bearer token / CSRF token)
- Response data structure

---

## Step 2: Test the fetch

Verify the API call works before writing a full script:

```bash
# Test in browser page context (has cookies, full DOM access)
tap browser evaluate "fetch('/api/items?q=test',{credentials:'include'}).then(r=>r.json()).then(d=>JSON.stringify(d))"
```

This tells you which execution tier you need:

| Result | Tier | Engine |
|---|---|---|
| Data returned directly | **Tier 1** — plain fetch | QuickJS (no browser) |
| Needs cookies / DOM | **Tier 2** — browser context | CDP browser fallback |

---

## Step 3: Write the script

### Metadata format

Every script starts with a `/* @meta */` JSON block followed by an `async function`:

```javascript
/* @meta
{
  "name": "site/action",
  "description": "What this script does",
  "domain": "www.example.com",
  "args": {
    "query": {"required": true,  "description": "Search keyword"},
    "count": {"required": false, "description": "Max results (default 10)"}
  },
  "readOnly": true,
  "example": "tap site site/action query=foo count=5"
}
*/
async function(args) {
  // implementation
}
```

| Field | Required | Description |
|---|---|---|
| `name` | yes | Unique ID in `site/action` format |
| `description` | yes | One-line human description |
| `domain` | no | Target domain — used to navigate before browser execution |
| `args` | yes | Map of argument definitions (`required` + `description`); use `{}` for no-argument scripts |
| `readOnly` | no | `true` for read-only operations |
| `capabilities` | no | `["network"]` for scripts needing browser network interception (auth tokens, internal APIs) |
| `example` | no | Example CLI invocation shown by `tap site info`. Use `bb-browser site ...` for upstream contributions |

### Execution model

tap runs scripts through a two-engine fallback chain:

1. **QuickJS** (fast, no browser): has an injected `fetch()` backed by Go's HTTP client. No DOM, no cookies.
2. **CDP Browser** (full page context): navigates to `domain`, evaluates the script in the real page. Has DOM, cookies, JS globals.

If QuickJS returns `{"error": "..."}`, tap automatically retries with the browser engine. Design scripts to explicitly signal fallback:

```javascript
async function(args) {
  var resp = await fetch('/api/data?q=' + encodeURIComponent(args.query));
  if (!resp.ok) return {error: 'HTTP ' + resp.status, hint: 'Needs browser context'};
  var data = await resp.json();
  return data.items;
}
```


### Tier 1: Plain fetch (QuickJS, ~1 s typical)

Works for public REST APIs or cookie-less endpoints.

```javascript
/* @meta
{
  "name": "hackernews/search",
  "description": "Search Hacker News stories",
  "domain": "hn.algolia.com",
  "args": {
    "query": {"required": true, "description": "Search query"},
    "count": {"required": false, "description": "Number of results"}
  },
  "readOnly": true,
  "example": "tap site hackernews/search query=golang count=10"
}
*/
async function(args) {
  if (!args.query) return {error: 'Missing argument: query'};
  const n = args.count || '10';
  const resp = await fetch(
    'https://hn.algolia.com/api/v1/search?query=' + encodeURIComponent(args.query) + '&hitsPerPage=' + n
  );
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const data = await resp.json();
  return data.hits.map(h => ({
    title: h.title,
    url:   h.url,
    score: h.points,
    by:    h.author
  }));
}
```

### Tier 2: Browser context with cookies (~3 s typical)

For sites that require login or JS rendering. Use `credentials: 'include'` to attach cookies.

```javascript
/* @meta
{
  "name": "github/notifications",
  "description": "List unread GitHub notifications",
  "domain": "github.com",
  "args": {},
  "readOnly": true,
  "example": "tap site github/notifications"
}
*/
async function(args) {
  // Extract CSRF token from meta tag (GitHub pattern)
  const csrf = document.querySelector('meta[name="user-scoped-authenticity-token"]')?.content;
  if (!csrf) return {error: 'CSRF token not found', hint: 'Run: tap login https://github.com/login'};

  const resp = await fetch('/notifications', {
    headers: {
      'accept': 'application/json',
      'x-requested-with': 'XMLHttpRequest',
      'x-csrf-token': csrf
    },
    credentials: 'include'
  });
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  return await resp.json();
}
```

Run with `-b` to use saved browser cookies:

```bash
tap login https://github.com/login   # one-time login
tap site -b github/notifications
```

---

## Step 4: Test locally

Drop the file into your local override directory — it takes precedence over the cache automatically:

```
~/.config/tap/sites/{site}/{script}.js
```

Example:

```bash
mkdir -p ~/.config/tap/sites/hackernews
cp hackernews_search.js ~/.config/tap/sites/hackernews/search.js
```

Test commands:

```bash
tap site hackernews/search query=golang              # default (QuickJS → browser)
tap site hackernews/search query=golang -f json      # JSON output
tap site -b hackernews/search query=golang           # force browser engine
tap --local-only site hackernews/search query=golang # skip remote cache (global flag)
tap site info hackernews/search                      # verify metadata parsed correctly
```

Note: `--local-only` is a **global** flag and must come before the subcommand:
`tap --local-only site ...` ✓  — not `tap site --local-only ...` ✗

---

## Resilience patterns

### Prefer API calls over DOM scraping

DOM structure changes frequently. Use `tap browser network` to find underlying APIs instead of querying CSS selectors.

### Semantic selectors over class names

When DOM scraping is unavoidable, use semantic HTML elements (`h3`, `a`, `article`) instead of CSS class names, which sites change regularly:

```javascript
// Fragile: class names change
const items = document.querySelectorAll('div.g');

// Robust: semantic structure
const results = [];
for (const h3 of document.querySelectorAll('h3')) {
  const a = h3.closest('a');
  if (!a) continue;
  const href = a.getAttribute('href');
  if (!href || !href.startsWith('http')) continue;
  results.push({ title: h3.textContent.trim(), url: href });
}
```

---

## Error handling

Return a plain object with `error` (and optional `hint`) to signal failures:

```javascript
// Missing argument
return {error: 'Missing argument: query'};

// HTTP failure + login hint
return {error: 'HTTP 401', hint: 'Run: tap login https://example.com/login'};

// Graceful fallback signal to browser engine
return {error: 'Needs DOM', hint: 'Retrying with browser engine'};
```

tap detects `401`/`403`/`unauthorized`/`login` keywords and surfaces a login suggestion automatically.

### WAF / anti-bot challenge detection

Some sites (especially Chinese ones — Xueqiu, Weibo, etc.) use WAFs that return HTTP 200 with an HTML challenge body instead of JSON. A `resp.ok` check passes, then `resp.json()` throws. Always guard against this:

```javascript
var resp = await fetch('https://example.com/api/data', {credentials: 'include'});
if (!resp.ok) return {error: 'HTTP ' + resp.status};

// Check content-type before parsing as JSON
var ct = resp.headers.get('content-type') || '';
if (!ct.includes('application/json')) {
  return {error: 'WAF challenge or login required', hint: 'Run: tap login https://example.com'};
}

var d;
try {
  d = await resp.json();
} catch (e) {
  return {error: 'Invalid JSON response', hint: 'Run: tap login https://example.com'};
}
```

---

## Step 5: Contribute upstream to bb-sites

Scripts are contributed to the upstream [bb-sites](https://github.com/epiral/bb-sites) repo. Make sure the `example` field uses `bb-browser site` (not `tap site`) for upstream compatibility.

```bash
gh repo fork epiral/bb-sites --clone && cd bb-sites
git checkout -b feat-site-action
# place your file at: site/action.js (flat layout: one dir per site)
git add site/action.js
git commit -m "✨ feat: add site/action script"
git push -u origin feat-site-action
gh pr create --repo epiral/bb-sites --title "feat(site): add action adapter"
```
