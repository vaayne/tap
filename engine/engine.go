// Package engine provides execution engines for running site scripts.
// It supports a fallback chain: QuickJS (fast, no browser) → CDP Browser (full browser).
package engine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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
// If all engines fail, returns an error that names each attempted engine.
func RunScript(ctx context.Context, engines []Engine, s *script.Script, args map[string]string, opts RunOpts) (any, error) {
	failures := make([]engineFailure, 0, len(engines))
	for _, e := range engines {
		result, err := e.Run(ctx, s, args, opts)
		if err != nil {
			failures = append(failures, engineFailure{name: e.Name(), err: err})
			log.Printf("%s failed: %v", e.Name(), err)
			continue
		}

		if errMsg := extractError(result); errMsg != "" {
			err = fmt.Errorf("returned error response: %s", errMsg)
			failures = append(failures, engineFailure{name: e.Name(), err: err})
			log.Printf("%s failed: %v, trying next engine", e.Name(), err)
			continue
		}

		return result, nil
	}
	return nil, engineFailuresError(failures)
}

// engineFailure records one attempted engine and why it failed.
type engineFailure struct {
	name string
	err  error
}

func engineFailuresError(failures []engineFailure) error {
	if len(failures) == 0 {
		return errors.New("no engines configured")
	}
	if len(failures) == 1 {
		return fmt.Errorf("%s failed: %w", strings.ToLower(failures[0].name), failures[0].err)
	}
	return allEnginesFailedError{failures: failures}
}

type allEnginesFailedError struct {
	failures []engineFailure
}

func (e allEnginesFailedError) Error() string {
	var b strings.Builder
	b.WriteString("all engines failed:")
	for _, failure := range e.failures {
		_, _ = fmt.Fprintf(&b, "\n  %s: %v", strings.ToLower(failure.name), failure.err)
	}
	return b.String()
}

func (e allEnginesFailedError) Unwrap() []error {
	causes := make([]error, 0, len(e.failures))
	for _, failure := range e.failures {
		causes = append(causes, failure.err)
	}
	return causes
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
