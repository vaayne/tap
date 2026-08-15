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

func TestRunPassesArgumentsWithoutShellAndReturnsData(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	bin := filepath.Join(dir, "agent-browser")
	script := `#!/bin/sh
printf '%s\n' "$@" > "$ARGS_FILE"
printf '%s' '{"success":true,"data":{"snapshot":"tree"},"error":null}'
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARGS_FILE", argsFile)

	data, err := New(bin).Run(context.Background(), "snapshot", "-i")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"snapshot":"tree"}` {
		t.Fatalf("data = %s", data)
	}
	assertFile(t, argsFile, "snapshot\n-i\n--json\n")
}

func TestRunRequiresCommand(t *testing.T) {
	_, err := New("unused").Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "command required") {
		t.Fatalf("error = %v, want command required", err)
	}
}

func TestRunReturnsStructuredErrorFromNonzeroCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	bin := filepath.Join(t.TempDir(), "agent-browser")
	script := `#!/bin/sh
printf '%s' '{"success":false,"error":"Element not found: @missing"}'
exit 1
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := New(bin).Run(context.Background(), "click", "@missing")
	if err == nil || err.Error() != "Element not found: @missing" {
		t.Fatalf("error = %v", err)
	}
}

func TestPathPrefersBundledSibling(t *testing.T) {
	dir := t.TempDir()
	name := "agent-browser"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(dir, name)
	if err := os.WriteFile(bin, []byte("fixture"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvBinary, "")
	client := New("")
	client.siblingDir = dir
	path, err := client.Path()
	if err != nil {
		t.Fatal(err)
	}
	if path != bin {
		t.Fatalf("path = %q, want sibling %q", path, bin)
	}
}

func TestPathExplicitOverrideBeatsBundledSibling(t *testing.T) {
	dir := t.TempDir()
	name := "agent-browser"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	sibling := filepath.Join(dir, name)
	if err := os.WriteFile(sibling, []byte("fixture"), 0o755); err != nil {
		t.Fatal(err)
	}
	override := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(override, []byte("fixture"), 0o755); err != nil {
		t.Fatal(err)
	}
	client := New(override)
	client.siblingDir = dir
	path, err := client.Path()
	if err != nil {
		t.Fatal(err)
	}
	if path != override {
		t.Fatalf("path = %q, want explicit override %q", path, override)
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

func TestOpenAndEvalReturnsStructuredBatchFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	bin := filepath.Join(t.TempDir(), "agent-browser")
	script := `#!/bin/sh
cat >/dev/null
printf '%s' '[{"success":true,"result":{}},{"success":false,"result":null,"error":"Evaluation error: TypeError: Failed to fetch"}]'
exit 1
`
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := New(bin).OpenAndEval(context.Background(), "https://example.com", "1", nil)
	if err == nil {
		t.Fatal("expected batch failure")
	}
	if !strings.Contains(err.Error(), "TypeError: Failed to fetch") {
		t.Fatalf("error = %q, want structured evaluation error", err)
	}
	if strings.Contains(err.Error(), "exit status 1") {
		t.Fatalf("error leaked opaque process status: %q", err)
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
