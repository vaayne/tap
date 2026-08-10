// Package assets contains browser-side helpers embedded in Tap releases.
package assets

import _ "embed"

// DefuddleBrowser is Defuddle's full browser bundle. It includes Markdown
// conversion and is evaluated in the active agent-browser tab.
//
//go:embed defuddle.browser.js
var DefuddleBrowser string
