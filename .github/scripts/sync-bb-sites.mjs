#!/usr/bin/env node
/**
 * Import bb-sites plus optional higher-priority local directories and POST the
 * resulting catalog to Tap. Tap-specific metadata normalization is applied
 * only to the first (bb-sites) directory; upstream source files stay intact.
 *
 * Usage: node sync-bb-sites.mjs [--compat <file>] <bb-sites-dir> [<override-dir> ...]
 *
 * Env:
 *   TAP_SCRIPTS_SECRET - shared secret for X-Tap-Secret header
 *   TAP_API_URL        - batch endpoint (default: https://tap.vaayne.com/api/batch)
 */

import { createHash } from "node:crypto";
import { readFile, readdir } from "node:fs/promises";
import { basename, join } from "node:path";
import { pathToFileURL } from "node:url";

const metaPattern = /\/\*\s*@meta\s*\n([\s\S]*?)\*\//;
const allowedMatchFields = new Set(["domain", "executionDomain", "startPath"]);
const allowedSetFields = new Set(["executionDomain", "startPath"]);

export function parseMeta(content) {
  const match = content.match(metaPattern);
  if (!match) return null;
  try {
    return JSON.parse(match[1]);
  } catch {
    return null;
  }
}

export function applyMetadataOverride(content, override) {
  for (const field of Object.keys(override)) {
    if (!allowedSetFields.has(field)) {
      throw new Error(`unsupported compatibility field: ${field}`);
    }
  }

  const match = content.match(metaPattern);
  const meta = parseMeta(content);
  if (!match || !meta) throw new Error("cannot override invalid @meta block");

  const normalized = `/* @meta\n${JSON.stringify({ ...meta, ...override }, null, 2)}\n*/`;
  return (
    content.slice(0, match.index) +
    normalized +
    content.slice(match.index + match[0].length)
  );
}

export function applyCompatibilityPolicy(content, policy) {
  const policyFields = Object.keys(policy);
  if (
    policyFields.some((field) => field !== "match" && field !== "set") ||
    !policy.match ||
    !policy.set
  ) {
    throw new Error("compatibility policy requires only match and set objects");
  }
  const meta = parseMeta(content);
  if (!meta) throw new Error("cannot match invalid @meta block");
  for (const [field, expected] of Object.entries(policy.match)) {
    if (!allowedMatchFields.has(field)) {
      throw new Error(`unsupported compatibility match field: ${field}`);
    }
    if (meta[field] !== expected) {
      throw new Error(
        `stale compatibility policy: expected ${field}=${JSON.stringify(expected)}, got ${JSON.stringify(meta[field])}`,
      );
    }
  }
  return applyMetadataOverride(content, policy.set);
}

export async function loadCompatibility(path) {
  if (!path) return {};
  const manifest = JSON.parse(await readFile(path, "utf8"));
  if (
    manifest.version !== 1 ||
    !manifest.scripts ||
    Array.isArray(manifest.scripts)
  ) {
    throw new Error(`invalid compatibility manifest: ${path}`);
  }
  return manifest.scripts;
}

/** Discover all site/action.js scripts under a directory. */
export async function discoverScripts(dir, compatibility = {}) {
  const scripts = [];
  let entries;
  try {
    entries = await readdir(dir, { withFileTypes: true });
  } catch {
    return scripts;
  }

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    const site = entry.name;
    if (site.startsWith(".") || site === "node_modules") continue;

    const siteDir = join(dir, site);
    const files = await readdir(siteDir);
    for (const file of files) {
      if (!file.endsWith(".js")) continue;
      const filePath = join(siteDir, file);
      let content = await readFile(filePath, "utf8");
      let meta = parseMeta(content);
      if (!meta) {
        console.warn(`Skipping ${site}/${file}: no valid @meta block`);
        continue;
      }
      const name = meta.name || `${site}/${basename(file, ".js")}`;

      if (compatibility[name]) {
        content = applyCompatibilityPolicy(content, compatibility[name]);
        meta = parseMeta(content);
      }

      scripts.push({
        name,
        site,
        content,
        hash: createHash("sha256").update(content).digest("hex"),
        description: meta.description || "",
        domain: meta.domain || "",
        args: JSON.stringify(meta.args || {}),
        capabilities: meta.capabilities || [],
        example: meta.example || "",
        readOnly: meta.readOnly ?? true,
      });
    }
  }

  return scripts;
}

function parseArgs(args) {
  let compatibilityPath = "";
  if (args[0] === "--compat") {
    if (!args[1]) throw new Error("--compat requires a manifest path");
    compatibilityPath = args[1];
    args = args.slice(2);
  }
  if (args.length === 0) {
    throw new Error(
      "Usage: node sync-bb-sites.mjs [--compat <file>] <bb-sites-dir> [<override-dir> ...]",
    );
  }
  return { compatibilityPath, sitesDirs: args };
}

export async function buildCatalog(sitesDirs, compatibility) {
  const byName = new Map();
  for (const [index, dir] of sitesDirs.entries()) {
    // Compatibility policy belongs to the imported bb-sites source. Later Tap
    // directories remain authoritative overrides and are never rewritten.
    const found = await discoverScripts(dir, index === 0 ? compatibility : {});
    for (const script of found) byName.set(script.name, script);
    console.log(`${dir}: ${found.length} scripts`);

    if (index === 0) {
      const imported = new Set(found.map((script) => script.name));
      const missing = Object.keys(compatibility).filter(
        (name) => !imported.has(name),
      );
      if (missing.length > 0) {
        throw new Error(
          `stale bb-sites compatibility entries: ${missing.join(", ")}`,
        );
      }
    }
  }
  return [...byName.values()];
}

export async function main(args = process.argv.slice(2)) {
  const { compatibilityPath, sitesDirs } = parseArgs(args);
  const secret = process.env.TAP_SCRIPTS_SECRET;
  if (!secret) throw new Error("TAP_SCRIPTS_SECRET is required");

  const compatibility = await loadCompatibility(compatibilityPath);
  const scripts = await buildCatalog(sitesDirs, compatibility);
  if (scripts.length === 0) throw new Error("No scripts found");

  const apiUrl = process.env.TAP_API_URL || "https://tap.vaayne.com/api/batch";
  console.log(`Total: ${scripts.length} scripts, posting to ${apiUrl}`);
  const response = await fetch(apiUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Tap-Secret": secret },
    body: JSON.stringify({ scripts }),
  });
  const body = await response.text();
  if (!response.ok) {
    throw new Error(`Batch update failed: HTTP ${response.status}\n${body}`);
  }
  console.log(`Success: ${body}`);
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch((error) => {
    console.error(error.message);
    process.exitCode = 1;
  });
}
