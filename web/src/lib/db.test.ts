import { beforeEach, describe, expect, it, vi } from "vitest"

type MockStmt = {
  sql: string
  bindings: unknown[]
  bind: (...args: unknown[]) => MockStmt
}

type MockDB = {
  prepare: (sql: string) => MockStmt
  batch: (stmts: MockStmt[]) => Promise<{ results: unknown[] }[]>
}

const batchCalls: MockStmt[][] = []

const DB: MockDB = {
  prepare(sql: string) {
    return {
      sql,
      bindings: [],
      bind(...args: unknown[]) {
        this.bindings = args
        return this
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

describe("batchUpdate", () => {
  beforeEach(() => {
    batchCalls.length = 0
  })

  it("recreates the staging table before loading scripts", async () => {
    const { batchUpdate } = await import("./db")

    await batchUpdate([
      {
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
      },
    ])

    expect(batchCalls[0]?.map((stmt) => stmt.sql)).toEqual([
      "DROP TABLE IF EXISTS _scripts_staging",
      expect.stringContaining("CREATE TABLE _scripts_staging"),
    ])

    expect(batchCalls.at(-1)?.map((stmt) => stmt.sql)).toContain(
      "DROP TABLE _scripts_staging",
    )
  })
})
