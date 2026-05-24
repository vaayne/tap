package browser

import "testing"

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

func TestAgentBrowserDownloadURL(t *testing.T) {
	got, err := agentBrowserDownloadURL("darwin", "arm64", "0.27.0")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/vercel-labs/agent-browser/releases/download/v0.27.0/agent-browser-darwin-arm64"
	if got != want {
		t.Fatalf("download URL = %q, want %q", got, want)
	}

	got, err = agentBrowserDownloadURL("windows", "amd64", "0.27.0")
	if err != nil {
		t.Fatal(err)
	}
	want = "https://github.com/vercel-labs/agent-browser/releases/download/v0.27.0/agent-browser-win32-x64.exe"
	if got != want {
		t.Fatalf("download URL = %q, want %q", got, want)
	}
}
