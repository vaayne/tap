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
	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return err
	}
	if isAttachedMode() {
		ab.Attached = true
		ab.SessionName = ""
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

	return printStatusHuman(result)
}

func printStatusHuman(result map[string]any) error {
	c := true

	if attached, ok := result["attached"].(bool); ok && attached {
		fmt.Printf("%s Attached mode\n", bold(c, "Context:"))
	}

	if session, ok := result["session"].(map[string]any); ok {
		if name, ok := session["name"].(string); ok && name != "" {
			fmt.Printf("%s %s\n", bold(c, "Session:"), name)
		}
	}
	if url, ok := result["url"].(string); ok && url != "" {
		fmt.Printf("%s %s\n", bold(c, "URL:"), url)
	}

	if tabs, ok := result["tabs"].(map[string]any); ok {
		if list, ok := tabs["tabs"].([]any); ok && len(list) > 0 {
			fmt.Printf("%s %d tabs\n", bold(c, "Tabs:"), len(list))
		}
	} else {
		fmt.Println("No current tab selected.")
		fmt.Println("Run: tap browser open <url>")
	}

	return nil
}
