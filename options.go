package tap

import "time"

// options holds the resolved configuration for a Client.
type options struct {
	sitesDir    string
	wsURL       string
	profileDir  string
	forceBrowser bool
	headless     bool
	timeout      time.Duration
}

func defaultOptions() options {
	return options{
		headless: true,
	}
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

// WithForceBrowser skips QuickJS and runs scripts directly in Chrome.
func WithForceBrowser(force bool) Option {
	return func(o *options) {
		o.forceBrowser = force
	}
}

// WithHeadless sets whether Chrome runs in headless mode (default: true).
func WithHeadless(headless bool) Option {
	return func(o *options) {
		o.headless = headless
	}
}

// WithTimeout sets the execution timeout for scripts and fetches.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}
