package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/agentbrowser"
)

func doctorCmd() *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Check Tap and its agent-browser runtime dependency",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "fix",
				Usage: "Delegate repairs to agent-browser doctor --fix",
			},
		},
		Action: runDoctor,
	}
}

func runDoctor(ctx context.Context, cmd *cli.Command) error {
	color := useColor(cmd)
	ok := checkMark(color)
	fail := failMark(color)
	fmt.Printf("%s tap %s\n", ok, bold(color, version))

	client := agentbrowser.New(cmd.String("agent-browser"))
	path, err := client.Path()
	if err != nil {
		fmt.Printf("%s agent-browser not found\n", fail)
		fmt.Printf("  %s\n", dim(color, "Install: npm install -g agent-browser && agent-browser install"))
		return err
	}
	version, err := client.Version(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s\n", ok, version)
	fmt.Printf("  %s\n", dim(color, path))

	if err := client.Doctor(ctx, cmd.Bool("fix")); err != nil {
		fmt.Printf("%s agent-browser runtime check failed\n", fail)
		return err
	}
	fmt.Printf("%s agent-browser runtime ready\n", ok)
	return nil
}

func checkMark(color bool) string { return green(color, "✓") }

func failMark(color bool) string {
	if !color {
		return "✗"
	}
	return "\033[31m✗\033[0m"
}
