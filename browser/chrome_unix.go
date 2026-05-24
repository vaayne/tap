//go:build !windows

package browser

import (
	"os/exec"
	"runtime"
)

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

func chromeLookPathNames() []string {
	names := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium-browser",
		"chromium",
	}
	if runtime.GOOS == "darwin" {
		if p, err := exec.LookPath("Google Chrome"); err == nil {
			return []string{p}
		}
	}
	return names
}
