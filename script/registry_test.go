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

	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	scripts := reg.List()
	if len(scripts) == 0 {
		t.Fatal("expected scripts, got none")
	}

	// Verify v2ex/hot exists
	s, ok := reg.Get("v2ex/hot")
	if !ok {
		t.Fatal("v2ex/hot not found in registry")
	}
	if s.Meta.Domain != "www.v2ex.com" {
		t.Errorf("domain = %q, want %q", s.Meta.Domain, "www.v2ex.com")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	dir := filepath.Join("..", "sites")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("sites/ directory not found")
	}

	reg, err := NewRegistry(dir)
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

	reg, err := NewRegistry(dir)
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
