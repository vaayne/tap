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

Your browser must have remote debugging enabled.

Examples:
  tap attach chrome                    Auto-discover from DevToolsActivePort
  tap attach chrome --browser-url http://localhost:9222`,
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runAttachChrome(ctx, cmd)
		},
	}
}

func runAttachChrome(ctx context.Context, cmd *cli.Command) error {
	url := cmd.String("browser-url")
	if url == "" {
		if portFile := cmd.String("port-file"); portFile != "" {
			return fmt.Errorf("port-file auto-discovery not yet implemented with agent-browser")
		}
		return fmt.Errorf("no Chrome URL provided; use --browser-url http://localhost:9222")
	}

	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return err
	}
	ab.SessionName = ""
	_, _, err = ab.Exec(ctx, "connect", url)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	if err := setAttachedMode(); err != nil {
		return fmt.Errorf("set attached mode: %w", err)
	}

	c := true
	fmt.Fprintf(os.Stderr, "%s Attached to Chrome\n", green(c, "✓"))
	fmt.Fprintf(os.Stderr, "  URL: %s\n", url)
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
	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return err
	}
	ab.SessionName = ""

	sessionOut, _, _ := ab.Exec(ctx, "session", "list", "--json")
	urlOut, _, _ := ab.Exec(ctx, "get", "url", "--json")

	var sessionEnv browser.AgentBrowserEnvelope[[]any]
	_ = json.Unmarshal(sessionOut, &sessionEnv)

	var urlEnv browser.AgentBrowserEnvelope[string]
	_ = json.Unmarshal(urlOut, &urlEnv)

	attached := isAttachedMode()

	result := map[string]any{
		"attached": attached,
		"sessions": sessionEnv.Data,
		"url":      urlEnv.Data,
	}

	if cmd.Bool("json") {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return nil
	}

	c := true
	if attached {
		fmt.Printf("%s Attached mode active\n", green(c, "✓"))
	} else {
		fmt.Println("No attachment active.")
		fmt.Println()
		fmt.Println("To attach:")
		fmt.Println("  tap attach chrome --browser-url http://localhost:9222")
		return nil
	}
	if urlEnv.Success && urlEnv.Data != "" {
		fmt.Printf("%s Current URL: %s\n", bold(c, "URL:"), urlEnv.Data)
	}
	if sessions := sessionEnv.Data; len(sessions) > 0 {
		fmt.Printf("%s %d sessions\n", bold(c, "Sessions:"), len(sessions))
	}
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
	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return err
	}
	ab.SessionName = ""
	_, _, _ = ab.Exec(ctx, "close")

	_ = clearAttachedMode()

	c := true
	fmt.Fprintf(os.Stderr, "%s Cleared attached Chrome metadata\n", green(c, "✓"))
	return nil
}
