package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp/kb"
	"github.com/urfave/cli/v3"
	defuddle "github.com/vaayne/go-defuddle"
	"github.com/vaayne/tap/browser"
)

// resolveKeyName maps human-readable key names to chromedp/kb constants.
func resolveKeyName(name string) string {
	// Handle modifier combinations like Ctrl+a
	if strings.Contains(name, "+") {
		parts := strings.SplitN(name, "+", 2)
		modifier := resolveKeyName(parts[0])
		key := resolveKeyName(parts[1])
		return modifier + key + modifier
	}

	keyMap := map[string]string{
		"enter":      kb.Enter,
		"return":     kb.Enter,
		"tab":        kb.Tab,
		"escape":     kb.Escape,
		"esc":        kb.Escape,
		"backspace":  kb.Backspace,
		"delete":     kb.Delete,
		"space":      " ",
		"arrowup":    kb.ArrowUp,
		"arrowdown":  kb.ArrowDown,
		"arrowleft":  kb.ArrowLeft,
		"arrowright": kb.ArrowRight,
		"home":       kb.Home,
		"end":        kb.End,
		"pageup":     kb.PageUp,
		"pagedown":   kb.PageDown,
		"ctrl":       kb.Control,
		"control":    kb.Control,
		"alt":        kb.Alt,
		"shift":      kb.Shift,
		"meta":       kb.Meta,
	}
	if v, ok := keyMap[strings.ToLower(name)]; ok {
		return v
	}
	return name
}

func browserNavigateCmd() *cli.Command {
	return &cli.Command{
		Name:      "navigate",
		Usage:     "Navigate a tracked browser tab to a URL",
		ArgsUsage: "<url>",
		Flags:     browserActionFlags(false),
		Description: `Navigate the resolved tracked tab to the given URL.

Examples:
  tap browser navigate https://example.com
  tap browser navigate --session dev --tab main https://example.com`,
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

Output formatting follows the same pretty/json/raw conventions used by 'tap site'.

Examples:
  tap browser evaluate "document.title"
  tap browser evaluate "document.querySelectorAll('a').length"
  tap browser evaluate "JSON.stringify(performance.timing)" -f json`,
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
session name, tab name, and current timestamp.

Examples:
  tap browser screenshot
  tap browser screenshot --output page.png
  tap browser screenshot --session dev --tab main`,
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
		Description: `Extract clean, readable content from the current page using defuddle.
Strips navigation, ads, scripts, and boilerplate — returns only the main content
as Markdown (default) or JSON with metadata.

Optionally scope to a CSS selector to extract from a specific section.

This is the most token-efficient way to read page content — far cheaper
than evaluating outerHTML.

Examples:
  tap browser text                       Extract main content as Markdown
  tap browser text "#article"            Extract from a specific section
  tap browser text -f json               Include metadata (title, word count)`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			rt, err := mgr.ResolveTarget(ctx, cmd.String("session"), cmd.String("tab"))
			if err != nil {
				return err
			}
			html, pageURL, err := browser.GetHTMLTarget(ctx, rt.DebugURL, rt.TargetID, sel)
			if err != nil {
				return err
			}
			if html == "" {
				fmt.Fprintln(os.Stderr, "No content found.")
				return nil
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
			// Default: print markdown directly
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
		Description: `Capture a compact semantic tree from the current page and assign stable refs
for interactive elements (e.g. @e1, @e2). Refs can be reused in click/type/fill/select.

Examples:
  tap browser snapshot
  tap browser snapshot --interactive -f json
  tap browser click @e3
  tap browser type @e1 "hello"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.Snapshot(ctx, cmd.String("session"), cmd.String("tab"), browser.SnapshotOptions{
				InteractiveOnly: cmd.Bool("interactive"),
				Selector:        cmd.String("selector"),
				Depth:           int64(cmd.Int("depth")),
				Mode:            "auto",
			})
			if err != nil {
				return err
			}
			if len(result.Nodes) == 0 {
				fmt.Fprintln(os.Stderr, "No snapshot nodes captured.")
				return nil
			}
			return printResult(cmd, result)
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
tab name, and timestamp.

Examples:
  tap browser pdf
  tap browser pdf --output report.pdf --landscape
  tap browser pdf --scale 0.8 --background=false`,
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
		ArgsUsage: "<selector|@eN> <value> [<selector|@eN> <value> ...]",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:  "submit",
			Usage: "CSS selector of element to click after filling (e.g. button[type=submit])",
		}),
		Description: `Fill one or more form fields by CSS selector, then optionally submit.

Arguments are target/value pairs. Targets can be CSS selectors or snapshot refs.
Values are set using React-compatible
native setters with proper input/change event dispatch, so this works with
React, Vue, Angular, and vanilla HTML forms.

Examples:
  tap browser fill "#username" "myuser"
  tap browser fill @e1 "me@example.com"
  tap browser fill "#email" "me@example.com" "#password" "secret" --submit "button[type=submit]"
  tap browser fill "input[type=search]" "query text"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) == 0 || len(args)%2 != 0 {
				return fmt.Errorf("arguments must be selector/value pairs (got %d args)", len(args))
			}

			var inputs []browser.FillInput
			for i := 0; i < len(args); i += 2 {
				inputs = append(inputs, browser.FillInput{
					Target: args[i],
					Value:  args[i+1],
				})
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")
			submitSelector := cmd.String("submit")

			if err := mgr.FillElements(ctx, sessionName, tabName, inputs, submitSelector); err != nil {
				return err
			}

			for _, f := range inputs {
				fmt.Fprintf(os.Stderr, "Filled %s\n", f.Target)
			}
			if submitSelector != "" {
				fmt.Fprintf(os.Stderr, "Clicked %s\n", submitSelector)
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
		Description: `Send key events to the page. Common keys:
  Enter, Tab, Escape, Backspace, Delete, ArrowUp, ArrowDown, ArrowLeft, ArrowRight,
  Space, Home, End, PageUp, PageDown, F1-F12

For modifier combinations, separate with +: Ctrl+a, Ctrl+c, Ctrl+v, Shift+Tab

Regular text is sent as individual keystrokes.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			key := cmd.Args().First()
			if key == "" {
				return fmt.Errorf("key required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			keys := resolveKeyName(key)
			if err := mgr.Keypress(ctx, cmd.String("session"), cmd.String("tab"), keys); err != nil {
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
		Description: `Handle a pending JavaScript dialog (alert, confirm, prompt, onbeforeunload).

Unhandled dialogs block all CDP commands. Use this to dismiss them.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			accept := cmd.Bool("accept")
			text := cmd.String("text")
			if err := mgr.Dialog(ctx, cmd.String("session"), cmd.String("tab"), accept, text); err != nil {
				return err
			}
			if accept {
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
		Description: `Get, set, or clear cookies for the current page.

Unlike document.cookie, this uses CDP and can access httpOnly cookies.`,
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			cookies, err := mgr.GetCookies(ctx, cmd.String("session"), cmd.String("tab"))
			if err != nil {
				return err
			}
			if len(cookies) == 0 {
				fmt.Fprintln(os.Stderr, "No cookies found.")
				return nil
			}
			return printResult(cmd, cookies)
		},
	}
}

func browserCookiesSetCmd() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set a cookie",
		ArgsUsage: "<name> <value>",
		Description: `Set a browser cookie via CDP (can set httpOnly cookies).

Examples:
  tap browser cookies set "session_id" "abc123"
  tap browser cookies set "token" "xyz" --domain .example.com --path /api`,
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.SetCookie(ctx, cmd.String("session"), cmd.String("tab"),
				args[0], args[1], cmd.String("domain"), cmd.String("path")); err != nil {
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
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.ClearCookies(ctx, cmd.String("session"), cmd.String("tab")); err != nil {
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
		Description: `Dispatch a real mouse click (mouseMoved → mousePressed → mouseReleased)
on the first visible element matching the CSS selector.

Unlike JavaScript .click(), this triggers hover states and works with
sites that listen on mousedown or have hover-triggered menus.

Examples:
  tap browser click "button.submit"
  tap browser click @e3
  tap browser click "a[href='/login']"`,
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
			if err := mgr.ClickElement(ctx, cmd.String("session"), cmd.String("tab"), sel); err != nil {
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
		Description: `Focus the element matching the CSS selector and send individual
keyDown/keyUp events for each character — behaving like a real user typing.

Use this instead of 'fill' when the site validates per-keystroke input
or has anti-bot detection.

Examples:
  tap browser type "#search" "golang tutorials"
  tap browser type @e1 "hello world"
  tap browser type "input[name=q]" "hello world"`,
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
			if err := mgr.TypeElement(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1]); err != nil {
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
CSS :hover states and mouseenter/mouseover listeners.

Examples:
  tap browser hover "nav.menu > li:first-child"
  tap browser hover ".dropdown-trigger"`,
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

Use this to trigger lazy-loaded content or scroll-based UI updates.

Examples:
  tap browser scroll "#comments"
  tap browser scroll --y 1000
  tap browser scroll --x 0 --y 99999`,
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
		ArgsUsage: "<selector|@eN> <value>",
		Flags:     browserActionFlags(false),
		Description: `Select an option by value in a <select> element, dispatching
focus, input, and change events. Works with native HTML selects.

For custom dropdown components (React Select, etc.), use 'click' instead.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser select <selector|@eN> <value>")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.SelectElement(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1]); err != nil {
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
		Description: `Navigate the resolved tracked tab back one entry in browser history.

Examples:
  tap browser back
  tap browser back --session dev --tab main`,
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
		Description: `Navigate the resolved tracked tab forward one entry in browser history.

Examples:
  tap browser forward
  tap browser forward --session dev --tab main`,
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
		Description: `Reload the current page of the resolved tracked tab.

Examples:
  tap browser reload
  tap browser reload --session dev --tab main`,
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
