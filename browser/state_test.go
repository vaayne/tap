package browser

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestStateSessionResolution(t *testing.T) {
	state := NewState()
	now := time.Date(2026, 4, 1, 8, 0, 0, 0, time.UTC)

	alpha, err := NewLocalSession("alpha", filepath.Join(t.TempDir(), "alpha"), true, now)
	if err != nil {
		t.Fatalf("NewLocalSession failed: %v", err)
	}
	beta, err := NewRemoteSession("beta", "wss://example.com/devtools/browser/1", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("NewRemoteSession failed: %v", err)
	}

	if err := state.CreateSession(alpha); err != nil {
		t.Fatalf("CreateSession(alpha) failed: %v", err)
	}
	if err := state.CreateSession(beta); err != nil {
		t.Fatalf("CreateSession(beta) failed: %v", err)
	}

	if session, err := state.ResolveSession(""); err != nil {
		t.Fatalf("ResolveSession(selected) failed: %v", err)
	} else if session.Name != "alpha" {
		t.Fatalf("ResolveSession(selected) = %q, want alpha", session.Name)
	}

	if err := state.SelectSession("beta"); err != nil {
		t.Fatalf("SelectSession(beta) failed: %v", err)
	}
	if session, err := state.ResolveSession(""); err != nil {
		t.Fatalf("ResolveSession(selected after select) failed: %v", err)
	} else if session.Name != "beta" {
		t.Fatalf("ResolveSession(selected after select) = %q, want beta", session.Name)
	}

	state.SelectedSession = ""
	if _, err := state.ResolveSession(""); !errors.Is(err, ErrAmbiguousSession) {
		t.Fatalf("ResolveSession(ambiguous) error = %v, want ErrAmbiguousSession", err)
	}

	if session, err := state.ResolveSession("alpha"); err != nil {
		t.Fatalf("ResolveSession(alpha) failed: %v", err)
	} else if session.Name != "alpha" {
		t.Fatalf("ResolveSession(alpha) = %q, want alpha", session.Name)
	}
}

func TestSessionTabSelectionAndDeletion(t *testing.T) {
	state := NewState()
	now := time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC)

	session, err := NewLocalSession("alpha", filepath.Join(t.TempDir(), "alpha"), false, now)
	if err != nil {
		t.Fatalf("NewLocalSession failed: %v", err)
	}
	if err := state.CreateSession(session); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	first, err := NewTab("first", "target-1", "https://example.com/1", now)
	if err != nil {
		t.Fatalf("NewTab(first) failed: %v", err)
	}
	second, err := NewTab("second", "target-2", "https://example.com/2", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("NewTab(second) failed: %v", err)
	}

	if err := state.UpsertTab("alpha", first); err != nil {
		t.Fatalf("UpsertTab(first) failed: %v", err)
	}
	if err := state.UpsertTab("alpha", second); err != nil {
		t.Fatalf("UpsertTab(second) failed: %v", err)
	}

	if session.SelectedTab != "first" {
		t.Fatalf("SelectedTab after first insert = %q, want first", session.SelectedTab)
	}
	if err := state.SelectTab("alpha", "second"); err != nil {
		t.Fatalf("SelectTab(second) failed: %v", err)
	}
	if session.SelectedTab != "second" {
		t.Fatalf("SelectedTab after select = %q, want second", session.SelectedTab)
	}

	if err := state.DeleteTab("alpha", "second"); err != nil {
		t.Fatalf("DeleteTab(second) failed: %v", err)
	}
	if session.SelectedTab != "first" {
		t.Fatalf("SelectedTab after delete = %q, want first", session.SelectedTab)
	}

	resolved, err := session.ResolveTab("")
	if err != nil {
		t.Fatalf("ResolveTab(selected) failed: %v", err)
	}
	if resolved.Name != "first" {
		t.Fatalf("ResolveTab(selected) = %q, want first", resolved.Name)
	}
}

func TestResolveTabRequiresSelectionWhenMultipleLiveTabsExist(t *testing.T) {
	session, err := NewLocalSession("alpha", filepath.Join(t.TempDir(), "alpha"), true, time.Now().UTC())
	if err != nil {
		t.Fatalf("NewLocalSession failed: %v", err)
	}

	first, _ := NewTab("first", "target-1", "", time.Now().UTC())
	second, _ := NewTab("second", "target-2", "", time.Now().UTC().Add(time.Minute))
	session.Tabs[first.Name] = first
	session.Tabs[second.Name] = second
	session.SelectedTab = ""

	if _, err := session.ResolveTab(""); !errors.Is(err, ErrAmbiguousTab) {
		t.Fatalf("ResolveTab(ambiguous) error = %v, want ErrAmbiguousTab", err)
	}
}

func TestReconcileSessionMarksMissingTargetsStale(t *testing.T) {
	state := NewState()
	now := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	session, err := NewLocalSession("alpha", filepath.Join(t.TempDir(), "alpha"), true, now)
	if err != nil {
		t.Fatalf("NewLocalSession failed: %v", err)
	}
	tab, err := NewTab("first", "target-1", "https://example.com", now)
	if err != nil {
		t.Fatalf("NewTab failed: %v", err)
	}
	session.Tabs[tab.Name] = tab
	session.SelectedTab = tab.Name

	if err := state.CreateSession(session); err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if err := state.ReconcileSession("alpha", []string{"target-2"}, now.Add(5*time.Minute)); err != nil {
		t.Fatalf("ReconcileSession failed: %v", err)
	}

	if tab.Status != TabStatusStale {
		t.Fatalf("tab status = %q, want stale", tab.Status)
	}
	if tab.TargetID != "" {
		t.Fatalf("tab target id = %q, want empty", tab.TargetID)
	}
	if session.SelectedTab != "" {
		t.Fatalf("SelectedTab after reconcile = %q, want empty", session.SelectedTab)
	}
}

func TestValidateNameRules(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "empty", value: "", want: false},
		{name: "space", value: "bad name", want: false},
		{name: "slash", value: "bad/name", want: false},
		{name: "valid", value: "good-name_1", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSessionName(tc.value)
			if tc.want && err != nil {
				t.Fatalf("ValidateSessionName(%q) failed: %v", tc.value, err)
			}
			if !tc.want && err == nil {
				t.Fatalf("ValidateSessionName(%q) succeeded, want error", tc.value)
			}
		})
	}
}
