package browser

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	// StateVersion is the current on-disk schema version for browser metadata.
	StateVersion = 1

	// EnvStateRoot overrides the default durable state directory.
	EnvStateRoot = "TAP_BROWSER_STATE_DIR"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,62}$`)

var (
	ErrNoSessions        = errors.New("no browser sessions found")
	ErrAmbiguousSession  = errors.New("browser session selection is ambiguous")
	ErrSessionNotFound   = errors.New("browser session not found")
	ErrNoTabs            = errors.New("no tracked tabs found")
	ErrAmbiguousTab      = errors.New("browser tab selection is ambiguous")
	ErrTabNotFound       = errors.New("browser tab not found")
	ErrClosedTabSelected = errors.New("browser tab is closed")
	ErrStaleTabSelected  = errors.New("browser tab is stale")
)

// Mode identifies how a session connects to Chrome.
type Mode string

const (
	ModeLocal  Mode = "local"
	ModeRemote Mode = "remote"
)

// Operation is a user-facing browser action supported by the session manager.
type Operation string

const (
	OperationSessionNew   Operation = "session_new"
	OperationSessionClose Operation = "session_close"
	OperationTabNew       Operation = "tab_new"
	OperationTabClose     Operation = "tab_close"
	OperationNavigate     Operation = "navigate"
	OperationEvaluate     Operation = "evaluate"
	OperationScreenshot   Operation = "screenshot"
)

// TabStatus tracks whether a named tab can be used directly.
type TabStatus string

const (
	TabStatusLive   TabStatus = "live"
	TabStatusStale  TabStatus = "stale"
	TabStatusClosed TabStatus = "closed"
)

// Capability describes whether an operation is supported for a session mode.
type Capability struct {
	Supported bool   `json:"supported"`
	Notes     string `json:"notes,omitempty"`
}

// LocalConfig freezes local launch settings into session metadata.
type LocalConfig struct {
	ProfileDir string `json:"profile_dir"`
	Headless   bool   `json:"headless"`
}

// RemoteConfig freezes the remote reconnect endpoint into session metadata.
type RemoteConfig struct {
	WSURL string `json:"ws_url"`
}

// ProcessRecord stores the launch markers needed for later ownership checks.
type ProcessRecord struct {
	PID            int       `json:"pid,omitempty"`
	DebugURL       string    `json:"debug_url,omitempty"`
	OwnershipToken string    `json:"ownership_token,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
}

// TabRecord stores the durable metadata for a tracked browser target.
type TabRecord struct {
	Name       string    `json:"name"`
	TargetID   string    `json:"target_id,omitempty"`
	URL        string    `json:"url,omitempty"`
	Status     TabStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
}

// SessionRecord stores one persistent browser session and its tracked tabs.
type SessionRecord struct {
	Name             string                   `json:"name"`
	Mode             Mode                     `json:"mode"`
	Local            *LocalConfig             `json:"local,omitempty"`
	Remote           *RemoteConfig            `json:"remote,omitempty"`
	Process          *ProcessRecord           `json:"process,omitempty"`
	SelectedTab      string                   `json:"selected_tab,omitempty"`
	Tabs             map[string]*TabRecord    `json:"tabs,omitempty"`
	Capabilities     map[Operation]Capability `json:"capabilities,omitempty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	LastReconciledAt time.Time                `json:"last_reconciled_at,omitempty"`
}

// State is the durable browser session metadata written to disk.
type State struct {
	Version         int                       `json:"version"`
	SelectedSession string                    `json:"selected_session,omitempty"`
	Sessions        map[string]*SessionRecord `json:"sessions"`
}

// NewState returns an empty metadata state with initialized maps.
func NewState() *State {
	return &State{
		Version:  StateVersion,
		Sessions: make(map[string]*SessionRecord),
	}
}

var localCapabilities = map[Operation]Capability{
	OperationSessionNew: {
		Supported: true,
		Notes:     "Launch a managed local Chrome with a dedicated profile and debugging endpoint.",
	},
	OperationSessionClose: {
		Supported: true,
		Notes:     "Terminate the managed browser after ownership checks, then remove local metadata and profile state.",
	},
	OperationTabNew: {
		Supported: true,
		Notes:     "Create a tracked target within the managed local browser.",
	},
	OperationTabClose: {
		Supported: true,
		Notes:     "Close the tracked target and drop its metadata.",
	},
	OperationNavigate: {
		Supported: true,
		Notes:     "Navigate an explicit or selected tracked tab.",
	},
	OperationEvaluate: {
		Supported: true,
		Notes:     "Run JavaScript on an explicit or selected tracked tab.",
	},
	OperationScreenshot: {
		Supported: true,
		Notes:     "Capture a screenshot from an explicit or selected tracked tab.",
	},
}

var remoteCapabilities = map[Operation]Capability{
	OperationSessionNew: {
		Supported: true,
		Notes:     "Bind the session to an explicit --ws-url and validate connection/auth/TLS up front.",
	},
	OperationSessionClose: {
		Supported: true,
		Notes:     "Remove tap metadata only. Remote browser processes are never terminated by tap.",
	},
	OperationTabNew: {
		Supported: true,
		Notes:     "Supported when the remote endpoint allows target creation; failures must be returned clearly.",
	},
	OperationTabClose: {
		Supported: true,
		Notes:     "Supported when the remote endpoint allows target closure; tap removes metadata only after confirmation.",
	},
	OperationNavigate: {
		Supported: true,
		Notes:     "Operate through the persisted session endpoint, ignoring later global --ws-url overrides.",
	},
	OperationEvaluate: {
		Supported: true,
		Notes:     "Operate through the persisted session endpoint, ignoring later global --ws-url overrides.",
	},
	OperationScreenshot: {
		Supported: true,
		Notes:     "Operate through the persisted session endpoint, ignoring later global --ws-url overrides.",
	},
}

// CapabilitiesForMode returns the documented capability contract for the given mode.
func CapabilitiesForMode(mode Mode) map[Operation]Capability {
	switch mode {
	case ModeLocal:
		return localCapabilities
	case ModeRemote:
		return remoteCapabilities
	default:
		return nil
	}
}

// NewLocalSession builds a validated local session record.
func NewLocalSession(name string, profileDir string, headless bool, now time.Time) (*SessionRecord, error) {
	if err := ValidateSessionName(name); err != nil {
		return nil, err
	}
	if strings.TrimSpace(profileDir) == "" {
		return nil, errors.New("local browser profile directory is required")
	}
	return &SessionRecord{
		Name:         name,
		Mode:         ModeLocal,
		Local:        &LocalConfig{ProfileDir: profileDir, Headless: headless},
		Tabs:         make(map[string]*TabRecord),
		Capabilities: CapabilitiesForMode(ModeLocal),
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
	}, nil
}

// NewRemoteSession builds a validated remote session record.
func NewRemoteSession(name string, wsURL string, now time.Time) (*SessionRecord, error) {
	if err := ValidateSessionName(name); err != nil {
		return nil, err
	}
	if strings.TrimSpace(wsURL) == "" {
		return nil, errors.New("remote browser WebSocket URL is required")
	}
	return &SessionRecord{
		Name:         name,
		Mode:         ModeRemote,
		Remote:       &RemoteConfig{WSURL: wsURL},
		Tabs:         make(map[string]*TabRecord),
		Capabilities: CapabilitiesForMode(ModeRemote),
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
	}, nil
}

// NewTab builds a validated tracked tab record.
func NewTab(name string, targetID string, url string, now time.Time) (*TabRecord, error) {
	if err := ValidateTabName(name); err != nil {
		return nil, err
	}
	return &TabRecord{
		Name:       name,
		TargetID:   strings.TrimSpace(targetID),
		URL:        strings.TrimSpace(url),
		Status:     TabStatusLive,
		CreatedAt:  now.UTC(),
		UpdatedAt:  now.UTC(),
		LastSeenAt: now.UTC(),
	}, nil
}

// ValidateSessionName checks the user-facing session naming rules.
func ValidateSessionName(name string) error {
	return validateName("session", name)
}

// ValidateTabName checks the user-facing tab naming rules.
func ValidateTabName(name string) error {
	return validateName("tab", name)
}

func validateName(kind string, name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("%s name is required", kind)
	}
	if trimmed != name {
		return fmt.Errorf("%s name cannot have leading or trailing whitespace", kind)
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("%s name must match %q", kind, namePattern.String())
	}
	return nil
}

// Normalize repairs nil maps and empty versions loaded from disk.
func (s *State) Normalize() {
	if s.Version == 0 {
		s.Version = StateVersion
	}
	if s.Sessions == nil {
		s.Sessions = make(map[string]*SessionRecord)
	}
	for name, session := range s.Sessions {
		if session == nil {
			delete(s.Sessions, name)
			continue
		}
		if session.Name == "" {
			session.Name = name
		}
		if session.Tabs == nil {
			session.Tabs = make(map[string]*TabRecord)
		}
		for tabName, tab := range session.Tabs {
			if tab == nil {
				delete(session.Tabs, tabName)
				continue
			}
			if tab.Name == "" {
				tab.Name = tabName
			}
			if tab.Status == "" {
				tab.Status = TabStatusLive
			}
		}
	}
}

// Validate checks that the in-memory state satisfies the Phase 1 metadata contract.
func (s *State) Validate() error {
	if s.Version != StateVersion {
		return fmt.Errorf("unsupported browser state version %d", s.Version)
	}
	if s.SelectedSession != "" {
		if _, ok := s.Sessions[s.SelectedSession]; !ok {
			return fmt.Errorf("selected session %q is missing", s.SelectedSession)
		}
	}
	for name, session := range s.Sessions {
		if err := ValidateSessionName(name); err != nil {
			return err
		}
		if session.Name != name {
			return fmt.Errorf("session %q must match map key", name)
		}
		if err := session.validate(); err != nil {
			return fmt.Errorf("session %q: %w", name, err)
		}
	}
	return nil
}

func (s *SessionRecord) validate() error {
	if err := ValidateSessionName(s.Name); err != nil {
		return err
	}
	switch s.Mode {
	case ModeLocal:
		if s.Local == nil {
			return errors.New("local session config is required")
		}
		if s.Remote != nil {
			return errors.New("local session cannot have remote config")
		}
		if strings.TrimSpace(s.Local.ProfileDir) == "" {
			return errors.New("local session profile directory is required")
		}
	case ModeRemote:
		if s.Remote == nil {
			return errors.New("remote session config is required")
		}
		if s.Local != nil {
			return errors.New("remote session cannot have local config")
		}
		if strings.TrimSpace(s.Remote.WSURL) == "" {
			return errors.New("remote session WebSocket URL is required")
		}
	default:
		return fmt.Errorf("unsupported session mode %q", s.Mode)
	}
	if s.Tabs == nil {
		s.Tabs = make(map[string]*TabRecord)
	}
	if s.SelectedTab != "" {
		tab, ok := s.Tabs[s.SelectedTab]
		if !ok {
			return fmt.Errorf("selected tab %q is missing", s.SelectedTab)
		}
		if tab.Status == TabStatusClosed {
			return fmt.Errorf("selected tab %q is closed", s.SelectedTab)
		}
		if tab.Status == TabStatusStale {
			return fmt.Errorf("selected tab %q is stale", s.SelectedTab)
		}
	}
	for name, tab := range s.Tabs {
		if err := ValidateTabName(name); err != nil {
			return err
		}
		if tab.Name != name {
			return fmt.Errorf("tab %q must match map key", name)
		}
		if err := tab.validate(); err != nil {
			return fmt.Errorf("tab %q: %w", name, err)
		}
	}
	return nil
}

func (t *TabRecord) validate() error {
	if err := ValidateTabName(t.Name); err != nil {
		return err
	}
	switch t.Status {
	case TabStatusLive, TabStatusStale, TabStatusClosed:
		return nil
	default:
		return fmt.Errorf("unsupported tab status %q", t.Status)
	}
}

// CreateSession adds a new named session and selects it if it is the first one.
func (s *State) CreateSession(session *SessionRecord) error {
	if session == nil {
		return errors.New("session is required")
	}
	if err := session.validate(); err != nil {
		return err
	}
	if _, exists := s.Sessions[session.Name]; exists {
		return fmt.Errorf("session %q already exists", session.Name)
	}
	s.Sessions[session.Name] = session
	if s.SelectedSession == "" {
		s.SelectedSession = session.Name
	}
	return nil
}

// DeleteSession removes a session and clears the selected session when needed.
func (s *State) DeleteSession(name string) error {
	if _, ok := s.Sessions[name]; !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, name)
	}
	delete(s.Sessions, name)
	if s.SelectedSession == name {
		s.SelectedSession = ""
	}
	return nil
}

// SelectSession persists the default session used when --session is omitted.
func (s *State) SelectSession(name string) error {
	if _, ok := s.Sessions[name]; !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, name)
	}
	s.SelectedSession = name
	return nil
}

// ResolveSession returns the explicit session, the selected session, or the only
// available session when there is exactly one.
func (s *State) ResolveSession(name string) (*SessionRecord, error) {
	if name != "" {
		session, ok := s.Sessions[name]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrSessionNotFound, name)
		}
		return session, nil
	}
	if s.SelectedSession != "" {
		return s.Sessions[s.SelectedSession], nil
	}
	switch len(s.Sessions) {
	case 0:
		return nil, ErrNoSessions
	case 1:
		for _, session := range s.Sessions {
			return session, nil
		}
	}
	return nil, fmt.Errorf("%w: use --session or 'tap browser session select <name>'", ErrAmbiguousSession)
}

// UpsertTab creates or updates a tracked tab. The first live tab becomes selected.
func (s *State) UpsertTab(sessionName string, tab *TabRecord) error {
	if tab == nil {
		return errors.New("tab is required")
	}
	session, ok := s.Sessions[sessionName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, sessionName)
	}
	if err := tab.validate(); err != nil {
		return err
	}
	if session.Tabs == nil {
		session.Tabs = make(map[string]*TabRecord)
	}
	if existing, ok := session.Tabs[tab.Name]; ok {
		tab.CreatedAt = existing.CreatedAt
	}
	session.Tabs[tab.Name] = tab
	session.UpdatedAt = time.Now().UTC()
	if session.SelectedTab == "" && tab.Status == TabStatusLive {
		session.SelectedTab = tab.Name
	}
	return nil
}

// DeleteTab removes a tracked tab and advances selection to the next live tab.
func (s *State) DeleteTab(sessionName string, tabName string) error {
	session, ok := s.Sessions[sessionName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, sessionName)
	}
	if _, ok := session.Tabs[tabName]; !ok {
		return fmt.Errorf("%w: %s", ErrTabNotFound, tabName)
	}
	delete(session.Tabs, tabName)
	if session.SelectedTab == tabName {
		session.SelectedTab = session.nextLiveTabName()
	}
	session.UpdatedAt = time.Now().UTC()
	return nil
}

// SelectTab persists the default live tab used when --tab is omitted.
func (s *State) SelectTab(sessionName string, tabName string) error {
	session, ok := s.Sessions[sessionName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, sessionName)
	}
	tab, ok := session.Tabs[tabName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrTabNotFound, tabName)
	}
	if tab.Status == TabStatusClosed {
		return fmt.Errorf("%w: %s", ErrClosedTabSelected, tabName)
	}
	if tab.Status == TabStatusStale {
		return fmt.Errorf("%w: %s", ErrStaleTabSelected, tabName)
	}
	session.SelectedTab = tabName
	session.UpdatedAt = time.Now().UTC()
	return nil
}

// ResolveTab returns an explicit tab, the selected tab, or the only live tracked
// tab when exactly one live tab exists.
func (s *SessionRecord) ResolveTab(name string) (*TabRecord, error) {
	if s == nil {
		return nil, ErrNoTabs
	}
	if s.Tabs == nil {
		s.Tabs = make(map[string]*TabRecord)
	}
	if name != "" {
		tab, ok := s.Tabs[name]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrTabNotFound, name)
		}
		return tab, nil
	}
	if s.SelectedTab != "" {
		tab, ok := s.Tabs[s.SelectedTab]
		if ok {
			return tab, nil
		}
	}
	liveTabs := s.liveTabs()
	switch len(liveTabs) {
	case 0:
		return nil, ErrNoTabs
	case 1:
		return liveTabs[0], nil
	}
	return nil, fmt.Errorf("%w: use --tab or 'tap browser tab select <name>'", ErrAmbiguousTab)
}

// ReconcileSession updates tab liveness after reloading or reconnecting a browser session.
func (s *State) ReconcileSession(sessionName string, liveTargetIDs []string, now time.Time) error {
	session, ok := s.Sessions[sessionName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, sessionName)
	}

	liveTargets := make(map[string]struct{}, len(liveTargetIDs))
	for _, targetID := range liveTargetIDs {
		targetID = strings.TrimSpace(targetID)
		if targetID == "" {
			continue
		}
		liveTargets[targetID] = struct{}{}
	}

	for _, tab := range session.Tabs {
		if tab.Status == TabStatusClosed {
			continue
		}
		if _, ok := liveTargets[tab.TargetID]; ok && tab.TargetID != "" {
			tab.Status = TabStatusLive
			tab.LastSeenAt = now.UTC()
			tab.UpdatedAt = now.UTC()
			continue
		}
		tab.Status = TabStatusStale
		tab.TargetID = ""
		tab.UpdatedAt = now.UTC()
	}

	if session.SelectedTab != "" {
		selected, ok := session.Tabs[session.SelectedTab]
		if !ok || selected.Status != TabStatusLive {
			session.SelectedTab = ""
		}
	}
	session.LastReconciledAt = now.UTC()
	session.UpdatedAt = now.UTC()
	return nil
}

func (s *SessionRecord) nextLiveTabName() string {
	liveTabs := s.liveTabs()
	if len(liveTabs) == 0 {
		return ""
	}
	return liveTabs[0].Name
}

func (s *SessionRecord) liveTabs() []*TabRecord {
	var tabs []*TabRecord
	for _, tab := range s.Tabs {
		if tab.Status != TabStatusLive {
			continue
		}
		tabs = append(tabs, tab)
	}
	sort.Slice(tabs, func(i int, j int) bool {
		if tabs[i].CreatedAt.Equal(tabs[j].CreatedAt) {
			return tabs[i].Name < tabs[j].Name
		}
		return tabs[i].CreatedAt.Before(tabs[j].CreatedAt)
	})
	return tabs
}
