package browser

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ChromeInfo holds information about an installed Chrome browser.
type ChromeInfo struct {
	Path    string
	Version string
}

// DetectChrome finds a Chrome/Chromium installation and returns its path and version.
// Returns nil if no Chrome is found.
func DetectChrome() *ChromeInfo {
	path, err := findChrome()
	if err != nil {
		return nil
	}

	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return &ChromeInfo{Path: path}
	}
	return &ChromeInfo{
		Path:    path,
		Version: parseVersion(string(out)),
	}
}

func parseVersion(output string) string {
	output = strings.TrimSpace(output)
	parts := strings.Fields(output)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func findChrome() (string, error) {
	for _, name := range chromeLookPathNames() {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	for _, p := range chromeFallbackPaths() {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", errors.New("chrome binary not found: install Google Chrome or Chromium")
}
