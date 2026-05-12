/* @meta
{
  "name": "twitter/getxapi-tweet-detail",
  "description": "Fetch tweet detail by ID via getxapi.com",
  "domain": "api.getxapi.com",
  "args": {
    "id": {"required": true, "description": "Tweet ID"}
  },
  "runtime": "http",
  "env": {
    "GET_X_API_KEY": {"required": true, "description": "API key for getxapi.com"}
  },
  "headers": {
    "Authorization": "Bearer ${GET_X_API_KEY}"
  },
  "readOnly": true,
  "example": "tap site twitter/getxapi-tweet-detail '2019264360682778716'"
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id', hint: 'Provide a tweet ID'};

  const resp = await fetch(`https://api.getxapi.com/twitter/tweet/detail?id=${encodeURIComponent(args.id)}`, {
    headers: {'accept': 'application/json'}
  });

  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  const body = await resp.json();
  if (body.status !== 'success') return {error: body.msg || 'API error'};

  return body.data;
}
