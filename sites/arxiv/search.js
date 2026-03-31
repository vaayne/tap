/* @meta
{
  "name": "arxiv/search",
  "description": "Search arXiv papers by query",
  "domain": "arxiv.org",
  "args": {
    "query": {"type": "string", "required": true, "description": "Search query"},
    "count": {"type": "number", "default": 10, "description": "Number of results (max 50)"}
  },
  "readOnly": true,
  "example": "bb-browser site arxiv/search \"large language model\""
}
*/

async function(args) {
  const query = args.query;
  if (!query) return {error: 'query is required'};
  const count = Math.min(args.count || 10, 50);

  const url = `https://export.arxiv.org/api/query?search_query=all:${encodeURIComponent(query)}&start=0&max_results=${count}`;
  const resp = await fetch(url);
  if (!resp.ok) return {error: 'HTTP ' + resp.status};

  const xml = await resp.text();
  const parser = new DOMParser();
  const doc = parser.parseFromString(xml, 'text/xml');

  const entries = doc.querySelectorAll('entry');
  const papers = Array.from(entries).map(entry => {
    const getText = (tag) => {
      const el = entry.querySelector(tag);
      return el ? el.textContent.trim() : '';
    };

    const authors = Array.from(entry.querySelectorAll('author > name'))
      .map(n => n.textContent.trim());

    const linkEl = entry.querySelector('link[title="pdf"]');
    const pdfLink = linkEl ? linkEl.getAttribute('href') : '';

    const absLink = entry.querySelector('id')?.textContent.trim() || '';

    const categories = Array.from(entry.querySelectorAll('category'))
      .map(c => c.getAttribute('term'))
      .filter(Boolean);

    const arxivId = absLink.replace('http://arxiv.org/abs/', '');

    return {
      id: arxivId,
      title: getText('title').replace(/\s+/g, ' '),
      abstract: getText('summary').replace(/\s+/g, ' ').substring(0, 500),
      authors: authors,
      published: getText('published').substring(0, 10),
      categories: categories,
      url: absLink,
      pdf: pdfLink
    };
  });

  const totalResults = doc.querySelector('totalResults')?.textContent || '0';

  return {
    query: query,
    totalResults: parseInt(totalResults),
    count: papers.length,
    papers: papers
  };
}
