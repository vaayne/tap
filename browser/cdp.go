package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
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

// GetHTMLTarget returns the outerHTML of the element matching sel, or the full
// page HTML if sel is empty.
func GetHTMLTarget(ctx context.Context, debugURL string, targetID string, sel string) (string, string, error) {
	if sel == "" {
		sel = "html"
	}
	js := fmt.Sprintf(`
(() => {
  const el = document.querySelector(%q);
  if (!el) return {html: "", url: location.href};
  return {html: el.outerHTML, url: location.href};
})()
`, sel)
	var result struct {
		HTML string `json:"html"`
		URL  string `json:"url"`
	}
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", "", fmt.Errorf("get html target: %w", err)
	}
	return result.HTML, result.URL, nil
}

// KeypressTarget sends key events to the page (not a specific element).
// The keys string uses chromedp/kb constants: "\r" for Enter, "\t" for Tab,
// "\u001b" for Escape, etc. Regular characters are sent as-is.
func KeypressTarget(ctx context.Context, debugURL string, targetID string, keys string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.KeyEvent(keys),
	); err != nil {
		return fmt.Errorf("keypress target: %w", err)
	}
	return nil
}

// DialogTarget accepts or dismisses a pending JavaScript dialog (alert/confirm/prompt).
// For prompt dialogs, promptText is entered before accepting.
func DialogTarget(ctx context.Context, debugURL string, targetID string, accept bool, promptText string) error {
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			p := page.HandleJavaScriptDialog(accept)
			if promptText != "" {
				p = p.WithPromptText(promptText)
			}
			return p.Do(ctx)
		}),
	)
	if err != nil {
		return fmt.Errorf("dialog target: %w", err)
	}
	return nil
}

// CookieEntry represents a browser cookie for JSON output.
type CookieEntry struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	SameSite string  `json:"sameSite,omitempty"`
}

// GetCookiesTarget returns all cookies for the current page.
func GetCookiesTarget(ctx context.Context, debugURL string, targetID string) ([]CookieEntry, error) {
	var cookies []CookieEntry
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			raw, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, c := range raw {
				cookies = append(cookies, CookieEntry{
					Name:     c.Name,
					Value:    c.Value,
					Domain:   c.Domain,
					Path:     c.Path,
					Expires:  c.Expires,
					HTTPOnly: c.HTTPOnly,
					Secure:   c.Secure,
					Session:  c.Session,
					SameSite: string(c.SameSite),
				})
			}
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("get cookies target: %w", err)
	}
	return cookies, nil
}

// SetCookieTarget sets a cookie on the page.
func SetCookieTarget(ctx context.Context, debugURL string, targetID string, name, value, domain, path string) error {
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			p := network.SetCookie(name, value)
			if domain != "" {
				p = p.WithDomain(domain)
			}
			if path != "" {
				p = p.WithPath(path)
			}
			return p.Do(ctx)
		}),
	)
	if err != nil {
		return fmt.Errorf("set cookie target: %w", err)
	}
	return nil
}

// ClearCookiesTarget deletes all cookies for the current page.
func ClearCookiesTarget(ctx context.Context, debugURL string, targetID string) error {
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			raw, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, c := range raw {
				if err := network.DeleteCookies(c.Name).WithDomain(c.Domain).WithPath(c.Path).Do(ctx); err != nil {
					return fmt.Errorf("delete cookie %q: %w", c.Name, err)
				}
			}
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("clear cookies target: %w", err)
	}
	return nil
}

// PDFTarget saves the current page as PDF and returns the bytes.
func PDFTarget(ctx context.Context, debugURL string, targetID string, landscape bool, printBackground bool, scale float64) ([]byte, error) {
	var buf []byte
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			params := page.PrintToPDF().
				WithLandscape(landscape).
				WithPrintBackground(printBackground)
			if scale > 0 {
				params = params.WithScale(scale)
			}
			data, _, err := params.Do(ctx)
			if err != nil {
				return err
			}
			buf = data
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("pdf target: %w", err)
	}
	return buf, nil
}

// ClickTarget dispatches a real mouse click on the first element matching sel.
// chromedp.Click sends the full mouseMoved → mousePressed → mouseReleased sequence.
func ClickTarget(ctx context.Context, debugURL string, targetID string, sel string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Click(sel, chromedp.ByQuery, chromedp.NodeVisible),
	); err != nil {
		return fmt.Errorf("click target: %w", err)
	}
	return nil
}

// TypeTarget sends individual key events to the element matching sel.
// Unlike FillTarget which sets .value directly, this dispatches keyDown/keyUp
// per character — behaving like a real user typing.
func TypeTarget(ctx context.Context, debugURL string, targetID string, sel string, text string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Click(sel, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.SendKeys(sel, text, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("type target: %w", err)
	}
	return nil
}

// HoverTarget moves the mouse to the center of the element matching sel,
// dispatching mouseMoved events that trigger CSS :hover states and mouseenter listeners.
func HoverTarget(ctx context.Context, debugURL string, targetID string, sel string) error {
	var nodes []*cdp.Node
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Nodes(sel, &nodes, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(nodes) == 0 {
				return fmt.Errorf("no element matching %q", sel)
			}
			box, err := dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
			if err != nil {
				return fmt.Errorf("get box model: %w", err)
			}
			c := box.Content
			x := (c[0] + c[2] + c[4] + c[6]) / 4
			y := (c[1] + c[3] + c[5] + c[7]) / 4
			return chromedp.MouseEvent(input.MouseMoved, x, y).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("hover target: %w", err)
	}
	return nil
}

// ScrollTarget scrolls the element matching sel into view. If sel is empty,
// scrolls to the given x,y pixel coordinates.
func ScrollTarget(ctx context.Context, debugURL string, targetID string, sel string, x, y float64) error {
	var actions []chromedp.Action
	if sel != "" {
		actions = append(actions, chromedp.ScrollIntoView(sel, chromedp.ByQuery))
	} else {
		js := fmt.Sprintf("window.scrollTo(%f, %f)", x, y)
		var ignore any
		actions = append(actions, chromedp.Evaluate(js, &ignore))
	}
	if err := withTarget(ctx, debugURL, targetID, actions...); err != nil {
		return fmt.Errorf("scroll target: %w", err)
	}
	return nil
}

// SelectTarget selects an option in a <select> element by value, dispatching
// focus, input, and change events.
func SelectTarget(ctx context.Context, debugURL string, targetID string, sel string, value string) error {
	js := fmt.Sprintf(`
(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  if (el.tagName.toLowerCase() !== "select") throw new Error("element is not a <select>");
  el.focus();
  el.value = %q;
  el.dispatchEvent(new Event("input", { bubbles: true }));
  el.dispatchEvent(new Event("change", { bubbles: true }));
  return true;
})()
`, sel, sel, value)
	var ok bool
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &ok, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return fmt.Errorf("select target: %w", err)
	}
	return nil
}

// WaitForTarget waits until the element matching sel is visible in the tab.
// The caller controls the deadline via ctx.
func WaitForTarget(ctx context.Context, debugURL string, targetID string, sel string, timeout time.Duration) error {
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	if err := withTarget(waitCtx, debugURL, targetID,
		chromedp.WaitVisible(sel, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("wait for target: %w", err)
	}
	return nil
}

// BackTarget navigates the tab backwards in history.
func BackTarget(ctx context.Context, debugURL string, targetID string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.NavigateBack(),
	); err != nil {
		return fmt.Errorf("back target: %w", err)
	}
	return nil
}

// ForwardTarget navigates the tab forwards in history.
func ForwardTarget(ctx context.Context, debugURL string, targetID string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.NavigateForward(),
	); err != nil {
		return fmt.Errorf("forward target: %w", err)
	}
	return nil
}

// ReloadTarget reloads the current page in the tab.
func ReloadTarget(ctx context.Context, debugURL string, targetID string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Reload(),
	); err != nil {
		return fmt.Errorf("reload target: %w", err)
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
