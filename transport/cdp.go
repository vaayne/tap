package transport

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

// TargetInfo holds metadata about a CDP target (browser tab).
type TargetInfo struct {
	TargetID string
	Title    string
	URL      string
	Type     string
}

// ListTargets enumerates page targets in a browser reachable at debugURL.
func ListTargets(ctx context.Context, debugURL string) ([]TargetInfo, error) {
	bctx, cancel := withBrowser(ctx, debugURL)
	defer cancel()

	// Ensure the browser connection is established.
	if err := chromedp.Run(bctx); err != nil {
		return nil, fmt.Errorf("list targets: connect: %w", err)
	}

	infos, err := target.GetTargets().Do(bctx)
	if err != nil {
		return nil, fmt.Errorf("list targets: %w", err)
	}

	var out []TargetInfo
	for _, ti := range infos {
		if ti.Type != "page" {
			continue
		}
		out = append(out, TargetInfo{
			TargetID: string(ti.TargetID),
			Title:    ti.Title,
			URL:      ti.URL,
			Type:     ti.Type,
		})
	}
	return out, nil
}

// CreateTarget creates a new browser tab navigated to url and returns its target ID.
func CreateTarget(ctx context.Context, debugURL string, url string) (string, error) {
	bctx, cancel := withBrowser(ctx, debugURL)
	defer cancel()

	// Ensure the browser connection is established.
	if err := chromedp.Run(bctx); err != nil {
		return "", fmt.Errorf("create target: connect: %w", err)
	}

	id, err := target.CreateTarget(url).Do(bctx)
	if err != nil {
		return "", fmt.Errorf("create target: %w", err)
	}
	return string(id), nil
}

// CloseTarget closes the browser tab identified by targetID.
func CloseTarget(ctx context.Context, debugURL string, targetID string) error {
	bctx, cancel := withBrowser(ctx, debugURL)
	defer cancel()

	// Ensure the browser connection is established.
	if err := chromedp.Run(bctx); err != nil {
		return fmt.Errorf("close target: connect: %w", err)
	}

	if err := target.CloseTarget(target.ID(targetID)).Do(bctx); err != nil {
		return fmt.Errorf("close target: %w", err)
	}
	return nil
}

// NavigateTarget navigates an existing browser tab to url and waits for the body to be ready.
func NavigateTarget(ctx context.Context, debugURL string, targetID string, url string) error {
	return withTarget(ctx, debugURL, targetID,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	)
}

// EvalTarget evaluates JavaScript in the context of an existing browser tab
// and returns the result.
func EvalTarget(ctx context.Context, debugURL string, targetID string, js string) (any, error) {
	var result any
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("eval target: %w", err)
	}
	return result, nil
}

// ScreenshotTarget captures a full-page screenshot of an existing browser tab
// and returns the PNG bytes.
func ScreenshotTarget(ctx context.Context, debugURL string, targetID string) ([]byte, error) {
	var buf []byte
	err := withTarget(ctx, debugURL, targetID,
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		return nil, fmt.Errorf("screenshot target: %w", err)
	}
	return buf, nil
}

// withBrowser connects to debugURL at the browser level and returns contexts for CDP commands.
func withBrowser(ctx context.Context, debugURL string) (context.Context, context.CancelFunc) {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, debugURL, chromedp.NoModifyURL)
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	return taskCtx, func() { taskCancel(); allocCancel() }
}

// withTarget connects to debugURL, attaches to the specific target, runs the actions, and cleans up.
func withTarget(ctx context.Context, debugURL string, targetID string, actions ...chromedp.Action) error {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, debugURL, chromedp.NoModifyURL)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(target.ID(targetID)))
	defer taskCancel()

	return chromedp.Run(taskCtx, actions...)
}
