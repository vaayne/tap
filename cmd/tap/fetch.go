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
		Usage:     "Fetch and extract clean content from a URL",
		ArgsUsage: "<url>",
		Description: `Fetch a URL and extract the main content as clean Markdown, stripping ads,
navigation, scripts, and boilerplate via go-defuddle.

Examples:
  tap fetch https://example.com/article          Clean Markdown output
  tap fetch --json https://example.com/article   Full metadata as JSON
  tap fetch -b https://example.com               Use browser for JS-rendered pages
  tap fetch --lp https://example.com             Use Lightpanda (fast headless, no auth)`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON with full metadata",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("URL required")
			}

			url := cmd.Args().First()

			client, err := newClient(ctx, cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			opts := &fetch.Options{
				Markdown:   true,
				UseBrowser: cmd.Bool("browser"),
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
