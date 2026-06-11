package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/accessibility"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// LocatorKind identifies which semantic locator strategy to use.
type LocatorKind string

const (
	LocatorRole        LocatorKind = "role"
	LocatorText        LocatorKind = "text"
	LocatorLabel       LocatorKind = "label"
	LocatorPlaceholder LocatorKind = "placeholder"
	LocatorAlt         LocatorKind = "alt"
	LocatorTitle       LocatorKind = "title"
	LocatorTestID      LocatorKind = "testid"
	LocatorFirst       LocatorKind = "first"
	LocatorLast        LocatorKind = "last"
	LocatorNth         LocatorKind = "nth"
)

// FindAction is the operation to perform once an element is located.
type FindAction string

const (
	FindActionClick   FindAction = "click"
	FindActionFill    FindAction = "fill"
	FindActionType    FindAction = "type"
	FindActionHover   FindAction = "hover"
	FindActionFocus   FindAction = "focus"
	FindActionCheck   FindAction = "check"
	FindActionUncheck FindAction = "uncheck"
	FindActionText    FindAction = "text"
)

// FindLocator describes a semantic element locator.
type FindLocator struct {
	Kind LocatorKind

	// Role locator fields.
	Role string
	Name string // accessible name filter (--name flag)

	// Text locator fields.
	Text  string
	Exact bool // --exact flag for text matching

	// Label, placeholder, alt, title, testid: all use Query.
	Query string

	// Nth/First/Last: Index and CSSSelector.
	Index       int    // 0-based index; -1 means last
	CSSSelector string // the CSS selector for nth/first/last
}

// FindTarget locates an element using the given locator, then performs action
// with optional value. Returns textContent when action is FindActionText.
func FindTarget(ctx context.Context, debugURL, targetID string, loc FindLocator, action FindAction, value string) (string, error) {
	var result string
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		backendNodeID, err := resolveLocator(ctx, loc)
		if err != nil {
			return err
		}
		out, err := dispatchFindAction(ctx, backendNodeID, action, value)
		if err != nil {
			return err
		}
		result = out
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("find %s: %w", loc.Kind, err)
	}
	return result, nil
}

// resolveLocator maps a FindLocator to a BackendNodeID using the appropriate
// strategy for each kind.
func resolveLocator(ctx context.Context, loc FindLocator) (cdp.BackendNodeID, error) {
	switch loc.Kind {
	case LocatorRole:
		return resolveByRole(ctx, loc.Role, loc.Name)
	case LocatorText:
		return resolveByJS(ctx, jsLocatorText(loc.Text, loc.Exact), fmt.Sprintf("text=%q", loc.Text))
	case LocatorLabel:
		return resolveByJS(ctx, jsLocatorLabel(loc.Query), fmt.Sprintf("label=%q", loc.Query))
	case LocatorPlaceholder:
		return resolveByJS(ctx, jsLocatorAttr("placeholder", loc.Query), fmt.Sprintf("placeholder=%q", loc.Query))
	case LocatorAlt:
		return resolveByJS(ctx, jsLocatorAttr("alt", loc.Query), fmt.Sprintf("alt=%q", loc.Query))
	case LocatorTitle:
		return resolveByJS(ctx, jsLocatorAttr("title", loc.Query), fmt.Sprintf("title=%q", loc.Query))
	case LocatorTestID:
		return resolveByJS(ctx, jsLocatorAttr("data-testid", loc.Query), fmt.Sprintf("testid=%q", loc.Query))
	case LocatorFirst:
		return resolveByJS(ctx, jsLocatorNth(loc.CSSSelector, 0), fmt.Sprintf("first %q", loc.CSSSelector))
	case LocatorLast:
		return resolveByJS(ctx, jsLocatorLast(loc.CSSSelector), fmt.Sprintf("last %q", loc.CSSSelector))
	case LocatorNth:
		return resolveByJS(ctx, jsLocatorNth(loc.CSSSelector, loc.Index), fmt.Sprintf("nth(%d) %q", loc.Index, loc.CSSSelector))
	default:
		return 0, fmt.Errorf("unknown locator kind %q", loc.Kind)
	}
}

// resolveByRole uses the AX tree to find the first node matching role (and
// optionally accessible name), then returns its BackendNodeID.
//
// Role matching is case-insensitive. Name matching uses substring containment
// (also case-insensitive). If multiple nodes match, the first is returned.
func resolveByRole(ctx context.Context, role, name string) (cdp.BackendNodeID, error) {
	if err := accessibility.Enable().Do(ctx); err != nil {
		return 0, fmt.Errorf("enable accessibility: %w", err)
	}
	defer func() { _ = accessibility.Disable().Do(ctx) }()

	nodes, err := accessibility.GetFullAXTree().Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("get ax tree: %w", err)
	}

	wantRole := strings.ToLower(strings.TrimSpace(role))
	wantName := strings.ToLower(strings.TrimSpace(name))

	for _, n := range nodes {
		nodeRole := strings.ToLower(normalizeAXValue(n.Role))
		if nodeRole != wantRole {
			continue
		}
		if wantName != "" {
			nodeName := strings.ToLower(normalizeAXValue(n.Name))
			if !strings.Contains(nodeName, wantName) {
				continue
			}
		}
		if n.BackendDOMNodeID == 0 {
			continue
		}
		return cdp.BackendNodeID(n.BackendDOMNodeID), nil
	}

	desc := fmt.Sprintf("role=%q", role)
	if name != "" {
		desc += fmt.Sprintf(" name=%q", name)
	}
	return 0, fmt.Errorf("no element found: %s", desc)
}

// resolveByJS runs a JS locator expression inside the page that returns a DOM
// node reference, then uses dom.DescribeNode to obtain the BackendNodeID.
// The expression must return the element or null/undefined.
func resolveByJS(ctx context.Context, jsExpr string, desc string) (cdp.BackendNodeID, error) {
	obj, ex, err := runtime.Evaluate(jsExpr).
		WithReturnByValue(false).
		WithAwaitPromise(false).
		Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("locator js: %w", err)
	}
	if ex != nil {
		return 0, fmt.Errorf("locator js exception: %s", ex.Text)
	}
	if obj == nil || obj.ObjectID == "" || obj.Type == "undefined" ||
		(obj.Type == "object" && obj.Subtype == "null") {
		return 0, fmt.Errorf("no element found: %s", desc)
	}

	nodeInfo, err := dom.DescribeNode().WithObjectID(obj.ObjectID).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("describe node: %w", err)
	}
	if nodeInfo.BackendNodeID == 0 {
		return 0, fmt.Errorf("element has no backend node ID: %s", desc)
	}
	return nodeInfo.BackendNodeID, nil
}

// ---------------------------------------------------------------------------
// JS locator expressions — each returns the first matching DOM element or null.
// ---------------------------------------------------------------------------

// jsLocatorText returns a JS expression that finds the first element whose
// visible text contains (or exactly matches) the given text.
func jsLocatorText(text string, exact bool) string {
	if exact {
		return fmt.Sprintf(`(function(){
  var text = %q;
  var all = document.querySelectorAll("*");
  for (var i = 0; i < all.length; i++) {
    var el = all[i];
    if (el.children.length === 0 && el.textContent.trim() === text) return el;
    if ((el.tagName === "INPUT" || el.tagName === "TEXTAREA") && (el.value||"").trim() === text) return el;
  }
  return null;
})()`, text)
	}
	lowerText := strings.ToLower(text)
	return fmt.Sprintf(`(function(){
  var text = %q;
  var all = document.querySelectorAll("*");
  for (var i = 0; i < all.length; i++) {
    var el = all[i];
    if (el.children.length === 0 && el.textContent.toLowerCase().includes(text)) return el;
    if ((el.tagName === "INPUT" || el.tagName === "TEXTAREA") && (el.value||"").toLowerCase().includes(text)) return el;
  }
  return null;
})()`, lowerText)
}

// jsLocatorLabel returns a JS expression that finds the first form element
// associated with a label whose text contains the given string. It checks
// label[for=id], wrapping label elements, and aria-label attributes.
func jsLocatorLabel(labelText string) string {
	lower := strings.ToLower(labelText)
	return fmt.Sprintf(`(function(){
  var want = %q;
  var labels = document.querySelectorAll("label");
  for (var i = 0; i < labels.length; i++) {
    var lbl = labels[i];
    if (!lbl.textContent.toLowerCase().includes(want)) continue;
    if (lbl.htmlFor) {
      var el = document.getElementById(lbl.htmlFor);
      if (el) return el;
    }
    var inp = lbl.querySelector("input,textarea,select,button");
    if (inp) return inp;
  }
  var all = document.querySelectorAll("[aria-label]");
  for (var j = 0; j < all.length; j++) {
    if ((all[j].getAttribute("aria-label")||"").toLowerCase().includes(want)) return all[j];
  }
  return null;
})()`, lower)
}

// jsLocatorAttr returns a JS expression that finds the first element with
// the given attribute value containing the query string (case-insensitive).
func jsLocatorAttr(attr, query string) string {
	lower := strings.ToLower(query)
	return fmt.Sprintf(`(function(){
  var want = %q;
  var all = document.querySelectorAll("[%s]");
  for (var i = 0; i < all.length; i++) {
    if ((all[i].getAttribute(%q)||"").toLowerCase().includes(want)) return all[i];
  }
  return null;
})()`, lower, attr, attr)
}

// jsLocatorNth returns a JS expression that returns the n-th (0-based) element
// matching the CSS selector.
func jsLocatorNth(css string, n int) string {
	return fmt.Sprintf(`(function(){
  var els = document.querySelectorAll(%q);
  return els[%d] || null;
})()`, css, n)
}

// jsLocatorLast returns a JS expression that returns the last element matching
// the CSS selector.
func jsLocatorLast(css string) string {
	return fmt.Sprintf(`(function(){
  var els = document.querySelectorAll(%q);
  return els.length > 0 ? els[els.length - 1] : null;
})()`, css)
}

// ---------------------------------------------------------------------------
// Action dispatch — all functions operate on an already-attached CDP context.
// ---------------------------------------------------------------------------

// dispatchFindAction performs action on the element identified by backendNodeID.
// Returns non-empty string only when action == FindActionText.
func dispatchFindAction(ctx context.Context, backendNodeID cdp.BackendNodeID, action FindAction, value string) (string, error) {
	switch action {
	case FindActionClick:
		return "", findClick(ctx, backendNodeID)
	case FindActionFill:
		return "", findFill(ctx, backendNodeID, value)
	case FindActionType:
		return "", findType(ctx, backendNodeID, value)
	case FindActionHover:
		return "", findHover(ctx, backendNodeID)
	case FindActionFocus:
		return "", findFocusJS(ctx, backendNodeID)
	case FindActionCheck:
		return "", findSetChecked(ctx, backendNodeID, true)
	case FindActionUncheck:
		return "", findSetChecked(ctx, backendNodeID, false)
	case FindActionText:
		return findGetText(ctx, backendNodeID)
	default:
		return "", fmt.Errorf("unknown action %q", action)
	}
}

// findClick scrolls the element into view and dispatches a real mouse click.
func findClick(ctx context.Context, backendNodeID cdp.BackendNodeID) error {
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
	if err := input.DispatchMouseEvent(input.MousePressed, x, y).WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
		return err
	}
	return input.DispatchMouseEvent(input.MouseReleased, x, y).WithButton(input.Left).WithClickCount(1).Do(ctx)
}

// findFill sets the value of a form field via React-compatible native setter.
// This is the inline equivalent of FillTargetByBackendNodeID but runs inside an
// already-attached CDP context (no withTarget wrapper needed).
func findFill(ctx context.Context, backendNodeID cdp.BackendNodeID, value string) error {
	obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return fmt.Errorf("resolve node: %w", err)
	}
	fn := `function(v){
  this.focus();
  var tag = this.tagName ? this.tagName.toLowerCase() : "";
  if (tag === "select") {
    this.value = v;
    this.dispatchEvent(new Event("change", { bubbles: true }));
    return true;
  }
  var setter = (tag === "textarea")
    ? Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value") && Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value").set
    : Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value") && Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value").set;
  if (setter) setter.call(this, v); else this.value = v;
  this.dispatchEvent(new Event("input", { bubbles: true }));
  this.dispatchEvent(new Event("change", { bubbles: true }));
  return true;
}`
	_, ex, err := runtime.CallFunctionOn(fmt.Sprintf(`function(){ return (%s).call(this, %q) }`, fn, value)).
		WithObjectID(obj.ObjectID).
		WithAwaitPromise(true).
		WithReturnByValue(true).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("fill: %w", err)
	}
	if ex != nil {
		return fmt.Errorf("fill: javascript exception")
	}
	return nil
}

// findType focuses the element via JS then inserts text using input.InsertText.
func findType(ctx context.Context, backendNodeID cdp.BackendNodeID, text string) error {
	obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return fmt.Errorf("resolve node: %w", err)
	}
	_, ex, err := runtime.CallFunctionOn(`function(){ this.focus(); return true; }`).
		WithObjectID(obj.ObjectID).
		WithReturnByValue(true).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("focus for type: %w", err)
	}
	if ex != nil {
		return fmt.Errorf("focus for type: javascript exception")
	}
	return input.InsertText(text).Do(ctx)
}

// findHover moves the mouse to the element centre using box model coordinates.
func findHover(ctx context.Context, backendNodeID cdp.BackendNodeID) error {
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
}

// findFocusJS focuses the element via JS CallFunctionOn.
func findFocusJS(ctx context.Context, backendNodeID cdp.BackendNodeID) error {
	obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return fmt.Errorf("resolve node: %w", err)
	}
	_, ex, err := runtime.CallFunctionOn(`function(){ this.focus(); return true; }`).
		WithObjectID(obj.ObjectID).
		WithReturnByValue(true).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("focus: %w", err)
	}
	if ex != nil {
		return fmt.Errorf("focus: javascript exception")
	}
	return nil
}

// findSetChecked checks or unchecks a checkbox/radio input via JS.
func findSetChecked(ctx context.Context, backendNodeID cdp.BackendNodeID, checked bool) error {
	obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return fmt.Errorf("resolve node: %w", err)
	}
	wantStr := "true"
	if !checked {
		wantStr = "false"
	}
	fn := fmt.Sprintf(`function(){
  var want = %s;
  if (this.checked !== want) {
    this.click();
    this.dispatchEvent(new Event("change", { bubbles: true }));
  }
  return true;
}`, wantStr)
	_, ex, err := runtime.CallFunctionOn(fn).
		WithObjectID(obj.ObjectID).
		WithReturnByValue(true).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("set checked: %w", err)
	}
	if ex != nil {
		return fmt.Errorf("set checked: javascript exception")
	}
	return nil
}

// findGetText returns the trimmed textContent of the element.
func findGetText(ctx context.Context, backendNodeID cdp.BackendNodeID) (string, error) {
	obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("resolve node: %w", err)
	}
	res, ex, err := runtime.CallFunctionOn(`function(){ return (this.textContent||"").trim(); }`).
		WithObjectID(obj.ObjectID).
		WithReturnByValue(true).
		Do(ctx)
	if err != nil {
		return "", fmt.Errorf("get text: %w", err)
	}
	if ex != nil {
		return "", fmt.Errorf("get text: javascript exception")
	}
	if res == nil || len(res.Value) == 0 {
		return "", nil
	}
	var text string
	if err := json.Unmarshal(res.Value, &text); err != nil {
		return strings.Trim(string(res.Value), `"`), nil
	}
	return text, nil
}
