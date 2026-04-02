//go:build !windows

package browser

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// platformSysProcAttr returns SysProcAttr to run Chrome in its own process
// group so it survives tap exit on unix systems.
func platformSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

// isProcessAlive checks whether a process with the given PID exists.
func isProcessAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

// killProcessPlatform terminates a process using SIGTERM first, falling back
// to SIGKILL after a 5-second grace period.
func killProcessPlatform(pid int) error {
	// Send SIGTERM for a graceful shutdown.
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return fmt.Errorf("send SIGTERM to %d: %w", pid, err)
	}

	// Wait up to 5 seconds for the process to exit.
	deadline := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			// Process didn't exit in time — force kill.
			if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
				if errors.Is(err, syscall.ESRCH) {
					return nil
				}
				return fmt.Errorf("send SIGKILL to %d: %w", pid, err)
			}
			return nil
		case <-ticker.C:
			if !isProcessAlive(pid) {
				return nil // process exited
			}
		}
	}
}

// chromeFallbackPaths returns well-known absolute Chrome/Chromium paths for
// macOS and Linux.
func chromeFallbackPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	default:
		return []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
		}
	}
}

// chromeLookPathNames returns binary names to search via exec.LookPath.
func chromeLookPathNames() []string {
	names := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium-browser",
		"chromium",
	}
	// On macOS also try the .app bundle name via open.
	if runtime.GOOS == "darwin" {
		if p, err := exec.LookPath("Google Chrome"); err == nil {
			return []string{p}
		}
	}
	return names
}
