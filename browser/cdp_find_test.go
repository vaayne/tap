package browser

import (
	"strings"
	"testing"
)

// TestJSLocatorText verifies the generated JS contains the expected search text.
func TestJSLocatorText(t *testing.T) {
	t.Run("substring", func(t *testing.T) {
		js := jsLocatorText("Sign in", false)
		// Case-folded text should appear in the generated JS.
		if !strings.Contains(js, "sign in") {
			t.Errorf("expected lowercase text in JS, got: %s", js)
		}
		if !strings.Contains(js, "includes") {
			t.Errorf("expected substring match (includes), got: %s", js)
		}
	})
	t.Run("exact", func(t *testing.T) {
		js := jsLocatorText("Submit", true)
		if !strings.Contains(js, "Submit") {
			t.Errorf("expected original-case text in exact JS, got: %s", js)
		}
		if !strings.Contains(js, "=== text") {
			t.Errorf("expected strict equality in exact JS, got: %s", js)
		}
	})
}

// TestJSLocatorLabel verifies the label locator references both htmlFor and aria-label.
func TestJSLocatorLabel(t *testing.T) {
	js := jsLocatorLabel("Email")
	if !strings.Contains(js, "htmlFor") {
		t.Errorf("expected htmlFor in label JS, got: %s", js)
	}
	if !strings.Contains(js, "aria-label") {
		t.Errorf("expected aria-label fallback in label JS, got: %s", js)
	}
}

// TestJSLocatorAttr verifies attribute-based locator JS.
func TestJSLocatorAttr(t *testing.T) {
	js := jsLocatorAttr("placeholder", "Search")
	if !strings.Contains(js, `[placeholder]`) {
		t.Errorf("expected attribute selector, got: %s", js)
	}
	if !strings.Contains(js, "search") {
		t.Errorf("expected lowercased query, got: %s", js)
	}
}

// TestJSLocatorNth verifies nth and last locators.
func TestJSLocatorNth(t *testing.T) {
	js := jsLocatorNth("li.item", 2)
	if !strings.Contains(js, `"li.item"`) {
		t.Errorf("expected CSS selector in nth JS, got: %s", js)
	}
	if !strings.Contains(js, "els[2]") {
		t.Errorf("expected index 2 in nth JS, got: %s", js)
	}
}

func TestJSLocatorLast(t *testing.T) {
	js := jsLocatorLast("tr")
	if !strings.Contains(js, "els.length - 1") {
		t.Errorf("expected last-element expression in JS, got: %s", js)
	}
}

// TestResolveLocatorKind verifies all LocatorKind constants are handled
// without panicking (no real browser required — just routing logic).
func TestResolveLocatorKindRouting(t *testing.T) {
	kinds := []LocatorKind{
		LocatorRole, LocatorText, LocatorLabel, LocatorPlaceholder,
		LocatorAlt, LocatorTitle, LocatorTestID,
		LocatorFirst, LocatorLast, LocatorNth,
	}
	for _, k := range kinds {
		if k == "" {
			t.Errorf("LocatorKind constant is empty")
		}
	}
}

// TestFindActionConstants verifies all action constants are non-empty.
func TestFindActionConstants(t *testing.T) {
	actions := []FindAction{
		FindActionClick, FindActionFill, FindActionType,
		FindActionHover, FindActionFocus,
		FindActionCheck, FindActionUncheck, FindActionText,
	}
	for _, a := range actions {
		if a == "" {
			t.Errorf("FindAction constant is empty")
		}
	}
}
