package impack

import (
	"image"
	"testing"
)

func TestArea(t *testing.T) {
	r := image.Rect(0, 0, 10, 10)
	a := area(r)

	if a != 100 {
		t.Errorf("Wrong area %d. Correct: 100.", a)
	}
}

func TestAspectRatio(t *testing.T) {
	r := image.Rect(0, 0, 120, 80)
	ar := aspectRatio(r)
	mustBe := float64(120) / 80

	if ar != mustBe {
		t.Errorf("Wrong aspect ratio %f. Correct: %f", ar, mustBe)
	}
}

func TestIntersectAny(t *testing.T) {
	rectsInt := []image.Rectangle{
		image.Rect(2, 1, 4, 3),
		image.Rect(2, 4, 4, 6),
	}
	rectsNoInt := []image.Rectangle{
		image.Rect(5, 1, 7, 4),
		image.Rect(5, 5, 7, 6),
	}
	r := image.Rect(1, 2, 3, 5)

	if !intersectsAny(r, rectsInt, len(rectsInt)) {
		t.Errorf("No intersection where must be.")
	}

	if intersectsAny(r, rectsNoInt, len(rectsNoInt)) {
		t.Errorf("Intersection where must not be.")
	}
}

func TestWastedArea(t *testing.T) {
	rects := []image.Rectangle{
		image.Rect(5, 1, 7, 4),
		image.Rect(5, 5, 7, 6),
	}
	r := image.Rect(1, 2, 3, 5)
	w := wastedArea(r, rects, len(rects))
	mustBe := 16
		
	if w != mustBe {
		t.Errorf("Wrong wasted area %d. Correct: %d", w, mustBe)
	}
}
