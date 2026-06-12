package browser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// HoverTargetByBackendNodeID moves the mouse to the center of the node
// identified by backendNodeID, dispatching a mouseMoved event.
func HoverTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		if err := dom.ScrollIntoViewIfNeeded().WithBackendNodeID(backendNodeID).Do(ctx); err != nil {
			return fmt.Errorf("scroll into view: %w", err)
		}
		box, err := dom.GetBoxModel().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("get box model: %w", err)
		}
		q := box.Content
		x := (q[0] + q[2] + q[4] + q[6]) / 4
		y := (q[1] + q[3] + q[5] + q[7]) / 4
		return input.DispatchMouseEvent(input.MouseMoved, x, y).Do(ctx)
	}))
}

// DblClickTarget dispatches a double-click (clickCount=2) on the first
// visible element matching sel.
func DblClickTarget(ctx context.Context, debugURL, targetID, sel string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.DoubleClick(sel, chromedp.ByQuery, chromedp.NodeVisible),
	); err != nil {
		return fmt.Errorf("dblclick target: %w", err)
	}
	return nil
}

// DblClickTargetByBackendNodeID dispatches a double-click via box model on a
// node resolved by its BackendNodeID.
func DblClickTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		if err := dom.ScrollIntoViewIfNeeded().WithBackendNodeID(backendNodeID).Do(ctx); err != nil {
			return fmt.Errorf("scroll into view: %w", err)
		}
		box, err := dom.GetBoxModel().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("get box model: %w", err)
		}
		q := box.Content
		x := (q[0] + q[2] + q[4] + q[6]) / 4
		y := (q[1] + q[3] + q[5] + q[7]) / 4
		if err := input.DispatchMouseEvent(input.MouseMoved, x, y).Do(ctx); err != nil {
			return err
		}
		if err := input.DispatchMouseEvent(input.MousePressed, x, y).
			WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
			return err
		}
		if err := input.DispatchMouseEvent(input.MouseReleased, x, y).
			WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
			return err
		}
		if err := input.DispatchMouseEvent(input.MousePressed, x, y).
			WithButton(input.Left).WithClickCount(2).Do(ctx); err != nil {
			return err
		}
		return input.DispatchMouseEvent(input.MouseReleased, x, y).
			WithButton(input.Left).WithClickCount(2).Do(ctx)
	}))
}

// FocusTarget focuses the first visible element matching sel.
func FocusTarget(ctx context.Context, debugURL, targetID, sel string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Focus(sel, chromedp.ByQuery, chromedp.NodeVisible),
	); err != nil {
		return fmt.Errorf("focus target: %w", err)
	}
	return nil
}

// FocusTargetByBackendNodeID focuses the node identified by backendNodeID.
func FocusTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		nodeID, err := resolveNodeID(ctx, backendNodeID)
		if err != nil {
			return err
		}
		if err := dom.Focus().WithNodeID(nodeID).Do(ctx); err != nil {
			return fmt.Errorf("focus by backend node: %w", err)
		}
		return nil
	}))
}

// checkboxJS is a React-compatible snippet that reads and conditionally
// toggles a checkbox to match the desired state, then dispatches change/input.
const checkboxJS = `
(function(sel, want) {
  var el = document.querySelector(sel);
  if (!el) throw new Error("element not found: " + sel);
  if (el.type !== "checkbox") throw new Error("element is not a checkbox: " + sel);
  if (el.checked !== want) {
    el.click();
  }
  return el.checked;
})(%q, %v)
`

// CheckTarget ensures the checkbox matching sel is checked.
func CheckTarget(ctx context.Context, debugURL, targetID, sel string) error {
	js := fmt.Sprintf(checkboxJS, sel, true)
	var checked bool
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &checked, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return fmt.Errorf("check target: %w", err)
	}
	if !checked {
		return fmt.Errorf("check target: element is still unchecked after click")
	}
	return nil
}

// UncheckTarget ensures the checkbox matching sel is unchecked.
func UncheckTarget(ctx context.Context, debugURL, targetID, sel string) error {
	js := fmt.Sprintf(checkboxJS, sel, false)
	var checked bool
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(js, &checked, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithReturnByValue(true).WithAwaitPromise(true)
		}),
	); err != nil {
		return fmt.Errorf("uncheck target: %w", err)
	}
	if checked {
		return fmt.Errorf("uncheck target: element is still checked after click")
	}
	return nil
}

// checkboxByNodeJS is a React-compatible snippet invoked via CallFunctionOn.
const checkboxByNodeJS = `function(want){
  if (this.type !== "checkbox") throw new Error("element is not a checkbox");
  if (this.checked !== want) {
    this.click();
  }
  return this.checked;
}`

// CheckTargetByBackendNodeID ensures a checkbox is checked.
func CheckTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(fmt.Sprintf(`function(){ return (%s).call(this, true) }`, checkboxByNodeJS)).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("check by backend node: %w", err)
		}
		if ex != nil {
			return fmt.Errorf("check by backend node: javascript exception")
		}
		var checked bool
		if val != nil {
			if jsonErr := json.Unmarshal(val.Value, &checked); jsonErr == nil && !checked {
				return fmt.Errorf("check by backend node: element is still unchecked after click")
			}
		}
		return nil
	}))
}

// UncheckTargetByBackendNodeID ensures a checkbox is unchecked.
func UncheckTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		val, ex, err := runtime.CallFunctionOn(fmt.Sprintf(`function(){ return (%s).call(this, false) }`, checkboxByNodeJS)).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("uncheck by backend node: %w", err)
		}
		if ex != nil {
			return fmt.Errorf("uncheck by backend node: javascript exception")
		}
		var checked bool
		if val != nil {
			if jsonErr := json.Unmarshal(val.Value, &checked); jsonErr == nil && checked {
				return fmt.Errorf("uncheck by backend node: element is still checked after click")
			}
		}
		return nil
	}))
}

// ScrollIntoViewTarget scrolls the first element matching sel into view.
func ScrollIntoViewTarget(ctx context.Context, debugURL, targetID, sel string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ScrollIntoView(sel, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("scrollintoview target: %w", err)
	}
	return nil
}

// ScrollIntoViewTargetByBackendNodeID scrolls the node identified by
// backendNodeID into view.
func ScrollIntoViewTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		if err := dom.ScrollIntoViewIfNeeded().WithBackendNodeID(backendNodeID).Do(ctx); err != nil {
			return fmt.Errorf("scrollintoview by backend node: %w", err)
		}
		return nil
	}))
}

// UploadTarget sets files on the <input type=file> matching sel via
// DOM.setFileInputFiles.
func UploadTarget(ctx context.Context, debugURL, targetID, sel string, files []string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.SetUploadFiles(sel, files, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("upload target: %w", err)
	}
	return nil
}

// UploadTargetByBackendNodeID sets files on a file input node via
// DOM.setFileInputFiles using BackendNodeID.
func UploadTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID, files []string) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		if err := dom.SetFileInputFiles(files).WithBackendNodeID(backendNodeID).Do(ctx); err != nil {
			return fmt.Errorf("upload by backend node: %w", err)
		}
		return nil
	}))
}

// DragTarget performs a mouse-based drag from the center of srcSel to the
// center of dstSel using a sequence of mouseMoved→mousePressed→moves→mouseReleased.
func DragTarget(ctx context.Context, debugURL, targetID, srcSel, dstSel string) error {
	var srcNodes, dstNodes []*cdp.Node
	return withTarget(ctx, debugURL, targetID,
		chromedp.Nodes(srcSel, &srcNodes, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.Nodes(dstSel, &dstNodes, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(srcNodes) == 0 {
				return fmt.Errorf("no element matching src selector %q", srcSel)
			}
			if len(dstNodes) == 0 {
				return fmt.Errorf("no element matching dst selector %q", dstSel)
			}
			srcBox, err := dom.GetBoxModel().WithNodeID(srcNodes[0].NodeID).Do(ctx)
			if err != nil {
				return fmt.Errorf("get src box model: %w", err)
			}
			dstBox, err := dom.GetBoxModel().WithNodeID(dstNodes[0].NodeID).Do(ctx)
			if err != nil {
				return fmt.Errorf("get dst box model: %w", err)
			}
			sc := srcBox.Content
			sx := (sc[0] + sc[2] + sc[4] + sc[6]) / 4
			sy := (sc[1] + sc[3] + sc[5] + sc[7]) / 4
			dc := dstBox.Content
			dx := (dc[0] + dc[2] + dc[4] + dc[6]) / 4
			dy := (dc[1] + dc[3] + dc[5] + dc[7]) / 4

			// Move to source, press, interpolate to destination, release.
			if err := input.DispatchMouseEvent(input.MouseMoved, sx, sy).Do(ctx); err != nil {
				return err
			}
			if err := input.DispatchMouseEvent(input.MousePressed, sx, sy).
				WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
				return err
			}
			const steps = 10
			for i := 1; i <= steps; i++ {
				t := float64(i) / float64(steps)
				mx := sx + (dx-sx)*t
				my := sy + (dy-sy)*t
				if err := input.DispatchMouseEvent(input.MouseMoved, mx, my).
					WithButton(input.Left).WithButtons(1).Do(ctx); err != nil {
					return err
				}
			}
			return input.DispatchMouseEvent(input.MouseReleased, dx, dy).
				WithButton(input.Left).WithClickCount(1).Do(ctx)
		}),
	)
}

// MouseMoveTarget dispatches a mouseMoved event at (x, y).
func MouseMoveTarget(ctx context.Context, debugURL, targetID string, x, y float64) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MouseMoved, x, y).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("mouse move target: %w", err)
	}
	return nil
}

// MouseDownTarget dispatches a mousePressed event at the current cursor
// position with the given button (left|right|middle).
func MouseDownTarget(ctx context.Context, debugURL, targetID string, button input.MouseButton) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MousePressed, 0, 0).
				WithButton(button).WithClickCount(1).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("mouse down target: %w", err)
	}
	return nil
}

// MouseUpTarget dispatches a mouseReleased event with the given button.
func MouseUpTarget(ctx context.Context, debugURL, targetID string, button input.MouseButton) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MouseReleased, 0, 0).
				WithButton(button).WithClickCount(1).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("mouse up target: %w", err)
	}
	return nil
}

// MouseWheelTarget dispatches a mouseWheel event with the given deltas (dy
// scrolls vertically, dx scrolls horizontally).
func MouseWheelTarget(ctx context.Context, debugURL, targetID string, dy, dx float64) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MouseWheel, 0, 0).
				WithDeltaY(dy).WithDeltaX(dx).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("mouse wheel target: %w", err)
	}
	return nil
}

// KeyboardTypeTarget sends per-character keyDown/char/keyUp events for text
// using chromedp.KeyEvent — same as real typing.
func KeyboardTypeTarget(ctx context.Context, debugURL, targetID, text string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.KeyEvent(text),
	); err != nil {
		return fmt.Errorf("keyboard type target: %w", err)
	}
	return nil
}

// KeyboardInsertTarget calls Input.insertText — no key events, instant paste.
func KeyboardInsertTarget(ctx context.Context, debugURL, targetID, text string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.InsertText(text).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("keyboard insert target: %w", err)
	}
	return nil
}

// KeydownTarget holds a key down (sends a rawKeyDown event).
func KeydownTarget(ctx context.Context, debugURL, targetID, key string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchKeyEvent(input.KeyRawDown).WithKey(key).WithCode(key).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("keydown target: %w", err)
	}
	return nil
}

// KeyupTarget releases a held key (sends a keyUp event).
func KeyupTarget(ctx context.Context, debugURL, targetID, key string) error {
	if err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchKeyEvent(input.KeyUp).WithKey(key).WithCode(key).Do(ctx)
		}),
	); err != nil {
		return fmt.Errorf("keyup target: %w", err)
	}
	return nil
}

// resolveNodeID is a small helper that maps a BackendNodeID to a NodeID via
// DOM.describeNode, which is cheaper than a full DOM.resolveNode.
func resolveNodeID(ctx context.Context, backendNodeID cdp.BackendNodeID) (cdp.NodeID, error) {
	node, err := dom.DescribeNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("describe node: %w", err)
	}
	return node.NodeID, nil
}
