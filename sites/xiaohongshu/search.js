/* @meta
{
  "name": "xiaohongshu/search",
  "description": "搜索小红书笔记",
  "domain": "www.xiaohongshu.com",
  "args": {
    "keyword": {"required": true, "description": "Search keyword"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site xiaohongshu/search 美食"
}
*/

async function(args) {
  if (!args.keyword) return {error: 'Missing argument: keyword'};

  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Not logged in?'};

  const searchStore = pinia._s.get('search');
  if (!searchStore) return {error: 'Search store not found', hint: 'Not logged in?'};

  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('search/notes')) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    searchStore.mutateSearchValue(args.keyword);
    await searchStore.loadMore();
    await new Promise(r => setTimeout(r, 500));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'Search failed', hint: 'Not logged in?'};
  const notes = (captured.data?.items || []).map(i => ({
    id: i.id, xsec_token: i.xsec_token,
    title: i.note_card?.display_title, type: i.note_card?.type,
    url: 'https://www.xiaohongshu.com/explore/' + i.id,
    author: i.note_card?.user?.nickname,
    likes: i.note_card?.interact_info?.liked_count
  }));
  return {keyword: args.keyword, count: notes.length, has_more: captured.data?.has_more, notes};
}
