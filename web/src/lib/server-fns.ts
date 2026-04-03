import { createServerFn } from "@tanstack/react-start"
import { listScripts, listSites, getScript } from "./db"
import type { ScriptListItem, ScriptDetail } from "./types"

export const fetchScriptList = createServerFn({ method: "GET" })
  .inputValidator(
    (input: { query?: string; sort?: string; site?: string }) => input,
  )
  .handler(async ({ data }) => {
    const sort = (data.sort as "popular" | "newest" | "name") || "popular"
    const [scripts, sites] = await Promise.all([
      listScripts({
        query: data.query || undefined,
        sort,
        site: data.site || undefined,
      }),
      listSites(),
    ])
    return { scripts, sites } as {
      scripts: ScriptListItem[]
      sites: string[]
    }
  })

export const fetchPopularScripts = createServerFn({ method: "GET" })
  .inputValidator((input: { limit?: number }) => input)
  .handler(async ({ data }) => {
    const limit = data.limit ?? 6
    const [scripts, sites] = await Promise.all([
      listScripts({ sort: "popular" }),
      listSites(),
    ])
    return {
      scripts: scripts.slice(0, limit),
      totalScripts: scripts.length,
      totalSites: sites.length,
    }
  })

export const fetchScriptDetail = createServerFn({ method: "GET" })
  .inputValidator((input: { name: string }) => input)
  .handler(async ({ data }): Promise<ScriptDetail | null> => {
    return getScript(data.name)
  })
