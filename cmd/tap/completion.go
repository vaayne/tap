package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/script"
)

func configureCompletionCommand(completion *cli.Command) {
	completion.Hidden = false
	completion.ArgsUsage = "<shell>"
	completion.Description = `Output shell completion for tap.

Examples:
  tap completion bash > ~/.local/share/bash-completion/completions/tap
  tap completion zsh > ~/.zfunc/_tap
  tap completion fish > ~/.config/fish/completions/tap.fish
  tap completion pwsh > ~/.config/powershell/tap.ps1`
}

func completeSiteRoot(ctx context.Context, cmd *cli.Command) {
	prefix, completeFlags := completionState(cmd)
	if completeFlags {
		cli.DefaultCompleteWithFlags(ctx, cmd)
		return
	}

	printVisibleCommands(cmd, prefix)
	printScriptCompletions(cmd, prefix)
}

func completeSiteScripts(ctx context.Context, cmd *cli.Command) {
	prefix, completeFlags := completionState(cmd)
	if completeFlags {
		cli.DefaultCompleteWithFlags(ctx, cmd)
		return
	}
	printScriptCompletions(cmd, prefix)
}

func completionState(cmd *cli.Command) (prefix string, completeFlags bool) {
	prefix = completionPrefix(cmd)
	return prefix, strings.HasPrefix(prefix, "-")
}

func completionPrefix(cmd *cli.Command) string {
	args := cmd.Args().Slice()
	if len(args) == 0 {
		return ""
	}
	return args[len(args)-1]
}

func printVisibleCommands(cmd *cli.Command, prefix string) {
	for _, sub := range cmd.Commands {
		if sub.Hidden {
			continue
		}
		if prefix != "" && !strings.HasPrefix(sub.Name, prefix) {
			continue
		}
		printCompletion(cmd, sub.Name, sub.Usage)
	}
}

func printScriptCompletions(cmd *cli.Command, prefix string) {
	reg, err := loadCompletionRegistry(cmd)
	if err != nil {
		return
	}
	for _, s := range reg.List() {
		if prefix != "" && !strings.HasPrefix(s.Meta.Name, prefix) {
			continue
		}
		printCompletion(cmd, s.Meta.Name, s.Meta.Description)
	}
}

func loadCompletionRegistry(cmd *cli.Command) (*script.Registry, error) {
	dir := cmd.String("sites-dir")
	if dir == "" {
		dir = defaultSitesDir()
	}

	overrideDir := defaultLocalOverrideDir()
	if cmd.Bool("local-only") {
		return loadRegistryDir(overrideDir)
	}

	reg, err := script.NewRegistryWithOverride(dir, overrideDir)
	if err == nil {
		return reg, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return loadRegistryDir(overrideDir)
}

func loadRegistryDir(dir string) (*script.Registry, error) {
	reg, err := script.NewRegistry(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return emptyRegistry(), nil
		}
		return nil, err
	}
	return reg, nil
}

func emptyRegistry() *script.Registry {
	return &script.Registry{}
}

func printCompletion(cmd *cli.Command, value, description string) {
	shell := os.Getenv("SHELL")
	if description != "" && (strings.HasSuffix(shell, "zsh") || strings.HasSuffix(shell, "fish")) {
		_, _ = fmt.Fprintf(cmd.Root().Writer, "%s:%s\n", value, description)
		return
	}
	_, _ = fmt.Fprintln(cmd.Root().Writer, value)
}
