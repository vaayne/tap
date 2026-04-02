//go:build !windows

package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type fileLock struct {
	file *os.File
}

// lockFile acquires an exclusive file lock. It retries with LOCK_NB to avoid
// blocking indefinitely if another process holds the lock and never releases it.
func lockFile(path string) (*fileLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	// Retry with non-blocking flock to avoid hanging forever.
	deadline := time.After(30 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &fileLock{file: file}, nil
		}
		if err != syscall.EWOULDBLOCK {
			_ = file.Close()
			return nil, fmt.Errorf("lock file: %w", err)
		}
		select {
		case <-deadline:
			_ = file.Close()
			return nil, fmt.Errorf("lock file %s: timed out after 30s", filepath.Base(path))
		case <-ticker.C:
		}
	}
}

func (l *fileLock) Unlock() error {
	if l == nil || l.file == nil {
		return nil
	}
	err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	closeErr := l.file.Close()
	if err != nil {
		return err
	}
	return closeErr
}
