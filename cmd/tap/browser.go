package main

import (
	"fmt"
	"os"
	"path/filepath"

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
				Sources: cli.EnvVars("TAP_BROWSER_STATE_ROOT"),
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
			withCategory("State", browserFormsCmd()),
			withCategory("State", browserCookiesCmd()),
			withCategory("Network", browserNetworkCmd()),
			withCategory("Emulation", browserSetCmd()),
			withCategory("Storage", browserStorageCmd()),
			withCategory("State", browserStateCmd()),
			withCategory("Auth", browserAuthCmd()),
			withCategory("Info", browserGetCmd()),
			withCategory("Performance", browserVitalsCmd()),
			withCategory("Diff", browserDiffCmd()),
		},
	}
}

func newAgentBrowser(cmd *cli.Command) (*browser.AgentBrowser, error) {
	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return nil, fmt.Errorf("agent-browser: %w", err)
	}
	if name := cmd.String("session"); name != "" {
		ab.SessionName = name
	}
	if isAttachedMode() {
		ab.Attached = true
		ab.SessionName = ""
	}
	return ab, nil
}

func browserActionFlags(includeOutput bool) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:   "session",
			Usage:  "Session name (maps to --session-name)",
			Hidden: true,
		},
	}
	if includeOutput {
		flags = append(flags, &cli.StringFlag{
			Name:  "output",
			Usage: "Write output to this file; defaults to a generated path when omitted",
		})
	}
	return flags
}

func withCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = cat
	return cmd
}

func attachedFlagPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "tap", "browser", "attached")
}

func isAttachedMode() bool {
	_, err := os.Stat(attachedFlagPath())
	return err == nil
}

func setAttachedMode() error {
	p := attachedFlagPath()
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	return os.WriteFile(p, []byte{}, 0o644)
}

func clearAttachedMode() error {
	_ = os.Remove(attachedFlagPath())
	return nil
}
