// Package fetch provides URL content extraction using go-defuddle.
// It fetches a web page via HTTP (with optional browser fallback) and
// extracts clean content (HTML or Markdown).
package fetch

import (
	"context"
	"fmt"
	"log"

	defuddle "github.com/vaayne/go-defuddle"
	"github.com/vaayne/tap/transport"
)

// Result holds the extracted content from a web page.
type Result struct {
	// Title is the page title.
	Title string `json:"title"`
	// Description is the meta description.
	Description string `json:"description"`
	// Domain is the hostname.
	Domain string `json:"domain"`
	// Author is the author name.
	Author string `json:"author"`
	// Published is the publish date.
	Published string `json:"published"`
	// Content is the extracted main content as clean HTML.
	Content string `json:"content"`
	// Markdown is the content converted to Markdown.
	Markdown string `json:"markdown,omitempty"`
	// WordCount is the word count of extracted content.
	WordCount int `json:"wordCount"`
}

// Fetcher extracts clean content from web pages.
type Fetcher struct {
	parser    *defuddle.Parser
	transport *transport.Transport
}

// New creates a new Fetcher backed by the given transport. Call Close() when done.
func New(tp *transport.Transport) (*Fetcher, error) {
	parser, err := defuddle.NewParser()
	if err != nil {
		return nil, fmt.Errorf("new parser: %w", err)
	}
	return &Fetcher{
		parser:    parser,
		transport: tp,
	}, nil
}

// Close releases resources.
func (f *Fetcher) Close() {
	if f.parser != nil {
		f.parser.Close()
	}
}

// Options controls fetch behavior.
type Options struct {
	// Markdown converts extracted HTML to Markdown.
	Markdown bool
	// UseBrowser forces browser-based fetching (level 2).
	UseBrowser bool
	// PauseFunc runs after browser navigation before HTML extraction.
	PauseFunc transport.PauseFunc
}

// Fetch retrieves a URL and extracts clean content.
// It tries HTTP first, falling back to browser if the result is poor.
func (f *Fetcher) Fetch(ctx context.Context, url string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{Markdown: true}
	}

	defOpts := &defuddle.Options{
		Markdown: opts.Markdown,
	}

	// If browser is forced, skip HTTP.
	if opts.UseBrowser {
		return f.fetchViaBrowser(ctx, url, opts, defOpts)
	}

	// Level 1: try direct HTTP.
	html, err := f.transport.GetHTML(ctx, url)
	if err == nil {
		result, parseErr := f.parse(html, url, defOpts)
		if parseErr == nil && hasContent(result) {
			return result, nil
		}
		if parseErr != nil {
			log.Printf("http fetch parse failed: %v, trying browser", parseErr)
		} else {
			log.Printf("http fetch returned poor content, trying browser")
		}
	} else {
		log.Printf("http fetch failed: %v, trying browser", err)
	}

	// Level 2: fallback to browser.
	return f.fetchViaBrowser(ctx, url, opts, defOpts)
}

func (f *Fetcher) fetchViaBrowser(ctx context.Context, url string, opts *Options, defOpts *defuddle.Options) (*Result, error) {
	html, err := f.transport.BrowseHTMLWithPause(ctx, url, opts.PauseFunc)
	if err != nil {
		return nil, fmt.Errorf("browser fetch: %w", err)
	}
	return f.parse(html, url, defOpts)
}

// ParseHTML extracts clean content from raw HTML without fetching.
func (f *Fetcher) ParseHTML(html, url string, opts *defuddle.Options) (*Result, error) {
	return f.parse(html, url, opts)
}

func (f *Fetcher) parse(html, url string, opts *defuddle.Options) (*Result, error) {
	dr, err := f.parser.Parse(html, url, opts)
	if err != nil {
		return nil, fmt.Errorf("defuddle parse: %w", err)
	}

	return &Result{
		Title:       dr.Title,
		Description: dr.Description,
		Domain:      dr.Domain,
		Author:      dr.Author,
		Published:   dr.Published,
		Content:     dr.Content,
		Markdown:    dr.Markdown,
		WordCount:   dr.WordCount,
	}, nil
}

// hasContent checks if a result has meaningful extracted content.
func hasContent(r *Result) bool {
	return r.Content != "" || r.Markdown != ""
}
