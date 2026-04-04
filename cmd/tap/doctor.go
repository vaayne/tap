package main

import (
	"context"
	"fmt"
	"log"
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
				Usage: "Download missing browser dependencies",
			},
			&cli.BoolFlag{
				Name:  "update",
				Usage: "Update Lightpanda to the latest nightly build",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runDoctor(ctx, cmd)
		},
	}
}

func runDoctor(ctx context.Context, cmd *cli.Command) error {
	// Suppress library log output in doctor; we print our own messages.
	log.SetOutput(nopWriter{})

	color := useColor(cmd)
	install := cmd.Bool("install")
	update := cmd.Bool("update")

	ok := checkMark(color)
	fail := failMark(color)
	warn := warnMark(color)

	// --- tap version ---
	fmt.Printf("%s tap %s\n", ok, bold(color, version))

	// --- Chrome ---
	chrome := browser.DetectChrome()
	if chrome != nil {
		v := chrome.Version
		if v == "" {
			v = "unknown version"
		}
		fmt.Printf("%s Chrome %s\n", ok, v)
		fmt.Printf("  %s\n", dim(color, chrome.Path))
	} else {
		fmt.Printf("%s Chrome not found\n", fail)
		fmt.Printf("  %s\n", dim(color, "Install Chrome or use --lightpanda as an alternative"))
	}

	// --- Lightpanda ---
	lp := browser.NewLightpanda("", "")
	if update {
		return updateLightpanda(ctx, lp, color, ok)
	}

	if lp.Installed() {
		meta, _ := lp.ReadMeta()
		if meta != nil {
			age := time.Since(meta.DownloadedAt)
			ageStr := formatAge(age)
			fmt.Printf("%s Lightpanda installed (%s ago)\n", ok, ageStr)
			fmt.Printf("  %s\n", dim(color, "Downloaded: "+meta.DownloadedAt.Format(time.RFC3339)))
			if age > 7*24*time.Hour {
				fmt.Printf("  %s Run %s to get the latest nightly\n", warn, bold(color, "tap doctor --update"))
			}
		} else {
			fmt.Printf("%s Lightpanda installed (unknown age)\n", warn)
			fmt.Printf("  %s Run %s to get the latest nightly\n", warn, bold(color, "tap doctor --update"))
		}
	} else if install {
		return installLightpanda(ctx, lp, color, ok)
	} else {
		fmt.Printf("%s Lightpanda not installed\n", warn)
		fmt.Printf("  %s Run %s to download it\n", dim(color, "→"), bold(color, "tap doctor --install"))
	}

	return nil
}

func installLightpanda(ctx context.Context, lp *browser.Lightpanda, color bool, ok string) error {
	fmt.Print("  Downloading Lightpanda... ")
	if err := lp.EnsureInstalled(ctx); err != nil {
		return fmt.Errorf("install lightpanda: %w", err)
	}
	fmt.Println("done")
	fmt.Printf("%s Lightpanda installed\n", ok)
	return nil
}

func updateLightpanda(ctx context.Context, lp *browser.Lightpanda, color bool, ok string) error {
	fmt.Print("  Updating Lightpanda to latest nightly... ")
	if err := lp.Update(ctx); err != nil {
		return fmt.Errorf("update lightpanda: %w", err)
	}
	fmt.Println("done")
	fmt.Printf("%s Lightpanda updated\n", ok)
	return nil
}

func checkMark(color bool) string {
	return green(color, "✓")
}

func failMark(color bool) string {
	if !color {
		return "✗"
	}
	return "\033[31m✗\033[0m"
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
