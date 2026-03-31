/* @meta
{
  "name": "xiaohongshu/note",
  "description": "获取小红书笔记详情（标题、正文、互动数据）",
  "domain": "www.xiaohongshu.com",
  "args": {
    "note_id": {"required": true, "description": "Note ID or URL"}
  },
  "capabilities": ["network"],
  "readOnly": true,
  "example": "bb-browser site xiaohongshu/note 69aa7160000000001b01634d"
}
*/

async function(args) {
  if (!args.note_id) return {error: 'Missing argument: note_id'};

  let noteId = args.note_id;
  const urlMatch = noteId.match(/explore\/([a-f0-9]+)/);
  if (urlMatch) noteId = urlMatch[1];

  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Not logged in?'};

  const noteStore = pinia._s.get('note');
  if (!noteStore) return {error: 'Note store not found', hint: 'Not logged in?'};

  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('/feed') && b?.includes?.(noteId)) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    await noteStore.getNoteDetailByNoteId(noteId);
    await new Promise(r => setTimeout(r, 500));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'Note fetch failed', hint: 'Note may be deleted or restricted'};
  const note = captured.data?.items?.[0]?.note_card;
  if (!note) return {error: 'Note not found in response'};
  return {
    note_id: noteId, title: note.title, desc: note.desc, type: note.type,
    url: 'https://www.xiaohongshu.com/explore/' + noteId,
    author: note.user?.nickname, author_id: note.user?.user_id,
    likes: note.interact_info?.liked_count, comments: note.interact_info?.comment_count,
    collects: note.interact_info?.collected_count, shares: note.interact_info?.share_count,
    tags: note.tag_list?.map(t => t.name),
    images: note.image_list?.map(img => img.info_list?.[0]?.url),
    created_time: note.time
  };
}
