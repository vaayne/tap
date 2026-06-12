package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
)

// HoverElement moves the mouse to an element by CSS selector or snapshot ref (@eN).
func (m *Manager) HoverElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "hover")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := HoverTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("hover: %w", err)
		}
		return nil
	}
	if err := HoverTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("hover: %w", err)
	}
	return nil
}

// DblClickElement double-clicks an element by CSS selector or snapshot ref.
func (m *Manager) DblClickElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "dblclick")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := DblClickTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("dblclick: %w", err)
		}
		return nil
	}
	if err := DblClickTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("dblclick: %w", err)
	}
	return nil
}

// FocusElement focuses an element by CSS selector or snapshot ref.
func (m *Manager) FocusElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "focus")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := FocusTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("focus: %w", err)
		}
		return nil
	}
	if err := FocusTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("focus: %w", err)
	}
	return nil
}

// CheckElement ensures a checkbox is checked. Accepts CSS selector or snapshot ref.
func (m *Manager) CheckElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "check")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := CheckTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("check: %w", err)
		}
		return nil
	}
	if err := CheckTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("check: %w", err)
	}
	return nil
}

// UncheckElement ensures a checkbox is unchecked. Accepts CSS selector or snapshot ref.
func (m *Manager) UncheckElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "uncheck")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := UncheckTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("uncheck: %w", err)
		}
		return nil
	}
	if err := UncheckTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("uncheck: %w", err)
	}
	return nil
}

// ScrollIntoViewElement scrolls an element into view. Accepts CSS selector or snapshot ref.
func (m *Manager) ScrollIntoViewElement(ctx context.Context, sessionName, tabName, arg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "scrollintoview")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := ScrollIntoViewTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID); err != nil {
			return fmt.Errorf("scrollintoview: %w", err)
		}
		return nil
	}
	if err := ScrollIntoViewTarget(ctx, rt.DebugURL, rt.TargetID, selector); err != nil {
		return fmt.Errorf("scrollintoview: %w", err)
	}
	return nil
}

// UploadFiles sets files on a file input element. Accepts CSS selector or snapshot ref.
func (m *Manager) UploadFiles(ctx context.Context, sessionName, tabName, arg string, files []string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "upload")
	if err != nil {
		return err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return err
	}
	if backendNodeID > 0 {
		if err := UploadTargetByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID, files); err != nil {
			return fmt.Errorf("upload: %w", err)
		}
		return nil
	}
	if err := UploadTarget(ctx, rt.DebugURL, rt.TargetID, selector, files); err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	return nil
}

// Drag performs a mouse-based drag from srcArg to dstArg. Both accept CSS
// selectors (snapshot refs are CSS-selector-only for drag since two elements
// are needed).
func (m *Manager) Drag(ctx context.Context, sessionName, tabName, srcArg, dstArg string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "drag")
	if err != nil {
		return err
	}
	// Drag always operates on CSS selectors; resolve @eN refs to selector hints.
	srcSel, srcNodeID, err := m.resolveElementArg(ctx, rt, srcArg)
	if err != nil {
		return err
	}
	dstSel, dstNodeID, err := m.resolveElementArg(ctx, rt, dstArg)
	if err != nil {
		return err
	}
	// For drag with backendNodeIDs we fall back to selector hints since
	// DragTarget uses CSS selector-based box model resolution.
	if srcNodeID > 0 {
		srcSel = selectorFromNodeID(srcSel, srcNodeID)
	}
	if dstNodeID > 0 {
		dstSel = selectorFromNodeID(dstSel, dstNodeID)
	}
	if err := DragTarget(ctx, rt.DebugURL, rt.TargetID, srcSel, dstSel); err != nil {
		return fmt.Errorf("drag: %w", err)
	}
	return nil
}

// selectorFromNodeID returns the selector hint when available, otherwise falls
// back to a data-attribute pseudo-selector for @eN refs. In practice the
// Manager always populates SelectorHint for resolved refs, so the fallback is
// a safety net.
func selectorFromNodeID(selectorHint string, _ cdp.BackendNodeID) string {
	if selectorHint != "" {
		return selectorHint
	}
	return "*"
}

// MouseMove dispatches a mouseMoved event at (x, y).
func (m *Manager) MouseMove(ctx context.Context, sessionName, tabName string, x, y float64) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "mouse move")
	if err != nil {
		return err
	}
	if err := MouseMoveTarget(ctx, rt.DebugURL, rt.TargetID, x, y); err != nil {
		return fmt.Errorf("mouse move: %w", err)
	}
	return nil
}

// MouseDown dispatches a mousePressed event with the given button.
func (m *Manager) MouseDown(ctx context.Context, sessionName, tabName string, button input.MouseButton) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "mouse down")
	if err != nil {
		return err
	}
	if err := MouseDownTarget(ctx, rt.DebugURL, rt.TargetID, button); err != nil {
		return fmt.Errorf("mouse down: %w", err)
	}
	return nil
}

// MouseUp dispatches a mouseReleased event with the given button.
func (m *Manager) MouseUp(ctx context.Context, sessionName, tabName string, button input.MouseButton) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "mouse up")
	if err != nil {
		return err
	}
	if err := MouseUpTarget(ctx, rt.DebugURL, rt.TargetID, button); err != nil {
		return fmt.Errorf("mouse up: %w", err)
	}
	return nil
}

// MouseWheel dispatches a mouseWheel event with the given deltas.
func (m *Manager) MouseWheel(ctx context.Context, sessionName, tabName string, dy, dx float64) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "mouse wheel")
	if err != nil {
		return err
	}
	if err := MouseWheelTarget(ctx, rt.DebugURL, rt.TargetID, dy, dx); err != nil {
		return fmt.Errorf("mouse wheel: %w", err)
	}
	return nil
}

// KeyboardType sends per-character key events for text (like real typing).
func (m *Manager) KeyboardType(ctx context.Context, sessionName, tabName, text string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "keyboard type")
	if err != nil {
		return err
	}
	if err := KeyboardTypeTarget(ctx, rt.DebugURL, rt.TargetID, text); err != nil {
		return fmt.Errorf("keyboard type: %w", err)
	}
	return nil
}

// KeyboardInsert inserts text instantly via Input.insertText (no key events).
func (m *Manager) KeyboardInsert(ctx context.Context, sessionName, tabName, text string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "keyboard insert")
	if err != nil {
		return err
	}
	if err := KeyboardInsertTarget(ctx, rt.DebugURL, rt.TargetID, text); err != nil {
		return fmt.Errorf("keyboard insert: %w", err)
	}
	return nil
}

// Keydown holds a key down.
func (m *Manager) Keydown(ctx context.Context, sessionName, tabName, key string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "keydown")
	if err != nil {
		return err
	}
	if err := KeydownTarget(ctx, rt.DebugURL, rt.TargetID, key); err != nil {
		return fmt.Errorf("keydown: %w", err)
	}
	return nil
}

// Keyup releases a held key.
func (m *Manager) Keyup(ctx context.Context, sessionName, tabName, key string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "keyup")
	if err != nil {
		return err
	}
	if err := KeyupTarget(ctx, rt.DebugURL, rt.TargetID, key); err != nil {
		return fmt.Errorf("keyup: %w", err)
	}
	return nil
}
