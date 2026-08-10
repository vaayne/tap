import {
  HeadContent,
  Link,
  Outlet,
  Scripts,
  createRootRoute,
} from "@tanstack/react-router"

import Header from "../components/Header"

import appCss from "../styles.css?url"

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      { title: "tap — Reusable site programs for agent-browser" },
      {
        name: "description",
        content:
          "Reusable site programs and readable web extraction powered by agent-browser.",
      },
    ],
    links: [{ rel: "stylesheet", href: appCss }],
  }),
  component: RootComponent,
  notFoundComponent: NotFound,
})

function NotFound() {
  return (
    <main className="flex flex-1 flex-col items-center justify-center gap-4 py-24">
      <h1 className="text-4xl font-bold">404</h1>
      <p className="text-muted-foreground">Page not found.</p>
      <Link to="/" className="text-primary underline underline-offset-4 hover:text-primary/80">
        Go home
      </Link>
    </main>
  )
}

function RootComponent() {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className="min-h-screen bg-background font-sans antialiased">
        <Header />
        <Outlet />
        <Scripts />
      </body>
    </html>
  )
}
