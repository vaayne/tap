//go:build windows

package browser

import "context"

// scanElectronProcesses is not yet implemented on Windows.
// Returns an empty list without error.
func scanElectronProcesses(_ context.Context) ([]ElectronProcess, error) {
	return nil, nil
}

// extractDebugPort parses --remote-debugging-port=PORT from a command line.
// Shared definition for Windows (unix version lives in electron_unix.go).
func extractDebugPort(_ string) int {
	return 0
}
