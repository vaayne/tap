package main

import (
	"context"
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
		Usage: "Capture and intercept network requests",
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
				Usage: "Glob pattern to match request URLs",
			},
			&cli.StringFlag{
				Name:  "method",
				Usage: "HTTP method(s) to match, comma-separated",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Maximum time to wait for a matching request",
				Value: 30 * time.Second,
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			args := []string{"network", "wait", "--json"}
			if p := cmd.String("url-pattern"); p != "" {
				args = append(args, "--url-pattern", p)
			}
			if m := cmd.String("method"); m != "" {
				args = append(args, "--method", m)
			}
			waitCtx, cancel := context.WithTimeout(ctx, cmd.Duration("timeout"))
			defer cancel()
			out, _, err := ab.Exec(waitCtx, args...)
			if err != nil {
				return err
			}
			var envelope browser.AgentBrowserEnvelope[map[string]any]
			if err := json.Unmarshal(out, &envelope); err != nil {
				return fmt.Errorf("parse network wait: %w", err)
			}
			if !envelope.Success {
				return fmt.Errorf("network wait: %s", envelope.Error)
			}
			return printResult(cmd, envelope.Data)
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			requestID := cmd.Args().First()
			if requestID == "" {
				return fmt.Errorf("request ID required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			out, _, err := ab.Exec(ctx, "network", "body", requestID, "--json")
			if err != nil {
				return err
			}
			format := cmd.String("format")
			if format == formatJSON {
				fmt.Println(string(out))
			} else {
				var envelope browser.AgentBrowserEnvelope[map[string]any]
				_ = json.Unmarshal(out, &envelope)
				if data, ok := envelope.Data["body"].(string); ok {
					_, _ = os.Stdout.Write([]byte(data))
				} else {
					_, _ = os.Stdout.Write(out)
				}
			}
			return nil
		},
	}
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
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Maximum duration to capture (0 = until interrupted)",
				Value: 0,
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			args := []string{"network", "log"}
			if p := cmd.String("url-pattern"); p != "" {
				args = append(args, "--url-pattern", p)
			}
			if m := cmd.String("method"); m != "" {
				args = append(args, "--method", m)
			}
			logCtx := ctx
			var logCancel context.CancelFunc
			if timeout := cmd.Duration("timeout"); timeout > 0 {
				logCtx, logCancel = context.WithTimeout(ctx, timeout)
				defer logCancel()
			}
			out, _, err := ab.Exec(logCtx, args...)
			if err != nil {
				return err
			}
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if line = strings.TrimSpace(line); line != "" {
					fmt.Println(line)
				}
			}
			return nil
		},
	}
}

func browserNetworkInterceptCmd() *cli.Command {
	return &cli.Command{
		Name:  "intercept",
		Usage: "Set request interception rules",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:  "url-pattern",
				Usage: "Glob pattern to match request URLs",
			},
			&cli.StringFlag{
				Name:  "method",
				Usage: "HTTP method(s) to match, comma-separated",
			},
			&cli.BoolFlag{
				Name:  "block",
				Usage: "Block matching requests",
			},
			&cli.StringFlag{
				Name:  "respond",
				Usage: "Mock response body",
			},
			&cli.IntFlag{
				Name:  "status",
				Usage: "Mock response HTTP status code",
				Value: 200,
			},
			&cli.StringFlag{
				Name:  "content-type",
				Usage: "Mock response Content-Type header",
				Value: "application/json",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			args := []string{"network", "intercept"}
			if p := cmd.String("url-pattern"); p != "" {
				args = append(args, "--url-pattern", p)
			}
			if m := cmd.String("method"); m != "" {
				args = append(args, "--method", m)
			}
			if cmd.Bool("block") {
				args = append(args, "--block")
			}
			if body := cmd.String("respond"); body != "" {
				args = append(args, "--respond", body, "--status", fmt.Sprintf("%d", cmd.Int("status")), "--content-type", cmd.String("content-type"))
			}
			_, _, err = ab.Exec(ctx, args...)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Interception rules set")
			<-ctx.Done()
			return nil
		},
	}
}

func browserNetworkClearCmd() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Remove all network interception rules",
		Flags: browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "network", "clear")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Interception rules cleared")
			return nil
		},
	}
}
