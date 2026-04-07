package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ResolveDebugURL accepts either a CDP browser WebSocket URL or an HTTP(S)
// DevTools base URL (for example http://127.0.0.1:9222) and returns the browser
// WebSocket URL that chromedp expects.
func ResolveDebugURL(ctx context.Context, endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("debug endpoint is empty")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse debug endpoint: %w", err)
	}

	switch u.Scheme {
	case "ws", "wss":
		return endpoint, nil
	case "http", "https":
		base := strings.TrimRight(endpoint, "/")
		if u.Path != "" && u.Path != "/" {
			u.Path = ""
			u.RawQuery = ""
			u.Fragment = ""
			base = strings.TrimRight(u.String(), "/")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/json/version", nil)
		if err != nil {
			return "", fmt.Errorf("build debug endpoint request: %w", err)
		}

		resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
		if err != nil {
			return "", fmt.Errorf("resolve debug endpoint: %w", err)
		}
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("resolve debug endpoint: status %d", resp.StatusCode)
		}

		var payload struct {
			WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return "", fmt.Errorf("decode debug endpoint response: %w", err)
		}
		if strings.TrimSpace(payload.WebSocketDebuggerURL) == "" {
			return "", fmt.Errorf("debug endpoint response missing webSocketDebuggerUrl")
		}
		return payload.WebSocketDebuggerURL, nil
	default:
		return "", fmt.Errorf("unsupported debug endpoint scheme %q", u.Scheme)
	}
}
