/* @meta
{
  "name": "hackernews/top",
  "description": "获取 Hacker News 当前热门帖子",
  "domain": "news.ycombinator.com",
  "args": {
    "count": {"required": false, "description": "Number of posts (default: 20, max: 50)"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site hackernews/top 10"
}
*/

async function(args) {
  const count = Math.min(parseInt(args.count) || 20, 50);
  // HN Firebase API is public, no auth needed
  const idsResp = await fetch('https://hacker-news.firebaseio.com/v0/topstories.json');
  if (!idsResp.ok) return {error: 'HTTP ' + idsResp.status};
  const ids = await idsResp.json();
  const topIds = ids.slice(0, count);

  const items = await Promise.all(topIds.map(async id => {
    const resp = await fetch('https://hacker-news.firebaseio.com/v0/item/' + id + '.json');
    return await resp.json();
  }));

  return {
    count: items.length,
    posts: items.map((item, i) => ({
      rank: i + 1, id: item.id, title: item.title,
      url: item.url || null,
      hn_url: 'https://news.ycombinator.com/item?id=' + item.id,
      author: item.by, score: item.score,
      comments: item.descendants || 0, time: item.time
    }))
  };
}
