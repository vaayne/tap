/* @meta
{
  "name": "twitter/xquik-article",
  "description": "Fetch X / Twitter article by tweet ID via Xquik",
  "domain": "xquik.com",
  "args": {
    "id": {"required": true, "description": "Tweet / article ID"}
  },
  "runtime": "http",
  "env": {
    "XQUIK_API_KEY": {"required": true, "description": "API key for Xquik"}
  },
  "headers": {
    "x-api-key": "${XQUIK_API_KEY}",
    "xquik-api-contract": "2026-04-29"
  },
  "readOnly": true,
  "example": "tap site twitter/xquik-article '1905545699552375179'"
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id', hint: 'Provide a tweet / article ID'};

  const resp = await fetch(`https://xquik.com/api/v1/x/articles/${encodeURIComponent(args.id)}`, {
    headers: {'accept': 'application/json'}
  });

  const text = await resp.text();
  let body = {};
  try {
    body = JSON.parse(text);
  } catch (_) {
    body = {};
  }

  if (!resp.ok) {
    return {error: body.message || body.error || ('HTTP ' + resp.status)};
  }

  return body;
}
