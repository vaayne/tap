package tap

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func testSitesDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("TAP_SITES_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot determine home dir")
		}
		dir = filepath.Join(home, ".config", "tap", "sites")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("sites dir %s does not exist, run 'tap site sync' first", dir)
	}
	return dir
}

func TestNew(t *testing.T) {
	dir := testSitesDir(t)
	client, err := New(WithSitesDir(dir))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	scripts := client.ListScripts()
	if len(scripts) == 0 {
		t.Fatal("expected scripts, got none")
	}

	s, ok := client.GetScript("v2ex/hot")
	if !ok {
		t.Fatal("v2ex/hot not found")
	}
	if s.Meta.Domain != "www.v2ex.com" {
		t.Errorf("domain = %q, want %q", s.Meta.Domain, "www.v2ex.com")
	}
}

func TestNew_NoSitesDir(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	if client.ListScripts() != nil {
		t.Error("expected nil scripts without sites dir")
	}
}

func TestRunScript_NotFound(t *testing.T) {
	dir := testSitesDir(t)
	client, err := New(WithSitesDir(dir))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	_, err = client.RunScript(context.Background(), "nonexistent/script", nil)
	if err == nil {
		t.Error("expected error for nonexistent script")
	}
}

func TestRunScript_MissingRequiredArg(t *testing.T) {
	dir := testSitesDir(t)
	client, err := New(WithSitesDir(dir))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	// twitter/search requires "query"
	_, err = client.RunScript(context.Background(), "twitter/search", nil)
	if err == nil {
		t.Error("expected error for missing required arg")
	}
}

func TestFetch(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	result, err := client.Fetch(context.Background(), "https://example.com", nil)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if result.Title == "" && result.Content == "" && result.Markdown == "" {
		t.Error("expected some content")
	}
}
