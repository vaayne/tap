package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/fetch"
)

func fetchCmd() *cli.Command {
	return &cli.Command{
		Name:      "fetch",
		Usage:     "Extract clean content from a URL or the current agent-browser tab",
		ArgsUsage: "[url]",
		Description: `Navigate to an optional URL and extract the main content with Defuddle.
With no URL, extract the active tab without navigating or creating a session.

Examples:
  tap fetch https://example.com/article          Clean Markdown output
  tap fetch --json https://example.com/article   Full metadata as JSON
	 tap fetch                                      Extract current tab`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON with full metadata",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			url := cmd.Args().First()

			client, err := newRuntimeClient(ctx, cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			opts := &fetch.Options{
				Markdown: true,
			}
			result, err := client.Fetch(ctx, url, opts)
			if err != nil {
				return err
			}

			if cmd.Bool("json") {
				out, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal result: %w", err)
				}
				fmt.Println(string(out))
			} else {
				if result.Title != "" {
					fmt.Printf("# %s\n\n", result.Title)
				}
				if result.Markdown != "" {
					fmt.Println(result.Markdown)
				} else if result.Content != "" {
					fmt.Println(result.Content)
				}
			}
			return nil
		},
	}
}
