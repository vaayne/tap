//go:build !(darwin && arm64) && !(darwin && amd64) && !(linux && arm64) && !(linux && amd64) && !(windows && amd64)

package browser

var embeddedAgentBrowser []byte
