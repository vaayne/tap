package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserNetworkCmd() *cli.Command {
	return &cli.Command{
		Name:  "network",
		Usage: "Capture and intercept network requests in a tracked browser tab",
		Description: `Network observation and interception for tracked browser tabs.

Uses the CDP Network domain for passive capture (wait, log, body) and
the Fetch domain for active interception (intercept, clear).`,
		Commands: []*cli.Command{
			browserNetworkWaitCmd(),
			browserNetworkBodyCmd(),
			browserNetworkLogCmd(),
			browserNetworkInterceptCmd(),
			browserNetworkClearCmd(),
		},
	}
}

func browserNetworkWaitCmd() *cli.Command {
	return &cli.Command{
		Name:  "wait",
		Usage: "Wait for a network request matching filters",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:  "url-pattern",
				Usage: "Glob pattern to match request URLs (* matches any chars including /)",
			},
			&cli.StringFlag{
				Name:  "method",
				Usage: "HTTP method(s) to match, comma-separated (e.g. GET,POST)",
			},
			&cli.StringFlag{
				Name:  "resource-type",
				Usage: "Resource type(s) to match, comma-separated (e.g. XHR,Fetch,Document)",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Maximum time to wait for a matching request",
				Value: 30 * time.Second,
			},
			&cli.BoolFlag{
				Name:  "body",
				Usage: "Include the response body in the output",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		),
		Description: `Block until a network request matching the given filters completes,
then print the captured request/response entry.

The --url-pattern flag uses glob syntax where * matches any characters
including path separators. For example:
  */api/*         matches https://example.com/api/v1/users
  *.ads.*         matches https://tracker.ads.example.com/pixel

Examples:
  tap browser network wait --url-pattern "*/api/search*"
  tap browser network wait --url-pattern "*/graphql" --method POST --body
  tap browser network wait --resource-type XHR,Fetch --timeout 10s`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			filter := buildNetworkFilter(cmd)
			timeout := cmd.Duration("timeout")
			includeBody := cmd.Bool("body")

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			waitCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			entry, err := mgr.NetworkWait(waitCtx, sessionName, tabName, filter, includeBody)
			if err != nil {
				return err
			}

			return printResult(cmd, entry)
		},
	}
}

func browserNetworkBodyCmd() *cli.Command {
	return &cli.Command{
		Name:      "body",
		Usage:     "Fetch the response body for a completed request by ID",
		ArgsUsage: "<requestId>",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json (base64-encoded), raw (binary to stdout)",
				Value:   formatRaw,
			},
		),
		Description: `Fetch and print the response body for a network request identified
by its request ID (from 'tap browser network wait' or 'tap browser network log' output).

Examples:
  tap browser network body "12345.67"
  tap browser network body "12345.67" --format json`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			requestID := cmd.Args().First()
			if requestID == "" {
				return fmt.Errorf("request ID required")
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			body, err := mgr.NetworkGetBody(ctx, sessionName, tabName, requestID)
			if err != nil {
				return err
			}

			format := cmd.String("format")
			switch format {
			case formatJSON:
				// Output as base64-encoded JSON string.
				encoded := base64.StdEncoding.EncodeToString(body)
				fmt.Printf("%q\n", encoded)
			default: // raw
				_, _ = os.Stdout.Write(body)
			}
			return nil
		},
	}
}

// buildNetworkFilter constructs a NetworkFilter from CLI flags.
func buildNetworkFilter(cmd *cli.Command) browser.NetworkFilter {
	var filter browser.NetworkFilter

	if p := cmd.String("url-pattern"); p != "" {
		filter.URLPattern = p
	}
	if m := cmd.String("method"); m != "" {
		filter.Methods = splitCSV(m)
	}
	if rt := cmd.String("resource-type"); rt != "" {
		filter.ResourceTypes = splitCSV(rt)
	}

	return filter
}

// splitCSV splits a comma-separated string into trimmed non-empty parts.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func browserNetworkLogCmd() *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "Stream captured network requests as NDJSON",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:  "url-pattern",
				Usage: "Glob pattern to match request URLs",
			},
			&cli.StringFlag{
				Name:  "method",
				Usage: "HTTP method(s) to match, comma-separated",
			},
			&cli.StringFlag{
				Name:  "resource-type",
				Usage: "Resource type(s) to match, comma-separated",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Maximum duration to capture (0 = until interrupted)",
				Value: 0,
			},
		),
		Description: `Enable the Network domain and stream completed request/response entries
as newline-delimited JSON (NDJSON) to stdout.

Runs until --timeout expires or the process is interrupted (Ctrl-C).

Examples:
  tap browser network log
  tap browser network log --url-pattern "*/api/*" --timeout 30s
  tap browser network log --resource-type XHR,Fetch`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			filter := buildNetworkFilter(cmd)

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			logCtx := ctx
			var logCancel context.CancelFunc
			if timeout := cmd.Duration("timeout"); timeout > 0 {
				logCtx, logCancel = context.WithTimeout(ctx, timeout)
				defer logCancel()
			}

			ch, cancel, err := mgr.NetworkLog(logCtx, sessionName, tabName, filter)
			if err != nil {
				return err
			}
			defer cancel()

			enc := json.NewEncoder(os.Stdout)
			for entry := range ch {
				if err := enc.Encode(entry); err != nil {
					return fmt.Errorf("write entry: %w", err)
				}
			}
			return nil
		},
	}
}

func browserNetworkInterceptCmd() *cli.Command {
	return &cli.Command{
		Name:  "intercept",
		Usage: "Set request interception rules (block, mock, or modify headers)",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:  "url-pattern",
				Usage: "Glob pattern to match request URLs",
			},
			&cli.StringFlag{
				Name:  "method",
				Usage: "HTTP method(s) to match, comma-separated",
			},
			&cli.StringFlag{
				Name:  "resource-type",
				Usage: "Resource type(s) to match, comma-separated",
			},
			&cli.BoolFlag{
				Name:  "block",
				Usage: "Block matching requests (mutually exclusive with --respond)",
			},
			&cli.StringSliceFlag{
				Name:  "header",
				Usage: `Add/override request header (repeatable, format "Key: Value")`,
			},
			&cli.StringFlag{
				Name:  "respond",
				Usage: "Mock response body (mutually exclusive with --block)",
			},
			&cli.IntFlag{
				Name:  "status",
				Usage: "Mock response HTTP status code (required with --respond)",
				Value: 200,
			},
			&cli.StringFlag{
				Name:  "content-type",
				Usage: "Mock response Content-Type header",
				Value: "application/json",
			},
		),
		Description: `Set Fetch domain interception rules on a tracked tab.

Rules are replace-all: each call replaces any previously set rules.
Pass a single rule per invocation. Use 'tap browser network clear' to
remove all rules.

This command keeps running to serve interception rules until interrupted
(Ctrl-C). Scripts should run it in the background.

Examples:
  # Block ad requests
  tap browser network intercept --block --url-pattern "*.ads.*"

  # Mock an API response
  tap browser network intercept --url-pattern "*/api/user" \\
    --respond '{"name":"test"}' --status 200

  # Add auth header to API requests
  tap browser network intercept --url-pattern "*/api/*" \\
    --header "Authorization: Bearer tok_abc123"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			filter := buildNetworkFilter(cmd)
			block := cmd.Bool("block")
			respondBody := cmd.String("respond")

			if block && respondBody != "" {
				return fmt.Errorf("--block and --respond are mutually exclusive")
			}

			rule := browser.InterceptRule{
				Filter: filter,
				Block:  block,
			}

			if respondBody != "" {
				rule.MockBody = respondBody
				rule.MockStatus = cmd.Int("status")
				rule.MockHeaders = map[string]string{
					"Content-Type": cmd.String("content-type"),
				}
			}

			// Parse --header flags.
			for _, h := range cmd.StringSlice("header") {
				parts := strings.SplitN(h, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid header format %q (expected \"Key: Value\")", h)
				}
				if rule.AddHeaders == nil {
					rule.AddHeaders = make(map[string]string)
				}
				rule.AddHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			if err := mgr.NetworkIntercept(ctx, sessionName, tabName, []browser.InterceptRule{rule}); err != nil {
				return err
			}

			if block {
				fmt.Fprintln(os.Stderr, "Blocking matching requests")
			} else if respondBody != "" {
				fmt.Fprintf(os.Stderr, "Mocking matching requests with status %d\n", rule.MockStatus)
			} else if len(rule.AddHeaders) > 0 {
				fmt.Fprintln(os.Stderr, "Adding headers to matching requests")
			}

			// Keep the process alive so the interception goroutine stays active.
			<-ctx.Done()
			return nil
		},
	}
}

func browserNetworkClearCmd() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Remove all Fetch domain interception rules",
		Flags: browserActionFlags(false),
		Description: `Disable the Fetch domain and remove all interception rules
from the resolved tracked tab.

This does not affect passive Network domain capture (log/wait).`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}

			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			if err := mgr.NetworkClearIntercept(ctx, sessionName, tabName); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Interception rules cleared")
			return nil
		},
	}
}
