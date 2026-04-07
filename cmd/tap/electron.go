package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func electronCmd() *cli.Command {
	return &cli.Command{
		Name:  "electron",
		Usage: "Debug Electron apps via Chrome DevTools Protocol",
		Description: `Connect tap to Electron apps for automation and debugging.

Electron exposes the same CDP protocol as Chrome. Use these commands to
discover running Electron apps, attach to them, or launch them with
debugging enabled. After connecting, all 'tap browser' commands work
against the session.

Examples:
  tap electron ps                              List running debuggable processes
  tap electron attach --port 9229              Attach to app on port 9229
  tap electron launch /path/to/MyApp           Launch app with debugging`,
		Commands: []*cli.Command{
			electronPsCmd(),
			electronAttachCmd(),
			electronLaunchCmd(),
			electronDiscoverCmd(),
		},
	}
}

func electronPsCmd() *cli.Command {
	return &cli.Command{
		Name:  "ps",
		Usage: "List running processes with CDP debug ports",
		Description: `Scans running processes for --remote-debugging-port arguments.

Matches Electron apps, Chrome instances, and any CEF-based app launched
with remote debugging enabled. Use the PORT shown with 'tap electron attach'.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			procs, err := browser.ScanElectronProcesses(ctx)
			if err != nil {
				return err
			}
			if len(procs) == 0 {
				fmt.Fprintln(os.Stderr, "No processes with --remote-debugging-port found.")
				fmt.Fprintln(os.Stderr, "Launch your app with --remote-debugging-port=PORT, then run 'tap electron attach'.")
				return nil
			}
			c := useColor(cmd)
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, bold(c, "PID\tPORT\tNAME"))
			for _, p := range procs {
				_, _ = fmt.Fprintf(w, "%d\t%d\t%s\n", p.PID, p.Port, p.Name)
			}
			return w.Flush()
		},
	}
}

func electronAttachCmd() *cli.Command {
	return &cli.Command{
		Name:  "attach",
		Usage: "Attach to a running Electron app via its debug port",
		Description: `Resolves the WebSocket debug URL from /json/version and creates a
named browser session bound to that endpoint.

The Electron app must have been started with --remote-debugging-port=PORT.
Once attached, use 'tap browser' commands with --session to control it.

Examples:
  tap electron attach --port 9229
  tap electron attach --port 9229 --session myapp
  tap browser tab list --session myapp
  tap browser screenshot --session myapp`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "CDP debug port the Electron app is listening on",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Usage:   "Session name to create",
				Value:   "electron",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			port := int(cmd.Int("port"))
			sessionName := cmd.String("session")

			wsURL, err := browser.ResolveElectronDebugURL(ctx, port)
			if err != nil {
				return fmt.Errorf("attach: %w", err)
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.CreateSession(ctx, sessionName, browser.ModeRemote, browser.SessionOptions{WSURL: wsURL}); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Session %q attached to Electron app at port %d\n", sessionName, port)
			fmt.Fprintf(os.Stderr, "Debug URL: %s\n", wsURL)
			return nil
		},
	}
}

func electronLaunchCmd() *cli.Command {
	return &cli.Command{
		Name:      "launch",
		Usage:     "Launch an Electron app with CDP debugging enabled",
		ArgsUsage: "<binary> [app-args...]",
		Description: `Launches the given binary with --remote-debugging-port=0 (OS assigns the port)
and creates a named session once the debug URL is available.

For macOS .app bundles, pass the inner binary path, not the .app directory:
  /Applications/MyApp.app/Contents/MacOS/MyApp

All arguments after the binary are forwarded to the app unchanged.

Examples:
  tap electron launch /Applications/MyApp.app/Contents/MacOS/MyApp
  tap electron launch electron --session dev -- /path/to/app
  tap browser screenshot --session dev`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Usage:   "Session name to create",
				Value:   "electron",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) == 0 {
				return fmt.Errorf("binary path required")
			}
			binaryPath := args[0]
			extra := args[1:]
			sessionName := cmd.String("session")

			// Resolve binary from PATH when not an absolute path.
			if !filepath.IsAbs(binaryPath) {
				resolved, err := exec.LookPath(binaryPath)
				if err != nil {
					return fmt.Errorf("binary %q not found: %w", binaryPath, err)
				}
				binaryPath = resolved
			}

			proc, err := browser.LaunchElectronApp(ctx, binaryPath, extra)
			if err != nil {
				return fmt.Errorf("launch: %w", err)
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.CreateSession(ctx, sessionName, browser.ModeRemote, browser.SessionOptions{WSURL: proc.DebugURL}); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Session %q created for Electron process %d\n", sessionName, proc.PID)
			fmt.Fprintf(os.Stderr, "Debug URL: %s\n", proc.DebugURL)
			return nil
		},
	}
}

func electronDiscoverCmd() *cli.Command {
	return &cli.Command{
		Name:  "discover",
		Usage: "Adopt live Electron windows as tracked tabs",
		Description: `Lists all live page targets in the session's Electron app and registers
each one as a tracked tab. Targets already tracked under a different name
are skipped. Tab names are derived from the window title (or 'window-N'
when the title is empty).

Run this after 'tap electron attach' or 'tap electron launch' to make
windows accessible to 'tap browser' commands.

Examples:
  tap electron discover --session electron
  tap browser screenshot --session electron`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Usage:   "Session name to discover into",
				Value:   "electron",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sessionName := cmd.String("session")

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			// Resolve session to get debug URL.
			session, err := mgr.GetSession(ctx, sessionName)
			if err != nil {
				return err
			}
			debugURL := ""
			if session.Process != nil && session.Process.DebugURL != "" {
				debugURL = session.Process.DebugURL
			} else if session.Remote != nil {
				debugURL = session.Remote.WSURL
			}
			if debugURL == "" {
				return fmt.Errorf("session %q has no debug URL", sessionName)
			}

			// List live targets via HTTP — more compatible with Electron apps
			// that do not support Target.getTargets over the browser WebSocket.
			targets, err := browser.ListTargetsHTTP(ctx, debugURL)
			if err != nil {
				return fmt.Errorf("list targets: %w", err)
			}
			if len(targets) == 0 {
				fmt.Fprintln(os.Stderr, "No live page targets found in session.")
				return nil
			}

			// Build set of already-tracked target IDs.
			tracked := make(map[string]bool, len(session.Tabs))
			for _, tab := range session.Tabs {
				tracked[tab.TargetID] = true
			}

			adopted := 0
			for i, t := range targets {
				if tracked[t.TargetID] {
					continue
				}
				// Derive a safe tab name from the window title.
				tabName := sanitizeTabName(t.Title)
				if tabName == "" {
					tabName = fmt.Sprintf("window-%d", i+1)
				}
				// Ensure uniqueness by appending index if needed.
				tabName = uniqueTabName(tabName, session.Tabs)

				if err := mgr.AdoptTab(ctx, sessionName, tabName, t.TargetID, t.URL); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not adopt target %s (%s): %v\n", t.TargetID, t.Title, err)
					continue
				}
				fmt.Fprintf(os.Stderr, "Adopted %q → tab %q (%s)\n", t.Title, tabName, t.URL)
				adopted++
			}

			if adopted == 0 {
				fmt.Fprintln(os.Stderr, "All live targets are already tracked.")
			} else {
				fmt.Fprintf(os.Stderr, "%d window(s) adopted. Use 'tap browser tab list --session %s' to see them.\n", adopted, sessionName)
			}
			return nil
		},
	}
}

// sanitizeTabName converts an arbitrary string to a valid session name by
// replacing spaces and special characters with hyphens and lowercasing.
func sanitizeTabName(title string) string {
	if title == "" {
		return ""
	}
	b := make([]byte, 0, len(title))
	for i := 0; i < len(title); i++ {
		c := title[i]
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
			b = append(b, c)
		case c >= 'A' && c <= 'Z':
			b = append(b, c+32) // lowercase
		case c == '-', c == '.', c == '_':
			b = append(b, c)
		default:
			if len(b) > 0 && b[len(b)-1] != '-' {
				b = append(b, '-')
			}
		}
	}
	// Trim trailing hyphen.
	for len(b) > 0 && b[len(b)-1] == '-' {
		b = b[:len(b)-1]
	}
	if len(b) > 32 {
		b = b[:32]
	}
	return string(b)
}

// uniqueTabName appends a numeric suffix to name until it is not present in tabs.
func uniqueTabName(name string, tabs map[string]*browser.TabRecord) string {
	if _, exists := tabs[name]; !exists {
		return name
	}
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d", name, i)
		if _, exists := tabs[candidate]; !exists {
			return candidate
		}
	}
}
