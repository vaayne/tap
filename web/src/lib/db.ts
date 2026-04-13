import { env } from "cloudflare:workers"
import type {
  ScriptRow,
  ScriptListItem,
  ScriptDetail,
  ScriptArg,
  SyncManifestItem,
  BatchScript,
} from "./types"

function parseArgs(argsJson: string): Record<string, ScriptArg> {
  try {
    return JSON.parse(argsJson) as Record<string, ScriptArg>
  } catch {
    return {}
  }
}

function parseCapabilities(json: string): string[] {
  try {
    const parsed = JSON.parse(json)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
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
    capabilities: parseCapabilities(row.capabilities),
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
    capabilities: parseCapabilities(script.capabilities),
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

/**
 * Replaces all scripts in D1 with the provided list.
 * Uses a DELETE-then-INSERT pattern batched through D1's batch API.
 */
export async function batchUpdate(
  scripts: BatchScript[],
): Promise<void> {
  const { DB } = env

  const scriptsTableInfo = await DB.prepare(
    "PRAGMA table_info(scripts)",
  ).all<{ name: string }>()
  const hasCapabilitiesColumn = scriptsTableInfo.results.some(
    (column) => column.name === "capabilities",
  )

  // D1 batch supports up to 100 statements; chunk inserts to stay safe.
  // To guarantee atomicity we insert into a temp table first, then
  // swap via DELETE + INSERT … SELECT in a single batch call.
  const CHUNK_SIZE = 95

  const dropTemp = DB.prepare("DROP TABLE IF EXISTS _scripts_staging")
  const createTemp = DB.prepare(
    `CREATE TABLE _scripts_staging (
      name TEXT PRIMARY KEY, site TEXT NOT NULL, description TEXT NOT NULL DEFAULT '',
      domain TEXT NOT NULL DEFAULT '', args TEXT NOT NULL DEFAULT '{}',
      read_only INTEGER NOT NULL DEFAULT 1, example TEXT NOT NULL DEFAULT '',
      capabilities TEXT NOT NULL DEFAULT '[]',
      content TEXT NOT NULL, hash TEXT NOT NULL,
      created_at TEXT NOT NULL DEFAULT (datetime('now')),
      updated_at TEXT NOT NULL DEFAULT (datetime('now'))
    )`,
  )

  // Phase 1: recreate staging table to avoid stale schema from prior runs
  await DB.batch([dropTemp, createTemp])

  const insertStmts = scripts.map((s) =>
    DB.prepare(
      `INSERT INTO _scripts_staging (name, site, description, domain, args, read_only, example, capabilities, content, hash)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    ).bind(
      s.name,
      s.site,
      s.description,
      s.domain,
      s.args,
      s.readOnly ? 1 : 0,
      s.example,
      JSON.stringify(s.capabilities ?? []),
      s.content,
      s.hash,
    ),
  )

  for (let i = 0; i < insertStmts.length; i += CHUNK_SIZE) {
    const chunk = insertStmts.slice(i, i + CHUNK_SIZE)
    await DB.batch(chunk)
  }

  const swapScripts = hasCapabilitiesColumn
    ? DB.prepare(
        `INSERT INTO scripts (name, site, description, domain, args, read_only, example, capabilities, content, hash, created_at, updated_at)
         SELECT name, site, description, domain, args, read_only, example, capabilities, content, hash, created_at, updated_at
         FROM _scripts_staging`,
      )
    : DB.prepare(
        `INSERT INTO scripts (name, site, description, domain, args, read_only, example, content, hash, created_at, updated_at)
         SELECT name, site, description, domain, args, read_only, example, content, hash, created_at, updated_at
         FROM _scripts_staging`,
      )

  // Phase 2: atomic swap — single batch so scripts is never empty
  await DB.batch([
    DB.prepare("DELETE FROM scripts"),
    swapScripts,
    DB.prepare("DROP TABLE _scripts_staging"),
  ])
}
