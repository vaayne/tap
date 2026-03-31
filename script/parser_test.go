package script

import (
	"testing"
)

func TestParse(t *testing.T) {
	content := `/* @meta
{
  "name": "test/hello",
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

	if s.Meta.Name != "test/hello" {
		t.Errorf("name = %q, want %q", s.Meta.Name, "test/hello")
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
  "name": "test/empty",
  "description": "empty"
}
*/`)
	if err == nil {
		t.Error("expected error for missing body")
	}
}

func TestParse_EmptyArgs(t *testing.T) {
	content := `/* @meta
{
  "name": "test/noargs",
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
