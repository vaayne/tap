// Package engine provides execution engines for running site scripts.
// It supports a fallback chain: QuickJS (fast, no browser) → CDP Browser (full browser).
package engine

import (
	"context"
	"fmt"
	"log"

	"github.com/vaayne/tap/script"
)

// Engine can execute a site script with arguments and return a JSON-compatible result.
type Engine interface {
	// Run executes a script with the given arguments.
	Run(ctx context.Context, s *script.Script, args map[string]string) (any, error)

	// Name returns the engine name for logging.
	Name() string

	// Close releases resources held by the engine.
	Close() error
}

// RunScript tries each engine in order, returning the first successful result.
// If all engines fail, returns the last error.
func RunScript(ctx context.Context, engines []Engine, s *script.Script, args map[string]string) (any, error) {
	var lastErr error
	for _, e := range engines {
		result, err := e.Run(ctx, s, args)
		if err == nil {
			return result, nil
		}
		lastErr = err
		log.Printf("%s failed: %v", e.Name(), err)
	}
	return nil, fmt.Errorf("all engines failed: %w", lastErr)
}
