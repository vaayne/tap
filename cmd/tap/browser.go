package main

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:               "browser",
		Usage:              "Browser automation and page interaction",
		CustomHelpTemplate: browserHelpTemplate,
		Description: `Open pages, interact with elements, and extract content from websites.

Quick start:
  tap browser open <url>                 Open or navigate to a URL
  tap browser tabs                       List open tabs
  tap browser switch <tab-id>            Switch to a different tab
  tap browser click "#submit"            Click an element
  tap browser text                       Extract readable page content
  tap browser screenshot                 Capture the current page
  tap browser snapshot --interactive     Map interactive elements to @eN refs
  tap browser click @e3                  Act on a ref from the snapshot

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
			withCategory("Tabs", browserOpenCmd()),
			withCategory("Tabs", browserTabsCmd()),
			withCategory("Tabs", browserSwitchCmd()),
			withCategory("Tabs", browserCloseTabCmd()),
			withCategory("Tabs", browserStatusCmd()),
			withCategory("Page actions", browserNavigateCmd()),
			withCategory("Page actions", browserBackCmd()),
			withCategory("Page actions", browserForwardCmd()),
			withCategory("Page actions", browserReloadCmd()),
			withCategory("Queries", browserTextCmd()),
			withCategory("Queries", browserSnapshotCmd()),
			withCategory("Queries", browserEvaluateCmd()),
			withCategory("Queries", browserScreenshotCmd()),
			withCategory("Queries", browserPDFCmd()),
			withCategory("Queries", browserGetCmd()),
			withCategory("Queries", browserIsCmd()),
			withCategory("Page actions", browserClickCmd()),
			withCategory("Page actions", browserDblclickCmd()),
			withCategory("Page actions", browserTypeCmd()),
			withCategory("Page actions", browserFillCmd()),
			withCategory("Page actions", browserHoverCmd()),
			withCategory("Page actions", browserFocusCmd()),
			withCategory("Page actions", browserCheckCmd()),
			withCategory("Page actions", browserUncheckCmd()),
			withCategory("Page actions", browserScrollCmd()),
			withCategory("Page actions", browserScrollIntoViewCmd()),
			withCategory("Page actions", browserSelectCmd()),
			withCategory("Find", browserFindCmd()),
			withCategory("Page actions", browserWaitCmd()),
			withCategory("Page actions", browserKeypressCmd()),
			withCategory("Page actions", browserKeydownCmd()),
			withCategory("Page actions", browserKeyupCmd()),
			withCategory("Page actions", browserKeyboardCmd()),
			withCategory("Page actions", browserMouseCmd()),
			withCategory("Page actions", browserDragCmd()),
			withCategory("Page actions", browserUploadCmd()),
			withCategory("Page actions", browserDialogCmd()),
			withCategory("Storage & state", browserFormsCmd()),
			withCategory("Storage & state", browserCookiesCmd()),
			withCategory("Storage & state", browserStorageCmd()),
			withCategory("Storage & state", browserStateCmd()),
			withCategory("Emulation", browserSetCmd()),
			withCategory("Network", browserNetworkCmd()),
		},
	}
}

const browserHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}

   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`

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
