package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SnapshotOptions controls browser snapshot generation.
type SnapshotOptions struct {
	InteractiveOnly bool
	Selector        string
	Depth           int64
	Mode            string // auto | ax
}

// SnapshotResult is an AI-friendly semantic snapshot of the current page.
type SnapshotResult struct {
	URL         string         `json:"url"`
	Title       string         `json:"title,omitempty"`
	Mode        string         `json:"mode"`
	GeneratedAt time.Time      `json:"generatedAt"`
	DocumentKey string         `json:"documentKey"`
	Nodes       []SnapshotNode `json:"nodes"`
	Refs        []SnapshotRef  `json:"refs"`
}

// SnapshotNode is one node in the compact semantic tree.
type SnapshotNode struct {
	Ref         string   `json:"ref,omitempty"`
	Role        string   `json:"role,omitempty"`
	Name        string   `json:"name,omitempty"`
	Value       string   `json:"value,omitempty"`
	Description string   `json:"description,omitempty"`
	States      []string `json:"states,omitempty"`
	Children    []int    `json:"children,omitempty"`
}

// SnapshotRef maps an element ref to backend IDs for follow-up actions.
type SnapshotRef struct {
	Ref              string `json:"ref"`
	BackendDOMNodeID int64  `json:"backendDomNodeId,omitempty"`
	AXNodeID         string `json:"axNodeId,omitempty"`
	FrameID          string `json:"frameId,omitempty"`
	SelectorHint     string `json:"selectorHint,omitempty"`
	Role             string `json:"role,omitempty"`
	Name             string `json:"name,omitempty"`
}

func isElementRef(arg string) bool {
	return strings.HasPrefix(strings.TrimSpace(arg), "@e")
}

func (m *Manager) saveSnapshot(sessionName, tabName string, result *SnapshotResult) error {
	if result == nil {
		return fmt.Errorf("snapshot is required")
	}
	path := m.snapshotPath(sessionName, tabName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot dir: %w", err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encode snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	return nil
}

func (m *Manager) loadSnapshot(sessionName, tabName string) (*SnapshotResult, error) {
	path := m.snapshotPath(sessionName, tabName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snapshot: %w", err)
	}
	var out SnapshotResult
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &out, nil
}

func (m *Manager) snapshotPath(sessionName, tabName string) string {
	return filepath.Join(m.store.Root(), "snapshots", sessionName, tabName+".json")
}
