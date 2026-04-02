package main

import (
	"context"
	"encoding/base64"
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
