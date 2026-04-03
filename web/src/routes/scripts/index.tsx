import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useState, useEffect, useCallback } from "react"
import { Search, X } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import { ScriptCard } from "@/components/ScriptCard"
import { fetchScriptList } from "@/lib/server-fns"

type SearchParams = {
  q?: string
  sort?: string
  site?: string
}

export const Route = createFileRoute("/scripts/")({
  validateSearch: (search: Record<string, unknown>): SearchParams => ({
    q: (search.q as string) || undefined,
    sort: (search.sort as string) || undefined,
    site: (search.site as string) || undefined,
  }),
  loaderDeps: ({ search }) => search,
  loader: ({ deps }) =>
    fetchScriptList({
      data: {
        query: deps.q,
        sort: deps.sort || "popular",
        site: deps.site,
      },
    }),
  component: BrowsePage,
  pendingComponent: BrowsePending,
})

function BrowsePage() {
  const { scripts, sites } = Route.useLoaderData()
  const search = Route.useSearch()
  const navigate = useNavigate()

  const [query, setQuery] = useState(search.q ?? "")

  // Sync input with URL search param
  useEffect(() => {
    setQuery(search.q ?? "")
  }, [search.q])

  // Debounced search
  const updateSearch = useCallback(
    (updates: Partial<SearchParams>) => {
      navigate({
        to: "/scripts",
        search: (prev: SearchParams) => ({
          ...prev,
          ...updates,
        }),
        replace: true,
      })
    },
    [navigate],
  )

  useEffect(() => {
    const timer = setTimeout(() => {
      const trimmed = query.trim()
      if (trimmed !== (search.q ?? "")) {
        updateSearch({ q: trimmed || undefined })
      }
    }, 300)
    return () => clearTimeout(timer)
  }, [query, search.q, updateSearch])

  const currentSort = search.sort || "popular"
  const currentSite = search.site || ""

  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">Browse Scripts</h1>
        <p className="text-muted-foreground">
          {scripts.length} script{scripts.length !== 1 ? "s" : ""} found
        </p>
      </div>

      {/* Filters */}
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <div className="relative flex-1 min-w-[200px] max-w-sm">
          <Search className="absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search scripts…"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="pl-8"
          />
        </div>

        <Select
          value={currentSort}
          onValueChange={(val) => updateSearch({ sort: val || undefined })}
        >
          <SelectTrigger className="w-[140px]">
            <SelectValue placeholder="Sort by" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="popular">Popular</SelectItem>
            <SelectItem value="newest">Newest</SelectItem>
            <SelectItem value="name">Name</SelectItem>
          </SelectContent>
        </Select>

        <Select
          value={currentSite}
          onValueChange={(val) =>
            updateSearch({ site: val || undefined })
          }
        >
          <SelectTrigger className="w-[160px]">
            <SelectValue placeholder="All sites" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">All sites</SelectItem>
            {sites.map((s) => (
              <SelectItem key={s} value={s}>
                {s}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {(search.q || search.site) && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() =>
              navigate({
                to: "/scripts",
                search: { sort: search.sort },
              })
            }
          >
            <X className="size-3.5" />
            Clear filters
          </Button>
        )}
      </div>

      {/* Results */}
      {scripts.length === 0 ? (
        <div className="py-16 text-center text-muted-foreground">
          <p className="text-lg">No scripts found</p>
          <p className="text-sm">Try adjusting your search or filters.</p>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {scripts.map((script) => (
            <ScriptCard key={script.name} script={script} />
          ))}
        </div>
      )}
    </main>
  )
}

function BrowsePending() {
  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <div className="mb-8">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="mt-2 h-5 w-32" />
      </div>
      <div className="mb-6 flex gap-3">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-8 w-40" />
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-36 rounded-xl" />
        ))}
      </div>
    </main>
  )
}
