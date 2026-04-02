package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

func browserNavigateCmd() *cli.Command {
	return &cli.Command{
		Name:      "navigate",
		Usage:     "Navigate a tracked browser tab",
		ArgsUsage: "<url>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			url := cmd.Args().First()
			if url == "" {
				return fmt.Errorf("URL required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			if err := mgr.Navigate(ctx, sessionName, tabName, url); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Navigated to %s\n", url)
			return nil
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
			configureLogging(cmd)
			js := cmd.Args().First()
			if js == "" {
				return fmt.Errorf("JavaScript expression required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			result, err := mgr.Evaluate(ctx, sessionName, tabName, js)
			if err != nil {
				return err
			}
			return printResult(cmd, result)
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
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			data, err := mgr.Screenshot(ctx, sessionName, tabName)
			if err != nil {
				return err
			}

			outPath := cmd.String("output")
			if outPath == "" {
				s := sessionName
				if s == "" {
					s = "default"
				}
				t := tabName
				if t == "" {
					t = "default"
				}
				outPath = fmt.Sprintf("screenshot-%s-%s-%d.png", s, t, time.Now().Unix())
			}

			if err := os.WriteFile(outPath, data, 0o644); err != nil {
				return fmt.Errorf("write screenshot: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%s\n", outPath)
			return nil
		},
	}
}
