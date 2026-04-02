//go:build windows

package browser

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// platformSysProcAttr returns SysProcAttr to run Chrome in a new process group
// so it is not killed when tap exits.
func platformSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

// isProcessAlive checks whether a process with the given PID exists and is
// still running on Windows.
func isProcessAlive(pid int) bool {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h) //nolint:errcheck

	var exitCode uint32
	if err := windows.GetExitCodeProcess(h, &exitCode); err != nil {
		return false
	}
	// STILL_ACTIVE (259) means the process has not exited.
	return exitCode == 259
}

// killProcessPlatform terminates a process on Windows. It first tries a
// graceful shutdown via taskkill (sends WM_CLOSE), then falls back to
// forceful termination after a 5-second grace period.
func killProcessPlatform(pid int) error {
	pidStr := strconv.Itoa(pid)

	// Phase 1: graceful shutdown (WM_CLOSE via taskkill without /F).
	_ = exec.Command("taskkill", "/T", "/PID", pidStr).Run()

	// Wait up to 5 seconds for the process to exit.
	deadline := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			// Process didn't exit in time — force kill.
			return exec.Command("taskkill", "/F", "/T", "/PID", pidStr).Run()
		case <-ticker.C:
			if !isProcessAlive(pid) {
				return nil
			}
		}
	}
}

// chromeFallbackPaths returns well-known absolute Chrome paths on Windows,
// including per-user and system-wide install locations and registry lookup.
func chromeFallbackPaths() []string {
	var paths []string

	// Standard install locations.
	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	localAppData := os.Getenv("LOCALAPPDATA")

	for _, root := range []string{programFiles, programFilesX86, localAppData} {
		if root == "" {
			continue
		}
		paths = append(paths, filepath.Join(root, "Google", "Chrome", "Application", "chrome.exe"))
	}

	// Registry lookup: HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe
	if p := chromeFromRegistry(); p != "" {
		paths = append(paths, p)
	}

	return paths
}

// chromeLookPathNames returns binary names to search via exec.LookPath on Windows.
func chromeLookPathNames() []string {
	return []string{
		"chrome",
		"chrome.exe",
		"chromium",
		"chromium.exe",
	}
}

// chromeFromRegistry reads the default Chrome path from the Windows registry.
func chromeFromRegistry() string {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return ""
	}
	defer k.Close() //nolint:errcheck

	val, _, err := k.GetStringValue("")
	if err != nil {
		return ""
	}
	return val
}
