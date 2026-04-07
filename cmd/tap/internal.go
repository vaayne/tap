package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func internalCmd() *cli.Command {
	return &cli.Command{
		Name:   "internal",
		Usage:  "Internal tap entrypoints",
		Hidden: true,
		Commands: []*cli.Command{
			proxyDaemonCmd(),
		},
	}
}

func proxyDaemonCmd() *cli.Command {
	return &cli.Command{
		Name:   "proxy-daemon",
		Usage:  "Run the internal CDP proxy daemon",
		Hidden: true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "listen", Required: true, Hidden: true},
			&cli.StringFlag{Name: "upstream-ws-url", Required: true, Hidden: true},
			&cli.StringFlag{Name: "ownership-token", Required: true, Hidden: true},
			&cli.StringFlag{Name: "ready-file", Required: true, Hidden: true},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ln, err := net.Listen("tcp", cmd.String("listen"))
			if err != nil {
				return err
			}
			defer func() { _ = ln.Close() }()

			record := &browser.ProxyDaemonRecord{
				PID:            os.Getpid(),
				ListenAddr:     ln.Addr().String(),
				Endpoint:       browser.ProxyEndpointForListenAddr(ln.Addr().String()),
				UpstreamWSURL:  cmd.String("upstream-ws-url"),
				OwnershipToken: cmd.String("ownership-token"),
				State:          browser.AttachStateAttachedReady,
				Status:         browser.AttachStateAttachedReady,
				StartedAt:      time.Now().UTC(),
				UpdatedAt:      time.Now().UTC(),
			}
			data, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("encode proxy daemon ready file: %w", err)
			}
			if err := os.WriteFile(cmd.String("ready-file"), data, 0o600); err != nil {
				return fmt.Errorf("write proxy daemon ready file: %w", err)
			}

			proxy := browser.NewProxy(browser.ProxyConfig{
				ListenAddr:     record.ListenAddr,
				Upstream:       record.UpstreamWSURL,
				OwnershipToken: record.OwnershipToken,
			})
			return proxy.ServeListener(ctx, ln)
		},
	}
}
