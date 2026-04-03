import { Link } from "@tanstack/react-router"
import { TrendingUp } from "lucide-react"
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { ScriptListItem } from "@/lib/types"

export function ScriptCard({ script }: { script: ScriptListItem }) {
  const [site, command] = script.name.split("/")
  return (
    <Link
      to="/scripts/$site/$command"
      params={{ site, command }}
      className="block focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded-xl"
    >
      <Card className="h-full transition-shadow hover:shadow-md">
        <CardHeader>
          <div className="flex items-start justify-between gap-2">
            <CardTitle className="font-mono text-sm">
              {script.name}
            </CardTitle>
            <Badge variant="secondary">{script.site}</Badge>
          </div>
          <CardDescription className="line-clamp-2">
            {script.description || "No description"}
          </CardDescription>
        </CardHeader>
        <CardFooter className="text-xs text-muted-foreground">
          <TrendingUp className="size-3.5 mr-1" />
          {script.usageCount} {script.usageCount === 1 ? "use" : "uses"}
          {script.domain && (
            <span className="ml-auto truncate text-muted-foreground/60">
              {script.domain}
            </span>
          )}
        </CardFooter>
      </Card>
    </Link>
  )
}
