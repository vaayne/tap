// Package transport provides a shared network layer for fetching web content.
// It supports two levels: direct HTTP and browser-based (agent-browser).
package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vaayne/tap/browser"
)

// BrowserType identifies which browser backend to use.
// Deprecated: agent-browser is the only backend.
type BrowserType string

const (
	// BrowserChrome uses the system Chrome/Chromium (default).
	// Deprecated: agent-browser manages Chrome internally.
	BrowserChrome BrowserType = "chrome"
	// BrowserLightpanda uses the Lightpanda headless browser.
	// Deprecated: not supported with agent-browser.
	BrowserLightpanda BrowserType = "lightpanda"
)

// Config holds transport configuration.
type Config struct {
	// WSURL is the remote CDP WebSocket URL. If empty, a local browser is launched.
	WSURL string
	// ProfileDir is the Chrome user data directory for persistent cookies/storage.
	ProfileDir string
	// Headless controls whether Chrome runs in headless mode (default: true).
	Headless bool
	// Browser selects the agent-browser engine (default: "chrome").
	Browser BrowserType
}

// Transport provides shared HTTP and browser-based network access.
type Transport struct {
	config       Config
	http         *http.Client
	agentBrowser *browser.AgentBrowser
}

// New creates a new Transport with the given config.
func New(ctx context.Context, config Config) (*Transport, error) {
	ab, err := browser.NewAgentBrowser("")
	if err != nil {
		return nil, fmt.Errorf("agent-browser: %w", err)
	}
	ab.ProfileDir = config.ProfileDir
	ab.Headed = !config.Headless
	if config.Browser == BrowserLightpanda {
		ab.Engine = "lightpanda"
	}

	t := &Transport{
		config:       config,
		http:         newHTTPClient(),
		agentBrowser: ab,
	}

	return t, nil
}

// Close releases resources held by the transport.
func (t *Transport) Close() error {
	return nil
}

// AgentBrowser returns the underlying agent-browser adapter.
func (t *Transport) AgentBrowser() *browser.AgentBrowser {
	return t.agentBrowser
}

// GetHTML fetches a URL via direct HTTP and returns the response body as a string.
func (t *Transport) GetHTML(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}

// Do executes an HTTP request and returns the response.
// Caller is responsible for closing the response body.
func (t *Transport) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return t.http.Do(req.WithContext(ctx))
}

// PauseFunc is called after navigation to let the user interact with the
// browser (e.g. login, solve a CAPTCHA). It should block until the user is
// done. The context is cancelled if the parent context is cancelled.
type PauseFunc func(ctx context.Context) error

// BrowseHTML navigates to a URL in a browser and returns the rendered HTML.
func (t *Transport) BrowseHTML(ctx context.Context, url string) (string, error) {
	return t.BrowseHTMLWithPause(ctx, url, nil)
}

// BrowseHTMLWithPause is like BrowseHTML but calls pauseFn after navigation.
func (t *Transport) BrowseHTMLWithPause(ctx context.Context, url string, pauseFn PauseFunc) (string, error) {
	if err := t.agentBrowser.Open(ctx, url, browser.OpenOpts{}); err != nil {
		return "", fmt.Errorf("browse html: %w", err)
	}
	if pauseFn != nil {
		if err := pauseFn(ctx); err != nil {
			return "", fmt.Errorf("pause: %w", err)
		}
	}
	html, err := t.agentBrowser.GetHTML(ctx)
	if err != nil {
		return "", fmt.Errorf("browse html: %w", err)
	}
	return html, nil
}

// BrowseEval navigates to a URL in a browser and evaluates JavaScript.
func (t *Transport) BrowseEval(ctx context.Context, url string, js string, headers map[string]string) (any, error) {
	return t.BrowseEvalWithPause(ctx, url, js, nil, headers)
}

// BrowseEvalWithPause is like BrowseEval but calls pauseFn after navigation.
func (t *Transport) BrowseEvalWithPause(ctx context.Context, url string, js string, pauseFn PauseFunc, headers map[string]string) (any, error) {
	preserveNativeFetch := `window.__nativeFetch = window.fetch.bind(window);`

	tmpfile, err := os.CreateTemp("", "tap-init-*.js")
	if err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	if _, err := tmpfile.WriteString(preserveNativeFetch); err != nil {
		_ = tmpfile.Close()
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	if err := tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}

	wrappedJS := fmt.Sprintf(
		`(function(){ const fetch = window.__nativeFetch || window.fetch; return %s; })()`,
		js,
	)

	openOpts := browser.OpenOpts{InitScript: tmpfile.Name()}
	if len(headers) > 0 {
		openOpts.Headers = headers
	}
	if err := t.agentBrowser.Open(ctx, url, openOpts); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	if pauseFn != nil {
		if err := pauseFn(ctx); err != nil {
			return nil, fmt.Errorf("pause: %w", err)
		}
	}
	result, err := t.agentBrowser.Eval(ctx, wrappedJS)
	if err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}
	return result, nil
}

// BrowseInteractive navigates to a URL and keeps the browser open until
// pauseFn returns. This is used by the "login" command to let users interact
// with a site (login, solve CAPTCHAs) while cookies are persisted in the
// Chrome profile directory.
func (t *Transport) BrowseInteractive(ctx context.Context, url string, pauseFn PauseFunc) error {
	if err := t.agentBrowser.Open(ctx, url, browser.OpenOpts{Headed: true}); err != nil {
		return fmt.Errorf("browse interactive: %w", err)
	}
	if pauseFn != nil {
		if err := pauseFn(ctx); err != nil {
			return fmt.Errorf("pause: %w", err)
		}
	}
	return nil
}
