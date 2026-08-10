package script

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	// Use the real sites/ dir
	dir := filepath.Join("..", "sites")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("sites/ directory not found")
	}

	reg, err := NewRegistry(Source{Path: dir, Type: ScriptSourceCache})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	scripts := reg.List()
	if len(scripts) == 0 {
		t.Fatal("expected scripts, got none")
	}

	// Verify exa/search exists
	s, ok := reg.Get("exa/search")
	if !ok {
		t.Fatal("exa/search not found in registry")
	}
	if s.Meta.Domain != "mcp.exa.ai" {
		t.Errorf("domain = %q, want %q", s.Meta.Domain, "mcp.exa.ai")
	}
	if len(s.Meta.Headers) != 1 {
		t.Errorf("headers count = %d, want 1", len(s.Meta.Headers))
	} else if s.Meta.Headers["X-API-Key"] != "${EXA_API_KEY}" {
		t.Errorf("headers[X-API-Key] = %q, want %q", s.Meta.Headers["X-API-Key"], "${EXA_API_KEY}")
	}
	if s.Source != ScriptSourceCache {
		t.Errorf("source = %d, want ScriptSourceCache (%d)", s.Source, ScriptSourceCache)
	}
}

func TestRegistryNameComesFromPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wrong-meta", "actual.js")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `/* @meta {"name":"ignored/name","description":"test","domain":"example.com"} */
async function () { return 1 }`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	reg, err := NewRegistry(Source{Path: dir})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Get("wrong-meta/actual"); !ok {
		t.Fatal("path-derived name not registered")
	}
	if _, ok := reg.Get("ignored/name"); ok {
		t.Fatal("metadata name must be ignored")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	dir := filepath.Join("..", "sites")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("sites/ directory not found")
	}

	reg, err := NewRegistry(Source{Path: dir, Type: ScriptSourceCache})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	_, ok := reg.Get("nonexistent/script")
	if ok {
		t.Error("expected not found")
	}
}

func TestRegistry_ListSorted(t *testing.T) {
	dir := filepath.Join("..", "sites")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("sites/ directory not found")
	}

	reg, err := NewRegistry(Source{Path: dir, Type: ScriptSourceCache})
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	scripts := reg.List()
	for i := 1; i < len(scripts); i++ {
		if scripts[i].Meta.Name < scripts[i-1].Meta.Name {
			t.Errorf("not sorted: %q came after %q", scripts[i].Meta.Name, scripts[i-1].Meta.Name)
		}
	}
}
