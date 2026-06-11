package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

// browserWaitCmd returns the "tap browser wait" command.
//
// Exactly one wait mode must be specified per invocation:
//
//	tap browser wait <selector>                          element visible (default)
//	tap browser wait <ms-or-duration>                   pure time wait
//	tap browser wait <selector> --state visible|hidden|attached|detached
//	tap browser wait --text "Welcome"                   body text contains substring
//	tap browser wait --url "**/dash"                    URL glob match
//	tap browser wait --load load|domcontentloaded|networkidle
//	tap browser wait --fn "js expression"               poll until truthy
//	tap browser wait --timeout 30s                      default timeout 30s
func browserWaitCmd() *cli.Command {
	return &cli.Command{
		Name:      "wait",
		Usage:     "Wait for a page condition to become true",
		ArgsUsage: "[<selector-or-duration>]",
		Flags: append(browserActionFlags(false),
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Maximum time to wait (default 30s)",
				Value: 30 * time.Second,
			},
			&cli.StringFlag{
				Name:  "state",
				Usage: "Element state: visible (default), hidden, attached, detached",
				Value: "visible",
			},
			&cli.StringFlag{
				Name:  "text",
				Usage: "Wait until document.body.innerText contains this substring",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: `Wait until location.href matches this glob (supports * and **)`,
			},
			&cli.StringFlag{
				Name:  "load",
				Usage: "Wait for page load event: load, domcontentloaded, networkidle",
			},
			&cli.StringFlag{
				Name:  "fn",
				Usage: "Poll until this JS expression evaluates to a truthy value",
			},
		),
		Description: `Wait for a page condition before proceeding.

Modes (exactly one must be active):

  POSITIONAL ARGUMENT
    <selector>          Wait for CSS selector to become visible (default state).
    <ms>                Pure time wait when the argument is a plain integer (ms).
    <duration>          Pure time wait when the argument is a Go duration (e.g. 2s).

  FLAGS
    --state             Change element state: visible (default), hidden, attached, detached.
                        Requires a positional CSS selector argument.
    --text              Wait until document.body.innerText contains the substring.
    --url               Wait until location.href matches the glob pattern.
                        Supports * (any chars except /) and ** (any chars including /).
    --load              Wait for a named page-load event:
                          load             — fires when the load event completes
                          domcontentloaded — fires when DOMContentLoaded fires
                          networkidle      — no new network requests for ~500 ms
    --fn                Poll until the JS expression returns a truthy value.

  --timeout             Maximum wait duration (default 30s). Applies to all modes.

Examples:
  tap browser wait "#login-form"
  tap browser wait ".spinner" --state hidden
  tap browser wait 2000
  tap browser wait 1.5s
  tap browser wait --text "Welcome back"
  tap browser wait --url "**/dashboard"
  tap browser wait --load networkidle
  tap browser wait --fn "window.__ready === true"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			return runBrowserWait(ctx, cmd)
		},
	}
}

func runBrowserWait(ctx context.Context, cmd *cli.Command) error {
	timeout := cmd.Duration("timeout")
	textFlag := cmd.String("text")
	urlFlag := cmd.String("url")
	loadFlag := cmd.String("load")
	fnFlag := cmd.String("fn")
	stateFlag := cmd.String("state")
	arg := cmd.Args().First()

	// Count how many modes are active to detect conflicts.
	modeCount := 0
	if textFlag != "" {
		modeCount++
	}
	if urlFlag != "" {
		modeCount++
	}
	if loadFlag != "" {
		modeCount++
	}
	if fnFlag != "" {
		modeCount++
	}
	if arg != "" {
		modeCount++
	}

	if modeCount > 1 {
		return fmt.Errorf("wait: only one mode may be active per invocation; got multiple arguments/flags")
	}
	if modeCount == 0 {
		return fmt.Errorf("wait: specify a selector, duration, --text, --url, --load, or --fn")
	}

	// Pure time wait: integer (ms) or Go duration string.
	if arg != "" {
		if d, ok := parseDurationArg(arg); ok {
			return timeSleep(ctx, d)
		}
		// Not a duration — treat as CSS selector below.
	}

	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}
	session := cmd.String("session")
	tab := cmd.String("tab")

	switch {
	case textFlag != "":
		if err := mgr.WaitForText(ctx, session, tab, textFlag, timeout); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Text %q found\n", textFlag)

	case urlFlag != "":
		if err := mgr.WaitForURL(ctx, session, tab, urlFlag, timeout); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "URL matched %q\n", urlFlag)

	case loadFlag != "":
		ls, err := parseLoadState(loadFlag)
		if err != nil {
			return err
		}
		if err := mgr.WaitForLoad(ctx, session, tab, ls, timeout); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Load state %q reached\n", loadFlag)

	case fnFlag != "":
		if err := mgr.WaitForFn(ctx, session, tab, fnFlag, timeout); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "JS condition is truthy")

	default:
		// Positional arg is a CSS selector.
		sel := arg
		es, err := parseElementState(stateFlag)
		if err != nil {
			return err
		}
		if err := mgr.WaitForElement(ctx, session, tab, sel, es, timeout); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Element %s is %s\n", sel, stateFlag)
	}

	return nil
}

// timeSleep blocks for d or until ctx is cancelled.
func timeSleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
	}
	fmt.Fprintf(os.Stderr, "Waited %s\n", d)
	return nil
}

// parseDurationArg tries to interpret s as a time wait:
//   - plain integer → milliseconds
//   - Go duration string (e.g. "1.5s", "500ms") → parsed duration
//
// Returns (duration, true) on success or (0, false) if s is not a duration.
func parseDurationArg(s string) (time.Duration, bool) {
	// Plain integer → milliseconds.
	if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Duration(ms) * time.Millisecond, true
	}
	// Go duration string.
	if d, err := time.ParseDuration(s); err == nil {
		return d, true
	}
	return 0, false
}

// parseElementState maps the --state flag value to an ElementState constant.
func parseElementState(s string) (browser.ElementState, error) {
	switch s {
	case "visible", "":
		return browser.ElementVisible, nil
	case "hidden":
		return browser.ElementHidden, nil
	case "attached":
		return browser.ElementAttached, nil
	case "detached":
		return browser.ElementDetached, nil
	default:
		return browser.ElementVisible, fmt.Errorf("wait: unknown --state %q; choose visible, hidden, attached, or detached", s)
	}
}

// parseLoadState maps the --load flag value to a LoadState constant.
func parseLoadState(s string) (browser.LoadState, error) {
	switch s {
	case string(browser.LoadStateLoad):
		return browser.LoadStateLoad, nil
	case string(browser.LoadStateDOMContentLoaded):
		return browser.LoadStateDOMContentLoaded, nil
	case string(browser.LoadStateNetworkIdle):
		return browser.LoadStateNetworkIdle, nil
	default:
		return browser.LoadStateLoad, fmt.Errorf("wait: unknown --load %q; choose load, domcontentloaded, or networkidle", s)
	}
}
