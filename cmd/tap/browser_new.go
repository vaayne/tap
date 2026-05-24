package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// passthroughCmd creates a simple command that forwards arguments to agent-browser.
func passthroughCmd(name, usage string, prefixArgs ...string) *cli.Command {
	return &cli.Command{
		Name:            name,
		Usage:           usage,
		ArgsUsage:       "[args...]",
		SkipFlagParsing: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ab, err := newAgentBrowser(cmd)
			if err != nil {
				return err
			}
			rawArgs, sessionName := extractPassthroughSession(cmd.Args().Slice())
			if sessionName != "" {
				ab.SessionName = sessionName
				ab.Attached = false
			}
			execArgs := append([]string{}, prefixArgs...)
			execArgs = append(execArgs, rawArgs...)
			out, _, err := ab.Exec(ctx, execArgs...)
			if err != nil {
				return err
			}
			if len(out) > 0 {
				fmt.Println(string(out))
			}
			return nil
		},
	}
}

func extractPassthroughSession(args []string) ([]string, string) {
	out := make([]string, 0, len(args))
	var session string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--session" {
			if i+1 < len(args) {
				session = args[i+1]
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "--session=") {
			session = strings.TrimPrefix(arg, "--session=")
			continue
		}
		out = append(out, arg)
	}
	return out, session
}

func browserSetCmd() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Configure browser emulation settings",
		Commands: []*cli.Command{
			passthroughCmd("viewport", "Set viewport size", "set", "viewport"),
			passthroughCmd("device", "Emulate a device", "set", "device"),
			passthroughCmd("geo", "Set geolocation", "set", "geo"),
			passthroughCmd("offline", "Toggle offline mode", "set", "offline"),
			passthroughCmd("headers", "Set extra HTTP headers", "set", "headers"),
			passthroughCmd("credentials", "Set HTTP auth credentials", "set", "credentials"),
			passthroughCmd("media", "Set CSS media type", "set", "media"),
		},
	}
}

func browserStorageCmd() *cli.Command {
	return &cli.Command{
		Name:  "storage",
		Usage: "Manage local/session storage",
		Commands: []*cli.Command{
			passthroughCmd("local", "Manage localStorage", "storage", "local"),
			passthroughCmd("session", "Manage sessionStorage", "storage", "session"),
		},
	}
}

func browserStateCmd() *cli.Command {
	return &cli.Command{
		Name:  "state",
		Usage: "Manage saved browser states",
		Commands: []*cli.Command{
			passthroughCmd("save", "Save current state", "state", "save"),
			passthroughCmd("load", "Load saved state", "state", "load"),
			passthroughCmd("list", "List saved states", "state", "list"),
			passthroughCmd("show", "Show state details", "state", "show"),
			passthroughCmd("clear", "Clear saved states", "state", "clear"),
		},
	}
}

func browserAuthCmd() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Manage authentication state",
		Commands: []*cli.Command{
			passthroughCmd("save", "Save auth credentials", "auth", "save"),
			passthroughCmd("login", "Perform login", "auth", "login"),
			passthroughCmd("list", "List saved auth", "auth", "list"),
			passthroughCmd("show", "Show auth details", "auth", "show"),
			passthroughCmd("delete", "Delete saved auth", "auth", "delete"),
		},
	}
}

func browserGetCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get page properties and values",
		Commands: []*cli.Command{
			passthroughCmd("text", "Get element or page text", "get", "text"),
			passthroughCmd("html", "Get element or page HTML", "get", "html"),
			passthroughCmd("value", "Get input value", "get", "value"),
			passthroughCmd("attr", "Get element attribute", "get", "attr"),
			passthroughCmd("title", "Get page title", "get", "title"),
			passthroughCmd("url", "Get current URL", "get", "url"),
			passthroughCmd("count", "Count elements matching selector", "get", "count"),
			passthroughCmd("box", "Get element bounding box", "get", "box"),
			passthroughCmd("styles", "Get computed styles", "get", "styles"),
		},
	}
}

func browserVitalsCmd() *cli.Command {
	return passthroughCmd("vitals", "Get Web Vitals metrics", "vitals")
}

func browserDiffCmd() *cli.Command {
	return &cli.Command{
		Name:  "diff",
		Usage: "Compare snapshots, screenshots, or URLs",
		Commands: []*cli.Command{
			passthroughCmd("snapshot", "Compare page snapshots", "diff", "snapshot"),
			passthroughCmd("screenshot", "Compare screenshots", "diff", "screenshot"),
			passthroughCmd("url", "Compare URL responses", "diff", "url"),
		},
	}
}
