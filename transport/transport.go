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

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Config holds transport configuration.
type Config struct {
	// WSURL is the remote CDP WebSocket URL. If empty, a local Chrome is launched.
	WSURL string
	// ProfileDir is the Chrome user data directory for persistent cookies/storage.
	ProfileDir string
	// Headless controls whether Chrome runs in headless mode (default: true).
	Headless bool
}

// Transport provides shared HTTP and browser-based network access.
type Transport struct {
	config Config
	http   *http.Client
}

// New creates a new Transport with the given config.
func New(config Config) *Transport {
	return &Transport{
		config: config,
		http:   &http.Client{},
	}
}

// Close releases resources held by the transport.
func (t *Transport) Close() error {
	return nil
}

// GetHTML fetches a URL via direct HTTP and returns the response body as a string.
func (t *Transport) GetHTML(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

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

// BrowseHTML navigates to a URL in a browser and returns the rendered HTML.
func (t *Transport) BrowseHTML(ctx context.Context, url string) (string, error) {
	bctx, cancel := t.newBrowserContext(ctx)
	defer cancel()

	var html string
	if err := chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &html),
	); err != nil {
		return "", fmt.Errorf("browse html: %w", err)
	}

	return html, nil
}

// BrowseEval navigates to a URL in a browser and evaluates JavaScript.
func (t *Transport) BrowseEval(ctx context.Context, url string, js string) (any, error) {
	bctx, cancel := t.newBrowserContext(ctx)
	defer cancel()

	var result any
	if err := chromedp.Run(bctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return nil, fmt.Errorf("browse eval: %w", err)
	}

	return result, nil
}

func (t *Transport) newBrowserContext(parent context.Context) (context.Context, context.CancelFunc) {
	if t.config.WSURL != "" {
		ctx, cancel1 := chromedp.NewRemoteAllocator(parent, t.config.WSURL, chromedp.NoModifyURL)
		ctx, cancel2 := chromedp.NewContext(ctx)
		return ctx, func() { cancel2(); cancel1() }
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", t.config.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)

	profileDir := t.config.ProfileDir
	if profileDir == "" {
		profileDir = defaultProfileDir()
	}
	opts = append(opts, chromedp.UserDataDir(profileDir))

	ctx, cancel1 := chromedp.NewExecAllocator(parent, opts...)
	ctx, cancel2 := chromedp.NewContext(ctx)
	return ctx, func() { cancel2(); cancel1() }
}

func defaultProfileDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "tap", "chrome-profile-"+os.Getenv("USER"))
}
