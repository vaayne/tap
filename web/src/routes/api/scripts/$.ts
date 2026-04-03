import { createFileRoute } from "@tanstack/react-router"
import { getScript, getScriptContent } from "@/lib/db"

export const Route = createFileRoute("/api/scripts/$")({  
  server: {
    handlers: {
      GET: async ({ request, params }) => {
        try {
          const splat = params._splat
          if (!splat) {
            return Response.json(
              { error: "Script name is required" },
              { status: 400 },
            )
          }

          // Check if requesting raw content: /api/scripts/google/search/content
          if (splat.endsWith("/content")) {
            const name = splat.replace(/\/content$/, "")
            const content = await getScriptContent(name)
            if (!content) {
              return Response.json(
                { error: `Script '${name}' not found` },
                { status: 404 },
              )
            }
            return new Response(content, {
              headers: {
                "Content-Type": "application/javascript",
                "Cache-Control": "public, max-age=300",
              },
            })
          }

          // Otherwise return full script detail
          const script = await getScript(splat)
          if (!script) {
            return Response.json(
              { error: `Script '${splat}' not found` },
              { status: 404 },
            )
          }

          return Response.json(script)
        } catch (error) {
          console.error("Error getting script:", error)
          return Response.json(
            { error: "Failed to get script" },
            { status: 500 },
          )
        }
      },
    },
  },
})
