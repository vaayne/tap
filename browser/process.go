package browser

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Chrome binary names searched via exec.LookPath (Linux).
var chromePathNames = []string{
	"google-chrome",
	"google-chrome-stable",
	"chromium-browser",
	"chromium",
}

// Known absolute paths checked as a fallback on macOS.
var chromeDarwinPaths = []string{
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Chromium.app/Contents/MacOS/Chromium",
}

// Known absolute paths checked as a fallback on Linux.
var chromeLinuxPaths = []string{
	"/usr/bin/google-chrome",
	"/usr/bin/google-chrome-stable",
	"/usr/bin/chromium-browser",
	"/usr/bin/chromium",
}

// findChrome discovers the first available Chrome or Chromium binary.
func findChrome() (string, error) {
	// Try PATH lookup first (works on all platforms).
	for _, name := range chromePathNames {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}

	// Fall back to well-known absolute paths.
	var fallbacks []string
	switch runtime.GOOS {
	case "darwin":
		fallbacks = chromeDarwinPaths
	default:
		fallbacks = chromeLinuxPaths
	}
	for _, p := range fallbacks {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", errors.New("chrome binary not found: install Google Chrome or Chromium")
}

// LaunchBrowser starts a managed Chrome process with remote debugging enabled.
// The returned ProcessRecord contains the PID, debug WebSocket URL, and an
// ownership token that later calls can use for verification.
func LaunchBrowser(ctx context.Context, config LocalConfig) (*ProcessRecord, error) {
	chromePath, err := findChrome()
	if err != nil {
		return nil, err
	}

	// Generate a unique ownership token (16 random bytes, hex-encoded).
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate ownership token: %w", err)
	}
	ownershipToken := hex.EncodeToString(tokenBytes)

	// Ensure the profile directory exists.
	if err := os.MkdirAll(config.ProfileDir, 0o755); err != nil {
		return nil, fmt.Errorf("create browser profile dir: %w", err)
	}

	// Build Chrome args.
	args := []string{
		"--remote-debugging-port=0",
		"--user-data-dir=" + config.ProfileDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-gpu",
	}
	if config.Headless {
		args = append(args, "--headless=new")
	}

	cmd := exec.Command(chromePath, args...)

	// Run Chrome in its own process group so it survives tap exit.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Pipe stderr to parse the DevTools debug URL.
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("pipe chrome stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start chrome: %w", err)
	}

	// Parse the debug URL from Chrome's stderr output with a timeout.
	debugURL, err := parseDebugURL(stderrPipe, 10*time.Second)
	if err != nil {
		// Best-effort cleanup: kill and reap to avoid zombies and close pipes.
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, fmt.Errorf("parse chrome debug URL: %w", err)
	}

	// Capture PID before Release, which sets Pid to -1.
	pid := cmd.Process.Pid

	// Release the process handle so Go does not accumulate zombies.
	// Chrome runs detached (Setpgid) and is managed via PID from metadata.
	_ = cmd.Process.Release()

	return &ProcessRecord{
		PID:            pid,
		DebugURL:       debugURL,
		OwnershipToken: ownershipToken,
		StartedAt:      time.Now().UTC(),
	}, nil
}

// parseDebugURL reads from r line by line looking for the DevTools listening
// message. It returns the WebSocket URL or an error if the timeout expires.
func parseDebugURL(r io.Reader, timeout time.Duration) (string, error) {
	type result struct {
		url string
		err error
	}
	ch := make(chan result, 1)

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "DevTools listening on ") {
				ch <- result{url: strings.TrimPrefix(line, "DevTools listening on ")}
				return
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- result{err: fmt.Errorf("read stderr: %w", err)}
			return
		}
		ch <- result{err: errors.New("chrome stderr closed without debug URL")}
	}()

	select {
	case res := <-ch:
		return res.url, res.err
	case <-time.After(timeout):
		return "", errors.New("timed out waiting for chrome debug URL")
	}
}

// CheckProcess verifies that the process described by record is alive and
// its debug endpoint is reachable.
func CheckProcess(record *ProcessRecord) error {
	if record == nil || record.PID <= 0 {
		return errors.New("no process record to check")
	}

	// Signal 0 checks whether the process exists without sending a real signal.
	if err := syscall.Kill(record.PID, 0); err != nil {
		return fmt.Errorf("process %d is not alive: %w", record.PID, err)
	}

	// Verify the debug endpoint is reachable.
	httpURL, err := debugURLToHTTP(record.DebugURL)
	if err != nil {
		return fmt.Errorf("parse debug URL: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(httpURL + "/json/version")
	if err != nil {
		return fmt.Errorf("debug endpoint unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Drain the body to allow connection reuse.
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("debug endpoint returned status %d", resp.StatusCode)
	}

	return nil
}

// KillProcess terminates the Chrome process described by record. It sends
// SIGTERM first and falls back to SIGKILL after a 5-second grace period.
func KillProcess(record *ProcessRecord) error {
	if record == nil || record.PID <= 0 {
		return nil
	}

	// Check whether the process is still alive.
	if err := syscall.Kill(record.PID, 0); err != nil {
		return nil // already dead
	}

	// Verify we own this process by checking the debug endpoint.
	// If the endpoint is unreachable (PID reuse), skip termination.
	if record.DebugURL != "" {
		if err := CheckProcess(record); err != nil {
			return nil // can't verify ownership — skip kill
		}
	}

	// Send SIGTERM for a graceful shutdown.
	if err := syscall.Kill(record.PID, syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return fmt.Errorf("send SIGTERM to %d: %w", record.PID, err)
	}

	// Wait up to 5 seconds for the process to exit.
	deadline := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			// Process didn't exit in time — force kill.
			if err := syscall.Kill(record.PID, syscall.SIGKILL); err != nil {
				// If ESRCH, the process exited between our check and kill.
				if errors.Is(err, syscall.ESRCH) {
					return nil
				}
				return fmt.Errorf("send SIGKILL to %d: %w", record.PID, err)
			}
			return nil
		case <-ticker.C:
			if err := syscall.Kill(record.PID, 0); err != nil {
				return nil // process exited
			}
		}
	}
}

// debugURLToHTTP converts a DevTools WebSocket URL to its HTTP equivalent.
// For example, ws://127.0.0.1:9222/devtools/browser/GUID becomes
// http://127.0.0.1:9222.
func debugURLToHTTP(debugURL string) (string, error) {
	if debugURL == "" {
		return "", errors.New("debug URL is empty")
	}

	u, err := url.Parse(debugURL)
	if err != nil {
		return "", fmt.Errorf("parse debug URL: %w", err)
	}

	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return "", fmt.Errorf("unexpected debug URL scheme %q", u.Scheme)
	}

	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}
