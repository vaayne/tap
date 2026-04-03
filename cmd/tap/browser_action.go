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

func browserNavigateCmd() *cli.Command {
	return &cli.Command{
		Name:      "navigate",
		Usage:     "Navigate a tracked browser tab",
		ArgsUsage: "<url>",
		Flags:     browserActionFlags(false),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			url := cmd.Args().First()
			if url == "" {
				return fmt.Errorf("URL required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			if err := mgr.Navigate(ctx, sessionName, tabName, url); err != nil {
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
		Description: `Evaluate JavaScript in the resolved tracked tab.

Output formatting follows the same pretty/json/raw conventions used by 'tap site'.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			js := cmd.Args().First()
			if js == "" {
				return fmt.Errorf("JavaScript expression required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			result, err := mgr.Evaluate(ctx, sessionName, tabName, js)
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
		Description: `Capture a screenshot from the resolved tracked tab.

When --output is omitted, tap will generate a deterministic file path from the
session name, tab name, and current timestamp.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			result, err := mgr.Screenshot(ctx, sessionName, tabName)
			if err != nil {
				return err
			}

			outPath := cmd.String("output")
			if outPath == "" {
				outPath = fmt.Sprintf("screenshot-%s-%s-%d.png", result.SessionName, result.TabName, time.Now().Unix())
			}

			if err := os.WriteFile(outPath, result.Data, 0o644); err != nil {
				return fmt.Errorf("write screenshot: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%s\n", outPath)
			return nil
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
		Description: `Save the current page of the resolved tracked tab as a PDF file.

When --output is omitted, tap generates a file name from the session,
tab name, and timestamp.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.PDF(ctx, cmd.String("session"), cmd.String("tab"),
				cmd.Bool("landscape"), cmd.Bool("background"), cmd.Float64("scale"))
			if err != nil {
				return err
			}
			outPath := cmd.String("output")
			if outPath == "" {
				outPath = fmt.Sprintf("page-%s-%s-%d.pdf", result.SessionName, result.TabName, time.Now().Unix())
			}
			if err := os.WriteFile(outPath, result.Data, 0o644); err != nil {
				return fmt.Errorf("write pdf: %w", err)
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
		Description: `List all fillable form elements (inputs, textareas, selects, buttons)
in the resolved tracked tab.

Each element includes its best CSS selector, type, name, placeholder, current
value, associated label, and role (text, toggle, select, submit). Use the
reported selectors with 'tap browser fill'.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			fields, err := mgr.Forms(ctx, sessionName, tabName)
			if err != nil {
				return err
			}
			if len(fields) == 0 {
				fmt.Fprintln(os.Stderr, "No fillable form elements found.")
				return nil
			}
			return printResult(cmd, fields)
		},
	}
}

func browserFillCmd() *cli.Command {
	return &cli.Command{
		Name:      "fill",
		Usage:     "Fill form fields in a tracked browser tab",
		ArgsUsage: "<selector> <value> [<selector> <value> ...]",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:  "submit",
			Usage: "CSS selector of element to click after filling (e.g. button[type=submit])",
		}),
		Description: `Fill one or more form fields by CSS selector, then optionally submit.

Arguments are selector/value pairs. Values are set using React-compatible
native setters with proper input/change event dispatch, so this works with
React, Vue, Angular, and vanilla HTML forms.

Examples:
  tap browser fill "#username" "myuser"
  tap browser fill "#email" "me@example.com" "#password" "secret" --submit "button[type=submit]"
  tap browser fill "input[type=search]" "query text"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) == 0 || len(args)%2 != 0 {
				return fmt.Errorf("arguments must be selector/value pairs (got %d args)", len(args))
			}

			var fields []browser.FillField
			for i := 0; i < len(args); i += 2 {
				fields = append(fields, browser.FillField{
					Selector: args[i],
					Value:    args[i+1],
				})
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			submitSelector := cmd.String("submit")

			if err := mgr.Fill(ctx, sessionName, tabName, fields, submitSelector); err != nil {
				return err
			}

			for _, f := range fields {
				fmt.Fprintf(os.Stderr, "Filled %s\n", f.Selector)
			}
			if submitSelector != "" {
				fmt.Fprintf(os.Stderr, "Clicked %s\n", submitSelector)
			}
			return nil
		},
	}
}

func browserClickCmd() *cli.Command {
	return &cli.Command{
		Name:      "click",
		Usage:     "Click an element by CSS selector",
		ArgsUsage: "<selector>",
		Flags:     browserActionFlags(false),
		Description: `Dispatch a real mouse click (mouseMoved → mousePressed → mouseReleased)
on the first visible element matching the CSS selector.

Unlike JavaScript .click(), this triggers hover states and works with
sites that listen on mousedown or have hover-triggered menus.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Click(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
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
		ArgsUsage: "<selector> <text>",
		Flags:     browserActionFlags(false),
		Description: `Focus the element matching the CSS selector and send individual
keyDown/keyUp events for each character — behaving like a real user typing.

Use this instead of 'fill' when the site validates per-keystroke input
or has anti-bot detection.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser type <selector> <text>")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Type(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1]); err != nil {
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
		Description: `Move the mouse to the center of the first visible element matching
the CSS selector. Dispatches real mouseMoved events that trigger
CSS :hover states and mouseenter/mouseover listeners.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Hover(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
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
		Description: `Scroll the element matching the CSS selector into view. If no selector
is provided, scroll to the absolute pixel position given by --x and --y.

Use this to trigger lazy-loaded content or scroll-based UI updates.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			x := cmd.Float64("x")
			y := cmd.Float64("y")
			if sel == "" && x == 0 && y == 0 {
				return fmt.Errorf("provide a CSS selector or --x/--y position")
			}
			if err := mgr.Scroll(ctx, cmd.String("session"), cmd.String("tab"), sel, x, y); err != nil {
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
		ArgsUsage: "<selector> <value>",
		Flags:     browserActionFlags(false),
		Description: `Select an option by value in a <select> element, dispatching
focus, input, and change events. Works with native HTML selects.

For custom dropdown components (React Select, etc.), use 'click' instead.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser select <selector> <value>")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Select(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1]); err != nil {
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
		Description: `Wait until the first element matching the CSS selector becomes visible.
Uses CDP's built-in visibility polling — more reliable than JS-based polling.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("CSS selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			timeout := cmd.Duration("timeout")
			if err := mgr.WaitFor(ctx, cmd.String("session"), cmd.String("tab"), sel, timeout); err != nil {
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Back(ctx, cmd.String("session"), cmd.String("tab")); err != nil {
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Forward(ctx, cmd.String("session"), cmd.String("tab")); err != nil {
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.Reload(ctx, cmd.String("session"), cmd.String("tab")); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Reloaded")
			return nil
		},
	}
}
