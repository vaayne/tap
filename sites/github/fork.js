/* @meta
{
  "name": "github/fork",
  "description": "Fork a GitHub repository",
  "domain": "github.com",
  "args": {
    "repo": {"required": true, "description": "Repository to fork (owner/repo)"}
  },
  "capabilities": ["network"],
  "readOnly": false,
  "example": "bb-browser site github/fork epiral/bb-sites"
}
*/

async function(args) {
  if (!args.repo) return {error: 'Missing argument: repo'};

  const resp = await fetch('https://api.github.com/repos/' + args.repo + '/forks', {
    method: 'POST',
    credentials: 'include',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({})
  });

  if (!resp.ok) {
    const status = resp.status;
    if (status === 401 || status === 403) return {error: 'HTTP ' + status, hint: 'Not logged in to GitHub'};
    if (status === 404) return {error: 'Repo not found: ' + args.repo};
    return {error: 'HTTP ' + status};
  }

  const fork = await resp.json();
  return {
    full_name: fork.full_name,
    url: fork.html_url,
    clone_url: fork.clone_url
  };
}
