package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vaayne/tap/script"
	"github.com/vaayne/tap/transport"
)

// Browser executes scripts in a real Chrome browser via CDP.
type Browser struct {
	transport *transport.Transport
	pauseFn   transport.PauseFunc
}

// NewBrowser creates a new Browser engine backed by the given transport.
func NewBrowser(tp *transport.Transport, pauseFn transport.PauseFunc) *Browser {
	return &Browser{transport: tp, pauseFn: pauseFn}
}

func (b *Browser) Name() string { return "Browser" }
func (b *Browser) Close() error { return nil }

func (b *Browser) Run(ctx context.Context, s *script.Script, args map[string]string) (any, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	js := fmt.Sprintf("(%s)(%s)", s.Body, string(argsJSON))

	navURL := "about:blank"
	if s.Meta.Domain != "" {
		navURL = "https://" + s.Meta.Domain
	}

	return b.transport.BrowseEvalWithPause(ctx, navURL, js, b.pauseFn)
}
