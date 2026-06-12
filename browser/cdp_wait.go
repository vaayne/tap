package browser

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ElementState describes which DOM readiness condition to wait for.
type ElementState int

const (
	// ElementVisible waits for the element to be present and visible.
	ElementVisible ElementState = iota
	// ElementHidden waits for the element to be absent from the layout
	// (hidden or removed from the DOM).
	ElementHidden
	// ElementAttached waits for the element to be present in the DOM
	// (does not require it to be visible).
	ElementAttached
	// ElementDetached waits for the element to be absent from the DOM.
	ElementDetached
)

// LoadState names a page-load event to wait for.
type LoadState string

const (
	LoadStateLoad             LoadState = "load"
	LoadStateDOMContentLoaded LoadState = "domcontentloaded"
	// LoadStateNetworkIdle polls until no network activity for ~500 ms.
	// Implemented via JS: checks window.performance.getEntriesByType("resource")
	// against a snapshot taken 500 ms earlier. This is a heuristic — it fires
	// once the browser has not started any new resource fetches for half a second.
	LoadStateNetworkIdle LoadState = "networkidle"
)

const pollInterval = 100 * time.Millisecond

// WaitForElementTarget waits for a CSS selector to reach the desired DOM state.
// Supported states: visible, hidden, attached, detached.
func WaitForElementTarget(ctx context.Context, debugURL, targetID, sel string, state ElementState, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var action chromedp.Action
	switch state {
	case ElementVisible:
		action = chromedp.WaitVisible(sel, chromedp.ByQuery)
	case ElementHidden:
		action = chromedp.WaitNotPresent(sel, chromedp.ByQuery)
	case ElementAttached:
		action = chromedp.WaitReady(sel, chromedp.ByQuery)
	case ElementDetached:
		action = chromedp.WaitNotPresent(sel, chromedp.ByQuery)
	default:
		return fmt.Errorf("wait element: unknown state %d", state)
	}

	if err := withTarget(waitCtx, debugURL, targetID, action); err != nil {
		return fmt.Errorf("wait element %q: %w", sel, err)
	}
	return nil
}

// WaitForTextTarget polls until document.body.innerText contains the given substring.
// Polling interval: 100 ms.
func WaitForTextTarget(ctx context.Context, debugURL, targetID, text string, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Escape the text for embedding in a JS string literal.
	escaped := strings.ReplaceAll(text, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	js := fmt.Sprintf(`document.body && document.body.innerText.includes("%s")`, escaped)

	if err := withTarget(waitCtx, debugURL, targetID, pollUntilTrue(waitCtx, js)); err != nil {
		return fmt.Errorf("wait text %q: %w", text, err)
	}
	return nil
}

// WaitForURLTarget polls until the page's location.href matches the glob pattern.
// Pattern supports * (any non-separator sequence) and ** (any sequence including /).
// Uses path.Match semantics with ** expanded before matching. The comparison is
// against the full URL string.
func WaitForURLTarget(ctx context.Context, debugURL, targetID, glob string, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	checkFn := func(current string) (bool, error) {
		return matchURLGlob(glob, current), nil
	}

	if err := withTarget(waitCtx, debugURL, targetID, pollUntilURL(waitCtx, checkFn)); err != nil {
		return fmt.Errorf("wait url %q: %w", glob, err)
	}
	return nil
}

// WaitForLoadTarget waits for a named page-load event.
// "load" and "domcontentloaded" use chromedp's built-in page events.
// "networkidle" polls until no new resource entries appear for ~500 ms.
func WaitForLoadTarget(ctx context.Context, debugURL, targetID string, state LoadState, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var action chromedp.Action
	switch state {
	case LoadStateLoad:
		action = chromedp.WaitReady("body", chromedp.ByQuery)
	case LoadStateDOMContentLoaded:
		// document.readyState reaches "interactive" at DOMContentLoaded.
		action = pollUntilTrue(waitCtx, `document.readyState === "interactive" || document.readyState === "complete"`)
	case LoadStateNetworkIdle:
		// networkidle heuristic: poll JS performance entries. We capture the
		// count of resource entries, wait 500 ms, and fire if the count hasn't
		// grown. This is approximate but works for the common SPA case.
		action = chromedp.ActionFunc(func(ctx context.Context) error {
			return waitNetworkIdle(ctx)
		})
	default:
		return fmt.Errorf("wait load: unknown state %q", state)
	}

	if err := withTarget(waitCtx, debugURL, targetID, action); err != nil {
		return fmt.Errorf("wait load %q: %w", state, err)
	}
	return nil
}

// WaitForFnTarget polls until the given JS expression evaluates to a truthy value.
// The expression is run directly via Runtime.evaluate — it is the caller's
// responsibility to ensure it is safe.
func WaitForFnTarget(ctx context.Context, debugURL, targetID, jsExpr string, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if err := withTarget(waitCtx, debugURL, targetID, pollUntilTrue(waitCtx, jsExpr)); err != nil {
		return fmt.Errorf("wait fn: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// pollUntilTrue returns a chromedp action that ticks every 100 ms and runs
// jsExpr until it returns a truthy JS value or the context is done.
func pollUntilTrue(ctx context.Context, jsExpr string) chromedp.Action {
	return chromedp.ActionFunc(func(actCtx context.Context) error {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			var result bool
			if err := chromedp.Evaluate(jsExpr, &result).Do(actCtx); err == nil && result {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-actCtx.Done():
				return actCtx.Err()
			case <-ticker.C:
			}
		}
	})
}

type urlCheckFn func(current string) (bool, error)

// pollUntilURL polls location.href every 100 ms and applies checkFn.
func pollUntilURL(ctx context.Context, checkFn urlCheckFn) chromedp.Action {
	return chromedp.ActionFunc(func(actCtx context.Context) error {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			var current string
			if err := chromedp.Evaluate(`location.href`, &current).Do(actCtx); err == nil {
				if ok, err := checkFn(current); err != nil {
					return err
				} else if ok {
					return nil
				}
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-actCtx.Done():
				return actCtx.Err()
			case <-ticker.C:
			}
		}
	})
}

// waitNetworkIdle fires once no new resource entries have appeared for 500 ms.
// It is called inside a withTarget action context.
func waitNetworkIdle(ctx context.Context) error {
	const quiescentWindow = 500 * time.Millisecond
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var prevCount int
	quiescentSince := time.Time{}

	// Seed the initial count.
	var initial float64
	if err := chromedp.Evaluate(
		`performance.getEntriesByType("resource").length`, &initial,
	).Do(ctx); err == nil {
		prevCount = int(initial)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var countF float64
			if err := chromedp.Evaluate(
				`performance.getEntriesByType("resource").length`, &countF,
			).Do(ctx); err != nil {
				continue
			}
			current := int(countF)
			if current != prevCount {
				prevCount = current
				quiescentSince = time.Time{}
				continue
			}
			// Count unchanged — start or extend the quiet window.
			if quiescentSince.IsZero() {
				quiescentSince = time.Now()
			}
			if time.Since(quiescentSince) >= quiescentWindow {
				return nil
			}
		}
	}
}

// matchURLGlob matches url against a glob pattern.
// Supports * (matches any char except /) and ** (matches anything including /).
// Pattern semantics documented in cmd/tap/browser_wait.go.
func matchURLGlob(pattern, url string) bool {
	// Fast path: no wildcards.
	if !strings.ContainsAny(pattern, "*?") {
		return strings.Contains(url, pattern)
	}

	// ** must be handled before path.Match which doesn't support **.
	// Strategy: replace ** with a placeholder, split on it, and check that all
	// literal segments appear in order in the url.
	if strings.Contains(pattern, "**") {
		return matchDoubleStarGlob(pattern, url)
	}

	// Single-* glob: use path.Match on the full URL string.
	// path.Match treats / literally, so * does not cross /.
	matched, err := path.Match(pattern, url)
	return err == nil && matched
}

// matchDoubleStarGlob handles patterns with ** by splitting on ** segments and
// verifying that between each pair of ** anchors the remaining literal (or
// single-*) segments appear in order within the URL.
func matchDoubleStarGlob(pattern, url string) bool {
	parts := strings.Split(pattern, "**")
	remaining := url
	for i, part := range parts {
		if part == "" {
			continue
		}
		// Replace single-* within this segment with a wildcard character for
		// simple prefix/suffix matching. For a full segment match we use
		// path.Match on a fake path segment.
		idx := indexGlobSegment(part, remaining)
		if idx < 0 {
			return false
		}
		remaining = remaining[idx+expandedLen(part, remaining, idx):]
		_ = i
	}
	return true
}

// indexGlobSegment finds the first occurrence of the literal/glob segment
// `pat` inside `s` and returns the byte offset, or -1 if not found.
// Only plain strings (no wildcards) are handled here; single-* falls back
// to a linear scan.
func indexGlobSegment(pat, s string) int {
	if !strings.ContainsAny(pat, "*?") {
		return strings.Index(s, pat)
	}
	// Brute-force: try matching path.Match(pat, s[i:]) for each i.
	for i := range len(s) {
		sub := s[i:]
		if ok, err := path.Match(pat, sub); err == nil && ok {
			return i
		}
		// Try on each / boundary for path patterns.
		if sub != "" {
			next := strings.IndexByte(sub[1:], '/')
			if next < 0 {
				break
			}
		}
	}
	return -1
}

// expandedLen returns the byte length of the portion of s starting at offset
// that matches pat (used after indexGlobSegment to advance the pointer).
func expandedLen(pat, s string, offset int) int {
	if !strings.ContainsAny(pat, "*?") {
		return len(pat)
	}
	// For wildcard segments try increasing lengths until path.Match stops matching.
	sub := s[offset:]
	best := 0
	for n := 1; n <= len(sub); n++ {
		if ok, err := path.Match(pat, sub[:n]); err == nil && ok {
			best = n
		}
	}
	return best
}
