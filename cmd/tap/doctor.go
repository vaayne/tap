package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func doctorCmd() *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Check and manage browser dependencies",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "install",
				Usage: "Extract embedded agent-browser and let it install browser dependencies",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runDoctor(ctx, cmd)
		},
	}
}

func runDoctor(ctx context.Context, cmd *cli.Command) error {
	log.SetOutput(nopWriter{})

	color := useColor(cmd)
	install := cmd.Bool("install")

	ok := checkMark(color)
	warn := warnMark(color)

	fmt.Printf("%s tap %s\n", ok, bold(color, version))

	agentInstall := browser.NewAgentBrowserInstall("")
	if install {
		action := "Extracting"
		if agentInstall.Installed() {
			action = "Refreshing"
		}
		fmt.Printf("  %s embedded agent-browser... ", action)
		if err := agentInstall.Update(ctx); err != nil {
			return fmt.Errorf("extract agent-browser: %w", err)
		}
		fmt.Println("done")

		path, err := browser.ResolveAgentBrowserPath()
		if err != nil {
			return err
		}
		fmt.Printf("  Running agent-browser install...\n")
		installCmd := exec.CommandContext(ctx, path, "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("agent-browser install: %w", err)
		}
	}

	if path, err := browser.ResolveAgentBrowserPath(); err == nil {
		fmt.Printf("%s agent-browser embedded (%s)\n", ok, browser.AgentBrowserVersion)
		fmt.Printf("  %s\n", dim(color, path))
		if meta, _ := agentInstall.ReadMeta(); meta != nil {
			age := time.Since(meta.InstalledAt)
			fmt.Printf("  %s\n", dim(color, "Extracted: "+meta.InstalledAt.Format(time.RFC3339)+" ("+formatAge(age)+" ago)"))
		}
	} else {
		fmt.Printf("%s agent-browser unavailable\n", warn)
		fmt.Printf("  %s\n", dim(color, err.Error()))
	}

	fmt.Printf("%s Chrome managed by agent-browser\n", ok)
	fmt.Printf("  %s\n", dim(color, "Run tap doctor --install to let agent-browser install browser dependencies"))

	return nil
}

func checkMark(color bool) string {
	return green(color, "✓")
}

func warnMark(color bool) string {
	return yellow(color, "!")
}

func formatAge(d time.Duration) string {
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
