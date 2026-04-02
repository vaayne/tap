package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func browserTabCmd() *cli.Command {
	return &cli.Command{
		Name:  "tab",
		Usage: "Manage tracked tabs within a persistent browser session",
		Description: `Tap manages only tracked tabs recorded in session metadata.

Closing or losing a tracked target updates its stored status. If the selected tab
is closed, tap selects the next remaining live tracked tab by creation order. If
no live tracked tabs remain, the session has no selected tab.`,
		Commands: []*cli.Command{
			browserTabNewCmd(),
			browserTabListCmd(),
			browserTabSelectCmd(),
			browserTabCloseCmd(),
		},
	}
}

func browserTabNewCmd() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a tracked browser tab",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "Initial URL for the new tab",
				Value: "about:blank",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			tabName := cmd.Args().First()
			if tabName == "" {
				return fmt.Errorf("tab name required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			url := cmd.String("url")
			if err := mgr.CreateTab(ctx, sessionName, tabName, url); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Tab %q created\n", tabName)
			return nil
		},
	}
}

func browserTabListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List tracked tabs for a session",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabs, err := mgr.ListTabs(ctx, sessionName)
			if err != nil {
				return err
			}
			if len(tabs) == 0 {
				fmt.Fprintln(os.Stderr, "No tracked tabs found.")
				return nil
			}

			// Determine which tab is selected.
			session, err := mgr.GetSession(ctx, sessionName)
			selectedTab := ""
			if err == nil && session != nil {
				selectedTab = session.SelectedTab
			}

			c := useColor(cmd)
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, bold(c, "NAME\tSTATUS\tURL\tSELECTED"))
			for _, tab := range tabs {
				status := string(tab.Status)
				switch tab.Status {
				case browser.TabStatusLive:
					status = green(c, status)
				case browser.TabStatusStale:
					status = yellow(c, status)
				}
				sel := ""
				if tab.Name == selectedTab {
					sel = green(c, "*")
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", tab.Name, status, tab.URL, sel)
			}
			return w.Flush()
		},
	}
}

func browserTabSelectCmd() *cli.Command {
	return &cli.Command{
		Name:      "select",
		Usage:     "Set the default live tracked tab for a session",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			tabName := cmd.Args().First()
			if tabName == "" {
				return fmt.Errorf("tab name required")
			}
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			if err := mgr.SelectTab(ctx, sessionName, tabName); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Tab %q selected\n", tabName)
			return nil
		},
	}
}

func browserTabCloseCmd() *cli.Command {
	return &cli.Command{
		Name:      "close",
		Usage:     "Close and remove a tracked browser tab",
		ArgsUsage: "[name]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "session",
				Usage: "Browser session name; defaults to the selected session when omitted",
			},
		},
		Description: `Close the resolved tracked tab and remove its metadata.

If the closed tab was selected, tap promotes the next remaining live tracked tab
by creation order. If none remain, selected-tab state is cleared.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			sessionName := cmd.String("session")
			tabName := cmd.Args().First()
			if err := mgr.CloseTab(ctx, sessionName, tabName); err != nil {
				return err
			}
			if tabName == "" {
				tabName = "(resolved)"
			}
			fmt.Fprintf(os.Stderr, "Tab %q closed\n", tabName)
			return nil
		},
	}
}
