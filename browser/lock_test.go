//go:build !windows

package browser

import (
	"testing"
	"time"
)

func TestSessionLockSerializesConcurrentAccess(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	firstLocked := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondDone := make(chan struct{})

	go func() {
		if err := store.WithSessionLock("alpha", func() error {
			close(firstLocked)
			<-releaseFirst
			return nil
		}); err != nil {
			t.Errorf("first WithSessionLock failed: %v", err)
		}
	}()

	<-firstLocked

	go func() {
		if err := store.WithSessionLock("alpha", func() error {
			close(secondDone)
			return nil
		}); err != nil {
			t.Errorf("second WithSessionLock failed: %v", err)
		}
	}()

	select {
	case <-secondDone:
		t.Fatal("second session lock acquired before first lock released")
	case <-time.After(75 * time.Millisecond):
	}

	close(releaseFirst)

	select {
	case <-secondDone:
	case <-time.After(2 * time.Second):
		t.Fatal("second session lock did not acquire after release")
	}
}
