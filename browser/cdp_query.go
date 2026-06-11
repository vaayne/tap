package browser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// BoundingBox holds the position and dimensions of an element.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// QueryTextTarget returns the textContent of the first element matching sel.
func QueryTextTarget(ctx context.Context, debugURL, targetID, sel string) (string, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  return el.textContent;
})()`, sel, sel)
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query text target: %w", err)
	}
	return result, nil
}

// QueryTextByBackendNodeID returns the textContent of an element by backend node ID.
func QueryTextByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(`function(){ return this.textContent; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query text by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("query text target: %w", err)
	}
	return result, nil
}

// QueryHTMLTarget returns the innerHTML of the first element matching sel.
func QueryHTMLTarget(ctx context.Context, debugURL, targetID, sel string) (string, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  return el.innerHTML;
})()`, sel, sel)
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query html target: %w", err)
	}
	return result, nil
}

// QueryHTMLByBackendNodeID returns the innerHTML of an element by backend node ID.
func QueryHTMLByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(`function(){ return this.innerHTML; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query html by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("query html target: %w", err)
	}
	return result, nil
}

// QueryValueTarget returns the value property of the first element matching sel.
func QueryValueTarget(ctx context.Context, debugURL, targetID, sel string) (string, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  return el.value !== undefined ? String(el.value) : "";
})()`, sel, sel)
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query value target: %w", err)
	}
	return result, nil
}

// QueryValueByBackendNodeID returns the value property of an element by backend node ID.
func QueryValueByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(`function(){ return this.value !== undefined ? String(this.value) : ""; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query value by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("query value target: %w", err)
	}
	return result, nil
}

// QueryAttrTarget returns the value of attr on the first element matching sel.
func QueryAttrTarget(ctx context.Context, debugURL, targetID, sel, attr string) (string, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  const v = el.getAttribute(%q);
  return v !== null ? v : "";
})()`, sel, sel, attr)
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query attr target: %w", err)
	}
	return result, nil
}

// QueryAttrByBackendNodeID returns the value of attr on an element by backend node ID.
func QueryAttrByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID, attr string) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		fn := fmt.Sprintf(`function(){ const v = this.getAttribute(%q); return v !== null ? v : ""; }`, attr)
		val, ex, err := runtime.CallFunctionOn(fn).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query attr by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("query attr target: %w", err)
	}
	return result, nil
}

// QueryTitleTarget returns the document.title of the current page.
func QueryTitleTarget(ctx context.Context, debugURL, targetID string) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(`document.title`, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query title target: %w", err)
	}
	return result, nil
}

// QueryURLTarget returns the current location.href of the page.
func QueryURLTarget(ctx context.Context, debugURL, targetID string) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(`location.href`, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true)
		}),
	)
	if err != nil {
		return "", fmt.Errorf("query url target: %w", err)
	}
	return result, nil
}

// QueryCountTarget returns the number of elements matching sel.
func QueryCountTarget(ctx context.Context, debugURL, targetID, sel string) (int, error) {
	js := fmt.Sprintf(`document.querySelectorAll(%q).length`, sel)
	var result float64 // JSON numbers unmarshal as float64
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true)
		}),
	)
	if err != nil {
		return 0, fmt.Errorf("query count target: %w", err)
	}
	return int(result), nil
}

// QueryBoxTarget returns the bounding box of the first element matching sel.
func QueryBoxTarget(ctx context.Context, debugURL, targetID, sel string) (*BoundingBox, error) {
	var nodes []*cdp.Node
	var box BoundingBox
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Nodes(sel, &nodes, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(nodes) == 0 {
				return fmt.Errorf("no element matching %q", sel)
			}
			model, err := dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
			if err != nil {
				return fmt.Errorf("get box model: %w", err)
			}
			q := model.Border
			// border quad: TL, TR, BR, BL (x,y pairs)
			box.X = q[0]
			box.Y = q[1]
			box.Width = q[2] - q[0]
			box.Height = q[7] - q[1]
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("query box target: %w", err)
	}
	return &box, nil
}

// QueryBoxByBackendNodeID returns the bounding box of an element by backend node ID.
func QueryBoxByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (*BoundingBox, error) {
	var box BoundingBox
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		model, err := dom.GetBoxModel().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("get box model: %w", err)
		}
		q := model.Border
		box.X = q[0]
		box.Y = q[1]
		box.Width = q[2] - q[0]
		box.Height = q[7] - q[1]
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("query box target: %w", err)
	}
	return &box, nil
}

// QueryStylesTarget returns the computed styles of the first element matching sel
// as a map of property name to value.
func QueryStylesTarget(ctx context.Context, debugURL, targetID, sel string) (map[string]string, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) throw new Error("element not found: " + %q);
  const cs = window.getComputedStyle(el);
  const out = {};
  for (let i = 0; i < cs.length; i++) {
    const prop = cs[i];
    out[prop] = cs.getPropertyValue(prop);
  }
  return out;
})()`, sel, sel)
	var result map[string]string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("query styles target: %w", err)
	}
	return result, nil
}

// QueryStylesByBackendNodeID returns the computed styles of an element by backend node ID.
func QueryStylesByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (map[string]string, error) {
	var result map[string]string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		fn := `function(){
  const cs = window.getComputedStyle(this);
  const out = {};
  for (let i = 0; i < cs.length; i++) {
    const prop = cs[i];
    out[prop] = cs.getPropertyValue(prop);
  }
  return out;
}`
		val, ex, err := runtime.CallFunctionOn(fn).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query styles by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("query styles target: %w", err)
	}
	return result, nil
}

// QueryVisibleTarget returns true if the first element matching sel is visible.
// Visible means: in DOM, not display:none/visibility:hidden, and non-zero dimensions.
func QueryVisibleTarget(ctx context.Context, debugURL, targetID, sel string) (bool, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) return false;
  const cs = window.getComputedStyle(el);
  if (cs.display === "none" || cs.visibility === "hidden" || cs.opacity === "0") return false;
  const r = el.getBoundingClientRect();
  return r.width > 0 && r.height > 0;
})()`, sel)
	var result bool
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return false, fmt.Errorf("query visible target: %w", err)
	}
	return result, nil
}

// QueryVisibleByBackendNodeID returns true if the element identified by backend node ID is visible.
func QueryVisibleByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (bool, error) {
	var result bool
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		fn := `function(){
  const cs = window.getComputedStyle(this);
  if (cs.display === "none" || cs.visibility === "hidden" || cs.opacity === "0") return false;
  const r = this.getBoundingClientRect();
  return r.width > 0 && r.height > 0;
}`
		val, ex, err := runtime.CallFunctionOn(fn).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query visible by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return false, fmt.Errorf("query visible target: %w", err)
	}
	return result, nil
}

// QueryEnabledTarget returns true if the first element matching sel is not disabled.
func QueryEnabledTarget(ctx context.Context, debugURL, targetID, sel string) (bool, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) return false;
  return !el.disabled;
})()`, sel)
	var result bool
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return false, fmt.Errorf("query enabled target: %w", err)
	}
	return result, nil
}

// QueryEnabledByBackendNodeID returns true if the element identified by backend node ID is not disabled.
func QueryEnabledByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (bool, error) {
	var result bool
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(`function(){ return !this.disabled; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query enabled by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return false, fmt.Errorf("query enabled target: %w", err)
	}
	return result, nil
}

// QueryCheckedTarget returns true if the first element matching sel is checked.
func QueryCheckedTarget(ctx context.Context, debugURL, targetID, sel string) (bool, error) {
	js := fmt.Sprintf(`(() => {
  const el = document.querySelector(%q);
  if (!el) return false;
  return !!el.checked;
})()`, sel)
	var result bool
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return false, fmt.Errorf("query checked target: %w", err)
	}
	return result, nil
}

// QueryCheckedByBackendNodeID returns true if the element identified by backend node ID is checked.
func QueryCheckedByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) (bool, error) {
	var result bool
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(`function(){ return !!this.checked; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return err
		}
		if ex != nil {
			return fmt.Errorf("query checked by backend node: javascript exception")
		}
		if val != nil && val.Value != nil {
			if err := json.Unmarshal(val.Value, &result); err != nil {
				return fmt.Errorf("unmarshal result: %w", err)
			}
		}
		return nil
	}))
	if err != nil {
		return false, fmt.Errorf("query checked target: %w", err)
	}
	return result, nil
}
