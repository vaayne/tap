package browser

import "testing"

func TestTargetInfoFields(t *testing.T) {
	ti := TargetInfo{
		TargetID: "ABC123",
		Title:    "Example Page",
		URL:      "https://example.com",
		Type:     "page",
	}

	if ti.TargetID != "ABC123" {
		t.Fatalf("TargetID = %q, want ABC123", ti.TargetID)
	}
	if ti.Title != "Example Page" {
		t.Fatalf("Title = %q, want Example Page", ti.Title)
	}
	if ti.URL != "https://example.com" {
		t.Fatalf("URL = %q, want https://example.com", ti.URL)
	}
	if ti.Type != "page" {
		t.Fatalf("Type = %q, want page", ti.Type)
	}
}
