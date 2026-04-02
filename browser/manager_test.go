package browser

import (
	"context"
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
	err := mgr.CreateSession(context.Background(), "", ModeLocal, SessionOptions{})
	if err == nil {
		t.Fatal("CreateSession with empty name should return error")
	}

	// Invalid name (spaces) should fail.
	err = mgr.CreateSession(context.Background(), "bad name", ModeLocal, SessionOptions{})
	if err == nil {
		t.Fatal("CreateSession with invalid name should return error")
	}
}

func TestManagerCreateSessionRemoteInvalidEndpoint(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	// Unreachable endpoint should fail.
	err := mgr.CreateSession(context.Background(), "remote1", ModeRemote, SessionOptions{
		WSURL: "ws://127.0.0.1:1/devtools/browser/unreachable",
	})
	if err == nil {
		t.Fatal("CreateSession with unreachable remote endpoint should return error")
	}
}

func TestManagerListSessionsEmpty(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	list, err := mgr.ListSessions(context.Background())
	if err != nil {
		t.Fatalf("ListSessions error: %v", err)
	}
	if len(list.Sessions) != 0 {
		t.Fatalf("ListSessions = %d sessions, want 0", len(list.Sessions))
	}
}

func TestManagerSelectSession(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	// Manually insert a session into the store.
	err := store.Update(func(state *State) error {
		session, err := NewLocalSession("alpha", store.Root()+"/profiles/alpha", true, now)
		if err != nil {
			return err
		}
		session.Process = &ProcessRecord{DebugURL: "ws://127.0.0.1:9222/devtools/browser/abc"}
		return state.CreateSession(session)
	})
	if err != nil {
		t.Fatalf("seed session: %v", err)
	}

	// SelectSession should succeed.
	if err := mgr.SelectSession(context.Background(), "alpha"); err != nil {
		t.Fatalf("SelectSession error: %v", err)
	}

	// Verify selected session is correct.
	session, err := mgr.GetSession(context.Background(), "")
	if err != nil {
		t.Fatalf("GetSession error: %v", err)
	}
	if session.Name != "alpha" {
		t.Fatalf("selected session = %q, want alpha", session.Name)
	}
}

func TestManagerSelectSessionNotFound(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)

	err := mgr.SelectSession(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("SelectSession with nonexistent name should return error")
	}
}

func TestManagerGetSessionResolution(t *testing.T) {
	store := testStore(t)
	mgr := NewManager(store)
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	// Insert two sessions and select one.
	err := store.Update(func(state *State) error {
		s1, err := NewLocalSession("first", store.Root()+"/profiles/first", true, now)
		if err != nil {
			return err
		}
		if err := state.CreateSession(s1); err != nil {
			return err
		}

		s2, err := NewRemoteSession("second", "wss://remote:9222/devtools/browser/1", now)
		if err != nil {
			return err
		}
		if err := state.CreateSession(s2); err != nil {
			return err
		}

		// First session is auto-selected; switch to second.
		return state.SelectSession("second")
	})
	if err != nil {
		t.Fatalf("seed sessions: %v", err)
	}

	// GetSession with empty name should resolve to the selected session.
	session, err := mgr.GetSession(context.Background(), "")
	if err != nil {
		t.Fatalf("GetSession error: %v", err)
	}
	if session.Name != "second" {
		t.Fatalf("GetSession(\"\") = %q, want second", session.Name)
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
