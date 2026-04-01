package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func browserTabCmd() *cli.Command {
	return &cli.Command{
		Name:  "tab",
		Usage: "Manage tracked tabs within a persistent browser session",
		Description: `Tap manages only tracked tabs recorded in session metadata.

Closing or losing a tracked target updates its stored status. If the selected tab
is closed, tap selects the next remaining live tracked tab by creation order. If
no live tracked tabs remain, the session has no selected tab.`,
		Commands: []*cli.Command{
			browserTabNewCmd(),
			browserTabListCmd(),
			browserTabSelectCmd(),
			browserTabCloseCmd(),
		},
	}
}

func browserTabNewCmd() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a tracked browser tab",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "Initial URL for the new tab",
				Value: "about:blank",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("tab name required")
			}
			return ensureBrowserPhase3("tap browser tab new", cmd)
		},
	}
}

func browserTabListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List tracked tabs for a session",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser tab list", cmd)
		},
	}
}

func browserTabSelectCmd() *cli.Command {
	return &cli.Command{
		Name:      "select",
		Usage:     "Set the default live tracked tab for a session",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("tab name required")
			}
			return ensureBrowserPhase3("tap browser tab select", cmd)
		},
	}
}

func browserTabCloseCmd() *cli.Command {
	return &cli.Command{
		Name:      "close",
		Usage:     "Close and remove a tracked browser tab",
		ArgsUsage: "[name]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Description: `Close the resolved tracked tab and remove its metadata.

If the closed tab was selected, tap promotes the next remaining live tracked tab
by creation order. If none remain, selected-tab state is cleared.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser tab close", cmd)
		},
	}
}
