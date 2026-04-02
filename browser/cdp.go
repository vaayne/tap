package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
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

	var out []TargetInfo
	err := chromedp.Run(bctx, chromedp.ActionFunc(func(ctx context.Context) error {
		infos, err := target.GetTargets().Do(ctx)
		if err != nil {
			return err
		}
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
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("list targets: %w", err)
	}
	return out, nil
}

// CreateTarget creates a new browser tab navigated to url and returns its target ID.
func CreateTarget(ctx context.Context, debugURL string, url string) (string, error) {
	bctx, cancel := withBrowser(ctx, debugURL)
	defer cancel()

	var id target.ID
	err := chromedp.Run(bctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		id, err = target.CreateTarget(url).Do(ctx)
		return err
	}))
	if err != nil {
		return "", fmt.Errorf("create target: %w", err)
	}
	return string(id), nil
}

// CloseTarget closes the browser tab identified by targetID.
func CloseTarget(ctx context.Context, debugURL string, targetID string) error {
	bctx, cancel := withBrowser(ctx, debugURL)
	defer cancel()

	// Use the browser-level executor because chromedp's tab-level Target.Execute
	// intercepts and rejects CloseTarget commands.
	return chromedp.Run(bctx, chromedp.ActionFunc(func(ctx context.Context) error {
		// Use bctx (the outer chromedp context) to reach the Browser executor.
		// The inner ctx from ActionFunc is bound to the tab-level Target executor
		// which intercepts CloseTarget, so we need the browser-level one.
		c := chromedp.FromContext(bctx)
		if c == nil || c.Browser == nil {
			return fmt.Errorf("close target: no browser connection")
		}
		browserCtx := cdp.WithExecutor(ctx, c.Browser)
		if err := target.CloseTarget(target.ID(targetID)).Do(browserCtx); err != nil {
			return fmt.Errorf("close target: %w", err)
		}
		return nil
	}))
}

// NavigateTarget navigates an existing browser tab to url and waits for the body to be ready.
func NavigateTarget(ctx context.Context, debugURL string, targetID string, url string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		return fmt.Errorf("navigate target: %w", err)
	}
	return nil
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
// It detaches from the target without closing it so the tab survives across calls.
func withTarget(ctx context.Context, debugURL string, targetID string, actions ...chromedp.Action) error {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, debugURL, chromedp.NoModifyURL)

	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(target.ID(targetID)))

	err := chromedp.Run(taskCtx, actions...)

	// Clear TargetID BEFORE cancel so chromedp's cancel handler does not
	// close the tab. We attach to an existing tab we don't own.
	if c := chromedp.FromContext(taskCtx); c != nil && c.Target != nil {
		c.Target.TargetID = ""
	}
	taskCancel()
	allocCancel()

	return err
}
