import { createFileRoute } from "@tanstack/react-router"
import { reportUsage } from "@/lib/db"

export const Route = createFileRoute("/api/usage")({
  server: {
    handlers: {
      POST: async ({ request }) => {
        try {
          const body = (await request.json()) as { script?: string }
          const scriptName = body?.script

          if (!scriptName || typeof scriptName !== "string") {
            return Response.json(
              { error: "Missing required field: script" },
              { status: 400 },
            )
          }

          await reportUsage(scriptName)

          return Response.json(
            { ok: true },
            { status: 201 },
          )
        } catch (error) {
          console.error("Error reporting usage:", error)
          return Response.json(
            { error: "Failed to report usage" },
            { status: 500 },
          )
        }
      },
    },
  },
})
