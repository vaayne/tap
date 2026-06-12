package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
	"github.com/vaayne/tap/browser"
	"github.com/vaayne/tap/transport"
)

var version = "dev"

func main() {
	_ = godotenv.Load()

	app := newApp()

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

func newApp() *cli.Command {
	return &cli.Command{
		Name:                            "tap",
		Usage:                           "Tap into any website from your terminal",
		Version:                         version,
		EnableShellCompletion:           true,
		ConfigureShellCompletionCommand: configureCompletionCommand,
		Description: `Tap runs site-specific JS scripts against websites and extracts clean content
from URLs. Scripts execute in QuickJS (fast, no browser) with automatic fallback
to a real Chrome browser when auth or JS rendering is needed.

Quick start:
  tap site list                          List available site scripts
  tap site hackernews/top                Run a site script
  tap fetch https://example.com          Extract clean readable content
  tap attach chrome                      Reuse your existing Chrome
  tap browser open https://example.com   Open a page in the default browser context
  tap browser open https://github.com/login --show
  tap status                             Show the active browser context + current tab

Use 'tap <command> --help' for details on any command.`,
		Flags: globalFlags(),
		Commands: []*cli.Command{
			siteCmd(),
			fetchCmd(),
			browserCmd(),
			attachCmd(),
			statusCmd(),
			doctorCmd(),
			upgradeCmd(),
			skillCmd(),
			internalCmd(),
			docsCmd(),
		},
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
			Aliases: []string{"browser-url", "devtools-url"},
			Usage:   "Remote Chrome DevTools endpoint override (advanced)",
			Sources: cli.EnvVars("TAP_WS_URL"),
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:    "profile-dir",
			Usage:   "Chrome profile directory override (advanced)",
			Sources: cli.EnvVars("TAP_PROFILE_DIR"),
			Hidden:  true,
		},
		&cli.BoolFlag{
			Name:    "browser",
			Aliases: []string{"b"},
			Usage:   "Force browser execution, skip QuickJS",
			Sources: cli.EnvVars("TAP_BROWSER"),
			Hidden:  true,
		},
		&cli.BoolFlag{
			Name:    "no-headless",
			Aliases: []string{"show"},
			Usage:   "Run browser in visible mode (advanced)",
			Hidden:  true,
		},
		&cli.BoolFlag{
			Name:    "lightpanda",
			Aliases: []string{"lp"},
			Usage:   "Use Lightpanda instead of Chrome (advanced)",
			Sources: cli.EnvVars("TAP_LIGHTPANDA"),
			Hidden:  true,
		},
		&cli.BoolFlag{
			Name:   "pause",
			Usage:  "Pause after navigation for manual interaction (advanced)",
			Hidden: true,
		},
		&cli.DurationFlag{
			Name:    "delay",
			Aliases: []string{"wait"},
			Usage:   "Wait after navigation (advanced)",
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:   "wait-selector",
			Usage:  "Wait until a CSS selector becomes visible (advanced)",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "wait-js",
			Usage:  "Wait until a JavaScript expression becomes truthy (advanced)",
			Hidden: true,
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout (advanced)",
			Value:   0,
			Sources: cli.EnvVars("TAP_TIMEOUT"),
			Hidden:  true,
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

	resolvedOpts, err := resolveBrowserClientOptions(ctx, cmd, false)
	if err != nil {
		return nil, err
	}
	opts = append(opts, resolvedOpts...)

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
	if cmd.Bool("no-headless") || cmd.Bool("show") || cmd.Bool("pause") {
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

func resolveBrowserClientOptions(ctx context.Context, cmd *cli.Command, forceManagedDefault bool) ([]tap.Option, error) {
	var opts []tap.Option

	if url := firstString(cmd, "browser-url", "ws-url", "devtools-url"); url != "" {
		return append(opts, tap.WithWSURL(url)), nil
	}
	if dir := firstString(cmd, "profile-dir"); dir != "" {
		return append(opts, tap.WithProfileDir(dir)), nil
	}

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return nil, err
	}
	defaultContext, err := mgr.DefaultContext(ctx)
	if err != nil {
		return nil, err
	}
	store, err := newBrowserStore(cmd)
	if err != nil {
		return nil, err
	}
	state, err := store.Load()
	if err != nil {
		return nil, err
	}
	if defaultContext != nil && defaultContext.Kind == browser.DefaultContextAttached {
		if state.ProxyDaemon == nil {
			return nil, fmt.Errorf("attached Chrome is stale: proxy daemon metadata is missing (run 'tap attach chrome')")
		}
		health := browser.CheckProxyDaemon(ctx, state.ProxyDaemon)
		if err := persistProxyDaemonHealth(store, state.ProxyDaemon, health); err == nil && !health.Healthy {
			defaultContext.Reason = health.Reason
			defaultContext.Stale = true
		}
		if !health.Healthy {
			return nil, fmt.Errorf("attached Chrome is stale: %s (run 'tap attach chrome')", health.Reason)
		}
	}
	if defaultContext != nil && defaultContext.Stale {
		return nil, fmt.Errorf("default browser context %q is stale: %s (run 'tap attach chrome')", defaultContext.SessionName, defaultContext.Reason)
	}

	session, err := mgr.GetSession(ctx, "")
	if err == nil {
		switch {
		case session.Remote != nil && session.Remote.WSURL != "":
			return append(opts, tap.WithWSURL(session.Remote.WSURL)), nil
		case session.Local != nil && session.Local.ProfileDir != "":
			return append(opts, tap.WithProfileDir(session.Local.ProfileDir)), nil
		}
	}
	if defaultContext != nil {
		return nil, fmt.Errorf("resolve default browser context %q: %w", defaultContext.SessionName, err)
	}
	if forceManagedDefault {
		return append(opts, tap.WithProfileDir(defaultManagedProfileDir(cmd))), nil
	}
	return opts, nil
}

func defaultManagedProfileDir(cmd *cli.Command) string {
	return filepath.Join(browserStateRoot(cmd), "profiles", browser.DefaultSessionName)
}

func firstString(cmd *cli.Command, names ...string) string {
	for _, name := range names {
		if value := cmd.String(name); value != "" {
			return value
		}
	}
	return ""
}

func browserClientFlags(includeLightpanda bool) []cli.Flag {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "browser",
			Aliases: []string{"b"},
			Usage:   "Use the browser path and reuse the resolved browser context",
		},
		&cli.BoolFlag{
			Name:    "show",
			Aliases: []string{"no-headless"},
			Usage:   "Run the browser visibly",
		},
		&cli.DurationFlag{
			Name:    "wait",
			Aliases: []string{"delay"},
			Usage:   "Wait a fixed duration after navigation before continuing",
		},
		&cli.StringFlag{
			Name:  "wait-selector",
			Usage: "Wait until a CSS selector becomes visible before continuing",
		},
		&cli.StringFlag{
			Name:  "wait-js",
			Usage: "Wait until a JavaScript expression becomes truthy before continuing",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout; 0 means no timeout",
		},
		&cli.StringFlag{
			Name:    "browser-url",
			Aliases: []string{"ws-url", "devtools-url"},
			Usage:   "One-shot DevTools endpoint override",
		},
		&cli.StringFlag{
			Name:  "profile-dir",
			Usage: "One-shot browser profile override",
		},
	}
	if includeLightpanda {
		flags = append(flags, &cli.BoolFlag{
			Name:    "lightpanda",
			Aliases: []string{"lp"},
			Usage:   "Use Lightpanda instead of Chrome (implies --browser)",
		})
	}
	return flags
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
