package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func loginCmd() *cli.Command {
	return &cli.Command{
		Name:      "login",
		Usage:     "Open a browser to log in or interact with a site",
		ArgsUsage: "<url>",
		Description: `Opens a visible browser window and navigates to the given URL.
The browser stays open until you press Enter in the terminal,
giving you time to log in, solve CAPTCHAs, or perform any
manual interaction. Cookies are saved to the Chrome profile
directory so subsequent tap commands are authenticated.

Examples:
  tap login https://github.com/login
  tap login https://twitter.com/i/flow/login
  tap login --profile-dir ~/.tap/work https://internal.corp.com`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("URL required (e.g. tap login https://github.com/login)")
			}

			url := normalizeURL(cmd.Args().First())

			// Force visible browser for login
			client, err := newClientWithOverrides(ctx, cmd, true)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()

			color := useColor(cmd)
			fmt.Printf("%s %s\n", bold(color, "Opening browser:"), url)
			fmt.Printf("%s\n", dim(color, "Press Enter when done to close the browser and save cookies..."))

			return client.Login(ctx, url, waitForEnter)
		},
	}
}
