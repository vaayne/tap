package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

const DefaultAgentBrowserSession = "default"

type AgentBrowser struct {
	Path        string
	SessionName string
	ProfileDir  string
	Headed      bool
	Attached    bool
	Engine      string
}

type OpenOpts struct {
	Headed     bool
	Headers    map[string]string
	InitScript string
}

type ExecResult struct {
	Stdout json.RawMessage
	Stderr string
}

func NewAgentBrowser(path string) (*AgentBrowser, error) {
	if path == "" {
		resolved, err := ResolveAgentBrowserPath()
		if err != nil {
			return nil, err
		}
		path = resolved
	}
	return &AgentBrowser{Path: path, SessionName: DefaultAgentBrowserSession}, nil
}

func (a *AgentBrowser) Exec(ctx context.Context, args ...string) (json.RawMessage, string, error) {
	return a.exec(ctx, nil, args...)
}

func (a *AgentBrowser) exec(ctx context.Context, stdin []byte, args ...string) (json.RawMessage, string, error) {
	cmdArgs := a.commandArgs(args...)
	cmd := exec.CommandContext(ctx, a.Path, cmdArgs...)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, stderr.String(), fmt.Errorf("agent-browser %s: %w: %s", strings.Join(cmdArgs, " "), err, strings.TrimSpace(stderr.String()))
	}
	return json.RawMessage(bytes.TrimSpace(stdout.Bytes())), stderr.String(), nil
}

func (a *AgentBrowser) commandArgs(args ...string) []string {
	out := make([]string, 0, len(args)+8)
	out = append(out, args...)
	if !slices.Contains(out, "--json") {
		out = append(out, "--json")
	}
	if !a.Attached && a.SessionName != "" && !hasFlag(out, "--session-name") {
		out = append(out, "--session-name", a.SessionName)
	}
	if a.ProfileDir != "" && !hasFlag(out, "--profile") {
		out = append(out, "--profile", a.ProfileDir)
	}
	if a.Headed && !slices.Contains(out, "--headed") {
		out = append(out, "--headed")
	}
	if a.Engine != "" && !hasFlag(out, "--engine") {
		out = append(out, "--engine", a.Engine)
	}
	return out
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag || strings.HasPrefix(arg, flag+"=") {
			return true
		}
	}
	return false
}

func (a *AgentBrowser) Open(ctx context.Context, url string, opts OpenOpts) error {
	args := []string{"open", url}
	if opts.Headed {
		args = append(args, "--headed")
	}
	if opts.InitScript != "" {
		args = append(args, "--init-script", opts.InitScript)
	}
	for name, value := range opts.Headers {
		args = append(args, "--headers", name+": "+value)
	}
	_, _, err := a.Exec(ctx, args...)
	return err
}

func (a *AgentBrowser) Eval(ctx context.Context, js string) (any, error) {
	out, stderr, err := a.exec(ctx, []byte(js), "eval", "--stdin")
	if err != nil {
		return nil, err
	}
	var envelope AgentBrowserEnvelope[map[string]any]
	if err := json.Unmarshal(out, &envelope); err == nil && envelope.Success {
		return envelope.Data["result"], nil
	}
	var value any
	if err := json.Unmarshal(out, &value); err != nil {
		return nil, fmt.Errorf("parse eval JSON: %w: %s", err, stderr)
	}
	return value, nil
}

func (a *AgentBrowser) GetHTML(ctx context.Context) (string, error) {
	value, err := a.Eval(ctx, "document.documentElement.outerHTML")
	if err != nil {
		return "", err
	}
	if s, ok := value.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("agent-browser html result is %T, not string", value)
}

func (a *AgentBrowser) Close(ctx context.Context) error {
	_, _, err := a.Exec(ctx, "close")
	return err
}
