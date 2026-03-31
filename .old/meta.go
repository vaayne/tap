package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ArgDef struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type ScriptMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Domain      string            `json:"domain"`
	Args        map[string]ArgDef `json:"args"`
	ReadOnly    bool              `json:"readOnly"`
	Example     string            `json:"example"`
}

type Script struct {
	Meta ScriptMeta
	Body string // the async function body
	Raw  string // full file content
	Path string // file path
}

func ParseScript(content string) (*Script, error) {
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

func parseMeta(content string) (*ScriptMeta, error) {
	start := strings.Index(content, "/* @meta")
	if start == -1 {
		return nil, fmt.Errorf("no @meta block found")
	}

	end := strings.Index(content[start:], "*/")
	if end == -1 {
		return nil, fmt.Errorf("unclosed @meta block")
	}

	block := content[start : start+end+2]
	// Extract JSON between @meta and */
	jsonStart := strings.Index(block, "{")
	jsonEnd := strings.LastIndex(block, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON found in @meta block")
	}

	jsonStr := block[jsonStart : jsonEnd+1]

	var meta ScriptMeta
	if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return nil, fmt.Errorf("unmarshal meta: %w", err)
	}

	return &meta, nil
}

func parseBody(content string) (string, error) {
	// Find the end of the @meta comment block
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
