/* @meta
{
  "name": "weixin/article",
  "description": "Fetch WeChat MP article by URL — title, author, content, images (HTTP-only, no browser)",
  "domain": "mp.weixin.qq.com",
  "args": {
    "url": {"required": true, "description": "WeChat article URL or short ID, e.g. https://mp.weixin.qq.com/s/xxxxx"}
  },
  "runtime": "http",
  "readOnly": true,
  "example": "tap site weixin/article url=https://mp.weixin.qq.com/s/DJtp4QUJtJHCnNU9yarsww"
}
*/

async function(args) {
  if (!args.url) {
    return {error: 'Missing argument: url', hint: 'Provide a WeChat article URL'};
  }

  let url = args.url.trim();
  if (!url.includes('://') && !url.startsWith('//')) {
    url = 'https://mp.weixin.qq.com/s/' + url;
  } else if (url.startsWith('//')) {
    url = 'https:' + url;
  }

  if (!url.includes('mp.weixin.qq.com')) {
    return {error: 'Invalid URL', hint: 'URL must be from mp.weixin.qq.com'};
  }

  // WeChat Mobile UA bypasses access restrictions without needing a browser
  const resp = await fetch(url, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Linux; U; Android 4.1.2; zh-cn; MI ONE Plus Build/JZO54K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.0.1.352',
      'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
      'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8'
    }
  });

  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  const html = await resp.text();

  if (html.includes('未知错误') && html.includes('暂无权限查看此页面内容')) {
    return {error: 'Access denied', hint: 'WeChat blocked this request. Try a different URL.'};
  }

  // The page is ~3MB of JS bundles. Work only on small relevant slices to avoid
  // running expensive operations on the full string.

  // Metadata variables live in a small <script> block containing 'var msg_title'.
  // Extract just that block (~20KB) for regex operations.
  const metaAnchor = 'var msg_title';
  const metaPos = html.indexOf(metaAnchor);
  const metaStart = metaPos !== -1 ? html.lastIndexOf('<script', metaPos) : -1;
  const metaEnd = metaPos !== -1 ? html.indexOf('</script>', metaPos) : -1;
  const metaBlock = (metaStart !== -1 && metaEnd !== -1)
    ? html.slice(metaStart, metaEnd)
    : (metaPos !== -1 ? html.slice(Math.max(0, metaPos - 500), metaPos + 10000) : '');

  // Article body lives in the div with id="js_content". Slice 300KB around it.
  const contentAnchor = 'id="js_content"';
  const contentIdx = html.indexOf(contentAnchor);
  const contentBlock = contentIdx !== -1
    ? html.slice(contentIdx, contentIdx + 300000)
    : '';

  function decodeEntities(str) {
    if (!str) return str;
    return str
      .replace(/&quot;/g, '"')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&amp;/g, '&')
      .replace(/&nbsp;/g, ' ')
      .replace(/&#39;/g, "'");
  }

  function extractStr(name) {
    let m = metaBlock.match(new RegExp('var\\s+' + name + '\\s*=\\s*htmlDecode\\(["\']([^"\']+)["\']\\)\\s*;'));
    if (m) return decodeEntities(m[1]);
    m = metaBlock.match(new RegExp('var\\s+' + name + '\\s*=\\s*["\']([^"\']+)["\']\\.html\\(false\\)\\s*;'));
    if (m) return decodeEntities(m[1]);
    m = metaBlock.match(new RegExp('var\\s+' + name + '\\s*=\\s*["\']([^"\']+)["\'][^;]*;'));
    if (m) return decodeEntities(m[1]);
    return null;
  }

  function extractNum(name) {
    const m = metaBlock.match(new RegExp('var\\s+' + name + '\\s*=\\s*["\']?([0-9]+)["\']?[^;]*;'));
    if (m) return parseInt(m[1], 10);
    return null;
  }

  const title = extractStr('msg_title');
  const author = extractStr('nickname');
  const description = extractStr('msg_desc');
  const coverImage = extractStr('msg_cdn_url');
  const sourceUrl = extractStr('msg_source_url');
  const canonicalUrl = extractStr('msg_link');
  const ct = extractNum('ct');
  const mid = extractStr('mid');
  const idx = extractStr('idx');
  const sn = extractStr('sn');
  const biz = extractStr('biz');

  // Extract content using indexOf (C-speed) instead of char-by-char scanning.
  function extractContent(block) {
    const marker = 'id="js_content"';
    const startIdx = block.indexOf(marker);
    if (startIdx === -1) return '';
    const tagEnd = block.indexOf('>', startIdx);
    if (tagEnd === -1) return '';
    const contentStart = tagEnd + 1;
    let depth = 1;
    let pos = contentStart;
    while (depth > 0) {
      const openIdx = block.indexOf('<div', pos);
      const closeIdx = block.indexOf('</div>', pos);
      if (closeIdx === -1) break;
      if (openIdx !== -1 && openIdx < closeIdx) {
        depth++;
        pos = openIdx + 4;
      } else {
        depth--;
        if (depth === 0) return block.slice(contentStart, closeIdx);
        pos = closeIdx + 6;
      }
    }
    return '';
  }

  const contentHtml = extractContent(contentBlock);

  let content = '';
  if (contentHtml) {
    content = contentHtml
      .replace(/<br\s*\/?>/gi, '\n')
      .replace(/<\/p>/gi, '\n')
      .replace(/<p[^>]*>/gi, '')
      .replace(/<section[^>]*>/gi, '\n')
      .replace(/<\/section>/gi, '')
      .replace(/<span[^>]*>/gi, '')
      .replace(/<\/span>/gi, '')
      .replace(/<strong[^>]*>/gi, '**')
      .replace(/<\/strong>/gi, '**')
      .replace(/<b[^>]*>/gi, '**')
      .replace(/<\/b>/gi, '**')
      .replace(/<em[^>]*>/gi, '*')
      .replace(/<\/em>/gi, '*')
      .replace(/<i[^>]*>/gi, '*')
      .replace(/<\/i>/gi, '*')
      .replace(/<h[1-6][^>]*>/gi, '\n')
      .replace(/<\/h[1-6]>/gi, '\n')
      .replace(/<img[^>]*>/gi, '')
      .replace(/<[^>]+>/g, '')
      .replace(/&nbsp;/g, ' ')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&amp;/g, '&')
      .replace(/&quot;/g, '"')
      .replace(/&#39;/g, "'")
      .replace(/\n{3,}/g, '\n\n')
      .trim();
  }

  const images = [];
  const imgRe = /<img[^>]+(?:data-src|src)=["']([^"']+)["']/g;
  let imgMatch;
  while ((imgMatch = imgRe.exec(contentHtml)) !== null) {
    const imgUrl = imgMatch[1];
    if (imgUrl && !imgUrl.startsWith('data:') && images.indexOf(imgUrl) === -1) {
      images.push(imgUrl);
    }
  }

  return {
    title,
    author,
    description,
    publish_time: ct ? new Date(ct * 1000).toISOString() : null,
    publish_timestamp: ct,
    cover_image: coverImage,
    source_url: sourceUrl,
    canonical_url: canonicalUrl,
    url,
    content,
    images: images.slice(0, 30),
    meta: {biz, mid, idx, sn}
  };
}
