import { env } from "cloudflare:workers"
import type {
  ScriptRow,
  ScriptListItem,
  ScriptDetail,
  ScriptArg,
  SyncManifestItem,
} from "./types"

function parseArgs(argsJson: string): Record<string, ScriptArg> {
  try {
    return JSON.parse(argsJson) as Record<string, ScriptArg>
  } catch {
    return {}
  }
}

function rowToListItem(
  row: ScriptRow,
  usageCount: number,
): ScriptListItem {
  return {
    name: row.name,
    site: row.site,
    description: row.description,
    domain: row.domain,
    readOnly: row.read_only === 1,
    example: row.example,
    args: parseArgs(row.args),
    usageCount,
    updatedAt: row.updated_at,
  }
}

type ListScriptsParams = {
  query?: string
  sort?: "popular" | "newest" | "name"
  site?: string
}

export async function listScripts(
  params: ListScriptsParams,
): Promise<ScriptListItem[]> {
  const { DB } = env
  const conditions: string[] = []
  const bindings: (string | number)[] = []

  if (params.query) {
    conditions.push(
      "(s.name LIKE ? ESCAPE '\\' OR s.description LIKE ? ESCAPE '\\')",
    )
    const escaped = params.query.replace(/[%_\\]/g, "\\$&")
    const pattern = `%${escaped}%`
    bindings.push(pattern, pattern)
  }

  if (params.site) {
    conditions.push("s.site = ?")
    bindings.push(params.site)
  }

  const whereClause =
    conditions.length > 0
      ? `WHERE ${conditions.join(" AND ")}`
      : ""

  let orderClause: string
  switch (params.sort) {
    case "popular":
      orderClause = "ORDER BY usage_count DESC, s.name ASC"
      break
    case "newest":
      orderClause = "ORDER BY s.updated_at DESC"
      break
    case "name":
    default:
      orderClause = "ORDER BY s.name ASC"
      break
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
  `

  const result = await DB.prepare(sql)
    .bind(...bindings)
    .all<ScriptRow & { usage_count: number }>()

  return result.results.map((row) =>
    rowToListItem(row, row.usage_count),
  )
}

/** Format Date as SQLite-compatible datetime string (matches datetime('now') output) */
function toSqliteDatetime(date: Date): string {
  return date.toISOString().replace("T", " ").replace(/\.\d{3}Z$/, "")
}

export async function getScript(
  name: string,
): Promise<ScriptDetail | null> {
  const { DB } = env
  const d7 = toSqliteDatetime(new Date(
    Date.now() - 7 * 24 * 60 * 60 * 1000,
  ))
  const d30 = toSqliteDatetime(new Date(
    Date.now() - 30 * 24 * 60 * 60 * 1000,
  ))

  const [scriptResult, totalResult, d7Result, d30Result] =
    await DB.batch([
      DB.prepare("SELECT * FROM scripts WHERE name = ?").bind(
        name,
      ),
      DB.prepare(
        "SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ?",
      ).bind(name),
      DB.prepare(
        "SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ? AND reported_at >= ?",
      ).bind(name, d7),
      DB.prepare(
        "SELECT COUNT(*) AS cnt FROM usage_events WHERE script_name = ? AND reported_at >= ?",
      ).bind(name, d30),
    ])

  const script = scriptResult.results[0] as
    | ScriptRow
    | undefined
  if (!script) return null

  const total = (totalResult.results[0] as { cnt: number })
    ?.cnt ?? 0
  const last7d = (d7Result.results[0] as { cnt: number })
    ?.cnt ?? 0
  const last30d = (d30Result.results[0] as { cnt: number })
    ?.cnt ?? 0

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
    usage: { total, last7d, last30d },
  }
}

export async function getScriptContent(
  name: string,
): Promise<string | null> {
  const { DB } = env
  const row = await DB.prepare(
    "SELECT content FROM scripts WHERE name = ?",
  )
    .bind(name)
    .first<{ content: string }>()
  return row?.content ?? null
}

export async function getSyncManifest(): Promise<
  SyncManifestItem[]
> {
  const { DB } = env
  const result = await DB.prepare(
    "SELECT name, hash, updated_at FROM scripts ORDER BY name",
  ).all<{ name: string; hash: string; updated_at: string }>()

  return result.results.map((row) => ({
    name: row.name,
    hash: row.hash,
    updatedAt: row.updated_at,
  }))
}

export async function reportUsage(
  scriptName: string,
): Promise<void> {
  const { DB } = env
  await DB.prepare(
    "INSERT INTO usage_events (script_name) VALUES (?)",
  )
    .bind(scriptName)
    .run()
}

export async function listSites(): Promise<string[]> {
  const { DB } = env
  const result = await DB.prepare(
    "SELECT DISTINCT site FROM scripts ORDER BY site",
  ).all<{ site: string }>()
  return result.results.map((row) => row.site)
}
