package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserProxyCmd() *cli.Command {
	return &cli.Command{
		Name:  "proxy",
		Usage: "Run a local CDP proxy in front of an existing browser",
		Description: `Runs a local HTTP+WebSocket CDP proxy that fronts an existing Chrome-compatible
browser debug endpoint. This lets tap attach multiple commands to a single user
browser while keeping one stable local endpoint.

Examples:
  tap browser proxy --upstream http://127.0.0.1:9222
  tap browser proxy --listen 127.0.0.1:9401 --upstream ws://127.0.0.1:9222/devtools/browser/abc

Then point tap at the proxy:
  TAP_WS_URL=http://127.0.0.1:9401 tap fetch -b https://example.com
  tap browser session new default --ws-url http://127.0.0.1:9401`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Usage: "Local listen address for the proxy",
				Value: "127.0.0.1:9401",
			},
			&cli.StringFlag{
				Name:  "upstream",
				Usage: "Upstream browser DevTools endpoint (WebSocket URL or HTTP base URL)",
				Value: "http://127.0.0.1:9222",
			},
			&cli.BoolFlag{
				Name:  "user-chrome",
				Usage: "Auto-discover the DevTools endpoint from the running user Chrome/Chromium profile",
			},
			&cli.StringFlag{
				Name:  "devtools-port-file",
				Usage: "Path to a DevToolsActivePort file to use instead of --upstream",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			upstream := cmd.String("upstream")
			if path := cmd.String("devtools-port-file"); path != "" {
				url, err := browser.ResolveDebugURLFromDevToolsFile(path)
				if err != nil {
					return fmt.Errorf("resolve devtools port file: %w", err)
				}
				upstream = url
				fmt.Fprintf(os.Stderr, "Resolved upstream from %s\n", path)
			} else if cmd.Bool("user-chrome") {
				url, path, err := browser.DiscoverUserChromeDebugURL()
				if err != nil {
					return err
				}
				upstream = url
				fmt.Fprintf(os.Stderr, "Discovered user Chrome via %s\n", path)
			}

			proxy := browser.NewProxy(browser.ProxyConfig{
				ListenAddr: cmd.String("listen"),
				Upstream:   upstream,
			})

			runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			fmt.Fprintf(os.Stderr, "Proxy listening on http://%s -> %s\n", cmd.String("listen"), upstream)
			fmt.Fprintf(os.Stderr, "Press Ctrl+C to stop.\n")
			return proxy.Serve(runCtx)
		},
	}
}
