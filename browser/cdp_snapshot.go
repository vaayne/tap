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

// SnapshotTarget builds an AX-powered semantic snapshot for a target tab.
func SnapshotTarget(ctx context.Context, debugURL, targetID string, opts SnapshotOptions) (*SnapshotResult, error) {
	if opts.Mode == "" {
		opts.Mode = "auto"
	}

	out := &SnapshotResult{Mode: "ax"}
	err := withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		if err := accessibility.Enable().Do(ctx); err != nil {
			return fmt.Errorf("enable accessibility: %w", err)
		}
		defer func() { _ = accessibility.Disable().Do(ctx) }()

		var nodes []*accessibility.Node
		var err error
		if opts.Depth > 0 {
			nodes, err = accessibility.GetFullAXTree().WithDepth(opts.Depth).Do(ctx)
		} else {
			nodes, err = accessibility.GetFullAXTree().Do(ctx)
		}
		if err != nil {
			return fmt.Errorf("get full ax tree: %w", err)
		}
		meta, err := EvalTarget(ctx, debugURL, targetID, `({
  url: location.href,
  title: document.title || "",
  key: [location.href, performance.timeOrigin || 0, document.body ? document.body.childElementCount : 0].join("|")
})`)
		if err != nil {
			return fmt.Errorf("snapshot metadata: %w", err)
		}
		if m, ok := meta.(map[string]any); ok {
			if v, ok := m["url"].(string); ok {
				out.URL = v
			}
			if v, ok := m["title"].(string); ok {
				out.Title = v
			}
			if v, ok := m["key"].(string); ok {
				out.DocumentKey = v
			}
		}
		buildSnapshotTree(out, nodes, opts)
		return nil
	}))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func buildSnapshotTree(out *SnapshotResult, axNodes []*accessibility.Node, opts SnapshotOptions) {
	idToNode := make(map[accessibility.NodeID]*accessibility.Node, len(axNodes))
	for _, n := range axNodes {
		idToNode[n.NodeID] = n
	}

	interactiveRoles := map[string]bool{
		"button": true, "link": true, "textbox": true, "searchbox": true, "checkbox": true,
		"radio": true, "combobox": true, "option": true, "menuitem": true, "tab": true,
	}
	structuralRoles := map[string]bool{
		"rootwebarea": true, "document": true, "main": true, "article": true,
		"heading": true, "list": true, "listitem": true, "dialog": true,
	}

	axToIndex := make(map[accessibility.NodeID]int, len(axNodes))
	nextRef := 1
	for _, n := range axNodes {
		role := normalizeAXValue(n.Role)
		name := normalizeAXValue(n.Name)
		if n.Ignored && !interactiveRoles[role] {
			continue
		}
		keep := structuralRoles[role] || interactiveRoles[role] || name != ""
		if opts.InteractiveOnly && !interactiveRoles[role] {
			keep = false
		}
		if !keep {
			continue
		}

		node := SnapshotNode{
			Role:        role,
			Name:        name,
			Value:       normalizeAXValue(n.Value),
			Description: normalizeAXValue(n.Description),
			States:      axStates(n.Properties),
		}
		if interactiveRoles[role] {
			node.Ref = fmt.Sprintf("@e%d", nextRef)
			nextRef++
			out.Refs = append(out.Refs, SnapshotRef{
				Ref:              node.Ref,
				BackendDOMNodeID: int64(n.BackendDOMNodeID),
				AXNodeID:         string(n.NodeID),
				FrameID:          string(n.FrameID),
				Role:             role,
				Name:             name,
			})
		}
		axToIndex[n.NodeID] = len(out.Nodes)
		out.Nodes = append(out.Nodes, node)
	}

	for _, n := range axNodes {
		parentIdx, ok := axToIndex[n.NodeID]
		if !ok {
			continue
		}
		for _, childID := range n.ChildIDs {
			if childIdx, ok := axToIndex[childID]; ok {
				out.Nodes[parentIdx].Children = append(out.Nodes[parentIdx].Children, childIdx)
			}
		}
	}
}

func normalizeAXValue(v *accessibility.Value) string {
	if v == nil {
		return ""
	}
	var anyv any
	if err := json.Unmarshal(v.Value, &anyv); err != nil {
		return ""
	}
	switch tv := anyv.(type) {
	case string:
		return strings.TrimSpace(tv)
	case bool, float64:
		return fmt.Sprintf("%v", tv)
	default:
		return ""
	}
}

func axStates(props []*accessibility.Property) []string {
	out := make([]string, 0, len(props))
	for _, p := range props {
		name := string(p.Name)
		if name == "disabled" || name == "checked" || name == "expanded" || name == "selected" {
			val := normalizeAXValue(p.Value)
			if val == "" {
				val = "true"
			}
			out = append(out, fmt.Sprintf("%s=%s", name, val))
		}
	}
	return out
}

func ClickTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID) error {
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
		if err := input.DispatchMouseEvent(input.MousePressed, x, y).WithButton(input.Left).WithClickCount(1).Do(ctx); err != nil {
			return err
		}
		return input.DispatchMouseEvent(input.MouseReleased, x, y).WithButton(input.Left).WithClickCount(1).Do(ctx)
	}))
}

func FillTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID, value string) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		fn := `function(v){
  this.focus();
  const tag = this.tagName ? this.tagName.toLowerCase() : "";
  if (tag === "select") {
    this.value = v;
    this.dispatchEvent(new Event("change", { bubbles: true }));
    return true;
  }
  const setter = (tag === "textarea")
    ? Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value")?.set
    : Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
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
			return fmt.Errorf("fill by backend node: %w", err)
		}
		if ex != nil {
			return fmt.Errorf("fill by backend node: javascript exception")
		}
		return nil
	}))
}

func TypeTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID, text string) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		if _, ex, err := runtime.CallFunctionOn(`function(){ this.focus(); return true; }`).
			WithObjectID(obj.ObjectID).
			WithReturnByValue(true).
			Do(ctx); err != nil {
			return err
		} else if ex != nil {
			return fmt.Errorf("focus by backend node: javascript exception")
		}
		return input.InsertText(text).Do(ctx)
	}))
}

func SelectTargetByBackendNodeID(ctx context.Context, debugURL, targetID string, backendNodeID cdp.BackendNodeID, value string) error {
	return withTarget(ctx, debugURL, targetID, chromedp.ActionFunc(func(ctx context.Context) error {
		obj, err := dom.ResolveNode().WithBackendNodeID(backendNodeID).Do(ctx)
		if err != nil {
			return fmt.Errorf("resolve node: %w", err)
		}
		fn := `function(v){
  if (!this.tagName || this.tagName.toLowerCase() !== "select") throw new Error("element is not select");
  this.focus();
  this.value = v;
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
			return err
		}
		if ex != nil {
			return fmt.Errorf("select by backend node: javascript exception")
		}
		return nil
	}))
}
