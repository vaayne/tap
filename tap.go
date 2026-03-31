// Package tap provides a unified API for interacting with web pages.
//
// Tap can run site scripts (with QuickJS → Browser fallback) and fetch
// clean content from URLs via go-defuddle. Both share a common transport
// layer for HTTP and browser-based network access.
//
// Basic usage:
//
//	client, err := tap.New(tap.WithSitesDir("./sites"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Run a site script
//	result, err := client.RunScript(ctx, "v2ex/hot", nil)
//
//	// Fetch clean content
//	content, err := client.Fetch(ctx, "https://example.com", nil)
package tap

import (
	"context"
	"fmt"
	"strings"

	"github.com/vaayne/tap/engine"
	"github.com/vaayne/tap/fetch"
	"github.com/vaayne/tap/script"
	"github.com/vaayne/tap/transport"
)

// Client is the main entry point for the tap library.
type Client struct {
	registry  *script.Registry
	engines   []engine.Engine
	fetcher   *fetch.Fetcher
	transport *transport.Transport
	opts      options
}

// New creates a new Client with the given options.
func New(optFns ...Option) (*Client, error) {
	opts := defaultOptions()
	for _, fn := range optFns {
		fn(&opts)
	}

	var reg *script.Registry
	if opts.sitesDir != "" {
		var err error
		reg, err = script.NewRegistry(opts.sitesDir)
		if err != nil {
			return nil, fmt.Errorf("load scripts: %w", err)
		}
	}

	tp := transport.New(transport.Config{
		WSURL:      opts.wsURL,
		ProfileDir: opts.profileDir,
		Headless:   opts.headless,
	})

	fetcher, err := fetch.New(tp)
	if err != nil {
		return nil, fmt.Errorf("new fetcher: %w", err)
	}

	var engines []engine.Engine
	if opts.forceBrowser {
		engines = []engine.Engine{
			engine.NewBrowser(tp),
		}
	} else {
		engines = []engine.Engine{
			engine.NewQuickJS(tp),
			engine.NewBrowser(tp),
		}
	}

	return &Client{
		registry:  reg,
		engines:   engines,
		fetcher:   fetcher,
		transport: tp,
		opts:      opts,
	}, nil
}

// Close releases all resources.
func (c *Client) Close() error {
	if c.fetcher != nil {
		c.fetcher.Close()
	}
	for _, e := range c.engines {
		_ = e.Close()
	}
	if c.transport != nil {
		_ = c.transport.Close()
	}
	return nil
}

// RunScript executes a site script by name with the given arguments.
// It tries QuickJS first, then falls back to the browser (unless --browser is set).
func (c *Client) RunScript(ctx context.Context, name string, args map[string]string) (any, error) {
	if c.registry == nil {
		return nil, fmt.Errorf("no sites directory configured")
	}

	s, ok := c.registry.Get(name)
	if !ok {
		return nil, &ScriptNotFoundError{Name: name, Available: c.scriptNames()}
	}

	if args == nil {
		args = make(map[string]string)
	}

	// Validate required args
	for argName, def := range s.Meta.Args {
		if def.Required {
			if _, ok := args[argName]; !ok {
				return nil, fmt.Errorf("missing required arg: %s (%s)", argName, def.Description)
			}
		}
	}

	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	return engine.RunScript(ctx, c.engines, s, args)
}

// Fetch retrieves a URL and extracts clean content using go-defuddle.
func (c *Client) Fetch(ctx context.Context, url string, opts *fetch.Options) (*fetch.Result, error) {
	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}
	return c.fetcher.Fetch(ctx, url, opts)
}

// ListScripts returns all available scripts sorted by name.
func (c *Client) ListScripts() []*script.Script {
	if c.registry == nil {
		return nil
	}
	return c.registry.List()
}

// GetScript returns a script by name.
func (c *Client) GetScript(name string) (*script.Script, bool) {
	if c.registry == nil {
		return nil, false
	}
	return c.registry.Get(name)
}

// scriptNames returns all registered script names for error suggestions.
func (c *Client) scriptNames() []string {
	scripts := c.ListScripts()
	names := make([]string, len(scripts))
	for i, s := range scripts {
		names[i] = s.Meta.Name
	}
	return names
}

// ScriptNotFoundError is returned when a script name doesn't match any registered script.
type ScriptNotFoundError struct {
	Name      string
	Available []string
}

func (e *ScriptNotFoundError) Error() string {
	return fmt.Sprintf("script not found: %s", e.Name)
}

// Suggestions returns script names similar to the requested name, ranked by relevance.
func (e *ScriptNotFoundError) Suggestions(max int) []string {
	type scored struct {
		name  string
		score int
	}
	var candidates []scored
	for _, name := range e.Available {
		if s := matchScore(e.Name, name); s > 0 {
			candidates = append(candidates, scored{name, s})
		}
	}
	// Sort by score descending
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

// matchScore returns a relevance score (0 = no match, higher = better).
func matchScore(query, target string) int {
	score := 0

	// Exact substring match is strong
	if strings.Contains(target, query) {
		score += 10
	}
	if strings.Contains(query, target) {
		score += 8
	}

	// Same site prefix is strong (e.g., "twiter/search" → "twitter/search")
	qSite, qAction := splitSlash(query)
	tSite, tAction := splitSlash(target)

	if qSite != "" && tSite != "" {
		if qSite == tSite {
			score += 20
		} else if editDistance(qSite, tSite) <= 2 {
			score += 15 // close typo in site name
		}
	}

	// Action part match
	if qAction != "" && tAction != "" {
		if qAction == tAction {
			score += 10
		} else if strings.Contains(tAction, qAction) || strings.Contains(qAction, tAction) {
			score += 5
		}
	}

	// Low edit distance on full name
	d := editDistance(query, target)
	if d <= 3 {
		score += (4 - d) * 3
	}

	return score
}

func splitSlash(s string) (string, string) {
	if i := strings.IndexByte(s, '/'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

// editDistance computes the Levenshtein distance between two strings.
func editDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[lb]
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
