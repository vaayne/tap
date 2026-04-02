package browser

import "testing"

func TestTargetTypePageConstant(t *testing.T) {
	if TargetTypePage != "page" {
		t.Fatalf("TargetTypePage = %q, want %q", TargetTypePage, "page")
	}
}
