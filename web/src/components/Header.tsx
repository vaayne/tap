import { Link } from "@tanstack/react-router"
import { Terminal, Search, Github } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function Header() {
  return (
    <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto flex h-14 max-w-6xl items-center gap-6 px-4">
        <Link to="/" className="flex items-center gap-2 font-semibold">
          <Terminal className="size-5" />
          <span className="text-base">tap</span>
        </Link>

        <nav className="flex items-center gap-1">
          <Button variant="ghost" size="sm" render={<Link to="/" />}>
            Home
          </Button>
          <Button variant="ghost" size="sm" render={<Link to="/scripts" />}>
            Scripts
          </Button>
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            className="text-muted-foreground"
            render={<Link to="/scripts" search={{ q: "" }} />}
          >
            <Search className="size-3.5" />
            <span className="hidden sm:inline">Search scripts…</span>
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            render={
              <a
                href="https://github.com/vaayne/tap"
                target="_blank"
                rel="noopener noreferrer"
              />
            }
          >
            <Github className="size-4" />
          </Button>
        </div>
      </div>
    </header>
  )
}
