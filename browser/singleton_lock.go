package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PrepareProfileDir removes stale Chrome singleton files before launch.
// A live local lock is left intact to avoid corrupting an active profile.
func PrepareProfileDir(profileDir string) error {
	lockPath := filepath.Join(profileDir, "SingletonLock")
	target, err := os.Readlink(lockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return clearSingletonFiles(profileDir)
	}

	host, pid, ok := parseSingletonLockTarget(target)
	localHost, hostErr := os.Hostname()
	if ok && hostErr == nil && host == localHost && isProcessAlive(pid) {
		return fmt.Errorf("profile %s is in use by another Chrome (pid %d); close it, or use --lp / tap attach chrome", profileDir, pid)
	}

	return clearSingletonFiles(profileDir)
}

func parseSingletonLockTarget(target string) (string, int, bool) {
	idx := strings.LastIndex(target, "-")
	if idx <= 0 || idx == len(target)-1 {
		return "", 0, false
	}
	pid, err := strconv.Atoi(target[idx+1:])
	if err != nil || pid <= 0 {
		return "", 0, false
	}
	return target[:idx], pid, true
}

func clearSingletonFiles(profileDir string) error {
	for _, name := range []string{"SingletonLock", "SingletonSocket", "SingletonCookie"} {
		if err := os.Remove(filepath.Join(profileDir, name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale chrome %s: %w", name, err)
		}
	}
	return nil
}
