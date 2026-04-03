-- Core scripts table
CREATE TABLE IF NOT EXISTS scripts (
  name        TEXT PRIMARY KEY,
  site        TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  domain      TEXT NOT NULL DEFAULT '',
  args        TEXT NOT NULL DEFAULT '{}',
  read_only   INTEGER NOT NULL DEFAULT 1,
  example     TEXT NOT NULL DEFAULT '',
  content     TEXT NOT NULL,
  hash        TEXT NOT NULL,
  created_at  TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_scripts_site ON scripts(site);

-- Raw usage events
CREATE TABLE IF NOT EXISTS usage_events (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  script_name TEXT NOT NULL,
  reported_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_usage_script ON usage_events(script_name);
CREATE INDEX IF NOT EXISTS idx_usage_time ON usage_events(reported_at);
