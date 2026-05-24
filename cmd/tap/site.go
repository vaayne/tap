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
		Name:          "site",
		Usage:         "Run site scripts",
		ArgsUsage:     "<script-name> [key=value ...]",
		ShellComplete: completeSiteRoot,
		Description: `Run site-specific JavaScript scripts that extract structured data from websites.

Scripts are organized as site/action (e.g. hackernews/top, google/search) and
execute in QuickJS with automatic browser fallback when cookies or DOM are needed.

Scripts auto-sync from the remote catalog every 24 hours into ~/.cache/tap/sites/.
Local overrides in ~/.config/tap/sites/ take precedence over cached scripts.

Examples:
  tap site list                              List all available scripts
  tap site hackernews/search query=golang    Run a script with arguments
  tap site -b twitter/search query=go        Run with browser cookies (auth)
  tap site info hackernews/search            Show script details and args
  tap site search "weather"                  Search the online catalog
  tap site sync                              Force-refresh the script cache
  tap site hackernews/top -f json            Output as JSON`,
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		}, browserClientFlags()...),
		Commands: []*cli.Command{
			siteRunCmd(),
			siteListCmd(),
			siteInfoCmd(),
			siteSearchCmd(),
			siteSyncCmd(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runSiteScript(ctx, cmd, cmd.Args().Slice())
		},
	}
}

func siteRunCmd() *cli.Command {
	return &cli.Command{
		Name:          "run",
		Usage:         "Run a site script",
		ArgsUsage:     "<script-name> [key=value ...]",
		ShellComplete: completeSiteScripts,
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		}, browserClientFlags()...),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runSiteScript(ctx, cmd, cmd.Args().Slice())
		},
	}
}

func runSiteScript(ctx context.Context, cmd *cli.Command, rawArgs []string) error {
	if len(rawArgs) == 0 {
		return fmt.Errorf("script name required. Run 'tap site list' to see available scripts")
	}

	scriptName := rawArgs[0]
	scriptArgs := parseArgs(rawArgs[1:])

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
					runtimeBadge := ""
					if s.Meta.Runtime != "" && s.Meta.Runtime != "auto" {
						runtimeBadge = dim(color, " ["+s.Meta.Runtime+"]")
					}
					fmt.Printf("  %-24s %s%s%s\n",
						green(color, actionName),
						s.Meta.Description,
						runtimeBadge,
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
		Name:          "info",
		Usage:         "Show detailed info for a script",
		ArgsUsage:     "<script-name>",
		ShellComplete: completeSiteScripts,
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

			if s.Meta.Runtime != "" && s.Meta.Runtime != "auto" {
				fmt.Printf("  %s  %s\n", bold(color, "Runtime:"), s.Meta.Runtime)
			}

			if s.Meta.Example != "" {
				fmt.Printf("  %s %s\n", bold(color, "Example:"), s.Meta.Example)
			}

			if len(s.Meta.Env) > 0 {
				fmt.Printf("\n  %s\n", bold(color, "Env:"))
				envNames := make([]string, 0, len(s.Meta.Env))
				for name := range s.Meta.Env {
					envNames = append(envNames, name)
				}
				sort.Strings(envNames)
				for _, envName := range envNames {
					def := s.Meta.Env[envName]
					req := dim(color, "optional")
					if def.Required {
						req = yellow(color, "required")
					}
					fmt.Printf("    %-16s %s  %s\n",
						green(color, envName),
						dim(color, "("+req+")"),
						def.Description,
					)
				}
			}

			if len(s.Meta.Headers) > 0 {
				fmt.Printf("\n  %s\n", bold(color, "Headers:"))
				headerKeys := make([]string, 0, len(s.Meta.Headers))
				for k := range s.Meta.Headers {
					headerKeys = append(headerKeys, k)
				}
				sort.Strings(headerKeys)
				for _, k := range headerKeys {
					v := s.Meta.Headers[k]
					if strings.Contains(v, "${") {
						v = dim(color, "(from env)")
					}
					fmt.Printf("    %-16s %s\n",
						green(color, k),
						v,
					)
				}
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
