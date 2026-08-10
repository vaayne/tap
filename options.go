package tap

import (
	"time"
)

// options holds the resolved configuration for a Client.
type options struct {
	sitesDir         string
	localOverrideDir string
	agentBrowser     string
	timeout          time.Duration
}

func defaultOptions() options {
	return options{}
}

// Option configures a Client.
type Option func(*options)

// WithLocalOverrideDir sets a directory that is checked before the main sites
// cache. Scripts found here shadow cached versions and are flagged as local
// overrides. Mirrors the path structure: {dir}/{site}/{script}.js
func WithLocalOverrideDir(dir string) Option {
	return func(o *options) {
		o.localOverrideDir = dir
	}
}

// WithSitesDir sets the directory containing site scripts.
func WithSitesDir(dir string) Option {
	return func(o *options) {
		o.sitesDir = dir
	}
}

// WithAgentBrowserBinary overrides the agent-browser executable. The default
// is TAP_AGENT_BROWSER, then agent-browser from PATH.
func WithAgentBrowserBinary(binary string) Option {
	return func(o *options) {
		o.agentBrowser = binary
	}
}

// WithTimeout sets the execution timeout for scripts and fetches.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}
