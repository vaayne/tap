package browser

import (
	"errors"
	"testing"
)

func TestProxyErrorResponseUsesOriginalID(t *testing.T) {
	msg := map[string]any{"id": int64(42)}
	err := &proxyClientMessageError{
		originalID: 7,
		sessionID:  "session-1",
		err:        errors.New("upstream browser not connected"),
	}

	resp := proxyErrorResponse(msg, err)
	if got, want := resp["id"], int64(7); got != want {
		t.Fatalf("proxyErrorResponse id = %#v, want %#v", got, want)
	}
	if got, want := resp["sessionId"], "session-1"; got != want {
		t.Fatalf("proxyErrorResponse sessionId = %#v, want %#v", got, want)
	}
}

func TestProxyErrorResponseFallsBackToMessageID(t *testing.T) {
	msg := map[string]any{"id": int64(42)}
	resp := proxyErrorResponse(msg, errors.New("boom"))
	if got, want := resp["id"], int64(42); got != want {
		t.Fatalf("proxyErrorResponse id = %#v, want %#v", got, want)
	}
	if _, ok := resp["sessionId"]; ok {
		t.Fatal("proxyErrorResponse should not set sessionId for generic errors")
	}
}
