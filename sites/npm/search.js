/* @meta
{
  "name": "npm/search",
  "description": "Search npm packages via registry API",
  "domain": "www.npmjs.com",
  "args": {
    "query": {"type": "string", "description": "Search query", "required": true},
    "count": {"type": "number", "description": "Number of results (default 20, max 250)", "default": 20}
  },
  "readOnly": true,
  "example": "bb-browser site npm/search \"react state management\""
}
*/

async function(args) {
  const query = args.query;
  if (!query) return {error: 'Missing required argument: query'};
  const size = Math.min(args.count || 20, 250);
  const url = `https://registry.npmjs.org/-/v1/search?text=${encodeURIComponent(query)}&size=${size}`;
  const resp = await fetch(url);
  if (!resp.ok) return {error: 'HTTP ' + resp.status};
  const data = await resp.json();
  const packages = (data.objects || []).map(obj => {
    const pkg = obj.package || {};
    const score = obj.score || {};
    return {
      name: pkg.name,
      version: pkg.version,
      description: (pkg.description || '').substring(0, 300),
      author: pkg.publisher?.username || pkg.author?.name || null,
      date: pkg.date,
      url: pkg.links?.npm || `https://www.npmjs.com/package/${pkg.name}`,
      homepage: pkg.links?.homepage || null,
      repository: pkg.links?.repository || null,
      score: Math.round((score.final || 0) * 100) / 100,
      searchScore: Math.round((obj.searchScore || 0) * 100) / 100,
      keywords: (pkg.keywords || []).slice(0, 8)
    };
  });
  return {total: data.total || packages.length, count: packages.length, packages};
}
