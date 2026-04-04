package browser

import (
	"os/exec"
	"strings"
)

// ChromeInfo holds information about an installed Chrome browser.
type ChromeInfo struct {
	Path    string
	Version string
}

// DetectChrome finds a Chrome/Chromium installation and returns its path and version.
// It reuses the same discovery logic as the browser launcher (findChrome).
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
	// "Google Chrome 125.0.6422.141" → "125.0.6422.141"
	output = strings.TrimSpace(output)
	parts := strings.Fields(output)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
