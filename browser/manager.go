package browser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// SessionOptions holds optional settings for session creation.
type SessionOptions struct {
	// Headless controls whether a local browser launches headlessly.
	Headless bool
	// WSURL is the remote browser WebSocket endpoint (remote mode only).
	WSURL string
}

// Manager coordinates browser session lifecycle, tab management, and
// browser actions using the metadata store and CDP transport layer.
type Manager struct {
	store *Store
}

// NewManager creates a session manager backed by the given store.
func NewManager(store *Store) *Manager {
	return &Manager{store: store}
}

// ---------------------------------------------------------------------------
// Session lifecycle
// ---------------------------------------------------------------------------

// CreateSession launches or connects to a browser and persists session metadata.
func (m *Manager) CreateSession(ctx context.Context, name string, mode Mode, opts SessionOptions) error {
	if err := ValidateSessionName(name); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	now := time.Now()

	switch mode {
	case ModeLocal:
		profileDir := filepath.Join(m.store.Root(), "profiles", name)

		proc, err := LaunchBrowser(ctx, LocalConfig{ProfileDir: profileDir, Headless: opts.Headless})
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		session, err := NewLocalSession(name, profileDir, opts.Headless, now)
		if err != nil {
			// Best-effort cleanup of the just-launched browser.
			_ = KillProcess(proc)
			return fmt.Errorf("create session: %w", err)
		}
		session.Process = proc

		if err := m.store.Update(func(state *State) error {
			return state.CreateSession(session)
		}); err != nil {
			_ = KillProcess(proc)
			return fmt.Errorf("create session: %w", err)
		}

	case ModeRemote:
		// Validate the remote endpoint is reachable.
		httpURL, err := debugURLToHTTP(opts.WSURL)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpURL+"/json/version", nil)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("create session: remote endpoint unreachable: %w", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("create session: remote endpoint returned status %d", resp.StatusCode)
		}

		session, err := NewRemoteSession(name, opts.WSURL, now)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		session.Process = &ProcessRecord{DebugURL: opts.WSURL}

		if err := m.store.Update(func(state *State) error {
			return state.CreateSession(session)
		}); err != nil {
			return fmt.Errorf("create session: %w", err)
		}

	default:
		return fmt.Errorf("create session: unsupported mode %q", mode)
	}

	return nil
}

// CloseSession terminates a browser session, kills the local process if
// applicable, and removes all related metadata.
func (m *Manager) CloseSession(_ context.Context, name string) error {
	resolved, err := m.resolveSessionName(name)
	if err != nil {
		return fmt.Errorf("close session: %w", err)
	}
	name = resolved
	return m.store.WithSessionLock(name, func() error {
		// Phase 1: read process info under state lock.
		var proc *ProcessRecord
		var profileDir string
		var isLocal bool
		err := m.store.Update(func(state *State) error {
			session, err := state.ResolveSession(name)
			if err != nil {
				return err
			}
			isLocal = session.Mode == ModeLocal
			if session.Process != nil {
				p := *session.Process
				proc = &p
			}
			if session.Local != nil {
				profileDir = session.Local.ProfileDir
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("close session: %w", err)
		}

		// Phase 2: kill the local browser process outside the state lock.
		if isLocal && proc != nil {
			// Best-effort: don't fail if the process is already dead.
			_ = KillProcess(proc)
			if profileDir != "" {
				_ = os.RemoveAll(profileDir)
			}
		}

		// Phase 3: atomically remove session from state.
		return m.store.Update(func(state *State) error {
			if err := state.DeleteSession(name); err != nil {
				return fmt.Errorf("close session: %w", err)
			}
			return nil
		})
	})
}

// ListSessions returns all tracked sessions sorted by name.
func (m *Manager) ListSessions(_ context.Context) ([]*SessionRecord, error) {
	state, err := m.store.Load()
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	sessions := make([]*SessionRecord, 0, len(state.Sessions))
	for _, s := range state.Sessions {
		sessions = append(sessions, s)
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Name < sessions[j].Name
	})
	return sessions, nil
}

// GetSession resolves a session by name (or falls back to the selected/only session).
func (m *Manager) GetSession(_ context.Context, name string) (*SessionRecord, error) {
	state, err := m.store.Load()
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	session, err := state.ResolveSession(name)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return session, nil
}

// SelectSession persists the default session used when --session is omitted.
func (m *Manager) SelectSession(_ context.Context, name string) error {
	if err := m.store.Update(func(state *State) error {
		return state.SelectSession(name)
	}); err != nil {
		return fmt.Errorf("select session: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Tab lifecycle
// ---------------------------------------------------------------------------

// CreateTab opens a new browser tab and tracks it in session metadata.
func (m *Manager) CreateTab(ctx context.Context, sessionName string, tabName string, url string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("create tab: %w", err)
	}
	sessionName = resolved

	// Phase 1: resolve debug URL under lock.
	var debugURL string
	err = m.store.UpdateSession(sessionName, func(_ *State, session *SessionRecord) error {
		du, err := resolveDebugURL(session)
		if err != nil {
			return fmt.Errorf("create tab: %w", err)
		}
		debugURL = du
		return nil
	})
	if err != nil {
		return err
	}

	// Phase 2: create CDP target outside the state lock.
	targetID, err := CreateTarget(ctx, debugURL, url)
	if err != nil {
		return fmt.Errorf("create tab: %w", err)
	}

	// Phase 3: persist tab metadata under lock.
	return m.store.UpdateSession(sessionName, func(state *State, _ *SessionRecord) error {
		now := time.Now()
		tab, err := NewTab(tabName, targetID, url, now)
		if err != nil {
			return fmt.Errorf("create tab: %w", err)
		}

		if err := state.UpsertTab(sessionName, tab); err != nil {
			return fmt.Errorf("create tab: %w", err)
		}
		return nil
	})
}

// CloseTab closes a browser tab and removes it from session metadata.
func (m *Manager) CloseTab(ctx context.Context, sessionName string, tabName string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("close tab: %w", err)
	}
	sessionName = resolved
	return m.store.UpdateSession(sessionName, func(state *State, session *SessionRecord) error {
		tab, err := session.ResolveTab(tabName)
		if err != nil {
			return fmt.Errorf("close tab: %w", err)
		}

		// Close the CDP target if the tab is still live.
		if tab.Status == TabStatusLive && tab.TargetID != "" {
			debugURL, err := resolveDebugURL(session)
			if err != nil {
				return fmt.Errorf("close tab: %w", err)
			}
			if err := CloseTarget(ctx, debugURL, tab.TargetID); err != nil {
				return fmt.Errorf("close tab: %w", err)
			}
		}

		if err := state.DeleteTab(sessionName, tab.Name); err != nil {
			return fmt.Errorf("close tab: %w", err)
		}
		return nil
	})
}

// ListTabs returns all tracked tabs for a session sorted by creation time.
func (m *Manager) ListTabs(_ context.Context, sessionName string) ([]*TabRecord, error) {
	state, err := m.store.Load()
	if err != nil {
		return nil, fmt.Errorf("list tabs: %w", err)
	}

	session, err := state.ResolveSession(sessionName)
	if err != nil {
		return nil, fmt.Errorf("list tabs: %w", err)
	}

	tabs := make([]*TabRecord, 0, len(session.Tabs))
	for _, t := range session.Tabs {
		tabs = append(tabs, t)
	}
	sort.Slice(tabs, func(i, j int) bool {
		if tabs[i].CreatedAt.Equal(tabs[j].CreatedAt) {
			return tabs[i].Name < tabs[j].Name
		}
		return tabs[i].CreatedAt.Before(tabs[j].CreatedAt)
	})
	return tabs, nil
}

// SelectTab persists the default tab used when --tab is omitted.
func (m *Manager) SelectTab(_ context.Context, sessionName string, tabName string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("select tab: %w", err)
	}
	sessionName = resolved
	return m.store.UpdateSession(sessionName, func(state *State, _ *SessionRecord) error {
		if err := state.SelectTab(sessionName, tabName); err != nil {
			return fmt.Errorf("select tab: %w", err)
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// Browser actions
// ---------------------------------------------------------------------------

// Navigate changes the URL of a tracked tab.
func (m *Manager) Navigate(ctx context.Context, sessionName string, tabName string, url string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("navigate: %w", err)
	}
	sessionName = resolved
	// Phase 1: resolve session/tab under lock, release before CDP I/O.
	var debugURL, targetID, resolvedSession, resolvedTab string
	err = m.store.UpdateSession(sessionName, func(_ *State, session *SessionRecord) error {
		tab, err := session.ResolveTab(tabName)
		if err != nil {
			return fmt.Errorf("navigate: %w", err)
		}
		if err := requireLiveTab(tab); err != nil {
			return fmt.Errorf("navigate: %w", err)
		}
		du, err := resolveDebugURL(session)
		if err != nil {
			return fmt.Errorf("navigate: %w", err)
		}
		debugURL = du
		targetID = tab.TargetID
		resolvedSession = session.Name
		resolvedTab = tab.Name
		return nil
	})
	if err != nil {
		return err
	}

	// Phase 2: CDP navigation outside any lock.
	if err := NavigateTarget(ctx, debugURL, targetID, url); err != nil {
		return fmt.Errorf("navigate: %w", err)
	}

	// Phase 3: persist URL update under lock.
	return m.store.UpdateSession(resolvedSession, func(_ *State, session *SessionRecord) error {
		if tab, ok := session.Tabs[resolvedTab]; ok {
			tab.URL = url
			tab.UpdatedAt = time.Now().UTC()
		}
		return nil
	})
}

// Evaluate runs JavaScript in a tracked tab and returns the result.
func (m *Manager) Evaluate(ctx context.Context, sessionName string, tabName string, js string) (any, error) {
	// Resolve session/tab under lock, then release before CDP I/O.
	debugURL, targetID, err := m.resolveTarget(sessionName, tabName, "evaluate")
	if err != nil {
		return nil, err
	}

	result, err := EvalTarget(ctx, debugURL, targetID, js)
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}
	return result, nil
}

// Screenshot captures a full-page PNG of a tracked tab.
func (m *Manager) Screenshot(ctx context.Context, sessionName string, tabName string) ([]byte, error) {
	// Resolve session/tab under lock, then release before CDP I/O.
	debugURL, targetID, err := m.resolveTarget(sessionName, tabName, "screenshot")
	if err != nil {
		return nil, err
	}

	buf, err := ScreenshotTarget(ctx, debugURL, targetID)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}
	return buf, nil
}

// ---------------------------------------------------------------------------
// Reconciliation
// ---------------------------------------------------------------------------

// Reconcile refreshes tab liveness by comparing metadata against live CDP targets.
func (m *Manager) Reconcile(ctx context.Context, sessionName string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("reconcile: %w", err)
	}
	sessionName = resolved
	return m.store.UpdateSession(sessionName, func(state *State, session *SessionRecord) error {
		debugURL, err := resolveDebugURL(session)
		if err != nil {
			return fmt.Errorf("reconcile: %w", err)
		}

		targets, err := ListTargets(ctx, debugURL)
		if err != nil {
			return fmt.Errorf("reconcile: %w", err)
		}

		targetIDs := make([]string, len(targets))
		for i, t := range targets {
			targetIDs[i] = t.TargetID
		}

		now := time.Now()
		if err := state.ReconcileSession(sessionName, targetIDs, now); err != nil {
			return fmt.Errorf("reconcile: %w", err)
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// resolveTarget resolves a session and tab under lock and returns the debug URL
// and target ID for use outside the lock during CDP I/O.
func (m *Manager) resolveTarget(sessionName string, tabName string, op string) (string, string, error) {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	sessionName = resolved
	var debugURL, targetID string
	err = m.store.WithSessionLock(sessionName, func() error {
		state, err := m.store.Load()
		if err != nil {
			return err
		}
		session, err := state.ResolveSession(sessionName)
		if err != nil {
			return err
		}
		tab, err := session.ResolveTab(tabName)
		if err != nil {
			return err
		}
		if err := requireLiveTab(tab); err != nil {
			return err
		}
		du, err := resolveDebugURL(session)
		if err != nil {
			return err
		}
		debugURL = du
		targetID = tab.TargetID
		return nil
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return debugURL, targetID, nil
}

// resolveSessionName resolves an optional session name to a concrete name
// before it reaches WithSessionLock or UpdateSession (which reject empty names).
func (m *Manager) resolveSessionName(name string) (string, error) {
	if name != "" {
		return name, nil
	}
	state, err := m.store.Load()
	if err != nil {
		return "", err
	}
	session, err := state.ResolveSession("")
	if err != nil {
		return "", err
	}
	return session.Name, nil
}

// resolveDebugURL extracts the CDP debug endpoint from session metadata.
func resolveDebugURL(session *SessionRecord) (string, error) {
	if session.Process != nil && session.Process.DebugURL != "" {
		return session.Process.DebugURL, nil
	}
	if session.Mode == ModeRemote && session.Remote != nil {
		return session.Remote.WSURL, nil
	}
	return "", fmt.Errorf("session %q has no debug endpoint", session.Name)
}

// requireLiveTab ensures a tab is in a usable state for browser actions.
func requireLiveTab(tab *TabRecord) error {
	if tab.Status != TabStatusLive {
		return fmt.Errorf("tab %q is %s, not live", tab.Name, tab.Status)
	}
	if tab.TargetID == "" {
		return fmt.Errorf("tab %q has no target ID", tab.Name)
	}
	return nil
}
