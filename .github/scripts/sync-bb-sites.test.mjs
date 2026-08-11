import assert from "node:assert/strict";
import { mkdtemp, mkdir, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import {
  applyMetadataOverride,
  applyCompatibilityPolicy,
  buildCatalog,
  discoverScripts,
  parseMeta,
} from "./sync-bb-sites.mjs";

const source = `/* @meta
{
  "name": "hackernews/top",
  "description": "HN top stories",
  "domain": "news.ycombinator.com",
  "args": {"count": {"required": false}}
}
*/
async function(args) {
  return fetch('https://hacker-news.firebaseio.com/v0/topstories.json');
}
`;

test("metadata overrides preserve upstream code and unrelated fields", () => {
  const normalized = applyMetadataOverride(source, {
    executionDomain: "hacker-news.firebaseio.com",
    startPath: "/v0/topstories.json",
  });
  const meta = parseMeta(normalized);

  assert.equal(meta.name, "hackernews/top");
  assert.equal(meta.description, "HN top stories");
  assert.deepEqual(meta.args, { count: { required: false } });
  assert.equal(meta.domain, "news.ycombinator.com");
  assert.equal(meta.executionDomain, "hacker-news.firebaseio.com");
  assert.equal(meta.startPath, "/v0/topstories.json");
  assert.match(
    normalized,
    /return fetch\('https:\/\/hacker-news\.firebaseio\.com/,
  );
});

test("metadata overrides reject fields outside the execution policy", () => {
  assert.throws(
    () =>
      applyMetadataOverride(source, { headers: { Authorization: "secret" } }),
    /unsupported compatibility field: headers/,
  );
});

test("compatibility policy fails when upstream metadata changes", () => {
  assert.throws(
    () =>
      applyCompatibilityPolicy(source, {
        match: { domain: "already-fixed.example" },
        set: { executionDomain: "hacker-news.firebaseio.com" },
      }),
    /stale compatibility policy: expected domain="already-fixed.example", got "news.ycombinator.com"/,
  );
});

test("Tap scripts derive their names from paths", async () => {
  const root = await mkdtemp(join(tmpdir(), "tap-native-sites-"));
  await mkdir(join(root, "example"), { recursive: true });
  await writeFile(
    join(root, "example", "search.js"),
    source.replace('  "name": "hackernews/top",\n', ""),
  );

  const scripts = await discoverScripts(root);
  assert.equal(scripts.length, 1);
  assert.equal(scripts[0].name, "example/search");
});

test("compatibility applies only to imported bb-sites, not Tap overrides", async () => {
  const root = await mkdtemp(join(tmpdir(), "tap-bb-sites-"));
  const upstream = join(root, "bb-sites");
  const local = join(root, "sites");
  await mkdir(join(upstream, "hackernews"), { recursive: true });
  await mkdir(join(local, "hackernews"), { recursive: true });
  await writeFile(join(upstream, "hackernews", "top.js"), source);
  const localSource = source.replace("HN top stories", "Tap override");
  await writeFile(join(local, "hackernews", "top.js"), localSource);

  const compatibility = {
    "hackernews/top": {
      match: { domain: "news.ycombinator.com" },
      set: {
        executionDomain: "hacker-news.firebaseio.com",
        startPath: "/v0/topstories.json",
      },
    },
  };
  const imported = await discoverScripts(upstream, compatibility);
  assert.equal(parseMeta(imported[0].content).domain, "news.ycombinator.com");
  assert.equal(
    parseMeta(imported[0].content).executionDomain,
    "hacker-news.firebaseio.com",
  );

  const catalog = await buildCatalog([upstream, local], compatibility);
  assert.equal(catalog.length, 1);
  assert.equal(catalog[0].description, "Tap override");
  assert.equal(parseMeta(catalog[0].content).domain, "news.ycombinator.com");
});

test("stale compatibility entries fail the import", async () => {
  const root = await mkdtemp(join(tmpdir(), "tap-bb-sites-empty-"));
  await assert.rejects(
    buildCatalog([root], {
      "missing/script": {
        match: { domain: "old.example.com" },
        set: { executionDomain: "new.example.com" },
      },
    }),
    /stale bb-sites compatibility entries: missing\/script/,
  );
});
