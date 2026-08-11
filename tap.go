// Package tap provides reusable site programs and content extraction powered
// by agent-browser.
package tap

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vaayne/tap/agentbrowser"
	"github.com/vaayne/tap/fetch"
	"github.com/vaayne/tap/script"
	"github.com/vaayne/tap/sites"
)

// Client is the main entry point for the tap library.
type Client struct {
	registry *script.Registry
	browser  *agentbrowser.Client
	fetcher  *fetch.Fetcher
	opts     options
}

// New creates a Client. It performs no browser startup or installation work.
func New(_ context.Context, optFns ...Option) (*Client, error) {
	opts := defaultOptions()
	for _, fn := range optFns {
		fn(&opts)
	}

	var registry *script.Registry
	if opts.sitesDir != "" {
		var err error
		registry, err = DefaultRegistry(opts.sitesDir, opts.localOverrideDir)
		if err != nil {
			return nil, fmt.Errorf("load scripts: %w", err)
		}
	}
	browser := agentbrowser.New(opts.agentBrowser)
	return &Client{
		registry: registry,
		browser:  browser,
		fetcher:  fetch.New(browser),
		opts:     opts,
	}, nil
}

// DefaultRegistry creates the standard registry with cache, built-in, and
// override sources in increasing priority order.
func DefaultRegistry(cacheDir, overrideDir string) (*script.Registry, error) {
	return script.NewRegistry(
		script.Source{Path: cacheDir, Type: script.ScriptSourceCache},
		script.Source{FS: sites.FS, Type: script.ScriptSourceBuiltin},
		script.Source{Path: overrideDir, Type: script.ScriptSourceOverride},
	)
}

// Close is a no-op. Tap never owns or closes agent-browser sessions.
func (c *Client) Close() error { return nil }

// RunScript opens the script's domain in the active agent-browser session and
// evaluates the site program there.
func (c *Client) RunScript(ctx context.Context, name string, args map[string]string) (any, error) {
	if c.registry == nil {
		return nil, fmt.Errorf("no sites directory configured")
	}
	s, ok := c.registry.Get(name)
	if !ok {
		return nil, &ScriptNotFoundError{Name: name, Available: c.scriptNames()}
	}
	if s.Source == script.ScriptSourceOverride {
		fmt.Fprintf(os.Stderr, "Using local script: %s\n", name)
	}
	if args == nil {
		args = make(map[string]string)
	}
	for argName, def := range s.Meta.Args {
		if def.Required {
			if _, exists := args[argName]; !exists {
				return nil, fmt.Errorf("missing required arg: %s (%s)", argName, def.Description)
			}
		}
	}
	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	navigationURL := s.Meta.ExecutionURL()
	headers := s.Meta.ResolveHeaders()
	program, err := siteProgram(s, args, headers)
	if err != nil {
		return nil, err
	}
	return c.browser.OpenAndEval(ctx, navigationURL, program, headers)
}

// Fetch extracts a URL through agent-browser. An empty URL reads the active tab
// without navigating or creating a browser session.
func (c *Client) Fetch(ctx context.Context, url string, opts *fetch.Options) (*fetch.Result, error) {
	if opts == nil {
		opts = &fetch.Options{Markdown: true}
	}
	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}
	return c.fetcher.Fetch(ctx, url, opts)
}

func siteProgram(s *script.Script, args, headers map[string]string) (string, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("marshal script args: %w", err)
	}
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return "", fmt.Errorf("marshal script headers: %w", err)
	}
	originJSON, err := json.Marshal(s.Meta.Origin())
	if err != nil {
		return "", fmt.Errorf("marshal script origin: %w", err)
	}
	return fmt.Sprintf(`(async () => {
  const __tapArgs = %s;
  const __tapHeaders = %s;
  const __tapOrigin = %s;
  if (location.origin !== __tapOrigin) {
    throw new Error("Tap execution origin mismatch: expected " + __tapOrigin + ", got " + location.origin);
  }
  const __tapNativeFetch = globalThis.fetch.bind(globalThis);
  const fetch = (input, init = {}) => {
    const url = new URL(input instanceof Request ? input.url : String(input), location.href);
    if (url.origin !== __tapOrigin) {
      throw new Error("Tap cross-origin fetch blocked: " + url.origin + " (declared origin: " + __tapOrigin + ")");
    }
    const headers = new Headers(input instanceof Request ? input.headers : undefined);
    new Headers(init.headers || {}).forEach((value, name) => headers.set(name, value));
    for (const [name, value] of Object.entries(__tapHeaders)) headers.set(name, value);
    return __tapNativeFetch(input, {...init, headers});
  };
  return await (%s)(__tapArgs);
})()`, argsJSON, headersJSON, originJSON, s.Body), nil
}

// ListScripts returns all available scripts sorted by name.
func (c *Client) ListScripts() []*script.Script {
	if c.registry == nil {
		return nil
	}
	return c.registry.List()
}

// ListScriptsOverrides returns scripts loaded from the local override directory.
func (c *Client) ListScriptsOverrides() []*script.Script {
	if c.registry == nil {
		return nil
	}
	return c.registry.ListOverrides()
}

// GetScript returns a script by its path-derived name.
func (c *Client) GetScript(name string) (*script.Script, bool) {
	if c.registry == nil {
		return nil, false
	}
	return c.registry.Get(name)
}

func (c *Client) scriptNames() []string {
	scripts := c.ListScripts()
	names := make([]string, len(scripts))
	for i, s := range scripts {
		names[i] = s.Meta.Name
	}
	return names
}

// ScriptNotFoundError is returned when a script name is absent.
type ScriptNotFoundError struct {
	Name      string
	Available []string
}

func (e *ScriptNotFoundError) Error() string {
	return fmt.Sprintf("script not found: %s", e.Name)
}

// Suggestions returns similar script names ranked by relevance.
func (e *ScriptNotFoundError) Suggestions(max int) []string {
	type scored struct {
		name  string
		score int
	}
	var candidates []scored
	for _, name := range e.Available {
		if score := matchScore(e.Name, name); score > 0 {
			candidates = append(candidates, scored{name, score})
		}
	}
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score > candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	result := make([]string, 0, max)
	for i := 0; i < len(candidates) && i < max; i++ {
		result = append(result, candidates[i].name)
	}
	return result
}

func matchScore(query, target string) int {
	score := 0
	if strings.Contains(target, query) {
		score += 10
	}
	if strings.Contains(query, target) {
		score += 8
	}
	qSite, qAction := splitSlash(query)
	tSite, tAction := splitSlash(target)
	if qSite != "" && tSite != "" {
		if qSite == tSite {
			score += 20
		} else if editDistance(qSite, tSite) <= 2 {
			score += 15
		}
	}
	if qAction != "" && tAction != "" {
		if qAction == tAction {
			score += 10
		} else if strings.Contains(tAction, qAction) || strings.Contains(qAction, tAction) {
			score += 5
		}
	}
	if distance := editDistance(query, target); distance <= 3 {
		score += (4 - distance) * 3
	}
	return score
}

func splitSlash(value string) (string, string) {
	if i := strings.IndexByte(value, '/'); i >= 0 {
		return value[:i], value[i+1:]
	}
	return value, ""
}

func editDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = min3(current[j-1]+1, previous[j]+1, previous[j-1]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
