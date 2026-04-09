package browser

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func startFakeUpstreamWS(t *testing.T) (string, func()) {
	t.Helper()
	upgrader := websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/devtools/browser/test" {
			http.NotFound(w, r)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}))
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/devtools/browser/test"
	return wsURL, srv.Close
}

func startProxyForTest(t *testing.T, upstream, token string) (*ProxyDaemonRecord, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	record := &ProxyDaemonRecord{
		PID:            os.Getpid(),
		ListenAddr:     ln.Addr().String(),
		Endpoint:       ProxyEndpointForListenAddr(ln.Addr().String()),
		UpstreamWSURL:  upstream,
		OwnershipToken: token,
	}
	proxy := NewProxy(ProxyConfig{ListenAddr: record.ListenAddr, Upstream: upstream, OwnershipToken: token})
	errCh := make(chan error, 1)
	go func() {
		errCh <- proxy.ServeListener(t.Context(), ln)
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	if _, err := WaitForProxyDaemonHealth(ctx, record); err != nil {
		t.Fatalf("WaitForProxyDaemonHealth failed: %v", err)
	}
	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = ShutdownProxyDaemon(shutdownCtx, record)
		select {
		case <-time.After(250 * time.Millisecond):
		case <-errCh:
		}
	}
	return record, cleanup
}

func TestStorePersistsProxyDaemonState(t *testing.T) {
	store := testStore(t)
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	want := &ProxyDaemonRecord{
		PID:               42,
		ListenAddr:        "127.0.0.1:12345",
		Endpoint:          "ws://127.0.0.1:12345/devtools/browser/proxy",
		UpstreamWSURL:     "ws://127.0.0.1:9222/devtools/browser/test",
		OwnershipToken:    "token-1",
		State:             AttachStateAttachedReady,
		Status:            AttachStateAttachedReady,
		LastHealthCheckAt: now,
		LastHealthyAt:     now,
		StartedAt:         now,
		UpdatedAt:         now,
	}
	if err := store.Update(func(state *State) error {
		state.ProxyDaemon = want
		return nil
	}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	reloaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if reloaded.ProxyDaemon == nil {
		t.Fatal("ProxyDaemon = nil, want record")
	}
	if reloaded.ProxyDaemon.Endpoint != want.Endpoint {
		t.Fatalf("Endpoint = %q, want %q", reloaded.ProxyDaemon.Endpoint, want.Endpoint)
	}
	if reloaded.ProxyDaemon.OwnershipToken != want.OwnershipToken {
		t.Fatalf("OwnershipToken = %q, want %q", reloaded.ProxyDaemon.OwnershipToken, want.OwnershipToken)
	}
}

func TestShouldReuseProxyDaemon(t *testing.T) {
	record := &ProxyDaemonRecord{UpstreamWSURL: "ws://127.0.0.1:9222/devtools/browser/test"}
	health := ProxyDaemonHealth{Healthy: true}
	if !ShouldReuseProxyDaemon(record, record.UpstreamWSURL, health) {
		t.Fatal("ShouldReuseProxyDaemon = false, want true")
	}
	if ShouldReuseProxyDaemon(record, "ws://127.0.0.1:9223/devtools/browser/test", health) {
		t.Fatal("ShouldReuseProxyDaemon should be false for changed upstream")
	}
	if ShouldReuseProxyDaemon(record, record.UpstreamWSURL, ProxyDaemonHealth{}) {
		t.Fatal("ShouldReuseProxyDaemon should be false for unhealthy daemon")
	}
}

func TestCheckProxyDaemonDetectsStaleReasons(t *testing.T) {
	health := CheckProxyDaemon(t.Context(), &ProxyDaemonRecord{
		PID:            999999,
		ListenAddr:     "127.0.0.1:1",
		Endpoint:       "ws://127.0.0.1:1/devtools/browser/proxy",
		UpstreamWSURL:  "ws://127.0.0.1:2/devtools/browser/test",
		OwnershipToken: "token",
	})
	if health.Healthy {
		t.Fatal("health.Healthy = true, want false")
	}
	if !strings.Contains(health.Reason, "daemon pid is not alive") {
		t.Fatalf("health.Reason = %q, want pid reason", health.Reason)
	}
	if !strings.Contains(health.Reason, "daemon endpoint is unreachable") {
		t.Fatalf("health.Reason = %q, want endpoint reason", health.Reason)
	}
}

func TestShutdownProxyDaemonRequiresOwnershipToken(t *testing.T) {
	upstream, closeUpstream := startFakeUpstreamWS(t)
	defer closeUpstream()
	record, cleanup := startProxyForTest(t, upstream, "correct-token")
	defer cleanup()

	wrong := *record
	wrong.OwnershipToken = "wrong-token"
	shutdownCtx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if err := ShutdownProxyDaemon(shutdownCtx, &wrong); err == nil {
		t.Fatal("ShutdownProxyDaemon with wrong token should fail")
	}
	health := CheckProxyDaemon(t.Context(), record)
	if !health.Healthy {
		t.Fatalf("health after wrong-token shutdown = %+v, want healthy", health)
	}

	shutdownCtx, cancel = context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if err := ShutdownProxyDaemon(shutdownCtx, record); err != nil {
		t.Fatalf("ShutdownProxyDaemon failed: %v", err)
	}
}
