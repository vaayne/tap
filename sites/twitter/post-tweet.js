/* @meta
{
  "name": "twitter/post-tweet",
  "description": "Post a tweet via the official X API",
  "domain": "api.x.com",
  "args": {
    "text": {"required": true, "description": "Tweet text"}
  },
  "runtime": "http",
  "env": {
    "X_ACCESS_TOKEN": {"required": true, "description": "User access token for the X API"}
  },
  "headers": {
    "Authorization": "Bearer ${X_ACCESS_TOKEN}"
  },
  "readOnly": false,
  "example": "tap site twitter/post-tweet 'Hello from the X API!'"
}
*/

async function(args) {
  if (!args.text) return {error: 'Missing argument: text', hint: 'Provide tweet text'};

  const resp = await fetch('https://api.x.com/2/tweets', {
    method: 'POST',
    headers: {
      'accept': 'application/json',
      'content-type': 'application/json'
    },
    body: JSON.stringify({text: args.text})
  });

  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  const body = await resp.json();
  return body.data || body;
}
