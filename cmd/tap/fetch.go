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

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

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
