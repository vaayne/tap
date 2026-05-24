//go:build windows

package browser

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func chromeFallbackPaths() []string {
	var paths []string

	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	localAppData := os.Getenv("LOCALAPPDATA")

	for _, root := range []string{programFiles, programFilesX86, localAppData} {
		if root == "" {
			continue
		}
		paths = append(paths, filepath.Join(root, "Google", "Chrome", "Application", "chrome.exe"))
	}

	if p := chromeFromRegistry(); p != "" {
		paths = append(paths, p)
	}

	return paths
}

func chromeLookPathNames() []string {
	return []string{
		"chrome",
		"chromium",
	}
}

func chromeFromRegistry() string {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return ""
	}
	defer func() { _ = k.Close() }()

	p, _, err := k.GetStringValue("")
	if err != nil {
		return ""
	}
	return p
}
