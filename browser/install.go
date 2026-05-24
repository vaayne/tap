package browser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	AgentBrowserVersion = "0.27.0"
	EnvAgentBrowser     = "TAP_AGENT_BROWSER"

	agentBrowserReleaseURL = "https://github.com/vercel-labs/agent-browser/releases/download"
)

type AgentBrowserInstall struct {
	binDir  string
	version string
}

type AgentBrowserMeta struct {
	DownloadedAt time.Time `json:"downloaded_at"`
	URL          string    `json:"url"`
	Version      string    `json:"version"`
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
	if path, err := exec.LookPath("agent-browser"); err == nil {
		return path, nil
	}
	install := NewAgentBrowserInstall("")
	if install.Installed() {
		return install.binPath(), nil
	}
	return "", errors.New("agent-browser not found: run tap doctor --install or set TAP_AGENT_BROWSER")
}

func (a *AgentBrowserInstall) Installed() bool {
	fi, err := os.Stat(a.binPath())
	return err == nil && fi.Size() > 0
}

func (a *AgentBrowserInstall) Update(ctx context.Context) error {
	return a.download(ctx)
}

func (a *AgentBrowserInstall) EnsureInstalled(ctx context.Context) error {
	if a.Installed() {
		return nil
	}
	return a.download(ctx)
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

func (a *AgentBrowserInstall) download(ctx context.Context) error {
	url, err := agentBrowserDownloadURL(runtime.GOOS, runtime.GOARCH, a.version)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(a.binDir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	bin := a.binPath()
	tmp := bin + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(tmp)
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write binary: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}
	if err := os.Rename(tmp, bin); err != nil {
		return fmt.Errorf("rename binary: %w", err)
	}
	return a.writeMeta(url)
}

func (a *AgentBrowserInstall) writeMeta(url string) error {
	meta := AgentBrowserMeta{DownloadedAt: time.Now().UTC(), URL: url, Version: a.version}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.metaPath(), data, 0o644)
}

func agentBrowserDownloadURL(goos, goarch, version string) (string, error) {
	platform, err := agentBrowserPlatform(goos, goarch)
	if err != nil {
		return "", err
	}
	asset := "agent-browser-" + platform
	if goos == "windows" {
		asset += ".exe"
	}
	return fmt.Sprintf("%s/v%s/%s", agentBrowserReleaseURL, version, asset), nil
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
