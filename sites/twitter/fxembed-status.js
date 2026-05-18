/* @meta
{
  "name": "twitter/fxembed-status",
  "description": "Fetch an X / Twitter post by ID via FxEmbed",
  "domain": "api.fxtwitter.com",
  "args": {
    "id": {"required": true, "description": "Tweet / post ID"},
    "about_account": {"required": false, "description": "Truthy value to include account metadata when available"},
    "lang": {"required": false, "description": "Target language for inline translations, e.g. en, es, zh-cn"}
  },
  "runtime": "http",
  "readOnly": true,
  "example": "tap site twitter/fxembed-status id=20 about_account=1"
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id', hint: 'Provide a tweet / post ID'};

  const params = [];
  if (args.about_account) params.push('about_account=' + encodeURIComponent(args.about_account));
  if (args.lang) params.push('lang=' + encodeURIComponent(args.lang));

  let url = `https://api.fxtwitter.com/2/status/${encodeURIComponent(args.id)}`;
  if (params.length) url += '?' + params.join('&');

  const resp = await fetch(url, {
    headers: {'accept': 'application/json'}
  });

  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  return await resp.json();
}
