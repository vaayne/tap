package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

// browserFindCmd returns the 'tap browser find' parent command with one
// subcommand per locator kind.
func browserFindCmd() *cli.Command {
	return &cli.Command{
		Name:  "find",
		Usage: "Locate elements by semantic attribute and perform an action",
		Description: `Find elements using semantic locators (role, text, label, …) and
interact with them without needing a CSS selector.

Supported actions:
  click    — real mouse click
  fill     — set value via React-compatible native setter (requires <value>)
  type     — simulate key-by-key typing (requires <value>)
  hover    — move mouse to element centre
  focus    — keyboard-focus the element
  check    — check a checkbox or radio
  uncheck  — uncheck a checkbox or radio
  text     — print the element's trimmed textContent

Multiple matches: the first matching element is used.

Examples:
  tap browser find role button click --name "Submit"
  tap browser find text "Sign in" click
  tap browser find label "Email" fill "me@example.com"
  tap browser find placeholder "Search…" type "golang"
  tap browser find testid "login-btn" click
  tap browser find first "li.item" text
  tap browser find nth 2 "li.item" click
  tap browser find last "tr" text`,
		Commands: []*cli.Command{
			findRoleCmd(),
			findTextCmd(),
			findLabelCmd(),
			findPlaceholderCmd(),
			findAltCmd(),
			findTitleCmd(),
			findTestIDCmd(),
			findFirstCmd(),
			findLastCmd(),
			findNthCmd(),
		},
	}
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// runFind resolves the manager/target and calls Manager.Find, printing output.
func runFind(ctx context.Context, cmd *cli.Command, loc browser.FindLocator, action browser.FindAction, value string) error {
	mgr, err := newBrowserManager(cmd)
	if err != nil {
		return err
	}
	result, err := mgr.Find(ctx, cmd.String("session"), cmd.String("tab"), loc, action, value)
	if err != nil {
		return err
	}
	if result != "" {
		fmt.Println(result)
	}
	return nil
}

// parseAction validates and returns the FindAction for the given string.
func parseAction(s string) (browser.FindAction, error) {
	switch browser.FindAction(s) {
	case browser.FindActionClick, browser.FindActionFill, browser.FindActionType,
		browser.FindActionHover, browser.FindActionFocus,
		browser.FindActionCheck, browser.FindActionUncheck, browser.FindActionText:
		return browser.FindAction(s), nil
	default:
		return "", fmt.Errorf("unknown action %q (use: click fill type hover focus check uncheck text)", s)
	}
}

// requireValue checks that fill/type actions have a value argument.
func requireValue(action browser.FindAction, value string) error {
	if (action == browser.FindActionFill || action == browser.FindActionType) && value == "" {
		return fmt.Errorf("action %q requires a <value> argument", action)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Subcommands
// ---------------------------------------------------------------------------

func findRoleCmd() *cli.Command {
	return &cli.Command{
		Name:      "role",
		Usage:     "Find by ARIA role",
		ArgsUsage: "<role> <action> [value]",
		Flags: append(browserActionFlags(false),
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Filter by accessible name (substring, case-insensitive)",
			},
		),
		Description: `Find the first element with the given ARIA role and perform action.

Use --name to narrow the match to a specific accessible name.

Examples:
  tap browser find role button click --name "Submit"
  tap browser find role textbox fill "hello@example.com" --name "Email"
  tap browser find role link click --name "Sign in"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find role <role> <action> [value]")
			}
			role, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorRole, Role: role, Name: cmd.String("name")}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find role=%q action=%s ok\n", role, action)
			}
			return nil
		},
	}
}

func findTextCmd() *cli.Command {
	return &cli.Command{
		Name:      "text",
		Usage:     "Find by visible text content",
		ArgsUsage: "<text> <action> [value]",
		Flags: append(browserActionFlags(false),
			&cli.BoolFlag{
				Name:  "exact",
				Usage: "Require an exact text match instead of substring",
			},
		),
		Description: `Find the first element whose textContent contains (or exactly matches)
the given text, then perform action.

Examples:
  tap browser find text "Sign in" click
  tap browser find text "Submit" click --exact
  tap browser find text "Welcome" text`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find text <text> <action> [value]")
			}
			text, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorText, Text: text, Exact: cmd.Bool("exact")}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find text=%q action=%s ok\n", text, action)
			}
			return nil
		},
	}
}

func findLabelCmd() *cli.Command {
	return &cli.Command{
		Name:      "label",
		Usage:     "Find a form element by its label text",
		ArgsUsage: "<label> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the first input/textarea/select/button associated with a label whose
text contains the given string, then perform action.

Resolution order: label[for=id] → wrapping <label> → aria-label.

Examples:
  tap browser find label "Email" fill "me@example.com"
  tap browser find label "Password" fill "secret"
  tap browser find label "Remember me" check`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find label <label> <action> [value]")
			}
			query, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorLabel, Query: query}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find label=%q action=%s ok\n", query, action)
			}
			return nil
		},
	}
}

func findPlaceholderCmd() *cli.Command {
	return &cli.Command{
		Name:      "placeholder",
		Usage:     "Find an input by its placeholder attribute",
		ArgsUsage: "<placeholder> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the first element whose placeholder attribute contains the given
string (case-insensitive), then perform action.

Examples:
  tap browser find placeholder "Search…" type "golang"
  tap browser find placeholder "Email address" fill "me@example.com"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find placeholder <placeholder> <action> [value]")
			}
			query, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorPlaceholder, Query: query}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find placeholder=%q action=%s ok\n", query, action)
			}
			return nil
		},
	}
}

func findAltCmd() *cli.Command {
	return &cli.Command{
		Name:      "alt",
		Usage:     "Find an element by its alt attribute",
		ArgsUsage: "<text> <action>",
		Flags:     browserActionFlags(false),
		Description: `Find the first element whose alt attribute contains the given string
(case-insensitive), then perform action.

Examples:
  tap browser find alt "company logo" click
  tap browser find alt "product image" hover`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find alt <text> <action>")
			}
			query, actionStr := args[0], args[1]
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorAlt, Query: query}
			if err := runFind(ctx, cmd, loc, action, ""); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find alt=%q action=%s ok\n", query, action)
			}
			return nil
		},
	}
}

func findTitleCmd() *cli.Command {
	return &cli.Command{
		Name:      "title",
		Usage:     "Find an element by its title attribute",
		ArgsUsage: "<text> <action>",
		Flags:     browserActionFlags(false),
		Description: `Find the first element whose title attribute contains the given string
(case-insensitive), then perform action.

Examples:
  tap browser find title "Close dialog" click
  tap browser find title "More options" hover`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find title <text> <action>")
			}
			query, actionStr := args[0], args[1]
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorTitle, Query: query}
			if err := runFind(ctx, cmd, loc, action, ""); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find title=%q action=%s ok\n", query, action)
			}
			return nil
		},
	}
}

func findTestIDCmd() *cli.Command {
	return &cli.Command{
		Name:      "testid",
		Usage:     "Find an element by its data-testid attribute",
		ArgsUsage: "<id> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the first element with a data-testid attribute containing the given
string (case-insensitive), then perform action.

Examples:
  tap browser find testid "login-btn" click
  tap browser find testid "email-input" fill "me@example.com"
  tap browser find testid "username" type "alice"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find testid <id> <action> [value]")
			}
			query, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorTestID, Query: query}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find testid=%q action=%s ok\n", query, action)
			}
			return nil
		},
	}
}

func findFirstCmd() *cli.Command {
	return &cli.Command{
		Name:      "first",
		Usage:     "Act on the first element matching a CSS selector",
		ArgsUsage: "<css-selector> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the first element matching the CSS selector, then perform action.
Equivalent to nth 0 <selector>.

Examples:
  tap browser find first "li.item" text
  tap browser find first "button[type=submit]" click`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find first <css-selector> <action> [value]")
			}
			sel, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorFirst, CSSSelector: sel}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find first=%q action=%s ok\n", sel, action)
			}
			return nil
		},
	}
}

func findLastCmd() *cli.Command {
	return &cli.Command{
		Name:      "last",
		Usage:     "Act on the last element matching a CSS selector",
		ArgsUsage: "<css-selector> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the last element matching the CSS selector, then perform action.

Examples:
  tap browser find last "tr" text
  tap browser find last "li.result" click`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser find last <css-selector> <action> [value]")
			}
			sel, actionStr := args[0], args[1]
			value := ""
			if len(args) > 2 {
				value = args[2]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			loc := browser.FindLocator{Kind: browser.LocatorLast, CSSSelector: sel}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find last=%q action=%s ok\n", sel, action)
			}
			return nil
		},
	}
}

func findNthCmd() *cli.Command {
	return &cli.Command{
		Name:      "nth",
		Usage:     "Act on the n-th (1-based) element matching a CSS selector",
		ArgsUsage: "<n> <css-selector> <action> [value]",
		Flags:     browserActionFlags(false),
		Description: `Find the n-th element matching the CSS selector (1-based), then perform
action. Use 1 for the first element, 2 for the second, etc.

Examples:
  tap browser find nth 2 "li.item" click
  tap browser find nth 3 "tr" text`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 3 {
				return fmt.Errorf("usage: tap browser find nth <n> <css-selector> <action> [value]")
			}
			n, err := strconv.Atoi(args[0])
			if err != nil || n < 1 {
				return fmt.Errorf("n must be a positive integer, got %q", args[0])
			}
			sel, actionStr := args[1], args[2]
			value := ""
			if len(args) > 3 {
				value = args[3]
			}
			action, err := parseAction(actionStr)
			if err != nil {
				return err
			}
			if err := requireValue(action, value); err != nil {
				return err
			}
			// Convert 1-based user input to 0-based internal index.
			loc := browser.FindLocator{Kind: browser.LocatorNth, CSSSelector: sel, Index: n - 1}
			if err := runFind(ctx, cmd, loc, action, value); err != nil {
				return err
			}
			if action != browser.FindActionText {
				fmt.Fprintf(os.Stderr, "find nth(%d)=%q action=%s ok\n", n, sel, action)
			}
			return nil
		},
	}
}
