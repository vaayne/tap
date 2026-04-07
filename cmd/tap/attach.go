package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func attachCmd() *cli.Command {
	return &cli.Command{
		Name:  "attach",
		Usage: "Connect tap to existing Chrome or Electron apps",
		Description: `Attach tap to an already-running browser or app for reuse.

This command establishes a persistent connection to:
- Your existing Chrome/Chromium (with remote debugging enabled)
- An Electron app with CDP debugging enabled

Once attached, all tap browser commands will use this context.

Examples:
  tap attach chrome                    Auto-discover your Chrome
  tap attach chrome --browser-url http://localhost:9222
  tap attach electron --port 9229      Attach to running Electron app
  tap attach electron --launch /path/to/app
  tap attach status                    Show attachment info
  tap attach clear                     Detach from external target`,
		Commands: []*cli.Command{
			attachChromeCmd(),
			attachElectronCmd(),
			attachStatusCmd(),
			attachClearCmd(),
		},
	}
}

func attachChromeCmd() *cli.Command {
	return &cli.Command{
		Name:  "chrome",
		Usage: "Attach to your existing Chrome/Chromium browser",
		Description: `Attach tap to an already-running Chrome or Chromium instance.

Your browser must have remote debugging enabled. Common ways:
- Chrome started with --remote-debugging-port=9222
- Chrome with "Developer Tools" open in user data dir

Examples:
  tap attach chrome                    Auto-discover from DevToolsActivePort
  tap attach chrome --browser-url http://localhost:9222
  tap attach chrome --port-file ~/Library/Application\ Support/Google/Chrome/DevToolsActivePort`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "browser-url",
				Aliases: []string{"url", "ws-url"},
				Usage:   "Chrome DevTools endpoint URL (WebSocket or HTTP base URL)",
			},
			&cli.StringFlag{
				Name:  "port-file",
				Usage: "Path to DevToolsActivePort file",
			},
			&cli.StringFlag{
				Name:  "listen",
				Usage: "Internal proxy listen address (advanced)",
				Value: browser.DefaultProxyListenAddr,
			},
			// Compatibility aliases
			&cli.StringFlag{
				Name:   "user-chrome",
				Usage:  "Compatibility: auto-discover user Chrome (now default)",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "devtools-port-file",
				Usage:  "Compatibility: explicit port file",
				Hidden: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runAttachChrome(ctx, cmd)
		},
	}
}

func runAttachChrome(ctx context.Context, cmd *cli.Command) error {
	var debugURL string
	var source string

	// Resolve debug URL from various sources
	if url := cmd.String("browser-url"); url != "" {
		debugURL = url
		source = "explicit URL"
	} else if path := cmd.String("port-file"); path != "" {
		url, err := browser.ResolveDebugURLFromDevToolsFile(path)
		if err != nil {
			return fmt.Errorf("resolve port file: %w", err)
		}
		debugURL = url
		source = fmt.Sprintf("port file: %s", path)
	} else if path := cmd.String("devtools-port-file"); path != "" {
		url, err := browser.ResolveDebugURLFromDevToolsFile(path)
		if err != nil {
			return fmt.Errorf("resolve port file: %w", err)
		}
		debugURL = url
		source = fmt.Sprintf("port file: %s", path)
	} else {
		// Auto-discover
		url, path, err := browser.DiscoverUserChromeDebugURL()
		if err != nil {
			return fmt.Errorf(`could not find a running Chrome with remote debugging enabled.

To use tap attach chrome, you need:
1. Chrome already running with --remote-debugging-port=9222, OR
2. Chrome with DevToolsActivePort in a standard location

Try one of:
  tap attach chrome --browser-url http://localhost:9222
  google-chrome --remote-debugging-port=9222 &
  tap attach chrome

Original error: %w`, err)
		}
		debugURL = url
		source = fmt.Sprintf("auto-discovered: %s", path)
	}

	// Resolve to WebSocket URL if HTTP base URL provided
	wsURL, err := browser.ResolveDebugURL(ctx, debugURL)
	if err != nil {
		return fmt.Errorf("resolve debug URL: %w", err)
	}

	// Create manager
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	// Create or update default session as remote attached session
	sessionName := browser.DefaultSessionName
	opts := browser.SessionOptions{WSURL: wsURL}

	// Check if session exists
	_, err = mgr.GetSession(ctx, sessionName)
	if err == nil {
		// Session exists - close it first
		if err := mgr.CloseSession(ctx, sessionName); err != nil {
			return fmt.Errorf("close existing session: %w", err)
		}
	}

	// Create new remote session
	if err := mgr.CreateSession(ctx, sessionName, browser.ModeRemote, opts); err != nil {
		return fmt.Errorf("create attached session: %w", err)
	}
	if err := mgr.SetDefaultContext(ctx, sessionName, browser.DefaultContextAttached); err != nil {
		return fmt.Errorf("set default context: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Attached to Chrome\n", green(c, "✓"))
	fmt.Fprintf(os.Stderr, "  Source: %s\n", source)
	fmt.Fprintf(os.Stderr, "  Session: %s\n", sessionName)

	// Auto-discover/adopt tabs from the attached browser
	if err := autoAdoptTabs(ctx, mgr, sessionName, wsURL); err != nil {
		// Non-fatal - just warn
		fmt.Fprintf(os.Stderr, "  Warning: could not auto-adopt tabs: %v\n", err)
	}

	return nil
}

func autoAdoptTabs(ctx context.Context, mgr *browser.Manager, sessionName, debugURL string) error {
	// List live targets via HTTP
	targets, err := browser.ListTargetsHTTP(ctx, debugURL)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return nil
	}

	// Adopt first target as current tab
	for i, t := range targets {
		tabName := fmt.Sprintf("tab-%d", i+1)
		if err := mgr.AdoptTab(ctx, sessionName, tabName, t.TargetID, t.URL); err != nil {
			continue // skip on error
		}
		if i == 0 {
			// Select first as current
			_ = mgr.SelectTab(ctx, sessionName, tabName)
		}
	}

	return nil
}

func attachElectronCmd() *cli.Command {
	return &cli.Command{
		Name:  "electron",
		Usage: "Attach to a running Electron app or launch one",
		Description: `Connect tap to an Electron app via Chrome DevTools Protocol.

The Electron app must expose a browser CDP endpoint (not Node inspector).

Examples:
  # Attach to already-running app
  tap attach electron --port 9229

  # Launch and attach
  tap attach electron --launch /Applications/MyApp.app/Contents/MacOS/MyApp`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Usage: "CDP debug port the Electron app is listening on",
			},
			&cli.StringFlag{
				Name:  "launch",
				Usage: "Path to Electron app binary to launch with debugging",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"session", "s"},
				Usage:   "Session name for this attachment",
				Value:   "electron",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runAttachElectron(ctx, cmd)
		},
	}
}

func runAttachElectron(ctx context.Context, cmd *cli.Command) error {
	port := int(cmd.Int("port"))
	launchPath := cmd.String("launch")
	sessionName := cmd.String("name")

	if port == 0 && launchPath == "" {
		return fmt.Errorf("either --port or --launch required")
	}

	var wsURL string
	var err error

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	if port > 0 {
		// Attach to existing
		wsURL, err = browser.ResolveElectronDebugURL(ctx, port)
		if err != nil {
			return fmt.Errorf("attach to port %d: %w", port, err)
		}
	} else {
		// Launch and attach
		// Extract app args after --
		args := cmd.Args().Slice()
		proc, err := browser.LaunchElectronApp(ctx, launchPath, args)
		if err != nil {
			return fmt.Errorf("launch: %w", err)
		}
		wsURL = proc.DebugURL
		// Port not directly available on ProcessRecord, but DebugURL contains it
	}

	// Close existing session if present
	_, err = mgr.GetSession(ctx, sessionName)
	if err == nil {
		if err := mgr.CloseSession(ctx, sessionName); err != nil {
			return fmt.Errorf("close existing session: %w", err)
		}
	}

	// Create remote session
	opts := browser.SessionOptions{WSURL: wsURL}
	if err := mgr.CreateSession(ctx, sessionName, browser.ModeRemote, opts); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	if err := mgr.SetDefaultContext(ctx, sessionName, browser.DefaultContextAttached); err != nil {
		return fmt.Errorf("set default context: %w", err)
	}

	// Auto-adopt tabs
	if err := autoAdoptTabs(ctx, mgr, sessionName, wsURL); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not auto-adopt tabs: %v\n", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Attached to Electron app\n", green(c, "✓"))
	fmt.Fprintf(os.Stderr, "  Session: %s\n", sessionName)

	return nil
}

func attachStatusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show attachment status",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runAttachStatus(ctx, cmd)
		},
	}
}

func runAttachStatus(ctx context.Context, cmd *cli.Command) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	defaultContext, _ := mgr.DefaultContext(ctx)
	session, err := mgr.GetSession(ctx, "")
	if err != nil {
		if cmd.Bool("json") {
			return printAttachStatusJSON(defaultContext, nil, "no_session")
		}
		fmt.Println("No attachment active.")
		fmt.Println()
		fmt.Println("To attach:")
		fmt.Println("  tap attach chrome")
		fmt.Println("  tap attach electron --port 9229")
		return nil
	}

	if cmd.Bool("json") {
		return printAttachStatusJSON(defaultContext, session, "")
	}

	return printAttachStatusHuman(defaultContext, session)
}

func printAttachStatusHuman(defaultContext *browser.DefaultContextRecord, session *browser.SessionRecord) error {
	c := true

	contextType := "Managed local browser"
	if defaultContext != nil {
		fmt.Printf("%s %s (%s)\n", bold(c, "Default context:"), defaultContext.SessionName, defaultContext.Kind)
		if defaultContext.Stale {
			fmt.Printf("%s %s\n", bold(c, "State:"), yellow(c, "stale"))
		}
	}
	if session.Remote != nil {
		contextType = "Attached remote browser"
	}

	fmt.Printf("%s %s\n", bold(c, "Context type:"), contextType)
	fmt.Printf("%s %s\n", bold(c, "Session:"), session.Name)
	fmt.Printf("%s %s\n", bold(c, "Mode:"), session.Mode)

	if session.Process != nil && session.Process.DebugURL != "" {
		fmt.Printf("%s %s\n", bold(c, "Process:"), session.Process.DebugURL)
	}
	if session.Remote != nil {
		fmt.Printf("%s %s\n", bold(c, "Remote URL:"), session.Remote.WSURL)
	}

	if len(session.Tabs) > 0 {
		fmt.Println()
		fmt.Printf("%s %d\n", bold(c, "Adopted tabs:"), len(session.Tabs))
		for name, tab := range session.Tabs {
			marker := ""
			if name == session.SelectedTab {
				marker = green(c, " (current)")
			}
			fmt.Printf("  %s: %s%s\n", name, tab.URL, marker)
		}
	}

	return nil
}

func printAttachStatusJSON(defaultContext *browser.DefaultContextRecord, session *browser.SessionRecord, errorState string) error {
	result := map[string]any{
		"error": errorState,
	}

	if defaultContext != nil {
		result["defaultContext"] = map[string]any{
			"sessionName": defaultContext.SessionName,
			"kind":        defaultContext.Kind,
			"stale":       defaultContext.Stale,
			"reason":      defaultContext.Reason,
		}
	}

	if session != nil {
		sessionMap := map[string]any{
			"name": session.Name,
			"mode": session.Mode,
		}
		if session.Remote != nil {
			sessionMap["remoteURL"] = session.Remote.WSURL
		}
		if session.Process != nil {
			sessionMap["processURL"] = session.Process.DebugURL
		}
		result["session"] = sessionMap
		result["tabs"] = len(session.Tabs)
		result["selectedTab"] = session.SelectedTab
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	return nil
}

func attachClearCmd() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Detach from external browser/app",
		Description: `Remove the attachment metadata.

This does not close the external browser or app, nor does it delete
cookies. It only stops tap from using this target as the default context.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runAttachClear(ctx, cmd)
		},
	}
}

func runAttachClear(ctx context.Context, cmd *cli.Command) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	defaultContext, err := mgr.DefaultContext(ctx)
	if err != nil || defaultContext == nil {
		fmt.Println("No attachment to clear.")
		return nil
	}

	session, err := mgr.GetSession(ctx, defaultContext.SessionName)
	if err != nil {
		_ = mgr.ClearDefaultContext(ctx)
		fmt.Println("Cleared stale attachment metadata.")
		return nil
	}

	// Only clear if it's a remote session (attached, not managed)
	if session.Remote != nil {
		if err := mgr.CloseSession(ctx, defaultContext.SessionName); err != nil {
			return fmt.Errorf("close session: %w", err)
		}
		_ = mgr.ClearDefaultContext(ctx)
		c := true
		fmt.Fprintf(os.Stderr, "%s Detached from external browser\n", green(c, "✓"))
	} else {
		fmt.Println("Default session is managed locally, not attached.")
		fmt.Println("To close: tap browser session close default")
	}

	return nil
}
