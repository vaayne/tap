//go:build darwin && amd64

package browser

import _ "embed"

//go:embed bin/agent-browser-darwin-x64
var embeddedAgentBrowser []byte
