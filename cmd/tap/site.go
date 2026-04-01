package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
)

func siteCmd() *cli.Command {
	return &cli.Command{
		Name:      "site",
		Usage:     "Run site scripts",
		ArgsUsage: "<script-name> [key=value ...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		},
		Commands: []*cli.Command{
			siteListCmd(),
			siteInfoCmd(),
			siteSearchCmd(),
			siteSyncCmd(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args()
			if args.Len() == 0 {
				return fmt.Errorf("script name required. Run 'tap site list' to see available scripts")
			}

			scriptName := args.First()
			scriptArgs := parseArgs(args.Tail())

			client, err := newClient(ctx, cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			if cmd.Bool("verbose") {
				mode := "auto (QuickJS → Browser)"
				if cmd.Bool("browser") {
					mode = "browser"
				}
				log.Printf("Running: %s [engine=%s]", scriptName, mode)
			}

			result, err := client.RunScript(ctx, scriptName, scriptArgs)
			if err != nil {
				return err
			}

			return printResult(cmd, result)
		},
	}
}

func siteListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all available scripts (grouped by site)",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			client, err := newClient(ctx, cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			scripts := client.ListScripts()
			if len(scripts) == 0 {
				fmt.Println("No scripts found.")
				return nil
			}

			color := useColor(cmd)
			groups := groupScripts(scripts)

			// Sort group names
			groupNames := make([]string, 0, len(groups))
			for name := range groups {
				groupNames = append(groupNames, name)
			}
			sort.Strings(groupNames)

			for i, groupName := range groupNames {
				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("%s\n", bold(color, groupName+"/"))
				for _, s := range groups[groupName] {
					argHints := formatArgHints(s, color)
					actionName := s.Meta.Name
					// Strip the group prefix for display
					if _, after, ok := strings.Cut(actionName, "/"); ok {
						actionName = after
					}
					fmt.Printf("  %-24s %s%s\n",
						green(color, actionName),
						s.Meta.Description,
						argHints,
					)
				}
			}

			fmt.Printf("\n%s\n", dim(color, fmt.Sprintf("%d scripts across %d sites", len(scripts), len(groups))))
			return nil
		},
	}
}

func siteInfoCmd() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Show detailed info for a script",
		ArgsUsage: "<script-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("script name required")
			}

			client, err := newClient(ctx, cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			name := cmd.Args().First()
			s, ok := client.GetScript(name)
			if !ok {
				return &tap.ScriptNotFoundError{Name: name, Available: scriptNames(client)}
			}

			color := useColor(cmd)

			fmt.Printf("%s\n", bold(color, s.Meta.Name))
			fmt.Printf("  %s\n\n", s.Meta.Description)

			fmt.Printf("  %s  %s\n", bold(color, "Domain:"), s.Meta.Domain)

			if s.Meta.Example != "" {
				fmt.Printf("  %s %s\n", bold(color, "Example:"), s.Meta.Example)
			}

			if len(s.Meta.Args) > 0 {
				fmt.Printf("\n  %s\n", bold(color, "Arguments:"))
				// Sort args for consistent output
				argNames := make([]string, 0, len(s.Meta.Args))
				for name := range s.Meta.Args {
					argNames = append(argNames, name)
				}
				sort.Strings(argNames)
				for _, argName := range argNames {
					def := s.Meta.Args[argName]
					req := dim(color, "optional")
					if def.Required {
						req = yellow(color, "required")
					}
					fmt.Printf("    %-16s %s  %s\n",
						green(color, argName),
						dim(color, "("+req+")"),
						def.Description,
					)
				}
			}

			fmt.Printf("\n  %s tap site %s", bold(color, "Usage:"), s.Meta.Name)
			for argName, def := range s.Meta.Args {
				if def.Required {
					fmt.Printf(" %s=<%s>", argName, argName)
				} else {
					fmt.Printf(" [%s=value]", argName)
				}
			}
			fmt.Println()

			return nil
		},
	}
}

func siteSearchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search scripts online by name or description",
		ArgsUsage: "<query>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("search query required")
			}

			query := strings.Join(cmd.Args().Slice(), " ")
			color := useColor(cmd)

			result, err := searchOnline(query)
			if err != nil {
				return err
			}

			if len(result.Scripts) == 0 {
				fmt.Printf("No scripts matching %q\n", query)
				return nil
			}

			for _, s := range result.Scripts {
				argHints := formatOnlineArgHints(s, color)
				fmt.Printf("  %-30s %s%s\n",
					green(color, s.Name),
					s.Description,
					argHints,
				)
			}
			fmt.Printf("\n%s\n", dim(color, fmt.Sprintf("%d result(s)", len(result.Scripts))))
			return nil
		},
	}
}

func siteSyncCmd() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync scripts from the remote catalog",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			dir := cmd.String("sites-dir")
			if dir == "" {
				dir = defaultSitesDir()
			}
			return syncScripts(dir, cmd.Bool("verbose"))
		},
	}
}
