package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelpShowsCompletionCommand(t *testing.T) {
	out, err := runTap(t, []string{"tap", "--help"}, "")
	if err != nil {
		t.Fatalf("run tap --help: %v", err)
	}
	if !strings.Contains(out, "completion") {
		t.Fatalf("expected help to mention completion command, got:\n%s", out)
	}
}

func TestCompletionCommandOutputsBashScript(t *testing.T) {
	out, err := runTapExternal(t, []string{"completion", "bash"}, "", nil)
	if err != nil {
		t.Fatalf("run tap completion bash: %v", err)
	}
	if !strings.Contains(out, "--generate-shell-completion") {
		t.Fatalf("expected bash completion script output, got:\n%s", out)
	}
}

func TestSiteShellCompletionUsesLocalRegistry(t *testing.T) {
	dir := t.TempDir()
	writeTestScript(t, filepath.Join(dir, "hackernews", "top.js"), testScript("hackernews/top", "Top stories"))

	out, err := runTapExternal(t, []string{"--sites-dir", dir, "site", "ha", "--generate-shell-completion"}, "/bin/bash", nil)
	if err != nil {
		t.Fatalf("run site shell completion: %v", err)
	}
	if !strings.Contains(out, "hackernews/top") {
		t.Fatalf("expected script completion from local registry, got:\n%s", out)
	}
}

func TestSiteShellCompletionFallsBackToOverrideDirWhenCacheMissing(t *testing.T) {
	configHome := t.TempDir()
	overrideDir := filepath.Join(configHome, "tap", "sites")
	writeTestScript(t, filepath.Join(overrideDir, "lobsters", "frontpage.js"), testScript("lobsters/frontpage", "Front page stories"))

	out, err := runTapExternal(
		t,
		[]string{"--sites-dir", filepath.Join(t.TempDir(), "missing-cache"), "site", "lo", "--generate-shell-completion"},
		"/bin/bash",
		[]string{"XDG_CONFIG_HOME=" + configHome},
	)
	if err != nil {
		t.Fatalf("run site shell completion with missing cache: %v", err)
	}
	if !strings.Contains(out, "lobsters/frontpage") {
		t.Fatalf("expected override script completion when cache is missing, got:\n%s", out)
	}
}

func runTap(t *testing.T, args []string, shell string) (string, error) {
	t.Helper()

	app := newApp()
	var out bytes.Buffer
	app.Writer = &out
	app.ErrWriter = &out

	oldArgs := os.Args
	oldShell := os.Getenv("SHELL")
	defer func() {
		os.Args = oldArgs
		if oldShell == "" {
			_ = os.Unsetenv("SHELL")
		} else {
			_ = os.Setenv("SHELL", oldShell)
		}
	}()

	os.Args = append([]string(nil), args...)
	if shell == "" {
		_ = os.Unsetenv("SHELL")
	} else {
		_ = os.Setenv("SHELL", shell)
	}

	err := app.Run(context.Background(), args)
	return out.String(), err
}

func runTapExternal(t *testing.T, args []string, shell string, extraEnv []string) (string, error) {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), extraEnv...)
	if shell != "" {
		cmd.Env = append(cmd.Env, "SHELL="+shell)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func testScript(name, description string) string {
	return fmt.Sprintf(`/* @meta
{
  "name": %q,
  "description": %q,
  "domain": "example.com"
}
*/

async function(args) {
  return {};
}`,
		name,
		description,
	)
}

func writeTestScript(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
