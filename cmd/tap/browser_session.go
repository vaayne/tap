package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func browserSessionCmd() *cli.Command {
	return &cli.Command{
		Name:  "session",
		Usage: "Create, inspect, select, list, and close persistent browser sessions",
		Description: `A session is one persistent browser instance managed by tap.

Local sessions own a dedicated Chrome profile directory and reconnect through
stored launch metadata. Remote sessions are bound to the explicit --ws-url used
when they are created; later global --ws-url overrides are ignored so reconnects
stay deterministic.`,
		Commands: []*cli.Command{
			browserSessionNewCmd(),
			browserSessionListCmd(),
			browserSessionInfoCmd(),
			browserSessionSelectCmd(),
			browserSessionCloseCmd(),
		},
	}
}

func browserSessionNewCmd() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a persistent browser session",
		ArgsUsage: "<name>",
		Description: `Create a local or remote persistent session.

Without --ws-url, tap will manage a local browser process for the session.
With --ws-url, tap persists the remote endpoint in session metadata and validates
connection/auth/TLS at creation time. Later commands reconnect through the saved
endpoint rather than a new global --ws-url value.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("session name required")
			}
			return ensureBrowserPhase3("tap browser session new", cmd)
		},
	}
}

func browserSessionListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List tracked browser sessions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser session list", cmd)
		},
	}
}

func browserSessionInfoCmd() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Show session metadata and tracked-tab status",
		ArgsUsage: "[name]",
		Description: `Show the persisted session configuration, selected-tab state, and tracked tabs.

Tracked tabs are reported as live, stale, or closed. Untracked live browser tabs
are not auto-adopted in v1.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser session info", cmd)
		},
	}
}

func browserSessionSelectCmd() *cli.Command {
	return &cli.Command{
		Name:      "select",
		Usage:     "Set the default browser session",
		ArgsUsage: "<name>",
		Description: `Persist the default session used by browser commands when --session is omitted.

If no session is selected and more than one session exists, browser actions fail
with guidance instead of guessing.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("session name required")
			}
			return ensureBrowserPhase3("tap browser session select", cmd)
		},
	}
}

func browserSessionCloseCmd() *cli.Command {
	return &cli.Command{
		Name:      "close",
		Usage:     "Close a persistent browser session",
		ArgsUsage: "[name]",
		Description: `Close the resolved session.

For local sessions, tap will terminate only the browser process it owns after
verifying launch ownership, then remove the session metadata and managed profile
directory. For remote sessions, tap removes its metadata only and never attempts
to kill the remote browser process.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser session close", cmd)
		},
	}
}
