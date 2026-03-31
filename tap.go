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
	})

	fetcher, err := fetch.New(tp)
	if err != nil {
		return nil, fmt.Errorf("new fetcher: %w", err)
	}

	engines := []engine.Engine{
		engine.NewQuickJS(tp),
		engine.NewBrowser(tp),
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
		e.Close()
	}
	if c.transport != nil {
		c.transport.Close()
	}
	return nil
}

// RunScript executes a site script by name with the given arguments.
// It tries QuickJS first, then falls back to the browser.
func (c *Client) RunScript(ctx context.Context, name string, args map[string]string) (any, error) {
	if c.registry == nil {
		return nil, fmt.Errorf("no sites directory configured")
	}

	s, ok := c.registry.Get(name)
	if !ok {
		return nil, fmt.Errorf("script not found: %s", name)
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

	return engine.RunScript(ctx, c.engines, s, args)
}

// Fetch retrieves a URL and extracts clean content using go-defuddle.
func (c *Client) Fetch(ctx context.Context, url string, opts *fetch.Options) (*fetch.Result, error) {
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
