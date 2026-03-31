package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Registry struct {
	scripts map[string]*Script // keyed by meta.name
	dir     string
}

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

		script, err := ParseScript(string(content))
		if err != nil {
			// Skip scripts that fail to parse
			return nil
		}

		script.Path = path
		r.scripts[script.Meta.Name] = script
		return nil
	})
}

func (r *Registry) Get(name string) (*Script, bool) {
	s, ok := r.scripts[name]
	return s, ok
}

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
