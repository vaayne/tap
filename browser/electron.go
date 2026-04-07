package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

// ElectronProcess describes a running process that exposes a CDP debug port.
type ElectronProcess struct {
	PID  int
	Name string // basename of the binary
	Port int
}

// ScanElectronProcesses returns all running processes that have
// --remote-debugging-port=PORT in their command line. It works for Electron
// apps, CEF-based apps, and any Chromium-derived binary launched with a debug
// port. Returns an empty slice (not an error) when no matching processes exist.
func ScanElectronProcesses(ctx context.Context) ([]ElectronProcess, error) {
	return scanElectronProcesses(ctx)
}

// ResolveElectronDebugURL queries /json/version at the given port and returns
// the WebSocket debug URL. Electron (and Chrome) publish the exact WS URL
// including the browser UUID path in the webSocketDebuggerUrl field.
func ResolveElectronDebugURL(ctx context.Context, port int) (string, error) {
	httpURL := fmt.Sprintf("http://127.0.0.1:%d/json/version", port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("query /json/version at port %d: %w", port, err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	var info struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("parse /json/version: %w", err)
	}
	if info.WebSocketDebuggerURL == "" {
		return "", fmt.Errorf("no webSocketDebuggerUrl in /json/version response at port %d", port)
	}
	return info.WebSocketDebuggerURL, nil
}

// LaunchElectronApp starts an Electron (or Electron-based) binary with
// --remote-debugging-port=0 prepended to the argument list. The OS assigns a
// free port; the debug WebSocket URL is extracted from stderr exactly as with
// Chrome. extra holds any additional arguments to pass to the binary.
//
// The returned ProcessRecord contains the PID and debug URL. The process is
// released (detached) so it survives tap exit; callers are responsible for
// tracking lifecycle via the PID and debug URL.
func LaunchElectronApp(ctx context.Context, binaryPath string, extra []string) (*ProcessRecord, error) {
	args := make([]string, 0, 1+len(extra))
	args = append(args, "--remote-debugging-port=0")
	args = append(args, extra...)

	// Use exec.Command (not CommandContext) so the Electron process is not
	// killed when the CLI context is cancelled — same pattern as LaunchBrowser.
	cmd := exec.Command(binaryPath, args...)
	cmd.SysProcAttr = platformSysProcAttr()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("pipe electron stderr: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start electron app: %w", err)
	}

	// parseDebugURL reuses the same logic as Chrome (both emit the same line).
	debugURL, err := parseDebugURL(stderrPipe, 15*time.Second)
	if c, ok := stderrPipe.(io.Closer); ok {
		_ = c.Close()
	}
	if err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, fmt.Errorf("parse electron debug URL: %w", err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return &ProcessRecord{
		PID:       pid,
		DebugURL:  debugURL,
		StartedAt: time.Now().UTC(),
	}, nil
}
