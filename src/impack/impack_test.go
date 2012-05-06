// ImPack - CSS sprites maker
// Copyright (C) 2012 Dmitry Bratus
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package impack

import (
	"image"
	"math/rand"
	"strconv"
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

func TestMakePlacements(t *testing.T) {
	pivot := image.Rect(3, 3, 6, 6)
	rect := image.Rect(0, 0, 2, 2)
	out := make([]image.Rectangle, 8)

	makePlacements(pivot, rect, out)

	for i, r := range out {
		t.Logf("%d: %s\n", i, r.String())

		if r.Dx() != 2 || r.Dy() != 2 {
			t.Errorf("Size of a placement is invalid.")
			t.Fail()
		}

		if !pivot.Intersect(r).Empty() {
			t.Errorf("Placement intersects the pivot.")
			t.Fail()
		}
	}
}

func TestArrange(t *testing.T) {
	rects := []image.Rectangle{
		image.Rect(0, 0, 4, 12),
		image.Rect(0, 0, 5, 15),
		image.Rect(0, 0, 2, 6),
		image.Rect(0, 0, 3, 9),
		image.Rect(0, 0, 1, 3),
	}
	areas := []int{4 * 12, 5 * 15, 2 * 6, 3 * 9, 1 * 3}

	union := Arrange(rects)

	str := ""

	for i := 0; i < union.Max.X; i++ {
		for j := 0; j < union.Max.Y; j++ {
			found := false

			for k := 0; k < len(rects); k++ {
				if !image.Rect(i, j, i+1, j+1).Intersect(rects[k]).Empty() {
					str += strconv.Itoa(k)
					found = true
					break
				}
			}

			if !found {
				str += "X"
			}
		}

		str += "\n"
	}

	t.Log(str)

	for i, r := range rects {
		if area(r) != areas[i] {
			t.Errorf("Size of a placement is invalid.")
			t.Fail()
		}
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
