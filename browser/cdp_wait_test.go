package browser

import (
	"testing"
)

func TestMatchURLGlob(t *testing.T) {
	tests := []struct {
		pattern string
		url     string
		want    bool
	}{
		// Exact substring match (no wildcards)
		{"example.com", "https://example.com/path", true},
		{"notfound", "https://example.com/path", false},

		// Single-star glob (does not cross /)
		{"https://example.com/*", "https://example.com/page", true},
		{"https://example.com/*", "https://example.com/a/b", false},

		// Double-star glob (crosses /)
		{"**/dashboard", "https://app.example.com/user/dashboard", true},
		{"**/dashboard", "https://app.example.com/other", false},
		{"https://**", "https://example.com/path", true},
		{"https://**/page", "https://example.com/nested/deeply/page", true},

		// Mixed
		{"**/user/**/profile", "https://app.example.com/user/123/profile", true},
		{"**/user/**/profile", "https://app.example.com/team/profile", false},
	}

	for _, tt := range tests {
		got := matchURLGlob(tt.pattern, tt.url)
		if got != tt.want {
			t.Errorf("matchURLGlob(%q, %q) = %v; want %v", tt.pattern, tt.url, got, tt.want)
		}
	}
}
