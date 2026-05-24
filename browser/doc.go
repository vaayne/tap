// Package browser provides the agent-browser adapter for tap.
//
// It manages the agent-browser binary (download, install, version check)
// and provides a thin Go wrapper around the agent-browser CLI with --json
// output. The adapter handles --session-name injection, stdin-based eval,
// and basic typed wrappers for common commands.
//
// Key types:
//
//   - AgentBrowser: the core adapter with Exec(), Open(), Eval(), GetHTML(), Close()
//   - AgentBrowserInstall: downloads and manages the agent-browser binary
//   - AgentBrowserEnvelope[T]: generic JSON response envelope
//
// Binary resolution order:
//  1. $TAP_AGENT_BROWSER env var
//  2. agent-browser on $PATH
//  3. ~/.cache/tap/agent-browser/agent-browser
//
// Session model:
//   - Default sessions use --session-name default for durable persistence
//   - Attached mode uses agent-browser connect (no --session-name)
package browser
