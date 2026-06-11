package browser

import (
	"context"
	"fmt"
)

// QueryText returns the textContent of the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryText(ctx context.Context, sessionName, tabName, arg string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get text")
	if err != nil {
		return "", err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return "", err
	}
	if backendNodeID > 0 {
		result, err := QueryTextByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return "", fmt.Errorf("get text: %w", err)
		}
		return result, nil
	}
	result, err := QueryTextTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return "", fmt.Errorf("get text: %w", err)
	}
	return result, nil
}

// QueryHTML returns the innerHTML of the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryHTML(ctx context.Context, sessionName, tabName, arg string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get html")
	if err != nil {
		return "", err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return "", err
	}
	if backendNodeID > 0 {
		result, err := QueryHTMLByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return "", fmt.Errorf("get html: %w", err)
		}
		return result, nil
	}
	result, err := QueryHTMLTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return "", fmt.Errorf("get html: %w", err)
	}
	return result, nil
}

// QueryValue returns the value property of the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryValue(ctx context.Context, sessionName, tabName, arg string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get value")
	if err != nil {
		return "", err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return "", err
	}
	if backendNodeID > 0 {
		result, err := QueryValueByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return "", fmt.Errorf("get value: %w", err)
		}
		return result, nil
	}
	result, err := QueryValueTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return "", fmt.Errorf("get value: %w", err)
	}
	return result, nil
}

// QueryAttr returns the value of attr on the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryAttr(ctx context.Context, sessionName, tabName, arg, attr string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get attr")
	if err != nil {
		return "", err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return "", err
	}
	if backendNodeID > 0 {
		result, err := QueryAttrByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID, attr)
		if err != nil {
			return "", fmt.Errorf("get attr: %w", err)
		}
		return result, nil
	}
	result, err := QueryAttrTarget(ctx, rt.DebugURL, rt.TargetID, selector, attr)
	if err != nil {
		return "", fmt.Errorf("get attr: %w", err)
	}
	return result, nil
}

// QueryTitle returns the document.title of the current page.
func (m *Manager) QueryTitle(ctx context.Context, sessionName, tabName string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get title")
	if err != nil {
		return "", err
	}
	result, err := QueryTitleTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return "", fmt.Errorf("get title: %w", err)
	}
	return result, nil
}

// QueryURL returns the current URL of the tracked tab.
func (m *Manager) QueryURL(ctx context.Context, sessionName, tabName string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get url")
	if err != nil {
		return "", err
	}
	result, err := QueryURLTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return "", fmt.Errorf("get url: %w", err)
	}
	return result, nil
}

// QueryCount returns the number of elements matching sel in the current page.
func (m *Manager) QueryCount(ctx context.Context, sessionName, tabName, sel string) (int, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get count")
	if err != nil {
		return 0, err
	}
	result, err := QueryCountTarget(ctx, rt.DebugURL, rt.TargetID, sel)
	if err != nil {
		return 0, fmt.Errorf("get count: %w", err)
	}
	return result, nil
}

// QueryBox returns the bounding box of the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryBox(ctx context.Context, sessionName, tabName, arg string) (*BoundingBox, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get box")
	if err != nil {
		return nil, err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return nil, err
	}
	if backendNodeID > 0 {
		result, err := QueryBoxByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return nil, fmt.Errorf("get box: %w", err)
		}
		return result, nil
	}
	result, err := QueryBoxTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return nil, fmt.Errorf("get box: %w", err)
	}
	return result, nil
}

// QueryStyles returns the computed styles of the element identified by arg (CSS selector or @eN ref).
func (m *Manager) QueryStyles(ctx context.Context, sessionName, tabName, arg string) (map[string]string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "get styles")
	if err != nil {
		return nil, err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return nil, err
	}
	if backendNodeID > 0 {
		result, err := QueryStylesByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return nil, fmt.Errorf("get styles: %w", err)
		}
		return result, nil
	}
	result, err := QueryStylesTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return nil, fmt.Errorf("get styles: %w", err)
	}
	return result, nil
}

// QueryVisible returns true if the element identified by arg (CSS selector or @eN ref) is visible.
func (m *Manager) QueryVisible(ctx context.Context, sessionName, tabName, arg string) (bool, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "is visible")
	if err != nil {
		return false, err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return false, err
	}
	if backendNodeID > 0 {
		result, err := QueryVisibleByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return false, fmt.Errorf("is visible: %w", err)
		}
		return result, nil
	}
	result, err := QueryVisibleTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return false, fmt.Errorf("is visible: %w", err)
	}
	return result, nil
}

// QueryEnabled returns true if the element identified by arg (CSS selector or @eN ref) is not disabled.
func (m *Manager) QueryEnabled(ctx context.Context, sessionName, tabName, arg string) (bool, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "is enabled")
	if err != nil {
		return false, err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return false, err
	}
	if backendNodeID > 0 {
		result, err := QueryEnabledByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return false, fmt.Errorf("is enabled: %w", err)
		}
		return result, nil
	}
	result, err := QueryEnabledTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return false, fmt.Errorf("is enabled: %w", err)
	}
	return result, nil
}

// QueryChecked returns true if the element identified by arg (CSS selector or @eN ref) is checked.
func (m *Manager) QueryChecked(ctx context.Context, sessionName, tabName, arg string) (bool, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "is checked")
	if err != nil {
		return false, err
	}
	selector, backendNodeID, err := m.resolveElementArg(ctx, rt, arg)
	if err != nil {
		return false, err
	}
	if backendNodeID > 0 {
		result, err := QueryCheckedByBackendNodeID(ctx, rt.DebugURL, rt.TargetID, backendNodeID)
		if err != nil {
			return false, fmt.Errorf("is checked: %w", err)
		}
		return result, nil
	}
	result, err := QueryCheckedTarget(ctx, rt.DebugURL, rt.TargetID, selector)
	if err != nil {
		return false, fmt.Errorf("is checked: %w", err)
	}
	return result, nil
}
