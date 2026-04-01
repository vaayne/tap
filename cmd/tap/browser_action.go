package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func browserNavigateCmd() *cli.Command {
	return &cli.Command{
		Name:      "navigate",
		Usage:     "Navigate a tracked browser tab",
		ArgsUsage: "<url>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("URL required")
			}
			return ensureBrowserPhase3("tap browser navigate", cmd)
		},
	}
}

func browserEvaluateCmd() *cli.Command {
	return &cli.Command{
		Name:      "evaluate",
		Usage:     "Run JavaScript in a tracked browser tab",
		ArgsUsage: "<javascript>",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Description: `Evaluate JavaScript in the resolved tracked tab.

Output formatting follows the same pretty/json/raw conventions used by 'tap site'.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("JavaScript expression required")
			}
			return ensureBrowserPhase3("tap browser evaluate", cmd)
		},
	}
}

func browserScreenshotCmd() *cli.Command {
	return &cli.Command{
		Name:      "screenshot",
		Usage:     "Capture a screenshot from a tracked browser tab",
		ArgsUsage: "",
		Flags:     browserActionFlags(true),
		Description: `Capture a screenshot from the resolved tracked tab.

When --output is omitted, tap will generate a deterministic file path from the
session name, tab name, and current timestamp.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return ensureBrowserPhase3("tap browser screenshot", cmd)
		},
	}
}
