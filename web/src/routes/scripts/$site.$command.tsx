import { createFileRoute, Link } from "@tanstack/react-router"
import {
  ArrowLeft,
  TrendingUp,
  Calendar,
  Globe,
  Lock,
  Unlock,
  Copy,
  Check,
} from "lucide-react"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Separator } from "@/components/ui/separator"
import { Skeleton } from "@/components/ui/skeleton"
import { fetchScriptDetail } from "@/lib/server-fns"
import type { ScriptDetail } from "@/lib/types"

export const Route = createFileRoute("/scripts/$site/$command")({
  loader: ({ params }) =>
    fetchScriptDetail({ data: { name: `${params.site}/${params.command}` } }),
  component: DetailPage,
  pendingComponent: DetailPending,
  notFoundComponent: () => (
    <main className="mx-auto max-w-4xl px-4 py-16 text-center">
      <p className="text-lg text-muted-foreground">Script not found</p>
      <Button variant="ghost" className="mt-4" render={<Link to="/scripts" />}>
        ← Back to scripts
      </Button>
    </main>
  ),
})

function DetailPage() {
  const script = Route.useLoaderData() as ScriptDetail | null

  if (!script) {
    return (
      <main className="mx-auto max-w-4xl px-4 py-16 text-center">
        <p className="text-lg text-muted-foreground">Script not found</p>
        <Button
          variant="ghost"
          className="mt-4"
          render={<Link to="/scripts" />}
        >
          ← Back to scripts
        </Button>
      </main>
    )
  }

  const args = Object.entries(script.args)

  return (
    <main className="mx-auto max-w-4xl px-4 py-8">
      {/* Back link */}
      <Button
        variant="ghost"
        size="sm"
        className="mb-4"
        render={<Link to="/scripts" />}
      >
        <ArrowLeft className="size-3.5" />
        Back
      </Button>

      {/* Header */}
      <div className="mb-6">
        <div className="flex flex-wrap items-center gap-2 mb-2">
          <h1 className="text-2xl font-bold font-mono">{script.name}</h1>
          <Badge variant="secondary">{script.site}</Badge>
          {script.readOnly ? (
            <Badge variant="outline">
              <Lock className="size-3" />
              Read-only
            </Badge>
          ) : (
            <Badge variant="outline">
              <Unlock className="size-3" />
              Read-write
            </Badge>
          )}
        </div>
        <p className="text-muted-foreground">
          {script.description || "No description"}
        </p>
        {script.domain && (
          <div className="mt-1 flex items-center gap-1 text-sm text-muted-foreground/70">
            <Globe className="size-3.5" />
            {script.domain}
          </div>
        )}
      </div>

      {/* Usage stats */}
      <div className="mb-6 grid grid-cols-3 gap-3">
        <UsageCard label="Last 7 days" count={script.usage.last7d} />
        <UsageCard label="Last 30 days" count={script.usage.last30d} />
        <UsageCard label="All time" count={script.usage.total} />
      </div>

      <Separator className="mb-6" />

      {/* Tabs */}
      <Tabs defaultValue="code">
        <TabsList>
          <TabsTrigger value="code">Code</TabsTrigger>
          {args.length > 0 && (
            <TabsTrigger value="args">
              Arguments ({args.length})
            </TabsTrigger>
          )}
          {script.example && (
            <TabsTrigger value="example">Example</TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="code" className="mt-4">
          <CodeBlock content={script.content} name={script.name} />
        </TabsContent>

        {args.length > 0 && (
          <TabsContent value="args" className="mt-4">
            <Card>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Required</TableHead>
                    <TableHead>Description</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {args.map(([name, arg]) => (
                    <TableRow key={name}>
                      <TableCell className="font-mono text-sm">
                        {name}
                      </TableCell>
                      <TableCell>
                        {arg.required ? (
                          <Badge variant="default" className="text-xs">
                            Required
                          </Badge>
                        ) : (
                          <span className="text-muted-foreground text-xs">
                            Optional
                          </span>
                        )}
                      </TableCell>
                      <TableCell className="text-muted-foreground whitespace-normal">
                        {arg.description}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Card>
          </TabsContent>
        )}

        {script.example && (
          <TabsContent value="example" className="mt-4">
            <Card>
              <CardContent>
                <pre className="overflow-x-auto text-sm font-mono whitespace-pre-wrap">
                  {script.example}
                </pre>
              </CardContent>
            </Card>
          </TabsContent>
        )}
      </Tabs>

      {/* Meta footer */}
      <Separator className="my-6" />
      <div className="flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <Calendar className="size-3" />
          Updated {formatDate(script.updatedAt)}
        </span>
        <span className="flex items-center gap-1">
          <Calendar className="size-3" />
          Created {formatDate(script.createdAt)}
        </span>
        <span className="font-mono">SHA-256: {script.hash.slice(0, 12)}…</span>
      </div>
    </main>
  )
}

function UsageCard({ label, count }: { label: string; count: number }) {
  return (
    <Card size="sm">
      <CardHeader className="pb-0">
        <CardTitle className="text-xs font-normal text-muted-foreground">
          {label}
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="flex items-center gap-1">
          <TrendingUp className="size-3.5 text-muted-foreground" />
          <span className="text-2xl font-bold">{count}</span>
        </div>
      </CardContent>
    </Card>
  )
}

function CodeBlock({
  content,
  name,
}: {
  content: string
  name: string
}) {
  const [copied, setCopied] = useState(false)

  function handleCopy() {
    navigator.clipboard.writeText(content)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between border-b py-2">
        <span className="font-mono text-xs text-muted-foreground">
          {name}.js
        </span>
        <Button variant="ghost" size="icon-xs" onClick={handleCopy}>
          {copied ? (
            <Check className="size-3" />
          ) : (
            <Copy className="size-3" />
          )}
        </Button>
      </CardHeader>
      <CardContent className="p-0">
        <pre className="overflow-x-auto p-4 text-sm font-mono leading-relaxed">
          <code>{content}</code>
        </pre>
      </CardContent>
    </Card>
  )
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr.replace(" ", "T") + "Z")
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  })
}

function DetailPending() {
  return (
    <main className="mx-auto max-w-4xl px-4 py-8">
      <Skeleton className="mb-4 h-7 w-16" />
      <Skeleton className="mb-2 h-8 w-64" />
      <Skeleton className="mb-6 h-5 w-96" />
      <div className="mb-6 grid grid-cols-3 gap-3">
        <Skeleton className="h-20 rounded-xl" />
        <Skeleton className="h-20 rounded-xl" />
        <Skeleton className="h-20 rounded-xl" />
      </div>
      <Skeleton className="h-64 rounded-xl" />
    </main>
  )
}
