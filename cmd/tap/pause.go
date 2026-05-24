package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
	"github.com/vaayne/tap/transport"
)

func hasPauseMode(cmd *cli.Command) bool {
	return pauseModeCount(cmd) > 0
}

func pauseModeCount(cmd *cli.Command) int {
	count := 0
	if cmd.Bool("pause") {
		count++
	}
	if cmd.Duration("delay") > 0 {
		count++
	}
	if cmd.String("wait-selector") != "" {
		count++
	}
	if cmd.String("wait-js") != "" {
		count++
	}
	return count
}

func resolvePauseFunc(cmd *cli.Command) (transport.PauseFunc, error) {
	if pauseModeCount(cmd) > 1 {
		return nil, fmt.Errorf("choose only one of --pause, --delay, --wait-selector, or --wait-js")
	}

	if cmd.Bool("pause") {
		return waitForEnter, nil
	}
	if d := cmd.Duration("delay"); d > 0 {
		return delayPause(d), nil
	}
	if selector := cmd.String("wait-selector"); selector != "" {
		return waitForSelector(selector), nil
	}
	if expr := cmd.String("wait-js"); expr != "" {
		return waitForJS(expr), nil
	}
	return nil, nil
}

func delayPause(d time.Duration) transport.PauseFunc {
	return func(ctx context.Context) error {
		timer := time.NewTimer(d)
		defer timer.Stop()

		select {
		case <-timer.C:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func waitForSelector(selector string) transport.PauseFunc {
	return func(ctx context.Context) error {
		path, err := browser.ResolveAgentBrowserPath()
		if err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, path, "wait", selector, "--session-name", "default", "--json")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

func waitForJS(expr string) transport.PauseFunc {
	return func(ctx context.Context) error {
		path, err := browser.ResolveAgentBrowserPath()
		if err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, path, "wait", "--fn", expr, "--session-name", "default", "--json")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

// waitForEnter blocks until the user presses Enter or the context is cancelled.
func waitForEnter(ctx context.Context) error {
	if !stdinIsTTY() {
		return fmt.Errorf("--pause requires an interactive terminal")
	}

	done := make(chan error, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		_, err := reader.ReadString('\n')
		if err != nil {
			done <- fmt.Errorf("read input: %w", err)
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func stdinIsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
