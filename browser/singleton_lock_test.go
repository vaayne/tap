package browser

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestPrepareProfileDirMissingLockNoop(t *testing.T) {
	dir := t.TempDir()

	if err := PrepareProfileDir(dir); err != nil {
		t.Fatalf("PrepareProfileDir() error = %v", err)
	}
}

func TestPrepareProfileDirLiveLocalPIDHeld(t *testing.T) {
	dir := t.TempDir()
	host, err := os.Hostname()
	if err != nil {
		t.Fatalf("Hostname() error = %v", err)
	}
	pid := os.Getpid()
	if err := os.Symlink(host+"-"+strconv.Itoa(pid), filepath.Join(dir, "SingletonLock")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	err = PrepareProfileDir(dir)
	if err == nil {
		t.Fatal("PrepareProfileDir() error = nil, want live-lock error")
	}
	if want := "in use by another Chrome (pid " + strconv.Itoa(pid) + ")"; !strings.Contains(err.Error(), want) {
		t.Fatalf("PrepareProfileDir() error = %q, want %q", err, want)
	}
	if _, statErr := os.Lstat(filepath.Join(dir, "SingletonLock")); statErr != nil {
		t.Fatalf("SingletonLock removed for live pid: %v", statErr)
	}
}

func TestPrepareProfileDirDeadPIDStale(t *testing.T) {
	dir := t.TempDir()
	host, err := os.Hostname()
	if err != nil {
		t.Fatalf("Hostname() error = %v", err)
	}
	for _, name := range []string{"SingletonLock", "SingletonSocket", "SingletonCookie"} {
		path := filepath.Join(dir, name)
		if name == "SingletonLock" {
			if err := os.Symlink(host+"-99999999", path); err != nil {
				t.Fatalf("Symlink() error = %v", err)
			}
			continue
		}
		if err := os.WriteFile(path, []byte("stale"), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", name, err)
		}
	}

	if err := PrepareProfileDir(dir); err != nil {
		t.Fatalf("PrepareProfileDir() error = %v", err)
	}
	assertSingletonFilesRemoved(t, dir)
}

func TestPrepareProfileDirGarbageTargetStale(t *testing.T) {
	dir := t.TempDir()
	if err := os.Symlink("not-a-chrome-lock", filepath.Join(dir, "SingletonLock")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}

	if err := PrepareProfileDir(dir); err != nil {
		t.Fatalf("PrepareProfileDir() error = %v", err)
	}
	assertSingletonFilesRemoved(t, dir)
}

func assertSingletonFilesRemoved(t *testing.T, dir string) {
	t.Helper()
	for _, name := range []string{"SingletonLock", "SingletonSocket", "SingletonCookie"} {
		_, err := os.Lstat(filepath.Join(dir, name))
		if !os.IsNotExist(err) {
			t.Fatalf("%s still exists or stat failed with non-ENOENT: %v", name, err)
		}
	}
}
