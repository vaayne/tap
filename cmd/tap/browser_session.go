package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ws-url",
				Usage: "Remote CDP WebSocket URL; omit for a managed local browser",
			},
			&cli.BoolFlag{
				Name:  "no-headless",
				Usage: "Run local browser in visible mode",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			name := cmd.Args().First()
			if name == "" {
				return fmt.Errorf("session name required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			mode := browser.ModeLocal
			opts := browser.SessionOptions{Headless: !cmd.Bool("no-headless")}
			if wsURL := cmd.String("ws-url"); wsURL != "" {
				mode = browser.ModeRemote
				opts.WSURL = wsURL
			}
			if err := mgr.CreateSession(ctx, name, mode, opts); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Session %q created (%s)\n", name, mode)
			return nil
		},
	}
}

func browserSessionListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List tracked browser sessions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			list, err := mgr.ListSessions(ctx)
			if err != nil {
				return err
			}
			if len(list.Sessions) == 0 {
				fmt.Fprintln(os.Stderr, "No browser sessions found.")
				return nil
			}

			c := useColor(cmd)
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, bold(c, "NAME\tMODE\tTABS\tSELECTED"))
			for _, s := range list.Sessions {
				sel := ""
				if s.Name == list.SelectedSession {
					sel = green(c, "*")
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", s.Name, s.Mode, len(s.Tabs), sel)
			}
			return w.Flush()
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
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			name := cmd.Args().First()
			session, err := mgr.GetSession(ctx, name)
			if err != nil {
				return err
			}

			c := useColor(cmd)
			_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", bold(c, "Name:"), session.Name)
			_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", bold(c, "Mode:"), session.Mode)
			debugURL := ""
			if session.Process != nil && session.Process.DebugURL != "" {
				debugURL = session.Process.DebugURL
			} else if session.Remote != nil {
				debugURL = session.Remote.WSURL
			}
			if debugURL != "" {
				_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", bold(c, "Debug URL:"), debugURL)
			}
			_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", bold(c, "Created:"), session.CreatedAt.Format("2006-01-02 15:04:05"))

			if len(session.Tabs) == 0 {
				_, _ = fmt.Fprintf(os.Stdout, "\n%s\n", dim(c, "No tracked tabs."))
				return nil
			}

			_, _ = fmt.Fprintf(os.Stdout, "\n%s\n", bold(c, "Tabs:"))
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "  NAME\tSTATUS\tURL")
			for _, tab := range session.Tabs {
				status := string(tab.Status)
				switch tab.Status {
				case browser.TabStatusLive:
					status = green(c, status)
				case browser.TabStatusStale:
					status = yellow(c, status)
				}
				sel := ""
				if tab.Name == session.SelectedTab {
					sel = " " + green(c, "*")
				}
				_, _ = fmt.Fprintf(w, "  %s%s\t%s\t%s\n", tab.Name, sel, status, tab.URL)
			}
			return w.Flush()
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
			configureLogging(cmd)
			name := cmd.Args().First()
			if name == "" {
				return fmt.Errorf("session name required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.SelectSession(ctx, name); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Session %q selected\n", name)
			return nil
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
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			name := cmd.Args().First()
			if err := mgr.CloseSession(ctx, name); err != nil {
				return err
			}
			if name == "" {
				name = "(resolved)"
			}
			fmt.Fprintf(os.Stderr, "Session %q closed\n", name)
			return nil
		},
	}
}
