/* @meta
{
  "name": "twitter/xquik-tweet-detail",
  "description": "Fetch tweet detail by ID via Xquik",
  "domain": "xquik.com",
  "args": {
    "id": {"required": true, "description": "Tweet ID"}
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
  "example": "tap site twitter/xquik-tweet-detail '2019264360682778716'"
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id', hint: 'Provide a tweet ID'};

  const resp = await fetch(`https://xquik.com/api/v1/x/tweets/${encodeURIComponent(args.id)}`, {
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
