import { createFileRoute } from "@tanstack/react-router"
import { getSyncManifest } from "@/lib/db"

export const Route = createFileRoute("/api/sync")({
  server: {
    handlers: {
      GET: async () => {
        try {
          const scripts = await getSyncManifest()
          return Response.json({ scripts })
        } catch (error) {
          console.error("Error getting sync manifest:", error)
          return Response.json(
            { error: "Failed to get sync manifest" },
            { status: 500 },
          )
        }
      },
    },
  },
})
