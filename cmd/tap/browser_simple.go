package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

// browserSimpleCmd returns commands that provide the simplified browser UX
// These are aliases or wrappers around existing browser session/tab commands
func browserOpenCmd() *cli.Command {
	return &cli.Command{
		Name:      "open",
		Usage:     "Open or navigate to a URL in the current browser context",
		ArgsUsage: "<url>",
		Description: `Open a URL in the default browser context.

If no current tab exists, creates one. If a current tab exists, navigates it.
Use --new-tab to force creation of a new tab.

Examples:
  tap browser open https://example.com
  tap browser open https://github.com --new-tab
  tap browser open https://example.com --show`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name (default: use active context)",
			},
			&cli.BoolFlag{
				Name:  "new-tab",
				Usage: "Create a new tab instead of navigating current tab",
			},
			&cli.BoolFlag{
				Name:    "show",
				Aliases: []string{"no-headless"},
				Usage:   "Run browser in visible mode",
			},
			&cli.DurationFlag{
				Name:  "wait",
				Usage: "Wait duration after navigation",
			},
			&cli.StringFlag{
				Name:  "wait-selector",
				Usage: "Wait until selector is visible",
			},
			&cli.DurationFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "Global timeout",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserOpen(ctx, cmd)
		},
	}
}

func runBrowserOpen(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("URL required")
	}
	url := cmd.Args().First()

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	sessionName := cmd.String("session")
	if sessionName == "" {
		sessionName = browser.DefaultSessionName
	}

	// Get or create session
	session, err := mgr.GetSession(ctx, sessionName)
	if err != nil {
		// Session doesn't exist - create managed local session
		opts := browser.SessionOptions{Headless: !cmd.Bool("show")}
		if err := mgr.CreateSession(ctx, sessionName, browser.ModeLocal, opts); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		session, err = mgr.GetSession(ctx, sessionName)
		if err != nil {
			return err
		}
	}

	// Determine if we need a new tab
	createNewTab := cmd.Bool("new-tab")
	var targetTab string

	if !createNewTab && session.SelectedTab != "" {
		// Check if selected tab is live
		if tab, ok := session.Tabs[session.SelectedTab]; ok && tab.Status == browser.TabStatusLive {
			targetTab = session.SelectedTab
		}
	}

	if targetTab == "" {
		// Need to create a new tab
		targetTab = generateNextTabName(session)
		if err := mgr.CreateTab(ctx, sessionName, targetTab, "about:blank"); err != nil {
			return fmt.Errorf("create tab: %w", err)
		}
		// Select it
		if err := mgr.SelectTab(ctx, sessionName, targetTab); err != nil {
			return err
		}
	}

	// Navigate the tab
	if err := mgr.Navigate(ctx, sessionName, targetTab, url); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

	// Handle wait flags
	if d := cmd.Duration("wait"); d > 0 {
		// Simple sleep - actual implementation would need proper wait
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(d):
		}
	}

	// Handle wait-selector if provided (simplified for Phase 1)
	if sel := cmd.String("wait-selector"); sel != "" {
		// TODO: Implement proper wait using Evaluate with polling
		// For now, just do a fixed delay as a placeholder
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Opened %s\n", green(c, "✓"), url)
	if targetTab != session.SelectedTab || createNewTab {
		fmt.Fprintf(os.Stderr, "  Tab: %s\n", targetTab)
	}

	return nil
}

func generateNextTabName(session *browser.SessionRecord) string {
	// Find highest tab-N number
	maxNum := 0
	for name := range session.Tabs {
		var num int
		if _, err := fmt.Sscanf(name, "tab-%d", &num); err == nil {
			if num > maxNum {
				maxNum = num
			}
		}
	}
	return fmt.Sprintf("tab-%d", maxNum+1)
}

// browserTabsCmd is a user-friendly alias for "browser tab list"
func browserTabsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tabs",
		Usage: "List browser tabs",
		Description: `Show all tabs in the current browser context.

This is a simplified view of tracked tabs with stable IDs.

Examples:
  tap browser tabs
  tap browser tabs --json`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserTabs(ctx, cmd)
		},
	}
}

func runBrowserTabs(ctx context.Context, cmd *cli.Command) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	sessionName := cmd.String("session")
	list, err := mgr.ListTabs(ctx, sessionName)
	if err != nil {
		return err
	}

	if cmd.Bool("json") {
		return printTabsJSON(list)
	}

	return printTabsHuman(list)
}

func printTabsHuman(list *browser.TabList) error {
	if len(list.Tabs) == 0 {
		fmt.Println("No tabs found.")
		fmt.Println("Run: tap browser open <url>")
		return nil
	}

	c := true
	fmt.Printf("%s %d tabs\n\n", bold(c, "Tabs:"), len(list.Tabs))

	fmt.Printf("%-10s %-20s %-40s %s\n", "ID", "STATUS", "URL", "CURRENT")
	fmt.Println(string(make([]byte, 80))) // separator line

	for _, tab := range list.Tabs {
		current := ""
		if tab.Name == list.SelectedTab {
			current = green(c, "*")
		}

		status := string(tab.Status)
		switch tab.Status {
		case browser.TabStatusLive:
			status = green(c, status)
		case browser.TabStatusStale:
			status = yellow(c, status)
		}

		url := tab.URL
		if len(url) > 38 {
			url = url[:35] + "..."
		}

		fmt.Printf("%-10s %-20s %-40s %s\n", tab.Name, status, url, current)
	}

	return nil
}

func printTabsJSON(list *browser.TabList) error {
	result := map[string]any{
		"selectedTab": list.SelectedTab,
		"count":       len(list.Tabs),
	}

	tabs := make([]map[string]any, 0, len(list.Tabs))
	for _, tab := range list.Tabs {
		tabs = append(tabs, map[string]any{
			"name":   tab.Name,
			"status": tab.Status,
			"url":    tab.URL,
			"target": tab.TargetID,
		})
	}
	result["tabs"] = tabs

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	return nil
}

// browserSwitchCmd is a user-friendly alias for "browser tab select"
func browserSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "Switch to a different tab",
		ArgsUsage: "<tab-id>",
		Description: `Set the current working tab.

Examples:
  tap browser switch tab-2
  tap browser switch tab-1`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserSwitch(ctx, cmd)
		},
	}
}

func runBrowserSwitch(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("tab ID required (e.g., tab-1, tab-2)")
	}
	tabName := cmd.Args().First()

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	sessionName := cmd.String("session")
	if err := mgr.SelectTab(ctx, sessionName, tabName); err != nil {
		return fmt.Errorf("switch tab: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Switched to %s\n", green(c, "✓"), tabName)
	return nil
}

// browserCloseTabCmd is a user-friendly alias for "browser tab close"
func browserCloseTabCmd() *cli.Command {
	return &cli.Command{
		Name:      "close-tab",
		Usage:     "Close a browser tab",
		ArgsUsage: "[tab-id]",
		Description: `Close a tab. If no tab ID is given, closes the current tab.

After closing, the next live tab becomes current.

Examples:
  tap browser close-tab        # Close current tab
  tap browser close-tab tab-2  # Close specific tab`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserCloseTab(ctx, cmd)
		},
	}
}

func runBrowserCloseTab(ctx context.Context, cmd *cli.Command) error {
	tabName := cmd.Args().First() // Empty means current/resolved tab

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	sessionName := cmd.String("session")

	// If no tab specified, resolve current
	if tabName == "" {
		session, err := mgr.GetSession(ctx, sessionName)
		if err != nil {
			return err
		}
		tabName = session.SelectedTab
		if tabName == "" {
			return fmt.Errorf("no current tab to close")
		}
	}

	if err := mgr.CloseTab(ctx, sessionName, tabName); err != nil {
		return fmt.Errorf("close tab: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Closed %s\n", green(c, "✓"), tabName)
	return nil
}

// browserStatusCmd is a user-friendly alias for "browser session info"
func browserStatusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show browser context and tab status",
		Description: `Display information about the current browser context
and selected tab.

Examples:
  tap browser status
  tap browser status --json`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserStatus(ctx, cmd)
		},
	}
}

func runBrowserStatus(ctx context.Context, cmd *cli.Command) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	sessionName := cmd.String("session")
	session, err := mgr.GetSession(ctx, sessionName)
	if err != nil {
		if cmd.Bool("json") {
			fmt.Println(`{"error": "no_session"}`)
			return nil
		}
		fmt.Println("No browser session active.")
		fmt.Println("Run: tap browser open <url>")
		return nil
	}

	if cmd.Bool("json") {
		return printStatusJSON(session, nil, "")
	}

	return printStatusHuman(session, nil)
}
