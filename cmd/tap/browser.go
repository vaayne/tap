package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/agentbrowser"
)

type browserExitError struct {
	code int
}

func (e *browserExitError) Error() string {
	return fmt.Sprintf("agent-browser exited with status %d", e.code)
}

func browserCmd() *cli.Command {
	return &cli.Command{
		Name:            "browser",
		Usage:           "Run an agent-browser command",
		ArgsUsage:       "[agent-browser args...]",
		SkipFlagParsing: true,
		Description: `Pass all arguments, input, output, environment, and exit status through to
agent-browser. Tap resolves the same bundled, configured, or PATH executable
used by its other commands and does not interpret browser commands.

Examples:
  tap browser --engine lightpanda open https://example.com
  tap browser snapshot -i
  tap browser click @e3
  tap browser eval --stdin`,
		Action: runBrowser,
	}
}

func runBrowser(ctx context.Context, cmd *cli.Command) error {
	path, err := agentbrowser.New(cmd.String("agent-browser")).Path()
	if err != nil {
		return err
	}

	child := exec.CommandContext(ctx, path, cmd.Args().Slice()...)
	child.Stdin = cmd.Root().Reader
	if child.Stdin == nil {
		child.Stdin = os.Stdin
	}
	child.Stdout = cmd.Root().Writer
	if child.Stdout == nil {
		child.Stdout = os.Stdout
	}
	child.Stderr = cmd.Root().ErrWriter
	if child.Stderr == nil {
		child.Stderr = os.Stderr
	}

	if err := child.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() >= 0 {
			return &browserExitError{code: exitErr.ExitCode()}
		}
		return fmt.Errorf("run agent-browser: %w", err)
	}
	return nil
}
