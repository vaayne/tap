package browser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
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

	// interceptMu guards interceptCancel for concurrent access.
	interceptMu sync.Mutex
	// interceptCancel tracks the active Fetch domain interception cancel func
	// per target (keyed by "session:tab"). When new rules are set, the previous
	// cancel is called first to prevent goroutine leaks.
	interceptCancel map[string]func()
}

// NewManager creates a session manager backed by the given store.
func NewManager(store *Store) *Manager {
	return &Manager{store: store, interceptCancel: make(map[string]func())}
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
		if err := checkDebugEndpoint(ctx, opts.WSURL); err != nil {
			return fmt.Errorf("create session: %w", err)
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
			// Only remove the profile after confirming the process is gone,
			// so we don't delete files Chrome is still writing to.
			if profileDir != "" && (proc.PID <= 0 || !isProcessAlive(proc.PID)) {
				removeProfileDir(profileDir)
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

// SessionList holds the result of listing sessions including the current selection.
type SessionList struct {
	Sessions        []*SessionRecord
	SelectedSession string
}

// ListSessions returns all tracked sessions sorted by name.
func (m *Manager) ListSessions(_ context.Context) (*SessionList, error) {
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
	return &SessionList{Sessions: sessions, SelectedSession: state.SelectedSession}, nil
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

	// Phase 1: resolve debug URL and check for duplicates under lock.
	var debugURL string
	err = m.store.UpdateSession(sessionName, func(_ *State, session *SessionRecord) error {
		if _, exists := session.Tabs[tabName]; exists {
			return fmt.Errorf("create tab: tab %q already exists in session %q", tabName, sessionName)
		}
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

	// Phase 3: persist tab metadata under lock. Re-check for duplicates
	// since another CreateTab could have raced between Phase 1 and Phase 3.
	// On any failure, best-effort close the orphaned CDP target.
	err = m.store.UpdateSession(sessionName, func(state *State, session *SessionRecord) error {
		if _, exists := session.Tabs[tabName]; exists {
			return fmt.Errorf("create tab: tab %q already exists in session %q", tabName, sessionName)
		}
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
	if err != nil {
		// Best-effort cleanup of the orphaned CDP target.
		_ = CloseTarget(ctx, debugURL, targetID)
		return err
	}
	return nil
}

// CloseTab closes a browser tab and removes it from session metadata.
func (m *Manager) CloseTab(ctx context.Context, sessionName string, tabName string) error {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return fmt.Errorf("close tab: %w", err)
	}
	sessionName = resolved

	// Phase 1: resolve tab info under lock.
	var debugURL, targetID, resolvedTab string
	var isLive bool
	err = m.store.UpdateSession(sessionName, func(_ *State, session *SessionRecord) error {
		tab, err := session.ResolveTab(tabName)
		if err != nil {
			return fmt.Errorf("close tab: %w", err)
		}
		resolvedTab = tab.Name
		isLive = tab.Status == TabStatusLive && tab.TargetID != ""
		if isLive {
			du, err := resolveDebugURL(session)
			if err != nil {
				return fmt.Errorf("close tab: %w", err)
			}
			debugURL = du
			targetID = tab.TargetID
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Phase 2: close the CDP target outside the lock.
	if isLive {
		if err := CloseTarget(ctx, debugURL, targetID); err != nil {
			return fmt.Errorf("close tab: %w", err)
		}
	}

	// Phase 3: remove tab metadata under lock.
	return m.store.UpdateSession(sessionName, func(state *State, _ *SessionRecord) error {
		if err := state.DeleteTab(sessionName, resolvedTab); err != nil {
			return fmt.Errorf("close tab: %w", err)
		}
		return nil
	})
}

// TabList holds the result of listing tabs including the current selection.
type TabList struct {
	Tabs        []*TabRecord
	SelectedTab string
}

// ListTabs returns all tracked tabs for a session sorted by creation time.
func (m *Manager) ListTabs(_ context.Context, sessionName string) (*TabList, error) {
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
	return &TabList{Tabs: tabs, SelectedTab: session.SelectedTab}, nil
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
		tab, ok := session.Tabs[resolvedTab]
		if !ok {
			// Tab was deleted between Phase 2 and Phase 3. Navigation
			// succeeded in the browser but we can't update metadata.
			return fmt.Errorf("navigate: tab %q was removed during navigation", resolvedTab)
		}
		tab.URL = url
		tab.UpdatedAt = time.Now().UTC()
		return nil
	})
}

// Evaluate runs JavaScript in a tracked tab and returns the result.
func (m *Manager) Evaluate(ctx context.Context, sessionName string, tabName string, js string) (any, error) {
	// Resolve session/tab under lock, then release before CDP I/O.
	rt, err := m.resolveTarget(sessionName, tabName, "evaluate")
	if err != nil {
		return nil, err
	}

	result, err := EvalTarget(ctx, rt.DebugURL, rt.TargetID, js)
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}
	return result, nil
}

// ScreenshotResult holds the screenshot data and resolved names.
type ScreenshotResult struct {
	Data        []byte
	SessionName string
	TabName     string
}

// Screenshot captures a full-page PNG of a tracked tab.
func (m *Manager) Screenshot(ctx context.Context, sessionName string, tabName string) (*ScreenshotResult, error) {
	// Resolve session/tab under lock, then release before CDP I/O.
	rt, err := m.resolveTarget(sessionName, tabName, "screenshot")
	if err != nil {
		return nil, err
	}

	buf, err := ScreenshotTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}
	return &ScreenshotResult{Data: buf, SessionName: rt.SessionName, TabName: rt.TabName}, nil
}

// Forms discovers fillable form elements in a tracked tab.
func (m *Manager) Forms(ctx context.Context, sessionName string, tabName string) ([]FormField, error) {
	rt, err := m.resolveTarget(sessionName, tabName, "forms")
	if err != nil {
		return nil, err
	}

	fields, err := FormsTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return nil, fmt.Errorf("forms: %w", err)
	}
	return fields, nil
}

// Fill sets values in form fields of a tracked tab.
func (m *Manager) Fill(ctx context.Context, sessionName string, tabName string, fields []FillField, submitSelector string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "fill")
	if err != nil {
		return err
	}

	if err := FillTarget(ctx, rt.DebugURL, rt.TargetID, fields, submitSelector); err != nil {
		return fmt.Errorf("fill: %w", err)
	}
	return nil
}

// PDFResult holds the PDF data and resolved names.
type PDFResult struct {
	Data        []byte
	SessionName string
	TabName     string
}

// PDF saves the current page as PDF from a tracked tab.
func (m *Manager) PDF(ctx context.Context, sessionName string, tabName string, landscape bool, printBackground bool, scale float64) (*PDFResult, error) {
	rt, err := m.resolveTarget(sessionName, tabName, "pdf")
	if err != nil {
		return nil, err
	}
	buf, err := PDFTarget(ctx, rt.DebugURL, rt.TargetID, landscape, printBackground, scale)
	if err != nil {
		return nil, fmt.Errorf("pdf: %w", err)
	}
	return &PDFResult{Data: buf, SessionName: rt.SessionName, TabName: rt.TabName}, nil
}

// Click dispatches a real mouse click on the element matching sel in a tracked tab.
func (m *Manager) Click(ctx context.Context, sessionName string, tabName string, sel string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "click")
	if err != nil {
		return err
	}
	if err := ClickTarget(ctx, rt.DebugURL, rt.TargetID, sel); err != nil {
		return fmt.Errorf("click: %w", err)
	}
	return nil
}

// Type sends individual key events to the element matching sel in a tracked tab.
func (m *Manager) Type(ctx context.Context, sessionName string, tabName string, sel string, text string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "type")
	if err != nil {
		return err
	}
	if err := TypeTarget(ctx, rt.DebugURL, rt.TargetID, sel, text); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	return nil
}

// Hover moves the mouse to the element matching sel in a tracked tab.
func (m *Manager) Hover(ctx context.Context, sessionName string, tabName string, sel string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "hover")
	if err != nil {
		return err
	}
	if err := HoverTarget(ctx, rt.DebugURL, rt.TargetID, sel); err != nil {
		return fmt.Errorf("hover: %w", err)
	}
	return nil
}

// Scroll scrolls to the element matching sel, or to absolute x,y if sel is empty.
func (m *Manager) Scroll(ctx context.Context, sessionName string, tabName string, sel string, x, y float64) error {
	rt, err := m.resolveTarget(sessionName, tabName, "scroll")
	if err != nil {
		return err
	}
	if err := ScrollTarget(ctx, rt.DebugURL, rt.TargetID, sel, x, y); err != nil {
		return fmt.Errorf("scroll: %w", err)
	}
	return nil
}

// Select selects an option by value in a <select> element in a tracked tab.
func (m *Manager) Select(ctx context.Context, sessionName string, tabName string, sel string, value string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "select")
	if err != nil {
		return err
	}
	if err := SelectTarget(ctx, rt.DebugURL, rt.TargetID, sel, value); err != nil {
		return fmt.Errorf("select: %w", err)
	}
	return nil
}

// WaitFor waits until the element matching sel is visible in a tracked tab.
func (m *Manager) WaitFor(ctx context.Context, sessionName string, tabName string, sel string, timeout time.Duration) error {
	rt, err := m.resolveTarget(sessionName, tabName, "wait")
	if err != nil {
		return err
	}
	if err := WaitForTarget(ctx, rt.DebugURL, rt.TargetID, sel, timeout); err != nil {
		return fmt.Errorf("wait: %w", err)
	}
	return nil
}

// Back navigates the tracked tab backwards in history.
func (m *Manager) Back(ctx context.Context, sessionName string, tabName string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "back")
	if err != nil {
		return err
	}
	if err := BackTarget(ctx, rt.DebugURL, rt.TargetID); err != nil {
		return fmt.Errorf("back: %w", err)
	}
	return nil
}

// Forward navigates the tracked tab forwards in history.
func (m *Manager) Forward(ctx context.Context, sessionName string, tabName string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "forward")
	if err != nil {
		return err
	}
	if err := ForwardTarget(ctx, rt.DebugURL, rt.TargetID); err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	return nil
}

// Reload reloads the current page in a tracked tab.
func (m *Manager) Reload(ctx context.Context, sessionName string, tabName string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "reload")
	if err != nil {
		return err
	}
	if err := ReloadTarget(ctx, rt.DebugURL, rt.TargetID); err != nil {
		return fmt.Errorf("reload: %w", err)
	}
	return nil
}

// Keypress sends key events to the page in a tracked tab.
func (m *Manager) Keypress(ctx context.Context, sessionName string, tabName string, keys string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "keypress")
	if err != nil {
		return err
	}
	if err := KeypressTarget(ctx, rt.DebugURL, rt.TargetID, keys); err != nil {
		return fmt.Errorf("keypress: %w", err)
	}
	return nil
}

// Dialog accepts or dismisses a pending JavaScript dialog in a tracked tab.
func (m *Manager) Dialog(ctx context.Context, sessionName string, tabName string, accept bool, promptText string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "dialog")
	if err != nil {
		return err
	}
	if err := DialogTarget(ctx, rt.DebugURL, rt.TargetID, accept, promptText); err != nil {
		return fmt.Errorf("dialog: %w", err)
	}
	return nil
}

// GetCookies returns all cookies for the current page in a tracked tab.
func (m *Manager) GetCookies(ctx context.Context, sessionName string, tabName string) ([]CookieEntry, error) {
	rt, err := m.resolveTarget(sessionName, tabName, "cookies get")
	if err != nil {
		return nil, err
	}
	cookies, err := GetCookiesTarget(ctx, rt.DebugURL, rt.TargetID)
	if err != nil {
		return nil, fmt.Errorf("cookies get: %w", err)
	}
	return cookies, nil
}

// SetCookie sets a cookie in a tracked tab.
func (m *Manager) SetCookie(ctx context.Context, sessionName string, tabName string, name, value, domain, path string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "cookies set")
	if err != nil {
		return err
	}
	if err := SetCookieTarget(ctx, rt.DebugURL, rt.TargetID, name, value, domain, path); err != nil {
		return fmt.Errorf("cookies set: %w", err)
	}
	return nil
}

// ClearCookies deletes all cookies for the current page in a tracked tab.
func (m *Manager) ClearCookies(ctx context.Context, sessionName string, tabName string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "cookies clear")
	if err != nil {
		return err
	}
	if err := ClearCookiesTarget(ctx, rt.DebugURL, rt.TargetID); err != nil {
		return fmt.Errorf("cookies clear: %w", err)
	}
	return nil
}

// NetworkWait blocks until a network request matching the filter completes in a
// tracked tab. If includeBody is true, the response body is fetched before
// returning. The caller controls the timeout via ctx.
func (m *Manager) NetworkWait(ctx context.Context, sessionName string, tabName string, filter NetworkFilter, includeBody bool) (*NetworkEntry, error) {
	rt, err := m.resolveTarget(sessionName, tabName, "network wait")
	if err != nil {
		return nil, err
	}

	entry, err := WaitForRequest(ctx, rt.DebugURL, rt.TargetID, filter, includeBody)
	if err != nil {
		return nil, fmt.Errorf("network wait: %w", err)
	}
	return entry, nil
}

// NetworkGetBody fetches the response body for a completed request by its
// request ID from a tracked tab.
func (m *Manager) NetworkGetBody(ctx context.Context, sessionName string, tabName string, requestID string) ([]byte, error) {
	rt, err := m.resolveTarget(sessionName, tabName, "network body")
	if err != nil {
		return nil, err
	}

	body, err := GetResponseBody(ctx, rt.DebugURL, rt.TargetID, requestID)
	if err != nil {
		return nil, fmt.Errorf("network body: %w", err)
	}
	return body, nil
}

// NetworkLog starts capturing network requests for a tracked tab and streams
// completed entries to the returned channel. Call cancel to stop capturing.
func (m *Manager) NetworkLog(ctx context.Context, sessionName string, tabName string, filter NetworkFilter) (<-chan NetworkEntry, func(), error) {
	rt, err := m.resolveTarget(sessionName, tabName, "network log")
	if err != nil {
		return nil, nil, err
	}

	ch, cancel, err := EnableNetworkLog(ctx, rt.DebugURL, rt.TargetID, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("network log: %w", err)
	}
	return ch, cancel, nil
}

// interceptKey returns the map key for tracking intercept cancel funcs.
func interceptKey(session, tab string) string {
	return session + ":" + tab
}

// cancelIntercept cancels and removes any active interception for the given key.
func (m *Manager) cancelIntercept(key string) {
	m.interceptMu.Lock()
	if prev, ok := m.interceptCancel[key]; ok {
		prev()
		delete(m.interceptCancel, key)
	}
	m.interceptMu.Unlock()
}

// NetworkIntercept sets Fetch domain interception rules on a tracked tab.
// Replaces any previously set rules (cancels the old interception goroutine).
func (m *Manager) NetworkIntercept(ctx context.Context, sessionName string, tabName string, rules []InterceptRule) error {
	rt, err := m.resolveTarget(sessionName, tabName, "network intercept")
	if err != nil {
		return err
	}

	key := interceptKey(rt.SessionName, rt.TabName)
	m.cancelIntercept(key)

	cancel, err := SetInterceptRules(ctx, rt.DebugURL, rt.TargetID, rules)
	if err != nil {
		return fmt.Errorf("network intercept: %w", err)
	}
	m.interceptMu.Lock()
	m.interceptCancel[key] = cancel
	m.interceptMu.Unlock()
	return nil
}

// NetworkClearIntercept removes all Fetch domain interception rules from a
// tracked tab.
func (m *Manager) NetworkClearIntercept(ctx context.Context, sessionName string, tabName string) error {
	rt, err := m.resolveTarget(sessionName, tabName, "network clear")
	if err != nil {
		return err
	}

	m.cancelIntercept(interceptKey(rt.SessionName, rt.TabName))

	if err := ClearIntercept(ctx, rt.DebugURL, rt.TargetID); err != nil {
		return fmt.Errorf("network clear: %w", err)
	}
	return nil
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

	// Phase 1: resolve debug URL under lock.
	var debugURL string
	err = m.store.UpdateSession(sessionName, func(_ *State, session *SessionRecord) error {
		du, err := resolveDebugURL(session)
		if err != nil {
			return fmt.Errorf("reconcile: %w", err)
		}
		debugURL = du
		return nil
	})
	if err != nil {
		return err
	}

	// Phase 2: list live CDP targets outside the lock.
	targets, err := ListTargets(ctx, debugURL)
	if err != nil {
		return fmt.Errorf("reconcile: %w", err)
	}

	targetIDs := make([]string, len(targets))
	for i, t := range targets {
		targetIDs[i] = t.TargetID
	}

	// Phase 3: reconcile metadata under lock.
	now := time.Now()
	return m.store.UpdateSession(sessionName, func(state *State, _ *SessionRecord) error {
		if err := state.ReconcileSession(sessionName, targetIDs, now); err != nil {
			return fmt.Errorf("reconcile: %w", err)
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// resolvedTarget holds the result of resolving a session and tab for CDP I/O.
type resolvedTarget struct {
	DebugURL    string
	TargetID    string
	SessionName string
	TabName     string
}

// resolveTarget resolves a session and tab under lock and returns the debug URL,
// target ID, and resolved names for use outside the lock during CDP I/O.
func (m *Manager) resolveTarget(sessionName string, tabName string, op string) (resolvedTarget, error) {
	resolved, err := m.resolveSessionName(sessionName)
	if err != nil {
		return resolvedTarget{}, fmt.Errorf("%s: %w", op, err)
	}
	sessionName = resolved
	var rt resolvedTarget
	rt.SessionName = sessionName
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
		rt.DebugURL = du
		rt.TargetID = tab.TargetID
		rt.TabName = tab.Name
		return nil
	})
	if err != nil {
		return resolvedTarget{}, fmt.Errorf("%s: %w", op, err)
	}
	return rt, nil
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

// removeProfileDir removes a Chrome profile directory with a retry loop to
// handle transient file locks (e.g., antivirus or indexing services on Windows).
func removeProfileDir(dir string) {
	for range 3 {
		if err := os.RemoveAll(dir); err == nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	// Final best-effort attempt.
	_ = os.RemoveAll(dir)
}
