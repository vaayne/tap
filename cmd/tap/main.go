package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
	"github.com/vaayne/tap/fetch"
	"github.com/vaayne/tap/script"
)

var version = "dev"

// output format constants
const (
	formatPretty = "pretty"
	formatJSON   = "json"
	formatRaw    = "raw"
)

func main() {
	_ = godotenv.Load()

	app := &cli.Command{
		Name:    "tap",
		Usage:   "Tap into any website from your terminal",
		Version: version,
		Flags:   globalFlags(),
		Commands: []*cli.Command{
			siteCmd(),
			fetchCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		if snf, ok := err.(*tap.ScriptNotFoundError); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", snf.Error())
			if suggestions := snf.Suggestions(5); len(suggestions) > 0 {
				fmt.Fprintf(os.Stderr, "\nDid you mean?\n")
				for _, s := range suggestions {
					fmt.Fprintf(os.Stderr, "  tap site %s\n", s)
				}
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "sites-dir",
			Usage:   "Directory containing site scripts",
			Value:   "./sites",
			Sources: cli.EnvVars("TAP_SITES_DIR"),
		},
		&cli.StringFlag{
			Name:    "ws-url",
			Usage:   "Remote CDP WebSocket URL",
			Sources: cli.EnvVars("TAP_WS_URL"),
		},
		&cli.StringFlag{
			Name:    "profile-dir",
			Usage:   "Chrome profile directory for persistent cookies",
			Sources: cli.EnvVars("TAP_PROFILE_DIR"),
		},
		&cli.BoolFlag{
			Name:    "browser",
			Aliases: []string{"b"},
			Usage:   "Force browser execution, skip QuickJS",
			Sources: cli.EnvVars("TAP_BROWSER"),
		},
		&cli.BoolFlag{
			Name:  "no-headless",
			Usage: "Run browser in visible mode (useful for debugging auth)",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Usage:   "Execution timeout (e.g., 30s, 2m)",
			Value:   0,
			Sources: cli.EnvVars("TAP_TIMEOUT"),
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging",
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "Suppress all log output",
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "Disable colored output",
			Sources: cli.EnvVars("NO_COLOR"),
		},
	}
}

func newClient(cmd *cli.Command) (*tap.Client, error) {
	var opts []tap.Option

	if dir := cmd.String("sites-dir"); dir != "" {
		opts = append(opts, tap.WithSitesDir(dir))
	}
	if url := cmd.String("ws-url"); url != "" {
		opts = append(opts, tap.WithWSURL(url))
	}
	if dir := cmd.String("profile-dir"); dir != "" {
		opts = append(opts, tap.WithProfileDir(dir))
	}
	if cmd.Bool("browser") {
		opts = append(opts, tap.WithForceBrowser(true))
	}
	if cmd.Bool("no-headless") {
		opts = append(opts, tap.WithHeadless(false))
	}
	if d := cmd.Duration("timeout"); d > 0 {
		opts = append(opts, tap.WithTimeout(d))
	}

	return tap.New(opts...)
}

// configureLogging sets up log output based on --verbose/--quiet flags.
func configureLogging(cmd *cli.Command) {
	if cmd.Bool("quiet") {
		log.SetOutput(nopWriter{})
	} else if !cmd.Bool("verbose") {
		// Default: suppress log output (only errors via stderr)
		log.SetOutput(nopWriter{})
	}
	// verbose: keep default log output (stderr)
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

// ---------- color helpers ----------

func useColor(cmd *cli.Command) bool {
	if cmd.Bool("no-color") {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func bold(color bool, s string) string {
	if !color {
		return s
	}
	return "\033[1m" + s + "\033[0m"
}

func dim(color bool, s string) string {
	if !color {
		return s
	}
	return "\033[2m" + s + "\033[0m"
}

func green(color bool, s string) string {
	if !color {
		return s
	}
	return "\033[32m" + s + "\033[0m"
}

func yellow(color bool, s string) string {
	if !color {
		return s
	}
	return "\033[33m" + s + "\033[0m"
}

func cyan(color bool, s string) string {
	if !color {
		return s
	}
	return "\033[36m" + s + "\033[0m"
}

// ---------- site command ----------

func siteCmd() *cli.Command {
	return &cli.Command{
		Name:      "site",
		Usage:     "Run site scripts",
		ArgsUsage: "<script-name> [key=value ...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json, pretty (default), raw",
				Value:   formatPretty,
			},
		},
		Commands: []*cli.Command{
			siteListCmd(),
			siteInfoCmd(),
			siteSearchCmd(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			args := cmd.Args()
			if args.Len() == 0 {
				return fmt.Errorf("script name required. Run 'tap site list' to see available scripts")
			}

			scriptName := args.First()
			scriptArgs := parseArgs(args.Tail())

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			if cmd.Bool("verbose") {
				mode := "auto (QuickJS → Browser)"
				if cmd.Bool("browser") {
					mode = "browser"
				}
				log.Printf("Running: %s [engine=%s]", scriptName, mode)
			}

			result, err := client.RunScript(ctx, scriptName, scriptArgs)
			if err != nil {
				return err
			}

			return printResult(cmd, result)
		},
	}
}

func siteListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all available scripts (grouped by site)",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			scripts := client.ListScripts()
			if len(scripts) == 0 {
				fmt.Println("No scripts found.")
				return nil
			}

			color := useColor(cmd)
			groups := groupScripts(scripts)

			// Sort group names
			groupNames := make([]string, 0, len(groups))
			for name := range groups {
				groupNames = append(groupNames, name)
			}
			sort.Strings(groupNames)

			for i, groupName := range groupNames {
				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("%s\n", bold(color, groupName+"/"))
				for _, s := range groups[groupName] {
					argHints := formatArgHints(s, color)
					actionName := s.Meta.Name
					// Strip the group prefix for display
					if _, after, ok := strings.Cut(actionName, "/"); ok {
						actionName = after
					}
					fmt.Printf("  %-24s %s%s\n",
						green(color, actionName),
						s.Meta.Description,
						argHints,
					)
				}
			}

			fmt.Printf("\n%s\n", dim(color, fmt.Sprintf("%d scripts across %d sites", len(scripts), len(groups))))
			return nil
		},
	}
}

func siteInfoCmd() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Show detailed info for a script",
		ArgsUsage: "<script-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("script name required")
			}

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			name := cmd.Args().First()
			s, ok := client.GetScript(name)
			if !ok {
				return &tap.ScriptNotFoundError{Name: name, Available: scriptNames(client)}
			}

			color := useColor(cmd)

			fmt.Printf("%s\n", bold(color, s.Meta.Name))
			fmt.Printf("  %s\n\n", s.Meta.Description)

			fmt.Printf("  %s  %s\n", bold(color, "Domain:"), s.Meta.Domain)

			if s.Meta.Example != "" {
				fmt.Printf("  %s %s\n", bold(color, "Example:"), s.Meta.Example)
			}

			if len(s.Meta.Args) > 0 {
				fmt.Printf("\n  %s\n", bold(color, "Arguments:"))
				// Sort args for consistent output
				argNames := make([]string, 0, len(s.Meta.Args))
				for name := range s.Meta.Args {
					argNames = append(argNames, name)
				}
				sort.Strings(argNames)
				for _, argName := range argNames {
					def := s.Meta.Args[argName]
					req := dim(color, "optional")
					if def.Required {
						req = yellow(color, "required")
					}
					fmt.Printf("    %-16s %s  %s\n",
						green(color, argName),
						dim(color, "("+req+")"),
						def.Description,
					)
				}
			}

			fmt.Printf("\n  %s tap site %s", bold(color, "Usage:"), s.Meta.Name)
			for argName, def := range s.Meta.Args {
				if def.Required {
					fmt.Printf(" %s=<%s>", argName, argName)
				} else {
					fmt.Printf(" [%s=value]", argName)
				}
			}
			fmt.Println()

			return nil
		},
	}
}

func siteSearchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search scripts by name or description",
		ArgsUsage: "<query>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("search query required")
			}

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			query := strings.ToLower(strings.Join(cmd.Args().Slice(), " "))
			scripts := client.ListScripts()
			color := useColor(cmd)

			var matches []*script.Script
			for _, s := range scripts {
				name := strings.ToLower(s.Meta.Name)
				desc := strings.ToLower(s.Meta.Description)
				domain := strings.ToLower(s.Meta.Domain)
				if strings.Contains(name, query) || strings.Contains(desc, query) || strings.Contains(domain, query) {
					matches = append(matches, s)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("No scripts matching %q\n", query)
				return nil
			}

			for _, s := range matches {
				argHints := formatArgHints(s, color)
				fmt.Printf("  %-30s %s%s\n",
					green(color, s.Meta.Name),
					s.Meta.Description,
					argHints,
				)
			}
			fmt.Printf("\n%s\n", dim(color, fmt.Sprintf("%d result(s)", len(matches))))
			return nil
		},
	}
}

// ---------- fetch command ----------

func fetchCmd() *cli.Command {
	return &cli.Command{
		Name:      "fetch",
		Usage:     "Fetch and extract clean content from a URL",
		ArgsUsage: "<url>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON with full metadata",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configureLogging(cmd)
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("URL required")
			}

			url := cmd.Args().First()

			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			opts := &fetch.Options{
				Markdown:   true,
				UseBrowser: cmd.Bool("browser"),
			}
			result, err := client.Fetch(ctx, url, opts)
			if err != nil {
				return err
			}

			if cmd.Bool("json") {
				out, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal result: %w", err)
				}
				fmt.Println(string(out))
			} else {
				if result.Title != "" {
					fmt.Printf("# %s\n\n", result.Title)
				}
				if result.Markdown != "" {
					fmt.Println(result.Markdown)
				} else if result.Content != "" {
					fmt.Println(result.Content)
				}
			}
			return nil
		},
	}
}

// ---------- helpers ----------

func printResult(cmd *cli.Command, result any) error {
	format := cmd.String("format")
	switch format {
	case formatRaw:
		out, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		fmt.Println(string(out))
	case formatJSON:
		out, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		fmt.Println(string(out))
	default: // pretty
		out, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		if useColor(cmd) {
			fmt.Println(colorizeJSON(string(out)))
		} else {
			fmt.Println(string(out))
		}
	}
	return nil
}

// colorizeJSON applies basic syntax highlighting to JSON output.
func colorizeJSON(s string) string {
	var b strings.Builder
	inString := false
	isKey := false
	escaped := false
	afterColon := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escaped {
			b.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' && inString {
			b.WriteByte(ch)
			escaped = true
			continue
		}

		if ch == '"' {
			if !inString {
				inString = true
				// Determine if this is a key (not after a colon on same context)
				isKey = !afterColon
				if isKey {
					b.WriteString("\033[36m\"") // cyan for keys
				} else {
					b.WriteString("\033[32m\"") // green for string values
				}
				afterColon = false
			} else {
				b.WriteString("\"\033[0m")
				inString = false
			}
			continue
		}

		if !inString {
			if ch == ':' {
				b.WriteByte(ch)
				afterColon = true
				continue
			}
			if ch == ',' || ch == '{' || ch == '[' {
				afterColon = false
			}
			// Colorize numbers
			if ch >= '0' && ch <= '9' || ch == '-' {
				b.WriteString("\033[33m") // yellow for numbers
				for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || s[i] == '-' || s[i] == 'e' || s[i] == 'E' || s[i] == '+') {
					b.WriteByte(s[i])
					i++
				}
				b.WriteString("\033[0m")
				i-- // back up one
				afterColon = false
				continue
			}
			// Colorize booleans and null
			for _, kw := range []string{"true", "false", "null"} {
				if i+len(kw) <= len(s) && s[i:i+len(kw)] == kw {
					b.WriteString("\033[33m" + kw + "\033[0m") // yellow
					i += len(kw) - 1
					afterColon = false
					goto next
				}
			}
		}

		b.WriteByte(ch)
	next:
	}
	return b.String()
}

func parseArgs(raw []string) map[string]string {
	args := make(map[string]string)
	for _, s := range raw {
		if k, v, ok := strings.Cut(s, "="); ok {
			args[k] = v
		}
	}
	return args
}

func formatArgHints(s *script.Script, color bool) string {
	if len(s.Meta.Args) == 0 {
		return ""
	}
	var parts []string
	for name, def := range s.Meta.Args {
		if def.Required {
			parts = append(parts, yellow(color, name+"*"))
		} else {
			parts = append(parts, dim(color, name))
		}
	}
	sort.Strings(parts) // consistent ordering
	return " [" + strings.Join(parts, ", ") + "]"
}

// groupScripts groups scripts by their site prefix (part before /).
func groupScripts(scripts []*script.Script) map[string][]*script.Script {
	groups := make(map[string][]*script.Script)
	for _, s := range scripts {
		group := s.Meta.Name
		if idx := strings.Index(group, "/"); idx != -1 {
			group = group[:idx]
		}
		groups[group] = append(groups[group], s)
	}
	return groups
}

func scriptNames(client *tap.Client) []string {
	scripts := client.ListScripts()
	names := make([]string, len(scripts))
	for i, s := range scripts {
		names[i] = s.Meta.Name
	}
	return names
}
