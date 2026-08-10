// Package agentbrowser provides the narrow subprocess boundary between Tap and
// the agent-browser CLI. Browser lifecycle and session ownership stay entirely
// inside agent-browser.
package agentbrowser

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultBinary = "agent-browser"
	EnvBinary     = "TAP_AGENT_BROWSER"
)

// Client invokes agent-browser. It deliberately has no session field:
// AGENT_BROWSER_SESSION and every other agent-browser setting are inherited
// from the process environment unchanged.
type Client struct {
	Binary     string
	explicit   bool
	siblingDir string
}

type envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

type evalData struct {
	Result json.RawMessage `json:"result"`
}

type batchResult struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage `json:"error"`
}

// New creates a thin client. Binary lookup remains lazy so registry-only
// commands such as `tap site list` work even before agent-browser is installed.
func New(binary string) *Client {
	explicit := binary != ""
	if binary == "" {
		binary = os.Getenv(EnvBinary)
		explicit = binary != ""
	}
	if binary == "" {
		binary = DefaultBinary
	}
	var siblingDir string
	if executable, err := os.Executable(); err == nil {
		siblingDir = filepath.Dir(executable)
	}
	return &Client{Binary: binary, explicit: explicit, siblingDir: siblingDir}
}

// Path resolves the configured executable without installing or downloading it.
func (c *Client) Path() (string, error) {
	if !c.explicit && c.siblingDir != "" {
		name := DefaultBinary
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		if path, err := exec.LookPath(filepath.Join(c.siblingDir, name)); err == nil {
			return path, nil
		}
	}
	path, err := exec.LookPath(c.Binary)
	if err != nil {
		return "", fmt.Errorf("agent-browser not found; install its native binary from https://github.com/vercel-labs/agent-browser/releases/latest, then run 'agent-browser install': %w", err)
	}
	return path, nil
}

// Version returns agent-browser's version string.
func (c *Client) Version(ctx context.Context) (string, error) {
	out, _, err := c.run(ctx, nil, "--version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Doctor delegates runtime diagnosis to agent-browser.
func (c *Client) Doctor(ctx context.Context, fix bool) error {
	args := []string{"doctor", "--offline", "--quick", "--json"}
	if fix {
		args = []string{"doctor", "--fix", "--json"}
	}
	_, err := c.runJSON(ctx, nil, args...)
	return err
}

// Open navigates the active tab. Session selection is inherited from the
// environment; Tap never adds or rewrites session flags.
func (c *Client) Open(ctx context.Context, url string) error {
	_, err := c.runJSON(ctx, nil, "open", url, "--json")
	return err
}

// Eval evaluates JavaScript through stdin, avoiding shell escaping and process
// argument leaks for expanded site headers.
func (c *Client) Eval(ctx context.Context, script string) (any, error) {
	data, err := c.runJSON(ctx, []byte(script), "eval", "--stdin", "--json")
	if err != nil {
		return nil, err
	}
	return decodeEval(data)
}

// OpenAndEval applies origin-scoped headers before the first navigation and
// evaluates a site program in one stdin-driven batch. Headers and JavaScript
// never enter process arguments. Headers introduced by Tap are cleared in a
// separate best-effort-safe command after the fail-fast batch, so a navigation
// failure cannot evaluate the script against an unrelated previous page.
func (c *Client) OpenAndEval(ctx context.Context, url, script string, headers map[string]string) (any, error) {
	open := []string{"open", url}
	if len(headers) > 0 {
		headersJSON, err := json.Marshal(headers)
		if err != nil {
			return nil, fmt.Errorf("encode site headers: %w", err)
		}
		open = append(open, "--headers", string(headersJSON))
	}
	commands := [][]string{
		open,
		{"eval", "--base64", base64.StdEncoding.EncodeToString([]byte(script))},
	}
	input, err := json.Marshal(commands)
	if err != nil {
		return nil, fmt.Errorf("encode agent-browser batch: %w", err)
	}
	out, _, batchErr := c.run(ctx, input, "batch", "--bail", "--json")
	var cleanupErr error
	if len(headers) > 0 {
		// agent-browser exposes no header getter/restore API. Clearing is safer
		// than leaving credentials in a shared session; switch to one-shot scoped
		// headers when agent-browser provides them.
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_, cleanupErr = c.runJSON(cleanupCtx, nil, "set", "headers", "{}", "--json")
	}
	if batchErr != nil || cleanupErr != nil {
		return nil, errors.Join(batchErr, cleanupErr)
	}
	var results []batchResult
	if err := json.Unmarshal(out, &results); err != nil {
		return nil, fmt.Errorf("decode agent-browser batch: %w", err)
	}
	if len(results) != len(commands) {
		return nil, fmt.Errorf("agent-browser batch returned %d results, want %d", len(results), len(commands))
	}
	for index, result := range results {
		if !result.Success {
			return nil, fmt.Errorf("agent-browser batch command %d: %s", index+1, decodeError(result.Error))
		}
	}
	return decodeEval(results[1].Result)
}

func decodeEval(data json.RawMessage) (any, error) {
	var result evalData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode agent-browser eval result: %w", err)
	}
	if len(result.Result) == 0 || bytes.Equal(result.Result, []byte("null")) {
		return nil, nil
	}
	var value any
	if err := json.Unmarshal(result.Result, &value); err != nil {
		return nil, fmt.Errorf("decode agent-browser eval value: %w", err)
	}
	return value, nil
}

// CurrentURL returns the active tab URL.
func (c *Client) CurrentURL(ctx context.Context) (string, error) {
	data, err := c.runJSON(ctx, nil, "get", "url", "--json")
	if err != nil {
		return "", err
	}
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("decode agent-browser URL: %w", err)
	}
	return result.URL, nil
}

// HasActiveSession checks daemon state without launching a browser. It is used
// by `tap fetch` with no URL so that reading the current tab has no hidden side
// effect.
func (c *Client) HasActiveSession(ctx context.Context) (bool, error) {
	data, err := c.runJSON(ctx, nil, "session", "list", "--json")
	if err != nil {
		return false, err
	}
	var result struct {
		Sessions []string `json:"sessions"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return false, fmt.Errorf("decode agent-browser sessions: %w", err)
	}
	want := os.Getenv("AGENT_BROWSER_SESSION")
	if want == "" {
		want = "default"
	}
	for _, session := range result.Sessions {
		if session == want {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) runJSON(ctx context.Context, stdin []byte, args ...string) (json.RawMessage, error) {
	out, stderr, err := c.run(ctx, stdin, args...)
	if err != nil {
		return nil, err
	}
	var response envelope
	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("decode agent-browser response: %w", err)
	}
	if !response.Success {
		message := decodeError(response.Error)
		if message == "" {
			message = strings.TrimSpace(string(stderr))
		}
		if message == "" {
			message = "command failed"
		}
		return nil, errors.New(message)
	}
	return response.Data, nil
}

func (c *Client) run(ctx context.Context, stdin []byte, args ...string) ([]byte, []byte, error) {
	path, err := c.Path()
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.CommandContext(ctx, path, args...)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, stderr.Bytes(), fmt.Errorf("agent-browser %s: %s", strings.Join(args, " "), message)
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

func decodeError(raw json.RawMessage) string {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return ""
	}
	var message string
	if json.Unmarshal(raw, &message) == nil {
		return message
	}
	var value struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(raw, &value) == nil {
		return value.Message
	}
	return strings.TrimSpace(string(raw))
}
