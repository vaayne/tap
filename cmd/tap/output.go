package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// output format constants
const (
	formatPretty = "pretty"
	formatJSON   = "json"
	formatRaw    = "raw"
)

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
