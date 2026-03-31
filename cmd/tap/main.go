package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
	"github.com/vaayne/tap/fetch"
	"github.com/vaayne/tap/script"
)

var version = "dev"

func main() {
	_ = godotenv.Load()

	app := &cli.Command{
		Name:    "tap",
		Usage:   "Tap into any website from your terminal",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "sites-dir",
				Usage:   "Directory containing site scripts",
				Value:   "./sites",
				Sources: cli.EnvVars("TAP_SITES_DIR"),
			},
			&cli.StringFlag{
				Name:    "ws-url",
				Usage:   "Remote CDP WebSocket URL",
				Sources: cli.EnvVars("TAP_WS_URL"),
			},
			&cli.StringFlag{
				Name:    "profile-dir",
				Usage:   "Chrome profile directory for persistent cookies",
				Sources: cli.EnvVars("TAP_PROFILE_DIR"),
			},
		},
		Commands: []*cli.Command{
			siteCmd(),
			fetchCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func newClient(cmd *cli.Command) (*tap.Client, error) {
	var opts []tap.Option

	if dir := cmd.String("sites-dir"); dir != "" {
		opts = append(opts, tap.WithSitesDir(dir))
	}
	if url := cmd.String("ws-url"); url != "" {
		opts = append(opts, tap.WithWSURL(url))
	}
	if dir := cmd.String("profile-dir"); dir != "" {
		opts = append(opts, tap.WithProfileDir(dir))
	}

	return tap.New(opts...)
}

func siteCmd() *cli.Command {
	return &cli.Command{
		Name:      "site",
		Usage:     "Run site scripts",
		ArgsUsage: "<script-name> [key=value ...]",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all available scripts",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					client, err := newClient(cmd)
					if err != nil {
						return err
					}
					defer client.Close()

					for _, s := range client.ListScripts() {
						argHints := formatArgHints(s)
						fmt.Printf("  %-30s %s%s\n", s.Meta.Name, s.Meta.Description, argHints)
					}
					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args()
			if args.Len() == 0 {
				return fmt.Errorf("script name required. Run 'tap site list' to see available scripts")
			}

			scriptName := args.First()
			scriptArgs := parseArgs(args.Tail())

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			log.Printf("Running: %s", scriptName)

			result, err := client.RunScript(ctx, scriptName, scriptArgs)
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal result: %w", err)
			}
			fmt.Println(string(out))
			return nil
		},
	}
}

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
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("URL required")
			}

			url := cmd.Args().First()

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			opts := &fetch.Options{Markdown: true}
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
				// Default: markdown output
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

func parseArgs(raw []string) map[string]string {
	args := make(map[string]string)
	for _, s := range raw {
		if k, v, ok := strings.Cut(s, "="); ok {
			args[k] = v
		}
	}
	return args
}

func formatArgHints(s *script.Script) string {
	if len(s.Meta.Args) == 0 {
		return ""
	}
	var parts []string
	for name, def := range s.Meta.Args {
		tag := name
		if def.Required {
			tag = name + "*"
		}
		parts = append(parts, tag)
	}
	return " [" + strings.Join(parts, ", ") + "]"
}
