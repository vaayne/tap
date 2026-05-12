// Package script handles parsing and discovery of site scripts.
package script

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ArgDef describes a single argument for a script.
type ArgDef struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// EnvDef describes a single environment variable dependency for a script.
type EnvDef struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// ScriptSource identifies where a script was loaded from.
type ScriptSource int

const (
	ScriptSourceCache    ScriptSource = iota // ~/.cache/tap/sites/
	ScriptSourceBuiltin                        // sites/ (embedded)
	ScriptSourceOverride                       // ~/.config/tap/sites/
)

// Meta holds the metadata extracted from a script's @meta block.
type Meta struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Domain       string            `json:"domain"`
	Args         map[string]ArgDef `json:"args"`
	ReadOnly     bool              `json:"readOnly"`
	Example      string            `json:"example"`
	Capabilities []string          `json:"capabilities"`
	Runtime      string            `json:"runtime"`
	AuthRequired bool              `json:"authRequired"`
	Headers      map[string]string `json:"headers"`
	Env          map[string]EnvDef `json:"env"`
}

// Script represents a parsed site script with metadata and function body.
type Script struct {
	Meta   Meta
	Body   string       // the async function body
	Raw    string       // full file content
	Path   string       // file path
	Source ScriptSource // where the script was loaded from
}

// Parse parses a script file content, extracting @meta JSON and the function body.
func Parse(content string) (*Script, error) {
	meta, err := parseMeta(content)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	body, err := parseBody(content)
	if err != nil {
		return nil, fmt.Errorf("parse body: %w", err)
	}

	return &Script{
		Meta: *meta,
		Body: body,
		Raw:  content,
	}, nil
}

func parseMeta(content string) (*Meta, error) {
	start := strings.Index(content, "/* @meta")
	if start == -1 {
		return nil, fmt.Errorf("no @meta block found")
	}

	end := strings.Index(content[start:], "*/")
	if end == -1 {
		return nil, fmt.Errorf("unclosed @meta block")
	}

	block := content[start : start+end+2]
	jsonStart := strings.Index(block, "{")
	jsonEnd := strings.LastIndex(block, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON found in @meta block")
	}

	jsonStr := block[jsonStart : jsonEnd+1]

	var meta Meta
	if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return nil, fmt.Errorf("unmarshal meta: %w", err)
	}

	return &meta, nil
}

// ValidateEnv checks that all required environment variables are set.
func (m *Meta) ValidateEnv() error {
	var missing []string
	for name, def := range m.Env {
		if def.Required {
			if _, ok := os.LookupEnv(name); !ok {
				missing = append(missing, fmt.Sprintf("%s (%s)", name, def.Description))
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

// ResolveHeaders copies Headers and interpolates ${ENV_VAR} values via os.Getenv.
// Headers referencing unset environment variables are skipped entirely.
func (m *Meta) ResolveHeaders() map[string]string {
	result := make(map[string]string, len(m.Headers))
	for k, v := range m.Headers {
		expanded, ok := expandEnv(v)
		if !ok {
			continue
		}
		result[k] = expanded
	}
	return result
}

func expandEnv(s string) (string, bool) {
	missing := false
	expanded := os.Expand(s, func(key string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			missing = true
			return ""
		}
		return val
	})
	return expanded, !missing
}

func parseBody(content string) (string, error) {
	end := strings.Index(content, "*/")
	if end == -1 {
		return "", fmt.Errorf("no @meta block end found")
	}

	body := strings.TrimSpace(content[end+2:])
	if body == "" {
		return "", fmt.Errorf("no function body found")
	}

	return body, nil
}
