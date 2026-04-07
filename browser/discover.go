package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// DiscoverUserChromeDebugURL resolves the browser WebSocket URL for an already
// running user Chrome/Chromium instance that has remote debugging enabled.
// It reads the DevToolsActivePort file from common user-data directories.
func DiscoverUserChromeDebugURL() (string, string, error) {
	for _, path := range devToolsActivePortCandidates() {
		url, err := ResolveDebugURLFromDevToolsFile(path)
		if err == nil {
			return url, path, nil
		}
		if !os.IsNotExist(err) {
			return "", path, err
		}
	}
	return "", "", fmt.Errorf("could not find DevToolsActivePort for a running user Chrome; enable remote debugging on your existing browser or pass --upstream")
}

// ResolveDebugURLFromDevToolsFile parses a DevToolsActivePort file and returns
// the browser WebSocket URL.
func ResolveDebugURLFromDevToolsFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return parseDevToolsActivePort(string(content))
}

func parseDevToolsActivePort(content string) (string, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid DevToolsActivePort content")
	}

	port, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || port <= 0 {
		return "", fmt.Errorf("invalid DevToolsActivePort port")
	}

	wsPath := strings.TrimSpace(lines[1])
	if !strings.HasPrefix(wsPath, "/") {
		return "", fmt.Errorf("invalid DevToolsActivePort websocket path")
	}

	return fmt.Sprintf("ws://127.0.0.1:%d%s", port, wsPath), nil
}

func devToolsActivePortCandidates() []string {
	home, _ := os.UserHomeDir()
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" && home != "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	switch runtime.GOOS {
	case "darwin":
		return existingCandidates(home,
			filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "DevToolsActivePort"),
			filepath.Join(home, "Library", "Application Support", "Google", "Chrome Beta", "DevToolsActivePort"),
			filepath.Join(home, "Library", "Application Support", "Chromium", "DevToolsActivePort"),
		)
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		return existingCandidates(localAppData,
			filepath.Join(localAppData, "Google", "Chrome", "User Data", "DevToolsActivePort"),
			filepath.Join(localAppData, "Google", "Chrome Beta", "User Data", "DevToolsActivePort"),
			filepath.Join(localAppData, "Chromium", "User Data", "DevToolsActivePort"),
		)
	default:
		return existingCandidates(xdgConfig,
			filepath.Join(xdgConfig, "google-chrome", "DevToolsActivePort"),
			filepath.Join(xdgConfig, "google-chrome-beta", "DevToolsActivePort"),
			filepath.Join(xdgConfig, "chromium", "DevToolsActivePort"),
		)
	}
}

func existingCandidates(root string, paths ...string) []string {
	if strings.TrimSpace(root) == "" {
		return nil
	}
	return paths
}
