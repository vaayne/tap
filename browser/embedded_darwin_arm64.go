//go:build darwin && arm64

package browser

import _ "embed"

//go:embed bin/agent-browser-darwin-arm64
var embeddedAgentBrowser []byte
