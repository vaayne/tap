package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

// browserSetCmd returns the "tap browser set" parent command with emulation subcommands.
func browserSetCmd() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Configure emulation overrides for the current tab",
		Description: `Persist and immediately apply emulation settings.
Settings survive across invocations — they are re-applied automatically every
time a tab is resolved for any browser command.

Subcommands:
  viewport <w> <h> [scale]   Set viewport dimensions and optional device scale factor
  device <name>              Emulate a device preset (e.g. "iPhone 14")
  geo <lat> <lng>            Override geolocation
  offline on|off             Toggle offline mode
  headers <json>             Set extra HTTP request headers ({"Name":"Value",...})
  media dark|light           Emulate prefers-color-scheme
  useragent <ua>             Override User-Agent string
  clear                      Remove all persisted emulation settings`,
		Commands: []*cli.Command{
			browserSetViewportCmd(),
			browserSetDeviceCmd(),
			browserSetGeoCmd(),
			browserSetOfflineCmd(),
			browserSetHeadersCmd(),
			browserSetMediaCmd(),
			browserSetUserAgentCmd(),
			browserSetClearCmd(),
		},
	}
}

func browserSetViewportCmd() *cli.Command {
	return &cli.Command{
		Name:      "viewport",
		Usage:     "Set viewport dimensions",
		ArgsUsage: "<width> <height> [scale]",
		Flags:     browserActionFlags(false),
		Description: `Override the browser viewport size and optional device scale factor.

Examples:
  tap browser set viewport 1280 720
  tap browser set viewport 390 844 3.0`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser set viewport <width> <height> [scale]")
			}
			w, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || w <= 0 {
				return fmt.Errorf("width must be a positive integer, got %q", args[0])
			}
			h, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil || h <= 0 {
				return fmt.Errorf("height must be a positive integer, got %q", args[1])
			}
			scale := 1.0
			if len(args) >= 3 {
				scale, err = strconv.ParseFloat(args[2], 64)
				if err != nil || scale <= 0 {
					return fmt.Errorf("scale must be a positive number, got %q", args[2])
				}
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{
				ViewportWidth:  w,
				ViewportHeight: h,
				ViewportScale:  scale,
			}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Viewport set to %dx%d (scale %.2g)\n", w, h, scale)
			return nil
		},
	}
}

func browserSetDeviceCmd() *cli.Command {
	return &cli.Command{
		Name:      "device",
		Usage:     "Emulate a device preset",
		ArgsUsage: "<name>",
		Flags:     browserActionFlags(false),
		Description: `Emulate a device using a chromedp/device preset. Sets viewport, user agent,
and touch emulation in one shot.

Common devices: "iPhone 14", "iPhone 14 Pro", "Pixel 5", "iPad Pro",
"Galaxy S8", "Nexus 5X", "Kindle Fire HDX"

Examples:
  tap browser set device "iPhone 14"
  tap browser set device "Pixel 5"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			name := strings.TrimSpace(cmd.Args().First())
			if name == "" {
				return fmt.Errorf("device name required")
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{DeviceName: name}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Device emulation set to %q\n", name)
			return nil
		},
	}
}

func browserSetGeoCmd() *cli.Command {
	return &cli.Command{
		Name:      "geo",
		Usage:     "Override geolocation",
		ArgsUsage: "<lat> <lng>",
		Flags:     browserActionFlags(false),
		Description: `Override the browser's geolocation to the given latitude and longitude.

Examples:
  tap browser set geo 37.7749 -122.4194
  tap browser set geo 51.5074 -0.1278`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args().Slice()
			if len(args) < 2 {
				return fmt.Errorf("usage: tap browser set geo <lat> <lng>")
			}
			lat, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("lat must be a number, got %q", args[0])
			}
			lng, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("lng must be a number, got %q", args[1])
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{GeoLat: &lat, GeoLng: &lng}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Geolocation set to (%.6f, %.6f)\n", lat, lng)
			return nil
		},
	}
}

func browserSetOfflineCmd() *cli.Command {
	return &cli.Command{
		Name:      "offline",
		Usage:     "Toggle offline mode",
		ArgsUsage: "on|off",
		Flags:     browserActionFlags(false),
		Description: `Emulate network offline (on) or restore connectivity (off).

Examples:
  tap browser set offline on
  tap browser set offline off`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			arg := strings.ToLower(strings.TrimSpace(cmd.Args().First()))
			var offline bool
			switch arg {
			case "on", "true", "1":
				offline = true
			case "off", "false", "0":
				offline = false
			default:
				return fmt.Errorf("expected on|off, got %q", arg)
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{Offline: &offline}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			if offline {
				fmt.Fprintln(os.Stderr, "Offline mode enabled")
			} else {
				fmt.Fprintln(os.Stderr, "Offline mode disabled")
			}
			return nil
		},
	}
}

func browserSetHeadersCmd() *cli.Command {
	return &cli.Command{
		Name:      "headers",
		Usage:     "Set extra HTTP request headers",
		ArgsUsage: "<json>",
		Flags:     browserActionFlags(false),
		Description: `Add extra HTTP headers sent with every request from the tab.
The argument must be a JSON object mapping header names to values.

Examples:
  tap browser set headers '{"Authorization":"Bearer token123"}'
  tap browser set headers '{"X-Custom-Header":"value","Accept-Language":"en-US"}'`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			raw := strings.TrimSpace(cmd.Args().First())
			if raw == "" {
				return fmt.Errorf("JSON headers object required")
			}
			var headers map[string]string
			if err := json.Unmarshal([]byte(raw), &headers); err != nil {
				return fmt.Errorf("parse headers JSON: %w", err)
			}
			if len(headers) == 0 {
				return fmt.Errorf("headers object must not be empty")
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{Headers: headers}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Extra headers set (%d header(s))\n", len(headers))
			return nil
		},
	}
}

func browserSetMediaCmd() *cli.Command {
	return &cli.Command{
		Name:      "media",
		Usage:     "Emulate prefers-color-scheme",
		ArgsUsage: "dark|light",
		Flags:     browserActionFlags(false),
		Description: `Force prefers-color-scheme to dark or light for the current tab.

Examples:
  tap browser set media dark
  tap browser set media light`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			scheme := strings.ToLower(strings.TrimSpace(cmd.Args().First()))
			if scheme != "dark" && scheme != "light" {
				return fmt.Errorf("expected dark|light, got %q", scheme)
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{MediaScheme: scheme}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "prefers-color-scheme set to %q\n", scheme)
			return nil
		},
	}
}

func browserSetUserAgentCmd() *cli.Command {
	return &cli.Command{
		Name:      "useragent",
		Usage:     "Override the User-Agent string",
		ArgsUsage: "<ua>",
		Flags:     browserActionFlags(false),
		Description: `Override the browser User-Agent sent with every request.

Examples:
  tap browser set useragent "Mozilla/5.0 (compatible; MyBot/1.0)"
  tap browser set useragent "Googlebot/2.1 (+http://www.google.com/bot.html)"`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			ua := strings.TrimSpace(cmd.Args().First())
			if ua == "" {
				return fmt.Errorf("user-agent string required")
			}

			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			e := &browser.EmulationSettings{UserAgent: ua}
			if err := mgr.SetEmulation(ctx, cmd.String("session"), cmd.String("tab"), e); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "User-Agent set to %q\n", ua)
			return nil
		},
	}
}

func browserSetClearCmd() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Remove all persisted emulation settings for the tab",
		Flags: browserActionFlags(false),
		Description: `Wipe all persisted emulation overrides for the current tab.
The browser's active overrides are not reversed until the next navigation or reload.

Examples:
  tap browser set clear`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			mgr, err := newBrowserManager(cmd)
			if err != nil {
				return err
			}
			if err := mgr.ClearEmulation(ctx, cmd.String("session"), cmd.String("tab")); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Emulation settings cleared")
			return nil
		},
	}
}
