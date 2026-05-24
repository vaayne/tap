package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
	defuddle "github.com/vaayne/go-defuddle"
	"github.com/vaayne/tap/browser"
)

func browserNavigateCmd() *cli.Command {
	return &cli.Command{
		Name:      "navigate",
		Usage:     "Navigate a tracked browser tab to a URL",
		ArgsUsage: "<url>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			url := cmd.Args().First()
			if url == "" {
				return fmt.Errorf("URL required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			if err := ab.Open(ctx, url, browser.OpenOpts{}); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Navigated to %s\n", url)
			return nil
		},
	}
}

func browserEvaluateCmd() *cli.Command {
	return &cli.Command{
		Name:      "evaluate",
		Usage:     "Run JavaScript in a tracked browser tab",
		ArgsUsage: "<javascript>",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			js := cmd.Args().First()
			if js == "" {
				return fmt.Errorf("JavaScript expression required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			result, err := ab.Eval(ctx, js)
			if err != nil {
				return err
			}
			return printResult(cmd, result)
		},
	}
}

func browserScreenshotCmd() *cli.Command {
	return &cli.Command{
		Name:      "screenshot",
		Usage:     "Capture a screenshot from a tracked browser tab",
		ArgsUsage: "",
		Flags:     browserActionFlags(true),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}

			outPath := cmd.String("output")
			if outPath == "" {
				outPath = fmt.Sprintf("screenshot-%d.png", time.Now().Unix())
			}

			_, _, err = ab.Exec(ctx, "screenshot", outPath)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", outPath)
			return nil
		},
	}
}

func browserTextCmd() *cli.Command {
	return &cli.Command{
		Name:      "text",
		Usage:     "Extract clean readable text from the page via defuddle",
		ArgsUsage: "[selector]",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()

			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}

			html, err := ab.GetHTML(ctx)
			if err != nil {
				return err
			}
			if html == "" {
				fmt.Fprintln(os.Stderr, "No content found.")
				return nil
			}

			if sel != "" {
				val, err := ab.Eval(ctx, fmt.Sprintf("document.querySelector(%q)?.outerHTML || ''", sel))
				if err != nil {
					return err
				}
				if s, ok := val.(string); ok {
					html = s
				}
			}

			// Get current URL for defuddle
			urlVal, _ := ab.Eval(ctx, "window.location.href")
			pageURL := ""
			if s, ok := urlVal.(string); ok {
				pageURL = s
			}

			parser, err := defuddle.NewParser()
			if err != nil {
				return fmt.Errorf("init parser: %w", err)
			}
			defer parser.Close()
			dr, err := parser.Parse(html, pageURL, &defuddle.Options{Markdown: true})
			if err != nil {
				return fmt.Errorf("parse content: %w", err)
			}

			format := cmd.String("format")
			if format == "json" || format == "raw" {
				return printResult(cmd, map[string]any{
					"title":       dr.Title,
					"description": dr.Description,
					"markdown":    dr.Markdown,
					"wordCount":   dr.WordCount,
					"url":         pageURL,
				})
			}
			content := dr.Markdown
			if content == "" {
				content = dr.Content
			}
			if dr.Title != "" {
				fmt.Printf("# %s\n\n", dr.Title)
			}
			fmt.Println(content)
			return nil
		},
	}
}

func browserSnapshotCmd() *cli.Command {
	return &cli.Command{
		Name:  "snapshot",
		Usage: "Capture an AI-friendly semantic page snapshot with stable refs",
		Flags: append(browserActionFlags(false),
			&cli.BoolFlag{
				Name:  "interactive",
				Usage: "Only include interactive nodes (button/link/input/etc)",
			},
			&cli.StringFlag{
				Name:  "selector",
				Usage: "Optional scope selector (reserved for future use)",
			},
			&cli.IntFlag{
				Name:  "depth",
				Usage: "Maximum AX tree depth to capture",
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
			args := []string{"snapshot", "--json"}
			if cmd.Bool("interactive") {
				args = append(args, "--interactive")
			}
			if sel := cmd.String("selector"); sel != "" {
				args = append(args, "--selector", sel)
			}
			if d := cmd.Int("depth"); d > 0 {
				args = append(args, "--depth", strconv.Itoa(d))
			}
			out, _, err := ab.Exec(ctx, args...)
			if err != nil {
				return err
			}
			var envelope browser.AgentBrowserEnvelope[map[string]any]
			if err := json.Unmarshal(out, &envelope); err != nil {
				return fmt.Errorf("parse snapshot: %w", err)
			}
			if !envelope.Success {
				return fmt.Errorf("snapshot: %s", envelope.Error)
			}
			return printResult(cmd, envelope.Data)
		},
	}
}

func browserPDFCmd() *cli.Command {
	return &cli.Command{
		Name:  "pdf",
		Usage: "Save the current page as PDF",
		Flags: append(browserActionFlags(true),
			&cli.BoolFlag{
				Name:  "landscape",
				Usage: "Use landscape orientation",
			},
			&cli.BoolFlag{
				Name:  "background",
				Usage: "Print background graphics",
				Value: true,
			},
			&cli.Float64Flag{
				Name:  "scale",
				Usage: "Scale of the webpage rendering (default 1.0)",
				Value: 1.0,
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}

			outPath := cmd.String("output")
			if outPath == "" {
				outPath = fmt.Sprintf("page-%d.pdf", time.Now().Unix())
			}

			args := []string{"pdf", outPath}
			if cmd.Bool("landscape") {
				args = append(args, "--landscape")
			}
			if !cmd.Bool("background") {
				args = append(args, "--no-background")
			}
			if s := cmd.Float64("scale"); s != 1.0 {
				args = append(args, "--scale", strconv.FormatFloat(s, 'f', -1, 64))
			}

			_, _, err = ab.Exec(ctx, args...)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s\n", outPath)
			return nil
		},
	}
}

func browserFormsCmd() *cli.Command {
	return &cli.Command{
		Name:  "forms",
		Usage: "Discover fillable form elements in a tracked browser tab",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			out, _, err := ab.Exec(ctx, "forms", "--json")
			if err != nil {
				return err
			}
			var envelope browser.AgentBrowserEnvelope[[]any]
			if err := json.Unmarshal(out, &envelope); err != nil {
				return fmt.Errorf("parse forms: %w", err)
			}
			if !envelope.Success {
				return fmt.Errorf("forms: %s", envelope.Error)
			}
			if len(envelope.Data) == 0 {
				fmt.Fprintln(os.Stderr, "No fillable form elements found.")
				return nil
			}
			return printResult(cmd, envelope.Data)
		},
	}
}

func browserFillCmd() *cli.Command {
	return &cli.Command{
		Name:      "fill",
		Usage:     "Fill form fields in a tracked browser tab",
		ArgsUsage: "<selector|@eN> <value> [<selector|@eN> <value> ...]",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:  "submit",
			Usage: "CSS selector of element to click after filling (e.g. button[type=submit])",
		}),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) == 0 || len(args)%2 != 0 {
				return fmt.Errorf("arguments must be selector/value pairs (got %d args)", len(args))
			}

			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}

			execArgs := []string{"fill"}
			for i := 0; i < len(args); i += 2 {
				execArgs = append(execArgs, args[i], args[i+1])
			}
			if submit := cmd.String("submit"); submit != "" {
				execArgs = append(execArgs, "--submit", submit)
			}

			_, _, err = ab.Exec(ctx, execArgs...)
			if err != nil {
				return err
			}

			for i := 0; i < len(args); i += 2 {
				fmt.Fprintf(os.Stderr, "Filled %s\n", args[i])
			}
			if submit := cmd.String("submit"); submit != "" {
				fmt.Fprintf(os.Stderr, "Clicked %s\n", submit)
			}
			return nil
		},
	}
}

func browserKeypressCmd() *cli.Command {
	return &cli.Command{
		Name:      "keypress",
		Usage:     "Send keyboard events to the page",
		ArgsUsage: "<key>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			key := cmd.Args().First()
			if key == "" {
				return fmt.Errorf("key required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "keypress", key)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Sent key: %s\n", key)
			return nil
		},
	}
}

func browserDialogCmd() *cli.Command {
	return &cli.Command{
		Name:  "dialog",
		Usage: "Accept or dismiss a JavaScript dialog",
		Flags: append(browserActionFlags(false),
			&cli.BoolFlag{
				Name:  "accept",
				Usage: "Accept the dialog (default: true)",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "text",
				Usage: "Text to enter for prompt dialogs",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			args := []string{"dialog"}
			if cmd.Bool("accept") {
				args = append(args, "--accept")
			} else {
				args = append(args, "--dismiss")
			}
			if text := cmd.String("text"); text != "" {
				args = append(args, "--text", text)
			}
			_, _, err = ab.Exec(ctx, args...)
			if err != nil {
				return err
			}
			if cmd.Bool("accept") {
				fmt.Fprintln(os.Stderr, "Dialog accepted")
			} else {
				fmt.Fprintln(os.Stderr, "Dialog dismissed")
			}
			return nil
		},
	}
}

func browserCookiesCmd() *cli.Command {
	return &cli.Command{
		Name:  "cookies",
		Usage: "Manage browser cookies (includes httpOnly)",
		Commands: []*cli.Command{
			browserCookiesGetCmd(),
			browserCookiesSetCmd(),
			browserCookiesClearCmd(),
		},
	}
}

func browserCookiesGetCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "List all cookies for the current page",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			out, _, err := ab.Exec(ctx, "cookies", "--json")
			if err != nil {
				return err
			}
			var envelope browser.AgentBrowserEnvelope[[]any]
			if err := json.Unmarshal(out, &envelope); err != nil {
				return fmt.Errorf("parse cookies: %w", err)
			}
			if !envelope.Success {
				return fmt.Errorf("cookies: %s", envelope.Error)
			}
			if len(envelope.Data) == 0 {
				fmt.Fprintln(os.Stderr, "No cookies found.")
				return nil
			}
			return printResult(cmd, envelope.Data)
		},
	}
}

func browserCookiesSetCmd() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set a cookie",
		ArgsUsage: "<name> <value>",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:  "domain",
				Usage: "Cookie domain",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Cookie path",
				Value: "/",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser cookies set <name> <value>")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			execArgs := []string{"cookies", "set", args[0], args[1]}
			if domain := cmd.String("domain"); domain != "" {
				execArgs = append(execArgs, "--domain", domain)
			}
			if path := cmd.String("path"); path != "" {
				execArgs = append(execArgs, "--path", path)
			}
			_, _, err = ab.Exec(ctx, execArgs...)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Cookie %q set\n", args[0])
			return nil
		},
	}
}

func browserCookiesClearCmd() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Delete all cookies for the current page",
		Flags: browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "cookies", "clear")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Cookies cleared")
			return nil
		},
	}
}

func browserClickCmd() *cli.Command {
	return &cli.Command{
		Name:      "click",
		Usage:     "Click an element by CSS selector or snapshot ref",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "click", sel)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Clicked %s\n", sel)
			return nil
		},
	}
}

func browserTypeCmd() *cli.Command {
	return &cli.Command{
		Name:      "type",
		Usage:     "Type text into an element with real key events",
		ArgsUsage: "<selector|@eN> <text>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser type <selector> <text>")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "type", args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Typed into %s\n", args[0])
			return nil
		},
	}
}

func browserHoverCmd() *cli.Command {
	return &cli.Command{
		Name:      "hover",
		Usage:     "Move mouse to an element to trigger hover state",
		ArgsUsage: "<selector>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "hover", sel)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Hovered %s\n", sel)
			return nil
		},
	}
}

func browserScrollCmd() *cli.Command {
	return &cli.Command{
		Name:      "scroll",
		Usage:     "Scroll to an element or position",
		ArgsUsage: "[selector]",
		Flags: append(browserActionFlags(false),
			&cli.Float64Flag{
				Name:  "x",
				Usage: "Absolute X pixel position (when no selector given)",
			},
			&cli.Float64Flag{
				Name:  "y",
				Usage: "Absolute Y pixel position (when no selector given)",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			x := cmd.Float64("x")
			y := cmd.Float64("y")
			if sel == "" && x == 0 && y == 0 {
				return fmt.Errorf("provide a CSS selector or --x/--y position")
			}
			args := []string{"scroll"}
			if sel != "" {
				args = append(args, sel)
			} else {
				args = append(args, "--x", strconv.FormatFloat(x, 'f', -1, 64), "--y", strconv.FormatFloat(y, 'f', -1, 64))
			}
			_, _, err = ab.Exec(ctx, args...)
			if err != nil {
				return err
			}
			if sel != "" {
				fmt.Fprintf(os.Stderr, "Scrolled to %s\n", sel)
			} else {
				fmt.Fprintf(os.Stderr, "Scrolled to (%s, %s)\n", strconv.FormatFloat(x, 'f', -1, 64), strconv.FormatFloat(y, 'f', -1, 64))
			}
			return nil
		},
	}
}

func browserSelectCmd() *cli.Command {
	return &cli.Command{
		Name:      "select",
		Usage:     "Select an option in a <select> element",
		ArgsUsage: "<selector|@eN> <value>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser select <selector|@eN> <value>")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "select", args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Selected %q in %s\n", args[1], args[0])
			return nil
		},
	}
}

func browserWaitCmd() *cli.Command {
	return &cli.Command{
		Name:      "wait",
		Usage:     "Wait for an element to become visible",
		ArgsUsage: "<selector>",
		Flags: append(browserActionFlags(false),
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "Max time to wait",
				Value: 30 * time.Second,
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			waitCtx, cancel := context.WithTimeout(ctx, cmd.Duration("timeout"))
			defer cancel()
			_, _, err = ab.Exec(waitCtx, "wait", sel)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Element %s is visible\n", sel)
			return nil
		},
	}
}

func browserBackCmd() *cli.Command {
	return &cli.Command{
		Name:  "back",
		Usage: "Navigate back in history",
		Flags: browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "back")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Navigated back")
			return nil
		},
	}
}

func browserForwardCmd() *cli.Command {
	return &cli.Command{
		Name:  "forward",
		Usage: "Navigate forward in history",
		Flags: browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "forward")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Navigated forward")
			return nil
		},
	}
}

func browserReloadCmd() *cli.Command {
	return &cli.Command{
		Name:  "reload",
		Usage: "Reload the current page",
		Flags: browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			_, _, err = ab.Exec(ctx, "reload")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Reloaded")
			return nil
		},
	}
}
