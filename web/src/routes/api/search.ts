import { createFileRoute } from "@tanstack/react-router"
import { listScripts, listSites } from "@/lib/db"

export const Route = createFileRoute("/api/search")({
  server: {
    handlers: {
      GET: async ({ request }) => {
        try {
          const url = new URL(request.url)
          const query = url.searchParams.get("q") ?? undefined
          const sort =
            (url.searchParams.get("sort") as
              | "popular"
              | "newest"
              | "name") ?? "name"
          const site = url.searchParams.get("site") ?? undefined

          const [scripts, sites] = await Promise.all([
            listScripts({ query, sort, site }),
            listSites(),
          ])

          return Response.json({ scripts, sites })
        } catch (error) {
          console.error("Error listing scripts:", error)
          return Response.json(
            { error: "Failed to list scripts" },
            { status: 500 },
          )
        }
      },
    },
  },
})
