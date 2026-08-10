// Package fetch extracts clean content from the active agent-browser page.
package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vaayne/tap/agentbrowser"
	"github.com/vaayne/tap/assets"
)

// Browser is the agent-browser surface needed for extraction.
type Browser interface {
	Open(context.Context, string) error
	Eval(context.Context, string) (any, error)
	CurrentURL(context.Context) (string, error)
	HasActiveSession(context.Context) (bool, error)
}

// Result holds content extracted by Defuddle in the browser page.
type Result struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Author      string `json:"author"`
	Published   string `json:"published"`
	Content     string `json:"content"`
	Markdown    string `json:"markdown,omitempty"`
	WordCount   int    `json:"wordCount"`
}

// Options controls fetch output.
type Options struct {
	Markdown bool
}

// Fetcher extracts the active browser document. It does not own or close the
// browser session.
type Fetcher struct {
	browser Browser
}

func New(browser Browser) *Fetcher {
	return &Fetcher{browser: browser}
}

// Fetch navigates when url is non-empty. With an empty URL it reads the current
// tab and refuses to create a browser session implicitly.
func (f *Fetcher) Fetch(ctx context.Context, url string, opts *Options) (*Result, error) {
	if url == "" {
		active, err := f.browser.HasActiveSession(ctx)
		if err != nil {
			return nil, fmt.Errorf("check current agent-browser session: %w", err)
		}
		if !active {
			return nil, fmt.Errorf("no active agent-browser session; open a page first")
		}
		current, err := f.browser.CurrentURL(ctx)
		if err != nil {
			return nil, fmt.Errorf("get current agent-browser tab: %w", err)
		}
		if current == "" || current == "about:blank" {
			return nil, fmt.Errorf("current agent-browser tab has no readable page")
		}
	} else if err := f.browser.Open(ctx, url); err != nil {
		return nil, fmt.Errorf("open page: %w", err)
	}

	value, err := f.browser.Eval(ctx, extractionScript())
	if err != nil {
		return nil, fmt.Errorf("extract page with Defuddle: %w", err)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal Defuddle result: %w", err)
	}
	var raw struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		Domain          string `json:"domain"`
		Author          string `json:"author"`
		Published       string `json:"published"`
		Content         string `json:"content"`
		ContentMarkdown string `json:"contentMarkdown"`
		WordCount       int    `json:"wordCount"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode Defuddle result: %w", err)
	}
	result := &Result{
		Title:       raw.Title,
		Description: raw.Description,
		Domain:      raw.Domain,
		Author:      raw.Author,
		Published:   raw.Published,
		Content:     raw.Content,
		Markdown:    raw.ContentMarkdown,
		WordCount:   raw.WordCount,
	}
	if opts != nil && !opts.Markdown {
		result.Markdown = ""
	}
	return result, nil
}

func extractionScript() string {
	var script strings.Builder
	script.Grow(len(assets.DefuddleBrowser) + 256)
	script.WriteString(assets.DefuddleBrowser)
	script.WriteString(`
;(async () => {
  const result = await new globalThis.Defuddle(document, {
    separateMarkdown: true
  }).parseAsync();
  return result;
})()
`)
	return script.String()
}

var _ Browser = (*agentbrowser.Client)(nil)
