package browser

import (
	"context"
	"fmt"
	"time"
)

// SetEmulation applies emulation settings to a tab immediately (via CDP) and
// persists them so they are re-applied on every subsequent resolveTarget call.
// Only the fields present in delta are written; pass a full EmulationSettings
// to replace all settings at once.
func (m *Manager) SetEmulation(ctx context.Context, sessionName, tabName string, delta *EmulationSettings) error {
	rt, err := m.resolveTarget(ctx, sessionName, tabName, "set emulation")
	if err != nil {
		return err
	}

	// Apply immediately via CDP.
	if err := ApplyEmulationTarget(ctx, rt.DebugURL, rt.TargetID, delta); err != nil {
		return fmt.Errorf("set emulation: %w", err)
	}

	// Persist by merging delta into the stored EmulationSettings.
	resolved, err := m.resolveSessionName(ctx, sessionName, false)
	if err != nil {
		return fmt.Errorf("set emulation: %w", err)
	}
	return m.store.UpdateSession(resolved, func(_ *State, session *SessionRecord) error {
		tab, ok := session.Tabs[rt.TabName]
		if !ok {
			return fmt.Errorf("set emulation: tab %q not found", rt.TabName)
		}
		if tab.Emulation == nil {
			tab.Emulation = &EmulationSettings{}
		}
		mergeEmulation(tab.Emulation, delta)
		tab.UpdatedAt = time.Now().UTC()
		return nil
	})
}

// ClearEmulation wipes all persisted emulation settings for the tab.
// It does not reverse already-applied CDP overrides (a page reload will clear them).
func (m *Manager) ClearEmulation(ctx context.Context, sessionName, tabName string) error {
	resolved, err := m.resolveSessionName(ctx, sessionName, false)
	if err != nil {
		return fmt.Errorf("clear emulation: %w", err)
	}

	// We need the resolved tab name — load from store.
	state, err := m.store.Load()
	if err != nil {
		return fmt.Errorf("clear emulation: %w", err)
	}
	session, err := state.ResolveSessionByPreference(resolved)
	if err != nil {
		return fmt.Errorf("clear emulation: %w", err)
	}
	tab, err := session.ResolveTab(tabName)
	if err != nil {
		return fmt.Errorf("clear emulation: %w", err)
	}
	resolvedTab := tab.Name

	return m.store.UpdateSession(resolved, func(_ *State, s *SessionRecord) error {
		t, ok := s.Tabs[resolvedTab]
		if !ok {
			return fmt.Errorf("clear emulation: tab %q not found", resolvedTab)
		}
		t.Emulation = nil
		t.UpdatedAt = time.Now().UTC()
		return nil
	})
}

// mergeEmulation writes non-zero fields from src into dst.
func mergeEmulation(dst, src *EmulationSettings) {
	if src == nil {
		return
	}
	if src.ViewportWidth != 0 {
		dst.ViewportWidth = src.ViewportWidth
	}
	if src.ViewportHeight != 0 {
		dst.ViewportHeight = src.ViewportHeight
	}
	if src.ViewportScale != 0 {
		dst.ViewportScale = src.ViewportScale
	}
	if src.DeviceName != "" {
		dst.DeviceName = src.DeviceName
		// A device preset implies viewport + UA; clear any independent overrides
		// that were previously set so re-apply is consistent.
		dst.ViewportWidth = 0
		dst.ViewportHeight = 0
		dst.ViewportScale = 0
		dst.UserAgent = ""
	}
	if src.GeoLat != nil {
		dst.GeoLat = src.GeoLat
	}
	if src.GeoLng != nil {
		dst.GeoLng = src.GeoLng
	}
	if src.Offline != nil {
		dst.Offline = src.Offline
	}
	if len(src.Headers) > 0 {
		if dst.Headers == nil {
			dst.Headers = make(map[string]string, len(src.Headers))
		}
		for k, v := range src.Headers {
			dst.Headers[k] = v
		}
	}
	if src.MediaScheme != "" {
		dst.MediaScheme = src.MediaScheme
	}
	if src.UserAgent != "" {
		dst.UserAgent = src.UserAgent
	}
}
