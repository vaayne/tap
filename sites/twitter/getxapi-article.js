/* @meta
{
  "description": "Fetch X / Twitter article by tweet ID via getxapi.com",
  "domain": "api.getxapi.com",
  "args": {
    "id": {
      "required": true,
      "description": "Tweet / article ID"
    }
  },
  "headers": {
    "Authorization": "Bearer ${GET_X_API_KEY}"
  }
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id', hint: 'Provide a tweet / article ID'};

  const resp = await fetch(`https://api.getxapi.com/twitter/tweet/article?id=${encodeURIComponent(args.id)}`, {
    headers: {'accept': 'application/json'}
  });

  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  const body = await resp.json();
  if (body.status !== 'success') return {error: body.msg || 'API error'};

  return body.article;
}
