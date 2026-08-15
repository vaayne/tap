package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
)

var version = "dev"

func main() {
	_ = godotenv.Load()
	configureHelpTemplates()
	if err := newApp().Run(context.Background(), os.Args); err != nil {
		if notFound, ok := err.(*tap.ScriptNotFoundError); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", notFound.Error())
			if suggestions := notFound.Suggestions(5); len(suggestions) > 0 {
				fmt.Fprintln(os.Stderr, "\nDid you mean?")
				for _, suggestion := range suggestions {
					fmt.Fprintf(os.Stderr, "  tap site %s\n", suggestion)
				}
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func configureHelpTemplates() {
	cli.ShowSubcommandHelp = func(cmd *cli.Command) error {
		template := cli.SubcommandHelpTemplate
		if cmd.CustomHelpTemplate != "" {
			template = cmd.CustomHelpTemplate
		}
		cli.HelpPrinter(cmd.Root().Writer, template, cmd)
		return nil
	}
}

func newApp() *cli.Command {
	return &cli.Command{
		Name:                            "tap",
		Usage:                           "Browser workflows, reusable site programs, and web extraction",
		Version:                         version,
		EnableShellCompletion:           true,
		ConfigureShellCompletionCommand: configureCompletionCommand,
		Description: `Tap runs host-side JavaScript workflows, discovers reusable site programs,
and extracts clean content through the active agent-browser session.

Quick start:
  tap site list
  tap site exa/search query="agent-browser" count=5
  tap fetch https://example.com
  tap fetch                              Extract the current agent-browser tab
  tap run workflow.js                    Run a JavaScript browser workflow
  tap doctor                             Check the runtime dependency

Tap inherits AGENT_BROWSER_SESSION and never creates, names, or closes sessions.`,
		Flags: globalFlags(),
		Commands: []*cli.Command{
			siteCmd(),
			fetchCmd(),
			runCmd(),
			doctorCmd(),
			upgradeCmd(),
			skillCmd(),
			docsCmd(),
		},
	}
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "sites-dir",
			Usage:   "Directory containing cached site scripts",
			Value:   defaultSitesDir(),
			Sources: cli.EnvVars("TAP_SITES_DIR"),
		},
		&cli.StringFlag{
			Name:    "agent-browser",
			Usage:   "agent-browser executable override",
			Sources: cli.EnvVars("TAP_AGENT_BROWSER"),
			Hidden:  true,
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout; 0 means agent-browser default",
			Sources: cli.EnvVars("TAP_TIMEOUT"),
			Hidden:  true,
		},
		&cli.BoolFlag{Name: "verbose", Usage: "Enable verbose logging"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "Suppress log output"},
		&cli.BoolFlag{
			Name:  "local-only",
			Usage: "Only use scripts from ~/.config/tap/sites/",
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "Disable colored output",
			Sources: cli.EnvVars("NO_COLOR"),
		},
	}
}

func newClient(ctx context.Context, cmd *cli.Command) (*tap.Client, error) {
	dir := cmd.String("sites-dir")
	if dir == "" {
		dir = defaultSitesDir()
	}
	overrideDir := defaultLocalOverrideDir()
	localOnly := cmd.Bool("local-only")
	if !localOnly {
		if err := ensureScripts(dir, cmd.Bool("verbose")); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: auto-sync failed: %v\n", err)
		}
	}

	var opts []tap.Option
	if localOnly {
		if err := os.MkdirAll(overrideDir, 0o755); err != nil {
			return nil, fmt.Errorf("create local override dir: %w", err)
		}
		opts = append(opts, tap.WithSitesDir(overrideDir))
	} else {
		opts = append(opts, tap.WithSitesDir(dir), tap.WithLocalOverrideDir(overrideDir))
	}
	if binary := cmd.String("agent-browser"); binary != "" {
		opts = append(opts, tap.WithAgentBrowserBinary(binary))
	}
	if timeout := cmd.Duration("timeout"); timeout > 0 {
		opts = append(opts, tap.WithTimeout(timeout))
	}
	return tap.New(ctx, opts...)
}

// newRuntimeClient skips site discovery and synchronization for commands that
// only need agent-browser. Fetch must remain usable when the registry is
// offline or absent.
func newRuntimeClient(ctx context.Context, cmd *cli.Command) (*tap.Client, error) {
	var opts []tap.Option
	if binary := cmd.String("agent-browser"); binary != "" {
		opts = append(opts, tap.WithAgentBrowserBinary(binary))
	}
	if timeout := cmd.Duration("timeout"); timeout > 0 {
		opts = append(opts, tap.WithTimeout(timeout))
	}
	return tap.New(ctx, opts...)
}

func configureLogging(cmd *cli.Command) {
	if cmd.Bool("quiet") || !cmd.Bool("verbose") {
		log.SetOutput(nopWriter{})
	}
}

type nopWriter struct{}

func (nopWriter) Write(data []byte) (int, error) { return len(data), nil }
