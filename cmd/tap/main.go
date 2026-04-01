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

	app := &cli.Command{
		Name:    "tap",
		Usage:   "Tap into any website from your terminal",
		Version: version,
		Flags:   globalFlags(),
		Commands: []*cli.Command{
			siteCmd(),
			fetchCmd(),
			loginCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		if snf, ok := err.(*tap.ScriptNotFoundError); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", snf.Error())
			if suggestions := snf.Suggestions(5); len(suggestions) > 0 {
				fmt.Fprintf(os.Stderr, "\nDid you mean?\n")
				for _, s := range suggestions {
					fmt.Fprintf(os.Stderr, "  tap site %s\n", s)
				}
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "sites-dir",
			Usage:   "Directory containing site scripts",
			Value:   defaultSitesDir(),
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
		&cli.BoolFlag{
			Name:    "browser",
			Aliases: []string{"b"},
			Usage:   "Force browser execution, skip QuickJS",
			Sources: cli.EnvVars("TAP_BROWSER"),
		},
		&cli.BoolFlag{
			Name:  "no-headless",
			Usage: "Run browser in visible mode (useful for debugging auth)",
		},
		&cli.BoolFlag{
			Name:  "pause",
			Usage: "Pause after navigation for manual interaction (login, CAPTCHA)",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout (e.g., 30s, 2m)",
			Value:   0,
			Sources: cli.EnvVars("TAP_TIMEOUT"),
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging",
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "Suppress all log output",
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "Disable colored output",
			Sources: cli.EnvVars("NO_COLOR"),
		},
	}
}

func newClient(cmd *cli.Command) (*tap.Client, error) {
	var opts []tap.Option

	dir := cmd.String("sites-dir")
	if dir == "" {
		dir = defaultSitesDir()
	}

	// Auto-sync if no local scripts exist
	if err := ensureScripts(dir, cmd.Bool("verbose")); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: auto-sync failed: %v\n", err)
	}

	opts = append(opts, tap.WithSitesDir(dir))
	if url := cmd.String("ws-url"); url != "" {
		opts = append(opts, tap.WithWSURL(url))
	}
	if dir := cmd.String("profile-dir"); dir != "" {
		opts = append(opts, tap.WithProfileDir(dir))
	}
	if cmd.Bool("browser") {
		opts = append(opts, tap.WithForceBrowser(true))
	}
	if cmd.Bool("no-headless") {
		opts = append(opts, tap.WithHeadless(false))
	}
	if cmd.Bool("pause") {
		opts = append(opts, tap.WithHeadless(false)) // --pause implies visible browser
		opts = append(opts, tap.WithForceBrowser(true)) // --pause implies browser mode
		opts = append(opts, tap.WithPause(terminalPause()))
	}
	if d := cmd.Duration("timeout"); d > 0 {
		opts = append(opts, tap.WithTimeout(d))
	}

	return tap.New(opts...)
}

// newClientWithOverrides creates a client with forced overrides (e.g. for login).
func newClientWithOverrides(cmd *cli.Command, forceVisible bool) (*tap.Client, error) {
	var opts []tap.Option

	dir := cmd.String("sites-dir")
	if dir == "" {
		dir = defaultSitesDir()
	}
	opts = append(opts, tap.WithSitesDir(dir))

	if url := cmd.String("ws-url"); url != "" {
		opts = append(opts, tap.WithWSURL(url))
	}
	if dir := cmd.String("profile-dir"); dir != "" {
		opts = append(opts, tap.WithProfileDir(dir))
	}
	if forceVisible {
		opts = append(opts, tap.WithHeadless(false))
	}
	if d := cmd.Duration("timeout"); d > 0 {
		opts = append(opts, tap.WithTimeout(d))
	}

	return tap.New(opts...)
}

// configureLogging sets up log output based on --verbose/--quiet flags.
func configureLogging(cmd *cli.Command) {
	if cmd.Bool("quiet") {
		log.SetOutput(nopWriter{})
	} else if !cmd.Bool("verbose") {
		// Default: suppress log output (only errors via stderr)
		log.SetOutput(nopWriter{})
	}
	// verbose: keep default log output (stderr)
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }
