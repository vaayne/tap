package browser

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// Web storage operations on Manager
// ---------------------------------------------------------------------------

// GetStorageAll returns all entries from the named storage in a tracked tab.
// storeName must be "localStorage" or "sessionStorage".
func (m *Manager) GetStorageAll(ctx context.Context, sessionName, tabName, storeName string) (map[string]string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "storage get-all")
	if err != nil {
		return nil, err
	}
	result, err := GetStorageAll(ctx, rt.DebugURL, rt.TargetID, storeName)
	if err != nil {
		return nil, fmt.Errorf("storage get-all: %w", err)
	}
	return result, nil
}

// GetStorageKey returns the value for a single key from the named storage.
func (m *Manager) GetStorageKey(ctx context.Context, sessionName, tabName, storeName, key string) (string, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "storage get")
	if err != nil {
		return "", err
	}
	val, err := GetStorageKey(ctx, rt.DebugURL, rt.TargetID, storeName, key)
	if err != nil {
		return "", fmt.Errorf("storage get: %w", err)
	}
	return val, nil
}

// SetStorageKey sets a key/value pair in the named storage.
func (m *Manager) SetStorageKey(ctx context.Context, sessionName, tabName, storeName, key, value string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "storage set")
	if err != nil {
		return err
	}
	if err := SetStorageKey(ctx, rt.DebugURL, rt.TargetID, storeName, key, value); err != nil {
		return fmt.Errorf("storage set: %w", err)
	}
	return nil
}

// ClearStorage removes all entries from the named storage.
func (m *Manager) ClearStorage(ctx context.Context, sessionName, tabName, storeName string) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "storage clear")
	if err != nil {
		return err
	}
	if err := ClearStorage(ctx, rt.DebugURL, rt.TargetID, storeName); err != nil {
		return fmt.Errorf("storage clear: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Auth state save/load
// ---------------------------------------------------------------------------

// SaveState exports all cookies (full browser context) and the localStorage of
// the current tab's origin into a Playwright storageState-compatible structure.
func (m *Manager) SaveState(ctx context.Context, sessionName, tabName string) (*StorageState, error) {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "state save")
	if err != nil {
		return nil, err
	}

	// Get all cookies (full browser context, not just current page).
	cookies, err := GetAllCookiesTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return nil, fmt.Errorf("state save: %w", err)
	}

	// Determine the origin of the current page.
	origin, err := GetCurrentOrigin(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return nil, fmt.Errorf("state save: %w", err)
	}

	// Capture localStorage for the current origin.
	lsMap, err := GetStorageAll(ctx, rt.DebugURL, rt.TargetID, "localStorage")
	if err != nil {
		return nil, fmt.Errorf("state save: %w", err)
	}

	lsEntries := make([]StorageEntry, 0, len(lsMap))
	for k, v := range lsMap {
		lsEntries = append(lsEntries, StorageEntry{Name: k, Value: v})
	}

	state := &StorageState{
		Cookies: cookies,
	}
	if len(lsEntries) > 0 || origin != "" {
		state.Origins = []OriginStorage{{
			Origin:       origin,
			LocalStorage: lsEntries,
		}}
	}
	return state, nil
}

// LoadState imports cookies and localStorage from a saved StorageState.
//
// Limitations:
//   - All cookies are applied globally regardless of the current page.
//   - localStorage entries are only restored for the origin matching the
//     current page. Origins that don't match the current page are skipped
//     with a warning (they can't be written cross-origin from JS).
func (m *Manager) LoadState(ctx context.Context, sessionName, tabName string, state *StorageState, warnFn func(string)) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "state load")
	if err != nil {
		return err
	}

	// Apply all cookies globally.
	if len(state.Cookies) > 0 {
		if err := SetAllCookiesTarget(ctx, rt.DebugURL, rt.TargetID, state.Cookies); err != nil {
			return fmt.Errorf("state load: %w", err)
		}
	}

	if len(state.Origins) == 0 {
		return nil
	}

	// Determine the current page origin so we can match entries.
	currentOrigin, err := GetCurrentOrigin(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return fmt.Errorf("state load: %w", err)
	}

	for _, os := range state.Origins {
		if os.Origin != currentOrigin {
			if warnFn != nil {
				warnFn(fmt.Sprintf("skipping localStorage for origin %q (current page is %q)", os.Origin, currentOrigin))
			}
			continue
		}
		for _, entry := range os.LocalStorage {
			if err := SetStorageKey(ctx, rt.DebugURL, rt.TargetID, "localStorage", entry.Name, entry.Value); err != nil {
				return fmt.Errorf("state load: localStorage[%q]: %w", entry.Name, err)
			}
		}
	}
	return nil
}
