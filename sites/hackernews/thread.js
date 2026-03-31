/* @meta
{
  "name": "hackernews/thread",
  "description": "获取 Hacker News 帖子的评论树",
  "domain": "news.ycombinator.com",
  "args": {
    "id": {"required": true, "description": "HN item ID or URL"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site hackernews/thread 12345678"
}
*/

async function(args) {
  if (!args.id) return {error: 'Missing argument: id'};
  let itemId = args.id;
  const urlMatch = itemId.match(/id=(\d+)/);
  if (urlMatch) itemId = urlMatch[1];

  const resp = await fetch('https://hacker-news.firebaseio.com/v0/item/' + itemId + '.json');
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const item = await resp.json();
  if (!item) return {error: 'Item not found', hint: 'Check the ID: ' + itemId};

  // Fetch comment tree (max 2 levels deep for performance)
  async function fetchComments(ids, depth) {
    if (!ids || ids.length === 0 || depth > 2) return [];
    const comments = await Promise.all(ids.slice(0, 30).map(async id => {
      const r = await fetch('https://hacker-news.firebaseio.com/v0/item/' + id + '.json');
      const c = await r.json();
      if (!c || c.deleted || c.dead) return null;
      return {
        id: c.id, author: c.by, text: c.text, time: c.time, depth,
        replies: await fetchComments(c.kids, depth + 1)
      };
    }));
    return comments.filter(Boolean);
  }

  const comments = await fetchComments(item.kids, 0);
  return {
    post: {id: item.id, title: item.title, url: item.url || null, hn_url: 'https://news.ycombinator.com/item?id=' + item.id, author: item.by, score: item.score, comments_count: item.descendants || 0, time: item.time, text: item.text},
    comments
  };
}
