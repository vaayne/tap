// Package fetch provides URL content extraction using go-defuddle.
// It fetches a web page via HTTP and extracts clean content (HTML or Markdown).
package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"

	defuddle "github.com/vaayne/go-defuddle"
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
	parser *defuddle.Parser
	client *http.Client
}

// New creates a new Fetcher. Call Close() when done.
func New() (*Fetcher, error) {
	parser, err := defuddle.NewParser()
	if err != nil {
		return nil, fmt.Errorf("new parser: %w", err)
	}
	return &Fetcher{
		parser: parser,
		client: &http.Client{},
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
}

// Fetch retrieves a URL and extracts clean content.
func (f *Fetcher) Fetch(ctx context.Context, url string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{Markdown: true}
	}

	html, err := f.fetchHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch html: %w", err)
	}

	defOpts := &defuddle.Options{
		Markdown: opts.Markdown,
	}

	dr, err := f.parser.Parse(html, url, defOpts)
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

func (f *Fetcher) fetchHTML(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}
