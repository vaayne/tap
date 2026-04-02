//go:build windows

package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:  "browser",
		Usage: "Manage persistent browser sessions (not supported on Windows)",
		Action: func(_ context.Context, _ *cli.Command) error {
			return cli.Exit("browser session management is not supported on Windows", 1)
		},
	}
}
