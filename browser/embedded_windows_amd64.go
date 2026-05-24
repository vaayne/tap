//go:build windows && amd64

package browser

import _ "embed"

//go:embed bin/agent-browser-win32-x64.exe
var embeddedAgentBrowser []byte
