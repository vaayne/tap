package browser

import (
	"testing"
	"time"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestManagerCreateSessionLocalValidation(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	// Empty name should fail validation.
	err := mgr.CreateSession(t.Context(), "", ModeLocal, SessionOptions{})
	if err == nil {
		t.Fatal("CreateSession with empty name should return error")
	}

	// Invalid name (spaces) should fail.
	err = mgr.CreateSession(t.Context(), "bad name", ModeLocal, SessionOptions{})
	if err == nil {
		t.Fatal("CreateSession with invalid name should return error")
	}
}

func TestManagerCreateSessionRemoteInvalidEndpoint(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	// Unreachable endpoint should fail.
	err := mgr.CreateSession(t.Context(), "remote1", ModeRemote, SessionOptions{
		WSURL: "ws://127.0.0.1:1/devtools/browser/unreachable",
	})
	if err == nil {
		t.Fatal("CreateSession with unreachable remote endpoint should return error")
	}
}

func TestManagerListSessionsEmpty(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	list, err := mgr.ListSessions(t.Context())
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(list.Sessions) != 0 {
		t.Fatalf("ListSessions = %d sessions, want 0", len(list.Sessions))
	}
}

func TestManagerGetSessionUsesPersistedDefaultContext(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	err := store.Update(func(state *State) error {
		def, err := NewLocalSession(DefaultSessionName, store.Root()+"/profiles/default", true, now)
		if err != nil {
			return err
		}
		other, err := NewRemoteSession("other", "wss://remote:9222/devtools/browser/1", now)
		if err != nil {
			return err
		}
		if err := state.CreateSession(def); err != nil {
			return err
		}
		if err := state.CreateSession(other); err != nil {
			return err
		}
		return state.SetDefaultContext("other", DefaultContextAttached, now)
	})
	if err != nil {
		t.Fatalf("seed sessions: %v", err)
	}

	session, err := mgr.GetSession(t.Context(), "")
	if err != nil {
		t.Fatalf("GetSession error: %v", err)
	}
	if session.Name != "other" {
		t.Fatalf("GetSession(\"\") = %q, want other", session.Name)
	}
}

func TestManagerGetSessionResolution(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	// Insert a "default" session and another session.
	err := store.Update(func(state *State) error {
		s1, err := NewLocalSession(DefaultSessionName, store.Root()+"/profiles/default", true, now)
		if err != nil {
			return err
		}
		if err := state.CreateSession(s1); err != nil {
			return err
		}

		s2, err := NewRemoteSession("other", "wss://remote:9222/devtools/browser/1", now)
		if err != nil {
			return err
		}
		return state.CreateSession(s2)
	})
	if err != nil {
		t.Fatalf("seed sessions: %v", err)
	}

	// GetSession with empty name should resolve to "default".
	session, err := mgr.GetSession(t.Context(), "")
	if err != nil {
		t.Fatalf("GetSession error: %v", err)
	}
	if session.Name != DefaultSessionName {
		t.Fatalf("GetSession(\"\") = %q, want %q", session.Name, DefaultSessionName)
	}

	// GetSession with explicit name should resolve correctly.
	session, err = mgr.GetSession(t.Context(), "other")
	if err != nil {
		t.Fatalf("GetSession(other) error: %v", err)
	}
	if session.Name != "other" {
		t.Fatalf("GetSession(other) = %q, want other", session.Name)
	}
}

func TestResolveDebugURL(t *testing.T) {
	t.Run("local with process debug URL", func(t *testing.T) {
		session := &SessionRecord{
			Name: "local1",
			Mode: ModeLocal,
			Process: &ProcessRecord{
				DebugURL: "ws://127.0.0.1:9222/devtools/browser/abc",
			},
		}
		got, err := resolveDebugURL(session)
		if err != nil {
			t.Fatalf("resolveDebugURL error: %v", err)
		}
		if got != "ws://127.0.0.1:9222/devtools/browser/abc" {
			t.Fatalf("resolveDebugURL = %q, want ws://127.0.0.1:9222/devtools/browser/abc", got)
		}
	})

	t.Run("remote with WSURL", func(t *testing.T) {
		session := &SessionRecord{
			Name:   "remote1",
			Mode:   ModeRemote,
			Remote: &RemoteConfig{WSURL: "wss://remote:9222/devtools/browser/xyz"},
		}
		got, err := resolveDebugURL(session)
		if err != nil {
			t.Fatalf("resolveDebugURL error: %v", err)
		}
		if got != "wss://remote:9222/devtools/browser/xyz" {
			t.Fatalf("resolveDebugURL = %q, want wss://remote:9222/devtools/browser/xyz", got)
		}
	})

	t.Run("no debug endpoint", func(t *testing.T) {
		session := &SessionRecord{
			Name: "empty",
			Mode: ModeLocal,
		}
		_, err := resolveDebugURL(session)
		if err == nil {
			t.Fatal("resolveDebugURL should return error when no debug endpoint")
		}
	})
}

func TestResolveTargetMarksUnreachableRemoteSessionStale(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	err := store.Update(func(state *State) error {
		session, err := NewRemoteSession("attached", "ws://127.0.0.1:1/devtools/browser/unreachable", now)
		if err != nil {
			return err
		}
		tab, err := NewTab("main", "target-1", "https://example.com", now)
		if err != nil {
			return err
		}
		session.Tabs[tab.Name] = tab
		session.SelectedTab = tab.Name
		session.Process = &ProcessRecord{DebugURL: session.Remote.WSURL}
		if err := state.CreateSession(session); err != nil {
			return err
		}
		return state.SetDefaultContext("attached", DefaultContextAttached, now)
	})
	if err != nil {
		t.Fatalf("seed remote session: %v", err)
	}

	if _, err := mgr.ResolveTarget(t.Context(), "", ""); err == nil {
		t.Fatal("ResolveTarget should fail for unreachable remote session")
	}

	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	session := state.Sessions["attached"]
	if session.Tabs["main"].Status != TabStatusStale {
		t.Fatalf("tab status = %q, want stale", session.Tabs["main"].Status)
	}
	if state.DefaultContext == nil || !state.DefaultContext.Stale {
		t.Fatal("default context should be marked stale")
	}
}

func TestRequireLiveTab(t *testing.T) {
	t.Run("live tab with target ID", func(t *testing.T) {
		tab := &TabRecord{Name: "t1", Status: TabStatusLive, TargetID: "target-1"}
		if err := requireLiveTab(tab); err != nil {
			t.Fatalf("requireLiveTab error: %v", err)
		}
	})

	t.Run("stale tab", func(t *testing.T) {
		tab := &TabRecord{Name: "t2", Status: TabStatusStale}
		if err := requireLiveTab(tab); err == nil {
			t.Fatal("requireLiveTab should return error for stale tab")
		}
	})

	t.Run("closed tab", func(t *testing.T) {
		tab := &TabRecord{Name: "t3", Status: TabStatusClosed}
		if err := requireLiveTab(tab); err == nil {
			t.Fatal("requireLiveTab should return error for closed tab")
		}
	})

	t.Run("live tab with empty target ID", func(t *testing.T) {
		tab := &TabRecord{Name: "t4", Status: TabStatusLive, TargetID: ""}
		if err := requireLiveTab(tab); err == nil {
			t.Fatal("requireLiveTab should return error for live tab with empty target ID")
		}
	})
}
