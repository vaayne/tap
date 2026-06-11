package browser

import (
	"context"
	"fmt"
	"time"
)

// WaitForElement waits until the element matching sel reaches the given DOM
// state in a tracked tab.
func (m *Manager) WaitForElement(ctx context.Context, sessionName, tabName, sel string, state ElementState, timeout time.Duration) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "wait element")
	if err != nil {
		return err
	}
	if err := WaitForElementTarget(ctx, rt.DebugURL, rt.TargetID, sel, state, timeout); err != nil {
		return fmt.Errorf("wait element: %w", err)
	}
	return nil
}

// WaitForText waits until document.body.innerText contains the given substring
// in a tracked tab.
func (m *Manager) WaitForText(ctx context.Context, sessionName, tabName, text string, timeout time.Duration) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "wait text")
	if err != nil {
		return err
	}
	if err := WaitForTextTarget(ctx, rt.DebugURL, rt.TargetID, text, timeout); err != nil {
		return fmt.Errorf("wait text: %w", err)
	}
	return nil
}

// WaitForURL waits until the page's location.href matches the glob pattern in a
// tracked tab.
func (m *Manager) WaitForURL(ctx context.Context, sessionName, tabName, glob string, timeout time.Duration) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "wait url")
	if err != nil {
		return err
	}
	if err := WaitForURLTarget(ctx, rt.DebugURL, rt.TargetID, glob, timeout); err != nil {
		return fmt.Errorf("wait url: %w", err)
	}
	return nil
}

// WaitForLoad waits for a named page-load event in a tracked tab.
func (m *Manager) WaitForLoad(ctx context.Context, sessionName, tabName string, state LoadState, timeout time.Duration) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "wait load")
	if err != nil {
		return err
	}
	if err := WaitForLoadTarget(ctx, rt.DebugURL, rt.TargetID, state, timeout); err != nil {
		return fmt.Errorf("wait load: %w", err)
	}
	return nil
}

// WaitForFn waits until the given JS expression evaluates to a truthy value in
// a tracked tab.
func (m *Manager) WaitForFn(ctx context.Context, sessionName, tabName, jsExpr string, timeout time.Duration) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "wait fn")
	if err != nil {
		return err
	}
	if err := WaitForFnTarget(ctx, rt.DebugURL, rt.TargetID, jsExpr, timeout); err != nil {
		return fmt.Errorf("wait fn: %w", err)
	}
	return nil
}
