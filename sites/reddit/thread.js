/* @meta
{
  "name": "reddit/thread",
  "description": "获取 Reddit 帖子的完整讨论树",
  "domain": "www.reddit.com",
  "args": {
    "url": {"required": true, "description": "Reddit post URL"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site reddit/thread https://www.reddit.com/r/LocalLLaMA/comments/1rrisqn/..."
}
*/

async function(args) {
  if (!args.url) return {error: 'Missing argument: url', hint: 'Provide a Reddit post URL'};
  let path = args.url.replace(/https?:\/\/[^/]*/, '').replace(/\?.*/, '').replace(/\/*$/, '/');
  // Normalize to /r/sub/comments/POST_ID/ — strip slug and any comment suffixes
  // Handles: .../comments/ID/slug/, .../comments/ID/comment/CID/, .../comments/ID/slug/CID/
  const m = path.match(/(\/r\/[^/]+\/comments\/[^/]+\/)/);
  if (m) path = m[1];
  const resp = await fetch(path + '.json?limit=500&depth=10&raw_json=1', {credentials: 'include'});
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const d = await resp.json();
  if (!d[0]?.data?.children?.[0]?.data) return {error: 'Unexpected response', hint: 'Post may be deleted or URL is incorrect'};
  const post = d[0].data.children[0].data;

  function flatten(children, depth) {
    let result = [];
    for (const child of children) {
      if (child.kind !== 't1') continue;
      const c = child.data;
      result.push({id: c.name, parent_id: c.parent_id, author: c.author, score: c.score, body: c.body, depth});
      if (c.replies?.data?.children)
        result = result.concat(flatten(c.replies.data.children, depth + 1));
    }
    return result;
  }

  const comments = flatten(d[1]?.data?.children || [], 0);
  return {
    post: {id: post.name, title: post.title, author: post.author, subreddit: post.subreddit_name_prefixed,
      score: post.score, num_comments: post.num_comments, selftext: post.selftext, url: post.url, created_utc: post.created_utc},
    comments_total: comments.length,
    comments
  };
}
