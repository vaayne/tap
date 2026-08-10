package agentbrowser

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestEvalUsesStdinAndInheritsSession(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	stdinFile := filepath.Join(dir, "stdin")
	envFile := filepath.Join(dir, "env")
	bin := filepath.Join(dir, "agent-browser")
	script := `#!/bin/sh
printf '%s\n' "$@" > "$ARGS_FILE"
cat > "$STDIN_FILE"
printf '%s' "$AGENT_BROWSER_SESSION" > "$ENV_FILE"
printf '%s' '{"success":true,"data":{"result":{"ok":true}},"error":null}'
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARGS_FILE", argsFile)
	t.Setenv("STDIN_FILE", stdinFile)
	t.Setenv("ENV_FILE", envFile)
	t.Setenv("AGENT_BROWSER_SESSION", "task-session")

	value, err := New(bin).Eval(context.Background(), "line 1\nline '2'")
	if err != nil {
		t.Fatal(err)
	}
	if value.(map[string]any)["ok"] != true {
		t.Fatalf("unexpected value: %#v", value)
	}
	assertFile(t, argsFile, "eval\n--stdin\n--json\n")
	assertFile(t, stdinFile, "line 1\nline '2'")
	assertFile(t, envFile, "task-session")
}

func TestOpenDoesNotAddSessionFlags(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	bin := filepath.Join(dir, "agent-browser")
	script := `#!/bin/sh
printf '%s\n' "$@" > "$ARGS_FILE"
printf '%s' '{"success":true,"data":{},"error":null}'
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARGS_FILE", argsFile)
	if err := New(bin).Open(context.Background(), "https://example.com"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "session") {
		t.Fatalf("Tap must not manage sessions: %s", data)
	}
	if string(data) != "open\nhttps://example.com\n--json\n" {
		t.Fatalf("unexpected args: %q", data)
	}
}

func TestOpenAndEvalStagesHeadersThroughBatchStdin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	stdinFile := filepath.Join(dir, "stdin")
	bin := filepath.Join(dir, "agent-browser")
	script := `#!/bin/sh
printf '%s\n' "$@" >> "$ARGS_FILE"
if [ "$1" = "batch" ]; then
  cat > "$STDIN_FILE"
  printf '%s' '[{"success":true,"result":{}},{"success":true,"result":{"result":{"ok":true}}}]'
else
  printf '%s' '{"success":true,"data":{},"error":null}'
fi
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARGS_FILE", argsFile)
	t.Setenv("STDIN_FILE", stdinFile)
	value, err := New(bin).OpenAndEval(
		context.Background(),
		"https://example.com",
		"line 1\nline 2",
		map[string]string{"Authorization": "Bearer secret"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if value.(map[string]any)["ok"] != true {
		t.Fatalf("unexpected value: %#v", value)
	}
	assertFile(t, argsFile, "batch\n--bail\n--json\nset\nheaders\n{}\n--json\n")
	data, err := os.ReadFile(stdinFile)
	if err != nil {
		t.Fatal(err)
	}
	var commands [][]string
	if err := json.Unmarshal(data, &commands); err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 {
		t.Fatalf("commands = %#v", commands)
	}
	if commands[0][0] != "open" || commands[0][2] != "--headers" {
		t.Fatalf("headers must precede navigation: %#v", commands[0])
	}
	if strings.Contains(string(mustRead(t, argsFile)), "secret") {
		t.Fatal("secret leaked into process args")
	}
}

func TestOpenAndEvalClearsHeadersAfterCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	logFile := filepath.Join(dir, "log")
	bin := filepath.Join(dir, "agent-browser")
	script := `#!/bin/sh
printf '%s\n' "$1" >> "$LOG_FILE"
if [ "$1" = "batch" ]; then
  sleep 5
else
  printf '%s' '{"success":true,"data":{},"error":null}'
fi
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LOG_FILE", logFile)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := New(bin).OpenAndEval(ctx, "https://example.com", "1", map[string]string{"X-Key": "secret"})
	if err == nil {
		t.Fatal("expected batch cancellation")
	}
	data := string(mustRead(t, logFile))
	if !strings.Contains(data, "set\n") {
		t.Fatalf("cleanup did not run after cancellation: %q", data)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, data, want)
	}
}
