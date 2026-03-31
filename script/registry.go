package script

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Registry indexes scripts by their meta name.
type Registry struct {
	scripts map[string]*Script
	dir     string
}

// NewRegistry scans a directory for .js script files and indexes them by meta.name.
func NewRegistry(dir string) (*Registry, error) {
	r := &Registry{
		scripts: make(map[string]*Script),
		dir:     dir,
	}
	if err := r.scan(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) scan() error {
	return filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
