#!/usr/bin/env node
/**
 * Parse site scripts and POST them to the tap web API batch endpoint.
 * Accepts one or more directories; later directories override earlier ones
 * when script names collide (matching the CLI registry priority).
 *
 * Usage: node sync-bb-sites.mjs <dir> [<dir> ...]
 *
 * Env:
 *   TAP_SCRIPTS_SECRET - shared secret for X-Tap-Secret header
 *   TAP_API_URL        - batch endpoint (default: https://tap.vaayne.com/api/batch)
 */

import { readdir, readFile } from "node:fs/promises"
import { join } from "node:path"
import { createHash } from "node:crypto"

const sitesDirs = process.argv.slice(2)
if (sitesDirs.length === 0) {
  console.error("Usage: node sync-bb-sites.mjs <dir> [<dir> ...]")
  process.exit(1)
}

const secret = process.env.TAP_SCRIPTS_SECRET
if (!secret) {
  console.error("TAP_SCRIPTS_SECRET is required")
  process.exit(1)
}

const apiUrl =
  process.env.TAP_API_URL || "https://tap.vaayne.com/api/batch"

/**
 * Parse the /* @meta ... * / block from a script file.
 */
function parseMeta(content) {
  const match = content.match(/\/\*\s*@meta\s*\n([\s\S]*?)\*\//)
  if (!match) return null
  try {
    return JSON.parse(match[1])
  } catch {
    return null
  }
}

/**
 * Discover all .js files under dir organized as site/action.js.
 */
async function discoverScripts(dir) {
  const scripts = []
  let entries
  try {
    entries = await readdir(dir, { withFileTypes: true })
  } catch {
    return scripts
  }

  for (const entry of entries) {
    if (!entry.isDirectory()) continue
    const site = entry.name
    if (site.startsWith(".") || site === "node_modules") continue

    const siteDir = join(dir, site)
    const files = await readdir(siteDir)

    for (const file of files) {
      if (!file.endsWith(".js")) continue
      const filePath = join(siteDir, file)
      const content = await readFile(filePath, "utf-8")
      const meta = parseMeta(content)
      if (!meta || !meta.name) {
        console.warn(`Skipping ${site}/${file}: no valid @meta block`)
        continue
      }

      const hash = createHash("sha256").update(content).digest("hex")

      scripts.push({
        name: meta.name,
        site,
        content,
        hash,
        description: meta.description || "",
        domain: meta.domain || "",
        args: JSON.stringify(meta.args || {}),
        capabilities: meta.capabilities || [],
        example: meta.example || "",
        readOnly: meta.readOnly ?? true,
      })
    }
  }

  return scripts
}

async function main() {
  const byName = new Map()
  for (const dir of sitesDirs) {
    const found = await discoverScripts(dir)
    for (const s of found) {
      byName.set(s.name, s)
    }
    console.log(`${dir}: ${found.length} scripts`)
  }
  const scripts = [...byName.values()]
  if (scripts.length === 0) {
    console.error("No scripts found")
    process.exit(1)
  }

  console.log(`Total: ${scripts.length} scripts, posting to ${apiUrl}`)

  const resp = await fetch(apiUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Tap-Secret": secret,
    },
    body: JSON.stringify({ scripts }),
  })

  const body = await resp.text()
  if (!resp.ok) {
    console.error(`Batch update failed: HTTP ${resp.status}`)
    console.error(body)
    process.exit(1)
  }

  console.log(`Success: ${body}`)
}

main()
