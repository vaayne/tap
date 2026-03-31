package tap

// options holds the resolved configuration for a Client.
type options struct {
	sitesDir   string
	wsURL      string
	profileDir string
}

func defaultOptions() options {
	return options{}
}

// Option configures a Client.
type Option func(*options)

// WithSitesDir sets the directory containing site scripts.
func WithSitesDir(dir string) Option {
	return func(o *options) {
		o.sitesDir = dir
	}
}

// WithWSURL sets the remote CDP WebSocket URL.
// If empty, a local Chrome is launched.
func WithWSURL(url string) Option {
	return func(o *options) {
		o.wsURL = url
	}
}

// WithProfileDir sets the Chrome user data directory for persistent cookies/storage.
// Defaults to ~/.cache/tap/chrome-profile-$USER.
func WithProfileDir(dir string) Option {
	return func(o *options) {
		o.profileDir = dir
	}
}
