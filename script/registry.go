package script

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Registry indexes scripts by their meta name.
type Registry struct {
	scripts          map[string]*Script
	dir              string
	localOverrideDir string // checked first; empty = disabled
}

// NewRegistry scans a directory for .js script files and indexes them by meta name.
func NewRegistry(dir string) (*Registry, error) {
	return NewRegistryWithOverride(dir, "")
}

// NewRegistryWithOverride is like NewRegistry but also checks localOverrideDir
// before the main cache dir. Scripts found there shadow the cached versions.
func NewRegistryWithOverride(dir, localOverrideDir string) (*Registry, error) {
	r := &Registry{
		scripts:          make(map[string]*Script),
		dir:              dir,
		localOverrideDir: localOverrideDir,
	}
	if err := r.scan(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) scan() error {
	// Load main cache dir first.
	if err := r.scanDir(r.dir, false); err != nil {
		return err
	}
	// Load local override dir second — overwrites any same-named cache entry.
	if r.localOverrideDir != "" {
		if err := r.scanDir(r.localOverrideDir, true); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) scanDir(dir string, isOverride bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Ignore missing override dir — it may not exist yet.
			if os.IsNotExist(err) && isOverride {
				return filepath.SkipDir
			}
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".js" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		s, err := Parse(string(content))
		if err != nil {
			// Skip scripts that fail to parse
			return nil
		}

		s.Path = path
		s.LocalOverride = isOverride
		r.scripts[s.Meta.Name] = s
		return nil
	})
}

// Get returns a script by its meta name.
func (r *Registry) Get(name string) (*Script, bool) {
	s, ok := r.scripts[name]
	return s, ok
}

// List returns all scripts sorted by name.
func (r *Registry) List() []*Script {
	scripts := make([]*Script, 0, len(r.scripts))
	for _, s := range r.scripts {
		scripts = append(scripts, s)
	}
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Meta.Name < scripts[j].Meta.Name
	})
	return scripts
}

// ListLocalOnly returns only scripts loaded from the local override directory.
func (r *Registry) ListLocalOnly() []*Script {
	scripts := make([]*Script, 0)
	for _, s := range r.scripts {
		if s.LocalOverride {
			scripts = append(scripts, s)
		}
	}
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Meta.Name < scripts[j].Meta.Name
	})
	return scripts
}
