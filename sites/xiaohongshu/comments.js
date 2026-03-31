/* @meta
{
  "name": "xiaohongshu/comments",
  "description": "获取小红书笔记的评论列表",
  "domain": "www.xiaohongshu.com",
  "args": {
    "note_id": {"required": true, "description": "Note ID"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site xiaohongshu/comments 69aa7160000000001b01634d"
}
*/

async function(args) {
  if (!args.note_id) return {error: 'Missing argument: note_id'};

  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Not logged in?'};

  // 先通过 note store 设置当前笔记（触发评论加载）
  const noteStore = pinia._s.get('note');
  if (!noteStore) return {error: 'Note store not found', hint: 'Not logged in?'};

  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('comment/page') && this.__url?.includes(args.note_id)) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    // 设置当前 noteId，触发评论加载
    noteStore.setCurrentNoteId(args.note_id);
    await noteStore.getNoteDetailByNoteId(args.note_id);
    await new Promise(r => setTimeout(r, 800));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'Comments fetch failed', hint: 'Not logged in?'};
  const comments = (captured.data?.comments || []).map(c => ({
    id: c.id, author: c.user_info?.nickname, author_id: c.user_info?.user_id,
    content: c.content, likes: c.like_count,
    sub_comment_count: c.sub_comment_count, created_time: c.create_time
  }));
  return {note_id: args.note_id, count: comments.length, has_more: captured.data?.has_more, comments};
}
