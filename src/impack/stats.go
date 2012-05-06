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
)

func unionOfRects(rects []image.Rectangle) image.Rectangle {
	if len(rects) == 0 {
		return image.ZR
	}

	union := rects[0]
	for i := 1; i < len(rects); i++ {
		union = union.Union(rects[i])
	}
	
	return union
}

func totalAreaOfRects(rects []image.Rectangle) int {
	a := 0
	for i := 0; i < len(rects); i++ {
		a += area(rects[i])
	}
	
	return a
}

func wastedAreaOfRects(rects []image.Rectangle) int {
	return area(unionOfRects(rects)) - totalAreaOfRects(rects)
}

func Stats(minSize, maxSize, iterationsCount int) (avgFillRate float64) {	
	rects := make([]image.Rectangle, 100)
		
	for i := 0; i < iterationsCount; i++ {
		
		//Initializing random rectangles.
		for j := 0 ; j < len(rects); j++ {
			rects[j] = image.Rect(0, 0, minSize + int(rand.Float64() * float64(maxSize - minSize)), minSize + int(rand.Float64() * float64(maxSize - minSize)))
		}
		
		Arrange(rects)
		
		unionArea := area(unionOfRects(rects))
		avgFillRate += float64(unionArea - wastedAreaOfRects(rects)) / float64(unionArea)
	}
	
	avgFillRate /= float64(iterationsCount)
	
	return
}