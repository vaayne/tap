package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:  "browser",
		Usage: "Manage persistent browser sessions and tracked tabs",
		Description: `Persistent browser automation lives under the "browser" namespace.

Session resolution order:
  1. --session
  2. the selected session from 'tap browser session select'
  3. the only available session, when exactly one exists

Tab resolution order:
  1. --tab
  2. the selected tab within the resolved session
  3. the only live tracked tab, when exactly one exists

Tracked tabs are named browser targets stored in tap metadata. Untracked live
browser tabs are ignored by default. When a tracked target disappears, tap marks
it stale and clears selected-tab state instead of silently adopting a new tab.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "state-root",
				Usage:   "Durable state directory for browser sessions and tabs",
				Sources: cli.EnvVars(browser.EnvStateRoot),
			},
		},
		Commands: []*cli.Command{
			browserSessionCmd(),
			browserTabCmd(),
			browserNavigateCmd(),
			browserEvaluateCmd(),
			browserScreenshotCmd(),
		},
	}
}

func browserActionNotImplemented(name string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, _ *cli.Command) error {
		return fmt.Errorf("%s is defined in Phase 1 but will be implemented in a later phase", name)
	}
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

func ensureBrowserPhase3(name string, cmd *cli.Command) error {
	configureLogging(cmd)
	root := browserStateRoot(cmd)
	if root == "" {
		return errors.New("browser state root could not be resolved")
	}
	return browserActionNotImplemented(name)(context.Background(), cmd)
}
