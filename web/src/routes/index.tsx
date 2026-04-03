import { createFileRoute, Link } from "@tanstack/react-router"
import {
  Terminal,
  Zap,
  Globe,
  Download,
  ArrowRight,
  Code,
  FileText,
  Layers,
  Github,
  ExternalLink,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { fetchPopularScripts } from "@/lib/server-fns"
import { ScriptCard } from "@/components/ScriptCard"

export const Route = createFileRoute("/")({
  loader: () => fetchPopularScripts({ data: { limit: 6 } }),
  component: HomePage,
})

function HomePage() {
  const { scripts, totalScripts, totalSites } = Route.useLoaderData()

  return (
    <main>
      {/* Hero */}
      <section className="relative overflow-hidden border-b">
        <div className="absolute inset-0 bg-gradient-to-b from-muted/50 to-background" />
        <div className="relative mx-auto max-w-6xl px-4 py-20 sm:py-32">
          <div className="mx-auto max-w-3xl text-center">
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border bg-background px-4 py-1.5 text-sm">
              <Zap className="size-3.5 text-yellow-500" />
              <span className="text-muted-foreground">
                QuickJS first, browser fallback when needed
              </span>
            </div>

            <h1 className="mb-4 text-4xl font-bold tracking-tight sm:text-6xl">
              🚰 Tap into any website
              <br />
              <span className="text-muted-foreground">from your terminal</span>
            </h1>

            <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground">
              A Go library and CLI that runs JavaScript scripts against real
              websites — fast via QuickJS, with full browser fallback. Also
              extracts clean content from any URL.
            </p>

            <div className="flex flex-wrap items-center justify-center gap-3">
              <Button size="lg" render={<Link to="/scripts" />}>
                <Terminal className="size-4" />
                Browse {totalScripts} Scripts
              </Button>
              <Button
                variant="outline"
                size="lg"
                render={
                  <a
                    href="https://github.com/vaayne/tap"
                    target="_blank"
                    rel="noopener noreferrer"
                  />
                }
              >
                <Github className="size-4" />
                GitHub
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Install */}
      <section className="border-b">
        <div className="mx-auto max-w-6xl px-4 py-16">
          <div className="mb-10 text-center">
            <h2 className="mb-2 text-2xl font-bold sm:text-3xl">
              Get started in seconds
            </h2>
            <p className="text-muted-foreground">
              Install the CLI or use as a Go library
            </p>
          </div>

          <div className="mx-auto grid max-w-3xl gap-4 sm:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <Download className="size-4" />
                  CLI
                </CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="overflow-x-auto rounded-lg bg-muted p-3 text-sm font-mono">
                  go install github.com/vaayne/tap/cmd/tap@latest
                </pre>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <Code className="size-4" />
                  Library
                </CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="overflow-x-auto rounded-lg bg-muted p-3 text-sm font-mono">
                  go get github.com/vaayne/tap
                </pre>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="border-b">
        <div className="mx-auto max-w-6xl px-4 py-16">
          <div className="mb-10 text-center">
            <h2 className="mb-2 text-2xl font-bold sm:text-3xl">
              How it works
            </h2>
            <p className="text-muted-foreground">
              Two commands, one shared transport layer
            </p>
          </div>

          <div className="mx-auto max-w-4xl">
            <pre className="overflow-x-auto rounded-xl border bg-muted/50 p-6 text-xs sm:text-sm font-mono leading-relaxed text-center">
{`                ┌─────────────────────────────────┐
                │       Shared Transport Layer     │
                │  Level 1: HTTP  │  Level 2: CDP  │
                └────────┬────────┴───────┬────────┘
                         │                │
          ┌──────────────┴──┐    ┌────────┴──────────┐
          │    tap site     │    │    tap fetch       │
          │  QuickJS → CDP  │    │  HTTP → CDP        │
          │  → structured   │    │  → defuddle        │
          │    JSON         │    │  → markdown/HTML   │
          └─────────────────┘    └───────────────────┘`}
            </pre>
          </div>

          <div className="mx-auto mt-10 grid max-w-4xl gap-6 sm:grid-cols-3">
            <FeatureCard
              icon={<Layers className="size-5" />}
              title="Transport Layer"
              description="Shared HTTP client and headless Chrome (CDP), configured once and used by all consumers."
            />
            <FeatureCard
              icon={<Terminal className="size-5" />}
              title="Site Scripts"
              description="Predefined recipes that fetch structured data. QuickJS first, Chrome fallback for pages needing cookies or DOM."
            />
            <FeatureCard
              icon={<FileText className="size-5" />}
              title="Content Extraction"
              description="Generic clean content from any URL via go-defuddle. Direct HTTP first, browser for JS-rendered pages."
            />
          </div>
        </div>
      </section>

      {/* Usage Examples */}
      <section className="border-b">
        <div className="mx-auto max-w-6xl px-4 py-16">
          <div className="mb-10 text-center">
            <h2 className="mb-2 text-2xl font-bold sm:text-3xl">
              Usage examples
            </h2>
            <p className="text-muted-foreground">
              Pipe to jq, combine with other tools
            </p>
          </div>

          <div className="mx-auto max-w-3xl space-y-4">
            <UsageExample
              label="Run site scripts"
              commands={[
                "tap site list",
                "tap site v2ex/hot",
                'tap site twitter/search query="claude code"',
                "tap site hackernews/top | jq '.stories[:3]'",
              ]}
            />
            <UsageExample
              label="Extract clean content"
              commands={[
                "tap fetch https://example.com/article",
                "tap fetch --json https://example.com/article",
              ]}
            />
          </div>
        </div>
      </section>

      {/* Stats bar */}
      <section className="border-b bg-muted/30">
        <div className="mx-auto grid max-w-6xl grid-cols-3 divide-x px-4">
          <StatItem value={totalScripts} label="Scripts" />
          <StatItem value={totalSites} label="Sites" />
          <StatItem value="Go" label="Language" />
        </div>
      </section>

      {/* Popular Scripts */}
      <section className="border-b">
        <div className="mx-auto max-w-6xl px-4 py-16">
          <div className="mb-8 flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold sm:text-3xl">
                Popular scripts
              </h2>
              <p className="mt-1 text-muted-foreground">
                {totalScripts} scripts across {totalSites} platforms
              </p>
            </div>
            <Button variant="outline" render={<Link to="/scripts" />}>
              View all
              <ArrowRight className="size-3.5" />
            </Button>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {scripts.map((script) => (
              <ScriptCard key={script.name} script={script} />
            ))}
          </div>
        </div>
      </section>

      {/* Credits */}
      <section className="border-b bg-muted/30">
        <div className="mx-auto max-w-6xl px-4 py-12 text-center">
          <h2 className="mb-3 text-xl font-bold">Standing on the shoulders of giants</h2>
          <p className="mb-6 text-muted-foreground">
            Based on{" "}
            <a
              href="https://github.com/epiral/bb-sites"
              target="_blank"
              rel="noopener noreferrer"
              className="underline underline-offset-4 hover:text-foreground"
            >
              bb-sites
            </a>
            {" "}and fully compatible with{" "}
            <a
              href="https://github.com/epiral/bb-browser"
              target="_blank"
              rel="noopener noreferrer"
              className="underline underline-offset-4 hover:text-foreground"
            >
              bb-browser
            </a>
          </p>
          <div className="flex flex-wrap items-center justify-center gap-3">
            <Badge variant="secondary" className="gap-1.5 px-3 py-1">
              <ExternalLink className="size-3" />
              <a
                href="https://github.com/epiral/bb-sites"
                target="_blank"
                rel="noopener noreferrer"
              >
                epiral/bb-sites
              </a>
            </Badge>
            <Badge variant="secondary" className="gap-1.5 px-3 py-1">
              <ExternalLink className="size-3" />
              <a
                href="https://github.com/epiral/bb-browser"
                target="_blank"
                rel="noopener noreferrer"
              >
                epiral/bb-browser
              </a>
            </Badge>
            <Badge variant="secondary" className="gap-1.5 px-3 py-1">
              <Github className="size-3" />
              <a
                href="https://github.com/vaayne/tap"
                target="_blank"
                rel="noopener noreferrer"
              >
                vaayne/tap
              </a>
            </Badge>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-8 text-center text-sm text-muted-foreground">
        <p>
          MIT License ·{" "}
          <a
            href="https://github.com/vaayne/tap"
            target="_blank"
            rel="noopener noreferrer"
            className="underline underline-offset-4 hover:text-foreground"
          >
            GitHub
          </a>
        </p>
      </footer>
    </main>
  )
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode
  title: string
  description: string
}) {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="mb-3 flex size-10 items-center justify-center rounded-lg bg-muted">
          {icon}
        </div>
        <h3 className="mb-1 font-semibold">{title}</h3>
        <p className="text-sm text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  )
}

function UsageExample({
  label,
  commands,
}: {
  label: string
  commands: string[]
}) {
  return (
    <div className="rounded-xl border bg-muted/30 overflow-hidden">
      <div className="border-b bg-muted/50 px-4 py-2 text-xs font-medium text-muted-foreground">
        {label}
      </div>
      <pre className="overflow-x-auto p-4 text-sm font-mono leading-relaxed">
        {commands.map((cmd, i) => (
          <div key={i}>
            <span className="select-none text-muted-foreground">$ </span>
            {cmd}
          </div>
        ))}
      </pre>
    </div>
  )
}

function StatItem({
  value,
  label,
}: {
  value: number | string
  label: string
}) {
  return (
    <div className="flex items-center justify-center gap-2 py-5 text-sm">
      <span className="text-lg font-bold">{value}</span>
      <span className="text-muted-foreground">{label}</span>
    </div>
  )
}
