// Package engine provides execution engines for running site scripts.
// It supports a fallback chain: QuickJS (fast, no browser) → CDP Browser (full browser).
package engine

import (
	"context"
	"fmt"
	"log"

	"github.com/vaayne/tap/script"
)

// RunOpts holds per-run configuration for script execution.
type RunOpts struct {
	Headers map[string]string // resolved meta headers (env vars interpolated)
}

// Engine can execute a site script with arguments and return a JSON-compatible result.
type Engine interface {
	// Run executes a script with the given arguments.
	Run(ctx context.Context, s *script.Script, args map[string]string, opts RunOpts) (any, error)

	// Name returns the engine name for logging.
	Name() string

	// Close releases resources held by the engine.
	Close() error
}

// RunScript tries each engine in order, returning the first successful result.
// If a result contains an "error" field (e.g. {"error":"HTTP 400"}), it is
// treated as a failure and the next engine is tried.
// If all engines fail, returns the last error.
func RunScript(ctx context.Context, engines []Engine, s *script.Script, args map[string]string, opts RunOpts) (any, error) {
	var lastErr error
	for _, e := range engines {
		result, err := e.Run(ctx, s, args, opts)
		if err != nil {
			lastErr = err
			log.Printf("%s failed: %v", e.Name(), err)
			continue
		}

		if errMsg := extractError(result); errMsg != "" {
			lastErr = fmt.Errorf("%s returned error response: %s", e.Name(), errMsg)
			log.Printf("%v, trying next engine", lastErr)
			continue
		}

		return result, nil
	}
	return nil, fmt.Errorf("all engines failed: %w", lastErr)
}

// extractError checks if a result looks like an error response.
// It detects JSON objects with an "error" key at the top level.
func extractError(result any) string {
	m, ok := result.(map[string]any)
	if !ok {
		return ""
	}
	errVal, ok := m["error"]
	if !ok {
		return ""
	}
	if s, ok := errVal.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", errVal)
}
