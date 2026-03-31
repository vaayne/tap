/* @meta
{
  "name": "pypi/search",
  "description": "搜索 Python 包",
  "domain": "pypi.org",
  "args": {
    "query": "搜索关键词",
    "page": "页码（默认 1）"
  },
  "readOnly": true,
  "example": "bb-browser site pypi/search \"machine learning\""
}
*/

async function(args) {
  const query = args.query || args._text;
  if (!query) return {error: 'Missing query. Usage: bb-browser site pypi/search "QUERY"'};
  const page = args.page || 1;
  const url = `https://pypi.org/search/?q=${encodeURIComponent(query)}&page=${page}`;
  const resp = await fetch(url);
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const html = await resp.text();
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, 'text/html');
  const snippets = doc.querySelectorAll('a.package-snippet');
  const packages = Array.from(snippets).map(el => ({
    name: (el.querySelector('.package-snippet__name') || {}).textContent?.trim(),
    version: (el.querySelector('.package-snippet__version') || {}).textContent?.trim(),
    description: (el.querySelector('.package-snippet__description') || {}).textContent?.trim(),
    url: el.href ? `https://pypi.org${el.getAttribute('href')}` : undefined
  }));
  return {query, page, count: packages.length, packages};
}
