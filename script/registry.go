package script

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Source describes one script source for the registry.
type Source struct {
	FS   fs.FS        // virtual filesystem (nil → use Path)
	Path string       // filesystem path (ignored if FS is set)
	Type ScriptSource // how to tag scripts from this source
}

// Registry indexes scripts by their meta name.
type Registry struct {
	scripts map[string]*Script
	sources []Source
}

// NewRegistry creates a Registry from one or more sources.
// Sources are scanned in order; later sources overwrite earlier ones
// with the same script name.
func NewRegistry(sources ...Source) (*Registry, error) {
	r := &Registry{
		scripts: make(map[string]*Script),
		sources: sources,
	}
	if err := r.scan(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) scan() error {
	for _, src := range r.sources {
		if src.FS != nil {
			if err := r.scanFS(src.FS, ".", src.Type); err != nil {
				return err
			}
		} else if src.Path != "" {
			if err := r.scanDir(src.Path, src.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Registry) scanDir(dir string, source ScriptSource) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // skip silently
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		s.Source = source
		r.scripts[s.Meta.Name] = s
		return nil
	})
}

func (r *Registry) scanFS(fsys fs.FS, root string, source ScriptSource) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == "." {
				return nil // missing root in FS — skip silently
			}
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".js" {
			return nil
		}

		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		s, err := Parse(string(content))
		if err != nil {
			// Skip scripts that fail to parse
			return nil
		}

		s.Path = path
		s.Source = source
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

// ListOverrides returns only scripts loaded from the local override directory.
func (r *Registry) ListOverrides() []*Script {
	scripts := make([]*Script, 0)
	for _, s := range r.scripts {
		if s.Source == ScriptSourceOverride {
			scripts = append(scripts, s)
		}
	}
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Meta.Name < scripts[j].Meta.Name
	})
	return scripts
}
