package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

// UserAgent is a realistic Chrome User-Agent sent with all direct HTTP requests.
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

// newHTTPClient returns an *http.Client whose TLS fingerprint mimics Chrome
// instead of the default Go fingerprint which is trivially detected by bots.
func newHTTPClient() *http.Client {
	dialer := &net.Dialer{}

	dialTLS := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, fmt.Errorf("dial: %w", err)
		}

		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("split host port: %w", err)
		}

		tlsConn := utls.UClient(conn, &utls.Config{
			ServerName: host,
		}, utls.HelloChrome_131)

		if err := tlsConn.HandshakeContext(ctx); err != nil {
			_ = tlsConn.Close()
			return nil, fmt.Errorf("tls handshake: %w", err)
		}

		return tlsConn, nil
	}

	h1 := &http.Transport{
		Proxy:          http.ProxyFromEnvironment,
		DialContext:    dialer.DialContext,
		DialTLSContext: dialTLS,
	}

	h2 := &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return dialTLS(ctx, network, addr)
		},
	}

	return &http.Client{Transport: &utlsTransport{h1: h1, h2: h2}}
}

// utlsTransport uses h2 for HTTPS (Chrome always negotiates h2) and h1 for
// plain HTTP. If a rare h1-only HTTPS server rejects h2, the fetch layer's
// existing browser fallback handles it.
type utlsTransport struct {
	h1 *http.Transport
	h2 *http2.Transport
}

func (t *utlsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" {
		return t.h1.RoundTrip(req)
	}
	return t.h2.RoundTrip(req)
}
