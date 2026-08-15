package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBrowserCommandPassesThroughProcessIO(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}

	bin := filepath.Join(t.TempDir(), "agent-browser")
	script := `#!/bin/sh
printf 'args:'
for arg in "$@"; do printf '\n<%s>' "$arg"; done
printf '\nstdin:'
cat
printf '\nengine=%s\n' "$AGENT_BROWSER_ENGINE"
printf 'browser warning\n' >&2
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TAP_AGENT_BROWSER", bin)
	t.Setenv("AGENT_BROWSER_ENGINE", "lightpanda")

	app := newApp()
	app.Reader = bytes.NewBufferString("document.title")
	var stdout, stderr bytes.Buffer
	app.Writer = &stdout
	app.ErrWriter = &stderr

	err := app.Run(context.Background(), []string{
		"tap", "browser", "--engine", "lightpanda", "eval", "value with spaces; echo unsafe",
	})
	if err != nil {
		t.Fatal(err)
	}
	wantStdout := "args:\n<--engine>\n<lightpanda>\n<eval>\n<value with spaces; echo unsafe>\nstdin:document.title\nengine=lightpanda\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout = %q, want %q", stdout.String(), wantStdout)
	}
	if stderr.String() != "browser warning\n" {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestBrowserCommandPreservesExitCode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}

	bin := filepath.Join(t.TempDir(), "agent-browser")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nprintf 'browser failed\\n' >&2\nexit 7\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TAP_AGENT_BROWSER", bin)

	app := newApp()
	app.Reader = bytes.NewReader(nil)
	app.Writer = &bytes.Buffer{}
	var stderr bytes.Buffer
	app.ErrWriter = &stderr

	err := app.Run(context.Background(), []string{"tap", "browser", "click", "@missing"})
	var exitErr *browserExitError
	if !errors.As(err, &exitErr) || exitErr.code != 7 {
		t.Fatalf("error = %v, want browser exit code 7", err)
	}
	if stderr.String() != "browser failed\n" {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
