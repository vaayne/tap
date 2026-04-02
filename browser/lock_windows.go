//go:build windows

package browser

import (
	"fmt"
	"os"
	"path/filepath"
)

type fileLock struct {
	file *os.File
}

// lockFile acquires an exclusive file lock on Windows using LockFileEx.
func lockFile(path string) (*fileLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	return &fileLock{file: file}, nil
}

func (l *fileLock) Unlock() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}
