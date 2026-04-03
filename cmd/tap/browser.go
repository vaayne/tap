package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:  "browser",
		Usage: "Manage persistent browser sessions and tracked tabs",
		Description: `Persistent browser automation with named sessions and tracked tabs.

Quick start:
  tap browser session new default          Start a headless browser
  tap browser tab new main --url <url>     Open a tracked tab
  tap browser text                         Extract clean content (Markdown)
  tap browser click "button.submit"        Interact with elements
  tap browser screenshot                   Capture the page

Commands are grouped by function:
  Session & Tab:  session, tab
  Navigation:     navigate, back, forward, reload
  Page Content:   text, evaluate, screenshot, pdf
  Interaction:    click, type, fill, hover, scroll, select, wait, keypress, dialog
  State:          forms, cookies
  Network:        network (wait, log, body, intercept, clear)

Session resolution: --session flag → selected session → the only session.
Tab resolution: --tab flag → selected tab → the only live tracked tab.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "state-root",
				Usage:   "Durable state directory for browser sessions and tabs",
				Sources: cli.EnvVars(browser.EnvStateRoot),
			},
		},
		Commands: []*cli.Command{
			withCategory("Session & Tab", browserSessionCmd()),
			withCategory("Session & Tab", browserTabCmd()),
			withCategory("Navigation", browserNavigateCmd()),
			withCategory("Navigation", browserBackCmd()),
			withCategory("Navigation", browserForwardCmd()),
			withCategory("Navigation", browserReloadCmd()),
			withCategory("Page Content", browserTextCmd()),
			withCategory("Page Content", browserEvaluateCmd()),
			withCategory("Page Content", browserScreenshotCmd()),
			withCategory("Page Content", browserPDFCmd()),
			withCategory("Interaction", browserClickCmd()),
			withCategory("Interaction", browserTypeCmd()),
			withCategory("Interaction", browserFillCmd()),
			withCategory("Interaction", browserHoverCmd()),
			withCategory("Interaction", browserScrollCmd()),
			withCategory("Interaction", browserSelectCmd()),
			withCategory("Interaction", browserWaitCmd()),
			withCategory("Interaction", browserKeypressCmd()),
			withCategory("Interaction", browserDialogCmd()),
			withCategory("State", browserFormsCmd()),
			withCategory("State", browserCookiesCmd()),
			withCategory("Network", browserNetworkCmd()),
		},
	}
}

func newBrowserManager(cmd *cli.Command) (*browser.Manager, error) {
	root := browserStateRoot(cmd)
	store, err := browser.NewStore(root)
	if err != nil {
		return nil, fmt.Errorf("init browser store: %w", err)
	}
	return browser.NewManager(store), nil
}

func browserActionFlags(includeOutput bool) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "session",
			Usage: "Browser session name; defaults to the selected session when omitted",
		},
		&cli.StringFlag{
			Name:  "tab",
			Usage: "Tracked tab name; defaults to the selected tab in the resolved session when omitted",
		},
	}
	if includeOutput {
		flags = append(flags, &cli.StringFlag{
			Name:  "output",
			Usage: "Write the screenshot to this file; defaults to a generated path when omitted",
		})
	}
	return flags
}

func withCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = cat
	return cmd
}

func browserStateRoot(cmd *cli.Command) string {
	root := cmd.String("state-root")
	if root != "" {
		return root
	}
	defaultRoot, err := browser.DefaultStateRoot()
	if err != nil {
		return ""
	}
	return defaultRoot
}
