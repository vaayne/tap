package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	stateFileName = "state.json"
	lockDirName   = "locks"
)

// Store manages durable browser metadata with atomic writes and file locks.
type Store struct {
	root string
}

// DefaultStateRoot returns the durable directory used for browser metadata.
func DefaultStateRoot() (string, error) {
	if root := os.Getenv(EnvStateRoot); root != "" {
		return root, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(configDir, "tap", "browser"), nil
}

// NewStore initializes a metadata store rooted at the provided directory.
func NewStore(root string) (*Store, error) {
	var err error
	if root == "" {
		root, err = DefaultStateRoot()
		if err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(filepath.Join(root, lockDirName), 0o755); err != nil {
		return nil, fmt.Errorf("create browser state dir: %w", err)
	}

	return &Store{root: root}, nil
}

// Root returns the store's durable state directory.
func (s *Store) Root() string {
	return s.root
}

// Load reads the current browser metadata state from disk.
func (s *Store) Load() (*State, error) {
	if s == nil {
		return nil, errors.New("browser store is required")
	}

	data, err := os.ReadFile(s.statePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewState(), nil
		}
		return nil, fmt.Errorf("read browser state: %w", err)
	}

	state := NewState()
	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("decode browser state: %w", err)
	}
	if err := state.Validate(); err != nil {
		return nil, fmt.Errorf("validate browser state: %w", err)
	}
	return state, nil
}

// Save writes the full metadata state atomically.
func (s *Store) Save(state *State) error {
	if s == nil {
		return errors.New("browser store is required")
	}
	if state == nil {
		return errors.New("browser state is required")
	}
	if err := state.Validate(); err != nil {
		return err
	}
	unlock, err := s.lock("state")
	if err != nil {
		return err
	}
	defer func() { _ = unlock() }()
	return s.writeState(state)
}

// Update takes an exclusive store lock, mutates the current state, and saves it.
func (s *Store) Update(fn func(*State) error) error {
	unlock, err := s.lock("state")
	if err != nil {
		return err
	}
	defer func() { _ = unlock() }()

	state, err := s.Load()
	if err != nil {
		return err
	}

	if err := fn(state); err != nil {
		return err
	}

	return s.writeState(state)
}

// WithSessionLock serializes higher-level runtime work for one session.
func (s *Store) WithSessionLock(sessionName string, fn func() error) error {
	if err := ValidateSessionName(sessionName); err != nil {
		return err
	}
	unlock, err := s.lock("session-" + sessionName)
	if err != nil {
		return err
	}
	defer func() { _ = unlock() }()
	return fn()
}

// UpdateSession combines session-scoped locking with durable store mutation.
func (s *Store) UpdateSession(sessionName string, fn func(*State, *SessionRecord) error) error {
	return s.WithSessionLock(sessionName, func() error {
		return s.Update(func(state *State) error {
			session, err := state.ResolveSession(sessionName)
			if err != nil {
				return err
			}
			return fn(state, session)
		})
	})
}

func (s *Store) writeState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode browser state: %w", err)
	}

	tmp, err := os.CreateTemp(s.root, "state-*.json")
	if err != nil {
		return fmt.Errorf("create temp browser state: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp browser state: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temp browser state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp browser state: %w", err)
	}
	if err := os.Rename(tmpName, s.statePath()); err != nil {
		return fmt.Errorf("replace browser state: %w", err)
	}
	return nil
}

func (s *Store) lock(name string) (func() error, error) {
	lock, err := lockFile(filepath.Join(s.root, lockDirName, name+".lock"))
	if err != nil {
		return nil, err
	}
	return lock.Unlock, nil
}

func (s *Store) statePath() string {
	return filepath.Join(s.root, stateFileName)
}
