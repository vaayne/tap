package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// browserGetCmd returns the `tap browser get` parent command with subcommands
// for querying page/element properties.
func browserGetCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Query element and page properties",
		Description: `Read-only element and page queries. Accepts CSS selectors or snapshot refs (@eN).

Examples:
  tap browser get text "h1"
  tap browser get html "#content"
  tap browser get value "input[name=q]"
  tap browser get attr "a.logo" href
  tap browser get title
  tap browser get url
  tap browser get count "li.item"
  tap browser get box "#sidebar"
  tap browser get styles "button.primary"`,
		Commands: []*cli.Command{
			browserGetTextCmd(),
			browserGetHTMLCmd(),
			browserGetValueCmd(),
			browserGetAttrCmd(),
			browserGetTitleCmd(),
			browserGetURLCmd(),
			browserGetCountCmd(),
			browserGetBoxCmd(),
			browserGetStylesCmd(),
		},
	}
}

// browserIsCmd returns the `tap browser is` parent command with subcommands
// for boolean element state checks.
func browserIsCmd() *cli.Command {
	return &cli.Command{
		Name:  "is",
		Usage: "Check element boolean state",
		Description: `Boolean element state checks. Prints "true" or "false" and exits 0.
Accepts CSS selectors or snapshot refs (@eN).

Examples:
  tap browser is visible "#modal"
  tap browser is enabled "button[type=submit]"
  tap browser is checked "input[type=checkbox]"
  tap browser is visible @e3`,
		Commands: []*cli.Command{
			browserIsVisibleCmd(),
			browserIsEnabledCmd(),
			browserIsCheckedCmd(),
		},
	}
}

// ---------------------------------------------------------------------------
// get subcommands
// ---------------------------------------------------------------------------

func browserGetTextCmd() *cli.Command {
	return &cli.Command{
		Name:      "text",
		Usage:     "Get the raw textContent of an element",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print the raw textContent of the first element matching the selector.

Unlike 'tap browser text' (which runs defuddle for clean markdown), this
returns the unprocessed DOM text — useful for reading specific values.

Examples:
  tap browser get text "h1"
  tap browser get text @e2`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryText(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetHTMLCmd() *cli.Command {
	return &cli.Command{
		Name:      "html",
		Usage:     "Get the innerHTML of an element",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print the innerHTML of the first element matching the selector.

Examples:
  tap browser get html "#article"
  tap browser get html @e1`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryHTML(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetValueCmd() *cli.Command {
	return &cli.Command{
		Name:      "value",
		Usage:     "Get the value property of an input element",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print the .value property of the first element matching the selector.
Works on input, textarea, and select elements.

Examples:
  tap browser get value "input[name=email]"
  tap browser get value @e4`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryValue(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetAttrCmd() *cli.Command {
	return &cli.Command{
		Name:      "attr",
		Usage:     "Get an attribute value from an element",
		ArgsUsage: "<selector|@eN> <attr>",
		Flags:     browserActionFlags(false),
		Description: `Print the value of the named attribute on the first element matching the selector.
Returns an empty string if the attribute is absent.

Examples:
  tap browser get attr "a.logo" href
  tap browser get attr "img#banner" src
  tap browser get attr @e1 aria-label`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser get attr <selector|@eN> <attr>")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryAttr(ctx, cmd.String("session"), cmd.String("tab"), args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetTitleCmd() *cli.Command {
	return &cli.Command{
		Name:  "title",
		Usage: "Get the current page title",
		Flags: browserActionFlags(false),
		Description: `Print the document.title of the current page.

Examples:
  tap browser get title`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryTitle(ctx, cmd.String("session"), cmd.String("tab"))
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetURLCmd() *cli.Command {
	return &cli.Command{
		Name:  "url",
		Usage: "Get the current page URL",
		Flags: browserActionFlags(false),
		Description: `Print the current location.href of the tracked tab.

Examples:
  tap browser get url`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryURL(ctx, cmd.String("session"), cmd.String("tab"))
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetCountCmd() *cli.Command {
	return &cli.Command{
		Name:      "count",
		Usage:     "Count elements matching a CSS selector",
		ArgsUsage: "<selector>",
		Flags:     browserActionFlags(false),
		Description: `Print the number of elements matching the CSS selector.
Note: count does not support @eN refs since snapshot refs address single elements.

Examples:
  tap browser get count "li.result"
  tap browser get count "input[type=checkbox]"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryCount(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserGetBoxCmd() *cli.Command {
	return &cli.Command{
		Name:      "box",
		Usage:     "Get the bounding box of an element",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print the bounding box (x, y, width, height in pixels) of the first element
matching the selector, as JSON.

Examples:
  tap browser get box "#sidebar"
  tap browser get box @e2`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryBox(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		},
	}
}

func browserGetStylesCmd() *cli.Command {
	return &cli.Command{
		Name:      "styles",
		Usage:     "Get the computed CSS styles of an element",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print the computed style properties of the first element matching the selector,
as a JSON object mapping property names to values.

Examples:
  tap browser get styles "button.primary"
  tap browser get styles @e1`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryStyles(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		},
	}
}

// ---------------------------------------------------------------------------
// is subcommands
// ---------------------------------------------------------------------------

func browserIsVisibleCmd() *cli.Command {
	return &cli.Command{
		Name:      "visible",
		Usage:     "Check if an element is visible",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print "true" if the element exists, is not hidden (display/visibility/opacity),
and has non-zero dimensions. Prints "false" otherwise. Always exits 0.

Examples:
  tap browser is visible "#modal"
  tap browser is visible @e3`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryVisible(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserIsEnabledCmd() *cli.Command {
	return &cli.Command{
		Name:      "enabled",
		Usage:     "Check if a form element is enabled (not disabled)",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print "true" if the element exists and its .disabled property is false.
Prints "false" if the element is disabled or not found. Always exits 0.

Examples:
  tap browser is enabled "button[type=submit]"
  tap browser is enabled @e5`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryEnabled(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}

func browserIsCheckedCmd() *cli.Command {
	return &cli.Command{
		Name:      "checked",
		Usage:     "Check if a checkbox or radio input is checked",
		ArgsUsage: "<selector|@eN>",
		Flags:     browserActionFlags(false),
		Description: `Print "true" if the element exists and its .checked property is true.
Prints "false" otherwise. Always exits 0.

Examples:
  tap browser is checked "input[type=checkbox]#agree"
  tap browser is checked @e7`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			sel := cmd.Args().First()
			if sel == "" {
				return fmt.Errorf("selector required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			result, err := mgr.QueryChecked(ctx, cmd.String("session"), cmd.String("tab"), sel)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
}
