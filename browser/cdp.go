package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

const TargetTypePage = "page"

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
			if ti.Type != TargetTypePage {
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

// FormField describes a fillable form element on the page.
type FormField struct {
	Selector    string `json:"selector"`
	Tag         string `json:"tag"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Value       string `json:"value,omitempty"`
	Label       string `json:"label,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Role        string `json:"role,omitempty"`
}

// FillField pairs a CSS selector with the value to fill.
type FillField struct {
	Selector string
	Value    string
}

// FormsTarget discovers fillable form elements in a browser tab.
func FormsTarget(ctx context.Context, debugURL string, targetID string) ([]FormField, error) {
	js := `
(() => {
  const els = document.querySelectorAll('input, textarea, select, button[type="submit"], button:not([type]), input[type="submit"]');
  return [...els].map(el => {
    const tag = el.tagName.toLowerCase();
    let selector = "";
    if (el.id) {
      selector = "#" + CSS.escape(el.id);
    } else if (el.name) {
      selector = tag + "[name=" + JSON.stringify(el.name) + "]";
      // Disambiguate when multiple elements share the same name (radio groups, checkboxes)
      const dupes = document.querySelectorAll(selector);
      if (dupes.length > 1) {
        if (el.value) {
          selector += "[value=" + JSON.stringify(el.value) + "]";
        } else {
          const idx = [...dupes].indexOf(el);
          if (idx > 0) {
            // nth-of-type won't work here; use a parent-scoped index
            const parent = el.parentElement;
            if (parent) {
              const siblings = [...parent.querySelectorAll(selector)];
              const sIdx = siblings.indexOf(el);
              if (sIdx >= 0) {
                let parentSel = parent.id ? "#" + CSS.escape(parent.id) : parent.tagName.toLowerCase();
                selector = parentSel + " " + selector + ":nth-of-type(" + (sIdx + 1) + ")";
              }
            }
          }
        }
      }
    } else if (el.placeholder) {
      selector = tag + "[placeholder=" + JSON.stringify(el.placeholder) + "]";
    } else if (el.type === "submit" || (tag === "button" && !el.type)) {
      const parent = el.parentElement;
      if (parent) {
        const siblings = [...parent.children].filter(c => c.tagName.toLowerCase() === tag);
        const idx = siblings.indexOf(el);
        if (idx >= 0) {
          // Build a unique path: parent selector + child nth-of-type
          let parentSel = "";
          if (parent.id) {
            parentSel = "#" + CSS.escape(parent.id);
          } else {
            parentSel = parent.tagName.toLowerCase();
            const gp = parent.parentElement;
            if (gp) {
              const pSiblings = [...gp.children].filter(c => c.tagName.toLowerCase() === parent.tagName.toLowerCase());
              const pIdx = pSiblings.indexOf(parent);
              if (pIdx >= 0) parentSel += ":nth-of-type(" + (pIdx + 1) + ")";
            }
          }
          selector = parentSel + " > " + tag + ":nth-of-type(" + (idx + 1) + ")";
        }
      }
    }

    let label = "";
    if (el.id) {
      const lbl = document.querySelector("label[for=" + JSON.stringify(el.id) + "]");
      if (lbl) label = lbl.textContent.trim();
    }
    if (!label && el.closest("label")) {
      label = el.closest("label").textContent.trim();
    }
    if (!label && el.getAttribute("aria-label")) {
      label = el.getAttribute("aria-label");
    }

    let role = "";
    if (tag === "button" || el.type === "submit") role = "submit";
    else if (tag === "select") role = "select";
    else if (tag === "textarea") role = "text";
    else if (["text","email","password","search","tel","url","number"].includes(el.type)) role = "text";
    else if (["checkbox","radio"].includes(el.type)) role = "toggle";
    else if (el.type === "hidden") role = "hidden";
    else if (el.type === "file") role = "file";

    return {
      selector: selector,
      tag: tag,
      type: el.type || "",
      name: el.name || "",
      id: el.id || "",
      placeholder: el.placeholder || "",
      value: tag === "select" ? el.value || "" : el.value || "",
      label: label,
      required: el.required || false,
      disabled: el.disabled || false,
      role: role,
    };
  }).filter(f => f.role !== "hidden" && f.selector !== "");
})()
`
	var fields []FormField
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &fields, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("forms target: %w", err)
	}
	return fields, nil
}

// FillTarget fills form fields in a browser tab using React-compatible value setting.
func FillTarget(ctx context.Context, debugURL string, targetID string, fields []FillField, submitSelector string) error {
	var actions []chromedp.Action

	for _, f := range fields {
		sel := f.Selector
		val := f.Value
		// Use a JS snippet that sets value via native setter and dispatches events.
		// This works with React, Vue, Angular, and vanilla HTML forms.
		js := fmt.Sprintf(`
(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  el.focus();
  const tag = el.tagName.toLowerCase();
  if (tag === "select") {
    el.value = %q;
    el.dispatchEvent(new Event("change", { bubbles: true }));
  } else if (el.type === "checkbox" || el.type === "radio") {
    const want = %q;
    if (want === "true" || want === "1" || want === "on") {
      if (!el.checked) el.click();
    } else {
      if (el.checked) el.click();
    }
  } else {
    const setter = (tag === "textarea")
      ? Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value")?.set
      : Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    if (setter) {
      setter.call(el, %q);
    } else {
      el.value = %q;
    }
    el.dispatchEvent(new Event("input", { bubbles: true }));
    el.dispatchEvent(new Event("change", { bubbles: true }));
  }
  return true;
})()
`, sel, sel, val, val, val, val)

		var ok bool
		actions = append(actions, chromedp.Evaluate(js, &ok, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}))
	}

	if submitSelector != "" {
		submitJS := fmt.Sprintf(`
(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("submit element not found: " + %q);
  el.click();
  return true;
})()
`, submitSelector, submitSelector)
		var ok bool
		actions = append(actions, chromedp.Evaluate(submitJS, &ok, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}))
	}

	if err := withTarget(ctx, debugURL, targetID, actions...); err != nil {
		return fmt.Errorf("fill target: %w", err)
	}
	return nil
}

// withTargetListen connects to debugURL, attaches to the specific target, and
// returns a long-lived context suitable for event listening. Unlike withTarget,
// it does NOT run actions or detach automatically — the caller controls the
// session lifetime.
//
// Usage:
//  1. Call withTargetListen to get taskCtx.
//  2. Register event listeners with chromedp.ListenTarget(taskCtx, ...).
//  3. Enable the domain with chromedp.Run(taskCtx, ...).
//  4. Wait/process events.
//  5. Call cancel — clears TargetID first (detach-without-close).
func withTargetListen(ctx context.Context, debugURL string, targetID string) (context.Context, func(), error) {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, debugURL, chromedp.NoModifyURL)
	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(target.ID(targetID)))

	// Run an empty action to force the CDP session to attach to the target.
	if err := chromedp.Run(taskCtx); err != nil {
		taskCancel()
		allocCancel()
		return nil, nil, fmt.Errorf("attach target for listen: %w", err)
	}

	cancel := func() {
		// Clear TargetID BEFORE cancel so chromedp's cancel handler does not
		// close the tab. Same detach-without-close trick as withTarget.
		if c := chromedp.FromContext(taskCtx); c != nil && c.Target != nil {
			c.Target.TargetID = ""
		}
		taskCancel()
		allocCancel()
	}

	return taskCtx, cancel, nil
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
