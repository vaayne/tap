/* @meta
{
  "name": "xiaohongshu/user_posts",
  "description": "获取小红书用户的笔记列表",
  "domain": "www.xiaohongshu.com",
  "args": {
    "user_id": {"required": true, "description": "User ID"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site xiaohongshu/user_posts 5a927d8411be10720ae9e1e4"
}
*/

async function(args) {
  if (!args.user_id) return {error: 'Missing argument: user_id'};

  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Not logged in?'};

  const userStore = pinia._s.get('user');
  if (!userStore) return {error: 'User store not found', hint: 'Not logged in?'};

  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('user_posted') && this.__url?.includes(args.user_id)) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    await userStore.fetchNotes({userId: args.user_id});
    await new Promise(r => setTimeout(r, 500));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'User posts fetch failed', hint: 'Not logged in?'};
  const notes = (captured.data?.notes || []).map(n => ({
    note_id: n.note_id, title: n.display_title, type: n.type,
    url: 'https://www.xiaohongshu.com/explore/' + n.note_id,
    likes: n.interact_info?.liked_count, time: n.last_update_time
  }));
  return {user_id: args.user_id, count: notes.length, has_more: captured.data?.has_more, notes};
}
