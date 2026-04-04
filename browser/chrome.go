package browser

import (
	"os/exec"
	"runtime"
	"strings"
)

// chromePaths lists common Chrome/Chromium binary names per OS.
var chromePaths = map[string][]string{
	"darwin": {
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
	},
	"linux": {
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
	},
	"windows": {
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
	},
}

// ChromeInfo holds information about an installed Chrome browser.
type ChromeInfo struct {
	Path    string
	Version string
}

// DetectChrome finds a Chrome/Chromium installation and returns its path and version.
// Returns nil if no Chrome is found.
func DetectChrome() *ChromeInfo {
	candidates := chromePaths[runtime.GOOS]

	for _, c := range candidates {
		path, err := exec.LookPath(c)
		if err != nil {
			// On macOS the paths are absolute, LookPath won't find them
			// unless they're in PATH — try direct stat for absolute paths.
			if strings.HasPrefix(c, "/") || strings.HasPrefix(c, `C:\`) {
				if _, statErr := exec.LookPath(c); statErr != nil {
					// Try running it directly to see if it exists.
					if out, runErr := exec.Command(c, "--version").Output(); runErr == nil {
						return &ChromeInfo{
							Path:    c,
							Version: parseVersion(string(out)),
						}
					}
				}
				continue
			}
			continue
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

	return nil
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
