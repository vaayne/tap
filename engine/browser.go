package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/vaayne/tap/script"
)

// BrowserConfig holds configuration for the CDP browser engine.
type BrowserConfig struct {
	// WSURL is the remote CDP WebSocket URL. If empty, a local Chrome is launched.
	WSURL string
	// ProfileDir is the Chrome user data directory for persistent cookies/storage.
	// Defaults to ~/.cache/tap/chrome-profile-$USER.
	ProfileDir string
}

// Browser executes scripts in a real Chrome browser via CDP.
type Browser struct {
	config BrowserConfig
}

// NewBrowser creates a new Browser engine with the given config.
func NewBrowser(config BrowserConfig) *Browser {
	return &Browser{config: config}
}

func (b *Browser) Name() string { return "Browser" }
func (b *Browser) Close() error { return nil }

func (b *Browser) Run(ctx context.Context, s *script.Script, args map[string]string) (any, error) {
	bctx, cancel := b.newContext(ctx)
	defer cancel()

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	js := fmt.Sprintf("(%s)(%s)", s.Body, string(argsJSON))

	navURL := "about:blank"
	if s.Meta.Domain != "" {
		navURL = "https://" + s.Meta.Domain
	}

	var result any
	if err := chromedp.Run(bctx,
		chromedp.Navigate(navURL),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp run: %w", err)
	}

	return result, nil
}

func (b *Browser) newContext(parent context.Context) (context.Context, context.CancelFunc) {
	if b.config.WSURL != "" {
		ctx, cancel1 := chromedp.NewRemoteAllocator(parent, b.config.WSURL, chromedp.NoModifyURL)
		ctx, cancel2 := chromedp.NewContext(ctx)
		return ctx, func() { cancel2(); cancel1() }
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)

	profileDir := b.config.ProfileDir
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
