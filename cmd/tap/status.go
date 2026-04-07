package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func statusCmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show the active browser context and current tab",
		Description: `Display the current default browser context, authentication state,
and active tab information.

This command helps you understand which browser context and tab
will be used by subsequent tap commands.

Examples:
  tap status              Show human-readable status
  tap status --json       Show machine-readable status`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runStatus(ctx, cmd)
		},
	}
}

func runStatus(ctx context.Context, cmd *cli.Command) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}

	// Get default session info
	session, err := mgr.GetSession(ctx, "")
	if err != nil {
		// No default session - show empty state
		if cmd.Bool("json") {
			return printStatusJSON(nil, nil, "no_default_session")
		}
		return printStatusEmpty()
	}

	// Get current tab
	var currentTab *browser.TabRecord
	if session.SelectedTab != "" {
		if tab, ok := session.Tabs[session.SelectedTab]; ok {
			currentTab = tab
		}
	}

	if cmd.Bool("json") {
		return printStatusJSON(session, currentTab, "")
	}

	return printStatusHuman(session, currentTab)
}

func printStatusEmpty() error {
	c := true // assume color for now
	fmt.Println("No active browser context.")
	fmt.Println()
	fmt.Printf("%s\n", bold(c, "Quick start:"))
	fmt.Println("  tap attach chrome         Attach to your existing Chrome")
	fmt.Println("  tap browser open <url>    Open a page in managed browser")
	return nil
}

func printStatusHuman(session *browser.SessionRecord, currentTab *browser.TabRecord) error {
	c := true

	// Context type
	contextType := "Managed local browser"
	if session.Remote != nil {
		contextType = "Attached remote browser"
	}

	fmt.Printf("%s %s\n", bold(c, "Browser context:"), contextType)
	fmt.Printf("%s %s\n", bold(c, "Session:"), session.Name)

	if session.Process != nil && session.Process.DebugURL != "" {
		fmt.Printf("%s %s\n", bold(c, "Debug URL:"), session.Process.DebugURL)
	} else if session.Remote != nil {
		fmt.Printf("%s %s\n", bold(c, "Remote URL:"), session.Remote.WSURL)
	}

	// Tab info
	fmt.Println()
	if currentTab != nil {
		fmt.Printf("%s %s\n", bold(c, "Current tab:"), currentTab.Name)
		fmt.Printf("%s %s\n", bold(c, "Status:"), string(currentTab.Status))
		if currentTab.URL != "" {
			fmt.Printf("%s %s\n", bold(c, "URL:"), currentTab.URL)
		}
	} else {
		fmt.Println("No current tab selected.")
		fmt.Println("Run: tap browser open <url>")
	}

	// Tab count
	liveCount := 0
	for _, tab := range session.Tabs {
		if tab.Status == browser.TabStatusLive {
			liveCount++
		}
	}
	if len(session.Tabs) > 0 {
		fmt.Println()
		fmt.Printf("%s %d total (%d live)\n", bold(c, "Tabs:"), len(session.Tabs), liveCount)
	}

	return nil
}

func printStatusJSON(session *browser.SessionRecord, currentTab *browser.TabRecord, errorState string) error {
	result := map[string]any{
		"error": errorState,
	}

	if session != nil {
		result["session"] = map[string]any{
			"name":      session.Name,
			"mode":      session.Mode,
			"createdAt": session.CreatedAt,
		}

		if session.Process != nil {
			result["process"] = map[string]any{
				"pid":      session.Process.PID,
				"debugURL": session.Process.DebugURL,
			}
		}

		if session.Remote != nil {
			result["remote"] = map[string]any{
				"wsURL": session.Remote.WSURL,
			}
		}

		tabs := make([]map[string]any, 0, len(session.Tabs))
		for name, tab := range session.Tabs {
			tabInfo := map[string]any{
				"name":   name,
				"status": tab.Status,
				"url":    tab.URL,
			}
			tabs = append(tabs, tabInfo)
		}
		result["tabs"] = tabs

		if currentTab != nil {
			result["currentTab"] = map[string]any{
				"name":   currentTab.Name,
				"status": currentTab.Status,
				"url":    currentTab.URL,
			}
		}
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal status: %w", err)
	}
	fmt.Println(string(out))
	return nil
}
