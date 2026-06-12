package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDocsDrift(t *testing.T) {
	got := renderDocs(newApp())

	// Locate docs/cli.md relative to this test file (cmd/tap/docs_test.go).
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file path")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	docsPath := filepath.Join(repoRoot, "docs", "cli.md")

	committed, err := os.ReadFile(docsPath)
	if err != nil {
		t.Fatalf("read %s: %v\n\nRun `mise run docs` to generate it.", docsPath, err)
	}

	if got != string(committed) {
		t.Errorf("docs/cli.md is out of date.\n\nRun `mise run docs` to regenerate it.\n\nDiff (got vs committed):\n%s",
			diffStrings(got, string(committed)))
	}
}

// diffStrings returns a simple line-level diff for test output, capped at 20 lines.
func diffStrings(got, want string) string {
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")

	var b strings.Builder
	n := max(len(gotLines), len(wantLines))
	differences := 0
	for i := 0; i < n && differences < 20; i++ {
		var g, w string
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if g != w {
			differences++
			fmt.Fprintf(&b, "line %d:\n  got:  %s\n  want: %s\n", i+1, g, w)
		}
	}
	if differences == 0 && len(gotLines) != len(wantLines) {
		fmt.Fprintf(&b, "line count differs: got %d, want %d\n", len(gotLines), len(wantLines))
	}
	return b.String()
}
