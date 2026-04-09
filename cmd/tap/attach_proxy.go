package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap/browser"
)

func newBrowserStore(cmd *cli.Command) (*browser.Store, error) {
	root := browserStateRoot(cmd)
	store, err := browser.NewStore(root)
	if err != nil {
		return nil, fmt.Errorf("init browser store: %w", err)
	}
	return store, nil
}

func ensureAttachedChromeProxy(ctx context.Context, cmd *cli.Command, upstreamWSURL, listenAddr string) (*browser.ProxyDaemonRecord, browser.ProxyDaemonHealth, bool, error) {
	store, err := newBrowserStore(cmd)
	if err != nil {
		return nil, browser.ProxyDaemonHealth{}, false, err
	}
	state, err := store.Load()
	if err != nil {
		return nil, browser.ProxyDaemonHealth{}, false, fmt.Errorf("load browser state: %w", err)
	}

	if existing := state.ProxyDaemon; existing != nil {
		health := browser.CheckProxyDaemon(ctx, existing)
		if existing.UpstreamWSURL == upstreamWSURL && health.Healthy {
			_ = persistProxyDaemonHealth(store, existing, health)
			return existing, health, true, nil
		}
		if existing.UpstreamWSURL != upstreamWSURL {
			health.Healthy = false
			health.Status = browser.AttachStateAttachedStale
			if health.Reason == "" {
				health.Reason = fmt.Sprintf("saved upstream changed from %s to %s", existing.UpstreamWSURL, upstreamWSURL)
			} else {
				health.Reason += fmt.Sprintf("; saved upstream changed from %s to %s", existing.UpstreamWSURL, upstreamWSURL)
			}
		}
		_ = persistProxyDaemonHealth(store, existing, health)
		_ = stopOwnedProxyDaemon(ctx, existing)
	}

	record, health, err := startAttachedChromeProxy(ctx, cmd, listenAddr, upstreamWSURL)
	if err != nil {
		return nil, health, false, err
	}
	if err := persistProxyDaemonRecord(store, record, health); err != nil {
		return nil, health, false, err
	}
	return record, health, false, nil
}

func persistProxyDaemonRecord(store *browser.Store, record *browser.ProxyDaemonRecord, health browser.ProxyDaemonHealth) error {
	return store.Update(func(state *browser.State) error {
		record.State = health.Status
		record.Status = health.Status
		record.LastHealthCheckAt = health.CheckedAt
		if health.Healthy {
			record.LastHealthyAt = health.CheckedAt
			record.LastError = ""
		} else {
			record.LastError = health.Reason
		}
		record.UpdatedAt = time.Now().UTC()
		state.ProxyDaemon = record
		if state.DefaultContext != nil && state.DefaultContext.Kind == browser.DefaultContextAttached {
			if health.Healthy {
				state.MarkDefaultContextHealthy(state.DefaultContext.SessionName, health.CheckedAt)
			} else {
				state.MarkDefaultContextStale(state.DefaultContext.SessionName, health.Reason, health.CheckedAt)
			}
		}
		return nil
	})
}

func persistProxyDaemonHealth(store *browser.Store, record *browser.ProxyDaemonRecord, health browser.ProxyDaemonHealth) error {
	return store.Update(func(state *browser.State) error {
		if state.ProxyDaemon == nil {
			return nil
		}
		state.ProxyDaemon.State = health.Status
		state.ProxyDaemon.Status = health.Status
		state.ProxyDaemon.LastHealthCheckAt = health.CheckedAt
		if health.Healthy {
			state.ProxyDaemon.LastHealthyAt = health.CheckedAt
			state.ProxyDaemon.LastError = ""
		} else {
			state.ProxyDaemon.LastError = health.Reason
		}
		state.ProxyDaemon.UpdatedAt = time.Now().UTC()
		if state.DefaultContext != nil && state.DefaultContext.Kind == browser.DefaultContextAttached {
			if health.Healthy {
				state.MarkDefaultContextHealthy(state.DefaultContext.SessionName, health.CheckedAt)
			} else {
				state.MarkDefaultContextStale(state.DefaultContext.SessionName, health.Reason, health.CheckedAt)
			}
		}
		return nil
	})
}

func stopOwnedProxyDaemon(ctx context.Context, record *browser.ProxyDaemonRecord) error {
	if record == nil {
		return nil
	}
	if err := browser.ShutdownProxyDaemon(ctx, record); err != nil {
		return err
	}
	waitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return browser.WaitForProxyDaemonExit(waitCtx, record)
}

func startAttachedChromeProxy(ctx context.Context, cmd *cli.Command, listenAddr, upstreamWSURL string) (*browser.ProxyDaemonRecord, browser.ProxyDaemonHealth, error) {
	token, err := browser.GenerateOwnershipToken()
	if err != nil {
		return nil, browser.ProxyDaemonHealth{}, err
	}

	record, err := spawnProxyDaemon(ctx, listenAddr, upstreamWSURL, token)
	if err != nil && listenAddr == browser.DefaultProxyListenAddr {
		record, err = spawnProxyDaemon(ctx, "127.0.0.1:0", upstreamWSURL, token)
	}
	if err != nil {
		return nil, browser.ProxyDaemonHealth{}, err
	}

	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	health, err := browser.WaitForProxyDaemonHealth(healthCtx, record)
	if err != nil {
		return record, health, fmt.Errorf("wait for proxy daemon health: %w", err)
	}
	return record, health, nil
}

func spawnProxyDaemon(ctx context.Context, listenAddr, upstreamWSURL, ownershipToken string) (*browser.ProxyDaemonRecord, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("resolve tap executable: %w", err)
	}

	readyFile := filepath.Join(os.TempDir(), fmt.Sprintf("tap-proxy-daemon-%d-%d.json", os.Getpid(), time.Now().UnixNano()))
	defer func() { _ = os.Remove(readyFile) }()

	args := []string{
		"internal", "proxy-daemon",
		"--listen", listenAddr,
		"--upstream-ws-url", upstreamWSURL,
		"--ownership-token", ownershipToken,
		"--ready-file", readyFile,
	}
	proc := exec.Command(exe, args...)
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", os.DevNull, err)
	}
	defer func() { _ = devNull.Close() }()
	proc.Stdin = devNull
	proc.Stdout = devNull
	proc.Stderr = devNull
	if err := proc.Start(); err != nil {
		return nil, fmt.Errorf("start proxy daemon: %w", err)
	}
	_ = proc.Process.Release()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(readyFile)
		if err == nil {
			var record browser.ProxyDaemonRecord
			if err := json.Unmarshal(data, &record); err != nil {
				return nil, fmt.Errorf("decode proxy daemon ready file: %w", err)
			}
			return &record, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return nil, fmt.Errorf("timed out waiting for proxy daemon startup")
}
