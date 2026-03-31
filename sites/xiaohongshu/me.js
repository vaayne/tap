/* @meta
{
  "name": "xiaohongshu/me",
  "description": "获取当前小红书登录用户信息",
  "domain": "www.xiaohongshu.com",
  "args": {},
  "capabilities": ["network"],
  "readOnly": true
}
*/

async function(args) {
  // 通过 pinia store 的 user action 获取，走页面完整签名链路
  const app = document.querySelector('#app')?.__vue_app__;
  const pinia = app?.config?.globalProperties?.$pinia;
  if (!pinia?._s) return {error: 'Page not ready', hint: 'Ensure xiaohongshu.com is fully loaded'};

  // 拦截 user/me 的 response
  let captured = null;
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;
  XMLHttpRequest.prototype.open = function(m, u) { this.__url = u; return origOpen.apply(this, arguments); };
  XMLHttpRequest.prototype.send = function(b) {
    if (this.__url?.includes('/user/me')) {
      const x = this;
      const orig = x.onreadystatechange;
      x.onreadystatechange = function() { if (x.readyState === 4 && !captured) { try { captured = JSON.parse(x.responseText); } catch {} } if (orig) orig.apply(this, arguments); };
    }
    return origSend.apply(this, arguments);
  };

  try {
    const userStore = pinia._s.get('user');
    if (userStore?.getUserInfo) await userStore.getUserInfo();
    else {
      // fallback: 直接触发 user/me 请求
      const feedStore = pinia._s.get('feed');
      if (feedStore) await feedStore.fetchFeeds();
    }
    await new Promise(r => setTimeout(r, 500));
  } finally {
    XMLHttpRequest.prototype.open = origOpen;
    XMLHttpRequest.prototype.send = origSend;
  }

  if (!captured?.success) return {error: captured?.msg || 'Failed to get user info', hint: 'Not logged in?'};
  const u = captured.data;
  return {nickname: u.nickname, red_id: u.red_id, desc: u.desc, gender: u.gender, userid: u.user_id, url: 'https://www.xiaohongshu.com/user/profile/' + u.user_id};
}
