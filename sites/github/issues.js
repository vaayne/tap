/* @meta
{
  "name": "github/issues",
  "description": "获取 GitHub 仓库的 issue 列表",
  "domain": "github.com",
  "args": {
    "repo": {"required": true, "description": "owner/repo format"},
    "state": {"required": false, "description": "open, closed, or all (default: open)"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site github/issues epiral/bb-browser"
}
*/

async function(args) {
  if (!args.repo) return {error: 'Missing argument: repo'};
  const state = args.state || 'open';
  const resp = await fetch('https://api.github.com/repos/' + args.repo + '/issues?state=' + state + '&per_page=30', {credentials: 'include'});
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const issues = await resp.json();
  return {
    repo: args.repo, state, count: issues.length,
    issues: issues.map(i => ({
      number: i.number, title: i.title, state: i.state,
      url: i.html_url,
      author: i.user?.login, labels: i.labels?.map(l => l.name),
      comments: i.comments, created_at: i.created_at,
      is_pr: !!i.pull_request
    }))
  };
}
