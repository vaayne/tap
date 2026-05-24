package browser

import (
	"context"
	"os"
	"testing"
)

func TestAgentBrowserPlatform(t *testing.T) {
	tests := []struct {
		goos string
		arch string
		want string
	}{
		{"darwin", "arm64", "darwin-arm64"},
		{"darwin", "amd64", "darwin-x64"},
		{"linux", "arm64", "linux-arm64"},
		{"linux", "amd64", "linux-x64"},
		{"windows", "amd64", "win32-x64"},
	}
	for _, tt := range tests {
		got, err := agentBrowserPlatform(tt.goos, tt.arch)
		if err != nil {
			t.Fatalf("agentBrowserPlatform(%q, %q): %v", tt.goos, tt.arch, err)
		}
		if got != tt.want {
			t.Fatalf("agentBrowserPlatform(%q, %q) = %q, want %q", tt.goos, tt.arch, got, tt.want)
		}
	}
}

func TestAgentBrowserExtract(t *testing.T) {
	dir := t.TempDir()
	install := NewAgentBrowserInstall(dir)
	if err := install.Update(context.Background()); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(install.binPath())
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() == 0 {
		t.Fatal("extracted binary is empty")
	}
	meta, err := install.ReadMeta()
	if err != nil {
		t.Fatal(err)
	}
	if meta == nil || meta.Source != "embedded" || meta.Version != AgentBrowserVersion {
		t.Fatalf("meta = %#v", meta)
	}
}
