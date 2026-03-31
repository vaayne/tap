package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	// Locate sites/ dir relative to executable
	sitesDir := filepath.Join(".", "sites")
	reg, err := NewRegistry(sitesDir)
	if err != nil {
		log.Fatalf("Failed to load scripts: %v", err)
	}

	cmd := os.Args[1]

	if cmd == "list" {
		listScripts(reg)
		return
	}

	// Treat cmd as script name, rest as key=value args
	script, ok := reg.Get(cmd)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown script: %s\n\n", cmd)
		fmt.Fprintf(os.Stderr, "Available scripts:\n")
		listScripts(reg)
		os.Exit(1)
	}

	args := parseArgs(os.Args[2:])

	// Validate required args
	for name, def := range script.Meta.Args {
		if def.Required {
			if _, ok := args[name]; !ok {
				log.Fatalf("Missing required arg: %s (%s)", name, def.Description)
			}
		}
	}

	log.Printf("Running: %s — %s", script.Meta.Name, script.Meta.Description)

	result, err := runScript(script, args)
	if err != nil {
		log.Fatalf("Failed to run script: %v", err)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(out))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  cdp <script-name> [key=value ...]\n")
	fmt.Fprintf(os.Stderr, "  cdp list\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  cdp v2ex/hot\n")
	fmt.Fprintf(os.Stderr, "  cdp twitter/search query=claude\n")
	fmt.Fprintf(os.Stderr, "  cdp bilibili/search keyword=编程\n")
}

func listScripts(reg *Registry) {
	for _, s := range reg.List() {
		argHints := ""
		if len(s.Meta.Args) > 0 {
			var parts []string
			for name, def := range s.Meta.Args {
				tag := name
				if def.Required {
					tag = name + "*"
				}
				parts = append(parts, tag)
			}
			argHints = " [" + strings.Join(parts, ", ") + "]"
		}
		fmt.Printf("  %-30s %s%s\n", s.Meta.Name, s.Meta.Description, argHints)
	}
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

func runScript(script *Script, args map[string]string) (any, error) {
	ctx, cancel := newBrowserContext()
	defer cancel()

	// Serialize args to JSON for the JS side
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	// Wrap: (async function(args) { ... })({...args...})
	js := fmt.Sprintf("(%s)(%s)", script.Body, string(argsJSON))

	var result any
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp run: %w", err)
	}

	return result, nil
}
