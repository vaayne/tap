package tap

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewWithoutSites(t *testing.T) {
	client, err := New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if client.ListScripts() != nil {
		t.Fatal("expected no registry")
	}
}

func TestRunScriptValidatesRequiredArgsBeforeBrowser(t *testing.T) {
	client := testClient(t, `/* @meta
{"description":"search","domain":"example.com","args":{"query":{"required":true,"description":"query"}}}
*/
async function(args) { return args; }`)
	_, err := client.RunScript(context.Background(), "test/run", nil)
	if err == nil || !strings.Contains(err.Error(), "missing required arg") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunScriptUsesAgentBrowserAndExpandsHeaders(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	t.Setenv("API_KEY", "test-key")
	client := testClient(t, `/* @meta
{"description":"search","domain":"example.com","args":{"query":{"required":true}},"headers":{"X-Key":"${API_KEY}"}}
*/
async function(args) { return {query: args.query}; }`)
	result, err := client.RunScript(context.Background(), "test/run", map[string]string{"query": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if result.(map[string]any)["ok"] != true {
		t.Fatalf("unexpected result: %#v", result)
	}
	args, err := os.ReadFile(os.Getenv("ARGS_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(args), "set") {
		t.Fatalf("site headers leaked into browser session: %s", args)
	}
	stdin, err := os.ReadFile(os.Getenv("STDIN_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	var commands [][]string
	if err := json.Unmarshal(stdin, &commands); err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 || len(commands[0]) != 2 {
		t.Fatalf("unexpected orchestration: %#v", commands)
	}
	decoded, err := base64.StdEncoding.DecodeString(commands[1][2])
	if err != nil {
		t.Fatal(err)
	}
	program := string(decoded)
	for _, want := range []string{`"query":"hello"`, `"X-Key":"test-key"`, `"example.com"`, "url.origin === __tapHeaderOrigin", "globalThis.fetch.bind"} {
		if !strings.Contains(program, want) {
			t.Fatalf("program missing %q", want)
		}
	}
	if strings.Contains(program, "cross-origin fetch blocked") {
		t.Fatal("Tap must not block cross-origin fetches")
	}
}

func testClient(t *testing.T, content string) *Client {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "test", "run.js")
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scriptPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	stdinFile := filepath.Join(dir, "stdin")
	bin := filepath.Join(dir, "agent-browser")
	fixture := `#!/bin/sh
printf '%s\n' "$@" > "$ARGS_FILE"
if [ "$1" = "batch" ]; then
  cat > "$STDIN_FILE"
  printf '%s' '[{"success":true,"result":{}},{"success":true,"result":{"result":{"ok":true}}}]'
else
  printf '%s' '{"success":true,"data":{},"error":null}'
fi
`
	if err := os.WriteFile(bin, []byte(fixture), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARGS_FILE", argsFile)
	t.Setenv("STDIN_FILE", stdinFile)
	client, err := New(context.Background(), WithSitesDir(dir), WithAgentBrowserBinary(bin))
	if err != nil {
		t.Fatal(err)
	}
	return client
}
