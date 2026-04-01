// Package browser manages alternative browser backends for CDP-based automation.
package browser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	lightpandaReleaseURL   = "https://github.com/lightpanda-io/browser/releases/download/nightly"
	lightpandaDefaultPort  = "9224"
	lightpandaTimeout      = "180"
)

// Lightpanda manages the Lightpanda browser process.
type Lightpanda struct {
	cmd    *exec.Cmd
	binDir string
	port   string
}

// NewLightpanda creates a new Lightpanda manager.
// binDir is where the binary is stored; if empty, defaults to ~/.cache/tap/lightpanda/.
// port is the CDP port; if empty, defaults to 9222.
func NewLightpanda(binDir, port string) *Lightpanda {
	if binDir == "" {
		home, _ := os.UserHomeDir()
		binDir = filepath.Join(home, ".cache", "tap", "lightpanda")
	}
	if port == "" {
		port = lightpandaDefaultPort
	}
	return &Lightpanda{binDir: binDir, port: port}
}

// WSURL returns the WebSocket URL for the running Lightpanda instance.
func (lp *Lightpanda) WSURL() string {
	return "ws://127.0.0.1:" + lp.port + "/"
}

// EnsureInstalled downloads the Lightpanda binary if it doesn't exist or is empty.
func (lp *Lightpanda) EnsureInstalled(ctx context.Context) error {
	bin := lp.binPath()
	if fi, err := os.Stat(bin); err == nil && fi.Size() > 0 {
		return nil // already installed
	}

	log.Println("downloading lightpanda browser")
	return lp.download(ctx)
}

// Running reports whether the Lightpanda server is currently running.
func (lp *Lightpanda) Running() bool {
	return lp.cmd != nil
}

// Start launches the Lightpanda server and waits for it to be ready.
// It is safe to call multiple times — subsequent calls are no-ops if already running.
func (lp *Lightpanda) Start(ctx context.Context) error {
	if lp.cmd != nil {
		return nil
	}

	if err := lp.EnsureInstalled(ctx); err != nil {
		return fmt.Errorf("ensure installed: %w", err)
	}

	// If the configured port is already in use, pick a free one.
	if isPortInUse(lp.port) {
		free, err := freePort()
		if err != nil {
			return fmt.Errorf("find free port: %w", err)
		}
		log.Printf("port %s in use, using %s", lp.port, free)
		lp.port = free
	}

	bin := lp.binPath()
	cmd := exec.CommandContext(ctx, bin, "--port", lp.port, "--timeout", lightpandaTimeout)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start lightpanda: %w", err)
	}

	// Monitor process exit to detect early failures (e.g. port already in use).
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()

	if err := waitForPortOrExit(ctx, lp.port, 10*time.Second, exited); err != nil {
		_ = cmd.Process.Kill()
		<-exited
		return err
	}

	lp.cmd = cmd
	log.Printf("lightpanda browser ready: %s", lp.WSURL())
	return nil
}

// Stop terminates the Lightpanda process.
func (lp *Lightpanda) Stop() {
	if lp.cmd == nil || lp.cmd.Process == nil {
		return
	}
	_ = lp.cmd.Process.Kill()
	_ = lp.cmd.Wait()
	lp.cmd = nil
	log.Println("lightpanda stopped")
}

// Cleanup removes the downloaded binary.
func (lp *Lightpanda) Cleanup() error {
	bin := lp.binPath()
	if err := os.Remove(bin); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove lightpanda binary: %w", err)
	}
	return nil
}

func (lp *Lightpanda) binPath() string {
	return filepath.Join(lp.binDir, "lightpanda")
}

func (lp *Lightpanda) download(ctx context.Context) error {
	url, err := lightpandaNightlyURL()
	if err != nil {
		return err
	}

	// Create dir if needed.
	if err := os.MkdirAll(lp.binDir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	bin := lp.binPath()
	tmp := bin + ".tmp"

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		f.Close()
		os.Remove(tmp) // clean up on failure; no-op after successful rename
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

	log.Printf("lightpanda browser downloaded: %s", bin)
	return nil
}

// lightpandaNightlyURL returns the download URL for the current OS/arch.
func lightpandaNightlyURL() (string, error) {
	var osName string
	switch runtime.GOOS {
	case "linux":
		osName = "linux"
	case "darwin":
		osName = "macos"
	default:
		return "", fmt.Errorf("unsupported OS for lightpanda: %s", runtime.GOOS)
	}

	var arch string
	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	default:
		return "", fmt.Errorf("unsupported arch for lightpanda: %s", runtime.GOARCH)
	}

	return fmt.Sprintf("%s/lightpanda-%s-%s", lightpandaReleaseURL, arch, osName), nil
}

func isPortInUse(port string) bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func freePort() (string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer ln.Close()
	addr := ln.Addr().(*net.TCPAddr)
	return fmt.Sprintf("%d", addr.Port), nil
}

// waitForPortOrExit polls until the port accepts connections or the process exits.
func waitForPortOrExit(ctx context.Context, port string, timeout time.Duration, exited <-chan error) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case err := <-exited:
			if err != nil {
				return fmt.Errorf("lightpanda exited: %w", err)
			}
			return fmt.Errorf("lightpanda exited unexpectedly")
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for port %s", port)
}
