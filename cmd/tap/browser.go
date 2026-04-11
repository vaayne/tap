package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:  "browser",
		Usage: "Browser automation and page interaction",
		Description: `Open pages, interact with elements, and extract content from websites.

Quick start:
  tap browser open <url>                 Open or navigate to a URL
  tap browser tabs                       List open tabs
  tap browser switch <tab-id>            Switch to a different tab
  tap browser click "#submit"            Click an element
  tap browser text                       Extract readable page content
  tap browser screenshot                 Capture the current page

For authenticated access:
  tap attach chrome                      Attach your existing Chrome
  tap browser open https://example.com/login --show

Default behavior:
- Uses the default browser context (managed or attached)
- Operates on the current tab (created automatically if needed)
- Use --new-tab to open in a new tab instead of navigating current`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "state-root",
				Usage:   "Durable state directory for browser sessions and tabs",
				Sources: cli.EnvVars(browser.EnvStateRoot),
			},
		},
		Commands: []*cli.Command{
			withCategory("Common", browserOpenCmd()),
			withCategory("Common", browserTabsCmd()),
			withCategory("Common", browserSwitchCmd()),
			withCategory("Common", browserCloseTabCmd()),
			withCategory("Common", browserStatusCmd()),
			withCategory("Navigation", browserNavigateCmd()),
			withCategory("Navigation", browserBackCmd()),
			withCategory("Navigation", browserForwardCmd()),
			withCategory("Navigation", browserReloadCmd()),
			withCategory("Page Content", browserTextCmd()),
			withCategory("Page Content", browserSnapshotCmd()),
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
			Name:   "session",
			Usage:  "Advanced session override",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "tab",
			Usage:  "Advanced tab override",
			Hidden: true,
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
