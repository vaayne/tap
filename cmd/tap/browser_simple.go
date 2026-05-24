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

	ab, err := newAgentBrowser(cmd)
	if err != nil {
		return err
	}

	if cmd.Bool("new-tab") {
		_, _, err := ab.Exec(ctx, "tab", "new", url)
		if err != nil {
			return fmt.Errorf("open new tab: %w", err)
		}
	} else {
		if err := ab.Open(ctx, url, browser.OpenOpts{Headed: cmd.Bool("show")}); err != nil {
			return fmt.Errorf("open: %w", err)
		}
	}

	if d := cmd.Duration("wait"); d > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(d):
		}
	}
	if sel := cmd.String("wait-selector"); sel != "" {
		_, _, err := ab.Exec(ctx, "wait", sel)
		if err != nil {
			return fmt.Errorf("wait selector: %w", err)
		}
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Opened %s\n", green(c, "✓"), url)
	return nil
}

func browserTabsCmd() *cli.Command {
	return &cli.Command{
		Name:  "tabs",
		Usage: "List browser tabs",
		Description: `Show all tabs in the current browser context.

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
	ab, err := newAgentBrowser(cmd)
	if err != nil {
		return err
	}
	out, _, err := ab.Exec(ctx, "tab", "--json")
	if err != nil {
		return err
	}

	var envelope browser.AgentBrowserEnvelope[map[string]any]
	if err := json.Unmarshal(out, &envelope); err != nil {
		return fmt.Errorf("parse tabs: %w", err)
	}
	if !envelope.Success {
		return fmt.Errorf("tabs: %s", envelope.Error)
	}

	if cmd.Bool("json") {
		pretty, _ := json.MarshalIndent(envelope.Data, "", "  ")
		fmt.Println(string(pretty))
		return nil
	}

	data := envelope.Data
	tabs, _ := data["tabs"].([]any)
	if len(tabs) == 0 {
		fmt.Println("No tabs found.")
		return nil
	}

	c := true
	fmt.Printf("%s %d tabs\n\n", bold(c, "Tabs:"), len(tabs))
	for _, t := range tabs {
		if m, ok := t.(map[string]any); ok {
			id, _ := m["id"].(string)
			url, _ := m["url"].(string)
			fmt.Printf("  %-6s %s\n", id, url)
		}
	}
	return nil
}

func browserSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "Switch to a different tab",
		ArgsUsage: "<tab-id>",
		Description: `Set the current working tab.

Examples:
  tap browser switch t2
  tap browser switch t1`,
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
		return fmt.Errorf("tab ID required (e.g., t1, t2)")
	}
	tabID := cmd.Args().First()

	ab, err := newAgentBrowser(cmd)
	if err != nil {
		return err
	}

	_, _, err = ab.Exec(ctx, "tab", tabID)
	if err != nil {
		return fmt.Errorf("switch tab: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Switched to %s\n", green(c, "✓"), tabID)
	return nil
}

func browserCloseTabCmd() *cli.Command {
	return &cli.Command{
		Name:      "close-tab",
		Usage:     "Close a browser tab",
		ArgsUsage: "<tab-id>",
		Description: `Close a specific tab.

Examples:
  tap browser close-tab t2`,
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
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("tab ID required")
	}
	tabID := cmd.Args().First()

	ab, err := newAgentBrowser(cmd)
	if err != nil {
		return err
	}

	_, _, err = ab.Exec(ctx, "tab", "close", tabID)
	if err != nil {
		return fmt.Errorf("close tab: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Closed %s\n", green(c, "✓"), tabID)
	return nil
}

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
	ab, err := newAgentBrowser(cmd)
	if err != nil {
		return err
	}

	sessionOut, _, _ := ab.Exec(ctx, "session", "--json")
	tabOut, _, _ := ab.Exec(ctx, "tab", "--json")
	urlOut, _, _ := ab.Exec(ctx, "get", "url", "--json")

	var sessionEnv browser.AgentBrowserEnvelope[map[string]any]
	_ = json.Unmarshal(sessionOut, &sessionEnv)

	var tabEnv browser.AgentBrowserEnvelope[map[string]any]
	_ = json.Unmarshal(tabOut, &tabEnv)

	var urlEnv browser.AgentBrowserEnvelope[string]
	_ = json.Unmarshal(urlOut, &urlEnv)

	result := map[string]any{
		"session": sessionEnv.Data,
		"tabs":    tabEnv.Data,
		"url":     urlEnv.Data,
	}
	if isAttachedMode() {
		result["attached"] = true
	}

	if cmd.Bool("json") {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return nil
	}

	c := true
	if name, ok := sessionEnv.Data["name"].(string); ok && name != "" {
		fmt.Printf("%s %s\n", bold(c, "Session:"), name)
	} else if sessionEnv.Success {
		fmt.Printf("%s %v\n", bold(c, "Session:"), sessionEnv.Data)
	}
	if urlEnv.Success && urlEnv.Data != "" {
		fmt.Printf("%s %s\n", bold(c, "URL:"), urlEnv.Data)
	}
	if tabs, ok := tabEnv.Data["tabs"].([]any); ok {
		fmt.Printf("%s %d\n", bold(c, "Tabs:"), len(tabs))
	}
	if isAttachedMode() {
		fmt.Printf("%s attached mode\n", bold(c, "Mode:"))
	}

	return nil
}
