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

This command establishes a persistent connection to your existing Chrome/Chromium.

Once attached, all tap browser commands will use this context.

Examples:
  tap attach chrome                    Auto-discover your Chrome
  tap attach chrome --browser-url http://localhost:9222
  tap attach status                    Show attachment info
  tap attach clear                     Detach from external target`,
		Commands: []*cli.Command{
			attachChromeCmd(),
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
				Name:   "listen",
				Usage:  "Internal proxy listen address (advanced)",
				Value:  browser.DefaultProxyListenAddr,
				Hidden: true,
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

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}
	store, err := newBrowserStore(cmd)
	if err != nil {
		return err
	}

	proxyRecord, health, reused, err := ensureAttachedChromeProxy(ctx, cmd, wsURL, cmd.String("listen"))
	if err != nil {
		return fmt.Errorf("ensure proxy daemon: %w", err)
	}

	sessionName := browser.DefaultSessionName
	if _, err := mgr.GetSession(ctx, sessionName); err == nil {
		if err := mgr.CloseSession(ctx, sessionName); err != nil {
			return fmt.Errorf("close existing session: %w", err)
		}
	}
	if err := mgr.CreateSession(ctx, sessionName, browser.ModeRemote, browser.SessionOptions{WSURL: proxyRecord.Endpoint}); err != nil {
		return fmt.Errorf("create attached session: %w", err)
	}
	if err := mgr.SetDefaultContext(ctx, sessionName, browser.DefaultContextAttached); err != nil {
		return fmt.Errorf("set default context: %w", err)
	}
	if err := persistProxyDaemonRecord(store, proxyRecord, health); err != nil {
		return fmt.Errorf("persist proxy daemon state: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Attached to Chrome\n", green(c, "✓"))
	fmt.Fprintf(os.Stderr, "  Source: %s\n", source)
	fmt.Fprintf(os.Stderr, "  Session: %s\n", sessionName)
	fmt.Fprintf(os.Stderr, "  Proxy: %s\n", proxyRecord.Endpoint)
	if reused {
		fmt.Fprintf(os.Stderr, "  Daemon: reused\n")
	} else {
		fmt.Fprintf(os.Stderr, "  Daemon: started (pid %d)\n", proxyRecord.PID)
	}

	// Auto-discover/adopt tabs from the attached browser.
	if err := autoAdoptTabs(ctx, mgr, sessionName, proxyRecord.Endpoint); err != nil {
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
	store, err := newBrowserStore(cmd)
	if err != nil {
		return err
	}
	state, err := store.Load()
	if err != nil {
		return err
	}

	defaultContext, _ := mgr.DefaultContext(ctx)
	var session *browser.SessionRecord
	if current, err := mgr.GetSession(ctx, ""); err == nil {
		session = current
	}

	var health browser.ProxyDaemonHealth
	if state.ProxyDaemon != nil {
		health = browser.CheckProxyDaemon(ctx, state.ProxyDaemon)
		_ = persistProxyDaemonHealth(store, state.ProxyDaemon, health)
	}

	if session == nil && state.ProxyDaemon == nil {
		if cmd.Bool("json") {
			return printAttachStatusJSON(defaultContext, nil, nil, browser.ProxyDaemonHealth{}, "no_session")
		}
		fmt.Println("No attachment active.")
		fmt.Println()
		fmt.Println("To attach:")
		fmt.Println("  tap attach chrome")
		return nil
	}

	if cmd.Bool("json") {
		return printAttachStatusJSON(defaultContext, session, state.ProxyDaemon, health, "")
	}
	return printAttachStatusHuman(defaultContext, session, state.ProxyDaemon, health)
}

func printAttachStatusHuman(defaultContext *browser.DefaultContextRecord, session *browser.SessionRecord, proxy *browser.ProxyDaemonRecord, health browser.ProxyDaemonHealth) error {
	c := true
	if defaultContext != nil {
		fmt.Printf("%s %s (%s)\n", bold(c, "Default context:"), defaultContext.SessionName, defaultContext.Kind)
		if defaultContext.Stale {
			fmt.Printf("%s %s\n", bold(c, "State:"), yellow(c, "stale"))
			if defaultContext.Reason != "" {
				fmt.Printf("%s %s\n", bold(c, "Reason:"), defaultContext.Reason)
			}
		}
	}
	if session != nil {
		fmt.Printf("%s %s\n", bold(c, "Session:"), session.Name)
		fmt.Printf("%s %s\n", bold(c, "Mode:"), session.Mode)
		if session.Remote != nil {
			fmt.Printf("%s %s\n", bold(c, "Remote URL:"), session.Remote.WSURL)
		}
	}
	if proxy != nil {
		fmt.Println()
		fmt.Printf("%s %s\n", bold(c, "Proxy daemon:"), proxy.Endpoint)
		fmt.Printf("%s %d\n", bold(c, "PID:"), proxy.PID)
		fmt.Printf("%s %s\n", bold(c, "Listen:"), proxy.ListenAddr)
		fmt.Printf("%s %s\n", bold(c, "Upstream:"), proxy.UpstreamWSURL)
		statusColor := green(c, health.Status)
		if !health.Healthy {
			statusColor = yellow(c, health.Status)
		}
		fmt.Printf("%s %s\n", bold(c, "Health:"), statusColor)
		if health.Reason != "" {
			fmt.Printf("%s %s\n", bold(c, "Stale reason:"), health.Reason)
		}
	}
	if session != nil && len(session.Tabs) > 0 {
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

func printAttachStatusJSON(defaultContext *browser.DefaultContextRecord, session *browser.SessionRecord, proxy *browser.ProxyDaemonRecord, health browser.ProxyDaemonHealth, errorState string) error {
	result := map[string]any{"error": errorState}
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
			"name":        session.Name,
			"mode":        session.Mode,
			"selectedTab": session.SelectedTab,
			"tabs":        len(session.Tabs),
		}
		if session.Remote != nil {
			sessionMap["remoteURL"] = session.Remote.WSURL
		}
		if session.Process != nil {
			sessionMap["processURL"] = session.Process.DebugURL
		}
		result["session"] = sessionMap
	}
	if proxy != nil {
		result["proxyDaemon"] = map[string]any{
			"pid":             proxy.PID,
			"listenAddr":      proxy.ListenAddr,
			"endpoint":        proxy.Endpoint,
			"upstreamWSURL":   proxy.UpstreamWSURL,
			"status":          health.Status,
			"healthy":         health.Healthy,
			"reason":          health.Reason,
			"lastHealthCheck": health.CheckedAt,
		}
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
	store, err := newBrowserStore(cmd)
	if err != nil {
		return err
	}
	state, err := store.Load()
	if err != nil {
		return err
	}

	defaultContext := state.DefaultContext
	if defaultContext == nil && state.ProxyDaemon == nil {
		fmt.Println("No attachment to clear.")
		return nil
	}
	if defaultContext != nil {
		if session, err := mgr.GetSession(ctx, defaultContext.SessionName); err == nil && session.Remote != nil {
			if err := mgr.CloseSession(ctx, defaultContext.SessionName); err != nil {
				return fmt.Errorf("close session: %w", err)
			}
		}
	}
	if state.ProxyDaemon != nil {
		_ = stopOwnedProxyDaemon(ctx, state.ProxyDaemon)
	}
	if err := store.Update(func(state *browser.State) error {
		state.ClearDefaultContext()
		state.ProxyDaemon = nil
		return nil
	}); err != nil {
		return fmt.Errorf("clear attachment metadata: %w", err)
	}
	c := true
	fmt.Fprintf(os.Stderr, "%s Cleared attached Chrome metadata\n", green(c, "✓"))
	return nil
}
