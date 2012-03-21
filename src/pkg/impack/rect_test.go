package impack

import (
	"image"
	"testing"
	"rand"
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
	rectsInt := []*image.Rectangle{
		&image.Rectangle{image.Point{2, 1}, image.Point{4, 3}},
		&image.Rectangle{image.Point{2, 4}, image.Point{4, 6}},
	}
	rectsNoInt := []*image.Rectangle{
		&image.Rectangle{image.Point{5, 1}, image.Point{7, 4}},
		&image.Rectangle{image.Point{5, 5}, image.Point{7, 6}},
	}
	r := image.Rect(1, 2, 3, 5)

	if !intersectsAny(r, rectsInt, len(rectsInt)) {
		t.Errorf("No intersection where must be.")
	}

	if intersectsAny(r, rectsNoInt, len(rectsNoInt)) {
		t.Errorf("Intersection where must not be.")
	}
}

func BenchmarkArrange(t *testing.B) {
	rects := make([]image.Rectangle, 100)
	
	minSize := 1
	maxSize := 100
	
	//Initializing random rectangles.
	for i := 0; i < len(rects); i++ {
		rects[i] = image.Rect(0, 0, minSize+int(rand.Float64()*float64(maxSize-minSize)), minSize+int(rand.Float64()*float64(maxSize-minSize)))
	}

	rectsWorking := make([]image.Rectangle, 100)

	for i := 0; i < 100; i++ {
		copy(rectsWorking, rects)
		Arrange(rectsWorking)
	}
}
