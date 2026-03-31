package engine

import (
	"context"
	"fmt"
	"testing"

	"github.com/vaayne/tap/script"
)

// mockEngine is a test engine that returns a fixed result or error.
type mockEngine struct {
	name   string
	result any
	err    error
}

func (m *mockEngine) Name() string { return m.name }
func (m *mockEngine) Close() error { return nil }
func (m *mockEngine) Run(_ context.Context, _ *script.Script, _ map[string]string) (any, error) {
	return m.result, m.err
}

func TestRunScript_FirstEngineSucceeds(t *testing.T) {
	engines := []Engine{
		&mockEngine{name: "fast", result: map[string]any{"ok": true}},
		&mockEngine{name: "slow", result: map[string]any{"ok": true}},
	}

	result, err := RunScript(context.Background(), engines, &script.Script{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]any)
	if m["ok"] != true {
		t.Errorf("result = %v, want ok=true", result)
	}
}

func TestRunScript_FallbackToSecond(t *testing.T) {
	engines := []Engine{
		&mockEngine{name: "fast", err: fmt.Errorf("unsupported")},
		&mockEngine{name: "slow", result: map[string]any{"fallback": true}},
	}

	result, err := RunScript(context.Background(), engines, &script.Script{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]any)
	if m["fallback"] != true {
		t.Errorf("result = %v, want fallback=true", result)
	}
}

func TestRunScript_AllFail(t *testing.T) {
	engines := []Engine{
		&mockEngine{name: "e1", err: fmt.Errorf("fail1")},
		&mockEngine{name: "e2", err: fmt.Errorf("fail2")},
	}

	_, err := RunScript(context.Background(), engines, &script.Script{}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
