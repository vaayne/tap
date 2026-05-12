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
	if s.Meta.Runtime != "http" {
		t.Errorf("runtime = %q, want %q", s.Meta.Runtime, "http")
	}
	if len(s.Meta.Env) != 1 {
		t.Errorf("env count = %d, want 1", len(s.Meta.Env))
	} else {
		def, ok := s.Meta.Env["EXA_API_KEY"]
		if !ok {
			t.Error("env missing EXA_API_KEY")
		} else {
			if def.Required {
				t.Error("env[EXA_API_KEY].required = true, want false")
			}
			if def.Description != "API key for Exa search (increases rate limit)" {
				t.Errorf("env[EXA_API_KEY].description = %q, want %q", def.Description, "API key for Exa search (increases rate limit)")
			}
		}
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
