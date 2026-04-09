package browser

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	AttachStateDetached      = "detached"
	AttachStateAttachedReady = "attached-ready"
	AttachStateAttachedStale = "attached-stale"
)

type ProxyDaemonHealth struct {
	CheckedAt         time.Time `json:"checked_at"`
	PIDAlive          bool      `json:"pid_alive"`
	DaemonReachable   bool      `json:"daemon_reachable"`
	UpstreamReachable bool      `json:"upstream_reachable"`
	Healthy           bool      `json:"healthy"`
	Status            string    `json:"status"`
	Reason            string    `json:"reason,omitempty"`
}

func GenerateOwnershipToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate ownership token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func ProxyEndpointForListenAddr(listenAddr string) string {
	return fmt.Sprintf("ws://%s/devtools/browser/proxy", listenAddr)
}

func ProxyHTTPBaseForListenAddr(listenAddr string) string {
	return fmt.Sprintf("http://%s", listenAddr)
}

func CheckProxyDaemon(ctx context.Context, record *ProxyDaemonRecord) ProxyDaemonHealth {
	health := ProxyDaemonHealth{CheckedAt: time.Now().UTC(), Status: AttachStateDetached}
	if record == nil {
		health.Reason = "no proxy daemon recorded"
		return health
	}

	var reasons []string
	health.Status = AttachStateAttachedStale
	if record.PID > 0 && isProcessAlive(record.PID) {
		health.PIDAlive = true
	} else {
		reasons = append(reasons, "daemon pid is not alive")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	if strings.TrimSpace(record.ListenAddr) != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ProxyHTTPBaseForListenAddr(record.ListenAddr)+"/healthz", nil)
		if err == nil {
			resp, err := client.Do(req)
			if err == nil {
				var payload struct {
					UpstreamReady bool `json:"upstreamReady"`
				}
				body, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				_ = json.Unmarshal(body, &payload)
				if resp.StatusCode == http.StatusOK {
					health.DaemonReachable = true
				}
				if payload.UpstreamReady {
					health.UpstreamReachable = true
				}
			} else {
				reasons = append(reasons, fmt.Sprintf("daemon health check failed: %v", err))
			}
		} else {
			reasons = append(reasons, fmt.Sprintf("daemon health request failed: %v", err))
		}
	} else {
		reasons = append(reasons, "daemon listen address is missing")
	}

	if !health.DaemonReachable {
		reasons = append(reasons, "daemon endpoint is unreachable")
	}
	if !health.UpstreamReachable {
		reasons = append(reasons, "upstream Chrome is unreachable")
	}
	if health.PIDAlive && health.DaemonReachable && health.UpstreamReachable {
		health.Healthy = true
		health.Status = AttachStateAttachedReady
		health.Reason = ""
		return health
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "proxy daemon is unhealthy")
	}
	health.Reason = strings.Join(uniqueStrings(reasons), "; ")
	return health
}

func ShouldReuseProxyDaemon(record *ProxyDaemonRecord, discoveredUpstreamWSURL string, health ProxyDaemonHealth) bool {
	return record != nil && record.UpstreamWSURL == strings.TrimSpace(discoveredUpstreamWSURL) && health.Healthy
}

func ShutdownProxyDaemon(ctx context.Context, record *ProxyDaemonRecord) error {
	if record == nil {
		return nil
	}
	if strings.TrimSpace(record.ListenAddr) == "" {
		return errors.New("proxy daemon listen address is missing")
	}
	if strings.TrimSpace(record.OwnershipToken) == "" {
		return errors.New("proxy daemon ownership token is missing")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ProxyHTTPBaseForListenAddr(record.ListenAddr)+"/shutdown", nil)
	if err != nil {
		return fmt.Errorf("build proxy shutdown request: %w", err)
	}
	req.Header.Set("X-Tap-Ownership-Token", record.OwnershipToken)
	resp, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("request proxy shutdown: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxy shutdown returned status %d", resp.StatusCode)
	}
	return nil
}

func WaitForProxyDaemonExit(ctx context.Context, record *ProxyDaemonRecord) error {
	if record == nil || record.PID <= 0 {
		return nil
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if !isProcessAlive(record.PID) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func WaitForProxyDaemonHealth(ctx context.Context, record *ProxyDaemonRecord) (ProxyDaemonHealth, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		health := CheckProxyDaemon(ctx, record)
		if health.Healthy {
			return health, nil
		}
		select {
		case <-ctx.Done():
			return health, ctx.Err()
		case <-ticker.C:
		}
	}
}

func IsListenAddrInUse(err error) bool {
	var opErr *net.OpError
	return errors.As(err, &opErr) && strings.Contains(strings.ToLower(opErr.Err.Error()), "address already in use") || strings.Contains(strings.ToLower(err.Error()), "address already in use")
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
