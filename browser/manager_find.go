package browser

import (
	"context"
	"fmt"
)

// Find locates an element using the given semantic locator and dispatches action.
// When action is FindActionText the trimmed textContent is returned; otherwise
// the return string is always empty.
//
// Multiple matches: the first matching element is used.
// No match: returns an error describing the locator that was searched.
func (m *Manager) Find(ctx context.Context, sessionName, tabName string, loc FindLocator, action FindAction, value string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "find")
	if err != nil {
		return "", err
	}
	result, err := FindTarget(ctx, rt.DebugURL, rt.TargetID, loc, action, value)
	if err != nil {
		return "", fmt.Errorf("find: %w", err)
	}
	return result, nil
}
