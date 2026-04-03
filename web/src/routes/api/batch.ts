import { createFileRoute } from "@tanstack/react-router"
import { verifySecret } from "@/lib/auth"
import { batchUpdate } from "@/lib/db"
import type { BatchScript } from "@/lib/types"

const MAX_PAYLOAD_BYTES = 500 * 1024 // 500 KB

export const Route = createFileRoute("/api/batch")({
  server: {
    handlers: {
      POST: async ({ request }) => {
        // Authenticate
        if (!verifySecret(request)) {
          return Response.json(
            { error: "Unauthorized" },
            { status: 401 },
          )
        }

        // Read body as text to enforce size limit regardless of Content-Length
        const text = await request.text()
        if (text.length > MAX_PAYLOAD_BYTES) {
          return Response.json(
            { error: "Payload too large (max 500KB)" },
            { status: 413 },
          )
        }

        // Parse body
        let body: { scripts?: unknown }
        try {
          body = JSON.parse(text)
        } catch {
          return Response.json(
            { error: "Invalid JSON body" },
            { status: 400 },
          )
        }

        // Validate shape
        if (!Array.isArray(body?.scripts)) {
          return Response.json(
            { error: "Missing or invalid 'scripts' array" },
            { status: 400 },
          )
        }

        const scripts = body.scripts as BatchScript[]

        // Basic per-script validation
        const names = new Set<string>()
        for (const s of scripts) {
          if (!s.name || typeof s.name !== "string") {
            return Response.json(
              { error: "Each script must have a string 'name'" },
              { status: 400 },
            )
          }
          if (names.has(s.name)) {
            return Response.json(
              { error: `Duplicate script name: ${s.name}` },
              { status: 400 },
            )
          }
          names.add(s.name)
        }

        // Perform batch update
        try {
          await batchUpdate(scripts)
        } catch (error) {
          console.error("batchUpdate failed:", error)
          return Response.json(
            { error: "Failed to update scripts" },
            { status: 500 },
          )
        }

        return Response.json({ success: true, count: scripts.length })
      },
    },
  },
})
