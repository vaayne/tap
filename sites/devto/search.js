/* @meta
{
  "name": "devto/search",
  "description": "Search Dev.to articles by keyword",
  "domain": "dev.to",
  "args": {
    "query": { "type": "string", "required": true, "description": "Search keyword or phrase" },
    "count": { "type": "number", "default": 10, "description": "Number of results (max 100)" }
  },
  "readOnly": true,
  "example": "bb-browser site devto/search \"rust programming\""
}
*/

async function(args) {
  const query = args.query;
  if (!query) return { error: 'query is required' };
  const count = Math.min(args.count || 10, 100);

  const url = 'https://dev.to/search/feed_content?per_page=' + count +
    '&page=0&search_fields=' + encodeURIComponent(query) + '&class_name=Article';

  const resp = await fetch(url);
  if (!resp.ok) return { error: 'HTTP ' + resp.status };

  const data = await resp.json();
  const articles = data.result || [];

  return {
    query: query,
    count: articles.length,
    articles: articles.map(a => ({
      title: a.title,
      url: a.path ? ('https://dev.to' + a.path) : null,
      description: (a.cloudinary_video_url ? '[video] ' : '') +
        (a.body_text || '').substring(0, 300),
      author: a.user ? a.user.name : null,
      username: a.user ? a.user.username : null,
      published_at: a.published_at_int ? new Date(a.published_at_int * 1000).toISOString() : null,
      reactions: a.public_reactions_count || 0,
      comments: a.comments_count || 0,
      tags: a.tag_list || [],
      reading_time: a.reading_time || null
    }))
  };
}
