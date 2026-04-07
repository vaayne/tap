package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
	"github.com/vaayne/tap/transport"
)

var version = "dev"

func main() {
	_ = godotenv.Load()

	app := &cli.Command{
		Name:    "tap",
		Usage:   "Tap into any website from your terminal",
		Version: version,
		Description: `Tap runs site-specific JS scripts against websites and extracts clean content
from URLs. Scripts execute in QuickJS (fast, no browser) with automatic fallback
to a real Chrome browser when auth or JS rendering is needed.

Quick start:
  tap fetch https://example.com          Extract article as clean Markdown
  tap fetch --json https://example.com   Full metadata as JSON
  tap site list                          List available site scripts
  tap site hackernews/top                Run a script
  tap site -b twitter/search query=go    Run with browser cookies (auth)
  tap login https://github.com/login     Log in via visible browser, save cookies

Browser automation:
  tap browser session new default        Start a persistent browser
  tap browser tab new main --url https://example.com
  tap browser text                       Extract clean page text (token-efficient)
  tap browser screenshot                 Capture the current page
  tap browser click "button.submit"      Interact with elements
  tap browser network wait --url-pattern "*/api/*" --body

Use 'tap <command> --help' for details on any command.`,
		Flags: globalFlags(),
		Commands: []*cli.Command{
			browserCmd(),
			electronCmd(),
			siteCmd(),
			fetchCmd(),
			loginCmd(),
			upgradeCmd(),
			doctorCmd(),
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
			Usage:   "Remote Chrome DevTools endpoint: WebSocket URL or HTTP base URL (e.g. ws://localhost:9222/devtools/browser/... or http://localhost:9222)",
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
			Name:    "lightpanda",
			Aliases: []string{"lp"},
			Usage:   "Use Lightpanda (lightweight headless browser) instead of Chrome (implies --browser)",
			Sources: cli.EnvVars("TAP_LIGHTPANDA"),
		},
		&cli.BoolFlag{
			Name:  "pause",
			Usage: "Pause after navigation for manual interaction (requires interactive terminal)",
		},
		&cli.DurationFlag{
			Name:  "delay",
			Usage: "Wait a fixed duration after navigation before continuing (implies --browser)",
		},
		&cli.StringFlag{
			Name:  "wait-selector",
			Usage: "Wait until a CSS selector becomes visible before continuing (implies --browser)",
		},
		&cli.StringFlag{
			Name:  "wait-js",
			Usage: "Wait until a JavaScript expression becomes truthy before continuing (implies --browser)",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout; 0 means no timeout (e.g. 30s, 2m)",
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
			Name:  "local-only",
			Usage: "Only use scripts from local override dir (~/.config/tap/sites/), skip cache",
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "Disable colored output",
			Sources: cli.EnvVars("NO_COLOR"),
		},
	}
}

func newClient(ctx context.Context, cmd *cli.Command) (*tap.Client, error) {
	var opts []tap.Option

	dir := cmd.String("sites-dir")
	if dir == "" {
		dir = defaultSitesDir()
	}

	localOnly := cmd.Bool("local-only")
	localOverrideDir := defaultLocalOverrideDir()

	// Auto-sync if no local scripts exist (skip when --local-only)
	if !localOnly {
		if err := ensureScripts(dir, cmd.Bool("verbose")); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: auto-sync failed: %v\n", err)
		}
	}

	if localOnly {
		// In local-only mode use the override dir as the sole scripts dir,
		// and disable the cache path so only local scripts are visible.
		if err := os.MkdirAll(localOverrideDir, 0o755); err != nil {
			return nil, fmt.Errorf("create local override dir: %w", err)
		}
		opts = append(opts, tap.WithSitesDir(localOverrideDir))
	} else {
		opts = append(opts, tap.WithSitesDir(dir))
		opts = append(opts, tap.WithLocalOverrideDir(localOverrideDir))
	}
	if url := cmd.String("ws-url"); url != "" {
		opts = append(opts, tap.WithWSURL(url))
	}
	if dir := cmd.String("profile-dir"); dir != "" {
		opts = append(opts, tap.WithProfileDir(dir))
	}
	pauseFn, err := resolvePauseFunc(cmd)
	if err != nil {
		return nil, err
	}

	if cmd.Bool("lightpanda") {
		opts = append(opts, tap.WithBrowserType(transport.BrowserLightpanda))
		opts = append(opts, tap.WithForceBrowser(true))
	}
	if cmd.Bool("browser") || hasPauseMode(cmd) {
		opts = append(opts, tap.WithForceBrowser(true))
	}
	if cmd.Bool("no-headless") || cmd.Bool("pause") {
		opts = append(opts, tap.WithHeadless(false))
	}
	if pauseFn != nil {
		opts = append(opts, tap.WithPause(pauseFn))
	}
	if d := cmd.Duration("timeout"); d > 0 {
		opts = append(opts, tap.WithTimeout(d))
	}

	return tap.New(ctx, opts...)
}

// newClientWithOverrides creates a client with forced overrides (e.g. for login).
func newClientWithOverrides(ctx context.Context, cmd *cli.Command, forceVisible bool) (*tap.Client, error) {
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
	if cmd.Bool("lightpanda") {
		opts = append(opts, tap.WithBrowserType(transport.BrowserLightpanda))
	}
	if forceVisible {
		opts = append(opts, tap.WithHeadless(false))
	}
	if d := cmd.Duration("timeout"); d > 0 {
		opts = append(opts, tap.WithTimeout(d))
	}

	return tap.New(ctx, opts...)
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
