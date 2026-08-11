// Package script handles parsing and discovery of site scripts.
package script

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

var envReferencePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// ArgDef describes a single argument for a script.
type ArgDef struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// ScriptSource identifies where a script was loaded from.
type ScriptSource int

const (
	ScriptSourceCache    ScriptSource = iota // ~/.cache/tap/sites/
	ScriptSourceBuiltin                      // sites/ (embedded)
	ScriptSourceOverride                     // ~/.config/tap/sites/
)

// Meta holds the metadata extracted from a script's @meta block.
type Meta struct {
	// Name is derived from the script's relative path by Registry. It is not part
	// of the 1.0 metadata format.
	Name string `json:"-"`
	// Description is a short human-readable summary shown in `tap site list`.
	Description string `json:"description"`
	// Domain is the script's catalog host and the default HTTPS execution host.
	Domain string `json:"domain"`
	// ExecutionDomain is a Tap-specific execution host override used by source
	// compatibility adapters. Site fetches are restricted to this exact origin.
	ExecutionDomain string `json:"executionDomain"`
	// StartPath is an optional same-origin navigation target. It is useful when
	// a domain root redirects away from the execution origin.
	StartPath string `json:"startPath"`
	// Args declares the named arguments the script accepts. Each key maps to an
	// ArgDef describing whether it is required and what it represents.
	Args map[string]ArgDef `json:"args"`
	// ReadOnly marks scripts that only read data and never mutate state.
	ReadOnly bool `json:"readOnly"`
	// Capabilities is reserved for future capability declarations.
	Capabilities []string `json:"capabilities"`
	// AuthRequired indicates the script needs browser-based authentication.
	AuthRequired bool `json:"authRequired"`
	// Headers are HTTP headers injected into every fetch() call made by the script.
	// Values may reference environment variables with ${VAR} syntax; headers whose
	// variable is unset are omitted entirely. See ResolveHeaders.
	Headers map[string]string `json:"headers"`
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
	if err := meta.validate(); err != nil {
		return nil, fmt.Errorf("validate meta: %w", err)
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

func (m *Meta) validate() error {
	if m.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if err := validateDomain("domain", m.Domain); err != nil {
		return err
	}
	if m.ExecutionDomain != "" {
		if err := validateDomain("executionDomain", m.ExecutionDomain); err != nil {
			return err
		}
	}
	domain := m.effectiveExecutionDomain()
	if m.StartPath != "" {
		start, err := url.Parse(m.StartPath)
		if err != nil || !strings.HasPrefix(m.StartPath, "/") || start.IsAbs() || start.Host != "" || start.Fragment != "" {
			return fmt.Errorf("startPath must be an absolute path on execution domain %q: %q", domain, m.StartPath)
		}
	}
	return nil
}

func validateDomain(field, domain string) error {
	if domain != strings.ToLower(domain) || strings.TrimSpace(domain) != domain {
		return fmt.Errorf("%s must be a lowercase hostname: %q", field, domain)
	}
	if net.ParseIP(domain) != nil {
		return fmt.Errorf("%s must be a hostname, not an IP address: %q", field, domain)
	}
	if len(domain) > 253 {
		return fmt.Errorf("%s exceeds 253 characters", field)
	}
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return fmt.Errorf("%s must be a fully qualified hostname: %q", field, domain)
	}
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return fmt.Errorf("invalid %s label in %q", field, domain)
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return fmt.Errorf("invalid character in %s %q", field, domain)
			}
		}
	}
	return nil
}

func (m *Meta) effectiveExecutionDomain() string {
	if m.ExecutionDomain != "" {
		return m.ExecutionDomain
	}
	return m.Domain
}

// Origin returns the exact origin available to site fetches.
func (m *Meta) Origin() string {
	return "https://" + m.effectiveExecutionDomain()
}

// ExecutionURL returns the same-origin page Tap opens before evaluation.
func (m *Meta) ExecutionURL() string {
	path := m.StartPath
	if path == "" {
		path = "/"
	}
	return m.Origin() + path
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

// EnvNames returns environment variables inferred from header templates.
func (m *Meta) EnvNames() []string {
	names := make(map[string]struct{})
	for _, value := range m.Headers {
		for _, match := range envReferencePattern.FindAllStringSubmatch(value, -1) {
			names[match[1]] = struct{}{}
		}
	}
	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)
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
