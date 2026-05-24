//go:build linux && arm64 && !musl

package browser

import _ "embed"

//go:embed bin/agent-browser-linux-arm64
var embeddedAgentBrowser []byte
