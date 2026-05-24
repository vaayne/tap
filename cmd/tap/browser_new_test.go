package main

import "testing"

func TestExtractPassthroughSession(t *testing.T) {
	args, session := extractPassthroughSession([]string{"--filter", "api", "--session", "dev", "--limit", "10"})
	if session != "dev" {
		t.Fatalf("session = %q, want dev", session)
	}
	want := []string{"--filter", "api", "--limit", "10"}
	if len(args) != len(want) {
		t.Fatalf("args = %v, want %v", args, want)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("args = %v, want %v", args, want)
		}
	}
}

func TestExtractPassthroughSessionEquals(t *testing.T) {
	args, session := extractPassthroughSession([]string{"--session=dev", "--abort"})
	if session != "dev" {
		t.Fatalf("session = %q, want dev", session)
	}
	if len(args) != 1 || args[0] != "--abort" {
		t.Fatalf("args = %v, want [--abort]", args)
	}
}
