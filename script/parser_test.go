package script

import (
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
  "description": "empty"
}
*/`)
	if err == nil {
		t.Error("expected error for missing body")
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
