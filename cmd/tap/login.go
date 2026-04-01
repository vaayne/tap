package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/transport"
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

			url := cmd.Args().First()
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			// Force visible browser for login
			client, err := newClientWithOverrides(cmd, true)
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

// waitForEnter blocks until the user presses Enter or the context is cancelled.
func waitForEnter(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// terminalPause returns a PauseFunc that waits for Enter on stdin.
func terminalPause() transport.PauseFunc {
	return waitForEnter
}
