package browser

import (
	"slices"
	"testing"
)

func TestAgentBrowserCommandArgs(t *testing.T) {
	ab := &AgentBrowser{Path: "agent-browser", SessionName: "dev", ProfileDir: "/tmp/profile", Headed: true}
	args := ab.commandArgs("open", "https://example.com")
	for _, want := range []string{"open", "https://example.com", "--json", "--session-name", "dev", "--profile", "/tmp/profile", "--headed"} {
		if !slices.Contains(args, want) {
			t.Fatalf("args %v missing %q", args, want)
		}
	}
}

func TestAgentBrowserCommandArgsAttachedSkipsSession(t *testing.T) {
	ab := &AgentBrowser{Path: "agent-browser", SessionName: "dev", Attached: true}
	args := ab.commandArgs("get", "url")
	if slices.Contains(args, "--session-name") {
		t.Fatalf("attached args included session name: %v", args)
	}
}

func TestAgentBrowserCommandArgsPreservesExplicitJSONAndSession(t *testing.T) {
	ab := &AgentBrowser{Path: "agent-browser", SessionName: "dev"}
	args := ab.commandArgs("session", "--json", "--session-name", "other")
	if got := count(args, "--json"); got != 1 {
		t.Fatalf("--json count = %d, want 1 in %v", got, args)
	}
	if got := count(args, "--session-name"); got != 1 {
		t.Fatalf("--session-name count = %d, want 1 in %v", got, args)
	}
}

func count(values []string, target string) int {
	var n int
	for _, value := range values {
		if value == target {
			n++
		}
	}
	return n
}
