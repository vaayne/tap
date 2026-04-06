package transport

import (
	"context"
	"fmt"
	"net"
	"net/http"

	utls "github.com/refraction-networking/utls"
)

// UserAgent is a realistic Chrome User-Agent sent with all direct HTTP requests.
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

// newHTTPClient returns an *http.Client whose TLS fingerprint mimics Chrome
// instead of the default Go fingerprint which is trivially detected by bots.
func newHTTPClient() *http.Client {
	dialer := &net.Dialer{}

	return &http.Client{
		Transport: &http.Transport{
			Proxy:       http.ProxyFromEnvironment,
			DialContext: dialer.DialContext,
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
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
			},
		},
	}
}
