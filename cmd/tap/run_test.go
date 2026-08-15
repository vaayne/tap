package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type fakeWorkflowBrowser struct {
	commands [][]string
	evals    []string
	err      error
}

func (b *fakeWorkflowBrowser) Run(_ context.Context, args ...string) (json.RawMessage, error) {
	b.commands = append(b.commands, append([]string(nil), args...))
	if b.err != nil {
		return nil, b.err
	}
	switch args[0] {
	case "open":
		return json.RawMessage(`{"url":"https://example.com/"}`), nil
	case "snapshot":
		return json.RawMessage(`{"snapshot":"tree"}`), nil
	default:
		return json.RawMessage(`{}`), nil
	}
}

func (b *fakeWorkflowBrowser) Eval(_ context.Context, script string) (any, error) {
	b.evals = append(b.evals, script)
	if b.err != nil {
		return nil, b.err
	}
	return "Example Domain", nil
}

func TestExecuteWorkflowDrivesBrowserAndPrintsOutput(t *testing.T) {
	browser := &fakeWorkflowBrowser{}
	var stdout, stderr bytes.Buffer
	source := `
const opened = await browser.open("https://example.com");
const page = await browser.snapshot("-i");
await browser.cmd("click", "@e1");
await browser.cmd("fill", "#query", "tap");
await browser.cmd("wait", 500);
const title = await browser.eval("document.title");
console.log(opened.url, page.snapshot, title);
`

	if err := executeWorkflow(context.Background(), browser, source, "workflow.js", &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	wantCommands := [][]string{
		{"open", "https://example.com"},
		{"snapshot", "-i"},
		{"click", "@e1"},
		{"fill", "#query", "tap"},
		{"wait", "500"},
	}
	if !reflect.DeepEqual(browser.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", browser.commands, wantCommands)
	}
	if !reflect.DeepEqual(browser.evals, []string{"document.title"}) {
		t.Fatalf("evals = %#v", browser.evals)
	}
	if stdout.String() != "https://example.com/ tree Example Domain\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestExecuteWorkflowCanCatchBrowserError(t *testing.T) {
	browser := &fakeWorkflowBrowser{err: errors.New("click failed")}
	var stdout bytes.Buffer
	source := `
try {
  await browser.cmd("click", "@e1");
} catch (error) {
  console.log(error.message);
}
`

	if err := executeWorkflow(context.Background(), browser, source, "workflow.js", &stdout, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "click failed\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestExecuteWorkflowReportsRejectedPromise(t *testing.T) {
	err := executeWorkflow(
		context.Background(),
		&fakeWorkflowBrowser{},
		`throw new Error("workflow failed");`,
		"workflow.js",
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err == nil || !strings.Contains(err.Error(), "workflow failed") {
		t.Fatalf("error = %v", err)
	}
}

func TestReadWorkflowFromStdin(t *testing.T) {
	source, filename, err := readWorkflow("", strings.NewReader("console.log('ok')"))
	if err != nil {
		t.Fatal(err)
	}
	if source != "console.log('ok')" || filename != "<stdin>" {
		t.Fatalf("source = %q, filename = %q", source, filename)
	}
}

func TestReadWorkflowFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workflow.js")
	if err := os.WriteFile(path, []byte(`await browser.cmd("wait", 10)`), 0o644); err != nil {
		t.Fatal(err)
	}
	source, filename, err := readWorkflow(path, strings.NewReader("ignored"))
	if err != nil {
		t.Fatal(err)
	}
	if source != `await browser.cmd("wait", 10)` || filename != path {
		t.Fatalf("source = %q, filename = %q", source, filename)
	}
}
