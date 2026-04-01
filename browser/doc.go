// Package browser defines the persistent browser session metadata model used by
// the planned `tap browser ...` workflow.
//
// Phase 1 establishes the durable state contract before the runtime is added:
//
//   - A session is one persistent browser instance, either local or remote.
//   - A tab is a named tracked browser target within a session.
//   - Stored metadata is the source of truth for names and selection defaults.
//   - Each command must reconcile metadata against live targets before acting.
//   - Untracked live browser tabs are ignored by default.
//   - Missing tracked targets become stale until they are recreated or removed.
//
// The package also defines the initial local-vs-remote capability matrix used
// by CLI help text and README documentation:
//
//   - Local sessions own their browser process and profile directory.
//   - Remote sessions bind metadata to an explicit CDP WebSocket endpoint.
//   - Remote session creation validates the endpoint up front.
//   - Remote session close removes tap metadata only and never kills the remote
//     browser process.
package browser
