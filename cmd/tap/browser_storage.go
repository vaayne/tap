package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

// ---------------------------------------------------------------------------
// storage subcommand group
// ---------------------------------------------------------------------------

func browserStorageCmd() *cli.Command {
	return &cli.Command{
		Name:  "storage",
		Usage: "Read and write browser web storage (localStorage / sessionStorage)",
		Description: `Inspect and manipulate localStorage and sessionStorage of the current tab.

Sub-commands:
  local                  Dump all localStorage entries as JSON
  local <key>            Print a single localStorage value
  local set <k> <v>      Set a localStorage entry
  local clear            Clear all localStorage entries
  session [...]          Same four forms for sessionStorage`,
		Commands: []*cli.Command{
			browserStorageTypeCmd("localStorage"),
			browserStorageTypeCmd("sessionStorage"),
		},
	}
}

// browserStorageTypeCmd builds the "local" or "session" sub-command group.
func browserStorageTypeCmd(storeName string) *cli.Command {
	// "localStorage" → "local", "sessionStorage" → "session"
	cmdName := "local"
	if storeName == "sessionStorage" {
		cmdName = "session"
	}
	return &cli.Command{
		Name:      cmdName,
		Usage:     fmt.Sprintf("Manage %s for the current tab", storeName),
		ArgsUsage: "[key]",
		Flags: append(browserActionFlags(false), &cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format: json, pretty (default), raw",
			Value:   formatPretty,
		}),
		Description: fmt.Sprintf(`Read or write %s entries for the current page.

Usage:
  tap browser %s                        Dump all entries as JSON
  tap browser %s <key>                  Print one value
  tap browser %s set <key> <value>      Set an entry
  tap browser %s clear                  Clear all entries`, storeName, cmdName, cmdName, cmdName, cmdName),
		Commands: []*cli.Command{
			{
				Name:        "set",
				Usage:       fmt.Sprintf("Set a %s entry", storeName),
				ArgsUsage:   "<key> <value>",
				Flags:       browserActionFlags(false),
				Description: fmt.Sprintf("Set a key-value pair in %s for the current page's origin.", storeName),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					configureLogging(cmd)
					args := cmd.Args().Slice()
					if len(args) < 2 {
						return fmt.Errorf("usage: tap browser %s set <key> <value>", cmdName)
					}
					mgr, err := newBrowserManager(cmd)
					if err != nil {
						return err
					}
					if err := mgr.SetStorageKey(ctx, cmd.String("session"), cmd.String("tab"), storeName, args[0], args[1]); err != nil {
						return err
					}
					fmt.Fprintf(os.Stderr, "%s[%q] set\n", storeName, args[0])
					return nil
				},
			},
			{
				Name:        "clear",
				Usage:       fmt.Sprintf("Clear all %s entries", storeName),
				Flags:       browserActionFlags(false),
				Description: fmt.Sprintf("Remove all entries from %s for the current page's origin.", storeName),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					configureLogging(cmd)
					mgr, err := newBrowserManager(cmd)
					if err != nil {
						return err
					}
					if err := mgr.ClearStorage(ctx, cmd.String("session"), cmd.String("tab"), storeName); err != nil {
						return err
					}
					fmt.Fprintf(os.Stderr, "%s cleared\n", storeName)
					return nil
				},
			},
		},
		// The default action handles: no args → dump all, one arg → get key.
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.String("tab")

			key := cmd.Args().First()
			if key == "" {
				// Dump all entries.
				entries, err := mgr.GetStorageAll(ctx, sessionName, tabName, storeName)
				if err != nil {
					return err
				}
				if len(entries) == 0 {
					fmt.Fprintf(os.Stderr, "%s is empty.\n", storeName)
					return nil
				}
				return printResult(cmd, entries)
			}

			// Single key lookup.
			val, err := mgr.GetStorageKey(ctx, sessionName, tabName, storeName, key)
			if err != nil {
				return err
			}
			if val == "" {
				fmt.Fprintf(os.Stderr, "%s[%q] is not set\n", storeName, key)
				return nil
			}
			fmt.Println(val)
			return nil
		},
	}
}

// ---------------------------------------------------------------------------
// state subcommand group
// ---------------------------------------------------------------------------

func browserStateCmd() *cli.Command {
	return &cli.Command{
		Name:  "state",
		Usage: "Save and load browser auth state (cookies + localStorage)",
		Description: `Export or import browser auth state in Playwright storageState format.

The state file is JSON:
  {
    "cookies": [{name, value, domain, path, expires, httpOnly, secure, sameSite}],
    "origins": [{"origin": "https://example.com", "localStorage": [{"name", "value"}]}]
  }

Limitations:
  save: captures localStorage only for the CURRENT tab's origin.
        All cookies from the entire browser context are included.
  load: cookies are applied globally; localStorage is only restored for
        origins matching the current page — other origins are skipped with a warning.`,
		Commands: []*cli.Command{
			browserStateSaveCmd(),
			browserStateLoadCmd(),
		},
	}
}

func browserStateSaveCmd() *cli.Command {
	return &cli.Command{
		Name:      "save",
		Usage:     "Export cookies and localStorage of the current origin to a file",
		ArgsUsage: "<path>",
		Flags:     browserActionFlags(false),
		Description: `Save all browser cookies and the localStorage of the current page's origin
to a JSON file in Playwright storageState format.

The file is written with 0600 permissions (it contains auth cookies).

Example:
  tap browser state save auth.json`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			path := cmd.Args().First()
			if path == "" {
				return fmt.Errorf("path required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			state, err := mgr.SaveState(ctx, cmd.String("session"), cmd.String("tab"))
			if err != nil {
				return err
			}
			data, err := json.MarshalIndent(state, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal state: %w", err)
			}
			if err := os.WriteFile(path, data, 0o600); err != nil {
				return fmt.Errorf("write state: %w", err)
			}
			cookieCount := len(state.Cookies)
			lsCount := 0
			for _, o := range state.Origins {
				lsCount += len(o.LocalStorage)
			}
			fmt.Fprintf(os.Stderr, "Saved %d cookies and %d localStorage entries to %s\n", cookieCount, lsCount, path)
			return nil
		},
	}
}

func browserStateLoadCmd() *cli.Command {
	return &cli.Command{
		Name:      "load",
		Usage:     "Import cookies and localStorage from a saved state file",
		ArgsUsage: "<path>",
		Flags:     browserActionFlags(false),
		Description: `Load browser auth state from a JSON file (Playwright storageState format).

All cookies are applied globally. localStorage is restored only for the origin
matching the current page — mismatched origins are skipped with a warning.

Example:
  tap browser state load auth.json`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			path := cmd.Args().First()
			if path == "" {
				return fmt.Errorf("path required")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read state: %w", err)
			}
			var state browser.StorageState
			if err := json.Unmarshal(data, &state); err != nil {
				return fmt.Errorf("parse state: %w", err)
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			warnFn := func(msg string) {
				fmt.Fprintf(os.Stderr, "warning: %s\n", msg)
			}
			if err := mgr.LoadState(ctx, cmd.String("session"), cmd.String("tab"), &state, warnFn); err != nil {
				return err
			}
			cookieCount := len(state.Cookies)
			lsCount := 0
			for _, o := range state.Origins {
				lsCount += len(o.LocalStorage)
			}
			fmt.Fprintf(os.Stderr, "Loaded %d cookies and %d localStorage entries from %s\n", cookieCount, lsCount, path)
			return nil
		},
	}
}
