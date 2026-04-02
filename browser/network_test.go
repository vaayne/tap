package browser

import (
	"testing"
)

func TestMatchURL(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		url     string
		want    bool
	}{
		// Empty pattern matches everything.
		{"empty pattern", "", "https://example.com/api/v1/users", true},
		{"empty pattern empty url", "", "", true},

		// Exact match.
		{"exact match", "https://example.com/api", "https://example.com/api", true},
		{"exact no match", "https://example.com/api", "https://example.com/other", false},

		// Single * matches across path segments (unlike filepath.Match).
		{"star prefix", "*/api/*", "https://example.com/api/v1/users", true},
		{"star suffix", "https://example.com/*", "https://example.com/api/v1/users", true},
		{"star middle", "https://*/api/*", "https://example.com/api/v1/users", true},
		{"star only", "*", "https://example.com/anything", true},
		{"star no match", "*/api/*", "https://example.com/other/v1", false},

		// Multiple stars.
		{"multi star", "*api*users*", "https://example.com/api/v1/users?page=1", true},
		{"multi star no match", "*api*admin*", "https://example.com/api/v1/users", false},

		// Pattern with dots (regex special chars).
		{"dot in pattern", "*.ads.*", "https://tracker.ads.example.com/pixel", true},
		{"dot literal", "*.ads.*", "https://tracker-ads-example.com/pixel", false},

		// Pattern with query strings.
		{"query match", "*/search?q=*", "https://example.com/search?q=hello", true},
		{"literal question mark", "*/api?key=*", "https://example.com/api?key=123", true},

		// Edge cases.
		{"url empty pattern nonempty", "https://example.com", "", false},
		{"pattern is star empty url", "*", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchURL(tt.pattern, tt.url)
			if got != tt.want {
				t.Errorf("matchURL(%q, %q) = %v, want %v", tt.pattern, tt.url, got, tt.want)
			}
		})
	}
}

func TestMatchesFilter(t *testing.T) {
	base := NetworkEntry{
		RequestID:    "1",
		URL:          "https://api.example.com/v1/users",
		Method:       "GET",
		Status:       200,
		ResourceType: "XHR",
	}

	tests := []struct {
		name   string
		entry  NetworkEntry
		filter NetworkFilter
		want   bool
	}{
		// Empty filter matches everything.
		{"empty filter", base, NetworkFilter{}, true},

		// URL pattern only.
		{"url match", base, NetworkFilter{URLPattern: "*/v1/users"}, true},
		{"url no match", base, NetworkFilter{URLPattern: "*/v2/*"}, false},

		// Method filter.
		{"method match", base, NetworkFilter{Methods: []string{"GET"}}, true},
		{"method case insensitive", base, NetworkFilter{Methods: []string{"get"}}, true},
		{"method no match", base, NetworkFilter{Methods: []string{"POST"}}, false},
		{"method multi", base, NetworkFilter{Methods: []string{"POST", "GET"}}, true},

		// ResourceType filter.
		{"type match", base, NetworkFilter{ResourceTypes: []string{"XHR"}}, true},
		{"type case insensitive", base, NetworkFilter{ResourceTypes: []string{"xhr"}}, true},
		{"type no match", base, NetworkFilter{ResourceTypes: []string{"Document"}}, false},
		{"type multi", base, NetworkFilter{ResourceTypes: []string{"Document", "XHR"}}, true},

		// Combined filters (all must match).
		{"url+method match", base, NetworkFilter{
			URLPattern: "*/users",
			Methods:    []string{"GET"},
		}, true},
		{"url+method partial", base, NetworkFilter{
			URLPattern: "*/users",
			Methods:    []string{"POST"},
		}, false},
		{"all three match", base, NetworkFilter{
			URLPattern:    "*api*",
			Methods:       []string{"GET"},
			ResourceTypes: []string{"XHR"},
		}, true},
		{"all three one fails", base, NetworkFilter{
			URLPattern:    "*api*",
			Methods:       []string{"GET"},
			ResourceTypes: []string{"Document"},
		}, false},

		// Entry with empty fields.
		{"empty method entry", NetworkEntry{URL: "https://x.com"}, NetworkFilter{Methods: []string{"GET"}}, false},
		{"empty type entry", NetworkEntry{URL: "https://x.com"}, NetworkFilter{ResourceTypes: []string{"XHR"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilter(tt.entry, tt.filter)
			if got != tt.want {
				t.Errorf("matchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
