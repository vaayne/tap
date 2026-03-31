/* @meta
{
  "name": "xiaohongshu/feed",
  "description": "获取小红书首页推荐 Feed",
  "domain": "www.xiaohongshu.com",
  "args": {},
  "capabilities": ["network"],
  "readOnly": true
}
*/

async function(args) {
  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Ensure xiaohongshu.com is fully loaded'};

  const feedStore = pinia._s.get('feed');
  if (!feedStore) return {error: 'Feed store not found', hint: 'Not logged in?'};

  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('homefeed') && !this.__url?.includes('category')) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    await feedStore.fetchFeeds();
    await new Promise(r => setTimeout(r, 500));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'Feed fetch failed', hint: 'Not logged in?'};
  const notes = (captured.data.items || []).map(item => ({
    id: item.id, xsec_token: item.xsec_token,
    title: item.note_card?.display_title, type: item.note_card?.type,
    url: 'https://www.xiaohongshu.com/explore/' + item.id,
    author: item.note_card?.user?.nickname, author_id: item.note_card?.user?.user_id,
    likes: item.note_card?.interact_info?.liked_count
  }));
  return {count: notes.length, has_more: !!captured.data.cursor_score, notes};
}
