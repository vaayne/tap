//go:build windows

package browser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/sys/windows"
)

type fileLock struct {
	file *os.File
}

// lockFile acquires an exclusive file lock on Windows using LockFileEx.
// It retries with LOCKFILE_FAIL_IMMEDIATELY to avoid blocking indefinitely.
func lockFile(path string) (*fileLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	// Overlay for LockFileEx — lock the entire file.
	ol := new(windows.Overlapped)

	deadline := time.After(30 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		err := windows.LockFileEx(
			windows.Handle(file.Fd()),
			windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
			0, // reserved
			1, // lock 1 byte
			0, // high-order bytes
			ol,
		)
		if err == nil {
			// Prevent GC from finalizing file before LockFileEx completes.
			runtime.KeepAlive(file)
			return &fileLock{file: file}, nil
		}

		// Only retry on ERROR_LOCK_VIOLATION (lock held by another process).
		// Any other error is unexpected and should be returned immediately.
		if !errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
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
	ol := new(windows.Overlapped)
	err := windows.UnlockFileEx(windows.Handle(l.file.Fd()), 0, 1, 0, ol)
	// Prevent GC from finalizing file before UnlockFileEx completes.
	runtime.KeepAlive(l.file)
	closeErr := l.file.Close()
	if err != nil {
		return err
	}
	return closeErr
}
