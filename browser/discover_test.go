package browser

import "testing"

func TestParseDevToolsActivePort(t *testing.T) {
	t.Run("valid content", func(t *testing.T) {
		got, err := parseDevToolsActivePort("9222\n/devtools/browser/abc\n")
		if err != nil {
			t.Fatalf("parseDevToolsActivePort returned error: %v", err)
		}
		want := "ws://127.0.0.1:9222/devtools/browser/abc"
		if got != want {
			t.Fatalf("parseDevToolsActivePort = %q, want %q", got, want)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		if _, err := parseDevToolsActivePort("abc\n/devtools/browser/abc\n"); err == nil {
			t.Fatal("parseDevToolsActivePort should reject invalid ports")
		}
	})

	t.Run("invalid websocket path", func(t *testing.T) {
		if _, err := parseDevToolsActivePort("9222\ndevtools/browser/abc\n"); err == nil {
			t.Fatal("parseDevToolsActivePort should reject invalid websocket paths")
		}
	})
}
