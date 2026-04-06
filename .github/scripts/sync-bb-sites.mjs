#!/usr/bin/env node
/**
 * Parse bb-sites scripts and POST them to the tap web API batch endpoint.
 *
 * Usage: node sync-bb-sites.mjs <bb-sites-dir>
 *
 * Env:
 *   TAP_SCRIPTS_SECRET - shared secret for X-Tap-Secret header
 *   TAP_API_URL        - batch endpoint (default: https://tap.vaayne.com/api/batch)
 */

import { readdir, readFile } from "node:fs/promises"
import { join, basename, dirname } from "node:path"
import { createHash } from "node:crypto"

const sitesDir = process.argv[2]
if (!sitesDir) {
  console.error("Usage: node sync-bb-sites.mjs <bb-sites-dir>")
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
 * Discover all .js files under sitesDir organized as site/action.js.
 */
async function discoverScripts() {
  const scripts = []
  const entries = await readdir(sitesDir, { withFileTypes: true })

  for (const entry of entries) {
    if (!entry.isDirectory()) continue
    const site = entry.name
    // Skip hidden dirs and common non-script dirs
    if (site.startsWith(".") || site === "node_modules") continue

    const siteDir = join(sitesDir, site)
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
      const action = basename(file, ".js")

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
  const scripts = await discoverScripts()
  if (scripts.length === 0) {
    console.error("No scripts found")
    process.exit(1)
  }

  console.log(`Found ${scripts.length} scripts, posting to ${apiUrl}`)

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
