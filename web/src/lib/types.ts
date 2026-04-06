// Script metadata parsed from @meta block in JS files
export type ScriptMeta = {
  name: string
  description: string
  domain: string
  args: Record<string, ScriptArg>
  readOnly: boolean
  example: string
}

export type ScriptArg = {
  required: boolean
  description: string
}

// Database row for scripts table
export type ScriptRow = {
  name: string
  site: string
  description: string
  domain: string
  args: string // JSON string
  read_only: number // 0 or 1
  example: string
  capabilities: string // JSON string, e.g. '["network"]'
  content: string
  hash: string
  created_at: string
  updated_at: string
}

// API response types
export type ScriptListItem = {
  name: string
  site: string
  description: string
  domain: string
  readOnly: boolean
  example: string
  capabilities: string[]
  args: Record<string, ScriptArg>
  usageCount: number
  updatedAt: string
}

export type ScriptDetail = ScriptListItem & {
  content: string
  hash: string
  createdAt: string
  usage: {
    total: number
    last7d: number
    last30d: number
  }
}

export type SyncManifestItem = {
  name: string
  hash: string
  updatedAt: string
}

// Payload type for batch update from bb-sites
export type BatchScript = {
  name: string
  site: string
  content: string
  hash: string
  description: string
  domain: string
  args: string // JSON string
  capabilities: string[] // e.g. ["network"]
  example: string
  readOnly: boolean
}

export type UsageEvent = {
  id: number
  script_name: string
  reported_at: string
}
