/* @meta
{
  "name": "reddit/search",
  "description": "搜索 Reddit 帖子",
  "domain": "www.reddit.com",
  "args": {
    "query": {"required": true, "description": "Search query"},
    "subreddit": {"required": false, "description": "Limit search to a subreddit (without r/ prefix)"},
    "sort": {"required": false, "description": "Sort order: relevance (default), hot, top, new, comments"},
    "time": {"required": false, "description": "Time filter: all (default), hour, day, week, month, year"},
    "count": {"required": false, "description": "Number of results (default 25, max 100)"},
    "after": {"required": false, "description": "Fullname of the last item for pagination"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site reddit/search \"claude code\" --sort top --time week"
}
*/

async function(args) {
  if (!args.query) return {error: 'Missing argument: query', hint: 'Provide a search query'};

  const count = Math.min(parseInt(args.count) || 25, 100);
  const sort = args.sort || 'relevance';
  const time = args.time || 'all';

  let url;
  if (args.subreddit) {
    url = '/r/' + args.subreddit + '/search.json?restrict_sr=on&q=' + encodeURIComponent(args.query);
  } else {
    url = '/search.json?q=' + encodeURIComponent(args.query);
  }
  url += '&sort=' + sort + '&t=' + time + '&limit=' + count + '&raw_json=1';
  if (args.after) url += '&after=' + args.after;

  const resp = await fetch(url, {credentials: 'include'});
  if (!resp.ok) return {error: 'HTTP ' + resp.status, hint: resp.status === 404 ? 'Subreddit not found' : 'API error'};
  const d = await resp.json();
  if (!d.data?.children) return {error: 'Unexpected response', hint: 'Reddit may be rate-limiting or returning a login page'};

  const posts = d.data.children.map(c => ({
    id: c.data.name,
    title: c.data.title,
    author: c.data.author,
    subreddit: c.data.subreddit_name_prefixed,
    score: c.data.score,
    num_comments: c.data.num_comments,
    created_utc: c.data.created_utc,
    url: c.data.url,
    permalink: 'https://www.reddit.com' + c.data.permalink,
    selftext_preview: (c.data.selftext || '').substring(0, 200),
    is_self: c.data.is_self,
    link_flair_text: c.data.link_flair_text || null
  }));

  return {
    query: args.query,
    subreddit: args.subreddit || null,
    sort,
    time,
    count: posts.length,
    after: d.data.after || null,
    posts
  };
}
