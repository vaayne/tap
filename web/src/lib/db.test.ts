import { beforeEach, describe, expect, it, vi } from "vitest"

type MockStmt = {
  sql: string
  bindings: unknown[]
  bind: (...args: unknown[]) => MockStmt
  all: <T>() => Promise<{ results: T[] }>
}

type MockDB = {
  prepare: (sql: string) => MockStmt
  batch: (stmts: MockStmt[]) => Promise<{ results: unknown[] }[]>
}

const batchCalls: MockStmt[][] = []
let scriptsTableColumns = [
  "name",
  "site",
  "description",
  "domain",
  "args",
  "read_only",
  "example",
  "capabilities",
  "content",
  "hash",
  "created_at",
  "updated_at",
]

const DB: MockDB = {
  prepare(sql: string) {
    return {
      sql,
      bindings: [],
      bind(...args: unknown[]) {
        this.bindings = args
        return this
      },
      async all<T>() {
        if (sql === "PRAGMA table_info(scripts)") {
          return {
            results: scriptsTableColumns.map((name) => ({ name })) as T[],
          }
        }

        return { results: [] as T[] }
      },
    }
  },
  async batch(stmts: MockStmt[]) {
    batchCalls.push(stmts)
    return stmts.map(() => ({ results: [] }))
  },
}

vi.mock("cloudflare:workers", () => ({
  env: { DB },
}))

const exampleScript = {
  name: "example",
  site: "demo",
  content: "console.log('ok')",
  hash: "abc123",
  description: "Example",
  domain: "example.com",
  args: "{}",
  capabilities: ["network"],
  example: "tap site example",
  readOnly: true,
}

describe("batchUpdate", () => {
  beforeEach(() => {
    batchCalls.length = 0
    scriptsTableColumns = [
      "name",
      "site",
      "description",
      "domain",
      "args",
      "read_only",
      "example",
      "capabilities",
      "content",
      "hash",
      "created_at",
      "updated_at",
    ]
  })

  it("recreates the staging table before loading scripts", async () => {
    const { batchUpdate } = await import("./db")

    await batchUpdate([exampleScript])

    expect(batchCalls[0]?.map((stmt) => stmt.sql)).toEqual([
      "DROP TABLE IF EXISTS _scripts_staging",
      expect.stringContaining("CREATE TABLE _scripts_staging"),
    ])

    expect(batchCalls.at(-1)?.map((stmt) => stmt.sql)).toContain(
      "DROP TABLE _scripts_staging",
    )
  })

  it("falls back when the remote scripts table has not been migrated", async () => {
    scriptsTableColumns = [
      "name",
      "site",
      "description",
      "domain",
      "args",
      "read_only",
      "example",
      "content",
      "hash",
      "created_at",
      "updated_at",
    ]

    const { batchUpdate } = await import("./db")

    await batchUpdate([exampleScript])

    const finalBatchSql = batchCalls.at(-1)?.map((stmt) => stmt.sql) ?? []
    const insertSql = finalBatchSql.find((sql) =>
      sql.startsWith("INSERT INTO scripts"),
    )

    expect(insertSql).toBeDefined()
    expect(insertSql).not.toContain("capabilities")
  })
})
