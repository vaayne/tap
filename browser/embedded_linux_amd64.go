//go:build linux && amd64 && !musl

package browser

import _ "embed"

//go:embed bin/agent-browser-linux-x64
var embeddedAgentBrowser []byte
