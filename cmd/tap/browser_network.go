package main

import "github.com/urfave/cli/v3"

func browserNetworkCmd() *cli.Command {
	return &cli.Command{
		Name:  "network",
		Usage: "Capture and intercept network requests",
		Commands: []*cli.Command{
			passthroughCmd("requests", "List captured network requests", "network", "requests"),
			passthroughCmd("route", "Set request interception rules", "network", "route"),
			passthroughCmd("unroute", "Remove request interception rules", "network", "unroute"),
			passthroughCmd("har", "Start or stop HAR capture", "network", "har"),
		},
	}
}
