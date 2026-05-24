package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/vaayne/tap/browser"
	"github.com/vaayne/tap/script"
	"github.com/vaayne/tap/transport"
)

// Browser executes scripts in a real browser via agent-browser.
type Browser struct {
	agentBrowser *browser.AgentBrowser
	pauseFn      transport.PauseFunc
}

// NewBrowser creates a new Browser engine backed by agent-browser.
func NewBrowser(ab *browser.AgentBrowser, pauseFn transport.PauseFunc) *Browser {
	return &Browser{agentBrowser: ab, pauseFn: pauseFn}
}

func (b *Browser) Name() string { return "Browser" }
func (b *Browser) Close() error  { return nil }

func (b *Browser) Run(ctx context.Context, s *script.Script, args map[string]string, opts RunOpts) (any, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	js := fmt.Sprintf("(%s)(%s)", s.Body, string(argsJSON))

	navURL := "about:blank"
	if s.Meta.Domain != "" {
		navURL = "https://" + s.Meta.Domain
	}

	// Preserve the native fetch before page scripts can override it.
	preserveNativeFetch := `window.__nativeFetch = window.fetch.bind(window);`
	tmpfile, err := os.CreateTemp("", "tap-init-*.js")
	if err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	if _, err := tmpfile.WriteString(preserveNativeFetch); err != nil {
		_ = tmpfile.Close()
		return nil, fmt.Errorf("write init script: %w", err)
	}
	if err := tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("close init script: %w", err)
	}

	openOpts := browser.OpenOpts{InitScript: tmpfile.Name()}
	if len(opts.Headers) > 0 {
		openOpts.Headers = opts.Headers
	}
	if err := b.agentBrowser.Open(ctx, navURL, openOpts); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	if b.pauseFn != nil {
		if err := b.pauseFn(ctx); err != nil {
			return nil, fmt.Errorf("pause: %w", err)
		}
	}
	wrappedJS := fmt.Sprintf(
		`(function(){ const fetch = window.__nativeFetch || window.fetch; return %s; })()`,
		js,
	)
	return b.agentBrowser.Eval(ctx, wrappedJS)
}
