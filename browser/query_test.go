package browser

import "testing"

// TestBoundingBox verifies that the BoundingBox struct is exported and usable.
func TestBoundingBox(t *testing.T) {
	box := &BoundingBox{X: 10, Y: 20, Width: 100, Height: 50}
	if box.X != 10 || box.Y != 20 || box.Width != 100 || box.Height != 50 {
		t.Fatalf("unexpected BoundingBox values: %+v", box)
	}
}
