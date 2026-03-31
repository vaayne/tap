package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
)

// newBrowserContext creates a chromedp context.
// If CDP_WS_URL is set, connects to the remote browser.
// Otherwise, launches a local headless Chrome with optional profile persistence via CDP_PROFILE_DIR.
func newBrowserContext() (context.Context, context.CancelFunc) {
	wsURL := os.Getenv("CDP_WS_URL")

	if wsURL != "" {
		ctx, cancel1 := chromedp.NewRemoteAllocator(context.Background(), wsURL, chromedp.NoModifyURL)
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

	profileDir := os.Getenv("CDP_PROFILE_DIR")
	if profileDir == "" {
		home, _ := os.UserHomeDir()
		profileDir = filepath.Join(home, ".cache", "cdp", "chrome-profile-"+os.Getenv("USER"))
	}
	opts = append(opts, chromedp.UserDataDir(profileDir))

	ctx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel2 := chromedp.NewContext(ctx)
	return ctx, func() { cancel2(); cancel1() }
}
