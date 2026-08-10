package fetch

import (
	"context"
	"strings"
	"testing"
)

type fakeBrowser struct {
	active  bool
	url     string
	opened  string
	script  string
	result  any
	openErr error
}

func (f *fakeBrowser) Open(_ context.Context, url string) error {
	f.opened = url
	return f.openErr
}
func (f *fakeBrowser) Eval(_ context.Context, script string) (any, error) {
	f.script = script
	return f.result, nil
}
func (f *fakeBrowser) CurrentURL(context.Context) (string, error) { return f.url, nil }
func (f *fakeBrowser) HasActiveSession(context.Context) (bool, error) {
	return f.active, nil
}

func TestFetchURLNavigatesAndExtracts(t *testing.T) {
	browser := &fakeBrowser{result: map[string]any{
		"title":           "Example",
		"content":         "<article>Hello</article>",
		"contentMarkdown": "Hello",
		"wordCount":       1,
	}}
	result, err := New(browser).Fetch(context.Background(), "https://example.com", &Options{Markdown: true})
	if err != nil {
		t.Fatal(err)
	}
	if browser.opened != "https://example.com" {
		t.Fatalf("opened %q", browser.opened)
	}
	if result.Markdown != "Hello" || result.Content == "" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !strings.Contains(browser.script, "globalThis.Defuddle") {
		t.Fatal("Defuddle browser bundle was not evaluated")
	}
}

func TestFetchCurrentTabDoesNotLaunchSession(t *testing.T) {
	browser := &fakeBrowser{}
	_, err := New(browser).Fetch(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "no active") {
		t.Fatalf("expected no-active-session error, got %v", err)
	}
	if browser.opened != "" {
		t.Fatalf("unexpected navigation: %s", browser.opened)
	}
}

func TestFetchCurrentTabRejectsBlankPage(t *testing.T) {
	browser := &fakeBrowser{active: true, url: "about:blank"}
	_, err := New(browser).Fetch(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "no readable page") {
		t.Fatalf("expected blank-page error, got %v", err)
	}
}
