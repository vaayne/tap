package tap

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vaayne/tap/engine"
)

func testSitesDir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("TAP_SITES_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot determine home dir")
		}
		dir = filepath.Join(home, ".cache", "tap", "sites")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("sites dir %s does not exist, run 'tap site sync' first", dir)
	}
	return dir
}

func TestNew(t *testing.T) {
	dir := testSitesDir(t)
	client, err := New(context.Background(), WithSitesDir(dir))
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
	client, err := New(context.Background())
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
	client, err := New(context.Background(), WithSitesDir(dir))
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
	client, err := New(context.Background(), WithSitesDir(dir))
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

func TestEnginesByRuntime(t *testing.T) {
	quickjs := engine.NewQuickJS(nil)
	browser := engine.NewBrowser(nil, nil)

	tests := []struct {
		name         string
		engines      []engine.Engine
		forceBrowser bool
		runtime      string
		want         []string
	}{
		{"normal auto", []engine.Engine{quickjs, browser}, false, "auto", []string{"QuickJS", "Browser"}},
		{"normal empty", []engine.Engine{quickjs, browser}, false, "", []string{"QuickJS", "Browser"}},
		{"normal http", []engine.Engine{quickjs, browser}, false, "http", []string{"QuickJS"}},
		{"normal browser", []engine.Engine{quickjs, browser}, false, "browser", []string{"Browser"}},
		{"normal lightpanda", []engine.Engine{quickjs, browser}, false, "lightpanda", []string{"Browser"}},
		{"normal unknown", []engine.Engine{quickjs, browser}, false, "unknown", []string{"QuickJS", "Browser"}},
		{"forceBrowser http", []engine.Engine{quickjs, browser}, true, "http", []string{"Browser"}},
		{"forceBrowser browser", []engine.Engine{quickjs, browser}, true, "browser", []string{"Browser"}},
		{"forceBrowser auto", []engine.Engine{quickjs, browser}, true, "auto", []string{"Browser"}},
		{"empty http", []engine.Engine{}, false, "http", []string{}},
		{"empty browser", []engine.Engine{}, false, "browser", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				engines: tt.engines,
				opts:    options{forceBrowser: tt.forceBrowser},
			}
			got := c.enginesByRuntime(tt.runtime)
			gotNames := make([]string, len(got))
			for i, e := range got {
				gotNames[i] = e.Name()
			}
			if !slices.Equal(gotNames, tt.want) {
				t.Errorf("enginesByRuntime(%q) = %v, want %v", tt.runtime, gotNames, tt.want)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	client, err := New(context.Background())
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
