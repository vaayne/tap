import {
  HeadContent,
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
      { title: "tap — Tap into any website from your terminal" },
      {
        name: "description",
        content:
          "A Go CLI and library that runs JavaScript scripts against real websites. 100+ scripts across 40+ sites. QuickJS first, browser fallback.",
      },
    ],
    links: [{ rel: "stylesheet", href: appCss }],
  }),
  component: RootComponent,
})

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
