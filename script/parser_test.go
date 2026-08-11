package script

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	content := `/* @meta
{
  "description": "A test script",
  "domain": "example.com",
  "args": {
    "query": {"required": true, "description": "Search query"},
    "count": {"required": false, "description": "Number of results"}
  },
  "readOnly": true,
  "example": "tap site test/hello query=foo"
}
*/

async function(args) {
  return {hello: args.query};
}`

	s, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if s.Meta.Name != "" {
		t.Errorf("name must be assigned by registry, got %q", s.Meta.Name)
	}
	if s.Meta.Description != "A test script" {
		t.Errorf("description = %q, want %q", s.Meta.Description, "A test script")
	}
	if s.Meta.Domain != "example.com" {
		t.Errorf("domain = %q, want %q", s.Meta.Domain, "example.com")
	}
	if !s.Meta.ReadOnly {
		t.Error("readOnly = false, want true")
	}
	if len(s.Meta.Args) != 2 {
		t.Fatalf("args count = %d, want 2", len(s.Meta.Args))
	}
	if !s.Meta.Args["query"].Required {
		t.Error("args[query].required = false, want true")
	}
	if s.Meta.Args["count"].Required {
		t.Error("args[count].required = true, want false")
	}
	if s.Body == "" {
		t.Error("body is empty")
	}
}

func TestParse_NoMeta(t *testing.T) {
	_, err := Parse("function() { return 1; }")
	if err == nil {
		t.Error("expected error for missing @meta")
	}
}

func TestParse_UnclosedMeta(t *testing.T) {
	_, err := Parse("/* @meta { \"name\": \"test\" }")
	if err == nil {
		t.Error("expected error for unclosed @meta")
	}
}

func TestParse_NoBody(t *testing.T) {
	_, err := Parse(`/* @meta
{
  "description": "empty",
  "domain": "example.com"
}
*/`)
	if err == nil {
		t.Error("expected error for missing body")
	}
}

func TestParse_RejectsInvalidDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain string
	}{
		{name: "missing"},
		{name: "scheme", domain: "https://example.com"},
		{name: "path", domain: "example.com/api"},
		{name: "port", domain: "example.com:8443"},
		{name: "uppercase", domain: "Example.com"},
		{name: "IP address", domain: "127.0.0.1"},
		{name: "single label", domain: "localhost"},
		{name: "leading hyphen", domain: "-api.example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := `/* @meta
{"description":"invalid domain","domain":"` + tt.domain + `","args":{}}
*/
async function(args) { return args; }`
			_, err := Parse(content)
			if err == nil || !strings.Contains(err.Error(), "domain") {
				t.Fatalf("Parse() error = %v, want domain validation error", err)
			}
		})
	}
}

func TestParse_ValidatesStartPath(t *testing.T) {
	valid := `/* @meta
{"description":"valid path","domain":"example.com","executionDomain":"api.example.com","startPath":"/api/bootstrap?format=json","args":{}}
*/
async function(args) { return args; }`
	script, err := Parse(valid)
	if err != nil {
		t.Fatal(err)
	}
	if got := script.Meta.ExecutionURL(); got != "https://api.example.com/api/bootstrap?format=json" {
		t.Fatalf("ExecutionURL() = %q", got)
	}

	for _, path := range []string{"api", "https://other.example/api", "//other.example/api", "/api#fragment"} {
		content := `/* @meta
{"description":"invalid path","domain":"example.com","startPath":"` + path + `","args":{}}
*/
async function(args) { return args; }`
		_, err := Parse(content)
		if err == nil || !strings.Contains(err.Error(), "startPath") {
			t.Fatalf("Parse(startPath=%q) error = %v", path, err)
		}
	}
}

func TestParse_RejectsInvalidExecutionDomain(t *testing.T) {
	content := `/* @meta
{"description":"invalid execution domain","domain":"example.com","executionDomain":"https://api.example.com","args":{}}
*/
async function(args) { return args; }`
	_, err := Parse(content)
	if err == nil || !strings.Contains(err.Error(), "executionDomain") {
		t.Fatalf("Parse() error = %v, want executionDomain validation error", err)
	}
}

func TestMeta_ResolveHeaders_AllSet(t *testing.T) {
	t.Setenv("API_KEY", "secret123")
	t.Setenv("USER_ID", "42")
	m := Meta{
		Headers: map[string]string{
			"X-API-Key": "${API_KEY}",
			"X-User-ID": "Bearer ${USER_ID}",
		},
	}
	resolved := m.ResolveHeaders()
	if len(resolved) != 2 {
		t.Fatalf("resolved count = %d, want 2", len(resolved))
	}
	if resolved["X-API-Key"] != "secret123" {
		t.Errorf("X-API-Key = %q, want %q", resolved["X-API-Key"], "secret123")
	}
	if resolved["X-User-ID"] != "Bearer 42" {
		t.Errorf("X-User-ID = %q, want %q", resolved["X-User-ID"], "Bearer 42")
	}
}

func TestMeta_ResolveHeaders_SkipMissing(t *testing.T) {
	// API_KEY is set, MISSING_VAR is not
	t.Setenv("API_KEY", "secret123")
	m := Meta{
		Headers: map[string]string{
			"X-API-Key": "${API_KEY}",
			"X-Missing": "${MISSING_VAR}",
			"X-Partial": "prefix-${MISSING_VAR}-suffix",
		},
	}
	resolved := m.ResolveHeaders()
	if len(resolved) != 1 {
		t.Fatalf("resolved count = %d, want 1", len(resolved))
	}
	if _, ok := resolved["X-API-Key"]; !ok {
		t.Error("X-API-Key should be present")
	}
	if _, ok := resolved["X-Missing"]; ok {
		t.Error("X-Missing should be skipped")
	}
	if _, ok := resolved["X-Partial"]; ok {
		t.Error("X-Partial should be skipped")
	}
}

func TestMeta_ResolveHeaders_EmptyHeaders(t *testing.T) {
	m := Meta{Headers: map[string]string{}}
	resolved := m.ResolveHeaders()
	if len(resolved) != 0 {
		t.Fatalf("resolved count = %d, want 0", len(resolved))
	}
}

func TestParse_EmptyArgs(t *testing.T) {
	content := `/* @meta
{
  "description": "no args",
  "domain": "example.com",
  "args": {}
}
*/

async function(args) { return {}; }`

	s, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(s.Meta.Args) != 0 {
		t.Errorf("args count = %d, want 0", len(s.Meta.Args))
	}
}
