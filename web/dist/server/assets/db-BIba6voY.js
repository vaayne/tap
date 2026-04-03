import { env } from "cloudflare:workers";
//#region src/lib/db.ts
function parseArgs(argsJson) {
	try {
		return JSON.parse(argsJson);
	} catch {
		return {};
	}
}
function rowToListItem(row, usageCount) {
	return {
		name: row.name,
		site: row.site,
		description: row.description,
		domain: row.domain,
		readOnly: row.read_only === 1,
		example: row.example,
		args: parseArgs(row.args),
		usageCount,
		updatedAt: row.updated_at
	};
}
async function listScripts(params) {
	const { DB } = env;
	const conditions = [];
	const bindings = [];
	if (params.query) {
		conditions.push("(s.name LIKE ? ESCAPE '\\' OR s.description LIKE ? ESCAPE '\\')");
		const pattern = `%${params.query.replace(/[%_\\]/g, "\\$&")}%`;
		bindings.push(pattern, pattern);
	}
	if (params.site) {
		conditions.push("s.site = ?");
		bindings.push(params.site);
	}
	const whereClause = conditions.length > 0 ? `WHERE ${conditions.join(" AND ")}` : "";
	let orderClause;
	switch (params.sort) {
		case "popular":
			orderClause = "ORDER BY usage_count DESC, s.name ASC";
			break;
		case "newest":
			orderClause = "ORDER BY s.updated_at DESC";
			break;
		default:
			orderClause = "ORDER BY s.name ASC";
			break;
	}
	const sql = `
    SELECT s.*, COALESCE(u.cnt, 0) AS usage_count
    FROM scripts s
    LEFT JOIN (
      SELECT script_name, COUNT(*) AS cnt
      FROM usage_events
      GROUP BY script_name
    ) u ON u.script_name = s.name
    ${whereClause}
    ${orderClause}
  `;
	return (await DB.prepare(sql).bind(...bindings).all()).results.map((row) => rowToListItem(row, row.usage_count));
}
/** Format Date as SQLite-compatible datetime string (matches datetime('now') output) */
function toSqliteDatetime(date) {
	return date.toISOString().replace("T", " ").replace(/\.\d{3}Z$/, "");
}
async function getScript(name) {
	const { DB } = env;
	const d7 = toSqliteDatetime(/* @__PURE__ */ new Date(Date.now() - 10080 * 60 * 1e3));
	const d30 = toSqliteDatetime(/* @__PURE__ */ new Date(Date.now() - 720 * 60 * 60 * 1e3));
	const [scriptResult, totalResult, d7Result, d30Result] = await DB.batch([
		DB.prepare("SELECT * FROM scripts WHERE name = ?").bind(name),
		DB.prepare("SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ?").bind(name),
		DB.prepare("SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ? AND reported_at >= ?").bind(name, d7),
		DB.prepare("SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ? AND reported_at >= ?").bind(name, d30)
	]);
	const script = scriptResult.results[0];
	if (!script) return null;
	const total = totalResult.results[0]?.cnt ?? 0;
	const last7d = d7Result.results[0]?.cnt ?? 0;
	const last30d = d30Result.results[0]?.cnt ?? 0;
	return {
		name: script.name,
		site: script.site,
		description: script.description,
		domain: script.domain,
		readOnly: script.read_only === 1,
		example: script.example,
		args: parseArgs(script.args),
		content: script.content,
		hash: script.hash,
		usageCount: total,
		createdAt: script.created_at,
		updatedAt: script.updated_at,
		usage: {
			total,
			last7d,
			last30d
		}
	};
}
async function getScriptContent(name) {
	const { DB } = env;
	return (await DB.prepare("SELECT content FROM scripts WHERE name = ?").bind(name).first())?.content ?? null;
}
async function getSyncManifest() {
	const { DB } = env;
	return (await DB.prepare("SELECT name, hash, updated_at FROM scripts ORDER BY name").all()).results.map((row) => ({
		name: row.name,
		hash: row.hash,
		updatedAt: row.updated_at
	}));
}
async function reportUsage(scriptName) {
	const { DB } = env;
	await DB.prepare("INSERT INTO usage_events (script_name) VALUES (?)").bind(scriptName).run();
}
async function listSites() {
	const { DB } = env;
	return (await DB.prepare("SELECT DISTINCT site FROM scripts ORDER BY site").all()).results.map((row) => row.site);
}
//#endregion
export { listSites as a, listScripts as i, getScriptContent as n, reportUsage as o, getSyncManifest as r, getScript as t };
