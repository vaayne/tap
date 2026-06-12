// Package transport provides a shared network layer for fetching web content.
// It supports two levels: direct HTTP and browser-based (CDP).
package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/vaayne/tap/browser"
)

// BrowserType identifies which browser backend to use.
type BrowserType string

const (
	// BrowserChrome uses the system Chrome/Chromium (default).
	BrowserChrome BrowserType = "chrome"
	// BrowserLightpanda uses the Lightpanda headless browser.
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
	// Browser selects the browser backend (default: "chrome").
	Browser BrowserType
}

// Transport provides shared HTTP and browser-based network access.
type Transport struct {
	config     Config
	http       *http.Client
	lightpanda *browser.Lightpanda

	// lpAllocCtx is a shared CDP allocator context for Lightpanda.
	// Created once at startup, reused for each browser context.
	lpAllocCtx    context.Context
	lpAllocCancel context.CancelFunc
}

// New creates a new Transport with the given config.
// If the Lightpanda browser backend is selected, it downloads (if needed)
// and starts the Lightpanda server eagerly so errors surface immediately.
func New(ctx context.Context, config Config) (*Transport, error) {
	if config.WSURL != "" {
		resolvedURL, err := browser.ResolveDebugURL(ctx, config.WSURL)
		if err != nil {
			return nil, fmt.Errorf("resolve debug endpoint: %w", err)
		}
		config.WSURL = resolvedURL
	}

	t := &Transport{
		config: config,
		http:   newHTTPClient(),
	}

	if config.Browser == BrowserLightpanda && config.WSURL == "" {
		lp := browser.NewLightpanda("", "")
		if err := lp.Start(ctx); err != nil {
			return nil, fmt.Errorf("start lightpanda: %w", err)
		}
		t.lightpanda = lp

		// Create a shared allocator for the Lightpanda lifetime.
		allocCtx, allocCancel := chromedp.NewRemoteAllocator(
			context.Background(), lp.WSURL(), chromedp.NoModifyURL,
		)
		t.lpAllocCtx = allocCtx
		t.lpAllocCancel = allocCancel
	}

	return t, nil
}

// Close releases resources held by the transport.
func (t *Transport) Close() error {
	if t.lpAllocCancel != nil {
		t.lpAllocCancel()
	}
	if t.lightpanda != nil {
		t.lightpanda.Stop()
	}
	return nil
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
	bctx, cancel, err := t.newBrowserContext(ctx)
	if err != nil {
		return "", err
	}
	defer cancel()

	if err := chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		return "", fmt.Errorf("browse html: %w", err)
	}

	if pauseFn != nil {
		if err := pauseFn(bctx); err != nil {
			return "", fmt.Errorf("pause: %w", err)
		}
	}

	var html string
	if err := chromedp.Run(bctx,
		chromedp.OuterHTML("html", &html),
	); err != nil {
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
	bctx, cancel, err := t.newBrowserContext(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	// Preserve the native fetch before page scripts can override it.
	// Some sites (e.g. GitHub) replace window.fetch with a custom
	// implementation that blocks cross-origin requests.
	preserveNativeFetch := `window.__nativeFetch = window.fetch.bind(window);`

	// Wrap the user script so that `fetch` resolves to the preserved native version.
	wrappedJS := fmt.Sprintf(
		`(function(){ const fetch = window.__nativeFetch || window.fetch; return %s; })()`,
		js,
	)

	actions := []chromedp.Action{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(preserveNativeFetch).Do(ctx)
			return err
		}),
	}
	if len(headers) > 0 {
		nh := make(network.Headers, len(headers))
		for k, v := range headers {
			nh[k] = v
		}
		actions = append(actions, network.SetExtraHTTPHeaders(nh))
	}
	actions = append(actions,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	)

	if err := chromedp.Run(bctx, actions...); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}

	if pauseFn != nil {
		if err := pauseFn(bctx); err != nil {
			return nil, fmt.Errorf("pause: %w", err)
		}
	}

	var result any
	if err := chromedp.Run(bctx,
		chromedp.Evaluate(wrappedJS, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}

	return result, nil
}

// BrowseInteractive navigates to a URL and keeps the browser open until
// pauseFn returns. This is used by the "login" command to let users interact
// with a site (login, solve CAPTCHAs) while cookies are persisted in the
// Chrome profile directory.
func (t *Transport) BrowseInteractive(ctx context.Context, url string, pauseFn PauseFunc) error {
	bctx, cancel, err := t.newBrowserContext(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	if err := chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		return fmt.Errorf("browse interactive: %w", err)
	}

	if pauseFn != nil {
		if err := pauseFn(bctx); err != nil {
			return fmt.Errorf("pause: %w", err)
		}
	}

	return nil
}

func (t *Transport) newBrowserContext(parent context.Context) (context.Context, context.CancelFunc, error) {
	// Remote CDP endpoint (explicit --ws-url or resolved from an HTTP DevTools base URL).
	if t.config.WSURL != "" {
		ctx, cancel1 := chromedp.NewRemoteAllocator(parent, t.config.WSURL, chromedp.NoModifyURL)
		ctx, cancel2 := chromedp.NewContext(ctx)
		return ctx, func() { cancel2(); cancel1() }, nil
	}

	// Lightpanda browser backend.
	if t.config.Browser == BrowserLightpanda {
		ctx, cancel := t.newLightpandaContext(parent)
		return ctx, cancel, nil
	}

	// Default: local Chrome.
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", t.config.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-web-security", true),
	)

	profileDir := t.config.ProfileDir
	if profileDir == "" {
		profileDir = defaultProfileDir()
	}
	opts = append(opts, chromedp.UserDataDir(profileDir))

	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create browser profile dir: %w", err)
	}
	if err := browser.PrepareProfileDir(profileDir); err != nil {
		return nil, nil, err
	}

	ctx, cancel1 := chromedp.NewExecAllocator(parent, opts...)
	ctx, cancel2 := chromedp.NewContext(ctx)
	return ctx, func() { cancel2(); cancel1() }, nil
}

func (t *Transport) newLightpandaContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel1 := chromedp.NewContext(t.lpAllocCtx)

	// Respect the caller's context deadline/cancellation.
	ctx, cancel2 := context.WithCancel(ctx)
	go func() {
		select {
		case <-parent.Done():
			cancel2()
		case <-ctx.Done():
		}
	}()

	return ctx, func() { cancel2(); cancel1() }
}

func defaultProfileDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "tap", "chrome-profile-"+os.Getenv("USER"))
}
