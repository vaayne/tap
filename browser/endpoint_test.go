package browser

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveDebugURLFromEndpoint(t *testing.T) {
	t.Run("passes through websocket endpoints", func(t *testing.T) {
		got, err := ResolveDebugURL(context.Background(), "ws://127.0.0.1:9222/devtools/browser/abc")
		if err != nil {
			t.Fatalf("ResolveDebugURL returned error: %v", err)
		}
		want := "ws://127.0.0.1:9222/devtools/browser/abc"
		if got != want {
			t.Fatalf("ResolveDebugURL = %q, want %q", got, want)
		}
	})

	t.Run("resolves browser websocket from http endpoint", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/json/version" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			_, _ = fmt.Fprint(w, `{"webSocketDebuggerUrl":"ws://127.0.0.1:9401/devtools/browser/proxy"}`)
		}))
		defer server.Close()

		got, err := ResolveDebugURL(context.Background(), server.URL)
		if err != nil {
			t.Fatalf("ResolveDebugURL returned error: %v", err)
		}
		want := "ws://127.0.0.1:9401/devtools/browser/proxy"
		if got != want {
			t.Fatalf("ResolveDebugURL = %q, want %q", got, want)
		}
	})

	t.Run("strips extra path from http endpoint", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/json/version" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			_, _ = fmt.Fprint(w, `{"webSocketDebuggerUrl":"ws://127.0.0.1:9222/devtools/browser/xyz"}`)
		}))
		defer server.Close()

		got, err := ResolveDebugURL(context.Background(), server.URL+"/devtools/browser/ignored")
		if err != nil {
			t.Fatalf("ResolveDebugURL returned error: %v", err)
		}
		want := "ws://127.0.0.1:9222/devtools/browser/xyz"
		if got != want {
			t.Fatalf("ResolveDebugURL = %q, want %q", got, want)
		}
	})

	t.Run("rejects unsupported schemes", func(t *testing.T) {
		if _, err := ResolveDebugURL(context.Background(), "ftp://example.com"); err == nil {
			t.Fatal("ResolveDebugURL should reject unsupported schemes")
		}
	})
}
