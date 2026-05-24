package browser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	AgentBrowserVersion = "0.27.0"
	EnvAgentBrowser     = "TAP_AGENT_BROWSER"
)

type AgentBrowserInstall struct {
	binDir  string
	version string
}

type AgentBrowserMeta struct {
	InstalledAt time.Time `json:"installed_at"`
	Source      string    `json:"source"`
	Version     string    `json:"version"`
}

func NewAgentBrowserInstall(binDir string) *AgentBrowserInstall {
	if binDir == "" {
		home, _ := os.UserHomeDir()
		binDir = filepath.Join(home, ".cache", "tap", "agent-browser")
	}
	return &AgentBrowserInstall{binDir: binDir, version: AgentBrowserVersion}
}

func ResolveAgentBrowserPath() (string, error) {
	if path := os.Getenv(EnvAgentBrowser); path != "" {
		return path, nil
	}
	install := NewAgentBrowserInstall("")
	if err := install.EnsureInstalled(context.Background()); err == nil {
		return install.binPath(), nil
	}
	if path, err := exec.LookPath("agent-browser"); err == nil {
		return path, nil
	}
	return "", errors.New("agent-browser embedded binary is unavailable for this platform; set TAP_AGENT_BROWSER")
}

func (a *AgentBrowserInstall) Installed() bool {
	fi, err := os.Stat(a.binPath())
	return err == nil && fi.Size() > 0
}

func (a *AgentBrowserInstall) Update(ctx context.Context) error {
	return a.extract(ctx)
}

func (a *AgentBrowserInstall) EnsureInstalled(ctx context.Context) error {
	if a.Installed() {
		return nil
	}
	return a.extract(ctx)
}

func (a *AgentBrowserInstall) ReadMeta() (*AgentBrowserMeta, error) {
	data, err := os.ReadFile(a.metaPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var meta AgentBrowserMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (a *AgentBrowserInstall) binPath() string {
	name := "agent-browser"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(a.binDir, name)
}

func (a *AgentBrowserInstall) metaPath() string {
	return a.binPath() + ".meta.json"
}

func (a *AgentBrowserInstall) extract(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if len(embeddedAgentBrowser) == 0 {
		return fmt.Errorf("embedded agent-browser is unavailable for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if err := os.MkdirAll(a.binDir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	bin := a.binPath()
	tmp := bin + ".tmp"
	if err := os.WriteFile(tmp, embeddedAgentBrowser, 0o755); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("write embedded binary: %w", err)
	}
	if err := os.Rename(tmp, bin); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename binary: %w", err)
	}
	return a.writeMeta("embedded")
}

func (a *AgentBrowserInstall) writeMeta(source string) error {
	meta := AgentBrowserMeta{InstalledAt: time.Now().UTC(), Source: source, Version: a.version}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.metaPath(), data, 0o644)
}

func agentBrowserPlatform(goos, goarch string) (string, error) {
	switch goos {
	case "darwin":
		if goarch == "arm64" {
			return "darwin-arm64", nil
		}
		if goarch == "amd64" {
			return "darwin-x64", nil
		}
	case "linux":
		if goarch == "arm64" {
			return "linux-arm64", nil
		}
		if goarch == "amd64" {
			return "linux-x64", nil
		}
	case "windows":
		if goarch == "amd64" {
			return "win32-x64", nil
		}
	}
	return "", fmt.Errorf("unsupported platform for agent-browser: %s/%s", goos, goarch)
}
