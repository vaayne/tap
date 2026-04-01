package fetch

import (
	"context"
	"testing"

	"github.com/vaayne/tap/transport"
)

func TestFetch_ParseHTML(t *testing.T) {
	tp, err := transport.New(context.Background(), transport.Config{})
	if err != nil {
		t.Fatalf("New transport: %v", err)
	}
	f, err := New(tp)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer f.Close()

	// Test with a real URL
	result, err := f.Fetch(context.Background(), "https://example.com", &Options{Markdown: true})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if result.Title == "" && result.Content == "" && result.Markdown == "" {
		t.Error("expected some content from example.com")
	}
}

func TestFetch_InvalidURL(t *testing.T) {
	tp, err := transport.New(context.Background(), transport.Config{})
	if err != nil {
		t.Fatalf("New transport: %v", err)
	}
	f, err := New(tp)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer f.Close()

	_, err = f.Fetch(context.Background(), "://bad-url", nil)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
