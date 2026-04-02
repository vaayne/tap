package browser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultStateRootUsesEnvOverride(t *testing.T) {
	root := t.TempDir()
	t.Setenv(EnvStateRoot, root)

	got, err := DefaultStateRoot()
	if err != nil {
		t.Fatalf("DefaultStateRoot failed: %v", err)
	}
	if got != root {
		t.Fatalf("DefaultStateRoot = %q, want %q", got, root)
	}
}

func TestStoreUpdatePersistsState(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	now := time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC)

	err = store.Update(func(state *State) error {
		session, err := NewLocalSession("alpha", filepath.Join(store.Root(), "profiles", "alpha"), true, now)
		if err != nil {
			return err
		}
		if err := state.CreateSession(session); err != nil {
			return err
		}
		tab, err := NewTab("first", "target-1", "https://example.com", now)
		if err != nil {
			return err
		}
		return state.UpsertTab("alpha", tab)
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	reloaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	session, err := reloaded.ResolveSession("alpha")
	if err != nil {
		t.Fatalf("ResolveSession failed: %v", err)
	}
	if session.SelectedTab != "first" {
		t.Fatalf("SelectedTab = %q, want first", session.SelectedTab)
	}
	if _, err := os.Stat(filepath.Join(store.Root(), stateFileName)); err != nil {
		t.Fatalf("state file missing: %v", err)
	}
}

func TestStoreUpdateSessionAppliesSessionScopedMutation(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	now := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	if err := store.Update(func(state *State) error {
		session, err := NewRemoteSession("remote", "wss://example.com/devtools/browser/1", now)
		if err != nil {
			return err
		}
		return state.CreateSession(session)
	}); err != nil {
		t.Fatalf("seed Update failed: %v", err)
	}

	if err := store.UpdateSession("remote", func(state *State, session *SessionRecord) error {
		tab, err := NewTab("docs", "target-9", "https://example.com/docs", now.Add(time.Minute))
		if err != nil {
			return err
		}
		return state.UpsertTab(session.Name, tab)
	}); err != nil {
		t.Fatalf("UpdateSession failed: %v", err)
	}

	reloaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	session, err := reloaded.ResolveSession("remote")
	if err != nil {
		t.Fatalf("ResolveSession failed: %v", err)
	}
	if _, ok := session.Tabs["docs"]; !ok {
		t.Fatalf("remote session tabs = %v, want docs", session.Tabs)
	}
}
